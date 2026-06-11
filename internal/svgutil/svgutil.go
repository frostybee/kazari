// Package svgutil encodes inline SVG markup as CSS-safe data URLs.
package svgutil

import "strings"

// dataURLEscaper percent-encodes the characters that break data URLs inside
// CSS url() values. Single quotes are left intact, so SVG attributes should
// use single quotes; double quotes are encoded for url("...") contexts.
var dataURLEscaper = strings.NewReplacer(
	"%", "%25",
	"<", "%3C",
	">", "%3E",
	"\"", "%22",
	"#", "%23",
	"{", "%7B",
	"}", "%7D",
)

// InlineSVGURL converts SVG markup into a data:image/svg+xml URL suitable
// for CSS properties such as background-image and mask-image.
func InlineSVGURL(svg string) string {
	return "data:image/svg+xml," + dataURLEscaper.Replace(svg)
}
