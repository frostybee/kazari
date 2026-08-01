// Package unwrap recovers source code, language, and meta information from
// code blocks found in HTML that static site generators have already built.
// It is the detection and extraction layer behind the kazari process command.
// Each supported generator family has an Unwrapper that recognizes its markup
// shape and recovers the original code text ready for RenderWithMeta.
package unwrap

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// BufferedToken is one HTML token captured together with its raw input bytes
// and the byte offset where those bytes start inside the tokenized input.
type BufferedToken struct {
	Type   html.TokenType
	Data   string // tag name for tag tokens, unescaped text for text tokens
	Attrs  []html.Attribute
	Raw    []byte
	Offset int
}

// TokenizeFragment tokenizes an entire HTML fragment into one flat token
// window. It performs no region discovery; callers hand the result to an
// Unwrapper or a chain. Offsets accumulate across the fragment so later
// phases can reuse the same bookkeeping when scanning whole files, but
// nothing in this package asserts Offset correctness beyond basic sanity.
func TokenizeFragment(src []byte) ([]BufferedToken, error) {
	z := html.NewTokenizer(bytes.NewReader(src))
	var tokens []BufferedToken
	offset := 0
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return tokens, nil
			}
			return nil, z.Err()
		}
		raw := append([]byte(nil), z.Raw()...)
		t := z.Token()
		tokens = append(tokens, BufferedToken{
			Type:   tt,
			Data:   t.Data,
			Attrs:  t.Attr,
			Raw:    raw,
			Offset: offset,
		})
		offset += len(raw)
	}
}
