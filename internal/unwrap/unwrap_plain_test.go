package unwrap

import "testing"

func TestPlainFallback_Fixtures(t *testing.T) {
	dirs := []string{
		"hugo-chroma-codefences-false",
		"mdbook-plain",
		"generic-plain-entities",
	}
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) { assertRegion(t, plainFallback{}, dir) })
	}
}

func TestPlainFallback_UnlabeledStillMatches(t *testing.T) {
	// A block with no detectable language still matches with an empty lang.
	// Whether unlabeled blocks get processed is the processor's
	// skipUnlabeled policy, not a detection decision.
	tokens, want := loadFixture(t, "edge-pre-no-language")
	region, ok := plainFallback{}.Match(tokens)
	if !ok {
		t.Fatal("unlabeled pre and code must still match")
	}
	if region.Lang != "" {
		t.Fatalf("lang = %q, want empty", region.Lang)
	}
	diffStrings(t, "code", region.Code, want.Code)
}

func TestPlainFallback_DeclinesPreWithoutCode(t *testing.T) {
	tokens, err := TokenizeFragment([]byte("<pre>ascii art\nonly</pre>"))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := (plainFallback{}).Match(tokens); ok {
		t.Fatal("a pre with no code element must fall through unclaimed")
	}
}

func TestEmptyCodeBlock_MatchesWithEmptyCode(t *testing.T) {
	tokens, want := loadFixture(t, "edge-empty-code-block")
	region, ok := plainFallback{}.Match(tokens)
	if !ok {
		t.Fatal("empty code block must still match")
	}
	if !want.Match {
		t.Fatal("fixture expectation must set match true")
	}
	diffStrings(t, "lang", region.Lang, want.Lang)
	diffStrings(t, "code", region.Code, want.Code)
}

func TestNbspIndentation_PreservedNotNormalized(t *testing.T) {
	// A non breaking space in recovered code stays a non breaking space.
	// Normalizing it could equally corrupt an intentional one inside a
	// string literal, so preservation is the documented behavior.
	tokens, want := loadFixture(t, "edge-nbsp-indentation")
	region, ok := plainFallback{}.Match(tokens)
	if !ok {
		t.Fatal("expected match")
	}
	diffStrings(t, "code", region.Code, want.Code)
}

func TestCRLF_PreservedInRecoveredCode(t *testing.T) {
	// The x/net/html tokenizer applies the WHATWG input preprocessing step,
	// which normalizes CRLF to LF in text before it ever reaches recovery.
	// Verified empirically against a fixture with literal CRLF bytes. This
	// is the friendly outcome, since the engine receives the same bytes a
	// markdown pipeline would produce, and this test pins it so a tokenizer
	// behavior change surfaces as a visible diff.
	tokens, want := loadFixture(t, "edge-crlf")
	region, ok := plainFallback{}.Match(tokens)
	if !ok {
		t.Fatal("expected match")
	}
	diffStrings(t, "code", region.Code, want.Code)
}
