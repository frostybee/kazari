package kazarinuri

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/frostybee/kazari"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func setup(t *testing.T) (*NuriHighlighter, context.Context, func()) {
	t.Helper()
	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		t.Fatal(err)
	}
	adapter := New(ctx, hl)
	return adapter, ctx, func() { hl.Close(ctx) }
}

func TestNew(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()
	if adapter == nil {
		t.Fatal("New() returned nil")
	}
}

func TestTokenize_GoCode(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	lines, err := adapter.Tokenize("package main\n\nfunc main() {}\n", "go", "github-dark")
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

func TestTokenize_ContentPreservation(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	code := "func main() { fmt.Println(42) }"
	lines, err := adapter.Tokenize(code, "go", "github-dark")
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

var hexColorRe = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func TestTokenize_ColorsAreHex(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	lines, err := adapter.Tokenize("func main() {}", "go", "github-dark")
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range lines {
		for j, tok := range line {
			if tok.Color != "" && !hexColorRe.MatchString(tok.Color) {
				t.Errorf("line %d token %d: Color %q is not valid hex", i, j, tok.Color)
			}
		}
	}
}

func TestTokenizeDual_MatchingBoundaries(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	light, dark, err := adapter.TokenizeDual("func main() {}", "go", "github-light", "github-dark")
	if err != nil {
		t.Fatal(err)
	}
	if len(light) != len(dark) {
		t.Fatalf("line count mismatch: light=%d dark=%d", len(light), len(dark))
	}
	for i := range light {
		if len(light[i]) != len(dark[i]) {
			t.Errorf("line %d token count mismatch: light=%d dark=%d", i, len(light[i]), len(dark[i]))
		}
	}
}

func TestTokenizeDual_DifferentColors(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	light, dark, err := adapter.TokenizeDual("func main() {}", "go", "github-light", "github-dark")
	if err != nil {
		t.Fatal(err)
	}

	differs := false
	for i := range light {
		for j := range light[i] {
			if light[i][j].Color != dark[i][j].Color {
				differs = true
			}
		}
	}
	if !differs {
		t.Error("expected light and dark tokens to have at least one different color")
	}
}

func TestGetThemeColors_KnownTheme(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	info, err := adapter.GetThemeColors("github-dark")
	if err != nil {
		t.Fatal(err)
	}
	if info.FG == "" {
		t.Error("expected non-empty FG")
	}
	if info.BG == "" {
		t.Error("expected non-empty BG")
	}
}

func TestGetThemeColors_SelectionAndLineNumber(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	info, err := adapter.GetThemeColors("github-dark")
	if err != nil {
		t.Fatal(err)
	}
	if info.SelectionBG == "" {
		t.Error("expected non-empty SelectionBG from VS Code theme")
	}
	if info.LineNumberFG == "" {
		t.Error("expected non-empty LineNumberFG from VS Code theme")
	}
}

func TestGetThemeColors_UnknownTheme(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	_, err := adapter.GetThemeColors("nonexistent-theme-xyz")
	if err == nil {
		t.Error("expected error for unknown theme")
	}
}

func TestGetLoadedLanguages(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	// Trigger a load so at least one language is registered
	adapter.Tokenize("package main", "go", "github-dark")

	langs := adapter.GetLoadedLanguages()
	if len(langs) == 0 {
		t.Error("expected non-empty language list after tokenizing")
	}
	found := false
	for _, l := range langs {
		if l == "go" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'go' in loaded languages, got %v", langs)
	}
}

func TestIntegration_RenderWithEngine(t *testing.T) {
	adapter, _, cleanup := setup(t)
	defer cleanup()

	engine := kazari.New(
		kazari.WithHighlighter(adapter),
		kazari.WithThemes("github-light", "github-dark"),
	)

	html, err := engine.Render("func main() {}", kazari.Options{Lang: "go", Title: "test.go"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "kazari-block") {
		t.Error("expected kazari-block wrapper")
	}
	if !strings.Contains(html, "--sl:") {
		t.Error("expected --sl CSS variable")
	}
	if !strings.Contains(html, "--sd:") {
		t.Error("expected --sd CSS variable")
	}
}
