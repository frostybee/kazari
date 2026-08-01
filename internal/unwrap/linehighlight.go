package unwrap

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// collectHighlightedLines walks a token window counting line wrapper spans
// and returns the 1-based numbers of lines whose wrapper also carries the hl
// class. The counter gates on the class token "line" because Chroma always
// pairs hl with line on the same span, which keeps counting correct whether
// or not any line is highlighted. The walk is deliberately independent of
// recoverText's gutter exclusion: line number spans use ln or lnt while line
// wrappers use line and hl, and conflating the two walks is an easy way to
// miscount. Chroma writes these classes only in classes mode; Hugo's default
// inline styles mode carries no classes, so no highlights are recoverable
// there and callers simply get an empty slice.
func collectHighlightedLines(tokens []BufferedToken) []int {
	line := 0
	var highlighted []int
	for _, tok := range tokens {
		if tok.Type != html.StartTagToken {
			continue
		}
		if !hasClass(tok.Attrs, "line") {
			continue
		}
		line++
		if hasClass(tok.Attrs, "hl") {
			highlighted = append(highlighted, line)
		}
	}
	return highlighted
}

// collapseRanges renders ascending line numbers as the brace range syntax
// the meta grammar parses, collapsing consecutive runs: 3,4,5,9 becomes
// "3-5,9". An empty input yields an empty string.
func collapseRanges(lines []int) string {
	if len(lines) == 0 {
		return ""
	}
	var parts []string
	start, prev := lines[0], lines[0]
	flush := func() {
		if start == prev {
			parts = append(parts, strconv.Itoa(start))
			return
		}
		parts = append(parts, strconv.Itoa(start)+"-"+strconv.Itoa(prev))
	}
	for _, n := range lines[1:] {
		if n == prev+1 {
			prev = n
			continue
		}
		flush()
		start, prev = n, n
	}
	flush()
	return strings.Join(parts, ",")
}

// buildMeta assembles the full meta string: the language token first, then a
// bare brace range token when any lines are highlighted. The bare braces
// parse as a MarkerMark line marker in internal/meta.
func buildMeta(lang string, highlighted []int) string {
	ranges := collapseRanges(highlighted)
	switch {
	case lang == "" && ranges == "":
		return ""
	case ranges == "":
		return lang
	case lang == "":
		return "{" + ranges + "}"
	default:
		return lang + " {" + ranges + "}"
	}
}
