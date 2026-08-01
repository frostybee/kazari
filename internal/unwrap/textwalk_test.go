package unwrap

import (
	"strings"
	"testing"
)

// walkFragment tokenizes a fragment and recovers text over the whole window.
func walkFragment(t *testing.T, fragment string) string {
	t.Helper()
	tokens, err := TokenizeFragment([]byte(fragment))
	if err != nil {
		t.Fatal(err)
	}
	return recoverText(tokens)
}

func TestIsGutterClass(t *testing.T) {
	for _, c := range []string{"ln", "lnt", "lntd", "gutter", "gl", "linenos", "lineno", "linenodiv", "line-numbers", "line-numbers-rows"} {
		if !isGutterElement(attrs(c)) {
			t.Errorf("%q must be a gutter class", c)
		}
	}
	for _, c := range []string{"line", "cl", "code", "token", "giallo-l"} {
		if isGutterElement(attrs(c)) {
			t.Errorf("%q must not be a gutter class", c)
		}
	}
}

func TestRecoverText_GutterExcluded(t *testing.T) {
	got := walkFragment(t, `<span class="ln">1</span>hello
<span class="ln">2</span>world`)
	if got != "hello\nworld" {
		t.Fatalf("got %q", got)
	}
}

func TestRecoverText_GutterSubtreeFullyExcluded(t *testing.T) {
	// Nested markup inside a gutter element stays excluded until the
	// gutter closes, then recovery resumes.
	got := walkFragment(t, `<td class="gutter"><pre>1
2</pre></td><td>code here</td>`)
	if got != "code here" {
		t.Fatalf("got %q", got)
	}
}

func TestRecoverText_BrBecomesNewline(t *testing.T) {
	if got := walkFragment(t, `a<br>b<br/>c`); got != "a\nb\nc" {
		t.Fatalf("got %q", got)
	}
}

func TestRecoverText_CommentsIgnored(t *testing.T) {
	if got := walkFragment(t, `a<!-- nope -->b`); got != "ab" {
		t.Fatalf("got %q", got)
	}
}

func TestRecoverText_EntitiesDecodedOnce(t *testing.T) {
	// The tokenizer already unescaped the text; a double decode would turn
	// the recovered ampersand entity back into a bare ampersand twice and
	// corrupt code that legitimately contains entity syntax.
	got := walkFragment(t, `&amp;amp; is the escaped form`)
	if got != "&amp; is the escaped form" {
		t.Fatalf("got %q", got)
	}
	if strings.Contains(got, "& is") {
		t.Fatal("text was unescaped twice")
	}
}
