package render

import (
	"strings"
	"testing"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/locale"
)

func defaultCfg() *config.Config {
	cfg := config.DefaultConfig()
	cfg.UIStrings = locale.Resolve("en-US", nil)
	return cfg
}

func simpleLine(content, lightColor string) TokenLine {
	return TokenLine{Tokens: []MergedToken{{Content: content, LightColor: lightColor}}}
}

func dualLine(content, lightColor, darkColor string) TokenLine {
	return TokenLine{Tokens: []MergedToken{{Content: content, LightColor: lightColor, DarkColor: darkColor}}}
}

func resolved() *config.ResolvedBlock {
	return &config.ResolvedBlock{
		Lang:            "go",
		Frame:           config.FrameCode,
		StartLineNumber: 1,
		RawCode:         "hello",
	}
}

// --- Frame tests ---

func TestRenderBlock_EditorFrame(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.Title = "main.go"
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	for _, want := range []string{
		`class="frame has-title"`,
		"kz-toolbar",
		`kz-title`,
		"main.go",
		"kz-lang",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in editor frame output", want)
		}
	}
}

func TestRenderBlock_EditorFrame_NoTitle(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `class="frame"`) {
		t.Error("expected frame class without has-title")
	}
	if strings.Contains(out, "has-title") {
		t.Error("should not have has-title when no title set")
	}
}

func TestRenderBlock_TerminalFrame(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.Frame = config.FrameTerminal
	r.Title = "Terminal"
	lines := []TokenLine{simpleLine("echo hi", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	for _, want := range []string{
		"is-terminal",
		"kz-terminal-header",
		"kz-terminal-dots",
		"Terminal",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in terminal frame output", want)
		}
	}
}

func TestRenderBlock_TerminalFrame_MinimalDots(t *testing.T) {
	cfg := defaultCfg()
	cfg.TerminalDotStyle = config.DotsMinimal
	r := resolved()
	r.Frame = config.FrameTerminal
	lines := []TokenLine{simpleLine("echo hi", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "kz-dots-minimal") {
		t.Error("expected kz-dots-minimal class")
	}
	if strings.Contains(out, "kz-terminal-dots") {
		t.Error("should not have colored dot spans in minimal mode")
	}
}

func TestRenderBlock_TerminalFrame_SrOnlyLabel(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.Frame = config.FrameTerminal
	r.Title = ""
	lines := []TokenLine{simpleLine("echo hi", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "sr-only") {
		t.Error("expected sr-only label for untitled terminal frame")
	}
}

func TestRenderBlock_FrameNone(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.Frame = config.FrameNone
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if strings.Contains(out, "<figure") {
		t.Error("frame none should not have <figure>")
	}
	if !strings.Contains(out, "<pre") {
		t.Error("frame none should have <pre>")
	}
}

// --- Line numbers ---

func TestRenderBlock_LineNumbers(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.LineNumbers = true
	lines := []TokenLine{simpleLine("line1", "#aaa"), simpleLine("line2", "#bbb")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `class="gutter"`) {
		t.Error("expected gutter div")
	}
	if !strings.Contains(out, `aria-hidden="true"`) {
		t.Error("expected aria-hidden on line numbers")
	}
	if !strings.Contains(out, ">1</div>") {
		t.Error("expected line number 1")
	}
	if !strings.Contains(out, ">2</div>") {
		t.Error("expected line number 2")
	}
}

func TestRenderBlock_LineNumbers_CustomStart(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.LineNumbers = true
	r.StartLineNumber = 10
	lines := []TokenLine{simpleLine("line", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, ">10</div>") {
		t.Error("expected line number starting at 10")
	}
}

func TestRenderBlock_LineNumbers_DynamicWidth(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.LineNumbers = true
	r.StartLineNumber = 1
	lines := make([]TokenLine, 100)
	for i := range lines {
		lines[i] = simpleLine("x", "#aaa")
	}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "--kz-ln-width:3ch") {
		t.Error("expected dynamic width 3ch for 100 lines")
	}
}

// --- Token styling ---

func TestRenderBlock_TokenColors(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{dualLine("hello", "#aaa", "#bbb")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "--sl:#aaa") {
		t.Error("expected --sl light color")
	}
	if !strings.Contains(out, "--sd:#bbb") {
		t.Error("expected --sd dark color")
	}
}

func TestRenderBlock_TokenBackground(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{{Tokens: []MergedToken{{
		Content: "hi", LightColor: "#aaa", DarkColor: "#bbb",
		LightBG: "#eee", DarkBG: "#222",
	}}}}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "--slbg:#eee") {
		t.Error("expected --slbg light background")
	}
	if !strings.Contains(out, "--sdbg:#222") {
		t.Error("expected --sdbg dark background")
	}
}

func TestRenderBlock_FontStyle(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{{Tokens: []MergedToken{{
		Content: "hi", LightColor: "#aaa", FontStyle: 1 | 2,
	}}}}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "--sfs:italic") {
		t.Error("expected --sfs:italic")
	}
	if !strings.Contains(out, "--sfw:bold") {
		t.Error("expected --sfw:bold")
	}
}

// --- Markers ---

func TestRenderBlock_LineMarkers(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.LineMarkers = []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 1, End: 1}}},
		{Type: config.MarkerIns, Lines: []config.LineRange{{Start: 2, End: 2}}},
		{Type: config.MarkerDel, Lines: []config.LineRange{{Start: 3, End: 3}}},
	}
	lines := []TokenLine{simpleLine("a", "#aaa"), simpleLine("b", "#aaa"), simpleLine("c", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "highlight mark") {
		t.Error("expected mark class")
	}
	if !strings.Contains(out, "highlight ins") {
		t.Error("expected ins class")
	}
	if !strings.Contains(out, "highlight del") {
		t.Error("expected del class")
	}
}

func TestRenderBlock_LabeledRange(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.LineMarkers = []config.LineMarker{
		{Type: config.MarkerIns, Lines: []config.LineRange{{Start: 1, End: 1}}, Label: "Added"},
	}
	lines := []TokenLine{simpleLine("x", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "tm-label") {
		t.Error("expected tm-label class")
	}
	if !strings.Contains(out, `data-label="Added"`) {
		t.Error("expected data-label attribute")
	}
}

// --- Focus ---

func TestRenderBlock_FocusLines(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.FocusLines = []config.LineRange{{Start: 2, End: 2}}
	lines := []TokenLine{simpleLine("a", "#aaa"), simpleLine("b", "#aaa"), simpleLine("c", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "has-focus") {
		t.Error("expected has-focus class on code element")
	}
	if !strings.Contains(out, "focused") {
		t.Error("expected focused class on focused line")
	}
}

// --- Copy and fullscreen buttons ---

func TestRenderBlock_CopyButton(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.RawCode = "hello\nworld"
	lines := []TokenLine{simpleLine("hello", "#aaa"), simpleLine("world", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "kz-copy-btn") {
		t.Error("expected copy button")
	}
	if !strings.Contains(out, "data-code=") {
		t.Error("expected data-code attribute")
	}
}

func TestRenderBlock_FullscreenButton(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "kz-fs-btn") {
		t.Error("expected fullscreen button")
	}
	if !strings.Contains(out, "kz-fs-exit-icon") {
		t.Error("expected exit fullscreen icon for icon switching")
	}
	if !strings.Contains(out, `aria-expanded="false"`) {
		t.Error("expected aria-expanded on fullscreen button")
	}
}

func TestRenderBlock_FullscreenFontControls(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	for _, want := range []string{"kz-font-controls", "kz-font-inc", "kz-font-dec", "kz-fs-hint"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestRenderBlock_NoFullscreenNoFontControls(t *testing.T) {
	cfg := defaultCfg()
	cfg.FullscreenButton = false
	r := resolved()
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	for _, absent := range []string{"kz-font-controls", "kz-font-inc", "kz-font-dec", "kz-fs-hint", "kz-fs-btn"} {
		if strings.Contains(out, absent) {
			t.Errorf("should not contain %q when fullscreen disabled", absent)
		}
	}
}

func TestRenderBlock_TerminalFullscreenFontControls(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.Frame = config.FrameTerminal
	lines := []TokenLine{simpleLine("$ echo hi", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	for _, want := range []string{"kz-font-controls", "kz-font-inc", "kz-font-dec", "kz-fs-btn", "kz-fs-hint"} {
		if !strings.Contains(out, want) {
			t.Errorf("terminal frame: expected %q in output", want)
		}
	}
}

// --- Wrap ---

func TestRenderBlock_Wrap(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.Wrap = true
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `class="wrap"`) {
		t.Error("expected wrap class on pre")
	}
}

// --- Data-lines ---

func TestRenderBlock_DataLines(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{simpleLine("a", "#aaa"), simpleLine("b", "#aaa"), simpleLine("c", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `data-lines="3"`) {
		t.Error("expected data-lines attribute with count 3")
	}
}

func TestRenderBlock_DataLines_Disabled(t *testing.T) {
	cfg := defaultCfg()
	cfg.DataLineCount = false
	r := resolved()
	lines := []TokenLine{simpleLine("a", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if strings.Contains(out, "data-lines") {
		t.Error("data-lines should not be present when disabled")
	}
}

// --- File icons ---

func TestRenderBlock_FileIcon(t *testing.T) {
	cfg := defaultCfg()
	cfg.FileIcons = true
	r := resolved()
	r.Title = "main.go"
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `data-ext="go"`) {
		t.Error("expected file icon with data-ext=go")
	}
}

func TestRenderBlock_FileIcon_CustomResolver(t *testing.T) {
	cfg := defaultCfg()
	cfg.FileIcons = true
	cfg.FileIconResolver = func(ext string) string {
		return "<img src=\"" + ext + ".svg\">"
	}
	r := resolved()
	r.Title = "app.py"
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "py.svg") {
		t.Error("expected custom resolver output")
	}
}

// --- Collapsible ---

func TestRenderBlock_CollapseThreshold(t *testing.T) {
	cfg := defaultCfg()
	cfg.Collapsible = &config.CollapsibleConfig{
		DefaultCollapsed: true,
	}
	r := resolved()
	r.CollapseThreshold = true
	r.CollapseSegments = []config.PreviewSegment{{Start: 1, End: 3}}
	lines := make([]TokenLine, 20)
	for i := range lines {
		lines[i] = simpleLine("x", "#aaa")
	}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "kz-collapsed") {
		t.Error("expected kz-collapsed class on wrapper")
	}
	if !strings.Contains(out, "kz-collapse-gradient") {
		t.Error("expected gradient overlay")
	}
	if !strings.Contains(out, "kz-collapse-btn") {
		t.Error("expected collapse button")
	}
	if !strings.Contains(out, "kz-hidden") {
		t.Error("expected hidden lines")
	}
}

func TestRenderBlock_CollapseRange_GitHub(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.CollapseRanges = []config.CollapseRange{
		{Start: 2, End: 3, LineCount: 2, Style: config.CollapseGithub},
	}
	lines := []TokenLine{simpleLine("a", "#aaa"), simpleLine("b", "#aaa"), simpleLine("c", "#aaa"), simpleLine("d", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, "<details class=\"kz-section\"") {
		t.Error("expected <details> for github collapse style")
	}
	if !strings.Contains(out, "<summary>") {
		t.Error("expected <summary> element")
	}
}

// --- Helpers ---

func TestDisplayLang(t *testing.T) {
	cases := []struct{ in, want string }{
		{"javascript", "JavaScript"},
		{"typescript", "TypeScript"},
		{"css", "CSS"},
		{"go", "Go"},
		{"python", "Python"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := displayLang(tc.in); got != tc.want {
			t.Errorf("displayLang(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFileExt(t *testing.T) {
	cases := []struct{ in, want string }{
		{"main.go", "go"},
		{"app.config.js", "js"},
		{"Makefile", ""},
		{"", ""},
		{"file.", ""},
	}
	for _, tc := range cases {
		if got := fileExt(tc.in); got != tc.want {
			t.Errorf("fileExt(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDigitCount(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{0, 1},
		{1, 1},
		{9, 1},
		{10, 2},
		{99, 2},
		{100, 3},
		{-1, 2},
		{-10, 3},
	}
	for _, tc := range cases {
		if got := digitCount(tc.n); got != tc.want {
			t.Errorf("digitCount(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestEncodeForDataCode(t *testing.T) {
	got := encodeForDataCode("line1\nline2\nline3")
	if strings.Contains(got, "\n") {
		t.Error("encoded output should not contain newlines")
	}
	if !strings.Contains(got, "\x7f") {
		t.Error("encoded output should contain DEL characters")
	}
}

// --- Per-block theme override tests ---

func TestRenderBlock_ThemeOverrideWrapper(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	r.ThemeOverrideStyle = "--kz-ovl-editor-bg:#282a36;--kz-ovl-editor-fg:#f8f8f2"
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `class="kazari-code kz-themed"`) {
		t.Error("wrapper should carry the kz-themed class")
	}
	if !strings.Contains(out, `style="--kz-ovl-editor-bg:#282a36;--kz-ovl-editor-fg:#f8f8f2"`) {
		t.Error("wrapper should carry the inline override style")
	}
}

func TestRenderBlock_ThemeOverrideWrapper_ComposesWithDataLinesAndCollapse(t *testing.T) {
	cfg := defaultCfg()
	cfg.DataLineCount = true
	r := resolved()
	r.ThemeOverrideStyle = "--kz-ovl-editor-bg:#282a36"
	r.CollapseThreshold = true
	r.CollapseConfig = &config.CollapsibleConfig{DefaultCollapsed: true}
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if !strings.Contains(out, `class="kazari-code kz-themed kz-collapsed"`) {
		t.Errorf("wrapper classes should compose, got output start %q", out[:120])
	}
	if !strings.Contains(out, `data-lines="1"`) {
		t.Error("data-lines should still be emitted")
	}
	if !strings.Contains(out, `style="--kz-ovl-editor-bg:#282a36"`) {
		t.Error("style attribute should still be emitted")
	}
}

func TestRenderBlock_NoThemeOverride_NoThemedClass(t *testing.T) {
	cfg := defaultCfg()
	r := resolved()
	lines := []TokenLine{simpleLine("hello", "#aaa")}

	out := RenderBlock(lines, r, cfg)

	if strings.Contains(out, "kz-themed") {
		t.Error("wrapper should not carry kz-themed without an override")
	}
	if strings.Contains(out, `<div class="kazari-code" style=`) {
		t.Error("wrapper should not carry a style attribute without an override")
	}
}

func TestBuildTokenStyle_PrefersResolvedMarkerBGs(t *testing.T) {
	cfg := defaultCfg()
	cfg.MinContrast = 5.5
	cfg.LightMarkerBGs = &config.MarkerEffectiveBGs{Mark: "#ffffff", Ins: "#ffffff", Del: "#ffffff"}

	baseResolved := resolved()
	overrideResolved := resolved()
	overrideResolved.LightMarkerBGs = &config.MarkerEffectiveBGs{Mark: "#000000", Ins: "#000000", Del: "#000000"}

	mt := config.MarkerMark
	tok := MergedToken{Content: "x", LightColor: "#888888"}

	pageStyle := buildTokenStyle(tok, &lineContext{
		resolved: baseResolved, cfg: cfg, contrastCache: map[string]string{},
	}, &mt)
	overrideStyle := buildTokenStyle(tok, &lineContext{
		resolved: overrideResolved, cfg: cfg, contrastCache: map[string]string{},
	}, &mt)

	if pageStyle == overrideStyle {
		t.Errorf("contrast adjustment should differ between page and override marker backgrounds, both produced %q", pageStyle)
	}
}
