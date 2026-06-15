package kazarichroma

import (
	"regexp"
	"strings"
	"testing"

	"github.com/frostybee/kazari"
)

func TestNew_Default(t *testing.T) {
	h := New()
	if h == nil {
		t.Fatal("New() returned nil")
	}
	if h.styleMap != nil {
		t.Error("default styleMap should be nil")
	}
}

func TestNew_WithStyleMap(t *testing.T) {
	m := map[string]string{"light": "github", "dark": "monokai"}
	h := New(WithStyleMap(m))
	if h.styleMap == nil {
		t.Fatal("styleMap should not be nil")
	}
	if h.styleMap["light"] != "github" {
		t.Errorf("styleMap[light] = %q, want github", h.styleMap["light"])
	}
}

func TestTokenize_GoCode(t *testing.T) {
	h := New()
	lines, err := h.Tokenize("package main\n\nfunc main() {}\n", "go", "monokai")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < 3 {
		t.Fatalf("expected >= 3 lines, got %d", len(lines))
	}
	hasColor := false
	for _, line := range lines {
		for _, tok := range line {
			if tok.Color != "" {
				hasColor = true
			}
		}
	}
	if !hasColor {
		t.Error("expected at least one token with a color")
	}
}

func TestTokenize_UnknownLang(t *testing.T) {
	h := New()
	_, err := h.Tokenize("hello", "nonexistent-lang-xyz", "monokai")
	if err == nil {
		t.Error("expected error for unknown language")
	}
}

func TestTokenize_EmptyCode(t *testing.T) {
	h := New()
	lines, err := h.Tokenize("", "go", "monokai")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) < 1 {
		t.Error("expected at least 1 line for empty input")
	}
}

func TestTokenize_MultiLine(t *testing.T) {
	h := New()
	code := "line1\nline2\nline3"
	lines, err := h.Tokenize(code, "text", "monokai")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

var hexColorRe = regexp.MustCompile(`^#[0-9a-f]{6}$`)

func TestTokenize_ColorsAreHex(t *testing.T) {
	h := New()
	lines, err := h.Tokenize("func main() {}", "go", "monokai")
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range lines {
		for j, tok := range line {
			if tok.Color != "" && !hexColorRe.MatchString(tok.Color) {
				t.Errorf("line %d token %d: Color %q is not valid hex", i, j, tok.Color)
			}
			if tok.BgColor != "" && !hexColorRe.MatchString(tok.BgColor) {
				t.Errorf("line %d token %d: BgColor %q is not valid hex", i, j, tok.BgColor)
			}
		}
	}
}

func TestTokenize_FontStyles(t *testing.T) {
	h := New()
	// GenericEmph is italic and GenericStrong is bold in monokai.
	// RST (reStructuredText) produces GenericEmph from *emphasis* and GenericStrong from **strong**.
	code := "*emphasis* and **strong**"
	lines, err := h.Tokenize(code, "rst", "monokai")
	if err != nil {
		t.Fatal(err)
	}
	hasItalic := false
	hasBold := false
	for _, line := range lines {
		for _, tok := range line {
			if tok.FontStyle&kazari.FontStyleItalic != 0 {
				hasItalic = true
			}
			if tok.FontStyle&kazari.FontStyleBold != 0 {
				hasBold = true
			}
		}
	}
	if !hasItalic {
		t.Error("expected at least one italic token (GenericEmph) in monokai rst")
	}
	if !hasBold {
		t.Error("expected at least one bold token (GenericStrong) in monokai rst")
	}
}

func TestTokenize_StyleMap(t *testing.T) {
	h := New(WithStyleMap(map[string]string{
		"my-light": "github",
	}))
	lines, err := h.Tokenize("x = 1", "python", "my-light")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestTokenize_ContentReconstruction(t *testing.T) {
	h := New()
	code := "func main() { fmt.Println(42) }"
	lines, err := h.Tokenize(code, "go", "monokai")
	if err != nil {
		t.Fatal(err)
	}
	var reconstructed strings.Builder
	for i, line := range lines {
		for _, tok := range line {
			reconstructed.WriteString(tok.Content)
		}
		if i < len(lines)-1 {
			reconstructed.WriteByte('\n')
		}
	}
	if reconstructed.String() != code {
		t.Errorf("content mismatch:\n  got:  %q\n  want: %q", reconstructed.String(), code)
	}
}

func TestGetThemeColors_LightStyle(t *testing.T) {
	h := New()
	// Use dracula which has both Text.Colour and Background set.
	// Many Pygments-derived styles (github, emacs, etc.) leave Text.Colour unset.
	info, err := h.GetThemeColors("dracula")
	if err != nil {
		t.Fatal(err)
	}
	if info.FG == "" {
		t.Error("expected non-empty FG for dracula style")
	}
	if info.BG == "" {
		t.Error("expected non-empty BG for dracula style")
	}
}

func TestGetThemeColors_DarkStyle(t *testing.T) {
	h := New()
	info, err := h.GetThemeColors("monokai")
	if err != nil {
		t.Fatal(err)
	}
	if info.BG == "" {
		t.Error("expected non-empty BG for monokai")
	}
}

func TestGetThemeColors_LineNumberFG(t *testing.T) {
	h := New()
	info, err := h.GetThemeColors("monokai")
	if err != nil {
		t.Fatal(err)
	}
	if info.LineNumberFG == "" {
		t.Error("expected synthesized LineNumberFG for monokai")
	}
}

func TestGetThemeColors_StyleMap(t *testing.T) {
	h := New(WithStyleMap(map[string]string{
		"my-dark": "monokai",
	}))
	info, err := h.GetThemeColors("my-dark")
	if err != nil {
		t.Fatal(err)
	}
	if info.BG == "" {
		t.Error("expected non-empty BG when resolved via style map")
	}
}

func TestGetLoadedLanguages(t *testing.T) {
	h := New()
	langs := h.GetLoadedLanguages()
	if len(langs) < 100 {
		t.Errorf("expected 100+ languages, got %d", len(langs))
	}
	found := false
	for _, l := range langs {
		if l == "Go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Go in loaded languages")
	}
}

func TestIntegration_RenderWithEngine(t *testing.T) {
	hl := New()
	engine := kazari.New(
		kazari.WithHighlighter(hl),
		kazari.WithThemes("monokai", ""),
	)

	html, err := engine.Render("func main() {}", kazari.Options{Lang: "go", Title: "test.go"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "kazari-block") {
		t.Error("expected kazari-block wrapper in output")
	}
	if !strings.Contains(html, "kz-line") {
		t.Error("expected kz-line divs in output")
	}
	if !strings.Contains(html, "--sl:") {
		t.Error("expected --sl CSS variable in token spans")
	}
}

func TestIntegration_DualTheme(t *testing.T) {
	hl := New(WithStyleMap(map[string]string{
		"light": "github",
		"dark":  "monokai",
	}))
	engine := kazari.New(
		kazari.WithHighlighter(hl),
		kazari.WithThemes("light", "dark"),
	)

	html, err := engine.Render("x = 1", kazari.Options{Lang: "python"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "--sl:") {
		t.Error("expected --sl (light) CSS variable")
	}
	if !strings.Contains(html, "--sd:") {
		t.Error("expected --sd (dark) CSS variable")
	}
}

func TestIntegration_CSS(t *testing.T) {
	hl := New()
	engine := kazari.New(
		kazari.WithHighlighter(hl),
		kazari.WithThemes("monokai", ""),
	)

	css := engine.CSS()
	if !strings.Contains(css, "--kz-") {
		t.Error("expected --kz- CSS variables in CSS output")
	}
}
