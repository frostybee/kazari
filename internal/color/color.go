package color

import (
	"fmt"
	"math"
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

	bgLum := GetLuminance(bg)

	// Try both lighter and darker directions, pick the one that meets contrast first.
	for _, lighten := range []bool{true, false} {
		for step := 0.05; step <= 1.0; step += 0.05 {
			var candidate string
			if lighten {
				cr := r + (1-r)*step
				cg := g + (1-g)*step
				cb := b + (1-b)*step
				candidate = ToHex(clamp01(cr), clamp01(cg), clamp01(cb))
			} else {
				cr := r * (1 - step)
				cg := g * (1 - step)
				cb := b * (1 - step)
				candidate = ToHex(clamp01(cr), clamp01(cg), clamp01(cb))
			}
			if GetColorContrast(candidate, bg) >= minContrast {
				return candidate
			}
		}
	}

	// If neither direction works well, choose based on background luminance.
	if bgLum > 0.5 {
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

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
