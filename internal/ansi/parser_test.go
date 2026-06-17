package ansi

import (
	"testing"
)

func TestParse_PlainText(t *testing.T) {
	lines := Parse("hello world")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if len(lines[0].Tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(lines[0].Tokens))
	}
	tok := lines[0].Tokens[0]
	if tok.Content != "hello world" {
		t.Errorf("content = %q, want %q", tok.Content, "hello world")
	}
	if tok.LightColor != "" || tok.DarkColor != "" {
		t.Errorf("expected no color, got light=%q dark=%q", tok.LightColor, tok.DarkColor)
	}
}

func TestParse_SingleColor(t *testing.T) {
	lines := Parse("\x1b[31mhello\x1b[0m")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	tok := lines[0].Tokens[0]
	if tok.Content != "hello" {
		t.Errorf("content = %q, want %q", tok.Content, "hello")
	}
	if tok.LightColor != "var(--kz-ansi-red)" {
		t.Errorf("light color = %q, want %q", tok.LightColor, "var(--kz-ansi-red)")
	}
	if tok.DarkColor != "var(--kz-ansi-red)" {
		t.Errorf("dark color = %q, want %q", tok.DarkColor, "var(--kz-ansi-red)")
	}
}

func TestParse_MultiColorLine(t *testing.T) {
	lines := Parse("\x1b[31mred\x1b[32mgreen\x1b[0m")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Content != "red" || tokens[0].LightColor != "var(--kz-ansi-red)" {
		t.Errorf("token 0: content=%q color=%q", tokens[0].Content, tokens[0].LightColor)
	}
	if tokens[1].Content != "green" || tokens[1].LightColor != "var(--kz-ansi-green)" {
		t.Errorf("token 1: content=%q color=%q", tokens[1].Content, tokens[1].LightColor)
	}
}

func TestParse_BoldAndColorCombined(t *testing.T) {
	lines := Parse("\x1b[1;33mwarn\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.Content != "warn" {
		t.Errorf("content = %q, want %q", tok.Content, "warn")
	}
	if tok.LightColor != "var(--kz-ansi-yellow)" {
		t.Errorf("color = %q, want yellow", tok.LightColor)
	}
	if tok.FontStyle&2 == 0 {
		t.Error("expected bold font style")
	}
}

func TestParse_StateAcrossLines(t *testing.T) {
	lines := Parse("\x1b[31mfirst\nsecond\x1b[0m")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Tokens[0].LightColor != "var(--kz-ansi-red)" {
		t.Error("line 1 should be red")
	}
	if lines[1].Tokens[0].LightColor != "var(--kz-ansi-red)" {
		t.Error("line 2 should still be red (state carries over)")
	}
}

func TestParse_Reset(t *testing.T) {
	lines := Parse("\x1b[31mred\x1b[0mplain")
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[1].LightColor != "" {
		t.Errorf("after reset, color should be empty, got %q", tokens[1].LightColor)
	}
}

func TestParse_256ColorForeground(t *testing.T) {
	lines := Parse("\x1b[38;5;196mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightColor != "#ff0000" {
		t.Errorf("256-color 196 = %q, want #ff0000", tok.LightColor)
	}
}

func TestParse_24BitRGB(t *testing.T) {
	lines := Parse("\x1b[38;2;255;128;0mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightColor != "#ff8000" {
		t.Errorf("24-bit color = %q, want #ff8000", tok.LightColor)
	}
}

func TestParse_BackgroundColor(t *testing.T) {
	lines := Parse("\x1b[42mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightBG != "var(--kz-ansi-green)" {
		t.Errorf("background = %q, want green", tok.LightBG)
	}
	if tok.DarkBG != "var(--kz-ansi-green)" {
		t.Errorf("dark bg = %q, want green", tok.DarkBG)
	}
}

func TestParse_FontStyles(t *testing.T) {
	tests := []struct {
		name  string
		code  string
		style int
	}{
		{"italic", "\x1b[3mtext\x1b[0m", 1},
		{"bold", "\x1b[1mtext\x1b[0m", 2},
		{"underline", "\x1b[4mtext\x1b[0m", 4},
		{"strikethrough", "\x1b[9mtext\x1b[0m", 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := Parse(tt.code)
			tok := lines[0].Tokens[0]
			if tok.FontStyle != tt.style {
				t.Errorf("font style = %d, want %d", tok.FontStyle, tt.style)
			}
		})
	}
}

func TestParse_UnsetFontStyle(t *testing.T) {
	lines := Parse("\x1b[1mbold\x1b[22mnormal\x1b[0m")
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].FontStyle&2 == 0 {
		t.Error("first token should be bold")
	}
	if tokens[1].FontStyle&2 != 0 {
		t.Error("second token should not be bold after SGR 22")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	lines := Parse("")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if len(lines[0].Tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(lines[0].Tokens))
	}
	if lines[0].Tokens[0].Content != "" {
		t.Errorf("content = %q, want empty", lines[0].Tokens[0].Content)
	}
}

func TestParse_MalformedSequence(t *testing.T) {
	lines := Parse("\x1b[Xmhello")
	tokens := lines[0].Tokens
	combined := ""
	for _, tok := range tokens {
		combined += tok.Content
	}
	if combined != "\x1b[Xmhello" {
		t.Errorf("malformed sequence should pass through, got %q", combined)
	}
}

func TestParse_BrightColors(t *testing.T) {
	lines := Parse("\x1b[90mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightColor != "var(--kz-ansi-bright-black)" {
		t.Errorf("bright black = %q, want var(--kz-ansi-bright-black)", tok.LightColor)
	}
}

func TestParse_EmptyParamsIsReset(t *testing.T) {
	lines := Parse("\x1b[31mred\x1b[mplain")
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[1].LightColor != "" {
		t.Errorf("empty params should reset, color = %q", tokens[1].LightColor)
	}
}

func TestParse_NoEmptyTokens(t *testing.T) {
	lines := Parse("\x1b[31m\x1b[1mtext\x1b[0m")
	tokens := lines[0].Tokens
	for i, tok := range tokens {
		if tok.Content == "" {
			t.Errorf("token %d has empty content", i)
		}
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token (consecutive escapes merged), got %d", len(tokens))
	}
	if tokens[0].LightColor != "var(--kz-ansi-red)" {
		t.Errorf("color = %q, want red", tokens[0].LightColor)
	}
	if tokens[0].FontStyle&2 == 0 {
		t.Error("expected bold")
	}
}

func TestParse_256ColorBackground(t *testing.T) {
	lines := Parse("\x1b[48;5;21mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightBG != "#0000ff" {
		t.Errorf("256-color bg 21 = %q, want #0000ff", tok.LightBG)
	}
}

func TestParse_24BitBackground(t *testing.T) {
	lines := Parse("\x1b[48;2;100;200;50mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightBG != "#64c832" {
		t.Errorf("24-bit bg = %q, want #64c832", tok.LightBG)
	}
}

func TestParse_DefaultForeground(t *testing.T) {
	lines := Parse("\x1b[31mred\x1b[39mdefault\x1b[0m")
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[1].LightColor != "" {
		t.Errorf("SGR 39 should reset fg, got %q", tokens[1].LightColor)
	}
}

func TestParse_DefaultBackground(t *testing.T) {
	lines := Parse("\x1b[42mgreen bg\x1b[49mdefault bg\x1b[0m")
	tokens := lines[0].Tokens
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[1].LightBG != "" {
		t.Errorf("SGR 49 should reset bg, got %q", tokens[1].LightBG)
	}
}

func TestParse_Color256_StandardRange(t *testing.T) {
	lines := Parse("\x1b[38;5;9mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightColor != "#ef2929" {
		t.Errorf("256-color index 9 = %q, want bright red #ef2929", tok.LightColor)
	}
}

func TestParse_Color256_Grayscale(t *testing.T) {
	lines := Parse("\x1b[38;5;232mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightColor != "#080808" {
		t.Errorf("256-color grayscale 232 = %q, want #080808", tok.LightColor)
	}
}

func TestParse_BrightBackground(t *testing.T) {
	lines := Parse("\x1b[101mtext\x1b[0m")
	tok := lines[0].Tokens[0]
	if tok.LightBG != "var(--kz-ansi-bright-red)" {
		t.Errorf("bright red bg = %q, want var(--kz-ansi-bright-red)", tok.LightBG)
	}
	if tok.DarkBG != "var(--kz-ansi-bright-red)" {
		t.Errorf("dark bright red bg = %q, want var(--kz-ansi-bright-red)", tok.DarkBG)
	}
}
