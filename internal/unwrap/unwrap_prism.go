package unwrap

import "strings"

// prismBarePre recognizes the Eleventy Prism plugin shape: the language
// class duplicated on both the pre and the code element. Requiring both
// keeps a hand written pre.language-x around a plain code element in the
// generic fallback where it belongs. Prism's optional line numbers plugin
// puts its span.line-numbers-rows gutter inside the pre but outside the
// code element, so walking only the code subtree already excludes it; the
// class stays in the shared gutter set for table shaped outputs where the
// gutter lives inside the recovered container.
type prismBarePre struct{}

func (prismBarePre) Name() string { return "prism-pre" }

func (prismBarePre) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "pre") {
		return Region{}, false
	}
	if _, ok := classWithPrefix(tokens[0].Attrs, "language-"); !ok {
		return Region{}, false
	}
	ci := findElement(tokens, 1, "code", nil)
	if ci < 0 {
		return Region{}, false
	}
	langToken, ok := classWithPrefix(tokens[ci].Attrs, "language-")
	if !ok {
		return Region{}, false
	}
	lang := strings.TrimPrefix(langToken, "language-")
	codeInner, _ := subtree(tokens, ci)
	code := strings.TrimRight(recoverText(codeInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}
