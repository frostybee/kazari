package kazari

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/theme"
)

// --- Enum parsers ---

func TestParseFrame(t *testing.T) {
	tests := map[string]Frame{
		"auto": FrameAuto, "code": FrameCode,
		"terminal": FrameTerminal, "none": FrameNone,
	}
	for s, want := range tests {
		got, err := parseFrame(s)
		if err != nil {
			t.Errorf("parseFrame(%q) error: %v", s, err)
		}
		if got != want {
			t.Errorf("parseFrame(%q) = %d, want %d", s, got, want)
		}
	}
	if _, err := parseFrame("invalid"); err == nil {
		t.Error("parseFrame(invalid) should error")
	}
}

func TestParseCollapseStyle(t *testing.T) {
	tests := map[string]CollapseStyle{
		"github": CollapseGithub, "collapsibleStart": CollapseCollapsibleStart,
		"collapsibleEnd": CollapseCollapsibleEnd, "collapsibleAuto": CollapseCollapsibleAuto,
	}
	for s, want := range tests {
		got, err := parseCollapseStyle(s)
		if err != nil {
			t.Errorf("parseCollapseStyle(%q) error: %v", s, err)
		}
		if got != want {
			t.Errorf("parseCollapseStyle(%q) = %d, want %d", s, got, want)
		}
	}
	if _, err := parseCollapseStyle("invalid"); err == nil {
		t.Error("parseCollapseStyle(invalid) should error")
	}
}

func TestParseLangIconMode(t *testing.T) {
	tests := map[string]LangIconMode{
		"none": LangIconNone, "iconOnly": LangIconOnly,
		"iconAndText": LangIconAndText,
	}
	for s, want := range tests {
		got, err := parseLangIconModeStr(s)
		if err != nil {
			t.Errorf("parseLangIconModeStr(%q) error: %v", s, err)
		}
		if got != want {
			t.Errorf("parseLangIconModeStr(%q) = %d, want %d", s, got, want)
		}
	}
	if _, err := parseLangIconModeStr("invalid"); err == nil {
		t.Error("parseLangIconModeStr(invalid) should error")
	}
}

func TestParseTerminalDotStyle(t *testing.T) {
	tests := map[string]TerminalDotStyle{
		"colored": DotsColored, "minimal": DotsMinimal,
	}
	for s, want := range tests {
		got, err := parseTerminalDotStyleStr(s)
		if err != nil {
			t.Errorf("parseTerminalDotStyleStr(%q) error: %v", s, err)
		}
		if got != want {
			t.Errorf("parseTerminalDotStyleStr(%q) = %d, want %d", s, got, want)
		}
	}
	if _, err := parseTerminalDotStyleStr("invalid"); err == nil {
		t.Error("parseTerminalDotStyleStr(invalid) should error")
	}
}

func TestParseDarkModeConfig(t *testing.T) {
	t.Run("selector", func(t *testing.T) {
		dm, err := parseDarkModeConfig(&DarkModeFileConfig{Kind: "selector", Selector: ".dark"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := dm.(selectorMode); !ok {
			t.Error("should be selectorMode")
		}
	})
	t.Run("mediaQuery", func(t *testing.T) {
		dm, err := parseDarkModeConfig(&DarkModeFileConfig{Kind: "mediaQuery"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := dm.(mediaQueryMode); !ok {
			t.Error("should be mediaQueryMode")
		}
	})
	t.Run("both", func(t *testing.T) {
		dm, err := parseDarkModeConfig(&DarkModeFileConfig{Kind: "both", Selector: ".dark"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := dm.(bothMode); !ok {
			t.Error("should be bothMode")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if _, err := parseDarkModeConfig(&DarkModeFileConfig{Kind: "invalid"}); err == nil {
			t.Error("should error on invalid kind")
		}
	})
}

// --- ParseConfig ---

func TestParseConfig_YAML(t *testing.T) {
	yaml := `
copyButton: true
tabWidth: 4
themes:
  light: "github-light"
  dark: "github-dark"
`
	fc, err := ParseConfig([]byte(yaml), "yaml")
	if err != nil {
		t.Fatal(err)
	}
	if fc.CopyButton == nil || !*fc.CopyButton {
		t.Error("copyButton should be true")
	}
	if fc.TabWidth == nil || *fc.TabWidth != 4 {
		t.Error("tabWidth should be 4")
	}
	if fc.Themes == nil || fc.Themes.Light != "github-light" {
		t.Error("themes.light should be github-light")
	}
}

func TestParseConfig_JSON(t *testing.T) {
	js := `{"copyButton": false, "tabWidth": 2}`
	fc, err := ParseConfig([]byte(js), "json")
	if err != nil {
		t.Fatal(err)
	}
	if fc.CopyButton == nil || *fc.CopyButton {
		t.Error("copyButton should be false")
	}
	if fc.TabWidth == nil || *fc.TabWidth != 2 {
		t.Error("tabWidth should be 2")
	}
}

func TestParseConfig_PartialConfig(t *testing.T) {
	yaml := `copyButton: true`
	fc, err := ParseConfig([]byte(yaml), "yaml")
	if err != nil {
		t.Fatal(err)
	}
	if fc.TabWidth != nil {
		t.Error("omitted tabWidth should be nil")
	}
	if fc.Minify != nil {
		t.Error("omitted minify should be nil")
	}
}

func TestParseConfig_BoolFalseVsOmitted(t *testing.T) {
	yaml := `copyButton: false`
	fc, err := ParseConfig([]byte(yaml), "yaml")
	if err != nil {
		t.Fatal(err)
	}
	if fc.CopyButton == nil {
		t.Fatal("explicit false should not be nil")
	}
	if *fc.CopyButton {
		t.Error("explicit false should be false")
	}
}

func TestParseConfig_UnknownFormat(t *testing.T) {
	if _, err := ParseConfig([]byte("{}"), "toml"); err == nil {
		t.Error("unknown format should error")
	}
}

func TestParseConfig_UnknownYAMLKey(t *testing.T) {
	yaml := `unknownField: true`
	if _, err := ParseConfig([]byte(yaml), "yaml"); err == nil {
		t.Error("unknown YAML key should error")
	}
}

func TestParseConfig_UnknownJSONKey(t *testing.T) {
	js := `{"unknownField": true}`
	if _, err := ParseConfig([]byte(js), "json"); err == nil {
		t.Error("unknown JSON key should error")
	}
}

// --- Validation ---

func TestValidation_NegativeTabWidth(t *testing.T) {
	tw := -1
	fc := &FileConfig{TabWidth: &tw}
	if err := validateFileConfig(fc); err == nil {
		t.Error("negative tabWidth should fail validation")
	}
}

func TestValidation_NegativeMinContrast(t *testing.T) {
	mc := -1.0
	fc := &FileConfig{MinContrast: &mc}
	if err := validateFileConfig(fc); err == nil {
		t.Error("negative minContrast should fail validation")
	}
}

func TestValidation_ZeroLineThreshold(t *testing.T) {
	lt := 0
	fc := &FileConfig{Collapsible: &CollapsibleFileConfig{LineThreshold: &lt}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("zero lineThreshold should fail validation")
	}
}

func TestValidation_InvalidFrame(t *testing.T) {
	f := "invalid"
	fc := &FileConfig{Defaults: &BlockDefaultsFileConfig{Frame: &f}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("invalid frame should fail validation")
	}
}

func TestValidation_InvalidTerminalDotStyle(t *testing.T) {
	s := "invalid"
	fc := &FileConfig{TerminalDotStyle: &s}
	if err := validateFileConfig(fc); err == nil {
		t.Error("invalid terminalDotStyle should fail validation")
	}
}

func TestValidation_MissingSelectorForSelectorKind(t *testing.T) {
	fc := &FileConfig{DarkMode: &DarkModeFileConfig{Kind: "selector"}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("missing selector for kind=selector should fail validation")
	}
}

func TestValidation_ValidConfig(t *testing.T) {
	tw := 4
	mc := 5.5
	f := "code"
	s := "colored"
	fc := &FileConfig{
		TabWidth:         &tw,
		MinContrast:      &mc,
		TerminalDotStyle: &s,
		Defaults:         &BlockDefaultsFileConfig{Frame: &f},
		DarkMode:         &DarkModeFileConfig{Kind: "selector", Selector: ".dark"},
	}
	if err := validateFileConfig(fc); err != nil {
		t.Errorf("valid config should pass: %v", err)
	}
}

// --- Style overrides ---

func TestParseStyleOverrides_String(t *testing.T) {
	raw := map[string]any{"radius": "0.5rem"}
	result, err := parseStyleOverrides(raw)
	if err != nil {
		t.Fatal(err)
	}
	sv := result["--kz-radius"]
	if sv.Value != "0.5rem" {
		t.Errorf("value = %q, want 0.5rem", sv.Value)
	}
	if sv.IsThemed() {
		t.Error("string value should not be themed")
	}
}

func TestParseStyleOverrides_Array(t *testing.T) {
	raw := map[string]any{"shadow": []any{"none", "0 2px 8px"}}
	result, err := parseStyleOverrides(raw)
	if err != nil {
		t.Fatal(err)
	}
	sv := result["--kz-shadow"]
	if sv.Dark != "none" || sv.Light != "0 2px 8px" {
		t.Errorf("themed values wrong: dark=%q light=%q", sv.Dark, sv.Light)
	}
}

func TestParseStyleOverrides_InvalidArrayLength(t *testing.T) {
	raw := map[string]any{"shadow": []any{"one"}}
	if _, err := parseStyleOverrides(raw); err == nil {
		t.Error("single-element array should error")
	}
}

func TestParseStyleOverrides_Map(t *testing.T) {
	raw := map[string]any{
		"editor-bg": map[string]any{"light": "#ffffff", "dark": "#1e293b"},
	}
	result, err := parseStyleOverrides(raw)
	if err != nil {
		t.Fatal(err)
	}
	sv := result["--kz-editor-bg"]
	if sv.Light != "#ffffff" {
		t.Errorf("light = %q, want #ffffff", sv.Light)
	}
	if sv.Dark != "#1e293b" {
		t.Errorf("dark = %q, want #1e293b", sv.Dark)
	}
}

func TestParseStyleOverrides_MapDarkOnly(t *testing.T) {
	raw := map[string]any{
		"editor-bg": map[string]any{"dark": "#1e293b"},
	}
	result, err := parseStyleOverrides(raw)
	if err != nil {
		t.Fatal(err)
	}
	sv := result["--kz-editor-bg"]
	if sv.Dark != "#1e293b" {
		t.Errorf("dark = %q, want #1e293b", sv.Dark)
	}
}

func TestParseStyleOverrides_MapEmpty(t *testing.T) {
	raw := map[string]any{
		"editor-bg": map[string]any{},
	}
	if _, err := parseStyleOverrides(raw); err == nil {
		t.Error("empty map should error")
	}
}

func TestParseStyleOverrides_InvalidType(t *testing.T) {
	raw := map[string]any{"shadow": 42}
	if _, err := parseStyleOverrides(raw); err == nil {
		t.Error("numeric value should error")
	}
}

func TestParseStyleOverrides_KeyNormalization(t *testing.T) {
	raw := map[string]any{
		"radius":       "1rem",
		"--kz-shadow":  "none",
		"--custom-var": "blue",
	}
	result, err := parseStyleOverrides(raw)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["--kz-radius"]; !ok {
		t.Error("bare name should be prefixed")
	}
	if _, ok := result["--kz-shadow"]; !ok {
		t.Error("already-prefixed should remain")
	}
	if _, ok := result["--custom-var"]; !ok {
		t.Error("custom prefix should remain")
	}
}

// --- FileConfigToOptions ---

func TestFileConfigToOptions_Booleans(t *testing.T) {
	f := false
	fc := &FileConfig{CopyButton: &f}
	opts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)
	if engine.cfg.CopyButton {
		t.Error("CopyButton should be false from file config")
	}
}

func TestFileConfigToOptions_PartialDefaults(t *testing.T) {
	wrap := true
	fc := &FileConfig{Defaults: &BlockDefaultsFileConfig{Wrap: &wrap}}
	opts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)
	if !engine.cfg.Defaults.Wrap {
		t.Error("Defaults.Wrap should be true from file config")
	}
	if !engine.cfg.Defaults.PreserveIndent {
		t.Error("Defaults.PreserveIndent should keep engine default (true)")
	}
}

func TestFileConfigToOptions_CommaLanguageKeys(t *testing.T) {
	wrap := true
	fc := &FileConfig{
		LanguageDefaults: map[string]BlockDefaultsFileConfig{
			"bash, sh, zsh": {Wrap: &wrap},
		},
	}
	opts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)

	for _, lang := range []string{"bash", "sh", "zsh"} {
		ld, ok := engine.cfg.LanguageDefaults[lang]
		if !ok {
			t.Errorf("language default for %q should exist", lang)
			continue
		}
		if !ld.Wrap {
			t.Errorf("language default for %q should have Wrap=true", lang)
		}
	}
}

func TestFileConfigToOptions_PartialCollapsible(t *testing.T) {
	lt := 25
	fc := &FileConfig{Collapsible: &CollapsibleFileConfig{LineThreshold: &lt}}
	opts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)

	if engine.cfg.Collapsible == nil {
		t.Fatal("Collapsible should be initialized")
	}
	if engine.cfg.Collapsible.LineThreshold != 25 {
		t.Errorf("LineThreshold = %d, want 25", engine.cfg.Collapsible.LineThreshold)
	}
	if engine.cfg.Collapsible.PreviewLines != 5 {
		t.Error("PreviewLines should keep default (5)")
	}
}

func TestFileConfigToOptions_StyleOverrides(t *testing.T) {
	fc := &FileConfig{
		StyleOverrides: map[string]any{
			"radius": "1rem",
			"shadow": []any{"none", "0 2px 8px rgba(0,0,0,0.1)"},
		},
	}
	opts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)

	if engine.cfg.StyleOverrides == nil {
		t.Fatal("StyleOverrides should be set")
	}
	if sv, ok := engine.cfg.StyleOverrides["--kz-radius"]; !ok || sv.Value != "1rem" {
		t.Error("radius override missing or wrong")
	}
	if sv, ok := engine.cfg.StyleOverrides["--kz-shadow"]; !ok || sv.Dark != "none" || sv.Light != "0 2px 8px rgba(0,0,0,0.1)" {
		t.Error("shadow themed override missing or wrong")
	}
}

func TestFileConfigToOptions_DarkMode(t *testing.T) {
	fc := &FileConfig{DarkMode: &DarkModeFileConfig{Kind: "mediaQuery"}}
	opts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)

	if engine.cfg.DarkMode.Kind != config.DarkModeMediaQueryKind {
		t.Error("DarkMode should be mediaQuery")
	}
}

// --- LoadConfig ---

func TestLoadConfig_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kazari.config.yaml")
	content := `
copyButton: false
tabWidth: 4
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	opts, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)

	if engine.cfg.CopyButton {
		t.Error("CopyButton should be false")
	}
	if engine.cfg.TabWidth != 4 {
		t.Errorf("TabWidth = %d, want 4", engine.cfg.TabWidth)
	}
}

func TestLoadConfig_UnsupportedExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	os.WriteFile(path, []byte(""), 0644)
	if _, err := LoadConfig(path); err == nil {
		t.Error("unsupported extension should error")
	}
}

// --- splitCommaKey ---

func TestSplitCommaKey(t *testing.T) {
	tests := map[string][]string{
		"bash":            {"bash"},
		"bash,sh,zsh":     {"bash", "sh", "zsh"},
		"bash, sh , zsh":  {"bash", "sh", "zsh"},
		"":                {},
	}
	for input, want := range tests {
		got := splitCommaKey(input)
		if len(got) != len(want) {
			t.Errorf("splitCommaKey(%q) = %v, want %v", input, got, want)
			continue
		}
		for i, g := range got {
			if g != want[i] {
				t.Errorf("splitCommaKey(%q)[%d] = %q, want %q", input, i, g, want[i])
			}
		}
	}
}

// --- KnownVarNames ---

func TestKnownVarNames(t *testing.T) {
	cfg := config.DefaultConfig()
	names := knownVarNamesFromConfig(cfg)

	if len(names) == 0 {
		t.Fatal("should return var names")
	}

	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}

	required := []string{"--kz-radius", "--kz-editor-bg", "--kz-editor-fg"}
	for _, r := range required {
		if !nameSet[r] {
			t.Errorf("should contain %q", r)
		}
	}

	for _, n := range names {
		if !strings.HasPrefix(n, "--kz-") {
			t.Errorf("unexpected var name without --kz- prefix: %q", n)
		}
	}
}

func knownVarNamesFromConfig(cfg *config.Config) []string {
	return theme.KnownVarNames(cfg)
}

// --- Composition tests (file config + programmatic options) ---

func compositionEngine(fileOpts []Option, programmatic ...Option) *Engine {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	base := []Option{
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	}
	all := append(base, fileOpts...)
	all = append(all, programmatic...)
	return New(all...)
}

func TestComposition_ScalarOverride(t *testing.T) {
	tw := 2
	fc := &FileConfig{TabWidth: &tw}
	fileOpts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	engine := compositionEngine(fileOpts, WithTabWidth(4))

	if engine.cfg.TabWidth != 4 {
		t.Errorf("TabWidth = %d, want 4 (programmatic should win)", engine.cfg.TabWidth)
	}
}

func TestComposition_StyleOverrideMerge(t *testing.T) {
	fc := &FileConfig{
		StyleOverrides: map[string]any{"radius": "0.5rem"},
	}
	fileOpts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	engine := compositionEngine(fileOpts,
		WithStyleOverrides(map[string]string{"shadow": "none"}),
	)

	if _, ok := engine.cfg.StyleOverrides["--kz-radius"]; !ok {
		t.Error("file config's radius override should be present")
	}
	if _, ok := engine.cfg.StyleOverrides["--kz-shadow"]; !ok {
		t.Error("programmatic shadow override should be present")
	}
}

func TestComposition_LanguageDefaultsMerge(t *testing.T) {
	ln := true
	fc := &FileConfig{
		LanguageDefaults: map[string]BlockDefaultsFileConfig{
			"go": {LineNumbers: &ln},
		},
	}
	fileOpts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	engine := compositionEngine(fileOpts,
		WithLanguageDefaults(map[string]BlockDefaults{
			"python": {LineNumbers: true},
		}),
	)

	if _, ok := engine.cfg.LanguageDefaults["go"]; !ok {
		t.Error("file config's 'go' language default should be present")
	}
	if _, ok := engine.cfg.LanguageDefaults["python"]; !ok {
		t.Error("programmatic 'python' language default should be present")
	}
}

func TestComposition_LanguageAliasesMerge(t *testing.T) {
	fc := &FileConfig{
		LanguageAliases: map[string]string{"js": "javascript"},
	}
	fileOpts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	engine := compositionEngine(fileOpts,
		WithLanguageAliases(map[string]string{"py": "python"}),
	)

	if engine.cfg.LanguageAliases["js"] != "javascript" {
		t.Error("file config's 'js' alias should be present")
	}
	if engine.cfg.LanguageAliases["py"] != "python" {
		t.Error("programmatic 'py' alias should be present")
	}
}

func TestComposition_UIStringsMerge(t *testing.T) {
	fc := &FileConfig{
		UIStrings: map[string]string{"copy.label": "Copy"},
	}
	fileOpts, err := FileConfigToOptions(fc)
	if err != nil {
		t.Fatal(err)
	}

	engine := compositionEngine(fileOpts,
		WithUIStrings(map[string]string{"expand.label": "Expand"}),
	)

	if engine.cfg.UIStringOverrides["copy.label"] != "Copy" {
		t.Error("file config's copy.label should be present")
	}
	if engine.cfg.UIStringOverrides["expand.label"] != "Expand" {
		t.Error("programmatic expand.label should be present")
	}
}

// --- WithConfigDir ---

func TestWithConfigDir_LoadsYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "kazari.config.yaml"), []byte("tabWidth: 4\n"), 0644)

	engine := compositionEngine(nil, WithConfigDir(dir))

	if engine.cfg.TabWidth != 4 {
		t.Errorf("TabWidth = %d, want 4", engine.cfg.TabWidth)
	}
}

func TestWithConfigDir_NoFile(t *testing.T) {
	dir := t.TempDir()

	engine := compositionEngine(nil, WithConfigDir(dir))

	if engine.cfg.TabWidth != 2 {
		t.Errorf("TabWidth = %d, want default 2", engine.cfg.TabWidth)
	}
}

func TestWithConfigDir_ProgrammaticOverrides(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "kazari.config.yaml"), []byte("tabWidth: 2\n"), 0644)

	engine := compositionEngine(nil, WithConfigDir(dir), WithTabWidth(4))

	if engine.cfg.TabWidth != 4 {
		t.Errorf("TabWidth = %d, want 4 (programmatic should win)", engine.cfg.TabWidth)
	}
}

func TestWithConfigDir_Priority(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "kazari.config.yaml"), []byte("tabWidth: 4\n"), 0644)
	os.WriteFile(filepath.Join(dir, "kazari.config.json"), []byte(`{"tabWidth": 8}`), 0644)

	engine := compositionEngine(nil, WithConfigDir(dir))

	if engine.cfg.TabWidth != 4 {
		t.Errorf("TabWidth = %d, want 4 (.yaml should win over .json)", engine.cfg.TabWidth)
	}
}

// --- A. Validation rules with no test coverage ---

func TestValidation_ZeroPreviewLines(t *testing.T) {
	pl := 0
	fc := &FileConfig{Collapsible: &CollapsibleFileConfig{PreviewLines: &pl}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("zero previewLines should fail validation")
	}
}

func TestValidation_InvalidCollapsibleStyle(t *testing.T) {
	s := "invalid"
	fc := &FileConfig{Collapsible: &CollapsibleFileConfig{Style: &s}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("invalid collapsible style should fail validation")
	}
}

func TestValidation_InvalidLanguageDefaultsFrame(t *testing.T) {
	f := "invalid"
	fc := &FileConfig{
		LanguageDefaults: map[string]BlockDefaultsFileConfig{
			"go": {Frame: &f},
		},
	}
	if err := validateFileConfig(fc); err == nil {
		t.Error("invalid languageDefaults frame should fail validation")
	}
}

// --- B. Validation rules with partial coverage ---

func TestValidation_MissingSelectorForBothKind(t *testing.T) {
	fc := &FileConfig{DarkMode: &DarkModeFileConfig{Kind: "both"}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("missing selector for kind=both should fail validation")
	}
}

func TestValidation_ZeroTabWidth(t *testing.T) {
	tw := 0
	fc := &FileConfig{TabWidth: &tw}
	if err := validateFileConfig(fc); err == nil {
		t.Error("zero tabWidth should fail validation")
	}
}

func TestValidation_InvalidLangIconMode(t *testing.T) {
	s := "invalid"
	fc := &FileConfig{LanguageIconMode: &s}
	if err := validateFileConfig(fc); err == nil {
		t.Error("invalid languageIconMode should fail validation")
	}
}

func TestValidation_InvalidDarkModeKind(t *testing.T) {
	fc := &FileConfig{DarkMode: &DarkModeFileConfig{Kind: "invalid", Selector: ".dark"}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("invalid darkMode.kind should fail validation")
	}
}

// --- C. Boundary valid values ---

func TestValidation_BoundaryValidValues(t *testing.T) {
	tw := 1
	mc := 0.0
	lt := 1
	pl := 1
	fc := &FileConfig{
		TabWidth:    &tw,
		MinContrast: &mc,
		Collapsible: &CollapsibleFileConfig{
			LineThreshold: &lt,
			PreviewLines:  &pl,
		},
	}
	if err := validateFileConfig(fc); err != nil {
		t.Errorf("minimum valid values should pass: %v", err)
	}
}

// --- D. ParseConfig error paths ---

func TestParseConfig_MalformedYAML(t *testing.T) {
	if _, err := ParseConfig([]byte("tabWidth: [bad"), "yaml"); err == nil {
		t.Error("malformed YAML should error")
	}
}

func TestParseConfig_MalformedJSON(t *testing.T) {
	if _, err := ParseConfig([]byte("{bad json}"), "json"); err == nil {
		t.Error("malformed JSON should error")
	}
}

func TestParseConfig_WrongTypeInYAML(t *testing.T) {
	if _, err := ParseConfig([]byte("tabWidth: abc"), "yaml"); err == nil {
		t.Error("wrong type for tabWidth should error")
	}
}

// --- E. Style override error paths ---

func TestParseStyleOverrides_ArrayNonStringElements(t *testing.T) {
	raw := map[string]any{"shadow": []any{42, "red"}}
	if _, err := parseStyleOverrides(raw); err == nil {
		t.Error("array with non-string elements should error")
	}
}

func TestParseStyleOverrides_MapNonStringValue(t *testing.T) {
	raw := map[string]any{
		"editor-bg": map[string]any{"light": 42},
	}
	if _, err := parseStyleOverrides(raw); err == nil {
		t.Error("map with non-string value should error")
	}
}

// --- F. LoadConfig edge cases ---

func TestLoadConfig_YMLExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kazari.config.yml")
	if err := os.WriteFile(path, []byte("tabWidth: 4\n"), 0644); err != nil {
		t.Fatal(err)
	}

	opts, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, opts...)

	if engine.cfg.TabWidth != 4 {
		t.Errorf("TabWidth = %d, want 4", engine.cfg.TabWidth)
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	if _, err := LoadConfig("/nonexistent/path/kazari.config.yaml"); err == nil {
		t.Error("non-existent file should error")
	}
}

// --- G. hangingIndent validation ---

func TestValidation_NegativeHangingIndent(t *testing.T) {
	hi := -1
	fc := &FileConfig{Defaults: &BlockDefaultsFileConfig{HangingIndent: &hi}}
	if err := validateFileConfig(fc); err == nil {
		t.Error("negative hangingIndent should fail validation")
	}
}

func TestValidation_ZeroHangingIndent(t *testing.T) {
	hi := 0
	fc := &FileConfig{Defaults: &BlockDefaultsFileConfig{HangingIndent: &hi}}
	if err := validateFileConfig(fc); err != nil {
		t.Errorf("zero hangingIndent should pass: %v", err)
	}
}

func TestValidation_NegativeLanguageDefaultsHangingIndent(t *testing.T) {
	hi := -1
	fc := &FileConfig{
		LanguageDefaults: map[string]BlockDefaultsFileConfig{
			"go": {HangingIndent: &hi},
		},
	}
	if err := validateFileConfig(fc); err == nil {
		t.Error("negative languageDefaults hangingIndent should fail validation")
	}
}

// --- H. minContrast upper bound ---

func TestValidation_MinContrastAboveMax(t *testing.T) {
	mc := 22.0
	fc := &FileConfig{MinContrast: &mc}
	if err := validateFileConfig(fc); err == nil {
		t.Error("minContrast above 21 should fail validation")
	}
}

func TestValidation_MinContrastAtMax(t *testing.T) {
	mc := 21.0
	fc := &FileConfig{MinContrast: &mc}
	if err := validateFileConfig(fc); err != nil {
		t.Errorf("minContrast=21 should pass: %v", err)
	}
}
