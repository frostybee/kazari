package kazari

import "github.com/frostybee/kazari/internal/svgutil"

// CreateInlineSVGURL converts SVG markup into a data:image/svg+xml URL for
// use in CSS custom properties (e.g. custom file icons or terminal icons).
// Use single quotes for SVG attributes; double quotes are percent-encoded.
func CreateInlineSVGURL(svg string) string {
	return svgutil.InlineSVGURL(svg)
}
