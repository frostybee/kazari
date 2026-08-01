package unwrap

import "testing"

func TestPrismBarePre_Fixtures(t *testing.T) {
	assertRegion(t, prismBarePre{}, "eleventy-prism")
}

func TestPrismBarePre_RequiresClassOnBothElements(t *testing.T) {
	// mdBook's playground pre carries no language class, so the block
	// belongs to the generic fallback, not to Prism.
	tokens, _ := loadFixture(t, "mdbook-plain")
	if _, ok := (prismBarePre{}).Match(tokens); ok {
		t.Fatal("pre without a language class must not match the Prism shape")
	}
}
