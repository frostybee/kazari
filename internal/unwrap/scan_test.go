package unwrap

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func mustDiscover(t *testing.T, src string) (Page, []byte) {
	t.Helper()
	b := []byte(src)
	tokens, err := TokenizeFragment(b)
	if err != nil {
		t.Fatal(err)
	}
	return Discover(b, tokens), b
}

func candidateBytes(src []byte, c Candidate) string {
	return string(src[c.ByteStart:c.ByteEnd])
}

func TestDiscover_TwoBlocksPerFile(t *testing.T) {
	page, src := mustDiscover(t, `<html><head><title>x</title></head><body>
<p>intro</p>
<pre><code class="language-go">a</code></pre>
<p>middle</p>
<pre><code class="language-js">b</code></pre>
</body></html>`)
	if len(page.Candidates) != 2 {
		t.Fatalf("candidates = %d, want 2", len(page.Candidates))
	}
	for i, c := range page.Candidates {
		if c.Kind != CandidateBarePre || c.Malformed {
			t.Fatalf("candidate %d: kind %v malformed %v", i, c.Kind, c.Malformed)
		}
		got := candidateBytes(src, c)
		if !strings.HasPrefix(got, "<pre>") || !strings.HasSuffix(got, "</pre>") {
			t.Fatalf("candidate %d bytes %q", i, got)
		}
	}
}

func TestDiscover_WrapperConsumesInnerPre(t *testing.T) {
	page, src := mustDiscover(t, `<body><div class="highlight"><pre class="chroma"><code class="language-go" data-lang="go">x</code></pre></div></body>`)
	if len(page.Candidates) != 1 {
		t.Fatalf("candidates = %d, want 1 (outer wins)", len(page.Candidates))
	}
	c := page.Candidates[0]
	if c.Kind != CandidateWrapper {
		t.Fatal("expected wrapper kind")
	}
	if got := candidateBytes(src, c); !strings.HasPrefix(got, `<div class="highlight">`) || !strings.HasSuffix(got, "</div>") {
		t.Fatalf("bytes %q", got)
	}
	// Root convention: the chain must accept the candidate slice as is.
	if _, name, ok := RunChain(DivWrapperChain, c.Tokens); !ok || name != "chroma-div" {
		t.Fatalf("chain on candidate tokens: ok=%v name=%q", ok, name)
	}
}

func TestDiscover_KazariOutputSuppressedAndCounted(t *testing.T) {
	page, _ := mustDiscover(t, `<body>
<div class="kazari-block not-content"><figure class="frame"><pre data-language="go"><code>old</code></pre></figure></div>
<pre><code class="language-go">fresh</code></pre>
</body>`)
	if page.Suppressed != 1 {
		t.Fatalf("suppressed = %d, want 1", page.Suppressed)
	}
	if len(page.Candidates) != 1 {
		t.Fatalf("candidates = %d, want only the sibling after the scope closed", len(page.Candidates))
	}
	if page.SuppressedUnclosed {
		t.Fatal("scope closed properly, must not report unclosed")
	}
}

func TestDiscover_UnclosedReservedScopeSuppressesToEOF(t *testing.T) {
	// A truly truncated file: the reserved scope never closes and no
	// ancestor closes either, so suppression runs to EOF. When an ancestor
	// does close, the scope ends with it, mirroring browser recovery; that
	// path is covered by the scope-ending test above.
	page, _ := mustDiscover(t, `<div class="kazari-block not-content"><pre data-language="go"><code>old</code></pre>
<pre><code class="language-go">never reached</code></pre>`)
	if len(page.Candidates) != 0 {
		t.Fatalf("candidates = %d, want 0 (fail closed)", len(page.Candidates))
	}
	if page.Suppressed != 2 {
		t.Fatalf("suppressed = %d, want 2", page.Suppressed)
	}
	if !page.SuppressedUnclosed {
		t.Fatal("must report the unclosed scope for the warning path")
	}
}

func TestDiscover_UnclosedWrapperDoesNotBlindRestOfFile(t *testing.T) {
	page, src := mustDiscover(t, `<html><head><title>x</title></head><body>
<div class="highlight"><pre><code class="language-go">inside</code></pre>
<p>later</p>
<pre><code class="language-js">after</code></pre>
</body></html>`)
	var malformed, barePres int
	for _, c := range page.Candidates {
		if c.Malformed {
			malformed++
			continue
		}
		if c.Kind == CandidateBarePre {
			barePres++
		}
	}
	if malformed != 1 {
		t.Fatalf("malformed = %d, want 1", malformed)
	}
	// Discovery resumes right after the unclosed root, so both the pre
	// inside the broken wrapper and the later block are still found.
	if barePres != 2 {
		t.Fatalf("bare pre candidates = %d, want 2", barePres)
	}
	if page.ScriptInsert != strings.Index(string(src), "</body>") {
		t.Fatal("body close offset lost after the malformed wrapper")
	}
}

func TestDiscover_UnmatchedWrapperReoffersInnerBlock(t *testing.T) {
	page, src := mustDiscover(t, `<body><div class="highlight-box"><p>note</p><pre><code class="language-go">real code</code></pre></div></body>`)
	if len(page.Candidates) != 1 || page.Candidates[0].Kind != CandidateWrapper {
		t.Fatalf("want one wrapper candidate, got %+v", page.Candidates)
	}
	outer := page.Candidates[0]
	if _, _, ok := RunChain(DivWrapperChain, outer.Tokens); ok {
		t.Fatal("the callout must not match any unwrapper")
	}
	inner := DiscoverWithin(outer.Tokens[1 : len(outer.Tokens)-1])
	if len(inner.Candidates) != 1 || inner.Candidates[0].Kind != CandidateBarePre {
		t.Fatalf("re-offer found %+v", inner.Candidates)
	}
	if got := candidateBytes(src, inner.Candidates[0]); !strings.Contains(got, "real code") {
		t.Fatalf("inner candidate bytes %q", got)
	}
}

func TestDiscover_InjectionOffsetsFullPage(t *testing.T) {
	src := `<html><head><title>x</title></head><body><pre><code>y</code></pre></body></html>`
	page, _ := mustDiscover(t, src)
	if page.LinkInsert != strings.Index(src, "</head>") {
		t.Fatalf("LinkInsert = %d", page.LinkInsert)
	}
	if page.ScriptInsert != strings.Index(src, "</body>") {
		t.Fatalf("ScriptInsert = %d", page.ScriptInsert)
	}
}

func TestDiscover_InjectionOffsetsHeadless(t *testing.T) {
	src := `<body class=x><pre><code>y</code></pre></body>`
	page, _ := mustDiscover(t, src)
	if page.LinkInsert != strings.Index(src, "><pre>")+1 {
		t.Fatalf("LinkInsert = %d, want just after the body start tag", page.LinkInsert)
	}
	if page.ScriptInsert != strings.Index(src, "</body>") {
		t.Fatalf("ScriptInsert = %d", page.ScriptInsert)
	}
}

func TestDiscover_InjectionOffsetsNoHeadNoBody(t *testing.T) {
	src := `<pre><code>y</code></pre>`
	page, _ := mustDiscover(t, src)
	if page.LinkInsert != 0 {
		t.Fatalf("LinkInsert = %d, want 0", page.LinkInsert)
	}
	if page.ScriptInsert != len(src) {
		t.Fatalf("ScriptInsert = %d, want EOF", page.ScriptInsert)
	}
}

func TestDiscover_BOMNoHead(t *testing.T) {
	src := "\xEF\xBB\xBF<pre><code>y</code></pre>"
	page, _ := mustDiscover(t, src)
	if page.LinkInsert != 3 {
		t.Fatalf("LinkInsert = %d, want 3 (after the BOM)", page.LinkInsert)
	}
}

func TestDiscover_GeneratorMeta(t *testing.T) {
	page, _ := mustDiscover(t, `<head><meta name="generator" content="Hugo 0.148.1"></head>`)
	if page.Generator != "Hugo 0.148.1" {
		t.Fatalf("generator = %q", page.Generator)
	}
}

func TestDiscover_ExistingAssetTagSpans(t *testing.T) {
	src := `<head><link rel="stylesheet" href="../kazari.css?v=abc" data-kazari="assets"></head><body><pre><code>x</code></pre><script src="../kazari.js?v=abc" data-kazari="assets"></script></body>`
	page, b := mustDiscover(t, src)
	if len(page.AssetTags) != 2 {
		t.Fatalf("asset tags = %d, want 2", len(page.AssetTags))
	}
	link, script := page.AssetTags[0], page.AssetTags[1]
	if link.Kind != "link" || link.Ref != "../kazari.css?v=abc" {
		t.Fatalf("link tag %+v", link)
	}
	if got := string(b[link.ByteStart:link.ByteEnd]); !strings.HasPrefix(got, "<link") || !strings.HasSuffix(got, ">") {
		t.Fatalf("link span %q", got)
	}
	if script.Kind != "script" || script.Ref != "../kazari.js?v=abc" {
		t.Fatalf("script tag %+v", script)
	}
	if got := string(b[script.ByteStart:script.ByteEnd]); !strings.HasSuffix(got, "</script>") {
		t.Fatalf("script span must include the close tag, got %q", got)
	}
}

func TestDiscover_CandidateTokensRootConvention(t *testing.T) {
	page, _ := mustDiscover(t, `<body><pre data-kazari="ignore"><code>x</code></pre></body>`)
	if len(page.Candidates) != 1 {
		t.Fatalf("candidates = %d", len(page.Candidates))
	}
	c := page.Candidates[0]
	if c.Tokens[0].Type != html.StartTagToken || c.Tokens[0].Data != "pre" {
		t.Fatal("candidate slice must start at the root start tag")
	}
	if got := SkipReason(c.Tokens); got != SkipKazariIgnore {
		t.Fatalf("SkipReason on candidate tokens = %q", got)
	}
}
