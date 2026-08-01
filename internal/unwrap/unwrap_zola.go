package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// Zola's class names are derived from the configured highlight theme (the
// pre carries the theme name, each line wrapper carries theme-l), so no
// literal class can be matched. The stable signature, verified live, is a
// code element with data-lang inside a pre whose inline style sets a
// background color. Both halves are required: data-lang alone is too weak
// and would claim unrelated hand written markup.

// zolaBarePre recognizes Zola output rooted at the pre itself.
type zolaBarePre struct{}

func (zolaBarePre) Name() string { return "zola-pre" }

func (zolaBarePre) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "pre") {
		return Region{}, false
	}
	return zolaExtract(tokens, 0)
}

// zolaDivWrapper recognizes Zola output inside a theme added wrapper div.
type zolaDivWrapper struct{}

func (zolaDivWrapper) Name() string { return "zola-div" }

func (zolaDivWrapper) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "div") {
		return Region{}, false
	}
	pi := findElement(tokens, 1, "pre", nil)
	if pi < 0 {
		return Region{}, false
	}
	return zolaExtract(tokens, pi)
}

func zolaExtract(tokens []BufferedToken, preIdx int) (Region, bool) {
	style, ok := attrValue(tokens[preIdx].Attrs, "style")
	if !ok || !strings.Contains(style, "background-color") {
		return Region{}, false
	}
	inner, _ := subtree(tokens, preIdx)
	ci := findElement(inner, 0, "code", func(a []html.Attribute) bool {
		_, has := attrValue(a, "data-lang")
		return has
	})
	if ci < 0 {
		return Region{}, false
	}
	lang, _ := attrValue(inner[ci].Attrs, "data-lang")
	codeInner, _ := subtree(inner, ci)
	code := strings.TrimRight(recoverText(codeInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}
