package kazari

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/frostybee/kazari/internal/color"
)

type mockHighlighter struct {
	lightTokens   [][]Token
	darkTokens    [][]Token
	themeInfo     ThemeInfo
	darkThemeInfo ThemeInfo
	lightThemeName string
}

func (m *mockHighlighter) Tokenize(code, lang, theme string) ([][]Token, error) {
	if m.darkTokens != nil && theme != m.lightTheme() {
		return m.darkTokens, nil
	}
	return m.lightTokens, nil
}

func (m *mockHighlighter) GetThemeColors(theme string) (ThemeInfo, error) {
	if m.darkThemeInfo.BG != "" && theme != m.lightTheme() {
		return m.darkThemeInfo, nil
	}
	return m.themeInfo, nil
}

func (m *mockHighlighter) lightTheme() string {
	if m.lightThemeName != "" {
		return m.lightThemeName
	}
	return "light-theme"
}

func (m *mockHighlighter) GetLoadedLanguages() []string {
	return []string{"go", "javascript", "bash"}
}

// dualMockHighlighter wraps mockHighlighter with the DualThemeTokenizer
// capability and counts calls on both paths.
type dualMockHighlighter struct {
	mockHighlighter
	dualCalls     int
	tokenizeCalls int
	lastLight     string
	lastDark      string
}

func (d *dualMockHighlighter) Tokenize(code, lang, theme string) ([][]Token, error) {
	d.tokenizeCalls++
	return d.mockHighlighter.Tokenize(code, lang, theme)
}

func (d *dualMockHighlighter) TokenizeDual(code, lang, lightTheme, darkTheme string) ([][]Token, [][]Token, error) {
	d.dualCalls++
	d.lastLight = lightTheme
	d.lastDark = darkTheme
	return d.lightTokens, d.darkTokens, nil
}

func newTestEngine(hl *mockHighlighter, opts ...Option) *Engine {
	base := []Option{
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithMinContrast(0),
	}
	return New(append(base, opts...)...)
}

// --- Core rendering tests ---

func TestNew(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", "dark-theme"),
	)
	if engine == nil {
		t.Fatal("New() returned nil")
	}
}

func TestRender_DualTheme(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "func", Color: "#cf222e"}, {Content: " main", Color: "#8250df"}},
		},
		darkTokens: [][]Token{
			{{Content: "func", Color: "#ff7b72"}, {Content: " main", Color: "#d2a8ff"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"), WithMinify(false), WithMinContrast(0))
	html, err := engine.Render("func main", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `class="kazari-block"`) {
		t.Error("missing kazari-block wrapper")
	}
	if !strings.Contains(html, "--sl:#cf222e") {
		t.Error("missing light color")
	}
	if !strings.Contains(html, "--sd:#ff7b72") {
		t.Error("missing dark color")
	}
}

func TestRender_DualThemeCapability(t *testing.T) {
	hl := &dualMockHighlighter{
		mockHighlighter: mockHighlighter{
			lightTokens: [][]Token{
				{{Content: "func", Color: "#cf222e"}, {Content: " main", Color: "#8250df"}},
			},
			darkTokens: [][]Token{
				{{Content: "func", Color: "#ff7b72"}, {Content: " main", Color: "#d2a8ff"}},
			},
			themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		},
	}

	engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"), WithMinify(false), WithMinContrast(0))
	html, err := engine.Render("func main", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if hl.dualCalls != 1 {
		t.Errorf("TokenizeDual calls: got %d, want 1", hl.dualCalls)
	}
	if hl.tokenizeCalls != 0 {
		t.Errorf("Tokenize calls: got %d, want 0 (capability path must replace both passes)", hl.tokenizeCalls)
	}
	if !strings.Contains(html, "--sl:#cf222e") {
		t.Error("missing light color")
	}
	if !strings.Contains(html, "--sd:#ff7b72") {
		t.Error("missing dark color")
	}
}

func TestRender_SingleThemeSkipsCapability(t *testing.T) {
	hl := &dualMockHighlighter{
		mockHighlighter: mockHighlighter{
			lightTokens: [][]Token{{{Content: "hello", Color: "#333333"}}},
			themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		},
	}

	engine := New(WithHighlighter(hl), WithThemes("light-theme", ""), WithMinify(false))
	html, err := engine.Render("hello", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if hl.dualCalls != 0 {
		t.Errorf("TokenizeDual calls: got %d, want 0 (single-theme must not use the capability)", hl.dualCalls)
	}
	if hl.tokenizeCalls != 1 {
		t.Errorf("Tokenize calls: got %d, want 1", hl.tokenizeCalls)
	}
	if strings.Contains(html, "--sd:") {
		t.Error("single-theme should not emit --sd")
	}
}

func TestRender_SingleTheme(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "--sd:") {
		t.Error("single-theme should not emit --sd")
	}
}

func TestRender_HTMLEscape(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "<div>&\"test\"</div>", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render(`<div>&"test"</div>`, Options{Lang: "html"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "&lt;div&gt;&amp;&#34;test&#34;&lt;/div&gt;") {
		t.Errorf("incorrect HTML escaping, got: %s", html)
	}
}

func TestRender_EmptyLine(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{},
			{{Content: "b", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("a\n\nb", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	count := strings.Count(html, `<div class="kz-line">`)
	if count != 3 {
		t.Errorf("expected 3 kz-line divs, got %d", count)
	}
}

func TestRender_FontStyles(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "italic", Color: "#000", FontStyle: FontStyleItalic}},
			{{Content: "bold", Color: "#000", FontStyle: FontStyleBold}},
			{{Content: "both", Color: "#000", FontStyle: FontStyleItalic | FontStyleBold | FontStyleUnderline}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("italic\nbold\nboth", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "--sfs:italic") {
		t.Error("missing italic")
	}
	if !strings.Contains(html, "--sfw:bold") {
		t.Error("missing bold")
	}
	if !strings.Contains(html, "--std:underline") {
		t.Error("missing underline")
	}
}

func TestRender_NoHighlighter(t *testing.T) {
	engine := New(WithMinify(false))
	html, err := engine.Render("hello world", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "hello world") {
		t.Error("plaintext fallback should contain the code")
	}
}

// --- Toolbar + Frame tests ---

func TestRender_Toolbar(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `<div class="kz-toolbar">`) {
		t.Error("missing toolbar")
	}
	if !strings.Contains(html, `kz-toolbar-left`) {
		t.Error("missing toolbar left section")
	}
	if !strings.Contains(html, `kz-toolbar-right`) {
		t.Error("missing toolbar right section")
	}
}

func TestRender_LanguageBadge_NoTitle(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLanguageBadge(true))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	// Badge should be in left section when no title
	if !strings.Contains(html, `<span class="kz-lang">Go</span>`) {
		t.Error("missing language badge in left section")
	}
}

func TestRender_TitleLeft_BadgeRight(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLanguageBadge(true))
	html, err := engine.Render("x", Options{Lang: "go", Title: "main.go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `<span class="kz-title">main.go</span>`) {
		t.Error("missing title in left section")
	}
	if !strings.Contains(html, `kz-toolbar-right`) {
		t.Error("missing toolbar right")
	}
	if !strings.Contains(html, `<span class="kz-lang">Go</span>`) {
		t.Error("missing language badge")
	}
}

func TestRender_EditorFrame(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("x", Options{Lang: "go", Title: "main.go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `<figure class="frame has-title"`) {
		t.Error("missing editor frame with has-title")
	}
}

func TestRender_TerminalFrame(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "echo hi", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("echo hi", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "is-terminal") {
		t.Error("bash should produce terminal frame")
	}
	if !strings.Contains(html, "kz-terminal-header") {
		t.Error("terminal frame should have terminal header, not toolbar")
	}
	if !strings.Contains(html, "kz-terminal-dots") {
		t.Error("terminal frame should have dots")
	}
	if strings.Contains(html, "kz-toolbar") {
		t.Error("terminal frame should not have editor toolbar")
	}
}

func TestTerminalDots_DefaultColored(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "echo hi", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("echo hi", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "kz-terminal-dots") {
		t.Error("default dot style should render colored dot spans")
	}
	if strings.Contains(html, "kz-dots-minimal") {
		t.Error("default dot style should not have kz-dots-minimal class")
	}
}

func TestTerminalDots_Minimal(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "echo hi", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithTerminalDotStyle(DotsMinimal))
	html, err := engine.Render("echo hi", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "kz-dots-minimal") {
		t.Error("minimal dot style should add kz-dots-minimal class to header")
	}
	if strings.Contains(html, "kz-terminal-dots") {
		t.Error("minimal dot style should not render colored dot spans")
	}
}

func TestCSS_BothDotStyles(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#000", BG: "#fff"}}
	engine := newTestEngine(hl)
	css := engine.CSS()

	if !strings.Contains(css, "--kz-terminal-dot-red") {
		t.Error("CSS should contain colored dot rules")
	}
	if !strings.Contains(css, "kz-dots-minimal") {
		t.Error("CSS should contain minimal dot rules")
	}
	if !strings.Contains(css, "--kz-terminal-icon") {
		t.Error("CSS should contain minimal dot variables")
	}
}

func TestRender_TerminalFrame_WithShebang(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "#!/bin/bash", Color: "#000"}},
			{{Content: "echo hi", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("#!/bin/bash\necho hi", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "is-terminal") {
		t.Error("bash with shebang should produce editor frame")
	}
}

func TestRender_FrameNone(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	f := FrameNone
	html, err := engine.Render("x", Options{Lang: "go", Frame: &f})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "<figure") {
		t.Error("frame=none should not produce figure")
	}
	if strings.Contains(html, "kz-toolbar") {
		t.Error("frame=none should not have toolbar")
	}
}

func TestRender_ExplicitFrameOverride(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	f := FrameTerminal
	html, err := engine.Render("x", Options{Lang: "go", Frame: &f})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "is-terminal") {
		t.Error("explicit frame=terminal should override auto-detection")
	}
}

// --- Copy button tests ---

func TestRender_CopyButton(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithCopyButton(true))
	html, err := engine.Render("line1\nline2", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "kz-copy-btn") {
		t.Error("missing copy button")
	}
	if !strings.Contains(html, "data-code=") {
		t.Error("missing data-code attribute")
	}
	if !strings.Contains(html, "<svg") {
		t.Error("missing copy SVG icon")
	}
	if !strings.Contains(html, `data-tooltip="Copy"`) {
		t.Error("missing copy tooltip attribute")
	}
}

func TestRender_NoCopyButton(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithCopyButton(false))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "kz-copy-btn") {
		t.Error("copy button should not be present when disabled")
	}
}

func TestRender_FullscreenButton(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithFullscreenButton(true))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "kz-fs-btn") {
		t.Error("missing fullscreen button")
	}
}

// --- Meta tests ---

func TestRenderWithMeta(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("x", `go title="main.go"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `data-language="go"`) {
		t.Error("meta string language not parsed")
	}
	if !strings.Contains(html, "main.go") {
		t.Error("meta string title not rendered")
	}
}

func TestRenderWithMeta_FrameOverride(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("x", `go frame="none"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if strings.Contains(html, "<figure") {
		t.Error("frame=none via meta should not produce figure")
	}
}

// --- CSS/JS tests ---

func TestCSS_Deterministic(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"), WithMinify(false))

	css1 := engine.CSS()
	css2 := engine.CSS()
	if css1 != css2 {
		t.Error("CSS() should be deterministic")
	}
	if css1 == "" {
		t.Error("CSS() should not be empty")
	}
}

func TestCSS_ContainsLayer(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCascadeLayer("kazari"))

	if !strings.Contains(engine.CSS(), "@layer kazari") {
		t.Error("CSS should contain @layer wrapper")
	}
}

func TestCSS_NoLayer(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCascadeLayer(""))

	if strings.Contains(engine.CSS(), "@layer") {
		t.Error("CSS should not contain @layer when disabled")
	}
}

func TestCSS_ContainsToolbarStyles(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false))

	css := engine.CSS()
	if !strings.Contains(css, ".kz-toolbar") {
		t.Error("CSS should contain toolbar styles")
	}
	if !strings.Contains(css, ".kz-lang") {
		t.Error("CSS should contain language badge styles")
	}
}

func TestCSS_ContainsCopyStyles(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCopyButton(true))

	if !strings.Contains(engine.CSS(), ".kz-copy-btn") {
		t.Error("CSS should contain copy button styles")
	}
}

func TestThemeCSS_ContainsVars(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false))

	css := engine.ThemeCSS()
	if !strings.Contains(css, "--kz-editor-bg") {
		t.Error("ThemeCSS should contain theme variables")
	}
	if !strings.Contains(css, ":root") {
		t.Error("ThemeCSS should contain root selector")
	}
}

func TestThemeCSS_NoStructuralCSS(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false))

	css := engine.ThemeCSS()
	if strings.Contains(css, ".kz-toolbar") {
		t.Error("ThemeCSS should not contain structural CSS")
	}
	if strings.Contains(css, ".kz-copy-btn") {
		t.Error("ThemeCSS should not contain structural CSS")
	}
	if strings.Contains(css, "grid-template") {
		t.Error("ThemeCSS should not contain layout CSS")
	}
}

func TestThemeCSS_CustomRoot(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithThemeCSSRoot(".my-scope"))

	css := engine.ThemeCSS()
	if !strings.Contains(css, ".my-scope") {
		t.Error("ThemeCSS should use custom CSS root")
	}
	if strings.Contains(css, ".kz-toolbar") {
		t.Error("ThemeCSS should not contain structural CSS")
	}
}

func TestThemeCSS_CascadeLayer(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCascadeLayer("kazari"))

	if !strings.Contains(engine.ThemeCSS(), "@layer kazari") {
		t.Error("ThemeCSS should respect cascade layer")
	}
}

func TestAssets_Hashing(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"))

	assets := engine.Assets()
	if len(assets.CSS.Hash) != 8 {
		t.Errorf("CSS hash should be 8 chars, got %q", assets.CSS.Hash)
	}
}

func TestJS_ContainsCopyHandler(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCopyButton(true))

	js := engine.JS()
	if !strings.Contains(js, "clipboard") {
		t.Error("JS should contain clipboard handler")
	}
	if !strings.Contains(js, "kz-copy-btn") {
		t.Error("JS should target kz-copy-btn")
	}
}

func TestJS_ContainsFullscreenHandler(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithFullscreenButton(true))

	if !strings.Contains(engine.JS(), "fullscreen") {
		t.Error("JS should contain fullscreen handler")
	}
}

func TestJS_EmptyWhenNoFeatures(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithCopyButton(false), WithFullscreenButton(false), WithWrapButton(false))

	if engine.JS() != "" {
		t.Error("JS should be empty when no features enabled")
	}
}

// --- Line number tests ---

func TestRender_NoLineNumbers_Default(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "kz-gutter") {
		t.Error("kz-gutter should not be present when line numbers disabled")
	}
	if strings.Contains(html, `class="kz-ln"`) {
		t.Error("ln should not be present when line numbers disabled")
	}
}

func TestRender_LineNumbers_Enabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
			{{Content: "c", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	html, err := engine.Render("a\nb\nc", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `<div class="kz-gutter">`) {
		t.Error("missing gutter div")
	}
	if !strings.Contains(html, `aria-hidden="true">1</div>`) {
		t.Error("missing line number 1")
	}
	if !strings.Contains(html, `aria-hidden="true">2</div>`) {
		t.Error("missing line number 2")
	}
	if !strings.Contains(html, `aria-hidden="true">3</div>`) {
		t.Error("missing line number 3")
	}
}

func TestRender_LineNumbers_StartLineNumber(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	ln := true
	start := 42
	html, err := engine.Render("a\nb", Options{Lang: "go", LineNumbers: &ln, StartLineNumber: &start})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `aria-hidden="true">42</div>`) {
		t.Error("missing line number 42")
	}
	if !strings.Contains(html, `aria-hidden="true">43</div>`) {
		t.Error("missing line number 43")
	}
}

func TestRender_LineNumbers_DynamicWidth(t *testing.T) {
	tokens := make([][]Token, 150)
	for i := range tokens {
		tokens[i] = []Token{{Content: "x", Color: "#000"}}
	}
	hl := &mockHighlighter{
		lightTokens: tokens,
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	code := strings.Repeat("x\n", 149) + "x"
	html, err := engine.Render(code, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "--kz-ln-width:3ch") {
		t.Error("should set --kz-ln-width:3ch for 150 lines")
	}
}

func TestRender_LineNumbers_NoWidthOverrideForSmallBlocks(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	html, err := engine.Render("a\nb", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "--kz-ln-width") {
		t.Error("should not set --kz-ln-width for small blocks (≤2 digits)")
	}
}

func TestRender_LineNumbers_AriaHidden(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("line numbers should have aria-hidden=\"true\"")
	}
}

func TestRender_LineNumbers_TerminalFrame(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "echo hi", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	html, err := engine.Render("echo hi", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "is-terminal") {
		t.Error("should render terminal frame")
	}
	if !strings.Contains(html, `<div class="kz-gutter">`) {
		t.Error("gutter should render inside terminal frame")
	}
}

func TestRenderWithMeta_LineNumbers(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("a\nb", `go showLineNumbers startLineNumber=10`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `<div class="kz-gutter">`) {
		t.Error("missing gutter from meta showLineNumbers")
	}
	if !strings.Contains(html, `aria-hidden="true">10</div>`) {
		t.Error("missing line number 10 from meta startLineNumber=10")
	}
	if !strings.Contains(html, `aria-hidden="true">11</div>`) {
		t.Error("missing line number 11")
	}
}

func TestRender_LineNumbers_StartLineNumber_Zero(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
			{{Content: "c", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	ln := true
	start := 0
	html, err := engine.Render("a\nb\nc", Options{Lang: "go", LineNumbers: &ln, StartLineNumber: &start})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `aria-hidden="true">0</div>`) {
		t.Error("missing line number 0")
	}
	if !strings.Contains(html, `aria-hidden="true">1</div>`) {
		t.Error("missing line number 1")
	}
	if !strings.Contains(html, `aria-hidden="true">2</div>`) {
		t.Error("missing line number 2")
	}
}

func TestRender_LineNumbers_StartLineNumber_Negative(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
			{{Content: "c", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	ln := true
	start := -5
	html, err := engine.Render("a\nb\nc", Options{Lang: "go", LineNumbers: &ln, StartLineNumber: &start})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `aria-hidden="true">-5</div>`) {
		t.Error("missing line number -5")
	}
	if !strings.Contains(html, `aria-hidden="true">-4</div>`) {
		t.Error("missing line number -4")
	}
	if !strings.Contains(html, `aria-hidden="true">-3</div>`) {
		t.Error("missing line number -3")
	}
}

func TestRender_LineNumbers_NegativeWidth(t *testing.T) {
	tokens := make([][]Token, 5)
	for i := range tokens {
		tokens[i] = []Token{{Content: "x", Color: "#000"}}
	}
	hl := &mockHighlighter{
		lightTokens: tokens,
		themeInfo:   ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl, WithLineNumbers(true))
	ln := true
	start := -100
	code := strings.Repeat("x\n", 4) + "x"
	html, err := engine.Render(code, Options{Lang: "go", LineNumbers: &ln, StartLineNumber: &start})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "--kz-ln-width:4ch") {
		t.Error("should set --kz-ln-width:4ch for start=-100 (4 chars including minus sign)")
	}
}

func TestRenderWithMeta_StartLineNumber_Zero(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
			{{Content: "c", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("a\nb\nc", `go showLineNumbers startLineNumber=0`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `aria-hidden="true">0</div>`) {
		t.Error("missing line number 0 from meta startLineNumber=0")
	}
	if !strings.Contains(html, `aria-hidden="true">1</div>`) {
		t.Error("missing line number 1")
	}
}

func TestRenderWithMeta_StartLineNumber_Negative(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{{Content: "b", Color: "#000"}},
			{{Content: "c", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("a\nb\nc", `go showLineNumbers startLineNumber=-5`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `aria-hidden="true">-5</div>`) {
		t.Error("missing line number -5 from meta startLineNumber=-5")
	}
	if !strings.Contains(html, `aria-hidden="true">-4</div>`) {
		t.Error("missing line number -4")
	}
}

// --- Text marker tests ---

func threeLineMock() *mockHighlighter {
	return &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "line one", Color: "#000"}},
			{{Content: "line two", Color: "#000"}},
			{{Content: "line three", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
}

func fiveLineMock() *mockHighlighter {
	return &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "line one", Color: "#000"}},
			{{Content: "line two", Color: "#000"}},
			{{Content: "line three", Color: "#000"}},
			{{Content: "line four", Color: "#000"}},
			{{Content: "line five", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
}

func TestRenderWithMeta_LineMarkers_Mark(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go {2}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `class="kz-line highlight mark"`) {
		t.Error("line 2 should have highlight mark class")
	}
	if strings.Count(html, "highlight") != 1 {
		t.Error("only line 2 should be highlighted")
	}
}

func TestRenderWithMeta_LineMarkers_Ins(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go ins={1}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `class="kz-line highlight ins"`) {
		t.Error("line 1 should have highlight ins class")
	}
}

func TestRenderWithMeta_LineMarkers_Del(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go del={3}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `class="kz-line highlight del"`) {
		t.Error("line 3 should have highlight del class")
	}
}

func TestRenderWithMeta_LineMarkers_Combined(t *testing.T) {
	engine := newTestEngine(fiveLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three\nline four\nline five", `go {1} ins={3} del={5}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `highlight mark`) {
		t.Error("missing mark class")
	}
	if !strings.Contains(html, `highlight ins`) {
		t.Error("missing ins class")
	}
	if !strings.Contains(html, `highlight del`) {
		t.Error("missing del class")
	}
}

func TestRenderWithMeta_LineMarkers_UnmarkedLinesPlain(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go {2}`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(html, `class="kz-line"`) != 2 {
		t.Errorf("expected 2 plain kz-line divs, got %d", strings.Count(html, `class="kz-line"`))
	}
}

func TestRenderWithMeta_LabeledRange(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go ins={"A":1-3}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `tm-label`) {
		t.Error("first line should have tm-label class")
	}
	if !strings.Contains(html, `data-label="A"`) {
		t.Error("first line should have data-label attribute")
	}
	if strings.Count(html, "tm-label") != 1 {
		t.Error("only the first line of the range should have the label")
	}
}

func TestRenderWithMeta_FocusLines_HasFocusOnCode(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go focus={2}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `class="has-focus"`) {
		t.Error("code element should have has-focus class")
	}
}

func TestRenderWithMeta_FocusLines_FocusedClass(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go focus={2}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `kz-line focused`) {
		t.Error("focused line should have focused class")
	}
	if strings.Count(html, `kz-line focused`) != 1 {
		t.Errorf("only 1 line should be focused, got %d", strings.Count(html, `kz-line focused`))
	}
}

func TestRenderWithMeta_NoFocus_NoHasFocusClass(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.RenderWithMeta("line one\nline two\nline three", `go`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "has-focus") {
		t.Error("should not have has-focus without focus lines")
	}
	if strings.Contains(html, "focused") {
		t.Error("should not have focused class without focus lines")
	}
}

func TestRender_LineMarkersViaOptions(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.Render("line one\nline two\nline three", Options{
		Lang: "go",
		LineMarkers: []LineMarker{
			{Type: MarkerIns, Lines: []Range{{Start: 2, End: 2}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `highlight ins`) {
		t.Error("line 2 should have ins class via Render() API")
	}
}

func TestRender_FocusLinesViaOptions(t *testing.T) {
	engine := newTestEngine(threeLineMock())
	html, err := engine.Render("line one\nline two\nline three", Options{
		Lang:       "go",
		FocusLines: []Range{{Start: 1, End: 1}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "has-focus") {
		t.Error("should have has-focus class")
	}
}

func TestRenderWithMeta_InlineMarker_MarkElement(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "const useState = true", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("const useState = true", `go "use"`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<mark>use</mark>") {
		t.Errorf("expected <mark>use</mark> in output, got: %s", html)
	}
}

func TestRenderWithMeta_InlineMarker_InsElement(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "hello world", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("hello world", `go ins="world"`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<ins>world</ins>") {
		t.Errorf("expected <ins>world</ins>, got: %s", html)
	}
}

func TestRenderWithMeta_InlineMarker_DelElement(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "hello world", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("hello world", `go del="hello"`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<del>hello</del>") {
		t.Errorf("expected <del>hello</del>, got: %s", html)
	}
}

func TestRenderWithMeta_InlineMarker_MultiTokenSpan(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{
				{Content: "fmt", Color: "#000"},
				{Content: ".", Color: "#111"},
				{Content: "Println", Color: "#222"},
			},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("fmt.Println", `go "fmt.Println"`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `class="open-end"`) {
		t.Error("first token should have open-end class")
	}
	if !strings.Contains(html, `class="open-start open-end"`) {
		t.Error("middle token should have both open classes")
	}
	if !strings.Contains(html, `class="open-start"`) && !strings.Contains(html, `open-start">`) {
		t.Error("last token should have open-start class")
	}
}

func TestRender_InlineMarkersViaOptions(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "hello world", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello world", Options{
		Lang:          "go",
		InlineMarkers: []InlineMarker{{Type: MarkerMark, Text: "world"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<mark>world</mark>") {
		t.Errorf("expected <mark>world</mark> via Render() API, got: %s", html)
	}
}

func TestRenderWithMeta_InlineMarker_SingleQuote(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "hello world", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("hello world", `go 'hello'`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<mark>hello</mark>") {
		t.Errorf("expected <mark>hello</mark> from single-quoted marker, got: %s", html)
	}
}

// --- Contrast adjustment tests ---

func TestRender_ContrastAdjustment_MarkedLine(t *testing.T) {
	// Use a token color very close to the effective mark background on white.
	// Mark bg rgba(255,200,0,0.12) on #ffffff composites to ~#fff8e0 (light yellow).
	// A light yellow token (#f0e8b0) has poor contrast against that bg.
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "test", Color: "#f0e8b0"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithMinContrast(4.5))
	html, err := engine.RenderWithMeta("test", `go {1}`)
	if err != nil {
		t.Fatal(err)
	}
	// The adjusted color should NOT be the original low-contrast color.
	if strings.Contains(html, "--sl:#f0e8b0") {
		t.Error("token color should be adjusted for contrast on marked line, but original color was kept")
	}
	if !strings.Contains(html, "--sl:") {
		t.Error("expected --sl: color to be present")
	}
}

func TestRender_ContrastAdjustment_UnmarkedLine(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "test", Color: "#f0e8b0"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithMinContrast(4.5))
	html, err := engine.RenderWithMeta("test", `go`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "--sl:#f0e8b0") {
		t.Error("token color should be adjusted on unmarked line, but original low-contrast color was kept")
	}
	if !strings.Contains(html, "--sl:") {
		t.Error("expected --sl: color to be present after adjustment")
	}
}

func TestRender_ContrastAdjustment_UnmarkedLine_DualTheme(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "test", Color: "#f0e8b0"}},
		},
		darkTokens: [][]Token{
			{{Content: "test", Color: "#3a3520"}},
		},
		themeInfo:      ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo:  ThemeInfo{FG: "#c9d1d9", BG: "#0d1117"},
		lightThemeName: "light",
	}
	engine := newTestEngine(hl,
		WithMinContrast(4.5),
		WithThemes("light", "dark"),
	)
	html, err := engine.RenderWithMeta("test", `go`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "--sl:#f0e8b0") {
		t.Error("light color should be adjusted on unmarked line")
	}
	if strings.Contains(html, "--sd:#3a3520") {
		t.Error("dark color should be adjusted on unmarked line")
	}
}

func TestRender_ContrastAdjustment_UnmarkedLine_AlreadyMeetsContrast(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "test", Color: "#000000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithMinContrast(4.5))
	html, err := engine.RenderWithMeta("test", `go`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "--sl:#000000") {
		t.Error("high-contrast token should not be changed")
	}
}

func TestRender_ContrastAdjustment_Disabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "test", Color: "#f0e8b0"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithMinContrast(0))
	html, err := engine.RenderWithMeta("test", `go {1}`)
	if err != nil {
		t.Fatal(err)
	}
	// With MinContrast=0, no adjustment should happen.
	if !strings.Contains(html, "--sl:#f0e8b0") {
		t.Error("token color should NOT be adjusted when MinContrast=0")
	}
}

func TestRender_ContrastAdjustment_DualTheme(t *testing.T) {
	// Light token has poor contrast on light mark bg; dark token has poor contrast on dark mark bg.
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "test", Color: "#f0e8b0"}},
		},
		darkTokens: [][]Token{
			{{Content: "test", Color: "#3a3520"}},
		},
		themeInfo:      ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo:  ThemeInfo{FG: "#c9d1d9", BG: "#0d1117"},
		lightThemeName: "light",
	}
	engine := newTestEngine(hl,
		WithMinContrast(4.5),
		WithThemes("light", "dark"),
	)
	html, err := engine.RenderWithMeta("test", `go {1}`)
	if err != nil {
		t.Fatal(err)
	}
	// Both colors should be adjusted.
	if strings.Contains(html, "--sl:#f0e8b0") {
		t.Error("light color should be adjusted")
	}
	if strings.Contains(html, "--sd:#3a3520") {
		t.Error("dark color should be adjusted")
	}
	if !strings.Contains(html, "--sl:") {
		t.Error("expected --sl: to be present")
	}
	if !strings.Contains(html, "--sd:") {
		t.Error("expected --sd: to be present")
	}
}

func TestCSS_ContainsMarkerVars(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)
	css := engine.CSS()
	vars := []string{
		"--kz-mark-bg", "--kz-mark-border", "--kz-mark-border-width",
		"--kz-ins-bg", "--kz-ins-border", "--kz-ins-indicator",
		"--kz-del-bg", "--kz-del-border", "--kz-del-indicator",
		"--kz-inline-mark-bg", "--kz-inline-mark-border",
		"--kz-focus-dimmed-opacity",
	}
	for _, v := range vars {
		if !strings.Contains(css, v) {
			t.Errorf("CSS missing variable: %s", v)
		}
	}
}

func TestCSS_ContainsMarkerStyles(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)
	css := engine.CSS()
	rules := []string{
		".highlight.mark",
		".highlight.ins",
		".highlight.del",
		"open-start",
		"open-end",
		"has-focus",
	}
	for _, r := range rules {
		if !strings.Contains(css, r) {
			t.Errorf("CSS missing rule/class: %s", r)
		}
	}
}

// --- Collapsible section tests ---

func makeMultiLineTokens(n int) [][]Token {
	lines := make([][]Token, n)
	for i := range lines {
		lines[i] = []Token{{Content: "line", Color: "#333"}}
	}
	return lines
}

func TestRender_Collapsible_ThresholdTriggered(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(20),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold:    15,
			PreviewLines:     8,
			DefaultCollapsed: true,
		}),
	)

	code := strings.Repeat("line\n", 19) + "line"
	html, err := engine.Render(code, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "kz-collapsed") {
		t.Error("missing kz-collapsed class on wrapper")
	}
	if !strings.Contains(html, "kz-hidden") {
		t.Error("missing kz-hidden class on hidden lines")
	}
	if !strings.Contains(html, "kz-collapse-gradient") {
		t.Error("missing gradient overlay")
	}
	if !strings.Contains(html, "kz-collapse-btn") {
		t.Error("missing collapse button")
	}
	if !strings.Contains(html, `aria-expanded="false"`) {
		t.Error("missing aria-expanded attribute")
	}
	if !strings.Contains(html, `data-expand="Show more"`) {
		t.Error("missing data-expand attribute")
	}
}

func TestRender_Collapsible_BelowThreshold(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(5),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	code := "a\nb\nc\nd\ne"
	html, err := engine.Render(code, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "kz-collapsed") {
		t.Error("should not have kz-collapsed class for short block")
	}
	if strings.Contains(html, "kz-hidden") {
		t.Error("should not have hidden lines for short block")
	}
}

func TestRenderWithMeta_Collapsible_ForceCollapse(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(5),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 3}),
	)

	code := "a\nb\nc\nd\ne"
	html, err := engine.RenderWithMeta(code, "go collapse")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, "kz-collapse-btn") {
		t.Error("collapse meta should force threshold collapse even below threshold")
	}
}

func TestRenderWithMeta_Collapsible_Nocollapse(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(20),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	code := strings.Repeat("line\n", 19) + "line"
	html, err := engine.RenderWithMeta(code, "go nocollapse")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if strings.Contains(html, "kz-collapse-btn") {
		t.Error("nocollapse should prevent threshold collapse")
	}
}

func TestRenderWithMeta_Collapsible_PerBlockThreshold_PreventsCollapse(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(16),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	code := strings.Repeat("line\n", 15) + "line"
	html, err := engine.RenderWithMeta(code, "go collapseThreshold=20")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if strings.Contains(html, "kz-collapse-btn") {
		t.Error("per-block threshold=20 should prevent collapse at 16 lines")
	}
}

func TestRenderWithMeta_Collapsible_PerBlockThreshold_TriggersCollapse(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(11),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 3}),
	)

	code := strings.Repeat("line\n", 10) + "line"
	html, err := engine.RenderWithMeta(code, "go collapseThreshold=10")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, "kz-collapse-btn") {
		t.Error("per-block threshold=10 should trigger collapse at 11 lines")
	}
}

func TestRenderWithMeta_Collapsible_PerBlockThreshold_NocollapseWins(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(20),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	code := strings.Repeat("line\n", 19) + "line"
	html, err := engine.RenderWithMeta(code, "go nocollapse collapseThreshold=5")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if strings.Contains(html, "kz-collapse-btn") {
		t.Error("nocollapse should win over per-block threshold")
	}
}

func TestRenderWithMeta_Collapsible_PerBlockThreshold_CollapseWins(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(5),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 3}),
	)

	code := "a\nb\nc\nd\ne"
	html, err := engine.RenderWithMeta(code, "go collapse collapseThreshold=999")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, "kz-collapse-btn") {
		t.Error("collapse should force collapse regardless of per-block threshold")
	}
}

func TestRender_Collapsible_PerBlockThreshold_StructuredAPI(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(16),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	threshold := 20
	code := strings.Repeat("line\n", 15) + "line"
	html, err := engine.Render(code, Options{
		Lang:     "go",
		Collapse: &CollapseOptions{Threshold: &threshold},
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "kz-collapse-btn") {
		t.Error("structured API threshold=20 should prevent collapse at 16 lines")
	}
}

func TestRender_Collapsible_ThresholdAfterFileNameExtraction(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(15),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	// 16 lines before extraction, 15 after the filename comment is removed.
	// The threshold check must use the post extraction count of 15, which
	// does not exceed the threshold.
	code := "// main.go\n" + strings.Repeat("line\n", 14) + "line"
	html, err := engine.Render(code, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "main.go") {
		t.Fatal("expected filename comment to be extracted into the title")
	}
	if strings.Contains(html, "kz-collapse-btn") {
		t.Error("threshold collapse must use the post extraction line count")
	}
}

func TestRenderWithMeta_Collapsible_RangeAfterFileNameExtraction(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(9),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 50, PreviewLines: 8}),
	)

	// 10 lines before extraction, 9 after the filename comment is removed.
	// The range end must be clamped to the post extraction line count,
	// collapsing lines 2 through 9 for a count of 8.
	code := "// main.go\na\nb\nc\nd\ne\nf\ng\nh\ni"
	html, err := engine.RenderWithMeta(code, "go collapse={2-10}")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, "8 collapsed lines") {
		t.Error("collapse range must be clamped to the post extraction line count")
	}
}

func TestRenderWithMeta_Collapsible_RangeBased(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold:  15,
			PreviewLines:   8,
			PreserveIndent: true,
		}),
	)

	code := "a\n    b\n    c\n    d\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, "go collapse={2-4}")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `<details class="kz-section">`) {
		t.Error("missing <details> for range-based collapse")
	}
	if !strings.Contains(html, "<summary>") {
		t.Error("missing <summary> element")
	}
	if !strings.Contains(html, "3 collapsed lines") {
		t.Error("missing summary text")
	}
	if !strings.Contains(html, `--kz-indent:4ch`) {
		t.Error("missing indent preservation")
	}
	if !strings.Contains(html, `class="expand"`) {
		t.Error("missing expand icon span")
	}
}

func TestRenderWithMeta_Collapsible_RangeWithLineNumbers(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(5),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)

	code := "a\nb\nc\nd\ne"
	html, err := engine.RenderWithMeta(code, "go showLineNumbers collapse={2-3}")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	// Summary should have empty gutter for alignment
	if !strings.Contains(html, `<div class="kz-gutter"><div class="kz-ln"></div></div>`) {
		t.Error("missing empty gutter placeholder in summary line")
	}
}

func TestRenderWithMeta_Collapsible_MultipleRanges(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(15),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)

	code := strings.Repeat("line\n", 14) + "line"
	html, err := engine.RenderWithMeta(code, "go collapse={2-4,8-10}")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	count := strings.Count(html, `<details class="kz-section">`)
	if count != 2 {
		t.Errorf("expected 2 <details> sections, got %d", count)
	}
}

func TestRender_Collapsible_CustomButtonText(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(20),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold:      15,
			PreviewLines:       8,
			DefaultCollapsed:   true,
			ExpandButtonText:   "Expand",
			CollapseButtonText: "Collapse",
		}),
	)

	code := strings.Repeat("line\n", 19) + "line"
	html, err := engine.Render(code, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `data-expand="Expand"`) {
		t.Error("missing custom expand text")
	}
	if !strings.Contains(html, `data-collapse="Collapse"`) {
		t.Error("missing custom collapse text")
	}
}

func TestCSS_ContainsCollapsibleStyles(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithCollapsible(CollapsibleConfig{LineThreshold: 15}),
	)
	css := engine.CSS()

	rules := []string{
		"kz-collapse-gradient",
		"kz-collapse-btn",
		"kz-hidden",
		"kz-section",
		"@media print",
	}
	for _, r := range rules {
		if !strings.Contains(css, r) {
			t.Errorf("CSS missing collapsible rule: %s", r)
		}
	}
}

func TestCSS_NoCollapsibleWhenDisabled(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)
	css := engine.CSS()

	if strings.Contains(css, "kz-collapse-gradient") {
		t.Error("collapsible CSS should not be included when feature is disabled")
	}
}

func TestJS_ContainsCollapsibleHandler(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithCollapsible(CollapsibleConfig{LineThreshold: 15}),
	)
	js := engine.JS()

	if !strings.Contains(js, "kz-collapse-btn") {
		t.Error("JS missing collapsible handler")
	}
	if !strings.Contains(js, "aria-expanded") {
		t.Error("JS missing aria-expanded toggle")
	}
}

func TestJS_NoCollapsibleWhenDisabled(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)
	js := engine.JS()

	if strings.Contains(js, "kz-collapse-btn") {
		t.Error("collapsible JS should not be included when feature is disabled")
	}
}

func TestRender_Collapsible_GapIndicator(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(30),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold:    15,
			PreviewLines:     8,
			DefaultCollapsed: true,
		}),
	)

	code := strings.Repeat("line\n", 29) + "line"
	// Mark line 12 — within 2× cap (16), should create gap indicator
	html, err := engine.Render(code, Options{
		Lang: "go",
		LineMarkers: []LineMarker{
			{Type: MarkerMark, Lines: []Range{{Start: 12, End: 12}}},
		},
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "kz-gap") {
		t.Error("missing gap indicator for non-contiguous preview segments")
	}
	if !strings.Contains(html, "kz-gap-indicator") {
		t.Error("missing gap indicator span")
	}
}

func TestRender_Collapsible_NoGapWithoutMarkers(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(30),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold:    15,
			PreviewLines:     8,
			DefaultCollapsed: true,
		}),
	)

	code := strings.Repeat("line\n", 29) + "line"
	html, err := engine.Render(code, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "kz-gap") {
		t.Error("should not have gap indicators without markers")
	}
}

func TestRender_Collapsible_BadgeFallback(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(50),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold:    15,
			PreviewLines:     8,
			DefaultCollapsed: true,
		}),
	)

	code := strings.Repeat("line\n", 49) + "line"
	// Mark lines 30-32 — beyond 2× cap (16), should trigger badge
	html, err := engine.Render(code, Options{
		Lang: "go",
		LineMarkers: []LineMarker{
			{Type: MarkerMark, Lines: []Range{{Start: 30, End: 32}}},
		},
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "(+3 highlighted)") {
		t.Error("missing badge fallback count on expand button")
	}
}

func TestCSS_ContainsGapStyles(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithCollapsible(CollapsibleConfig{LineThreshold: 15}),
	)
	css := engine.CSS()

	if !strings.Contains(css, "kz-gap") {
		t.Error("CSS missing gap indicator styles")
	}
}

// --- Collapse style tests ---

func TestRenderWithMeta_CollapseStyle_CollapsibleStart(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)

	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, `go collapse={2-4} collapseStyle="collapsible-start"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `<div class="kz-section collapsible-start">`) {
		t.Error("missing collapsible-start wrapper div")
	}
	if !strings.Contains(html, `<div class="content-lines">`) {
		t.Error("missing content-lines div")
	}
	if strings.Contains(html, `<details class="kz-section">`) {
		t.Error("should not have github-style details wrapper")
	}
}

func TestRenderWithMeta_CollapseStyle_CollapsibleEnd(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)

	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, `go collapse={2-4} collapseStyle="collapsible-end"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `<div class="kz-section collapsible-end">`) {
		t.Error("missing collapsible-end wrapper div")
	}
}

func TestRenderWithMeta_CollapseStyle_Auto_NotAtEnd(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)

	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, `go collapse={2-4} collapseStyle="collapsible-auto"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `collapsible-start`) {
		t.Error("auto should resolve to collapsible-start when range is not at end")
	}
}

func TestRenderWithMeta_CollapseStyle_Auto_AtEnd(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)

	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, `go collapse={8-10} collapseStyle="collapsible-auto"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `collapsible-end`) {
		t.Error("auto should resolve to collapsible-end when range reaches last line")
	}
}

func TestRenderWithMeta_CollapseStyle_GithubDefault(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100}),
	)

	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, `go collapse={2-4}`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `<details class="kz-section">`) {
		t.Error("default style should be github with details wrapper")
	}
	if strings.Contains(html, "content-lines") {
		t.Error("github style should not have content-lines div")
	}
}

func TestRender_CollapseStyle_EngineDefault(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{
			LineThreshold: 100,
			Style:         CollapseCollapsibleStart,
		}),
	)

	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	html, err := engine.RenderWithMeta(code, `go collapse={2-4}`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `collapsible-start`) {
		t.Error("engine-level Style should apply when no meta collapseStyle specified")
	}
}

func TestCSS_ContainsCollapsibleStartStyles(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithCollapsible(CollapsibleConfig{LineThreshold: 15}),
	)
	css := engine.CSS()

	rules := []string{
		"collapsible-start",
		"collapsible-end",
		"content-lines",
		"column-reverse",
	}
	for _, r := range rules {
		if !strings.Contains(css, r) {
			t.Errorf("CSS missing collapsible style rule: %s", r)
		}
	}
}

// --- Terminal comment stripping tests ---

func extractDataCode(t *testing.T, html string) string {
	t.Helper()
	idx := strings.Index(html, `data-code="`)
	if idx == -1 {
		t.Fatal("no data-code attribute found")
	}
	start := idx + len(`data-code="`)
	end := strings.Index(html[start:], `"`)
	return html[start : start+end]
}

func TestTerminalCommentStripping_DefaultEnabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "# Install deps", Color: "#6a737d"}},
			{{Content: "npm install", Color: "#24292f"}},
			{{Content: "# Start server", Color: "#6a737d"}},
			{{Content: "npm start", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#1e1e1e"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("# Install deps\nnpm install\n# Start server\nnpm start", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	dataCode := extractDataCode(t, html)
	if strings.Contains(dataCode, "Install deps") {
		t.Error("data-code should not contain comment text")
	}
	if !strings.Contains(dataCode, "npm install") {
		t.Error("data-code should contain commands")
	}
}

func TestTerminalCommentStripping_PreservesCommands(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "# comment", Color: "#6a737d"}},
			{{Content: "echo hello", Color: "#24292f"}},
			{{Content: "curl example.com", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#1e1e1e"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("# comment\necho hello\ncurl example.com", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	// data-code uses DEL (\x7f) for newlines
	if !strings.Contains(html, "echo hello") {
		t.Error("missing 'echo hello' in output")
	}
	if !strings.Contains(html, "curl example.com") {
		t.Error("missing 'curl example.com' in output")
	}
}

func TestTerminalCommentStripping_EditorFrameUnaffected(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "// a Go comment", Color: "#6a737d"}},
			{{Content: "func main() {}", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("// a Go comment\nfunc main() {}", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "a Go comment") {
		t.Error("editor frame should preserve comments in data-code")
	}
}

func TestTerminalCommentStripping_Disabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "# Install deps", Color: "#6a737d"}},
			{{Content: "npm install", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#1e1e1e"},
	}
	engine := newTestEngine(hl, WithTerminalCommentStripping(false))
	html, err := engine.Render("# Install deps\nnpm install", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	dataCode := extractDataCode(t, html)
	if !strings.Contains(dataCode, "Install deps") {
		t.Error("with stripping disabled, comments should remain in data-code")
	}
}

func TestTerminalCommentStripping_IndentedComments(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "  # indented comment", Color: "#6a737d"}},
			{{Content: "npm install", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#1e1e1e"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("  # indented comment\nnpm install", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	dataCode := extractDataCode(t, html)
	if strings.Contains(dataCode, "indented comment") {
		t.Error("indented comments should be stripped from data-code")
	}
}

func TestTerminalCommentStripping_InlineHashPreserved(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "echo hello # inline", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#1e1e1e"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("echo hello # inline", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "echo hello # inline") {
		t.Error("inline # should not be stripped — only full comment lines")
	}
}

func TestTerminalCommentStripping_CommentsStillRendered(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "# visible comment", Color: "#6a737d"}},
			{{Content: "npm install", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#1e1e1e"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("# visible comment\nnpm install", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	// Comments are stripped from data-code but still rendered in the visible HTML
	if !strings.Contains(html, "visible comment") {
		t.Error("comments should still be visible in the rendered code block")
	}
	// But the data-code attribute should not contain them
	// Extract data-code value: it uses DEL encoding
	idx := strings.Index(html, `data-code="`)
	if idx == -1 {
		t.Fatal("no data-code attribute found")
	}
	dataCodeStart := idx + len(`data-code="`)
	dataCodeEnd := strings.Index(html[dataCodeStart:], `"`)
	dataCode := html[dataCodeStart : dataCodeStart+dataCodeEnd]
	if strings.Contains(dataCode, "visible comment") {
		t.Error("data-code should not contain stripped comment text")
	}
}

// --- Mermaid pass-through tests ---

func TestMermaid_RendersRawCode(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "graph TD", Color: "#24292f"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("graph TD;\n  A-->B;", Options{Lang: "mermaid"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `<pre class="mermaid">`) {
		t.Error("mermaid block should render <pre class=\"mermaid\">")
	}
	if !strings.Contains(html, "A--&gt;B;") {
		t.Error("mermaid code should be HTML-escaped")
	}
}

func TestMermaid_NoFrameMarkup(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "graph TD", Color: "#24292f"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("graph TD;", Options{Lang: "mermaid"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "kazari-block") {
		t.Error("mermaid block should not contain kazari-block wrapper")
	}
	if strings.Contains(html, "kz-copy-btn") {
		t.Error("mermaid block should not contain copy button")
	}
	if strings.Contains(html, "data-code") {
		t.Error("mermaid block should not contain data-code attribute")
	}
}

func TestMermaid_DisabledRendersNormally(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "graph TD", Color: "#24292f"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithMermaidPassThrough(false))
	html, err := engine.Render("graph TD;", Options{Lang: "mermaid"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, `<pre class="mermaid">`) {
		t.Error("with mermaid disabled, should not render raw mermaid block")
	}
	if !strings.Contains(html, "kazari-block") {
		t.Error("with mermaid disabled, should render normally with kazari-block wrapper")
	}
}

func TestMermaid_NonMermaidUnaffected(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "func main", Color: "#24292f"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("func main() {}", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, `<pre class="mermaid">`) {
		t.Error("non-mermaid language should not produce mermaid block")
	}
}

// --- Regex marker integration tests ---

func TestRegexMarker_RendersInlineMarkup(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "func main() {}", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("func main() {}", `go /func/`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "<mark>") {
		t.Error("regex marker should produce <mark> element")
	}
}

func TestRegexMarker_CaptureGroup(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "yes and yep", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("yes and yep", `text /ye(s|p)/`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "<mark>") {
		t.Error("capture group regex should produce <mark> elements")
	}
}

// --- Per-block theme override tests ---

func TestThemeOverride_MetaParsed(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "func main", Color: "#24292f"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("func main", `go theme="dracula"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "kazari-block") {
		t.Error("should render normally with theme override")
	}
}

func TestThemeOverride_DefaultUnchanged(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#aaa"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "#aaa") {
		t.Error("default theme should produce expected token colors")
	}
}

func TestThemeOverride_GoAPI(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#aaa"}}},
		themeInfo:    ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	_, err := engine.Render("x", Options{Lang: "go", Theme: "dracula"})
	if err != nil {
		t.Fatalf("Render() with Theme option error: %v", err)
	}
}

// --- Hybrid diff tests ---

func TestDiff_StripsPrefixes(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "const x = 1;", Color: "#24292f"}},
			{{Content: "const y = 2;", Color: "#24292f"}},
			{{Content: "const z = 3;", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("  const x = 1;\n+ const y = 2;\n- const z = 3;", `diff lang="javascript"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if strings.Contains(html, "highlight ins") || strings.Contains(html, "highlight del") {
		// Markers are applied, good
	}
	if !strings.Contains(html, "const x") {
		t.Error("diff block should contain stripped code")
	}
}

func TestDiff_AppliesMarkers(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#24292f"}},
			{{Content: "b", Color: "#24292f"}},
			{{Content: "c", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("  a\n+ b\n- c", `diff lang="text"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "ins") {
		t.Error("+ lines should produce ins markers")
	}
	if !strings.Contains(html, "del") {
		t.Error("- lines should produce del markers")
	}
}

func TestDiff_NoLangRendersNormally(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "+ added", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("+ added", `diff`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "+ added") {
		t.Error("without lang= meta, diff should render as plain diff text")
	}
}

// --- Options API parity tests ---

func TestRender_PreserveIndent_PerBlock(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "    foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: true}))

	noPreserve := false
	html, err := engine.Render("    foo()", Options{Lang: "go", PreserveIndent: &noPreserve})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, `class="indent"`) {
		t.Error("PreserveIndent=false per-block should suppress indent spans")
	}
}

func TestRender_HangingIndent_PerBlock(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "  foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	noPreserve := false
	wrap := true
	hanging := 4
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: false}))
	html, err := engine.Render("  foo()", Options{
		Lang: "go", Wrap: &wrap, PreserveIndent: &noPreserve, HangingIndent: &hanging,
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `--kz-indent:4ch`) {
		t.Errorf("HangingIndent=4 per-block should produce --kz-indent:4ch, got: %s", html)
	}
}

func TestRender_CollapseStyle_PerBlock(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: makeMultiLineTokens(10),
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 100, PreserveIndent: true}),
	)
	code := "a\nb\nc\nd\ne\nf\ng\nh\ni\nj"
	style := CollapseCollapsibleStart
	html, err := engine.Render(code, Options{
		Lang: "go",
		Collapse: &CollapseOptions{
			Ranges: []Range{{Start: 2, End: 4}},
			Style:  &style,
		},
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `collapsible-start`) {
		t.Error("Collapse.Style should produce collapsible-start wrapper via Render()")
	}
}

func TestRender_DiffLang_PerBlock(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#24292f"}},
			{{Content: "b", Color: "#24292f"}},
			{{Content: "c", Color: "#24292f"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("  a\n+ b\n- c", Options{Lang: "diff", DiffLang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "ins") {
		t.Error("DiffLang via Render() should produce ins markers for + lines")
	}
	if !strings.Contains(html, "del") {
		t.Error("DiffLang via Render() should produce del markers for - lines")
	}
}

// --- ANSI rendering tests ---

func TestRender_ANSI_BasicColor(t *testing.T) {
	engine := newTestEngine(&mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	})
	html, err := engine.Render("\x1b[31mhello\x1b[0m", Options{Lang: "ansi"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "--sl:var(--kz-ansi-red)") {
		t.Error("expected red color CSS var in --sl custom property")
	}
	if !strings.Contains(html, "hello") {
		t.Error("expected token content")
	}
}

func TestRender_ANSI_TerminalFrame(t *testing.T) {
	engine := newTestEngine(&mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	})
	html, err := engine.Render("\x1b[32mok\x1b[0m", Options{Lang: "ansi"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "is-terminal") {
		t.Error("ANSI blocks should auto-detect as terminal frame")
	}
}

func TestRender_ANSI_NoHighlighter(t *testing.T) {
	engine := New(WithMinify(false))
	html, err := engine.Render("\x1b[34mblue\x1b[0m", Options{Lang: "ansi"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "--sl:var(--kz-ansi-blue)") {
		t.Error("ANSI rendering should work without a highlighter")
	}
}

func TestRender_ANSI_WithLineNumbers(t *testing.T) {
	engine := newTestEngine(&mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}, WithLineNumbers(true))
	html, err := engine.Render("\x1b[31mline1\x1b[0m\nline2", Options{Lang: "ansi"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `class="kz-ln"`) {
		t.Error("ANSI blocks should support line numbers")
	}
}

func TestRender_ANSI_WithMarkers(t *testing.T) {
	engine := newTestEngine(&mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	})
	html, err := engine.Render("line1\nline2\nline3", Options{
		Lang:        "ansi",
		LineMarkers: []LineMarker{{Type: MarkerMark, Lines: []Range{{Start: 2, End: 2}}}},
	})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "highlight mark") {
		t.Error("ANSI blocks should support line markers")
	}
}

func TestRender_ANSI_CopyButton(t *testing.T) {
	engine := newTestEngine(&mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}, WithCopyButton(true))
	html, err := engine.Render("\x1b[31mhello\x1b[0m", Options{Lang: "ansi"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "data-code=") {
		t.Error("ANSI blocks should have a copy button with data-code")
	}
}

func TestRenderWithMeta_ANSI(t *testing.T) {
	engine := newTestEngine(&mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	})
	html, err := engine.RenderWithMeta("\x1b[33myellow\x1b[0m", `ansi title="Output"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "--sl:var(--kz-ansi-yellow)") {
		t.Error("expected yellow color")
	}
	if !strings.Contains(html, "Output") {
		t.Error("expected title from meta string")
	}
}

func TestRender_ANSI_DualColorsSame(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
		WithMinContrast(0),
	)
	html, err := engine.Render("\x1b[31mred\x1b[0m", Options{Lang: "ansi"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "--sl:var(--kz-ansi-red)") {
		t.Error("expected light color")
	}
	if !strings.Contains(html, "--sd:var(--kz-ansi-red)") {
		t.Error("expected dark color to match light (ANSI colors are theme-independent)")
	}
}

// --- data-lines attribute tests ---

func TestRender_DataLineCount_Default(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "line1", Color: "#333"}},
			{{Content: "line2", Color: "#333"}},
			{{Content: "line3", Color: "#333"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("line1\nline2\nline3", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `data-lines="3"`) {
		t.Error("data-lines should be present by default with correct count")
	}
}

func TestRender_DataLineCount_Disabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "line1", Color: "#333"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDataLineCount(false))
	html, err := engine.Render("line1", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "data-lines") {
		t.Error("data-lines should not be present when disabled")
	}
}

func TestRender_DataLineCount_SingleLine(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "hello", Color: "#333"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `data-lines="1"`) {
		t.Error("single-line block should have data-lines=\"1\"")
	}
}

// --- Theme CSS root tests ---

func TestCSS_ThemeCSSRoot_Default(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	css := engine.CSS()
	if !strings.Contains(css, ":root {") {
		t.Error("default CSS root should be :root")
	}
}

func TestCSS_ThemeCSSRoot_Custom(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithThemeCSSRoot(".kazari-container"),
	)
	css := engine.CSS()
	if !strings.Contains(css, ".kazari-container {") {
		t.Error("custom CSS root should appear in output")
	}
	if strings.Contains(css, ":root {") {
		t.Error(":root should not appear when custom root is set")
	}
}

func TestCSS_ThemeCSSRoot_DarkModeComposition(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo:     ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo: ThemeInfo{FG: "#e1e4e8", BG: "#0d1117"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
		WithThemeCSSRoot(".container"),
		WithDarkMode(SelectorMode(".dark")),
	)
	css := engine.CSS()
	if !strings.Contains(css, ".container.dark {") {
		t.Error("dark mode selector should compose with custom root")
	}
}

// --- Theme customizer tests ---

func TestNew_ThemeCustomizer_ModifiesBG(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithThemeCustomizer(func(name string, colors ThemeInfo) ThemeInfo {
			colors.BG = "#0a0a0a"
			return colors
		}),
	)
	css := engine.CSS()
	if !strings.Contains(css, "#0a0a0a") {
		t.Error("customizer should modify the background color in CSS output")
	}
	if strings.Contains(css, "--kz-editor-bg:#ffffff") {
		t.Error("original editor BG should be replaced by customizer")
	}
}

func TestNew_ThemeCustomizer_ReceivesThemeName(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo:     ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo: ThemeInfo{FG: "#e1e4e8", BG: "#0d1117"},
	}
	var receivedNames []string
	New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithThemeCustomizer(func(name string, colors ThemeInfo) ThemeInfo {
			receivedNames = append(receivedNames, name)
			return colors
		}),
	)
	if len(receivedNames) != 2 {
		t.Fatalf("expected 2 customizer calls, got %d", len(receivedNames))
	}
	if receivedNames[0] != "light-theme" {
		t.Errorf("first call should be light theme, got %q", receivedNames[0])
	}
	if receivedNames[1] != "dark-theme" {
		t.Errorf("second call should be dark theme, got %q", receivedNames[1])
	}
}

func TestNew_ThemeCustomizer_NilIsNoOp(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)
	css := engine.CSS()
	if !strings.Contains(css, "#ffffff") {
		t.Error("without customizer, original BG should be preserved")
	}
}

func TestNew_ThemeCustomizer_BothThemes(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo:     ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo: ThemeInfo{FG: "#e1e4e8", BG: "#0d1117"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
		WithThemeCustomizer(func(name string, colors ThemeInfo) ThemeInfo {
			if name == "light-theme" {
				colors.BG = "#f0f0f0"
			} else {
				colors.BG = "#111111"
			}
			return colors
		}),
	)
	css := engine.CSS()
	if !strings.Contains(css, "#f0f0f0") {
		t.Error("customizer should modify light theme BG")
	}
	if !strings.Contains(css, "#111111") {
		t.Error("customizer should modify dark theme BG")
	}
}

// --- Locale / UI strings tests ---

func TestRender_Locale_CopyButton(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithLocale("fr-FR"),
	)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `data-tooltip="Copier"`) {
		t.Error("copy button should use French tooltip")
	}
	if !strings.Contains(html, `data-copied="Copié !"`) {
		t.Error("copy button should use French success text")
	}
}

func TestRender_Locale_FullscreenButton(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithLocale("ja-JP"),
	)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `aria-label="全画面"`) {
		t.Error("fullscreen button should use Japanese label")
	}
}

func TestRender_Locale_ThresholdOverlay(t *testing.T) {
	longCode := strings.Repeat("line\n", 20)
	hl := &mockHighlighter{
		lightTokens: func() [][]Token {
			lines := make([][]Token, 20)
			for i := range lines {
				lines[i] = []Token{{Content: "line", Color: "#333"}}
			}
			return lines
		}(),
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithLocale("fr-FR"),
		WithCollapsible(CollapsibleConfig{LineThreshold: 5, PreviewLines: 3, DefaultCollapsed: true}),
	)
	html, err := engine.Render(longCode, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "Afficher plus") {
		t.Error("threshold overlay should use French expand text")
	}
	if !strings.Contains(html, "Afficher moins") {
		t.Error("threshold overlay should use French collapse text")
	}
}

func TestRender_Locale_CustomOverride(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithUIStrings(map[string]string{"copy.label": "Copy code"}),
	)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "kz-copy-btn") {
		t.Error("copy button should be present")
	}
}

func TestRender_Locale_DefaultEnglish(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `data-tooltip="Copy"`) {
		t.Error("default should use English copy tooltip")
	}
	if !strings.Contains(html, `aria-label="Fullscreen"`) {
		t.Error("default should use English fullscreen label")
	}
}

func TestRender_Locale_CollapseConfigOverridesLocale(t *testing.T) {
	longCode := strings.Repeat("line\n", 20)
	hl := &mockHighlighter{
		lightTokens: func() [][]Token {
			lines := make([][]Token, 20)
			for i := range lines {
				lines[i] = []Token{{Content: "line", Color: "#333"}}
			}
			return lines
		}(),
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithLocale("fr-FR"),
		WithCollapsible(CollapsibleConfig{
			LineThreshold:    5,
			PreviewLines:     3,
			DefaultCollapsed: true,
			ExpandButtonText: "Reveal",
		}),
	)
	html, err := engine.Render(longCode, Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "Reveal") {
		t.Error("CollapseConfig text should override locale")
	}
}

// --- File icon tests ---

func TestRender_FileIcons_Present(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "go", Title: "main.go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `data-ext="go"`) {
		t.Error("icon span should have data-ext for .go files")
	}
	if !strings.Contains(html, `class="kz-file-icon"`) {
		t.Error("icon span should have kz-file-icon class")
	}
}

func TestRender_FileIcons_CorrectExt(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "javascript", Title: "app.config.js"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `data-ext="js"`) {
		t.Error("should extract last extension (js), not config")
	}
}

func TestRender_FileIcons_NoTitle(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "kz-file-icon") {
		t.Error("no icon when no title is set")
	}
}

func TestRender_FileIcons_NoExtension(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "go", Title: "Makefile"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "kz-file-icon") {
		t.Error("no icon when title has no extension")
	}
}

func TestRender_FileIcons_Disabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithFileIcons(false))
	html, err := engine.Render("hello", Options{Lang: "go", Title: "main.go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "kz-file-icon") {
		t.Error("no icon when file icons disabled")
	}
}

func TestRender_FileIcons_TerminalFrame(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "echo hi", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("echo hi", Options{Lang: "bash", Title: "terminal.sh"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "kz-file-icon") {
		t.Error("no icon for terminal frames")
	}
}

func TestRender_FileIcons_CustomResolver(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithFileIconResolver(func(ext string) string {
		return fmt.Sprintf(`<img class="custom-icon" src="/icons/%s.svg"/>`, ext)
	}))
	html, err := engine.Render("hello", Options{Lang: "go", Title: "main.go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `<img class="custom-icon" src="/icons/go.svg"/>`) {
		t.Error("custom resolver output should be used")
	}
	if strings.Contains(html, "kz-file-icon") {
		t.Error("default span should not appear when resolver is set")
	}
}

func TestCSS_FileIcons_VarsPresent(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	css := engine.CSS()
	if !strings.Contains(css, "--kz-file-icon-size") {
		t.Error("file icon CSS vars should be present when enabled")
	}
}

func TestCSS_FileIcons_VarsAbsent(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithFileIcons(false))
	css := engine.CSS()
	if strings.Contains(css, "--kz-file-icon-size") {
		t.Error("file icon CSS vars should not be present when disabled")
	}
}

// --- Word wrap tests ---

func TestRender_Wrap_ClassEmitted(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: true}))
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `<pre class="wrap"`) {
		t.Error("pre should carry wrap class when Wrap is enabled")
	}
}

func TestRender_Wrap_AbsentByDefault(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("hello", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, `class="wrap"`) {
		t.Error("pre should not carry wrap class by default")
	}
}

func TestRenderWithMeta_Wrap(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "hello", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.RenderWithMeta("hello", "go wrap")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, `<pre class="wrap"`) {
		t.Error("wrap meta option should emit wrap class on pre")
	}
}

func TestRender_Wrap_IndentVarOnIndentedLine(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "    foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: true}))
	html, err := engine.Render("    foo()", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `style="--kz-indent:4ch"`) {
		t.Errorf("indented line should carry --kz-indent:4ch, got: %s", html)
	}
	if !strings.Contains(html, `<span class="indent">    </span>`) {
		t.Error("leading whitespace should render inside span.indent")
	}
	if !strings.Contains(html, ">foo()</span>") {
		t.Error("trimmed token content should follow the indent span")
	}
}

func TestRender_Wrap_NoVarOnUnindentedLine(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: true}))
	html, err := engine.Render("foo()", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, "--kz-indent") {
		t.Error("unindented line should not carry --kz-indent")
	}
	if strings.Contains(html, `class="indent"`) {
		t.Error("unindented line should not emit span.indent")
	}
}

func TestRender_Wrap_HangingIndentAdds(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "  foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: true, HangingIndent: 3}))
	html, err := engine.Render("  foo()", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `style="--kz-indent:5ch"`) {
		t.Errorf("hanging indent should add to preserved indent (2+3=5), got: %s", html)
	}
}

func TestRender_Wrap_HangingOnlyWhenPreserveDisabled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "    foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: false, HangingIndent: 2}))
	html, err := engine.Render("    foo()", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `style="--kz-indent:2ch"`) {
		t.Errorf("PreserveIndent=false should use hanging indent only, got: %s", html)
	}
}

func TestRender_Wrap_WithLineNumbers(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "  foo()", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithDefaults(BlockDefaults{Wrap: true, PreserveIndent: true, LineNumbers: true}))
	html, err := engine.Render("  foo()", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `<pre class="wrap"`) {
		t.Error("wrap class should coexist with line numbers")
	}
	if !strings.Contains(html, `class="kz-gutter"`) {
		t.Error("line number gutter should render in wrap mode")
	}
	if !strings.Contains(html, `style="--kz-indent:2ch"`) {
		t.Error("indent var should render alongside line numbers")
	}
}

func TestCSS_WrapStyles(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false))
	css := engine.CSS()
	if !strings.Contains(css, "pre.wrap .kz-line .kz-code") {
		t.Error("CSS should contain wrap rules")
	}
	if !strings.Contains(css, "text-indent: calc(var(--kz-indent, 0ch) * -1)") {
		t.Error("CSS should contain negative text-indent for hanging alignment")
	}
	if !strings.Contains(css, ".indent") {
		t.Error("CSS should contain span.indent rule")
	}
}

// --- Token background tests ---

func TestCSS_TokenBackground_LightRule(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false))
	css := engine.CSS()
	if !strings.Contains(css, `span[style^="--"] { color: var(--sl, inherit); background-color: var(--slbg, transparent);`) {
		t.Error("CSS should contain unconditional light token properties rule with starts-with selector")
	}
}

func TestCSS_TokenBackground_DarkRulePerStrategy(t *testing.T) {
	darkRule := `span[style^="--"] { color: var(--sd, inherit); background-color: var(--sdbg, transparent);`

	t.Run("selector", func(t *testing.T) {
		hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
		engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"),
			WithDarkMode(SelectorMode(".dark")), WithMinify(false))
		if !strings.Contains(engine.CSS(), ".dark .kazari-block .kz-line "+darkRule) {
			t.Error("selector mode should scope dark token rule under .dark")
		}
	})

	t.Run("media", func(t *testing.T) {
		hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
		engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"),
			WithDarkMode(MediaQueryMode()), WithMinify(false))
		css := engine.CSS()
		if !strings.Contains(css, ".kazari-block .kz-line "+darkRule) {
			t.Error("media mode should emit dark token rule")
		}
	})

	t.Run("both", func(t *testing.T) {
		hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
		engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"),
			WithDarkMode(BothMode(".dark")), WithMinify(false))
		css := engine.CSS()
		if !strings.Contains(css, ".dark .kazari-block .kz-line "+darkRule) {
			t.Error("both mode should scope dark token rule under .dark")
		}
		if strings.Count(css, darkRule) < 2 {
			t.Error("both mode should emit dark token rule in media query and under selector")
		}
	})
}

func TestRender_TokenBackground_EndToEnd(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "FAIL", Color: "#ffffff", BgColor: "#ff0000"}}},
		darkTokens:  [][]Token{{{Content: "FAIL", Color: "#ffffff", BgColor: "#cc0000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"), WithMinify(false))
	html, err := engine.Render("FAIL", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, "--slbg:#ff0000") {
		t.Error("token light background should render as --slbg")
	}
	if !strings.Contains(html, "--sdbg:#cc0000") {
		t.Error("token dark background should render as --sdbg")
	}
}

// --- Unknown language fallback tests ---

type erroringHighlighter struct {
	mockHighlighter
}

func (e *erroringHighlighter) Tokenize(code, lang, theme string) ([][]Token, error) {
	return nil, fmt.Errorf("tokenize failed for %s", lang)
}

func TestRender_UnknownLanguage_PlaintextFallback(t *testing.T) {
	hl := &erroringHighlighter{
		mockHighlighter: mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}},
	}
	var warnings []string
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithWarningHandler(func(msg string) { warnings = append(warnings, msg) }),
	)
	html, err := engine.Render("some code", Options{Lang: "madeuplang"})
	if err != nil {
		t.Fatalf("unknown language should fall back to plaintext, got error: %v", err)
	}
	if !strings.Contains(html, "some code") {
		t.Error("plaintext fallback should render the raw code")
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
}

func TestRender_UnknownLanguage_WarningMessage(t *testing.T) {
	hl := &erroringHighlighter{
		mockHighlighter: mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}},
	}
	var got string
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithWarningHandler(func(msg string) { got = msg }),
	)
	if _, err := engine.Render("x", Options{Lang: "madeuplang"}); err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(got, `"madeuplang"`) {
		t.Errorf("warning should name the unknown language, got %q", got)
	}
	if !strings.Contains(got, "plaintext") {
		t.Errorf("warning should mention plaintext fallback, got %q", got)
	}
}

func TestRender_KnownLanguage_ErrorPropagates(t *testing.T) {
	hl := &erroringHighlighter{
		mockHighlighter: mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithWarningHandler(func(string) {}),
	)
	// "go" is in mockHighlighter.GetLoadedLanguages, so the error is real.
	if _, err := engine.Render("x", Options{Lang: "go"}); err == nil {
		t.Fatal("errors for known languages must propagate")
	}
}

func TestRender_UnknownLanguage_DefaultLogPath(t *testing.T) {
	hl := &erroringHighlighter{
		mockHighlighter: mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}},
	}
	engine := New(WithHighlighter(hl), WithThemes("light-theme", ""), WithMinify(false))
	html, err := engine.Render("safe", Options{Lang: "madeuplang"})
	if err != nil {
		t.Fatalf("default log path should not error: %v", err)
	}
	if !strings.Contains(html, "safe") {
		t.Error("plaintext fallback should render code on default log path")
	}
}

// --- Terminal sr-only label tests ---

func TestRender_TerminalSRLabel_Untitled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "npm install", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("npm install", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `<span class="sr-only">Terminal window</span>`) {
		t.Error("untitled terminal frame should carry sr-only label")
	}
}

func TestRender_TerminalSRLabel_AbsentWhenTitled(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "npm install", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("npm install", Options{Lang: "bash", Title: "deploy.sh"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(html, `class="sr-only"`) {
		t.Error("titled terminal frame should not carry sr-only label")
	}
}

func TestRender_TerminalSRLabel_Localized(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "npm install", Color: "#333333"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLocale("fr-FR"))
	html, err := engine.Render("npm install", Options{Lang: "bash"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(html, `<span class="sr-only">Fenêtre de terminal</span>`) {
		t.Error("sr-only label should localize")
	}
}

// --- Copy button visibility tests ---

func TestCSS_CopyButton_AlwaysVisible(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCopyButton(true))
	css := engine.CSS()
	if strings.Contains(css, "--kz-copy-idle-opacity") {
		t.Error("copy CSS should not use idle opacity hiding")
	}
	if strings.Contains(css, "@media (hover: hover)") && strings.Contains(css, "kz-copy-btn") {
		t.Error("copy button should not be gated behind hover media query")
	}
}

func TestCSS_CopyButton_RTLIsolation(t *testing.T) {
	hl := &mockHighlighter{themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"}}
	engine := New(WithHighlighter(hl), WithMinify(false), WithCopyButton(true))
	css := engine.CSS()
	if !strings.Contains(css, "direction: ltr") {
		t.Error("copy button should force ltr direction")
	}
	if !strings.Contains(css, "unicode-bidi: isolate") {
		t.Error("copy button should isolate bidi context")
	}
}

// --- Fold background tests ---

func TestCSS_FoldBG_DerivedCollapseVars(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff", FoldBG: "#54aeff"},
	}
	engine := newTestEngine(hl, WithCollapsible(CollapsibleConfig{LineThreshold: 20, PreviewLines: 10}))
	css := engine.CSS()
	if !strings.Contains(css, "--kz-collapse-closed-bg: #54aeff33") {
		t.Error("foldBG should derive --kz-collapse-closed-bg at alpha 0.2")
	}
	if !strings.Contains(css, "--kz-collapse-closed-border: #54aeff80") {
		t.Error("foldBG should derive --kz-collapse-closed-border at alpha 0.5")
	}
}

func TestCSS_FoldBG_AbsentUsesStaticDefaults(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithCollapsible(CollapsibleConfig{LineThreshold: 20, PreviewLines: 10}))
	css := engine.CSS()
	if !strings.Contains(css, "--kz-collapse-closed-bg: rgb(84 174 255 / 20%)") {
		t.Error("static collapse defaults should remain when foldBG is absent")
	}
	if strings.Count(css, "--kz-collapse-closed-bg:") != 1 {
		t.Error("no derived collapse vars should be emitted without foldBG")
	}
}

func TestCSS_FoldBG_CustomizerCanModify(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff", FoldBG: "#54aeff"},
	}
	engine := newTestEngine(hl,
		WithCollapsible(CollapsibleConfig{LineThreshold: 20, PreviewLines: 10}),
		WithThemeCustomizer(func(name string, ti ThemeInfo) ThemeInfo {
			ti.FoldBG = "#ff0000"
			return ti
		}),
	)
	css := engine.CSS()
	if !strings.Contains(css, "--kz-collapse-closed-bg: #ff000033") {
		t.Error("customizer should be able to override FoldBG before CSS generation")
	}
}

// --- Theme adjustments tests ---

func TestThemeAdjustments_HueAppliedToBG(t *testing.T) {
	hue := 145.0
	chroma := 0.05
	ti := applyThemeAdjustments(
		ThemeInfo{BG: "#3366cc", FG: "#ffffff"},
		&ThemeAdjustments{Hue: &hue, Chroma: &chroma},
	)
	if ti.BG == "#3366cc" {
		t.Fatal("BG should be tinted")
	}
	_, _, h, err := color.ToOKLCH(ti.BG)
	if err != nil {
		t.Fatalf("tinted BG %q does not parse: %v", ti.BG, err)
	}
	if h < 143 || h > 147 {
		t.Errorf("tinted BG hue = %f, want ~145", h)
	}
	if ti.FG != "#ffffff" {
		t.Error("default targets should leave foregrounds unchanged")
	}
}

func TestThemeAdjustments_ChromaOnly(t *testing.T) {
	chroma := 0.0
	ti := applyThemeAdjustments(
		ThemeInfo{BG: "#3366cc"},
		&ThemeAdjustments{Chroma: &chroma},
	)
	lBefore, _, _, _ := color.ToOKLCH("#3366cc")
	lAfter, c, _, err := color.ToOKLCH(ti.BG)
	if err != nil {
		t.Fatalf("tinted BG %q does not parse: %v", ti.BG, err)
	}
	if c > 0.005 {
		t.Errorf("chroma 0 should desaturate, got chroma %f", c)
	}
	if lAfter < lBefore-0.01 || lAfter > lBefore+0.01 {
		t.Errorf("lightness should be preserved: %f -> %f", lBefore, lAfter)
	}
}

func TestThemeAdjustments_TargetsSelective(t *testing.T) {
	hue := 30.0
	chroma := 0.08
	ti := applyThemeAdjustments(
		ThemeInfo{BG: "#1e1e2e", FG: "#3366cc", LineNumberFG: "#3366cc"},
		&ThemeAdjustments{Hue: &hue, Chroma: &chroma, Targets: AdjustForegrounds},
	)
	if ti.BG != "#1e1e2e" {
		t.Error("AdjustForegrounds should leave BG unchanged")
	}
	if ti.FG == "#3366cc" {
		t.Error("AdjustForegrounds should tint FG")
	}
	if ti.LineNumberFG == "#3366cc" {
		t.Error("AdjustForegrounds should tint LineNumberFG")
	}
}

func TestThemeAdjustments_NilFieldsUnchanged(t *testing.T) {
	in := ThemeInfo{BG: "#1e1e2e", FG: "#cdd6f4", SelectionBG: "#45475a"}
	if out := applyThemeAdjustments(in, &ThemeAdjustments{}); out != in {
		t.Error("adjustments with nil Hue and Chroma should be a no-op")
	}
	if out := applyThemeAdjustments(in, nil); out != in {
		t.Error("nil adjustments should be a no-op")
	}
}

func TestThemeAdjustments_CustomizerRunsAfter(t *testing.T) {
	hue := 145.0
	chroma := 0.05
	var customizerSawBG string
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#3366cc"},
	}
	engine := newTestEngine(hl,
		WithThemeAdjustments(ThemeAdjustments{Hue: &hue, Chroma: &chroma}),
		WithThemeCustomizer(func(name string, ti ThemeInfo) ThemeInfo {
			customizerSawBG = ti.BG
			ti.BG = "#123456"
			return ti
		}),
	)
	if customizerSawBG == "" || customizerSawBG == "#3366cc" {
		t.Errorf("customizer should receive the adjusted BG, saw %q", customizerSawBG)
	}
	if !strings.Contains(engine.CSS(), "--kz-editor-bg: #123456") {
		t.Error("customizer must get final say over adjusted colors")
	}
}

func TestRender_ThemeAdjustments_EndToEndCSS(t *testing.T) {
	hue := 250.0
	chroma := 0.03
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#222222"},
	}
	engine := newTestEngine(hl, WithThemeAdjustments(ThemeAdjustments{Hue: &hue, Chroma: &chroma}))
	css := engine.CSS()
	if strings.Contains(css, "--kz-editor-bg: #222222") {
		t.Error("editor BG should be tinted in generated CSS")
	}
	if !strings.Contains(css, "--kz-editor-fg: #24292f") {
		t.Error("editor FG should be untouched by default targets")
	}
}

// --- CreateInlineSVGURL tests ---

func TestCreateInlineSVGURL(t *testing.T) {
	got := CreateInlineSVGURL("<svg viewBox='0 0 1 1'/>")
	want := "data:image/svg+xml,%3Csvg viewBox='0 0 1 1'/%3E"
	if got != want {
		t.Errorf("CreateInlineSVGURL():\ngot:  %s\nwant: %s", got, want)
	}
}

// --- Per-block theme override visual theming tests ---

func overrideMock() *mockHighlighter {
	return &mockHighlighter{
		lightTokens:   [][]Token{{{Content: "x", Color: "#aaa"}}},
		themeInfo:     ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo: ThemeInfo{FG: "#f8f8f2", BG: "#282a36"},
	}
}

func TestThemeOverride_EmitsThemedWrapper(t *testing.T) {
	engine := newTestEngine(overrideMock())
	html, err := engine.RenderWithMeta("x", `go theme="dracula"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "kz-themed") {
		t.Error("override block should carry the kz-themed class")
	}
	if !strings.Contains(html, "--kz-ovl-editor-bg:#282a36") {
		t.Error("override block should carry the override background inline")
	}
	if !strings.Contains(html, "--kz-ovl-editor-fg:#f8f8f2") {
		t.Error("override block should carry the override foreground inline")
	}
	if strings.Contains(html, "--kz-ovd-") {
		t.Error("single theme engine should not emit dark slot vars")
	}
}

func TestThemeOverride_NoOverride_NoThemedClass(t *testing.T) {
	engine := newTestEngine(overrideMock())
	html, err := engine.RenderWithMeta("x", "go")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if strings.Contains(html, "kz-themed") {
		t.Error("blocks without an override should not carry kz-themed")
	}
}

func TestThemeOverride_SameAsEngineThemes_NoOp(t *testing.T) {
	engine := newTestEngine(overrideMock())
	html, err := engine.RenderWithMeta("x", `go theme="light-theme"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if strings.Contains(html, "kz-themed") {
		t.Error("override matching the engine themes should be a no-op")
	}
}

func TestThemeOverride_DualThemeEngine_EmitsBothSlots(t *testing.T) {
	hl := overrideMock()
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
	)
	html, err := engine.RenderWithMeta("x", `go theme="dracula"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "--kz-ovl-editor-bg:#282a36") {
		t.Error("should emit light slot override vars")
	}
	if !strings.Contains(html, "--kz-ovd-editor-bg:#282a36") {
		t.Error("dual theme engine should emit dark slot override vars")
	}
}

func TestThemeOverride_PartialCommaForm_NoDarkSlot(t *testing.T) {
	hl := overrideMock()
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
	)
	html, err := engine.RenderWithMeta("x", `go theme="dracula,"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if !strings.Contains(html, "--kz-ovl-editor-bg:#282a36") {
		t.Error("light slot vars should be present")
	}
	if strings.Contains(html, "--kz-ovd-") {
		t.Error("empty dark half of a comma override should emit no dark slot vars")
	}
}

func TestThemeOverride_SingleValueAppliesToDarkTokenization(t *testing.T) {
	hl := &dualMockHighlighter{
		mockHighlighter: mockHighlighter{
			lightTokens:   [][]Token{{{Content: "x", Color: "#aaa"}}},
			darkTokens:    [][]Token{{{Content: "x", Color: "#bbb"}}},
			themeInfo:     ThemeInfo{FG: "#24292f", BG: "#ffffff"},
			darkThemeInfo: ThemeInfo{FG: "#f8f8f2", BG: "#282a36"},
		},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
	)
	if _, err := engine.RenderWithMeta("x", `go theme="dracula"`); err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if hl.lastLight != "dracula" || hl.lastDark != "dracula" {
		t.Errorf("single value override should tokenize both modes with the override theme, got light=%q dark=%q", hl.lastLight, hl.lastDark)
	}
}

func TestThemeOverride_SingleThemePage_StaysSingle(t *testing.T) {
	hl := &dualMockHighlighter{
		mockHighlighter: mockHighlighter{
			lightTokens:   [][]Token{{{Content: "x", Color: "#aaa"}}},
			darkTokens:    [][]Token{{{Content: "x", Color: "#bbb"}}},
			themeInfo:     ThemeInfo{FG: "#24292f", BG: "#ffffff"},
			darkThemeInfo: ThemeInfo{FG: "#f8f8f2", BG: "#282a36"},
		},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)
	if _, err := engine.RenderWithMeta("x", `go theme="dracula"`); err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if hl.dualCalls != 0 {
		t.Error("single theme engine must not switch to dual tokenization on override")
	}
	if hl.tokenizeCalls == 0 {
		t.Error("expected single theme tokenization")
	}
}

type erroringThemeHighlighter struct {
	mockHighlighter
}

func (e *erroringThemeHighlighter) GetThemeColors(theme string) (ThemeInfo, error) {
	if theme != "light-theme" {
		return ThemeInfo{}, fmt.Errorf("unknown theme %q", theme)
	}
	return e.mockHighlighter.GetThemeColors(theme)
}

func TestThemeOverride_UnknownThemeWarnsAndSkips(t *testing.T) {
	hl := &erroringThemeHighlighter{
		mockHighlighter: mockHighlighter{
			lightTokens: [][]Token{{{Content: "x", Color: "#aaa"}}},
			themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		},
	}
	var warnings []string
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithWarningHandler(func(msg string) { warnings = append(warnings, msg) }),
	)

	html, err := engine.RenderWithMeta("x", `go theme="nope"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	if strings.Contains(html, "kz-themed") {
		t.Error("unknown override theme should not produce a themed wrapper")
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d: %v", len(warnings), warnings)
	}

	if _, err := engine.RenderWithMeta("x", `go theme="nope"`); err != nil {
		t.Fatalf("second render error: %v", err)
	}
	if len(warnings) != 1 {
		t.Errorf("cached failure should not warn again, got %d warnings", len(warnings))
	}
}

func TestThemeOverride_MarkerContrastUsesOverrideBG(t *testing.T) {
	hl := overrideMock()
	hl.lightTokens = [][]Token{{{Content: "x", Color: "#777777"}}}
	engine := newTestEngine(hl, WithMinContrast(5.5))

	extractSL := func(html string) string {
		idx := strings.Index(html, "--sl:")
		if idx < 0 {
			t.Fatal("no --sl token style found")
		}
		end := strings.IndexAny(html[idx:], ";\"")
		return html[idx : idx+end]
	}

	pageHTML, err := engine.RenderWithMeta("x", "go {1}")
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}
	overrideHTML, err := engine.RenderWithMeta("x", `go {1} theme="dracula"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if extractSL(pageHTML) == extractSL(overrideHTML) {
		t.Errorf("marker contrast should adjust against the override background, both produced %s", extractSL(pageHTML))
	}
}

func TestThemeOverride_ConcurrentRenders(t *testing.T) {
	hl := overrideMock()
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
	)

	var wg sync.WaitGroup
	metas := []string{`go theme="dracula"`, `go theme="nord"`, `go theme="a,b"`, "go"}
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(meta string) {
			defer wg.Done()
			if _, err := engine.RenderWithMeta("x", meta); err != nil {
				t.Errorf("concurrent render error: %v", err)
			}
		}(metas[i%len(metas)])
	}
	wg.Wait()
}

// --- normalizeVarName ---

func TestNormalizeVarName(t *testing.T) {
	tests := map[string]struct{ input, want string }{
		"bare name":       {"radius", "--kz-radius"},
		"already prefixed": {"--kz-shadow", "--kz-shadow"},
		"custom prefix":   {"--custom-var", "--custom-var"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := normalizeVarName(tc.input); got != tc.want {
				t.Errorf("normalizeVarName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// --- WithStyleOverrides merge ---

func TestWithStyleOverrides_Merge(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
		WithStyleOverrides(map[string]string{
			"radius": "0.5rem",
			"shadow": "none",
		}),
		WithStyleOverrides(map[string]string{
			"shadow":    "0 1px 4px rgba(0,0,0,0.1)",
			"font-size": "1rem",
		}),
	)
	css := engine.CSS()

	if !strings.Contains(css, "--kz-radius: 0.5rem") {
		t.Error("first call's unique key should be preserved")
	}
	if !strings.Contains(css, "--kz-shadow: 0 1px 4px rgba(0,0,0,0.1)") {
		t.Error("later call should win for duplicate key")
	}
	if strings.Contains(css, "--kz-shadow: none") {
		t.Error("earlier value for duplicate key should be overwritten")
	}
	if !strings.Contains(css, "--kz-font-size: 1rem") {
		t.Error("second call's unique key should be present")
	}
}

// --- WithThemedStyleOverrides ---

func TestWithThemedStyleOverrides(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		darkThemeInfo: ThemeInfo{FG: "#d4d4d4", BG: "#1e1e1e"},
	}
	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
		WithThemedStyleOverrides(map[string]StyleValue{
			"shadow": {Dark: "none", Light: "0 2px 8px rgba(0,0,0,0.1)"},
		}),
	)
	css := engine.CSS()

	if !strings.Contains(css, "--kz-shadow: 0 2px 8px rgba(0,0,0,0.1)") {
		t.Error("light value should appear in CSS")
	}
	if !strings.Contains(css, "--kz-shadow: none") {
		t.Error("dark value should appear in CSS")
	}
}

// --- Deep merge composition ---

func TestWithLanguageDefaults_Merge(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := newTestEngine(hl,
		WithLanguageDefaults(map[string]BlockDefaults{
			"go":     {LineNumbers: true},
			"python": {Wrap: true},
		}),
		WithLanguageDefaults(map[string]BlockDefaults{
			"python": {LineNumbers: true, Wrap: false},
			"rust":   {LineNumbers: true},
		}),
	)

	cfg := engine.cfg

	if _, ok := cfg.LanguageDefaults["go"]; !ok {
		t.Error("first call's unique key 'go' should be preserved")
	}
	if _, ok := cfg.LanguageDefaults["rust"]; !ok {
		t.Error("second call's unique key 'rust' should be present")
	}
	py := cfg.LanguageDefaults["python"]
	if !py.LineNumbers {
		t.Error("later call should win: python.LineNumbers should be true")
	}
	if py.Wrap {
		t.Error("later call should win: python.Wrap should be false")
	}
}

func TestWithLanguageAliases_Merge(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := newTestEngine(hl,
		WithLanguageAliases(map[string]string{
			"js": "javascript",
			"ts": "typescript",
		}),
		WithLanguageAliases(map[string]string{
			"ts": "tsx",
			"py": "python",
		}),
	)

	cfg := engine.cfg

	if cfg.LanguageAliases["js"] != "javascript" {
		t.Error("first call's unique key 'js' should be preserved")
	}
	if cfg.LanguageAliases["py"] != "python" {
		t.Error("second call's unique key 'py' should be present")
	}
	if cfg.LanguageAliases["ts"] != "tsx" {
		t.Error("later call should win for duplicate key 'ts'")
	}
}

func TestWithUIStrings_Merge(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := newTestEngine(hl,
		WithUIStrings(map[string]string{
			"copy.label":   "Copy",
			"copy.success": "Copied!",
		}),
		WithUIStrings(map[string]string{
			"copy.success":    "Done!",
			"expand.label":    "Expand",
		}),
	)

	cfg := engine.cfg

	if cfg.UIStringOverrides["copy.label"] != "Copy" {
		t.Error("first call's unique key should be preserved")
	}
	if cfg.UIStringOverrides["expand.label"] != "Expand" {
		t.Error("second call's unique key should be present")
	}
	if cfg.UIStringOverrides["copy.success"] != "Done!" {
		t.Error("later call should win for duplicate key")
	}
}

// --- Language icon mode ---

func TestRender_LangIcon_DefaultTextOnly(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `<span class="kz-lang">Go</span>`) {
		t.Error("default mode should emit text badge")
	}
	if strings.Contains(html, "kz-lang-icon") {
		t.Error("default mode should not emit icon slot")
	}
}

func TestRender_LangIcon_IconOnly(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLanguageIconMode(LangIconOnly))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `data-lang="go"`) {
		t.Error("icon-only mode should emit data-lang attribute")
	}
	if !strings.Contains(html, `kz-lang-icon`) {
		t.Error("icon-only mode should emit icon slot")
	}
	if strings.Contains(html, `<span class="kz-lang">`) {
		t.Error("icon-only mode should not emit text badge")
	}
}

func TestRender_LangIcon_IconAndText(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLanguageIconMode(LangIconAndText))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `kz-lang-icon`) {
		t.Error("icon+text mode should emit icon slot")
	}
	if !strings.Contains(html, `<span class="kz-lang">Go</span>`) {
		t.Error("icon+text mode should emit text badge")
	}
	iconIdx := strings.Index(html, "kz-lang-icon")
	textIdx := strings.Index(html, "kz-lang\">")
	if iconIdx > textIdx {
		t.Error("icon should appear before text")
	}
}

func TestRender_LangIcon_BadgeDisabledSuppressesBoth(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{{{Content: "x", Color: "#000"}}},
		themeInfo:   ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLanguageIconMode(LangIconOnly), WithLanguageBadge(false))
	html, err := engine.Render("x", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "kz-lang-icon") {
		t.Error("LanguageBadge=false should suppress icon")
	}
	if strings.Contains(html, "kz-lang") {
		t.Error("LanguageBadge=false should suppress text")
	}
}

// --- Links tests ---

func TestRender_Links_BasicLink(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "import fmt", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLinks(true))
	result, err := engine.Render("import @[fmt](https://pkg.go.dev/fmt)", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, `<a class="kz-link" href="https://pkg.go.dev/fmt" rel="noopener noreferrer">`) {
		t.Errorf("expected kz-link anchor, got: %s", result)
	}
	if !strings.Contains(result, "fmt</a>") && !strings.Contains(result, "fmt</span></a>") {
		t.Errorf("expected link text 'fmt' inside anchor, got: %s", result)
	}
}

func TestRender_Links_DisabledByDefault(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "import @[fmt](https://pkg.go.dev/fmt)", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	result, err := engine.Render("import @[fmt](https://pkg.go.dev/fmt)", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, "kz-link") {
		t.Errorf("links should not render when feature is disabled, got: %s", result)
	}
}

func TestRender_Links_StrippedFromRawCode(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "import fmt", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLinks(true), WithCopyButton(true))
	result, err := engine.Render("import @[fmt](https://pkg.go.dev/fmt)", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, `data-code="import @[fmt]`) {
		t.Errorf("copy button data should not contain link syntax, got: %s", result)
	}
	if !strings.Contains(result, `data-code="import fmt"`) {
		t.Errorf("copy button should contain cleaned code, got: %s", result)
	}
}

func TestRender_Links_MultiTokenSpan(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{
				{Content: "fmt", Color: "#000"},
				{Content: ".", Color: "#111"},
				{Content: "Println", Color: "#222"},
			},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLinks(true))
	result, err := engine.Render("@[fmt.Println](https://pkg.go.dev/fmt#Println)", Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	count := strings.Count(result, "kz-link")
	if count < 2 {
		t.Errorf("expected multiple kz-link anchors for multi-token link, got %d in: %s", count, result)
	}
}

func TestRender_Links_HrefEscaped(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "click here", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLinks(true))
	result, err := engine.Render(`@[click here](https://example.com/a&b<c>d"e)`, Options{Lang: "go"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result, `href="https://example.com/a&b<`) {
		t.Errorf("href should be HTML-escaped, got: %s", result)
	}
	if !strings.Contains(result, "&amp;") {
		t.Errorf("expected HTML entity in href, got: %s", result)
	}
}

func TestRender_Links_CSS_Included(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl, WithLinks(true))
	css := engine.CSS()
	if !strings.Contains(css, ".kz-link") {
		t.Errorf("CSS should include .kz-link rules when links enabled")
	}
}

func TestRender_Links_CSS_ExcludedWhenDisabled(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	engine := newTestEngine(hl)
	css := engine.CSS()
	if strings.Contains(css, ".kz-link") {
		t.Errorf("CSS should not include .kz-link rules when links disabled")
	}
}
