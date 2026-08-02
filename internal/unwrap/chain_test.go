package unwrap

import "testing"

func TestBarePreChain_ChromaClaimsBeforePlain(t *testing.T) {
	// The individual fixture tests exercise each unwrapper directly; only a
	// full chain run can catch an ordering bug where the generic fallback
	// claims a shape a specific unwrapper should own.
	tokens, want := loadFixture(t, "hugo-chroma-classes")
	region, name, ok := RunChain(BarePreChain, tokens[1:])
	if !ok {
		t.Fatal("expected a chain match")
	}
	if name != "chroma-pre" {
		t.Fatalf("winning unwrapper = %q, want chroma-pre", name)
	}
	diffStrings(t, "code", region.Code, want.Code)
}

func TestDivWrapperChain_RougeClaimsBeforePygments(t *testing.T) {
	// The realistic collision risk: Rouge and Pygments both wrap in an
	// element literally named highlight.
	tokens, want := loadFixture(t, "jekyll-rouge-kramdown")
	region, name, ok := RunChain(DivWrapperChain, tokens)
	if !ok {
		t.Fatal("expected a chain match")
	}
	if name != "rouge-div" {
		t.Fatalf("winning unwrapper = %q, want rouge-div", name)
	}
	diffStrings(t, "lang", region.Lang, want.Lang)
}

func TestDivWrapperChain_RejectsCalloutWithoutCodeSignature(t *testing.T) {
	tokens, want := loadFixture(t, "edge-div-highlight-no-code-signature")
	if _, name, ok := RunChain(DivWrapperChain, tokens); ok != want.Match {
		t.Fatalf("chain claim = %v via %q, want %v", ok, name, want.Match)
	}
}

func TestDataKzMeta_ShortCircuitsMetaSynthesis(t *testing.T) {
	// The attribute value replaces the synthesized meta verbatim while
	// language detection and source recovery still run normally.
	tokens, want := loadFixture(t, "edge-data-kz-meta")
	region, _, ok := RunChain(BarePreChain, tokens)
	if !ok {
		t.Fatal("expected a chain match")
	}
	diffStrings(t, "meta", region.Meta, want.Meta)
	diffStrings(t, "lang", region.Lang, want.Lang)
	diffStrings(t, "code", region.Code, want.Code)
}

func TestDataKzMeta_MarkersShortCircuit(t *testing.T) {
	tokens, want := loadFixture(t, "edge-data-kz-meta-markers")
	region, _, ok := RunChain(BarePreChain, tokens)
	if !ok {
		t.Fatal("expected a chain match")
	}
	diffStrings(t, "meta", region.Meta, want.Meta)
	diffStrings(t, "lang", region.Lang, want.Lang)
	diffStrings(t, "code", region.Code, want.Code)
}

func TestDataKzMeta_CollapseShortCircuit(t *testing.T) {
	tokens, want := loadFixture(t, "edge-data-kz-meta-collapse")
	region, _, ok := RunChain(BarePreChain, tokens)
	if !ok {
		t.Fatal("expected a chain match")
	}
	diffStrings(t, "meta", region.Meta, want.Meta)
	diffStrings(t, "lang", region.Lang, want.Lang)
	diffStrings(t, "code", region.Code, want.Code)
}

func TestHTMLCommentInSpanSoup_Dropped(t *testing.T) {
	tokens, want := loadFixture(t, "edge-html-comment-in-span-soup")
	region, _, ok := RunChain(DivWrapperChain, tokens)
	if !ok {
		t.Fatal("expected a chain match")
	}
	diffStrings(t, "code", region.Code, want.Code)
}
