package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// rougeDivWrapper recognizes Jekyll's kramdown output: an outer div carrying
// both a language class token and the highlighter-rouge token, wrapping
// div.highlight > pre.highlight > code. The language lives only in the outer
// class list; no data-lang exists anywhere in this shape.
type rougeDivWrapper struct{}

func (rougeDivWrapper) Name() string { return "rouge-div" }

func (rougeDivWrapper) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "div") {
		return Region{}, false
	}
	langToken, ok := classWithPrefix(tokens[0].Attrs, "language-")
	if !ok || !hasClass(tokens[0].Attrs, "highlighter-rouge") {
		return Region{}, false
	}
	pi := findElement(tokens, 1, "pre", func(a []html.Attribute) bool { return hasClass(a, "highlight") })
	if pi < 0 {
		return Region{}, false
	}
	inner, _ := subtree(tokens, pi)
	ci := findElement(inner, 0, "code", nil)
	if ci < 0 {
		// The presence of a code element is what separates Rouge's
		// div.highlight from the Pygments family, which never emits one.
		return Region{}, false
	}
	lang := strings.TrimPrefix(langToken, "language-")
	codeInner, _ := subtree(inner, ci)
	code := strings.TrimRight(recoverText(codeInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}

// rougeHighlightFigure recognizes the legacy Liquid highlight tag output:
// figure.highlight > pre > code with language class and data-lang.
type rougeHighlightFigure struct{}

func (rougeHighlightFigure) Name() string { return "rouge-figure" }

func (rougeHighlightFigure) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "figure") || !hasClass(tokens[0].Attrs, "highlight") {
		return Region{}, false
	}
	pi := findElement(tokens, 1, "pre", nil)
	if pi < 0 {
		return Region{}, false
	}
	inner, _ := subtree(tokens, pi)
	ci := findElement(inner, 0, "code", nil)
	if ci < 0 {
		return Region{}, false
	}
	lang, ok := attrValue(inner[ci].Attrs, "data-lang")
	if !ok {
		c, okc := classWithPrefix(inner[ci].Attrs, "language-")
		if !okc {
			return Region{}, false
		}
		lang = strings.TrimPrefix(c, "language-")
	}
	codeInner, _ := subtree(inner, ci)
	code := strings.TrimRight(recoverText(codeInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}

// rougeLinenosTable recognizes Rouge's line number table: a rouge-table with
// a gutter cell and a td.code cell holding the real content. It is tried
// before the other Rouge shapes because the table nests inside them.
type rougeLinenosTable struct{}

func (rougeLinenosTable) Name() string { return "rouge-table" }

func (rougeLinenosTable) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || tokens[0].Type != html.StartTagToken {
		return Region{}, false
	}
	if tokens[0].Data != "div" && tokens[0].Data != "figure" {
		return Region{}, false
	}
	ti := findElement(tokens, 1, "table", func(a []html.Attribute) bool { return hasClass(a, "rouge-table") })
	if ti < 0 {
		return Region{}, false
	}
	lang := rougeTableLang(tokens[:ti])
	inner, _ := subtree(tokens, ti)
	di := findElement(inner, 0, "td", func(a []html.Attribute) bool { return hasClass(a, "code") })
	if di < 0 {
		return Region{}, false
	}
	cell, _ := subtree(inner, di)
	pi := findElement(cell, 0, "pre", nil)
	if pi < 0 {
		return Region{}, false
	}
	preInner, _ := subtree(cell, pi)
	code := strings.TrimRight(recoverText(preInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}

// rougeTableLang finds the language among the elements that wrap the table:
// a data-lang or language class on a code element, or a language token on
// the region root.
func rougeTableLang(tokens []BufferedToken) string {
	for _, tok := range tokens {
		if tok.Type != html.StartTagToken || tok.Data != "code" {
			continue
		}
		if v, ok := attrValue(tok.Attrs, "data-lang"); ok {
			return v
		}
		if c, ok := classWithPrefix(tok.Attrs, "language-"); ok {
			return strings.TrimPrefix(c, "language-")
		}
	}
	if len(tokens) > 0 {
		if c, ok := classWithPrefix(tokens[0].Attrs, "language-"); ok {
			return strings.TrimPrefix(c, "language-")
		}
	}
	return ""
}
