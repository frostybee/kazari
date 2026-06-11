package kazari

import (
	"testing"

	"github.com/frostybee/kazari/internal/render"
)

func TestMergeTokens_FastPath_MatchingBoundaries(t *testing.T) {
	light := [][]Token{
		{{Content: "func", Color: "#cf222e"}, {Content: " main", Color: "#6639ba"}},
	}
	dark := [][]Token{
		{{Content: "func", Color: "#ff7b72"}, {Content: " main", Color: "#d2a8ff"}},
	}

	lines := mergeTokens(light, dark)
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	if len(lines[0].Tokens) != 2 {
		t.Fatalf("got %d tokens, want 2", len(lines[0].Tokens))
	}

	tok := lines[0].Tokens[0]
	if tok.Content != "func" {
		t.Errorf("Content = %q, want %q", tok.Content, "func")
	}
	if tok.LightColor != "#cf222e" {
		t.Errorf("LightColor = %q, want %q", tok.LightColor, "#cf222e")
	}
	if tok.DarkColor != "#ff7b72" {
		t.Errorf("DarkColor = %q, want %q", tok.DarkColor, "#ff7b72")
	}
}

func TestMergeTokens_FastPath_SingleTheme(t *testing.T) {
	light := [][]Token{
		{{Content: "hello", Color: "#111111", BgColor: "#eeeeee"}},
	}

	lines := mergeTokens(light, nil)
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}

	tok := lines[0].Tokens[0]
	if tok.LightColor != "#111111" {
		t.Errorf("LightColor = %q, want %q", tok.LightColor, "#111111")
	}
	if tok.DarkColor != "" {
		t.Errorf("DarkColor = %q, want empty (no dark theme)", tok.DarkColor)
	}
	if tok.LightBG != "#eeeeee" {
		t.Errorf("LightBG = %q, want %q", tok.LightBG, "#eeeeee")
	}
	if tok.DarkBG != "" {
		t.Errorf("DarkBG = %q, want empty", tok.DarkBG)
	}
}

func TestMergeTokens_FastPath_EmptyInput(t *testing.T) {
	lines := mergeTokens([][]Token{}, [][]Token{})
	if len(lines) != 0 {
		t.Errorf("got %d lines, want 0", len(lines))
	}
}

func TestMergeTokens_FastPath_EmptyLine(t *testing.T) {
	light := [][]Token{{}}
	dark := [][]Token{{}}

	lines := mergeTokens(light, dark)
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	if len(lines[0].Tokens) != 0 {
		t.Errorf("got %d tokens, want 0", len(lines[0].Tokens))
	}
}

func TestMergeTokens_FastPath_FontStylePreserved(t *testing.T) {
	light := [][]Token{
		{{Content: "italic", Color: "#aaa", FontStyle: FontStyleItalic | FontStyleBold}},
	}
	dark := [][]Token{
		{{Content: "italic", Color: "#bbb"}},
	}

	lines := mergeTokens(light, dark)
	tok := lines[0].Tokens[0]
	if tok.FontStyle != FontStyleItalic|FontStyleBold {
		t.Errorf("FontStyle = %d, want %d", tok.FontStyle, FontStyleItalic|FontStyleBold)
	}
}

func TestMergeTokens_FastPath_MultipleLines(t *testing.T) {
	light := [][]Token{
		{{Content: "line1", Color: "#aaa"}},
		{{Content: "line2", Color: "#bbb"}},
		{{Content: "line3", Color: "#ccc"}},
	}
	dark := [][]Token{
		{{Content: "line1", Color: "#111"}},
		{{Content: "line2", Color: "#222"}},
		{{Content: "line3", Color: "#333"}},
	}

	lines := mergeTokens(light, dark)
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	if lines[2].Tokens[0].DarkColor != "#333" {
		t.Errorf("line 3 DarkColor = %q, want %q", lines[2].Tokens[0].DarkColor, "#333")
	}
}

func TestMergeTokens_FastPath_DarkShorterThanLight(t *testing.T) {
	light := [][]Token{
		{{Content: "line1", Color: "#aaa"}},
		{{Content: "line2", Color: "#bbb"}},
	}
	dark := [][]Token{
		{{Content: "line1", Color: "#111"}},
	}

	lines := mergeTokens(light, dark)
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}
	if lines[0].Tokens[0].DarkColor != "#111" {
		t.Errorf("line 1 DarkColor = %q, want %q", lines[0].Tokens[0].DarkColor, "#111")
	}
	if lines[1].Tokens[0].DarkColor != "" {
		t.Errorf("line 2 DarkColor = %q, want empty (dark has no line 2)", lines[1].Tokens[0].DarkColor)
	}
}

func TestMergeTokens_SlowPath_DifferentBoundaries(t *testing.T) {
	// Light: "he" + "llo"  (2 tokens, different boundary than dark)
	// Dark:  "hel" + "lo"  (2 tokens)
	light := [][]Token{
		{{Content: "he", Color: "#aaa"}, {Content: "llo", Color: "#bbb"}},
	}
	dark := [][]Token{
		{{Content: "hel", Color: "#111"}, {Content: "lo", Color: "#222"}},
	}

	lines := mergeTokens(light, dark)
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	tokens := lines[0].Tokens

	// Verify content reconstruction
	var content string
	for _, tok := range tokens {
		content += tok.Content
	}
	if content != "hello" {
		t.Errorf("reconstructed content = %q, want %q", content, "hello")
	}

	// Every token should have both light and dark colors
	for i, tok := range tokens {
		if tok.LightColor == "" {
			t.Errorf("token %d (%q): missing LightColor", i, tok.Content)
		}
		if tok.DarkColor == "" {
			t.Errorf("token %d (%q): missing DarkColor", i, tok.Content)
		}
	}
}

func TestMergeTokens_SlowPath_FinerGranularity(t *testing.T) {
	// Light has one token, dark has three
	light := [][]Token{
		{{Content: "abcdef", Color: "#aaa"}},
	}
	dark := [][]Token{
		{{Content: "ab", Color: "#111"}, {Content: "cd", Color: "#222"}, {Content: "ef", Color: "#333"}},
	}

	lines := mergeTokens(light, dark)
	tokens := lines[0].Tokens
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}

	assertMerged(t, tokens[0], "ab", "#aaa", "#111")
	assertMerged(t, tokens[1], "cd", "#aaa", "#222")
	assertMerged(t, tokens[2], "ef", "#aaa", "#333")
}

func TestMergeTokens_SlowPath_LightLongerContent(t *testing.T) {
	// Light has more content than dark (different token counts triggers slow path)
	light := [][]Token{
		{{Content: "hello", Color: "#aaa"}, {Content: " world", Color: "#bbb"}},
	}
	dark := [][]Token{
		{{Content: "hello world", Color: "#111"}},
	}

	lines := mergeTokens(light, dark)
	tokens := lines[0].Tokens

	var content string
	for _, tok := range tokens {
		content += tok.Content
	}
	if content != "hello world" {
		t.Errorf("reconstructed content = %q, want %q", content, "hello world")
	}

	// First token should have both colors
	if tokens[0].LightColor != "#aaa" {
		t.Errorf("token 0 LightColor = %q, want %q", tokens[0].LightColor, "#aaa")
	}
	if tokens[0].DarkColor != "#111" {
		t.Errorf("token 0 DarkColor = %q, want %q", tokens[0].DarkColor, "#111")
	}
}

func TestMergeTokens_SlowPath_BgColors(t *testing.T) {
	light := [][]Token{
		{{Content: "ab", Color: "#aaa", BgColor: "#eee"}, {Content: "cd", Color: "#bbb", BgColor: "#ddd"}},
	}
	dark := [][]Token{
		{{Content: "abcd", Color: "#111", BgColor: "#222"}},
	}

	lines := mergeTokens(light, dark)
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("got %d tokens, want 2", len(tokens))
	}
	if tokens[0].LightBG != "#eee" || tokens[0].DarkBG != "#222" {
		t.Errorf("token 0: LightBG=%q DarkBG=%q", tokens[0].LightBG, tokens[0].DarkBG)
	}
	if tokens[1].LightBG != "#ddd" || tokens[1].DarkBG != "#222" {
		t.Errorf("token 1: LightBG=%q DarkBG=%q", tokens[1].LightBG, tokens[1].DarkBG)
	}
}

func TestPlaintextLines_Basic(t *testing.T) {
	lines := plaintextLines("line1\nline2\nline3")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	for i, want := range []string{"line1", "line2", "line3"} {
		if lines[i].Tokens[0].Content != want {
			t.Errorf("line %d = %q, want %q", i, lines[i].Tokens[0].Content, want)
		}
	}
}

func TestPlaintextLines_Empty(t *testing.T) {
	lines := plaintextLines("")
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	if lines[0].Tokens[0].Content != "" {
		t.Errorf("Content = %q, want empty", lines[0].Tokens[0].Content)
	}
}

func TestPlaintextLines_TrailingNewline(t *testing.T) {
	lines := plaintextLines("a\nb\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2 (trailing newline stripped)", len(lines))
	}
}

func TestSplitLines_NoTrailingNewline(t *testing.T) {
	lines := splitLines("a\nb")
	if len(lines) != 2 {
		t.Fatalf("got %d, want 2", len(lines))
	}
}

func TestSplitLines_TrailingNewline(t *testing.T) {
	lines := splitLines("a\nb\n")
	if len(lines) != 2 {
		t.Fatalf("got %d, want 2 (trailing empty stripped)", len(lines))
	}
}

func TestSplitLines_Empty(t *testing.T) {
	lines := splitLines("")
	if len(lines) != 1 {
		t.Fatalf("got %d, want 1", len(lines))
	}
	if lines[0] != "" {
		t.Errorf("got %q, want empty", lines[0])
	}
}

func TestExpandTabs(t *testing.T) {
	if got := expandTabs("a\tb", 4); got != "a    b" {
		t.Errorf("got %q, want %q", got, "a    b")
	}
}

func TestExpandTabs_NoTabs(t *testing.T) {
	input := "no tabs here"
	if got := expandTabs(input, 4); got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func assertMerged(t *testing.T, tok render.MergedToken, content, lightColor, darkColor string) {
	t.Helper()
	if tok.Content != content {
		t.Errorf("Content = %q, want %q", tok.Content, content)
	}
	if tok.LightColor != lightColor {
		t.Errorf("LightColor = %q, want %q", tok.LightColor, lightColor)
	}
	if tok.DarkColor != darkColor {
		t.Errorf("DarkColor = %q, want %q", tok.DarkColor, darkColor)
	}
}
