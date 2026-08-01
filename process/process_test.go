package process

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/frostybee/kazari"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

var updateGolden = os.Getenv("KAZARI_UPDATE_GOLDEN") == "1"

var (
	engineOnce sync.Once
	testEngine *kazari.Engine
)

// engine returns one shared default configured Engine backed by the real
// Nuri highlighter, so corpus goldens reflect production output.
func engine(t *testing.T) *kazari.Engine {
	t.Helper()
	engineOnce.Do(func() {
		ctx := context.Background()
		hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
		if err != nil {
			panic(err)
		}
		testEngine = kazari.New(kazari.WithHighlighter(kazarinuri.New(ctx, hl)))
	})
	return testEngine
}

func newProcessor(t *testing.T, cfg Config) *Processor {
	t.Helper()
	if cfg.Engine == nil {
		cfg.Engine = engine(t)
	}
	p, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// copyTree copies a testdata input tree into a temp root on the real disk.
func copyTree(t *testing.T, from, to string) {
	t.Helper()
	err := filepath.WalkDir(from, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(from, path)
		dest := filepath.Join(to, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		return os.WriteFile(dest, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestCorpus processes each corpus case end to end on the real filesystem,
// compares every HTML file against the committed golden tree, verifies the
// emitted assets, and then runs a second time asserting a complete no-op.
// Set KAZARI_UPDATE_GOLDEN=1 to regenerate goldens after an intentional
// rendering change; regenerated goldens must be reviewed by eye before
// committing, since the corpus is the spec.
func TestCorpus(t *testing.T) {
	cases, err := os.ReadDir(filepath.Join("testdata", "corpus"))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cases {
		if !c.IsDir() {
			continue
		}
		t.Run(c.Name(), func(t *testing.T) {
			caseDir := filepath.Join("testdata", "corpus", c.Name())
			root := t.TempDir()
			copyTree(t, filepath.Join(caseDir, "input"), root)

			p := newProcessor(t, Config{})
			result, rerr := p.Run(context.Background(), root)
			if rerr != nil {
				t.Fatalf("first run: %v", rerr)
			}
			for _, fr := range result.Files {
				if fr.Err != nil {
					t.Fatalf("%s: %v", fr.Path, fr.Err)
				}
			}

			// Assets must exist at the root and match the engine exactly.
			for _, name := range []string{"kazari.css", "kazari.js"} {
				data, aerr := os.ReadFile(filepath.Join(root, name))
				if aerr != nil {
					t.Fatalf("asset %s: %v", name, aerr)
				}
				want := p.cfg.Engine.CSS()
				if name == "kazari.js" {
					want = p.cfg.Engine.JS()
				}
				if string(data) != want {
					t.Fatalf("asset %s does not match engine output", name)
				}
			}

			goldenDir := filepath.Join(caseDir, "golden")
			if updateGolden {
				if uerr := os.RemoveAll(goldenDir); uerr != nil {
					t.Fatal(uerr)
				}
			}
			werr := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				rel, _ := filepath.Rel(root, path)
				if !strings.HasSuffix(rel, ".html") && !strings.HasSuffix(rel, ".htm") {
					return nil
				}
				got, rerr := os.ReadFile(path)
				if rerr != nil {
					return rerr
				}
				goldenPath := filepath.Join(goldenDir, rel)
				if updateGolden {
					if merr := os.MkdirAll(filepath.Dir(goldenPath), 0o755); merr != nil {
						return merr
					}
					return os.WriteFile(goldenPath, got, 0o644)
				}
				want, gerr := os.ReadFile(goldenPath)
				if gerr != nil {
					t.Fatalf("missing golden for %s (run with KAZARI_UPDATE_GOLDEN=1 once, then review)", rel)
				}
				if string(got) != string(want) {
					i := 0
					for i < len(got) && i < len(want) && got[i] == want[i] {
						i++
					}
					lo := i - 60
					if lo < 0 {
						lo = 0
					}
					t.Fatalf("%s differs from golden at byte %d\n got: %q\nwant: %q", rel, i,
						got[lo:min(i+60, len(got))], want[lo:min(i+60, len(want))])
				}
				return nil
			})
			if werr != nil {
				t.Fatal(werr)
			}

			// Second run must be a complete no-op.
			second, serr := p.Run(context.Background(), root)
			if serr != nil {
				t.Fatalf("second run: %v", serr)
			}
			if second.ChangedCount != 0 {
				for _, fr := range second.Files {
					if fr.Changed {
						t.Errorf("second run changed %s", fr.Path)
					}
				}
				t.Fatalf("second run ChangedCount = %d, want 0", second.ChangedCount)
			}
		})
	}
}

func TestCheckModeWritesNothing(t *testing.T) {
	fs := newMemFS()
	root := filepath.Join("site")
	page := filepath.Join(root, "index.html")
	fs.put(page, []byte(`<html><head><title>x</title></head><body><pre><code class="language-go">package main</code></pre></body></html>`))

	p := newProcessor(t, Config{Check: true, FS: fs})
	result, err := p.Run(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if fs.writeCount() != 0 {
		t.Fatalf("Check mode performed %d writes", fs.writeCount())
	}
	// Both assets would be created and the page would change.
	if result.ChangedCount != 3 {
		t.Fatalf("ChangedCount = %d, want 3", result.ChangedCount)
	}
	for _, a := range result.Assets {
		if a.Action != "created" {
			t.Fatalf("asset %s action %q, want created", a.Path, a.Action)
		}
	}
	if len(result.Files) != 1 || !result.Files[0].Changed {
		t.Fatal("page must report as would-change")
	}
}

func TestUnlabeledPolicy(t *testing.T) {
	src := []byte(`<html><body><pre><code>no language here</code></pre></body></html>`)
	for _, skip := range []bool{true, false} {
		fs := newMemFS()
		root := filepath.Join("site")
		fs.put(filepath.Join(root, "index.html"), src)
		p := newProcessor(t, Config{SkipUnlabeled: skip, FS: fs})
		result, err := p.Run(context.Background(), root)
		if err != nil {
			t.Fatal(err)
		}
		fr := result.Files[0]
		if skip {
			if fr.Changed || fr.BlocksRewritten != 0 {
				t.Fatal("SkipUnlabeled must leave the block untouched")
			}
			if len(fr.BlocksSkipped) != 1 || fr.BlocksSkipped[0] != "unlabeled" {
				t.Fatalf("skip reasons %v", fr.BlocksSkipped)
			}
		} else {
			if !fr.Changed || fr.BlocksRewritten != 1 {
				t.Fatal("default policy must render the unlabeled block as plain text")
			}
		}
	}
}

func TestStrayAssetsAndBinariesIgnored(t *testing.T) {
	fs := newMemFS()
	root := filepath.Join("site")
	fs.put(filepath.Join(root, "kazari.css"), []byte("stale css"))
	fs.put(filepath.Join(root, "kazari.js"), []byte("stale js"))
	fs.put(filepath.Join(root, "image.png"), []byte{0x89, 0x50, 0x4E, 0x47})
	fs.put(filepath.Join(root, "index.html"), []byte(`<html><body><p>no code</p></body></html>`))

	p := newProcessor(t, Config{FS: fs})
	result, err := p.Run(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Files) != 1 {
		t.Fatalf("processed %d files, want only index.html", len(result.Files))
	}
	// The stale assets get refreshed as assets, never parsed as pages.
	for _, a := range result.Assets {
		if a.Action != "updated" {
			t.Fatalf("asset %s action %q, want updated", a.Path, a.Action)
		}
	}
	if data, _ := fs.ReadFile(filepath.Join(root, "image.png")); string(data) != string([]byte{0x89, 0x50, 0x4E, 0x47}) {
		t.Fatal("binary file must never be touched")
	}
}

func TestMaxFileBytesGuard(t *testing.T) {
	fs := newMemFS()
	root := filepath.Join("site")
	big := `<html><body><pre><code class="language-go">` + strings.Repeat("x", 512) + `</code></pre></body></html>`
	fs.put(filepath.Join(root, "big.html"), []byte(big))

	p := newProcessor(t, Config{FS: fs, MaxFileBytes: 64})
	result, err := p.Run(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	fr := result.Files[0]
	if fr.Err == nil {
		t.Fatal("oversized file must carry an error")
	}
	if fr.Changed {
		t.Fatal("oversized file must stay untouched")
	}
	if data, _ := fs.ReadFile(filepath.Join(root, "big.html")); string(data) != big {
		t.Fatal("oversized file bytes must be preserved")
	}
}

// errHighlighter claims go as a loaded language and then fails to tokenize
// it, which is the only way to force a real render error out of the engine:
// unknown languages fall back to plaintext instead of erroring.
type errHighlighter struct{}

func (errHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
	return nil, os.ErrInvalid
}
func (errHighlighter) GetThemeColors(theme string) (kazari.ThemeInfo, error) {
	return kazari.ThemeInfo{FG: "#24292f", BG: "#ffffff"}, nil
}
func (errHighlighter) GetLoadedLanguages() []string { return []string{"go"} }

func TestRenderErrorLeavesBlockUntouched(t *testing.T) {
	fs := newMemFS()
	root := filepath.Join("site")
	src := []byte(`<html><body><pre><code class="language-go">package main</code></pre></body></html>`)
	fs.put(filepath.Join(root, "index.html"), src)

	p := newProcessor(t, Config{
		Engine: kazari.New(kazari.WithHighlighter(errHighlighter{})),
		FS:     fs,
	})
	result, err := p.Run(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	fr := result.Files[0]
	if fr.Changed || fr.BlocksRewritten != 0 {
		t.Fatal("render error must leave the block untouched")
	}
	if len(fr.BlocksSkipped) != 1 || fr.BlocksSkipped[0] != "render-error" {
		t.Fatalf("skip reasons %v", fr.BlocksSkipped)
	}
	if data, _ := fs.ReadFile(filepath.Join(root, "index.html")); string(data) != string(src) {
		t.Fatal("page bytes must be preserved on render error")
	}
}

func TestContextCancellation(t *testing.T) {
	fs := newMemFS()
	root := filepath.Join("site")
	for _, name := range []string{"a.html", "b.html", "c.html"} {
		fs.put(filepath.Join(root, name), []byte(`<html><body><pre><code class="language-go">package main</code></pre></body></html>`))
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := newProcessor(t, Config{FS: fs, Check: true})
	result, err := p.Run(ctx, root)
	if err == nil {
		t.Fatal("canceled context must surface an error")
	}
	if len(result.Files) != 3 {
		t.Fatalf("partial result must still carry %d slots, got %d", 3, len(result.Files))
	}
}

func TestAssetsBaseOverride(t *testing.T) {
	fs := newMemFS()
	root := filepath.Join("site")
	fs.put(filepath.Join(root, "sub", "page.html"), []byte(`<html><head><title>x</title></head><body><pre><code class="language-go">package main</code></pre></body></html>`))

	p := newProcessor(t, Config{FS: fs, AssetsBase: "/assets/"})
	if _, err := p.Run(context.Background(), root); err != nil {
		t.Fatal(err)
	}
	out, _ := fs.ReadFile(filepath.Join(root, "sub", "page.html"))
	if !strings.Contains(string(out), `href="/assets/kazari.css?v=`) {
		t.Fatalf("expected AssetsBase href, got: %.300s", out)
	}
	if strings.Contains(string(out), "../") {
		t.Fatal("AssetsBase must suppress relative path computation")
	}
}

func TestByteExactOutsideRegions(t *testing.T) {
	// Ugly but legal page: unquoted attributes, a script whose string
	// literal contains a close tag lookalike, comments adjacent to the
	// block. Everything outside the replaced region must survive byte
	// exactly.
	prefix := `<html><head><script>var s = "</div>" + '<pre>';</script></head><body class=dark data-x=1><!-- before -->`
	block := `<pre><code class="language-go">package main</code></pre>`
	suffix := `<!-- after --><p unquoted=yes>tail</p></body></html>`
	fs := newMemFS()
	root := filepath.Join("site")
	fs.put(filepath.Join(root, "index.html"), []byte(prefix+block+suffix))

	p := newProcessor(t, Config{FS: fs})
	if _, err := p.Run(context.Background(), root); err != nil {
		t.Fatal(err)
	}
	out, _ := fs.ReadFile(filepath.Join(root, "index.html"))
	s := string(out)
	// The stylesheet link lands before the head close, splitting the prefix.
	head := `<html><head><script>var s = "</div>" + '<pre>';</script>`
	if !strings.HasPrefix(s, head) {
		t.Fatalf("prefix before the link injection altered: %.200s", s)
	}
	if !strings.Contains(s, `<!-- before --><div class="kazari-block`) {
		t.Fatal("bytes between body start and the block altered")
	}
	if !strings.Contains(s, `<!-- after --><p unquoted=yes>tail</p>`) {
		t.Fatal("suffix bytes altered")
	}
	if strings.Contains(s, block) {
		t.Fatal("original block should have been replaced")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
