package kazarimd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/frostybee/kazari"
	"github.com/yuin/goldmark"
)

type mockHighlighter struct {
	lightTokens [][]kazari.Token
	darkTokens  [][]kazari.Token
	themeInfo   kazari.ThemeInfo
}

func (m *mockHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
	if theme == "dark-theme" && m.darkTokens != nil {
		return m.darkTokens, nil
	}
	lines := strings.Split(code, "\n")
	tokens := make([][]kazari.Token, len(lines))
	for i, line := range lines {
		if line == "" {
			tokens[i] = []kazari.Token{{Content: "", Color: ""}}
		} else {
			tokens[i] = []kazari.Token{{Content: line, Color: "#333333"}}
		}
	}
	return tokens, nil
}

func (m *mockHighlighter) GetThemeColors(theme string) (kazari.ThemeInfo, error) {
	return m.themeInfo, nil
}

func (m *mockHighlighter) GetLoadedLanguages() []string {
	return []string{"go", "python", "javascript", "bash"}
}

func newTestEngine(opts ...kazari.Option) *kazari.Engine {
	hl := &mockHighlighter{
		themeInfo: kazari.ThemeInfo{FG: "#24292f", BG: "#ffffff"},
	}
	base := []kazari.Option{
		kazari.WithHighlighter(hl),
		kazari.WithThemes("light-theme", ""),
		kazari.WithMinify(false),
	}
	return kazari.New(append(base, opts...)...)
}

func renderMarkdown(t *testing.T, engine *kazari.Engine, md string, withCodeGroups bool) string {
	t.Helper()
	exts := []goldmark.Extender{New(engine)}
	if withCodeGroups {
		exts = append(exts, CodeGroups(engine))
	}
	parser := goldmark.New(goldmark.WithExtensions(exts...))
	var buf bytes.Buffer
	if err := parser.Convert([]byte(md), &buf); err != nil {
		t.Fatalf("goldmark.Convert() error: %v", err)
	}
	return buf.String()
}

// --- Basic extension tests ---

func TestBasicFencedCodeBlock(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "```go\nfunc main() {}\n```\n", false)

	if !strings.Contains(html, `class="kazari-code"`) {
		t.Error("missing kazari-code wrapper")
	}
	if !strings.Contains(html, `<figure`) {
		t.Error("missing figure element")
	}
	if !strings.Contains(html, "func main()") {
		t.Error("missing code content")
	}
}

func TestNoCodeBlocks(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "# Hello\n\nSome text.\n", false)

	if strings.Contains(html, "kazari-code") {
		t.Error("kazari-code should not appear when no code blocks present")
	}
	if !strings.Contains(html, "<h1>Hello</h1>") {
		t.Error("heading should still render")
	}
}

func TestLanguageExtraction(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "```go\npackage main\n```\n", false)

	if !strings.Contains(html, `data-lang="go"`) {
		t.Error("missing data-lang attribute")
	}
}

func TestMetaPassThrough_Title(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "```go title=\"main.go\"\npackage main\n```\n", false)

	if !strings.Contains(html, "main.go") {
		t.Error("title should appear in output")
	}
	if !strings.Contains(html, "has-title") {
		t.Error("frame should have has-title class")
	}
}

func TestNoLanguage(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "```\nhello world\n```\n", false)

	if !strings.Contains(html, `class="kazari-code"`) {
		t.Error("should still render with kazari-code wrapper")
	}
	if !strings.Contains(html, "hello world") {
		t.Error("code content missing")
	}
}

func TestEmptyCodeBlock(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "```go\n```\n", false)

	if !strings.Contains(html, `class="kazari-code"`) {
		t.Error("empty block should still render wrapper")
	}
}

func TestMultipleBlocks(t *testing.T) {
	engine := newTestEngine()
	md := "```go\nfunc main() {}\n```\n\n```python\ndef main():\n    pass\n```\n"
	html := renderMarkdown(t, engine, md, false)

	count := strings.Count(html, `class="kazari-code"`)
	if count != 2 {
		t.Errorf("expected 2 kazari-code blocks, got %d", count)
	}
}

func TestMetaPassThrough_LineNumbers(t *testing.T) {
	engine := newTestEngine()
	html := renderMarkdown(t, engine, "```go showLineNumbers\nfunc main() {}\n```\n", false)

	if !strings.Contains(html, `class="ln"`) {
		t.Error("line numbers should be rendered")
	}
}

// --- Code group tests ---

func TestCodeGroup_Basic(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go title=\"main.go\"\nfunc main() {}\n```\n\n```python title=\"main.py\"\ndef main():\n    pass\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, `class="kazari-code kz-group"`) {
		t.Error("missing kz-group wrapper")
	}
	if !strings.Contains(html, `role="tablist"`) {
		t.Error("missing tablist")
	}
	if strings.Count(html, `role="tab"`) != 2 {
		t.Errorf("expected 2 tabs, got %d", strings.Count(html, `role="tab"`))
	}
	if strings.Count(html, `role="tabpanel"`) != 2 {
		t.Errorf("expected 2 panels, got %d", strings.Count(html, `role="tabpanel"`))
	}
}

func TestCodeGroup_TabLabels_Title(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go title=\"main.go\"\nfunc main() {}\n```\n\n```python title=\"app.py\"\nprint()\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, ">main.go</button>") {
		t.Error("tab should use title 'main.go'")
	}
	if !strings.Contains(html, ">app.py</button>") {
		t.Error("tab should use title 'app.py'")
	}
}

func TestCodeGroup_TabLabels_Language(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go\nfunc main() {}\n```\n\n```python\nprint()\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, ">Go</button>") {
		t.Error("tab should fallback to capitalized language 'Go'")
	}
	if !strings.Contains(html, ">Python</button>") {
		t.Error("tab should fallback to capitalized language 'Python'")
	}
}

func TestCodeGroup_SingleBlock(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go title=\"main.go\"\nfunc main() {}\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, `class="kazari-code kz-group"`) {
		t.Error("single-block group should still render as group")
	}
	if strings.Count(html, `role="tab"`) != 1 {
		t.Errorf("expected 1 tab, got %d", strings.Count(html, `role="tab"`))
	}
}

func TestCodeGroup_NoBlocks(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\nSome text but no code.\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if strings.Contains(html, "kz-group") {
		t.Error("empty code group should be omitted")
	}
}

func TestCodeGroup_HiddenPanels(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go title=\"a.go\"\nfunc a() {}\n```\n\n```go title=\"b.go\"\nfunc b() {}\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if strings.Count(html, `hidden`) < 1 {
		t.Error("non-active panels should have hidden attribute")
	}
	// First panel should NOT be hidden.
	firstPanel := strings.Index(html, `role="tabpanel"`)
	hiddenPos := strings.Index(html, `hidden`)
	if hiddenPos < firstPanel {
		t.Error("first panel should not be hidden")
	}
}

func TestCodeGroup_ActiveTab(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go title=\"a.go\"\nfunc a() {}\n```\n\n```go title=\"b.go\"\nfunc b() {}\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, `aria-selected="true"`) {
		t.Error("first tab should have aria-selected=true")
	}
	if strings.Count(html, `aria-selected="true"`) != 1 {
		t.Errorf("only one tab should be selected, got %d", strings.Count(html, `aria-selected="true"`))
	}
}

func TestCodeGroup_StandaloneOutside(t *testing.T) {
	engine := newTestEngine()
	md := "```go\nfunc standalone() {}\n```\n\n:::code-group\n\n```go title=\"grouped.go\"\nfunc grouped() {}\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, "standalone") {
		t.Error("standalone block should render")
	}
	if !strings.Contains(html, "kz-group") {
		t.Error("code group should render")
	}
	// Standalone block should NOT be inside a group.
	standalonePos := strings.Index(html, "standalone")
	groupPos := strings.Index(html, "kz-group")
	if standalonePos > groupPos {
		t.Error("standalone block should appear before code group")
	}
}

func TestCodeGroup_RovingTabindex(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go title=\"a.go\"\nfunc a() {}\n```\n\n```go title=\"b.go\"\nfunc b() {}\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if strings.Count(html, `tabindex="0"`) != 1 {
		t.Error("only active tab should have tabindex=0")
	}
	if strings.Count(html, `tabindex="-1"`) < 1 {
		t.Error("inactive tabs should have tabindex=-1")
	}
}

// --- CSS/JS conditional tests ---

func TestCSS_IncludesCodeGroupStyles(t *testing.T) {
	engine := newTestEngine()
	CodeGroups(engine) // enables code groups
	css := engine.CSS()

	if !strings.Contains(css, "kz-group") {
		t.Error("CSS should include code group styles when enabled")
	}
}

func TestCSS_ExcludesCodeGroupStyles(t *testing.T) {
	engine := newTestEngine()
	css := engine.CSS()

	if strings.Contains(css, "kz-group") {
		t.Error("CSS should not include code group styles when disabled")
	}
}

func TestJS_IncludesCodeGroupHandler(t *testing.T) {
	engine := newTestEngine()
	CodeGroups(engine)
	js := engine.JS()

	if !strings.Contains(js, "kz-group") {
		t.Error("JS should include code group handler when enabled")
	}
}

func TestJS_ExcludesCodeGroupHandler(t *testing.T) {
	engine := newTestEngine()
	js := engine.JS()

	if strings.Contains(js, "kz-group") {
		t.Error("JS should not include code group handler when disabled")
	}
}

// --- Tab sync tests ---

func TestCodeGroup_SyncAttribute(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group sync=\"language\"\n\n```go\nfunc main() {}\n```\n\n```python\nprint()\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, `data-sync="language"`) {
		t.Error("missing data-sync attribute")
	}
}

func TestCodeGroup_SyncAttribute_SingleQuote(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group sync='runtime'\n\n```go\nfmt.Println()\n```\n\n```python\nprint()\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, `data-sync="runtime"`) {
		t.Error("missing data-sync attribute for single-quote syntax")
	}
}

func TestCodeGroup_NoSync(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group\n\n```go\nfunc main() {}\n```\n\n```python\nprint()\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if strings.Contains(html, "data-sync") {
		t.Error("data-sync should not be present without sync attribute")
	}
	if !strings.Contains(html, `class="kazari-code kz-group"`) {
		t.Error("code group should still render without sync")
	}
}

func TestCodeGroup_SyncAttribute_HTMLEscaped(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group sync=\"a&b\"\n\n```go\nfunc main() {}\n```\n\n```python\nprint()\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	if !strings.Contains(html, `data-sync="a&amp;b"`) {
		t.Error("sync key should be HTML-escaped")
	}
}

func TestCodeGroup_SyncMultipleGroups(t *testing.T) {
	engine := newTestEngine()
	md := ":::code-group sync=\"lang\"\n\n```go\nfmt.Println()\n```\n\n```python\nprint()\n```\n\n:::\n\nSome text.\n\n:::code-group sync=\"lang\"\n\n```go\nimport \"fmt\"\n```\n\n```python\nimport os\n```\n\n:::\n"
	html := renderMarkdown(t, engine, md, true)

	count := strings.Count(html, `data-sync="lang"`)
	if count != 2 {
		t.Errorf("expected 2 groups with data-sync, got %d", count)
	}
}

type failingHighlighter struct {
	mockHighlighter
}

func (f *failingHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
	return nil, errors.New("tokenize failed")
}

func TestCodeGroup_RenderErrorPropagates(t *testing.T) {
	hl := &failingHighlighter{
		mockHighlighter: mockHighlighter{
			themeInfo: kazari.ThemeInfo{FG: "#24292f", BG: "#ffffff"},
		},
	}
	engine := kazari.New(
		kazari.WithHighlighter(hl),
		kazari.WithThemes("light-theme", ""),
		kazari.WithMinify(false),
	)
	exts := []goldmark.Extender{New(engine), CodeGroups(engine)}
	parser := goldmark.New(goldmark.WithExtensions(exts...))

	md := ":::code-group\n\n```go\nfunc main() {}\n```\n\n:::\n"
	var buf bytes.Buffer
	if err := parser.Convert([]byte(md), &buf); err == nil {
		t.Fatal("expected goldmark.Convert to return the render error from inside the code group")
	}
}
