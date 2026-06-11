package main

import (
	"context"
	"testing"

	"github.com/frostybee/kazari"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

// plainHighlighter hides NuriHighlighter's DualThemeTokenizer capability so
// the engine takes the two-pass fallback path. Methods are forwarded
// explicitly (embedding would forward TokenizeDual too).
type plainHighlighter struct {
	n *kazarinuri.NuriHighlighter
}

func (p plainHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
	return p.n.Tokenize(code, lang, theme)
}

func (p plainHighlighter) GetThemeColors(theme string) (kazari.ThemeInfo, error) {
	return p.n.GetThemeColors(theme)
}

func (p plainHighlighter) GetLoadedLanguages() []string {
	return p.n.GetLoadedLanguages()
}

// TestDualThemeEquivalence guards the TokenizeDual adapter's light/dark ↔
// default/ThemeStyles mapping: rendering through the single-pass capability
// path must produce byte-identical HTML to the two-pass fallback. Any color
// drift here means the mapping in NuriHighlighter.TokenizeDual is inverted
// or misaligned.
func TestDualThemeEquivalence(t *testing.T) {
	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		t.Fatal(err)
	}
	defer hl.Close(ctx)

	nuriHL := kazarinuri.New(ctx, hl)
	newEngine := func(h kazari.Highlighter) *kazari.Engine {
		return kazari.New(
			kazari.WithHighlighter(h),
			kazari.WithThemes("github-light", "github-dark"),
			kazari.WithMinify(false),
		)
	}
	dualEngine := newEngine(nuriHL)
	plainEngine := newEngine(plainHighlighter{n: nuriHL})

	samples := []struct {
		lang string
		code string
	}{
		{"go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\", 42)\n}\n"},
		{"javascript", "const x = 1;\nlet s = `tpl ${x}`;\n\nasync function f() { return x; }\n"},
		{"python", "def f(n):\n    return [i**2 for i in range(n)]\n"},
		{"markdown", "# Title\n\nSome **bold** and `code`.\n\n- item\n"},
		{"bash", "echo \"hello\" | grep -c h\n"},
		{"text", "plain text, no grammar\n"},
	}

	for _, s := range samples {
		opts := kazari.Options{Lang: s.lang, Title: "eq-test"}
		dual, err := dualEngine.Render(s.code, opts)
		if err != nil {
			t.Fatalf("%s (capability path): %v", s.lang, err)
		}
		plain, err := plainEngine.Render(s.code, opts)
		if err != nil {
			t.Fatalf("%s (two-pass path): %v", s.lang, err)
		}
		if dual != plain {
			t.Errorf("%s: capability and two-pass HTML differ", s.lang)
			max := len(dual)
			if len(plain) < max {
				max = len(plain)
			}
			for i := 0; i < max; i++ {
				if dual[i] != plain[i] {
					lo := i - 60
					if lo < 0 {
						lo = 0
					}
					hi := i + 60
					if hi > max {
						hi = max
					}
					t.Errorf("first diff at byte %d:\n  dual:  …%s…\n  plain: …%s…", i, dual[lo:hi], plain[lo:hi])
					break
				}
			}
		}
	}
}
