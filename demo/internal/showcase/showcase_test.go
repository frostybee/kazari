package showcase

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/frostybee/kazari"
)

type mockHighlighter struct{}

func (mockHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
	lines := strings.Split(code, "\n")
	tokens := make([][]kazari.Token, len(lines))
	for index, line := range lines {
		tokens[index] = []kazari.Token{{Content: line, Color: "#24292f"}}
	}
	return tokens, nil
}

func (mockHighlighter) GetThemeColors(theme string) (kazari.ThemeInfo, error) {
	if strings.Contains(theme, "dark") || theme == "dracula" {
		return kazari.ThemeInfo{FG: "#f0f0f0", BG: "#161b22", SelectionBG: "#264f78", LineNumberFG: "#8b949e", FoldBG: "#30363d"}, nil
	}
	return kazari.ThemeInfo{FG: "#24292f", BG: "#ffffff", SelectionBG: "#add6ff", LineNumberFG: "#57606a", FoldBG: "#ddf4ff"}, nil
}

func (mockHighlighter) GetLoadedLanguages() []string {
	return []string{"go", "javascript", "bash", "powershell", "text", "jsx", "python", "rust", "mermaid", "ansi", "diff"}
}

func TestBuildCatalogParityAndCompleteness(t *testing.T) {
	nuri, err := Build(Config{
		BackendName: "Nuri",
		HTMLFile:    "showcase.html",
		OtherName:   "Chroma",
		OtherHref:   "../chroma/showcase-chroma.html",
	}, mockHighlighter{})
	if err != nil {
		t.Fatalf("Build Nuri: %v", err)
	}
	chroma, err := Build(Config{
		BackendName: "Chroma",
		HTMLFile:    "showcase-chroma.html",
		OtherName:   "Nuri",
		OtherHref:   "../nuri/showcase.html",
	}, mockHighlighter{})
	if err != nil {
		t.Fatalf("Build Chroma: %v", err)
	}

	nuriIDs := exampleIDs(nuri.Page.Categories)
	chromaIDs := exampleIDs(chroma.Page.Categories)
	if !reflect.DeepEqual(nuriIDs, chromaIDs) {
		t.Fatalf("backend catalogs differ:\nNuri: %v\nChroma: %v", nuriIDs, chromaIDs)
	}
	if len(nuri.Page.Categories) != 8 {
		t.Fatalf("category count = %d, want 8", len(nuri.Page.Categories))
	}
	if len(nuriIDs) != 43 {
		t.Fatalf("example count = %d, want 43", len(nuriIDs))
	}

	seen := make(map[string]bool, len(nuriIDs))
	for _, category := range nuri.Page.Categories {
		for _, example := range category.Examples {
			if seen[example.ID] {
				t.Fatalf("duplicate example ID %q", example.ID)
			}
			seen[example.ID] = true
			if strings.TrimSpace(example.Description) == "" {
				t.Fatalf("example %q has no description", example.ID)
			}
			if len(example.Recipes) == 0 {
				t.Fatalf("example %q has no recipes", example.ID)
			}
			for _, recipe := range example.Recipes {
				if recipe.Label == "" || strings.TrimSpace(recipe.Code) == "" {
					t.Fatalf("example %q has incomplete recipe: %#v", example.ID, recipe)
				}
			}
		}
	}
}

func TestBuildUsesExternalAssetsAndBackendLinks(t *testing.T) {
	output, err := Build(Config{
		BackendName: "Nuri",
		HTMLFile:    "showcase.html",
		OtherName:   "Chroma",
		OtherHref:   "../chroma/showcase-chroma.html",
	}, mockHighlighter{})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	for _, expected := range []string{
		`<link rel="stylesheet" href="showcase.css">`,
		`<script defer src="showcase.js"></script>`,
		`href="../chroma/showcase-chroma.html"`,
		`data-backend-link`,
		`data-category-link="frames"`,
	} {
		if !strings.Contains(output.HTML, expected) {
			t.Errorf("HTML missing %q", expected)
		}
	}
	if strings.Contains(output.HTML, "<style") {
		t.Error("HTML contains an embedded style element")
	}
	if strings.Contains(output.HTML, "<script>") {
		t.Error("HTML contains an inline behavior script")
	}
	if strings.Contains(output.HTML, `type="module"`) {
		t.Error("HTML loads showcase.js as a module, which breaks direct file:// use")
	}

	for _, expected := range []string{".kazari-code", ".showcase-filters", ".recipe-disclosure", ".nav-category.is-active", ".dark .clear-filters", ".clear-filters:focus-visible"} {
		if !strings.Contains(output.CSS, expected) {
			t.Errorf("CSS missing %q", expected)
		}
	}
	for _, expected := range []string{"kz-copy-btn", "example-search", "data-copy-link", "data-recipe-tab", "setCopiedState", "updateScrollSpy", "aria-current"} {
		if !strings.Contains(output.JS, expected) {
			t.Errorf("JS missing %q", expected)
		}
	}
}

func TestGenerateWritesThreeFiles(t *testing.T) {
	dir := t.TempDir()
	config := Config{
		BackendName: "Chroma",
		HTMLFile:    "showcase-chroma.html",
		OtherName:   "Nuri",
		OtherHref:   "../nuri/showcase.html",
	}
	if err := Generate(dir, config, mockHighlighter{}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for _, name := range []string{"showcase-chroma.html", "showcase.css", "showcase.js"} {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("stat %s: %v", name, err)
		}
		if info.Size() == 0 {
			t.Fatalf("%s is empty", name)
		}
	}
}

func exampleIDs(categories []Category) []string {
	var ids []string
	for _, category := range categories {
		for _, example := range category.Examples {
			ids = append(ids, example.ID)
		}
	}
	return ids
}
