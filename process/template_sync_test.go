package process

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHugoHookTemplateInSync guards against drift between the canonical
// user facing render hook template and every copy of it in the repository.
// The e2e fixture must test exactly what users copy, and the example site
// under examples/hugo is the file users are most likely to copy from, so
// both are compared against the same canonical source.
func TestHugoHookTemplateInSync(t *testing.T) {
	canonicalPath := filepath.Join("..", "integrations", "hugo", "render-codeblock.html")
	canonical, err := os.ReadFile(canonicalPath)
	if err != nil {
		t.Fatalf("read canonical template: %v", err)
	}

	copies := []struct {
		name string
		path string
	}{
		{
			name: "e2e fixture",
			path: filepath.Join("testdata", "hugo-fixture", "hook-site", "layouts", "_markup", "render-codeblock.html"),
		},
		{
			name: "examples/hugo",
			path: filepath.Join("..", "examples", "hugo", "layouts", "_markup", "render-codeblock.html"),
		},
	}

	for _, c := range copies {
		t.Run(c.name, func(t *testing.T) {
			got, err := os.ReadFile(c.path)
			if err != nil {
				t.Fatalf("read %s template: %v", c.name, err)
			}
			diffParity(t, "render-codeblock.html canonical vs "+c.name, string(got), string(canonical))
		})
	}
}
