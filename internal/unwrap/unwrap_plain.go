package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// plainFallback recognizes the generic pre > code shape, the HTML5 spec
// convention emitted by mdBook, markdown-it, Hugo with codeFences disabled,
// and hand written pages. It runs last in the bare pre chain and matches
// even without a language class; the unlabeled block policy belongs to the
// processor configuration, not to detection. A pre without any code element
// falls through unclaimed. The one shape it deliberately declines is
// Expressive Code output, whose div.ec-line structure inside the code
// element carries no newline text nodes, so generic text recovery would
// silently glue every line together.
type plainFallback struct{}

func (plainFallback) Name() string { return "plain" }

func (plainFallback) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "pre") {
		return Region{}, false
	}
	ci := findElement(tokens, 1, "code", nil)
	if ci < 0 {
		return Region{}, false
	}
	codeInner, _ := subtree(tokens, ci)
	if findElement(codeInner, 0, "div", func(a []html.Attribute) bool { return hasClass(a, "ec-line") }) >= 0 {
		return Region{}, false
	}
	lang := ""
	if c, ok := classWithPrefix(tokens[ci].Attrs, "language-"); ok {
		lang = strings.TrimPrefix(c, "language-")
	} else if v, ok := attrValue(tokens[ci].Attrs, "data-lang"); ok {
		lang = v
	} else if c, ok := classWithPrefix(tokens[0].Attrs, "language-"); ok {
		lang = strings.TrimPrefix(c, "language-")
	}
	code := strings.TrimRight(recoverText(codeInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}
