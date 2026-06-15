package kazari

import "testing"

func TestStyleValue(t *testing.T) {
	tests := map[string]struct {
		sv        StyleValue
		isThemed  bool
		lightVal  string
		darkVal   string
	}{
		"universal": {
			sv:       StyleValue{Value: "0.5rem"},
			isThemed: false,
			lightVal: "0.5rem",
			darkVal:  "0.5rem",
		},
		"themed both": {
			sv:       StyleValue{Dark: "#111", Light: "#fff"},
			isThemed: true,
			lightVal: "#fff",
			darkVal:  "#111",
		},
		"dark only": {
			sv:       StyleValue{Dark: "none"},
			isThemed: true,
			lightVal: "",
			darkVal:  "none",
		},
		"light only": {
			sv:       StyleValue{Light: "bold"},
			isThemed: true,
			lightVal: "bold",
			darkVal:  "",
		},
		"empty": {
			sv:       StyleValue{},
			isThemed: false,
			lightVal: "",
			darkVal:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.sv.IsThemed(); got != tc.isThemed {
				t.Errorf("IsThemed() = %v, want %v", got, tc.isThemed)
			}
			if got := tc.sv.LightValue(); got != tc.lightVal {
				t.Errorf("LightValue() = %q, want %q", got, tc.lightVal)
			}
			if got := tc.sv.DarkValue(); got != tc.darkVal {
				t.Errorf("DarkValue() = %q, want %q", got, tc.darkVal)
			}
		})
	}
}
