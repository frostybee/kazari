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
	return []string{"go", "javascript"}
}

func TestNew(t *testing.T) {
	engine := New(
		WithHighlighter(&mockHighlighter{
			themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		}),
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

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
	)

	html, err := engine.Render("func main", Options{Lang: "go"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	// Verify structure
	if !strings.Contains(html, `<div class="kazari-code">`) {
		t.Error("missing kazari-code wrapper")
	}
	if !strings.Contains(html, `data-language="go"`) {
		t.Error("missing data-language attribute")
	}
	if !strings.Contains(html, `<div class="kz-line">`) {
		t.Error("missing kz-line div")
	}

	// Verify dual-theme token colors
	if !strings.Contains(html, "--sl:#cf222e") {
		t.Error("missing light color for 'func'")
	}
	if !strings.Contains(html, "--sd:#ff7b72") {
		t.Error("missing dark color for 'func'")
	}
	if !strings.Contains(html, "--sl:#8250df") {
		t.Error("missing light color for ' main'")
	}
	if !strings.Contains(html, "--sd:#d2a8ff") {
		t.Error("missing dark color for ' main'")
	}
}

func TestRender_SingleTheme(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "hello", Color: "#333333"}},
		},
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)

	html, err := engine.Render("hello", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "--sd:") {
		t.Error("single-theme mode should not emit --sd")
	}
	if !strings.Contains(html, "--sl:#333333") {
		t.Error("missing light color")
	}
}

func TestRender_HTMLEscape(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "<div>&\"test\"</div>", Color: "#000000"}},
		},
		themeInfo: ThemeInfo{FG: "#000000", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)

	html, err := engine.Render(`<div>&"test"</div>`, Options{Lang: "html"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if strings.Contains(html, "<div>&") {
		t.Error("HTML content should be escaped")
	}
	if !strings.Contains(html, "&lt;div&gt;&amp;&#34;test&#34;&lt;/div&gt;") {
		t.Errorf("incorrect HTML escaping, got: %s", html)
	}
}

func TestRender_EmptyLine(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "a", Color: "#000"}},
			{}, // empty line
			{{Content: "b", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)

	html, err := engine.Render("a\n\nb", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	// All 3 lines should produce kz-line divs
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

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)

	html, err := engine.Render("italic\nbold\nboth", Options{Lang: "text"})
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(html, "--sfs:italic") {
		t.Error("missing italic font style")
	}
	if !strings.Contains(html, "--sfw:bold") {
		t.Error("missing bold font weight")
	}
	if !strings.Contains(html, "--std:underline") {
		t.Error("missing underline text decoration")
	}
}

func TestRenderWithMeta(t *testing.T) {
	hl := &mockHighlighter{
		lightTokens: [][]Token{
			{{Content: "x", Color: "#000"}},
		},
		themeInfo: ThemeInfo{FG: "#000", BG: "#fff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", ""),
		WithMinify(false),
	)

	html, err := engine.RenderWithMeta("x", `go title="main.go"`)
	if err != nil {
		t.Fatalf("RenderWithMeta() error: %v", err)
	}

	if !strings.Contains(html, `data-language="go"`) {
		t.Error("meta string language not parsed")
	}
}

func TestCSS_Deterministic(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
	)

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
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
		WithCascadeLayer("kazari"),
	)

	css := engine.CSS()
	if !strings.Contains(css, "@layer kazari") {
		t.Error("CSS should contain @layer wrapper")
	}
}

func TestCSS_NoLayer(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
		WithMinify(false),
		WithCascadeLayer(""),
	)

	css := engine.CSS()
	if strings.Contains(css, "@layer") {
		t.Error("CSS should not contain @layer when disabled")
	}
}

func TestAssets_Hashing(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
	)

	assets := engine.Assets()

	if assets.CSS.Hash == "" {
		t.Error("CSS hash should not be empty")
	}
	if len(assets.CSS.Hash) != 8 {
		t.Errorf("CSS hash should be 8 chars, got %d", len(assets.CSS.Hash))
	}
	if !strings.HasPrefix(assets.CSS.Filename, "kazari-") {
		t.Error("CSS filename should start with kazari-")
	}
	if !strings.HasSuffix(assets.CSS.Filename, ".css") {
		t.Error("CSS filename should end with .css")
	}
}

func TestJS_EmptyInPhase1(t *testing.T) {
	hl := &mockHighlighter{
		themeInfo: ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}

	engine := New(
		WithHighlighter(hl),
		WithThemes("light-theme", "dark-theme"),
	)

	js := engine.JS()
	if js != "" {
		t.Error("JS() should be empty in Phase 1")
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
	if !strings.Contains(html, `<div class="kazari-code">`) {
		t.Error("plaintext should still have wrapper structure")
	}
}
