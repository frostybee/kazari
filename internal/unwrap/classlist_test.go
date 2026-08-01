package unwrap

import (
	"testing"

	"golang.org/x/net/html"
)

func attrs(class string) []html.Attribute {
	return []html.Attribute{{Key: "class", Val: class}}
}

func TestHasClass(t *testing.T) {
	cases := []struct {
		class string
		name  string
		want  bool
	}{
		{"language-css", "language-c", false},
		{"language-c", "language-c", true},
		{"highlight not-prose", "highlight", true},
		{"highlighter-rouge", "rouge", false},
		{"highlighter-rouge", "highlighter-rouge", true},
		{"", "highlight", false},
		{"line hl", "hl", true},
	}
	for _, c := range cases {
		if got := hasClass(attrs(c.class), c.name); got != c.want {
			t.Errorf("hasClass(%q, %q) = %v, want %v", c.class, c.name, got, c.want)
		}
	}
	if hasClass(nil, "highlight") {
		t.Error("hasClass on missing class attribute must be false")
	}
}

func TestClassWithPrefix(t *testing.T) {
	if tok, ok := classWithPrefix(attrs("language-rust edition2024"), "language-"); !ok || tok != "language-rust" {
		t.Fatalf("got %q, %v", tok, ok)
	}
	// A bare token equal to the prefix minus the dash must not count.
	if _, ok := classWithPrefix(attrs("highlight"), "highlight-"); ok {
		t.Fatal("highlight must not satisfy the highlight- prefix")
	}
	if _, ok := classWithPrefix(attrs("highlight-python3 notranslate"), "highlight-"); !ok {
		t.Fatal("highlight-python3 must satisfy the highlight- prefix")
	}
}

func TestAttrValue(t *testing.T) {
	a := []html.Attribute{{Key: "data-lang", Val: "go"}, {Key: "class", Val: "chroma"}}
	if v, ok := attrValue(a, "data-lang"); !ok || v != "go" {
		t.Fatalf("got %q, %v", v, ok)
	}
	if _, ok := attrValue(a, "style"); ok {
		t.Fatal("missing attribute must report false")
	}
}
