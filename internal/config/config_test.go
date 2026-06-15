package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	t.Run("boolean defaults", func(t *testing.T) {
		checks := map[string]struct{ got, want bool }{
			"CopyButton":               {cfg.CopyButton, true},
			"FullscreenButton":         {cfg.FullscreenButton, true},
			"FrameDetection":           {cfg.FrameDetection, true},
			"FileNameExtraction":       {cfg.FileNameExtraction, true},
			"LanguageBadge":            {cfg.LanguageBadge, true},
			"StyleReset":               {cfg.StyleReset, true},
			"ThemedScrollbars":         {cfg.ThemedScrollbars, true},
			"ThemedSelection":          {cfg.ThemedSelection, false},
			"ContentExclusion":         {cfg.ContentExclusion, false},
			"Minify":                   {cfg.Minify, true},
			"TerminalCommentStripping": {cfg.TerminalCommentStripping, true},
		}
		for name, c := range checks {
			if c.got != c.want {
				t.Errorf("%s = %v, want %v", name, c.got, c.want)
			}
		}
	})

	t.Run("string and numeric defaults", func(t *testing.T) {
		if cfg.LightTheme != "github-light" {
			t.Errorf("LightTheme = %q, want %q", cfg.LightTheme, "github-light")
		}
		if cfg.DarkTheme != "github-dark" {
			t.Errorf("DarkTheme = %q, want %q", cfg.DarkTheme, "github-dark")
		}
		if cfg.CascadeLayer != "kazari" {
			t.Errorf("CascadeLayer = %q, want %q", cfg.CascadeLayer, "kazari")
		}
		if cfg.TabWidth != 2 {
			t.Errorf("TabWidth = %d, want 2", cfg.TabWidth)
		}
		if cfg.MinContrast != 5.5 {
			t.Errorf("MinContrast = %f, want 5.5", cfg.MinContrast)
		}
		if cfg.TerminalDotStyle != DotsColored {
			t.Errorf("TerminalDotStyle = %d, want DotsColored(%d)", cfg.TerminalDotStyle, DotsColored)
		}
		if cfg.DarkMode.Kind != DarkModeSelectorKind {
			t.Errorf("DarkMode.Kind = %d, want DarkModeSelectorKind", cfg.DarkMode.Kind)
		}
		if cfg.DarkMode.Selector != ".dark" {
			t.Errorf("DarkMode.Selector = %q, want %q", cfg.DarkMode.Selector, ".dark")
		}
	})

	t.Run("pointer fields nil", func(t *testing.T) {
		if cfg.Collapsible != nil {
			t.Error("Collapsible should be nil")
		}
		if cfg.LightMarkerBGs != nil {
			t.Error("LightMarkerBGs should be nil")
		}
		if cfg.DarkMarkerBGs != nil {
			t.Error("DarkMarkerBGs should be nil")
		}
		if cfg.LanguageDefaults != nil {
			t.Error("LanguageDefaults should be nil")
		}
		if cfg.LanguageAliases != nil {
			t.Error("LanguageAliases should be nil")
		}
	})
}

func boolPtr(b bool) *bool { return &b }
func intPtr(i int) *int    { return &i }

func TestResolve_EngineDefaultsOnly(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Defaults.Wrap = true
	cfg.Defaults.LineNumbers = true

	resolved := cfg.Resolve("go", nil)

	if resolved.Lang != "go" {
		t.Errorf("Lang = %q, want %q", resolved.Lang, "go")
	}
	if !resolved.Wrap {
		t.Error("Wrap should inherit engine default (true)")
	}
	if !resolved.LineNumbers {
		t.Error("LineNumbers should inherit engine default (true)")
	}
	if resolved.StartLineNumber != 1 {
		t.Errorf("StartLineNumber = %d, want 1", resolved.StartLineNumber)
	}
}

func TestResolve_LanguageDefaultsOverrideEngine(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Defaults.Wrap = false
	cfg.LanguageDefaults = map[string]BlockDefaults{
		"bash,sh,zsh": {Wrap: true, LineNumbers: true},
	}

	resolved := cfg.Resolve("bash", nil)
	if !resolved.Wrap {
		t.Error("Wrap should be overridden to true by language defaults")
	}
	if !resolved.LineNumbers {
		t.Error("LineNumbers should be overridden to true by language defaults")
	}
}

func TestResolve_CommaKeyMatching(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageDefaults = map[string]BlockDefaults{
		"bash,sh,zsh": {Wrap: true},
	}

	for _, lang := range []string{"bash", "sh", "zsh"} {
		resolved := cfg.Resolve(lang, nil)
		if !resolved.Wrap {
			t.Errorf("Wrap should be true for %q (matched via comma-separated key)", lang)
		}
	}
}

func TestResolve_BlockOptsOverrideLanguageDefaults(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageDefaults = map[string]BlockDefaults{
		"bash": {Wrap: true, LineNumbers: true},
	}

	resolved := cfg.Resolve("bash", &BlockOptions{
		Wrap:        boolPtr(false),
		LineNumbers: boolPtr(false),
	})
	if resolved.Wrap {
		t.Error("Wrap should be overridden to false by block opts")
	}
	if resolved.LineNumbers {
		t.Error("LineNumbers should be overridden to false by block opts")
	}
}

func TestResolve_PointerFalseApplied(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Defaults.LineNumbers = true

	resolved := cfg.Resolve("go", &BlockOptions{
		LineNumbers: boolPtr(false),
	})
	if resolved.LineNumbers {
		t.Error("LineNumbers *bool = false should override the default true")
	}
}

func TestResolve_StartLineNumberZero(t *testing.T) {
	cfg := DefaultConfig()
	resolved := cfg.Resolve("go", &BlockOptions{
		StartLineNumber: intPtr(0),
	})
	if resolved.StartLineNumber != 0 {
		t.Errorf("StartLineNumber = %d, want 0 (pointer-based override)", resolved.StartLineNumber)
	}
}

func TestResolve_FrameAutoOverride(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Defaults.Frame = FrameCode

	resolved := cfg.Resolve("go", &BlockOptions{
		Frame: intPtr(FrameAuto),
	})
	if resolved.Frame != FrameAuto {
		t.Errorf("Frame = %d, want FrameAuto(%d)", resolved.Frame, FrameAuto)
	}
}

func TestResolve_NilLanguageDefaults(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageDefaults = nil

	resolved := cfg.Resolve("bash", nil)
	if resolved == nil {
		t.Fatal("Resolve should not return nil with nil LanguageDefaults")
	}
}

func TestResolve_NoLanguageMatch(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Defaults.Wrap = true
	cfg.LanguageDefaults = map[string]BlockDefaults{
		"python": {Wrap: false},
	}

	resolved := cfg.Resolve("go", nil)
	if !resolved.Wrap {
		t.Error("Wrap should keep engine default when language doesn't match")
	}
}

func TestResolveLanguage_AliasFound(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageAliases = map[string]string{
		"javascript": "js",
		"typescript": "ts",
	}

	if got := cfg.ResolveLanguage("javascript"); got != "js" {
		t.Errorf("ResolveLanguage(%q) = %q, want %q", "javascript", got, "js")
	}
}

func TestResolveLanguage_NoAlias(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageAliases = map[string]string{"javascript": "js"}

	if got := cfg.ResolveLanguage("python"); got != "python" {
		t.Errorf("ResolveLanguage(%q) = %q, want %q", "python", got, "python")
	}
}

func TestResolveLanguage_NilAliases(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageAliases = nil

	if got := cfg.ResolveLanguage("Go"); got != "go" {
		t.Errorf("ResolveLanguage(%q) = %q, want %q", "Go", got, "go")
	}
}

func TestResolveLanguage_MixedCase(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LanguageAliases = map[string]string{"javascript": "js"}

	if got := cfg.ResolveLanguage("JavaScript"); got != "js" {
		t.Errorf("ResolveLanguage(%q) = %q, want %q", "JavaScript", got, "js")
	}
}

func TestStyleValue(t *testing.T) {
	tests := map[string]struct {
		sv       StyleValue
		isThemed bool
		lightVal string
		darkVal  string
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
