package unwrap

import "testing"

func assertSkip(t *testing.T, dir string) {
	t.Helper()
	tokens, want := loadFixture(t, dir)
	if got := SkipReason(tokens); got != want.Skip {
		t.Fatalf("SkipReason = %q, want %q", got, want.Skip)
	}
}

func TestSkip_KzClass(t *testing.T) {
	assertSkip(t, "edge-already-processed-kz")
}

func TestSkip_DataKazariIgnore(t *testing.T) {
	assertSkip(t, "edge-data-kazari-ignore")
}

func TestSkip_Mermaid(t *testing.T) {
	assertSkip(t, "edge-mermaid")
}

func TestSkip_EmptyForOrdinaryBlocks(t *testing.T) {
	for _, dir := range []string{"hugo-chroma-classes", "generic-plain-entities"} {
		tokens, _ := loadFixture(t, dir)
		if got := SkipReason(tokens); got != "" {
			t.Fatalf("%s: SkipReason = %q, want empty", dir, got)
		}
	}
}
