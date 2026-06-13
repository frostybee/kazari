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

// --- Per-block theme override ---

func TestBlockOverrideStyle_DualTheme(t *testing.T) {
	cfg := testConfig()
	style := BlockOverrideStyle(cfg, darkColors, lightColors)

	if !strings.Contains(style, "--kz-ovl-editor-bg:#1e1e1e") {
		t.Error("should emit light slot editor bg from the override theme")
	}
	if !strings.Contains(style, "--kz-ovl-editor-fg:#d4d4d4") {
		t.Error("should emit light slot editor fg")
	}
	if !strings.Contains(style, "--kz-ovd-editor-bg:#ffffff") {
		t.Error("should emit dark slot editor bg")
	}
	if !strings.Contains(style, "--kz-ovl-toolbar-bg:") {
		t.Error("should emit toolbar vars")
	}
	if !strings.Contains(style, "--kz-ovl-ln-fg:") {
		t.Error("should emit line number fg")
	}
	if strings.HasSuffix(style, ";") {
		t.Error("style should not end with a trailing semicolon")
	}
}

func TestBlockOverrideStyle_SingleThemeEngine(t *testing.T) {
	cfg := testConfig()
	cfg.DarkTheme = ""
	style := BlockOverrideStyle(cfg, darkColors, ThemeColors{})

	if !strings.Contains(style, "--kz-ovl-editor-bg:#1e1e1e") {
		t.Error("should emit light slot vars")
	}
	if strings.Contains(style, "--kz-ovd-") {
		t.Error("single theme engine should not emit dark slot vars")
	}
}

func TestBlockOverrideStyle_UnusableLight(t *testing.T) {
	cfg := testConfig()
	if style := BlockOverrideStyle(cfg, ThemeColors{}, darkColors); style != "" {
		t.Errorf("unusable light colors should produce empty style, got %q", style)
	}
}

func TestBlockOverrideStyle_FeatureGates(t *testing.T) {
	cfg := testConfig()
	cfg.Collapsible = nil
	cfg.ThemedScrollbars = false
	style := BlockOverrideStyle(cfg, darkColors, lightColors)
	if strings.Contains(style, "collapse-btn") || strings.Contains(style, "scrollbar-thumb") {
		t.Error("gated clusters should be absent when features are off")
	}

	cfg.Collapsible = &config.CollapsibleConfig{LineThreshold: 10}
	cfg.ThemedScrollbars = true
	style = BlockOverrideStyle(cfg, darkColors, lightColors)
	if !strings.Contains(style, "--kz-ovl-collapse-btn-fg:") {
		t.Error("collapse button vars should be present when collapsible is enabled")
	}
	if !strings.Contains(style, "--kz-ovl-scrollbar-thumb:") {
		t.Error("scrollbar vars should be present when themed scrollbars are enabled")
	}
}

func TestTokenSwitchingCSS_ThemedMapping(t *testing.T) {
	cfg := testConfig()
	css := TokenSwitchingCSS(cfg)

	if !strings.Contains(css, ".kazari-code.kz-themed { --kz-editor-bg: var(--kz-ovl-editor-bg);") {
		t.Error("should contain light kz-themed mapping rule")
	}
	if strings.Contains(css, "--kz-collapse-gradient-end") {
		t.Error("gradient end re declaration should be absent when collapsible is disabled")
	}
	if !strings.Contains(css, ".dark .kazari-code.kz-themed { --kz-editor-bg: var(--kz-ovd-editor-bg, var(--kz-ovl-editor-bg));") {
		t.Error("dark mapping should be selector scoped with ovl fallbacks")
	}

	cfg.Collapsible = &config.CollapsibleConfig{LineThreshold: 10}
	css = TokenSwitchingCSS(cfg)
	if !strings.Contains(css, "--kz-collapse-gradient-end: var(--kz-editor-bg);") {
		t.Error("light mapping must re declare the collapse gradient end when collapsible is enabled")
	}
}

func TestTokenSwitchingCSS_ThemedMapping_DarkRuleSingleLine(t *testing.T) {
	cfg := testConfig()
	rule := themedDarkRule(cfg)
	if strings.Count(rule, "\n") != 1 || !strings.HasSuffix(rule, "}\n") {
		t.Errorf("dark mapping rule must be a single line for writeScopedRules, got %q", rule)
	}
}

func TestTokenSwitchingCSS_ThemedMapping_NoDarkTheme(t *testing.T) {
	cfg := testConfig()
	cfg.DarkTheme = ""
	css := TokenSwitchingCSS(cfg)

	if !strings.Contains(css, ".kazari-code.kz-themed { --kz-editor-bg: var(--kz-ovl-editor-bg);") {
		t.Error("light mapping must be present for single theme engines")
	}
	if strings.Contains(css, "--kz-ovd-") {
		t.Error("dark mapping should be absent when no dark theme is configured")
	}
}

func TestTokenSwitchingCSS_ThemedMapping_MediaQueryMode(t *testing.T) {
	cfg := testConfig()
	cfg.DarkMode.Kind = config.DarkModeMediaQueryKind
	css := TokenSwitchingCSS(cfg)

	mediaIdx := strings.Index(css, "@media (prefers-color-scheme: dark)")
	darkMapIdx := strings.Index(css, "var(--kz-ovd-editor-bg, var(--kz-ovl-editor-bg))")
	if mediaIdx < 0 || darkMapIdx < 0 || darkMapIdx < mediaIdx {
		t.Error("dark mapping should appear inside the media query block")
	}
}
