package unwrap

// Phase 0 scope boundary: TokenizeFragment tokenizes one self contained
// fragment. Multi region discovery, byte offset correctness across whole
// files, and splicing belong to the processor phase, so Offset gets only a
// consistency check here, not a positional contract.

import (
	"testing"

	"golang.org/x/net/html"
)

func TestTokenizeFragment_RawBytesRoundTrip(t *testing.T) {
	src := []byte(`<pre tabindex=0 class=chroma><code class="language-go">a &amp; b</code></pre>`)
	tokens, err := TokenizeFragment(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) == 0 {
		t.Fatal("no tokens")
	}
	total := 0
	for i, tok := range tokens {
		if tok.Offset != total {
			t.Fatalf("token %d offset = %d, want %d", i, tok.Offset, total)
		}
		total += len(tok.Raw)
	}
	if total != len(src) {
		t.Fatalf("raw bytes cover %d of %d input bytes", total, len(src))
	}
}

func TestTokenizeFragment_TextIsUnescaped(t *testing.T) {
	tokens, err := TokenizeFragment([]byte(`<code>1 &lt; 2 &amp;&amp; 3 &gt; 2</code>`))
	if err != nil {
		t.Fatal(err)
	}
	var text string
	for _, tok := range tokens {
		if tok.Type == html.TextToken {
			text += tok.Data
		}
	}
	if text != "1 < 2 && 3 > 2" {
		t.Fatalf("got %q", text)
	}
}
