package kazari

import (
	"strings"
	"testing"
)

type mockHighlighter struct {
	lightTokens [][]Token
	darkTokens  [][]Token
	themeInfo   ThemeInfo
}

func (m *mockHighlighter) Tokenize(code, lang, theme string) ([][]Token, error) {
	if theme == "dark-theme" && m.darkTokens != nil {
		return m.darkTokens, nil
	}
	return m.lightTokens, nil
}

func (m *mockHighlighter) GetThemeColors(theme string) (ThemeInfo, error) {
	return m.themeInfo, nil
}

func (m *mockHighlighter) GetLoadedLanguages() []string {
	return []string{"go", "javascript", "bash"}
}

func newTestEngine(hl *mockHighlighter, opts ...Option) *Engine {
	base := []Option{
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
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

	engine := New(WithHighlighter(hl), WithThemes("light-theme", "dark-theme"), WithMinify(false))
	html, err := engine.Render("func main", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, `<div class="kazari-code">`) {
		t.Error("missing kazari-code wrapper")
	}
	if !strings.Contains(html, "--sl:#cf222e") {
		t.Error("missing light color")
	}
	if !strings.Contains(html, "--sd:#ff7b72") {
		t.Error("missing dark color")
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
	// Badge should be in right section when title is set
	if !strings.Contains(html, `kz-toolbar-right`) {
		t.Error("missing toolbar right")
	}
	if !strings.Contains(html, `<span class="kz-lang">Go</span>`) {
		t.Error("missing language badge in right section")
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
	if !strings.Contains(html, "<span>Copy</span>") {
		t.Error("missing copy text label")
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
	engine := New(WithHighlighter(hl), WithCopyButton(false), WithFullscreenButton(false))

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

	if strings.Contains(html, "gutter") {
		t.Error("gutter should not be present when line numbers disabled")
	}
	if strings.Contains(html, `class="ln"`) {
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

	if !strings.Contains(html, `<div class="gutter">`) {
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
	html, err := engine.Render("a\nb", Options{Lang: "go", LineNumbers: &ln, StartLineNumber: 42})
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
	if !strings.Contains(html, `<div class="gutter">`) {
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

	if !strings.Contains(html, `<div class="gutter">`) {
		t.Error("missing gutter from meta showLineNumbers")
	}
	if !strings.Contains(html, `aria-hidden="true">10</div>`) {
		t.Error("missing line number 10 from meta startLineNumber=10")
	}
	if !strings.Contains(html, `aria-hidden="true">11</div>`) {
		t.Error("missing line number 11")
	}
}
