package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseHex parses a hex color string into r, g, b, a components (0-1 range).
// Supports #rgb, #rrggbb, #rrggbbaa formats (with or without #).
func ParseHex(hex string) (r, g, b, a float64, err error) {
	hex = strings.TrimPrefix(hex, "#")
	a = 1.0

	switch len(hex) {
	case 3:
		var ri, gi, bi uint8
		_, err = fmt.Sscanf(hex, "%1x%1x%1x", &ri, &gi, &bi)
		r = float64(ri*17) / 255
		g = float64(gi*17) / 255
		b = float64(bi*17) / 255
	case 6:
		var ri, gi, bi uint8
		_, err = fmt.Sscanf(hex, "%02x%02x%02x", &ri, &gi, &bi)
		r = float64(ri) / 255
		g = float64(gi) / 255
		b = float64(bi) / 255
	case 8:
		var ri, gi, bi, ai uint8
		_, err = fmt.Sscanf(hex, "%02x%02x%02x%02x", &ri, &gi, &bi, &ai)
		r = float64(ri) / 255
		g = float64(gi) / 255
		b = float64(bi) / 255
		a = float64(ai) / 255
	default:
		err = fmt.Errorf("invalid hex color length: %q", hex)
	}
	return
}

// ToHex converts r, g, b (0-1) to a "#rrggbb" string.
func ToHex(r, g, b float64) string {
	return fmt.Sprintf("#%02x%02x%02x",
		uint8(math.Round(r*255)),
		uint8(math.Round(g*255)),
		uint8(math.Round(b*255)),
	)
}

// ToHexRGBA converts r, g, b, a (0-1) to a "#rrggbbaa" string.
func ToHexRGBA(r, g, b, a float64) string {
	return fmt.Sprintf("#%02x%02x%02x%02x",
		uint8(math.Round(r*255)),
		uint8(math.Round(g*255)),
		uint8(math.Round(b*255)),
		uint8(math.Round(a*255)),
	)
}

// linearize converts an sRGB component to linear RGB.
func linearize(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// delinearize converts a linear RGB component back to sRGB.
func delinearize(c float64) float64 {
	if c <= 0.0031308 {
		return c * 12.92
	}
	return 1.055*math.Pow(c, 1.0/2.4) - 0.055
}

// GetLuminance returns the WCAG 2.1 relative luminance of a hex color (0-1).
func GetLuminance(hex string) float64 {
	r, g, b, _, err := ParseHex(hex)
	if err != nil {
		return 0
	}
	return 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b)
}

// IsLight returns true if the hex color has a relative luminance above 0.5.
func IsLight(hex string) bool {
	return GetLuminance(hex) > 0.5
}

// GetColorContrast returns the WCAG contrast ratio between two hex colors (1-21).
func GetColorContrast(color1, color2 string) float64 {
	l1 := GetLuminance(color1)
	l2 := GetLuminance(color2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// Lighten increases the luminance of a hex color by amount (0-1).
func Lighten(hex string, amount float64) string {
	r, g, b, _, err := ParseHex(hex)
	if err != nil {
		return hex
	}
	r = r + (1-r)*amount
	g = g + (1-g)*amount
	b = b + (1-b)*amount
	return ToHex(clamp01(r), clamp01(g), clamp01(b))
}

// Darken decreases the luminance of a hex color by amount (0-1).
func Darken(hex string, amount float64) string {
	r, g, b, _, err := ParseHex(hex)
	if err != nil {
		return hex
	}
	r = r * (1 - amount)
	g = g * (1 - amount)
	b = b * (1 - amount)
	return ToHex(clamp01(r), clamp01(g), clamp01(b))
}

// Mix blends two colors. amount=0 returns color1, amount=1 returns color2.
func Mix(color1, color2 string, amount float64) string {
	r1, g1, b1, _, err1 := ParseHex(color1)
	r2, g2, b2, _, err2 := ParseHex(color2)
	if err1 != nil || err2 != nil {
		return color1
	}
	r := r1 + (r2-r1)*amount
	g := g1 + (g2-g1)*amount
	b := b1 + (b2-b1)*amount
	return ToHex(r, g, b)
}

// SetAlpha returns the color with the given alpha (0-1) as #rrggbbaa.
func SetAlpha(hex string, alpha float64) string {
	r, g, b, _, err := ParseHex(hex)
	if err != nil {
		return hex
	}
	return ToHexRGBA(r, g, b, clamp01(alpha))
}

// OnBackground computes the opaque appearance of a semi-transparent color on a background.
func OnBackground(color, background string) string {
	r, g, b, a, err := ParseHex(color)
	if err != nil {
		return color
	}
	if a >= 1.0 {
		return color
	}
	br, bg, bb, _, err := ParseHex(background)
	if err != nil {
		return color
	}
	outR := r*a + br*(1-a)
	outG := g*a + bg*(1-a)
	outB := b*a + bb*(1-a)
	return ToHex(outR, outG, outB)
}

// EnsureContrastOnBackground adjusts a color's luminance to meet the minimum
// contrast ratio against the given background. Returns the adjusted hex color.
func EnsureContrastOnBackground(color, bg string, minContrast float64) string {
	if GetColorContrast(color, bg) >= minContrast {
		return color
	}

	r, g, b, _, err := ParseHex(color)
	if err != nil {
		return color
	}

	// Try both lighter and darker directions, pick the one that meets contrast first.
	// Use integer-indexed loop to avoid floating-point accumulation drift
	// (adding 0.05 nineteen times yields 1.00000000000000022, skipping the 20th step).
	for _, lighten := range []bool{true, false} {
		for i := 1; i <= 20; i++ {
			step := float64(i) * 0.05
			var cr, cg, cb float64
			if lighten {
				cr = r + (1-r)*step
				cg = g + (1-g)*step
				cb = b + (1-b)*step
			} else {
				cr = r * (1 - step)
				cg = g * (1 - step)
				cb = b * (1 - step)
			}
			candidate := ToHex(clamp01(cr), clamp01(cg), clamp01(cb))
			if GetColorContrast(candidate, bg) >= minContrast {
				return candidate
			}
		}
	}

	// Neither direction reached the threshold; return whichever extreme has higher contrast.
	if GetColorContrast("#000000", bg) > GetColorContrast("#ffffff", bg) {
		return "#000000"
	}
	return "#ffffff"
}

// SetLuminance sets the absolute luminance of a color (0-1).
func SetLuminance(hex string, targetLuminance float64) string {
	r, g, b, _, err := ParseHex(hex)
	if err != nil {
		return hex
	}

	currentLum := 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b)
	if currentLum == 0 {
		// Pure black — can only lighten uniformly.
		gray := delinearize(targetLuminance)
		return ToHex(gray, gray, gray)
	}

	ratio := targetLuminance / currentLum
	lr := linearize(r) * ratio
	lg := linearize(g) * ratio
	lb := linearize(b) * ratio

	return ToHex(
		clamp01(delinearize(lr)),
		clamp01(delinearize(lg)),
		clamp01(delinearize(lb)),
	)
}

// ToOKLCH converts a hex color to OKLCH components: lightness (0-1),
// chroma (0-0.4 typical), hue in degrees (0-360). Alpha is ignored.
// Uses the OKLab transform by Björn Ottosson.
func ToOKLCH(hex string) (l, c, h float64, err error) {
	r, g, b, _, err := ParseHex(hex)
	if err != nil {
		return 0, 0, 0, err
	}
	lab, labA, labB := linearSRGBToOKLab(linearize(r), linearize(g), linearize(b))
	c = math.Sqrt(labA*labA + labB*labB)
	h = math.Atan2(labB, labA) * 180 / math.Pi
	if h < 0 {
		h += 360
	}
	return lab, c, h, nil
}

// FromOKLCH converts OKLCH components back to a "#rrggbb" hex color.
// Colors outside the sRGB gamut are clamped by reducing chroma until they fit.
func FromOKLCH(l, c, h float64) string {
	hRad := h * math.Pi / 180
	for ; c >= 0; c -= 0.001 {
		labA := c * math.Cos(hRad)
		labB := c * math.Sin(hRad)
		r, g, b := okLabToLinearSRGB(l, labA, labB)
		if inUnitRange(r) && inUnitRange(g) && inUnitRange(b) {
			return ToHex(
				delinearize(clamp01(r)),
				delinearize(clamp01(g)),
				delinearize(clamp01(b)),
			)
		}
	}
	// Chroma 0 is always in gamut for l in [0,1]; clamp l as a last resort.
	r, g, b := okLabToLinearSRGB(clamp01(l), 0, 0)
	return ToHex(
		clamp01(delinearize(clamp01(r))),
		clamp01(delinearize(clamp01(g))),
		clamp01(delinearize(clamp01(b))),
	)
}

// SetHueChroma sets the OKLCH hue (degrees) and chroma of a hex color while
// preserving its lightness. Alpha is preserved.
func SetHueChroma(hex string, hue, chroma float64) string {
	l, _, _, err := ToOKLCH(hex)
	if err != nil {
		return hex
	}
	out := FromOKLCH(l, chroma, hue)
	_, _, _, a, _ := ParseHex(hex)
	if a < 1 {
		out = SetAlpha(out, a)
	}
	return out
}

func linearSRGBToOKLab(r, g, b float64) (okL, okA, okB float64) {
	l := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b

	lc := math.Cbrt(l)
	mc := math.Cbrt(m)
	sc := math.Cbrt(s)

	okL = 0.2104542553*lc + 0.7936177850*mc - 0.0040720468*sc
	okA = 1.9779984951*lc - 2.4285922050*mc + 0.4505937099*sc
	okB = 0.0259040371*lc + 0.7827717662*mc - 0.8086757660*sc
	return
}

func okLabToLinearSRGB(okL, okA, okB float64) (r, g, b float64) {
	lc := okL + 0.3963377774*okA + 0.2158037573*okB
	mc := okL - 0.1055613458*okA - 0.0638541728*okB
	sc := okL - 0.0894841775*okA - 1.2914855480*okB

	l := lc * lc * lc
	m := mc * mc * mc
	s := sc * sc * sc

	r = 4.0767416621*l - 3.3077115913*m + 0.2309699292*s
	g = -1.2684380046*l + 2.6097574011*m - 0.3413193965*s
	b = -0.0041960863*l - 0.7034186147*m + 1.7076147010*s
	return
}

func inUnitRange(v float64) bool {
	return v >= -1e-6 && v <= 1+1e-6
}

// ParseRGBA parses a CSS rgba(R,G,B,A) string into r, g, b, a components (0-1 range).
// R, G, B are integers 0-255; A is a float 0-1.
func ParseRGBA(s string) (r, g, b, a float64, err error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "rgba(") || !strings.HasSuffix(s, ")") {
		return 0, 0, 0, 0, fmt.Errorf("invalid rgba format: %q", s)
	}
	inner := s[5 : len(s)-1]
	parts := strings.Split(inner, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("rgba expects 4 components, got %d", len(parts))
	}
	ri, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid red: %w", err)
	}
	gi, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid green: %w", err)
	}
	bi, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid blue: %w", err)
	}
	af, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid alpha: %w", err)
	}
	return float64(ri) / 255, float64(gi) / 255, float64(bi) / 255, af, nil
}

// RGBAToHex converts a CSS rgba() string to a #rrggbbaa hex string.
func RGBAToHex(s string) (string, error) {
	r, g, b, a, err := ParseRGBA(s)
	if err != nil {
		return "", err
	}
	return ToHexRGBA(r, g, b, a), nil
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
