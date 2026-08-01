package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// pygmentsDivWrapper recognizes the Pygments family (Sphinx, Pelican): an
// outer div whose class carries the highlight- prefix with the language as
// suffix, wrapping div.highlight > pre with no code element at all. The
// absence of a code element is part of the positive signature; it is what
// distinguishes this family from Rouge, whose wrapper is also literally
// named highlight. The language suffix is kept verbatim (python3 stays
// python3); mapping aliases is engine configuration, not recovery.
type pygmentsDivWrapper struct{}

func (pygmentsDivWrapper) Name() string { return "pygments-div" }

func (pygmentsDivWrapper) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "div") {
		return Region{}, false
	}
	langToken, ok := classWithPrefix(tokens[0].Attrs, "highlight-")
	if !ok {
		return Region{}, false
	}
	lang := strings.TrimPrefix(langToken, "highlight-")

	// Sphinx linenos mode wraps everything in table.highlighttable with a
	// td.linenos gutter cell and the real div.highlight inside td.code.
	window := tokens[1:]
	if ti := findElement(tokens, 1, "table", func(a []html.Attribute) bool { return hasClass(a, "highlighttable") }); ti >= 0 {
		inner, _ := subtree(tokens, ti)
		di := findElement(inner, 0, "td", func(a []html.Attribute) bool { return hasClass(a, "code") })
		if di < 0 {
			return Region{}, false
		}
		window, _ = subtree(inner, di)
	}
	code, ok := pygmentsExtract(window)
	if !ok {
		return Region{}, false
	}
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}

// pygmentsExtract locates div.highlight > pre inside the window, requires
// the pre to hold no code element, and recovers its text. The leading empty
// span Pygments emits contributes nothing and needs no special case.
func pygmentsExtract(tokens []BufferedToken) (string, bool) {
	hi := findElement(tokens, 0, "div", func(a []html.Attribute) bool { return hasClass(a, "highlight") })
	if hi < 0 {
		return "", false
	}
	inner, _ := subtree(tokens, hi)
	pi := findElement(inner, 0, "pre", nil)
	if pi < 0 {
		return "", false
	}
	preInner, _ := subtree(inner, pi)
	if findElement(preInner, 0, "code", nil) >= 0 {
		return "", false
	}
	return strings.TrimRight(recoverText(preInner), "\n"), true
}
