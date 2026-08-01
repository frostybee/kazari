package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// Skip reasons reported by SkipReason. The processor records these per block
// so check runs can account for skipped regions without re-deriving why.
const (
	SkipKzClass      = "kz-class"
	SkipKazariIgnore = "data-kazari-ignore"
	SkipMermaid      = "mermaid"
)

// SkipReason reports why a candidate region must be left untouched, or an
// empty string when the region should proceed to unwrapping. These checks
// run before any chain dispatch because an already processed site hits them
// on nearly every candidate during a check run. Kazari's own output roots at
// class kazari-block with kz- prefixed classes inside, so both prefixes mark
// a block as already processed.
func SkipReason(tokens []BufferedToken) string {
	if len(tokens) == 0 {
		return ""
	}
	for _, c := range classList(tokens[0].Attrs) {
		if strings.HasPrefix(c, "kz-") || strings.HasPrefix(c, "kazari-") {
			return SkipKzClass
		}
	}
	if v, ok := attrValue(tokens[0].Attrs, "data-kazari"); ok && v == "ignore" {
		return SkipKazariIgnore
	}
	if isMermaid(tokens) {
		return SkipMermaid
	}
	return ""
}

// isMermaid detects mermaid blocks in the shapes generators and mermaid
// integrations converge on: a mermaid class on the container, a
// language-mermaid class, or data-lang set to mermaid.
func isMermaid(tokens []BufferedToken) bool {
	for _, tok := range tokens {
		if tok.Type != html.StartTagToken {
			continue
		}
		if tok.Data != "pre" && tok.Data != "code" && tok.Data != "div" {
			continue
		}
		if hasClass(tok.Attrs, "mermaid") || hasClass(tok.Attrs, "language-mermaid") {
			return true
		}
		if v, ok := attrValue(tok.Attrs, "data-lang"); ok && v == "mermaid" {
			return true
		}
	}
	return false
}
