package unwrap

import "testing"

func TestChromaDivWrapper_Fixtures(t *testing.T) {
	dirs := []string{
		"hugo-chroma-classes",
		"hugo-chroma-inline-styles",
		"hugo-chroma-linenos-inline",
		"hugo-chroma-lntable",
		"hugo-chroma-hl-lines",
		"hugo-chroma-linenostart",
		"edge-html-comment-in-span-soup",
	}
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) { assertRegion(t, chromaDivWrapper{}, dir) })
	}
}

func TestChromaBarePre_Fixtures(t *testing.T) {
	// Bare Chroma output is the classes fixture without Hugo's wrapper div,
	// so slicing off the root token reproduces it exactly. The trailing div
	// close token stays in the window and must be harmless.
	tokens, want := loadFixture(t, "hugo-chroma-classes")
	region, ok := chromaBarePre{}.Match(tokens[1:])
	if !ok {
		t.Fatal("expected bare pre.chroma to match")
	}
	diffStrings(t, "lang", region.Lang, want.Lang)
	diffStrings(t, "meta", region.Meta, want.Meta)
	diffStrings(t, "code", region.Code, want.Code)
}

func TestChromaBarePre_DeclinesPlainPre(t *testing.T) {
	tokens, _ := loadFixture(t, "hugo-chroma-codefences-false")
	if _, ok := (chromaBarePre{}).Match(tokens); ok {
		t.Fatal("plain pre without the chroma class must fall through to the generic fallback")
	}
}
