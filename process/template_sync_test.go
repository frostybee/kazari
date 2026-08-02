package process

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHugoHookTemplateInSync guards against drift between the canonical
// user facing render hook template and the copy the e2e fixture actually
// builds with. The fixture must always test exactly what users copy.
func TestHugoHookTemplateInSync(t *testing.T) {
	canonicalPath := filepath.Join("..", "integrations", "hugo", "render-codeblock.html")
	fixturePath := filepath.Join("testdata", "hugo-fixture", "hook-site", "layouts", "_markup", "render-codeblock.html")

	canonical, err := os.ReadFile(canonicalPath)
	if err != nil {
		t.Fatalf("read canonical template: %v", err)
	}
	fixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture template: %v", err)
	}
	diffParity(t, "render-codeblock.html canonical vs fixture", string(fixture), string(canonical))
}
