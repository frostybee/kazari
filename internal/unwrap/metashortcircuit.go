package unwrap

import "golang.org/x/net/html"

// kzMetaAttr finds a data-kz-meta attribute on the region root or on any pre
// or code element inside the region. A render hook template stashes the full
// fence info string there so features beyond the language survive the build.
// The value is returned verbatim; internal/meta parses it later inside
// RenderWithMeta. Attribute values arrive already entity unescaped from the
// tokenizer.
func kzMetaAttr(tokens []BufferedToken) (string, bool) {
	for i, tok := range tokens {
		if tok.Type != html.StartTagToken {
			continue
		}
		if i > 0 && tok.Data != "pre" && tok.Data != "code" {
			continue
		}
		if v, ok := attrValue(tok.Attrs, "data-kz-meta"); ok {
			return v, true
		}
	}
	return "", false
}
