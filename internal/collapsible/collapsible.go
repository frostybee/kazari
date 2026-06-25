package collapsible

import (
	"sort"
	"strings"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/locale"
)

// CollapseResult holds the resolved collapse state for a single code block.
type CollapseResult struct {
	Threshold       bool
	PreviewSegments []config.PreviewSegment
	BeyondCapCount  int // marked lines beyond 2× cap
	Ranges          []config.CollapseRange
}

// ResolveCollapse determines the collapse state for a code block.
func ResolveCollapse(
	lineCount int,
	spec *config.CollapseSpec,
	cfg *config.CollapsibleConfig,
	code string,
	markers []config.LineMarker,
	focusLines []config.LineRange,
) CollapseResult {
	result := CollapseResult{}

	// Threshold-based collapse
	if shouldThresholdCollapse(lineCount, spec, cfg) {
		result.Threshold = true
		preview := cfg.PreviewLines
		if preview <= 0 {
			preview = 8
		}
		segments, beyondCap := computePreviewSegments(preview, lineCount, markers, focusLines)
		result.PreviewSegments = segments
		result.BeyondCapCount = beyondCap
	}

	// Range-based collapse
	if spec != nil && len(spec.Ranges) > 0 {
		validated := validateRanges(spec.Ranges, lineCount)

		// Resolve style: per-block meta > engine config > github default
		baseStyle := config.CollapseGithub
		if cfg != nil && cfg.Style != config.CollapseGithub {
			baseStyle = cfg.Style
		}
		if spec.Style != nil {
			baseStyle = *spec.Style
		}

		for _, r := range validated {
			indent := computeMinIndent(code, r.Start, r.End)
			style := resolveCollapseStyle(baseStyle, r.End, lineCount)
			result.Ranges = append(result.Ranges, config.CollapseRange{
				Start:     r.Start,
				End:       r.End,
				LineCount: r.End - r.Start + 1,
				MinIndent: indent,
				Style:     style,
			})
		}
	}

	return result
}

// ShouldThresholdCollapse determines if threshold-based collapse applies.
func shouldThresholdCollapse(lineCount int, spec *config.CollapseSpec, cfg *config.CollapsibleConfig) bool {
	if cfg == nil {
		return false
	}
	if spec != nil && spec.Disabled {
		return false
	}
	if spec != nil && spec.Enabled {
		return true
	}
	threshold := cfg.LineThreshold
	if spec != nil && spec.Threshold != nil && *spec.Threshold > 0 {
		threshold = *spec.Threshold
	}
	if threshold <= 0 {
		threshold = 15
	}
	return lineCount > threshold
}

// ValidateRanges filters and normalizes collapse ranges.
// Drops reversed ranges, out-of-bounds ranges, and overlapping ranges (first wins).
// Returns sorted by Start.
func validateRanges(ranges []config.LineRange, lineCount int) []config.LineRange {
	var valid []config.LineRange
	for _, r := range ranges {
		if r.Start > r.End {
			continue
		}
		if r.Start < 1 || r.End < 1 {
			continue
		}
		if r.Start > lineCount {
			continue
		}
		// Clamp end to line count
		end := r.End
		if end > lineCount {
			end = lineCount
		}
		valid = append(valid, config.LineRange{Start: r.Start, End: end})
	}

	sort.Slice(valid, func(i, j int) bool {
		return valid[i].Start < valid[j].Start
	})

	// Remove overlapping ranges (first valid wins)
	var result []config.LineRange
	for _, r := range valid {
		if len(result) > 0 && r.Start <= result[len(result)-1].End {
			continue
		}
		result = append(result, r)
	}
	return result
}

// ComputePreviewSegments computes non-contiguous visible segments for threshold preview.
// Returns the segments to show and the count of marked lines beyond the 2× cap (for badge).
func computePreviewSegments(previewLines int, lineCount int, markers []config.LineMarker, focusLines []config.LineRange) ([]config.PreviewSegment, int) {
	base := previewLines
	if base >= lineCount {
		return []config.PreviewSegment{{Start: 1, End: lineCount}}, 0
	}

	maxCap := base * 2
	if maxCap > lineCount {
		maxCap = lineCount
	}

	markedSet := buildMarkedSet(markers, focusLines)

	if len(markedSet) == 0 {
		return []config.PreviewSegment{{Start: 1, End: base}}, 0
	}

	// Partition marked lines: within cap vs beyond cap
	var withinCap []int
	beyondCap := 0
	for line := range markedSet {
		if line <= base {
			continue // already in base preview
		}
		if line <= maxCap {
			withinCap = append(withinCap, line)
		} else {
			beyondCap++
		}
	}

	if len(withinCap) == 0 {
		return []config.PreviewSegment{{Start: 1, End: base}}, beyondCap
	}

	sort.Ints(withinCap)

	// Build segments: base preview + marked line clusters (±1 context)
	segments := []config.PreviewSegment{{Start: 1, End: base}}
	for _, line := range withinCap {
		start := line - 1
		if start < 1 {
			start = 1
		}
		end := line + 1
		if end > lineCount {
			end = lineCount
		}
		seg := config.PreviewSegment{Start: start, End: end}

		// Merge with previous segment if overlapping or adjacent
		last := &segments[len(segments)-1]
		if seg.Start <= last.End+1 {
			if seg.End > last.End {
				last.End = seg.End
			}
		} else {
			segments = append(segments, seg)
		}
	}

	return segments, beyondCap
}

func buildMarkedSet(markers []config.LineMarker, focusLines []config.LineRange) map[int]bool {
	markedSet := make(map[int]bool)
	for _, m := range markers {
		for _, lr := range m.Lines {
			for line := lr.Start; line <= lr.End; line++ {
				markedSet[line] = true
			}
		}
	}
	for _, lr := range focusLines {
		for line := lr.Start; line <= lr.End; line++ {
			markedSet[line] = true
		}
	}
	return markedSet
}

// ComputeMinIndent calculates the minimum indentation (in spaces) of lines
// in the range [startLine, endLine] (1-based), ignoring blank lines.
func computeMinIndent(code string, startLine, endLine int) int {
	lines := strings.Split(code, "\n")
	minIndent := -1

	for i := startLine - 1; i < endLine && i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(trimmed)
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent < 0 {
		return 0
	}
	return minIndent
}

// ResolveCollapseStyle resolves CollapseCollapsibleAuto to start or end
// based on whether the range ends at the last line of the block.
func resolveCollapseStyle(style config.CollapseStyle, rangeEnd, lineCount int) config.CollapseStyle {
	if style != config.CollapseCollapsibleAuto {
		return style
	}
	if rangeEnd >= lineCount {
		return config.CollapseCollapsibleEnd
	}
	return config.CollapseCollapsibleStart
}

// SummaryText returns the summary text for a collapsed section.
func SummaryText(lineCount int, strings *locale.UIStrings) string {
	return locale.FormatCollapsedLines(strings, lineCount)
}
