package color

import (
	"math"
	"testing"
)

func TestParseRGBA(t *testing.T) {
	tests := []struct {
		input      string
		r, g, b, a float64
		wantErr    bool
	}{
		{"rgba(255,200,0,0.12)", 1.0, 200.0 / 255, 0, 0.12, false},
		{"rgba(46,160,67,0.12)", 46.0 / 255, 160.0 / 255, 67.0 / 255, 0.12, false},
		{"rgba(248,81,73,0.12)", 248.0 / 255, 81.0 / 255, 73.0 / 255, 0.12, false},
		{"rgba(0,0,0,1)", 0, 0, 0, 1, false},
		{"rgba(255,255,255,0)", 1, 1, 1, 0, false},
		{"  rgba(255,255,255,0.5)  ", 1, 1, 1, 0.5, false},
		{"rgb(255,0,0)", 0, 0, 0, 0, true},
		{"not a color", 0, 0, 0, 0, true},
		{"rgba(255,0,0)", 0, 0, 0, 0, true}, // only 3 components
	}

	for _, tt := range tests {
		r, g, b, a, err := ParseRGBA(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseRGBA(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRGBA(%q) error: %v", tt.input, err)
			continue
		}
		if math.Abs(r-tt.r) > 0.001 || math.Abs(g-tt.g) > 0.001 ||
			math.Abs(b-tt.b) > 0.001 || math.Abs(a-tt.a) > 0.001 {
			t.Errorf("ParseRGBA(%q) = (%f,%f,%f,%f), want (%f,%f,%f,%f)",
				tt.input, r, g, b, a, tt.r, tt.g, tt.b, tt.a)
		}
	}
}

func TestRGBAToHex(t *testing.T) {
	hex, err := RGBAToHex("rgba(255,200,0,0.12)")
	if err != nil {
		t.Fatalf("RGBAToHex error: %v", err)
	}
	if hex != "#ffc8001f" {
		t.Errorf("RGBAToHex = %q, want #ffc8001f", hex)
	}
}

func TestOnBackground_WithComposited(t *testing.T) {
	markerHex, _ := RGBAToHex("rgba(255,200,0,0.12)")
	result := OnBackground(markerHex, "#ffffff")
	if result == "#ffffff" {
		t.Error("composited color should differ from pure white background")
	}
	if result == markerHex {
		t.Error("composited color should differ from semi-transparent input")
	}
}

func TestParseHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		r, g, b float64
		wantErr bool
	}{
		{"3-char", "#fff", 1, 1, 1, false},
		{"6-char", "#ff0000", 1, 0, 0, false},
		{"6-char no hash", "00ff00", 0, 1, 0, false},
		{"8-char alpha", "#ff000080", 1, 0, 0, false},
		{"invalid length", "#abcde", 0, 0, 0, true},
		{"invalid chars", "#gggggg", 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, _, err := ParseHex(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHex(%q) expected error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseHex(%q) error: %v", tt.input, err)
			}
			if math.Abs(r-tt.r) > 0.01 || math.Abs(g-tt.g) > 0.01 || math.Abs(b-tt.b) > 0.01 {
				t.Errorf("ParseHex(%q) = (%f,%f,%f), want (%f,%f,%f)", tt.input, r, g, b, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestParseHex_8Char_Alpha(t *testing.T) {
	_, _, _, a, err := ParseHex("#ff000080")
	if err != nil {
		t.Fatalf("ParseHex error: %v", err)
	}
	if math.Abs(a-0.502) > 0.01 {
		t.Errorf("alpha = %f, want ~0.502", a)
	}
}

func TestGetLuminance(t *testing.T) {
	whiteLum := GetLuminance("#ffffff")
	if whiteLum < 0.99 {
		t.Errorf("white luminance = %f, want ~1.0", whiteLum)
	}

	blackLum := GetLuminance("#000000")
	if blackLum > 0.01 {
		t.Errorf("black luminance = %f, want ~0.0", blackLum)
	}
}

func TestIsLight(t *testing.T) {
	if !IsLight("#ffffff") {
		t.Error("white should be light")
	}
	if IsLight("#000000") {
		t.Error("black should not be light")
	}
	if IsLight("#1e1e1e") {
		t.Error("dark editor bg should not be light")
	}
}

func TestGetColorContrast(t *testing.T) {
	ratio := GetColorContrast("#000000", "#ffffff")
	if math.Abs(ratio-21) > 0.1 {
		t.Errorf("black/white contrast = %f, want ~21", ratio)
	}

	same := GetColorContrast("#ff0000", "#ff0000")
	if math.Abs(same-1) > 0.01 {
		t.Errorf("same color contrast = %f, want 1.0", same)
	}

	sym1 := GetColorContrast("#336699", "#ffffff")
	sym2 := GetColorContrast("#ffffff", "#336699")
	if math.Abs(sym1-sym2) > 0.001 {
		t.Errorf("contrast not symmetric: %f vs %f", sym1, sym2)
	}
}

func TestLighten(t *testing.T) {
	lighter := Lighten("#000000", 0.5)
	if lighter == "#000000" {
		t.Error("lightening black should produce a different color")
	}

	unchanged := Lighten("#336699", 0)
	if unchanged != "#336699" {
		t.Errorf("Lighten with amount=0 should return original, got %q", unchanged)
	}
}

func TestDarken(t *testing.T) {
	darker := Darken("#ffffff", 0.5)
	if darker == "#ffffff" {
		t.Error("darkening white should produce a different color")
	}

	unchanged := Darken("#336699", 0)
	if unchanged != "#336699" {
		t.Errorf("Darken with amount=0 should return original, got %q", unchanged)
	}
}

func TestMix(t *testing.T) {
	start := Mix("#ff0000", "#0000ff", 0)
	if start != "#ff0000" {
		t.Errorf("Mix amount=0 should return color1, got %q", start)
	}

	end := Mix("#ff0000", "#0000ff", 1)
	if end != "#0000ff" {
		t.Errorf("Mix amount=1 should return color2, got %q", end)
	}

	mid := Mix("#000000", "#ffffff", 0.5)
	r, g, b, _, _ := ParseHex(mid)
	if math.Abs(r-0.5) > 0.01 || math.Abs(g-0.5) > 0.01 || math.Abs(b-0.5) > 0.01 {
		t.Errorf("Mix 50/50 black+white = %q, expected mid-gray", mid)
	}
}

func TestSetAlpha(t *testing.T) {
	result := SetAlpha("#ff0000", 0.5)
	if len(result) != 9 {
		t.Errorf("SetAlpha should produce 9-char hex (#rrggbbaa), got %q", result)
	}
	_, _, _, a, err := ParseHex(result)
	if err != nil {
		t.Fatalf("ParseHex error on SetAlpha result: %v", err)
	}
	if math.Abs(a-0.5) > 0.01 {
		t.Errorf("alpha = %f, want ~0.5", a)
	}
}

func TestSetLuminance(t *testing.T) {
	result := SetLuminance("#808080", 0.3)
	lum := GetLuminance(result)
	if math.Abs(lum-0.3) > 0.05 {
		t.Errorf("SetLuminance target=0.3, got luminance=%f", lum)
	}
}

func TestSetLuminance_PureBlack(t *testing.T) {
	result := SetLuminance("#000000", 0.5)
	if result == "#000000" {
		t.Error("SetLuminance on pure black should produce a non-black color")
	}
}

func TestEnsureContrastOnBackground(t *testing.T) {
	high := EnsureContrastOnBackground("#000000", "#ffffff", 4.5)
	if high != "#000000" {
		t.Errorf("already high contrast should be unchanged, got %q", high)
	}

	adjusted := EnsureContrastOnBackground("#777777", "#888888", 4.5)
	ratio := GetColorContrast(adjusted, "#888888")
	if ratio < 4.5 {
		t.Errorf("adjusted color contrast = %f, want >= 4.5", ratio)
	}
}

func TestEnsureContrastOnBackground_LightOnLight(t *testing.T) {
	adjusted := EnsureContrastOnBackground("#eeeeee", "#ffffff", 4.5)
	ratio := GetColorContrast(adjusted, "#ffffff")
	if ratio < 4.5 {
		t.Errorf("adjusted contrast = %f, want >= 4.5", ratio)
	}
}

func TestToHex(t *testing.T) {
	if got := ToHex(1, 0, 0); got != "#ff0000" {
		t.Errorf("ToHex(1,0,0) = %q, want #ff0000", got)
	}
	if got := ToHex(0, 0, 0); got != "#000000" {
		t.Errorf("ToHex(0,0,0) = %q, want #000000", got)
	}
}

func TestToHexRGBA(t *testing.T) {
	if got := ToHexRGBA(1, 0, 0, 0.5); got != "#ff000080" {
		t.Errorf("ToHexRGBA(1,0,0,0.5) = %q, want #ff000080", got)
	}
}

// --- OKLCH conversion ---

func TestOKLCH_RoundTrip(t *testing.T) {
	colors := []string{"#ff0000", "#00ff00", "#0000ff", "#1e1e2e", "#f5f5f5", "#808080", "#cf222e", "#54aeff"}
	for _, hex := range colors {
		l, c, h, err := ToOKLCH(hex)
		if err != nil {
			t.Fatalf("ToOKLCH(%q) error: %v", hex, err)
		}
		got := FromOKLCH(l, c, h)
		r1, g1, b1, _, _ := ParseHex(hex)
		r2, g2, b2, _, err := ParseHex(got)
		if err != nil {
			t.Fatalf("FromOKLCH produced unparseable color %q for %q", got, hex)
		}
		const tol = 1.5 / 255
		if absDiff(r1, r2) > tol || absDiff(g1, g2) > tol || absDiff(b1, b2) > tol {
			t.Errorf("round trip %q -> %q exceeds tolerance", hex, got)
		}
	}
}

func TestOKLCH_KnownColors(t *testing.T) {
	// sRGB red: L=0.6280, C=0.2577, H=29.23 per the OKLab reference values.
	l, c, h, err := ToOKLCH("#ff0000")
	if err != nil {
		t.Fatalf("ToOKLCH error: %v", err)
	}
	if absDiff(l, 0.6280) > 0.001 {
		t.Errorf("red lightness = %f, want ~0.6280", l)
	}
	if absDiff(c, 0.2577) > 0.001 {
		t.Errorf("red chroma = %f, want ~0.2577", c)
	}
	if absDiff(h, 29.23) > 0.1 {
		t.Errorf("red hue = %f, want ~29.23", h)
	}

	// White is achromatic with full lightness.
	l, c, _, _ = ToOKLCH("#ffffff")
	if absDiff(l, 1.0) > 0.001 {
		t.Errorf("white lightness = %f, want ~1.0", l)
	}
	if c > 0.001 {
		t.Errorf("white chroma = %f, want ~0", c)
	}

	// Black has zero lightness.
	l, _, _, _ = ToOKLCH("#000000")
	if l > 0.001 {
		t.Errorf("black lightness = %f, want ~0", l)
	}
}

func TestFromOKLCH_GamutClamping(t *testing.T) {
	// Chroma 0.4 at this lightness/hue is far outside sRGB; the result must
	// still parse and land in gamut.
	got := FromOKLCH(0.5, 0.4, 145)
	if _, _, _, _, err := ParseHex(got); err != nil {
		t.Fatalf("clamped color %q does not parse: %v", got, err)
	}
	l, c, h, _ := ToOKLCH(got)
	if absDiff(l, 0.5) > 0.02 {
		t.Errorf("clamping should preserve lightness, got %f", l)
	}
	if c >= 0.4 {
		t.Errorf("clamping should reduce chroma, got %f", c)
	}
	if absDiff(h, 145) > 2 {
		t.Errorf("clamping should preserve hue, got %f", h)
	}
}

func TestSetHueChroma_PreservesLightness(t *testing.T) {
	before, _, _, _ := ToOKLCH("#3366cc")
	out := SetHueChroma("#3366cc", 145, 0.1)
	after, c, h, err := ToOKLCH(out)
	if err != nil {
		t.Fatalf("SetHueChroma produced unparseable color %q", out)
	}
	if absDiff(before, after) > 0.01 {
		t.Errorf("lightness changed: %f -> %f", before, after)
	}
	if absDiff(h, 145) > 1 {
		t.Errorf("hue = %f, want ~145", h)
	}
	if absDiff(c, 0.1) > 0.01 {
		t.Errorf("chroma = %f, want ~0.1", c)
	}
}

func TestSetHueChroma_PreservesAlpha(t *testing.T) {
	out := SetHueChroma("#3366cc80", 200, 0.05)
	_, _, _, a, err := ParseHex(out)
	if err != nil {
		t.Fatalf("result %q does not parse: %v", out, err)
	}
	if absDiff(a, 0.5) > 0.01 {
		t.Errorf("alpha = %f, want ~0.5", a)
	}
}

func TestSetHueChroma_InvalidColorPassesThrough(t *testing.T) {
	if got := SetHueChroma("not-a-color", 100, 0.1); got != "not-a-color" {
		t.Errorf("invalid input should pass through, got %q", got)
	}
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}
