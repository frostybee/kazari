package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/collapsible"
	"github.com/frostybee/kazari/internal/config"
)

func buildCollapseRangeMap(ranges []config.CollapseRange) map[int]int {
	if len(ranges) == 0 {
		return nil
	}
	m := make(map[int]int)
	for i, cr := range ranges {
		for line := cr.Start; line <= cr.End; line++ {
			m[line] = i
		}
	}
	return m
}

func (lctx *lineContext) inCollapseRange(lineNum int) (int, bool) {
	if lctx.collapseRangeMap == nil {
		return 0, false
	}
	idx, ok := lctx.collapseRangeMap[lineNum]
	return idx, ok
}

func (lctx *lineContext) hasVisibleLineAfter(lineNum int) bool {
	for _, seg := range lctx.resolved.CollapseSegments {
		if seg.Start > lineNum {
			return true
		}
	}
	return false
}

func (lctx *lineContext) inCollapseRangeEnd(lineNum int) bool {
	idx, ok := lctx.collapseRangeMap[lineNum]
	if !ok {
		return false
	}
	return lineNum == lctx.resolved.CollapseRanges[idx].End
}

func buildThresholdVisibleSet(segments []config.PreviewSegment) map[int]bool {
	if len(segments) == 0 {
		return nil
	}
	m := make(map[int]bool)
	for _, seg := range segments {
		for line := seg.Start; line <= seg.End; line++ {
			m[line] = true
		}
	}
	return m
}

func renderSummaryLine(sb *strings.Builder, resolved *config.ResolvedBlock, cr config.CollapseRange, cfg *config.Config) {
	sb.WriteString("<summary>")
	sb.WriteString("<div class=\"kz-line\">")
	if resolved.LineNumbers {
		sb.WriteString("<div class=\"kz-gutter\"><div class=\"kz-ln\"></div></div>")
	}

	indentStyle := ""
	if cr.MinIndent > 0 && resolved.CollapseConfig != nil && resolved.CollapseConfig.PreserveIndent {
		indentStyle = fmt.Sprintf(" style=\"--kz-indent:%dch\"", cr.MinIndent)
	}

	sb.WriteString(fmt.Sprintf("<div class=\"kz-code\"%s>", indentStyle))
	sb.WriteString("<span class=\"expand\" aria-hidden=\"true\"></span>")
	sb.WriteString("<span class=\"collapse\" aria-hidden=\"true\"></span>")
	sb.WriteString(fmt.Sprintf("<span class=\"text\">%s</span>", collapsible.SummaryText(cr.LineCount, cfg.UIStrings)))
	sb.WriteString("</div>")
	sb.WriteString("</div>")
	sb.WriteString("</summary>")
}

func renderCollapseRangeOpen(sb *strings.Builder, resolved *config.ResolvedBlock, cr config.CollapseRange, cfg *config.Config) {
	switch cr.Style {
	case config.CollapseCollapsibleStart:
		sb.WriteString("<div class=\"kz-section collapsible-start\">")
		sb.WriteString("<details>")
		renderSummaryLine(sb, resolved, cr, cfg)
		sb.WriteString("</details>")
		sb.WriteString("<div class=\"content-lines\">")

	case config.CollapseCollapsibleEnd:
		sb.WriteString("<div class=\"kz-section collapsible-end\">")
		sb.WriteString("<details>")
		renderSummaryLine(sb, resolved, cr, cfg)
		sb.WriteString("</details>")
		sb.WriteString("<div class=\"content-lines\">")

	default:
		sb.WriteString("<details class=\"kz-section\">")
		renderSummaryLine(sb, resolved, cr, cfg)
	}
}

func renderCollapseRangeClose(sb *strings.Builder, cr config.CollapseRange) {
	switch cr.Style {
	case config.CollapseCollapsibleStart, config.CollapseCollapsibleEnd:
		sb.WriteString("</div></div>")
	default:
		sb.WriteString("</details>")
	}
}

func renderGapIndicator(sb *strings.Builder, resolved *config.ResolvedBlock) {
	sb.WriteString("<div class=\"kz-line kz-gap\">")
	if resolved.LineNumbers {
		sb.WriteString("<div class=\"kz-gutter\"><div class=\"kz-ln\"></div></div>")
	}
	sb.WriteString("<div class=\"kz-code\"><span class=\"kz-gap-indicator\" aria-hidden=\"true\">⋮</span><span class=\"sr-only\">Lines hidden</span></div>")
	sb.WriteString("</div>")
}

func renderHiddenLine(sb *strings.Builder, line TokenLine, lineNum int, lctx *lineContext) {
	classes := "kz-line kz-hidden"
	extraClasses, markerType, _ := resolveMarkerClasses(lineNum, lctx)
	classes += extraClasses

	sb.WriteString(fmt.Sprintf("<div class=\"%s\">", classes))
	if lctx.resolved.LineNumbers {
		sb.WriteString(fmt.Sprintf("<div class=\"kz-gutter\"><div class=\"kz-ln\" aria-hidden=\"true\">%d</div></div>", lineNum))
	}
	sb.WriteString("<div class=\"kz-code\">")
	for _, tok := range line.Tokens {
		if tok.Content == "" {
			continue
		}
		renderToken(sb, tok, lctx, markerType)
	}
	sb.WriteString("</div></div>")
}

func renderCollapseContentStart(sb *strings.Builder, resolved *config.ResolvedBlock) {
	if resolved.CollapseThreshold {
		sb.WriteString("<div class=\"kz-collapse-content\">")
	}
}

func renderCollapseContentEnd(sb *strings.Builder, resolved *config.ResolvedBlock) {
	if resolved.CollapseThreshold {
		sb.WriteString("<div class=\"kz-collapse-gradient\"></div>")
		sb.WriteString("</div>")
	}
}

func renderCollapseBar(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	expandText := cfg.UIStrings.ExpandButtonText
	collapseText := cfg.UIStrings.CollapseButtonText
	expandedAnnouncement := cfg.UIStrings.ExpandedAnnouncement
	collapsedAnnouncement := cfg.UIStrings.CollapsedAnnouncement
	if resolved.CollapseConfig != nil {
		if resolved.CollapseConfig.ExpandButtonText != "" {
			expandText = resolved.CollapseConfig.ExpandButtonText
		}
		if resolved.CollapseConfig.CollapseButtonText != "" {
			collapseText = resolved.CollapseConfig.CollapseButtonText
		}
		if resolved.CollapseConfig.ExpandedAnnouncement != "" {
			expandedAnnouncement = resolved.CollapseConfig.ExpandedAnnouncement
		}
		if resolved.CollapseConfig.CollapsedAnnouncement != "" {
			collapsedAnnouncement = resolved.CollapseConfig.CollapsedAnnouncement
		}
	}

	if resolved.CollapseBeyondCap > 0 {
		expandText = fmt.Sprintf("%s (+%d highlighted)", expandText, resolved.CollapseBeyondCap)
	}

	sb.WriteString("<div class=\"kz-collapse-bar\">")
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-collapse-btn\" aria-expanded=\"false\" data-expand=\"%s\" data-collapse=\"%s\" data-expanded-msg=\"%s\" data-collapsed-msg=\"%s\">%s</button>",
		html.EscapeString(expandText),
		html.EscapeString(collapseText),
		html.EscapeString(expandedAnnouncement),
		html.EscapeString(collapsedAnnouncement),
		html.EscapeString(expandText),
	))
	sb.WriteString("</div>")
	sb.WriteString("<div class=\"kz-sr-announce\" aria-live=\"polite\"></div>")
}
