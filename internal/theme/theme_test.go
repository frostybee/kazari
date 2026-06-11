package theme

import (
	"strings"
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func testConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.DarkMode = config.DarkModeConfig{
		Kind:     config.DarkModeSelectorKind,
		Selector: ".dark",
	}
	return cfg
}

var lightColors = ThemeColors{EditorBG: "#ffffff", EditorFG: "#24292f"}
var darkColors = ThemeColors{EditorBG: "#1e1e1e", EditorFG: "#d4d4d4"}

// --- GenerateVars dark mode strategies ---

func TestGenerateVars_SelectorMode(t *testing.T) {
	cfg := testConfig()
	css := GenerateVars(cfg, lightColors, darkColors)

	if !strings.Contains(css, ":root {") {
		t.Error("should contain :root block")
	}
	if !strings.Contains(css, ":root.dark {") {
		t.Error("should contain :root.dark block")
	}
	if strings.Contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("selector mode should not contain media query")
	}
}

func TestGenerateVars_MediaQueryMode(t *testing.T) {
	cfg := testConfig()
	cfg.DarkMode.Kind = config.DarkModeMediaQueryKind
	css := GenerateVars(cfg, lightColors, darkColors)

	if !strings.Contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("should contain media query")
	}
	if strings.Contains(css, ":root.dark {") {
		t.Error("media query mode should not contain .dark selector")
	}
}

func TestGenerateVars_BothMode(t *testing.T) {
	cfg := testConfig()
	cfg.DarkMode.Kind = config.DarkModeBothKind
	css := GenerateVars(cfg, lightColors, darkColors)

	if !strings.Contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("should contain media query")
	}
	if !strings.Contains(css, ":root.dark {") {
		t.Error("should contain .dark selector")
	}
}

// --- Conditional features ---

func TestGenerateVars_CollapsibleEnabled(t *testing.T) {
	cfg := testConfig()
	cfg.Collapsible = &config.CollapsibleConfig{LineThreshold: 15}
	css := GenerateVars(cfg, lightColors, darkColors)

	if !strings.Contains(css, "--kz-collapse-btn-bg") {
		t.Error("should contain collapsible CSS variables")
	}
}

func TestGenerateVars_CollapsibleNil(t *testing.T) {
	cfg := testConfig()
	cfg.Collapsible = nil
	css := GenerateVars(cfg, lightColors, darkColors)

	if strings.Contains(css, "--kz-collapse-btn-bg") {
		t.Error("should not contain collapsible CSS variables when nil")
	}
}

func TestGenerateVars_CodeGroups(t *testing.T) {
	cfg := testConfig()
	cfg.CodeGroups = true
	css := GenerateVars(cfg, lightColors, darkColors)

	if !strings.Contains(css, "--kz-group-tab-bg") {
		t.Error("should contain code group CSS variables")
	}
}

func TestGenerateVars_DotsMinimal(t *testing.T) {
	cfg := testConfig()
	cfg.TerminalDotStyle = config.DotsMinimal
	css := GenerateVars(cfg, lightColors, darkColors)

	if !strings.Contains(css, "--kz-terminal-dots-opacity") {
		t.Error("should contain minimal dots CSS variables")
	}
}

// --- buildThemeVars luminance detection ---

func TestBuildThemeVars_LightBG(t *testing.T) {
	cfg := testConfig()
	vars := buildThemeVars(ThemeColors{EditorBG: "#ffffff", EditorFG: "#24292f"}, cfg)

	varMap := make(map[string]string)
	for _, v := range vars {
		varMap[v.name] = v.value
	}
	if varMap["--kz-gutter-border-color"] != "rgba(0,0,0,0.1)" {
		t.Errorf("light BG gutter border = %q, want rgba(0,0,0,0.1)", varMap["--kz-gutter-border-color"])
	}
}

func TestBuildThemeVars_DarkBG(t *testing.T) {
	cfg := testConfig()
	vars := buildThemeVars(ThemeColors{EditorBG: "#1e1e1e", EditorFG: "#d4d4d4"}, cfg)

	varMap := make(map[string]string)
	for _, v := range vars {
		varMap[v.name] = v.value
	}
	if varMap["--kz-gutter-border-color"] != "rgba(255,255,255,0.1)" {
		t.Errorf("dark BG gutter border = %q, want rgba(255,255,255,0.1)", varMap["--kz-gutter-border-color"])
	}
}

func TestBuildThemeVars_CustomLineNumberFG(t *testing.T) {
	cfg := testConfig()
	vars := buildThemeVars(ThemeColors{EditorBG: "#ffffff", EditorFG: "#24292f", LineNumberFG: "#custom"}, cfg)

	varMap := make(map[string]string)
	for _, v := range vars {
		varMap[v.name] = v.value
	}
	if varMap["--kz-ln-fg"] != "#custom" {
		t.Errorf("custom LineNumberFG should be used, got %q", varMap["--kz-ln-fg"])
	}
}

func TestBuildThemeVars_SelectionBG(t *testing.T) {
	cfg := testConfig()

	withSelection := buildThemeVars(ThemeColors{EditorBG: "#ffffff", EditorFG: "#24292f", SelectionBG: "#abc123"}, cfg)
	found := false
	for _, v := range withSelection {
		if v.name == "--kz-selection-bg" {
			found = true
			if v.value != "#abc123" {
				t.Errorf("SelectionBG = %q, want #abc123", v.value)
			}
		}
	}
	if !found {
		t.Error("should include --kz-selection-bg when set")
	}

	withoutSelection := buildThemeVars(ThemeColors{EditorBG: "#ffffff", EditorFG: "#24292f"}, cfg)
	for _, v := range withoutSelection {
		if v.name == "--kz-selection-bg" {
			t.Error("should not include --kz-selection-bg when empty")
		}
	}
}

// --- TokenSwitchingCSS ---

func TestTokenSwitchingCSS_NoDarkTheme(t *testing.T) {
	cfg := testConfig()
	cfg.DarkTheme = ""
	css := TokenSwitchingCSS(cfg)

	if !strings.Contains(css, "color: var(--sl)") {
		t.Error("should contain light color rule")
	}
	if strings.Contains(css, "var(--sd)") {
		t.Error("should not contain dark color rule when no dark theme")
	}
}

func TestTokenSwitchingCSS_SelectorMode(t *testing.T) {
	cfg := testConfig()
	css := TokenSwitchingCSS(cfg)

	if !strings.Contains(css, ".dark .kazari-code") {
		t.Error("should contain .dark selector rule")
	}
	if !strings.Contains(css, "var(--sd)") {
		t.Error("should contain dark color variable")
	}
}

func TestTokenSwitchingCSS_MediaQueryMode(t *testing.T) {
	cfg := testConfig()
	cfg.DarkMode.Kind = config.DarkModeMediaQueryKind
	css := TokenSwitchingCSS(cfg)

	if !strings.Contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("should contain media query")
	}
}

func TestTokenSwitchingCSS_BothMode(t *testing.T) {
	cfg := testConfig()
	cfg.DarkMode.Kind = config.DarkModeBothKind
	css := TokenSwitchingCSS(cfg)

	if !strings.Contains(css, "@media (prefers-color-scheme: dark)") {
		t.Error("should contain media query")
	}
	if !strings.Contains(css, ".dark .kazari-code") {
		t.Error("should contain .dark selector rule")
	}
}
