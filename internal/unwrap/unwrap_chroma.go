package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// chromaDivWrapper recognizes Hugo's wrapped Chroma output. Two shapes are
// accepted under the div with the highlight class token: the classes mode
// shape whose pre carries the chroma class, and Hugo's default inline styles
// mode, which carries no Chroma classes at all. In the second shape the code
// element Hugo emits still carries the language class and data-lang, which
// keeps the signature strong enough to avoid claiming unrelated wrappers.
type chromaDivWrapper struct{}

func (chromaDivWrapper) Name() string { return "chroma-div" }

func (chromaDivWrapper) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "div") || !hasClass(tokens[0].Attrs, "highlight") {
		return Region{}, false
	}

	// Table line number mode: div.chroma > table.lntable with two td.lntd
	// cells. The first cell is the gutter; its pre has no code element and
	// holds only digits. Recovery must run inside the second cell only.
	if ti := findElement(tokens, 1, "table", func(a []html.Attribute) bool { return hasClass(a, "lntable") }); ti >= 0 {
		inner, _ := subtree(tokens, ti)
		first := findElement(inner, 0, "td", func(a []html.Attribute) bool { return hasClass(a, "lntd") })
		if first < 0 {
			return Region{}, false
		}
		_, after := subtree(inner, first)
		second := findElement(inner, after, "td", func(a []html.Attribute) bool { return hasClass(a, "lntd") })
		if second < 0 {
			return Region{}, false
		}
		cell, _ := subtree(inner, second)
		pi := findElement(cell, 0, "pre", nil)
		if pi < 0 {
			return Region{}, false
		}
		return chromaExtract(cell, pi)
	}

	pi := findElement(tokens, 1, "pre", nil)
	if pi < 0 {
		return Region{}, false
	}
	if !hasClass(tokens[pi].Attrs, "chroma") {
		// Inline styles mode has no chroma class; require a language
		// signature on the code element before claiming the region.
		inner, _ := subtree(tokens, pi)
		ci := findElement(inner, 0, "code", nil)
		if ci < 0 {
			return Region{}, false
		}
		if _, ok := attrValue(inner[ci].Attrs, "data-lang"); !ok {
			if _, ok := classWithPrefix(inner[ci].Attrs, "language-"); !ok {
				return Region{}, false
			}
		}
	}
	return chromaExtract(tokens, pi)
}

// chromaBarePre recognizes Chroma classes mode output without Hugo's wrapper
// div. It requires the chroma class so that plain fenced output from
// codeFences=false falls through to the generic fallback instead.
type chromaBarePre struct{}

func (chromaBarePre) Name() string { return "chroma-pre" }

func (chromaBarePre) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "pre") || !hasClass(tokens[0].Attrs, "chroma") {
		return Region{}, false
	}
	return chromaExtract(tokens, 0)
}

// chromaExtract pulls language, highlighted lines, and source text out of a
// Chroma pre element located at preIdx inside the token window.
func chromaExtract(tokens []BufferedToken, preIdx int) (Region, bool) {
	pre := tokens[preIdx]
	inner, _ := subtree(tokens, preIdx)
	ci := findElement(inner, 0, "code", nil)
	if ci < 0 {
		return Region{}, false
	}
	lang := ""
	if v, ok := attrValue(inner[ci].Attrs, "data-lang"); ok {
		lang = v
	} else if c, ok := classWithPrefix(inner[ci].Attrs, "language-"); ok {
		lang = strings.TrimPrefix(c, "language-")
	} else if v, ok := attrValue(pre.Attrs, "data-lang"); ok {
		lang = v
	}
	codeInner, _ := subtree(inner, ci)
	code := strings.TrimRight(recoverText(codeInner), "\n")
	highlighted := collectHighlightedLines(codeInner)
	return Region{Lang: lang, Meta: buildMeta(lang, highlighted), Code: code}, true
}
