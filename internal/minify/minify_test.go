package minify

import (
	"strings"
	"testing"
)

func TestCSS_Empty(t *testing.T) {
	if got := CSS(""); got != "" {
		t.Errorf("CSS empty input should return empty, got %q", got)
	}
}

func TestCSS_MinifiesWhitespace(t *testing.T) {
	input := "body {\n  color: red;\n  background: blue;\n}"
	got := CSS(input)
	if len(got) >= len(input) {
		t.Errorf("CSS output should be shorter than input: %d >= %d", len(got), len(input))
	}
}

func TestCSS_RemovesComments(t *testing.T) {
	input := "/* comment */ body { color: red; }"
	got := CSS(input)
	if strings.Contains(got, "comment") {
		t.Errorf("CSS comments should be removed, got %q", got)
	}
}

func TestCSS_AlreadyMinified(t *testing.T) {
	input := "body{color:red}"
	got := CSS(input)
	if len(got) > len(input) {
		t.Errorf("already minified CSS should not grow: %d > %d", len(got), len(input))
	}
}

func TestJS_Empty(t *testing.T) {
	if got := JS(""); got != "" {
		t.Errorf("JS empty input should return empty, got %q", got)
	}
}

func TestJS_MinifiesWhitespace(t *testing.T) {
	input := "function test() {\n  return 42;\n}"
	got := JS(input)
	if len(got) >= len(input) {
		t.Errorf("JS output should be shorter than input: %d >= %d", len(got), len(input))
	}
}

func TestJS_RemovesComments(t *testing.T) {
	input := "// comment\nvar x = 1;"
	got := JS(input)
	if strings.Contains(got, "comment") {
		t.Errorf("JS comments should be removed, got %q", got)
	}
}

func TestJS_AlreadyMinified(t *testing.T) {
	input := "var x=1;"
	got := JS(input)
	if len(got) > len(input) {
		t.Errorf("already minified JS should not grow: %d > %d", len(got), len(input))
	}
}
