package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// gutterClasses are class tokens that mark line number gutters across the
// supported generator families. Text inside an element carrying any of these
// is display chrome, not source code, and must be excluded from recovery.
// Sources: Chroma emits ln (inline mode), lnt and lntd (table mode); Rouge
// emits gutter, gl, lineno; Pygments emits linenos and linenodiv; the Prism
// line numbers plugin emits line-numbers and line-numbers-rows.
var gutterClasses = map[string]bool{
	"ln":                true,
	"lnt":               true,
	"lntd":              true,
	"gutter":            true,
	"gl":                true,
	"linenos":           true,
	"lineno":            true,
	"linenodiv":         true,
	"line-numbers":      true,
	"line-numbers-rows": true,
}

// voidElements never receive end tags and must not affect depth tracking.
var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

// isGutterElement reports whether the attributes carry a gutter class.
func isGutterElement(attrs []html.Attribute) bool {
	for _, c := range classList(attrs) {
		if gutterClasses[c] {
			return true
		}
	}
	return false
}

// recoverText concatenates the text content of a token window, excluding
// every subtree rooted at a gutter element. Text token Data arrives already
// entity unescaped from the tokenizer, so no second unescape ever runs. A br
// element becomes a newline, guarding against legacy generator output.
// Comments are ignored. Callers apply strings.TrimRight(result, "\n") once,
// mirroring the Goldmark path in goldmark/renderer.go.
func recoverText(tokens []BufferedToken) string {
	var b strings.Builder
	depth := 0
	skipAbove := -1 // when >= 0, suppress output until depth returns to it
	for _, tok := range tokens {
		switch tok.Type {
		case html.StartTagToken:
			if tok.Data == "br" {
				if skipAbove < 0 {
					b.WriteString("\n")
				}
				continue
			}
			if voidElements[tok.Data] {
				continue
			}
			depth++
			if skipAbove < 0 && isGutterElement(tok.Attrs) {
				skipAbove = depth - 1
			}
		case html.SelfClosingTagToken:
			if tok.Data == "br" && skipAbove < 0 {
				b.WriteString("\n")
			}
		case html.EndTagToken:
			if voidElements[tok.Data] {
				continue
			}
			if depth > 0 {
				depth--
			}
			if skipAbove >= 0 && depth <= skipAbove {
				skipAbove = -1
			}
		case html.TextToken:
			if skipAbove < 0 {
				b.WriteString(tok.Data)
			}
		}
	}
	return b.String()
}
