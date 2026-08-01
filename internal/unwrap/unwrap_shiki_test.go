package unwrap

import "testing"

func TestShikiAstroBarePre_Fixtures(t *testing.T) {
	assertRegion(t, shikiAstroBarePre{}, "astro-shiki")
}

func TestShikiDeclinesExpressiveCode(t *testing.T) {
	// Expressive Code carries data-language but not the astro-code class,
	// and its div.ec-line structure holds no newline text nodes. Claiming
	// it would glue every line together, so the whole chain must leave it
	// alone in this phase.
	tokens, want := loadFixture(t, "astro-expressive-code")
	if _, ok := (shikiAstroBarePre{}).Match(tokens); ok {
		t.Fatal("shiki unwrapper must decline the Expressive Code shape")
	}
	if _, _, ok := RunChain(BarePreChain, tokens); ok != want.Match {
		t.Fatalf("chain claim = %v, want %v", ok, want.Match)
	}
}
