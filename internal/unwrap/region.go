package unwrap

import "golang.org/x/net/html"

// Region is the extracted content of one recognized code block.
type Region struct {
	Lang string
	Meta string // full meta string with the language token first
	Code string // recovered source text, ready for RenderWithMeta
}

// Unwrapper recognizes one generator family's markup shape inside a buffered
// candidate region and extracts a Region from it. Match returns false when
// the tokens do not carry that family's signature, letting the next unwrapper
// in the chain try.
type Unwrapper interface {
	Name() string
	Match(tokens []BufferedToken) (Region, bool)
}

// isStartTag reports whether the token opens the given element.
func isStartTag(tok BufferedToken, tag string) bool {
	return tok.Type == html.StartTagToken && tok.Data == tag
}

// findElement returns the index of the first StartTagToken with the given
// tag name at or after from for which pred accepts its attributes. A nil
// pred accepts any attributes. It returns -1 when no such element exists.
func findElement(tokens []BufferedToken, from int, tag string, pred func([]html.Attribute) bool) int {
	for i := from; i < len(tokens); i++ {
		if tokens[i].Type != html.StartTagToken || tokens[i].Data != tag {
			continue
		}
		if pred == nil || pred(tokens[i].Attrs) {
			return i
		}
	}
	return -1
}

// subtree returns the tokens strictly inside the element opened at index i,
// plus the index just past its matching end tag. Depth is tracked on the
// element's own tag name only, so unbalanced markup in unrelated tags cannot
// desync the boundary. An unclosed element yields the rest of the window.
func subtree(tokens []BufferedToken, i int) ([]BufferedToken, int) {
	tag := tokens[i].Data
	depth := 1
	for j := i + 1; j < len(tokens); j++ {
		switch tokens[j].Type {
		case html.StartTagToken:
			if tokens[j].Data == tag {
				depth++
			}
		case html.EndTagToken:
			if tokens[j].Data == tag {
				depth--
				if depth == 0 {
					return tokens[i+1 : j], j + 1
				}
			}
		}
	}
	return tokens[i+1:], len(tokens)
}
