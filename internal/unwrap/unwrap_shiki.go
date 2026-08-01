package unwrap

import "strings"

// shikiAstroBarePre recognizes Astro's built in Shiki output: a pre with the
// astro-code class token, the language in data-language on the pre, lines as
// span.line, and class free token spans styled inline. The astro-code class
// requirement is the guard that keeps this unwrapper from claiming the
// Expressive Code shape used by Starlight, which also carries data-language
// but structures lines as div.ec-line and stays unclaimed in this phase.
type shikiAstroBarePre struct{}

func (shikiAstroBarePre) Name() string { return "shiki-astro" }

func (shikiAstroBarePre) Match(tokens []BufferedToken) (Region, bool) {
	if len(tokens) == 0 || !isStartTag(tokens[0], "pre") || !hasClass(tokens[0].Attrs, "astro-code") {
		return Region{}, false
	}
	lang, _ := attrValue(tokens[0].Attrs, "data-language")
	ci := findElement(tokens, 1, "code", nil)
	if ci < 0 {
		return Region{}, false
	}
	codeInner, _ := subtree(tokens, ci)
	code := strings.TrimRight(recoverText(codeInner), "\n")
	return Region{Lang: lang, Meta: buildMeta(lang, nil), Code: code}, true
}
