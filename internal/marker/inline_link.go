package marker

import "github.com/frostybee/kazari/internal/config"

// ProcessLinks splits tokens according to link annotations, producing segments
// with the Link field set. This reuses the same splitTokens machinery as inline markers.
func ProcessLinks(tokens []config.MergedToken, links []config.LinkAnnotation) []TokenWithSegments {
	if len(links) == 0 {
		return wrapTokens(tokens)
	}

	matches := make([]inlineMatch, len(links))
	for i, l := range links {
		matches[i] = inlineMatch{
			start: l.Start,
			end:   l.End,
			mtype: config.MarkerNone,
			link:  l.URL,
		}
	}

	return splitTokens(tokens, matches)
}

// ProcessInlineMarkersAndLinks handles both inline markers and link annotations
// on the same line. Inline markers are processed first, then link annotations
// are applied on top by further splitting segments.
func ProcessInlineMarkersAndLinks(tokens []config.MergedToken, markers []config.InlineMarker, links []config.LinkAnnotation) []TokenWithSegments {
	result := ProcessInlineMarkers(tokens, markers)
	if len(links) == 0 {
		return result
	}
	return applyLinksToSegmented(result, links)
}

func applyLinksToSegmented(tokenSegs []TokenWithSegments, links []config.LinkAnnotation) []TokenWithSegments {
	out := make([]TokenWithSegments, 0, len(tokenSegs))
	cursor := 0
	li := 0

	for _, ts := range tokenSegs {
		var newSegs []Segment
		for _, seg := range ts.Segments {
			segStart := cursor
			segEnd := cursor + len(seg.Content)

			var parts []Segment
			pos := segStart

			for li < len(links) && pos < segEnd {
				l := links[li]
				if l.Start >= segEnd {
					break
				}
				if l.End <= pos {
					li++
					continue
				}

				lStart := l.Start
				if lStart < pos {
					lStart = pos
				}
				lEnd := l.End
				if lEnd > segEnd {
					lEnd = segEnd
				}

				if lStart > pos {
					parts = append(parts, Segment{
						Content: seg.Content[pos-segStart : lStart-segStart],
						Marker:  seg.Marker,
					})
				}

				linkMarker := &InlineAnnotation{
					Type:      config.MarkerNone,
					OpenStart: l.Start < segStart,
					OpenEnd:   l.End > segEnd,
					Link:      l.URL,
				}
				if seg.Marker != nil {
					linkMarker.Type = seg.Marker.Type
					linkMarker.OpenStart = linkMarker.OpenStart || seg.Marker.OpenStart
					linkMarker.OpenEnd = linkMarker.OpenEnd || seg.Marker.OpenEnd
				}

				parts = append(parts, Segment{
					Content: seg.Content[lStart-segStart : lEnd-segStart],
					Marker:  linkMarker,
				})

				pos = lEnd
				if l.End <= segEnd {
					li++
				} else {
					break
				}
			}

			if pos < segEnd {
				parts = append(parts, Segment{
					Content: seg.Content[pos-segStart:],
					Marker:  seg.Marker,
				})
			}

			if len(parts) == 0 {
				parts = []Segment{seg}
			}

			newSegs = append(newSegs, parts...)
			cursor = segEnd
		}

		out = append(out, TokenWithSegments{Token: ts.Token, Segments: newSegs})
	}

	return out
}
