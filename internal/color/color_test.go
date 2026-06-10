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
