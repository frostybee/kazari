package marker

import (
	"regexp"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// InlineAnnotation describes how a segment of text is marked inline.
type InlineAnnotation struct {
	Type      config.MarkerType
	OpenStart bool // match started in a previous token
	OpenEnd   bool // match continues into the next token
}

// Segment is a piece of a token's content, optionally wrapped in an inline marker.
type Segment struct {
	Content string
	Marker  *InlineAnnotation // nil = plain text
}

// TokenWithSegments pairs an original token with its segmented content.
type TokenWithSegments struct {
	Token    config.MergedToken
	Segments []Segment
}

type inlineMatch struct {
	start    int // character offset in plain text (inclusive)
	end      int // character offset in plain text (exclusive)
	mtype    config.MarkerType
	priority int
}

// ProcessInlineMarkers finds all inline marker matches in a line of tokens
// and returns tokens split into segments with marker annotations.
func ProcessInlineMarkers(tokens []config.MergedToken, markers []config.InlineMarker) []TokenWithSegments {
	if len(markers) == 0 {
		return wrapTokens(tokens)
	}

	plainText := buildPlainText(tokens)
	if plainText == "" {
		return wrapTokens(tokens)
	}

	matches := findAllMatches(plainText, markers)
	if len(matches) == 0 {
		return wrapTokens(tokens)
	}

	matches = resolveInlineOverlaps(matches)
	return splitTokens(tokens, matches)
}

func buildPlainText(tokens []config.MergedToken) string {
	var sb strings.Builder
	for _, t := range tokens {
		sb.WriteString(t.Content)
	}
	return sb.String()
}

func findAllMatches(text string, markers []config.InlineMarker) []inlineMatch {
	var matches []inlineMatch
	for _, m := range markers {
		if m.Text == "" {
			continue
		}
		if m.IsRegex {
			matches = append(matches, findRegexMatches(text, m)...)
		} else {
			matches = append(matches, findLiteralMatches(text, m)...)
		}
	}
	return matches
}

func findLiteralMatches(text string, m config.InlineMarker) []inlineMatch {
	var matches []inlineMatch
	offset := 0
	for {
		idx := strings.Index(text[offset:], m.Text)
		if idx < 0 {
			break
		}
		start := offset + idx
		end := start + len(m.Text)
		matches = append(matches, inlineMatch{
			start:    start,
			end:      end,
			mtype:    m.Type,
			priority: int(m.Type),
		})
		offset = end
	}
	return matches
}

func findRegexMatches(text string, m config.InlineMarker) []inlineMatch {
	re, err := regexp.Compile(m.Text)
	if err != nil {
		return nil
	}
	var matches []inlineMatch
	for _, loc := range re.FindAllStringSubmatchIndex(text, -1) {
		start, end := loc[0], loc[1]
		if len(loc) >= 4 && loc[2] >= 0 {
			start, end = loc[2], loc[3]
		}
		matches = append(matches, inlineMatch{
			start:    start,
			end:      end,
			mtype:    m.Type,
			priority: int(m.Type),
		})
	}
	return matches
}

// resolveInlineOverlaps handles overlapping matches using priority (mark < del < ins).
// Matches are claimed in priority order. When a lower priority match overlaps
// ranges already claimed by higher priority matches, only the overlapping parts
// are removed and the surviving fragments are kept.
func resolveInlineOverlaps(matches []inlineMatch) []inlineMatch {
	if len(matches) <= 1 {
		return matches
	}

	ordered := make([]inlineMatch, len(matches))
	copy(ordered, matches)
	sortByPriority(ordered)

	var resolved []inlineMatch
	for _, m := range ordered {
		fragments := []inlineMatch{m}
		for _, r := range resolved {
			fragments = subtractClaimedRange(fragments, r.start, r.end)
			if len(fragments) == 0 {
				break
			}
		}
		resolved = append(resolved, fragments...)
	}

	sortMatches(resolved)
	return resolved
}

// subtractClaimedRange removes the interval from start (inclusive) to end
// (exclusive) out of each fragment, keeping the surviving pieces.
func subtractClaimedRange(fragments []inlineMatch, start, end int) []inlineMatch {
	var out []inlineMatch
	for _, f := range fragments {
		if f.end <= start || f.start >= end {
			out = append(out, f)
			continue
		}
		if f.start < start {
			out = append(out, inlineMatch{
				start: f.start, end: start,
				mtype: f.mtype, priority: f.priority,
			})
		}
		if f.end > end {
			out = append(out, inlineMatch{
				start: end, end: f.end,
				mtype: f.mtype, priority: f.priority,
			})
		}
	}
	return out
}

// sortByPriority orders matches by priority descending, then by start ascending,
// so higher priority matches claim their ranges first.
func sortByPriority(matches []inlineMatch) {
	for i := 1; i < len(matches); i++ {
		key := matches[i]
		j := i - 1
		for j >= 0 && (matches[j].priority < key.priority || (matches[j].priority == key.priority && matches[j].start > key.start)) {
			matches[j+1] = matches[j]
			j--
		}
		matches[j+1] = key
	}
}

func sortMatches(matches []inlineMatch) {
	// Simple insertion sort — match counts are small.
	for i := 1; i < len(matches); i++ {
		key := matches[i]
		j := i - 1
		for j >= 0 && (matches[j].start > key.start || (matches[j].start == key.start && matches[j].priority < key.priority)) {
			matches[j+1] = matches[j]
			j--
		}
		matches[j+1] = key
	}
}

func splitTokens(tokens []config.MergedToken, matches []inlineMatch) []TokenWithSegments {
	result := make([]TokenWithSegments, 0, len(tokens))
	cursor := 0 // character position in plain text
	mi := 0     // match index

	for _, tok := range tokens {
		tokStart := cursor
		tokEnd := cursor + len(tok.Content)
		cursor = tokEnd

		if tok.Content == "" {
			result = append(result, TokenWithSegments{
				Token:    tok,
				Segments: []Segment{{Content: ""}},
			})
			continue
		}

		var segments []Segment
		pos := tokStart // current position within the plain text

		for mi < len(matches) && pos < tokEnd {
			m := matches[mi]

			if m.start >= tokEnd {
				break // match is beyond this token
			}

			if m.end <= pos {
				mi++ // match already passed
				continue
			}

			// Clamp match to token boundaries.
			segStart := m.start
			if segStart < pos {
				segStart = pos
			}
			segEnd := m.end
			if segEnd > tokEnd {
				segEnd = tokEnd
			}

			// Emit unmarked text before the match.
			if segStart > pos {
				segments = append(segments, Segment{
					Content: tok.Content[pos-tokStart : segStart-tokStart],
				})
			}

			// Emit marked segment.
			segments = append(segments, Segment{
				Content: tok.Content[segStart-tokStart : segEnd-tokStart],
				Marker: &InlineAnnotation{
					Type:      m.mtype,
					OpenStart: m.start < tokStart, // match started before this token
					OpenEnd:   m.end > tokEnd,      // match continues past this token
				},
			})

			pos = segEnd
			if m.end <= tokEnd {
				mi++ // match fully consumed within this token
			} else {
				break // match continues into next token
			}
		}

		// Emit remaining unmarked text.
		if pos < tokEnd {
			segments = append(segments, Segment{
				Content: tok.Content[pos-tokStart:],
			})
		}

		if len(segments) == 0 {
			segments = []Segment{{Content: tok.Content}}
		}

		result = append(result, TokenWithSegments{Token: tok, Segments: segments})
	}

	return result
}

func wrapTokens(tokens []config.MergedToken) []TokenWithSegments {
	result := make([]TokenWithSegments, len(tokens))
	for i, tok := range tokens {
		result[i] = TokenWithSegments{
			Token:    tok,
			Segments: []Segment{{Content: tok.Content}},
		}
	}
	return result
}
