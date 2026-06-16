package link

import (
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func TestExtractLinks_NoLinks(t *testing.T) {
	code := "fmt.Println(\"hello\")"
	cleaned, links := ExtractLinks(code)
	if cleaned != code {
		t.Errorf("got %q, want %q", cleaned, code)
	}
	if links[0] != nil {
		t.Errorf("expected nil links for line 0, got %v", links[0])
	}
}

func TestExtractLinks_SingleLink(t *testing.T) {
	code := `import @[fmt](https://pkg.go.dev/fmt)`
	cleaned, links := ExtractLinks(code)
	if cleaned != `import fmt` {
		t.Errorf("cleaned = %q, want %q", cleaned, "import fmt")
	}
	want := []config.LinkAnnotation{{Start: 7, End: 10, URL: "https://pkg.go.dev/fmt"}}
	if len(links[0]) != 1 {
		t.Fatalf("got %d links, want 1", len(links[0]))
	}
	if links[0][0] != want[0] {
		t.Errorf("got %+v, want %+v", links[0][0], want[0])
	}
}

func TestExtractLinks_MultipleLinks(t *testing.T) {
	code := `@[foo](http://a) bar @[baz](http://b)`
	cleaned, links := ExtractLinks(code)
	if cleaned != "foo bar baz" {
		t.Errorf("cleaned = %q, want %q", cleaned, "foo bar baz")
	}
	if len(links[0]) != 2 {
		t.Fatalf("got %d links, want 2", len(links[0]))
	}
	if links[0][0].Start != 0 || links[0][0].End != 3 {
		t.Errorf("link 0: got start=%d end=%d, want start=0 end=3", links[0][0].Start, links[0][0].End)
	}
	if links[0][1].Start != 8 || links[0][1].End != 11 {
		t.Errorf("link 1: got start=%d end=%d, want start=8 end=11", links[0][1].Start, links[0][1].End)
	}
}

func TestExtractLinks_NoBracketWithoutAt(t *testing.T) {
	code := `[not a link](http://example.com)`
	cleaned, links := ExtractLinks(code)
	if cleaned != code {
		t.Errorf("cleaned = %q, want %q", cleaned, code)
	}
	if links[0] != nil {
		t.Errorf("expected nil links, got %v", links[0])
	}
}

func TestExtractLinks_MultiLine(t *testing.T) {
	code := "line1\n@[link](http://a)\nline3"
	cleaned, links := ExtractLinks(code)
	if cleaned != "line1\nlink\nline3" {
		t.Errorf("cleaned = %q, want %q", cleaned, "line1\nlink\nline3")
	}
	if links[0] != nil {
		t.Errorf("line 0 should have no links")
	}
	if len(links[1]) != 1 {
		t.Fatalf("line 1: got %d links, want 1", len(links[1]))
	}
	if links[1][0].Start != 0 || links[1][0].End != 4 {
		t.Errorf("line 1 link: got start=%d end=%d, want start=0 end=4", links[1][0].Start, links[1][0].End)
	}
	if links[2] != nil {
		t.Errorf("line 2 should have no links")
	}
}

func TestExtractLinks_UnsafeSchemeRejected(t *testing.T) {
	code := `@[click](javascript:alert(1))`
	cleaned, links := ExtractLinks(code)
	if cleaned != code {
		t.Errorf("unsafe link should be kept as literal text, got cleaned=%q", cleaned)
	}
	if links[0] != nil {
		t.Errorf("unsafe link should produce no annotations, got %v", links[0])
	}
}

func TestExtractLinks_DataSchemeRejected(t *testing.T) {
	code := `@[x](data:text/html,<script>alert(1)</script>)`
	cleaned, links := ExtractLinks(code)
	if cleaned != code {
		t.Errorf("data: scheme should be kept as literal text, got cleaned=%q", cleaned)
	}
	if links[0] != nil {
		t.Errorf("data: scheme should produce no annotations, got %v", links[0])
	}
}

func TestExtractLinks_ProtocolRelativeRejected(t *testing.T) {
	code := `@[click](//evil.com/payload)`
	cleaned, links := ExtractLinks(code)
	if cleaned != code {
		t.Errorf("protocol-relative URL should be kept as literal text, got cleaned=%q", cleaned)
	}
	if links[0] != nil {
		t.Errorf("protocol-relative URL should produce no annotations, got %v", links[0])
	}
}

func TestExtractLinks_RelativePathAllowed(t *testing.T) {
	code := `@[docs](/api/reference)`
	cleaned, links := ExtractLinks(code)
	if cleaned != "docs" {
		t.Errorf("cleaned = %q, want %q", cleaned, "docs")
	}
	if len(links[0]) != 1 || links[0][0].URL != "/api/reference" {
		t.Errorf("relative path should be allowed, got %v", links[0])
	}
}

func TestExtractLinks_URLWithSpecialChars(t *testing.T) {
	code := `@[docs](https://example.com/path?q=1&r=2#section)`
	cleaned, links := ExtractLinks(code)
	if cleaned != "docs" {
		t.Errorf("cleaned = %q, want %q", cleaned, "docs")
	}
	if links[0][0].URL != "https://example.com/path?q=1&r=2#section" {
		t.Errorf("URL = %q", links[0][0].URL)
	}
}
