package process

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestHugoE2E builds the two fixture sites under testdata/hugo-fixture with
// a real hugo binary, runs the processor over the combined output, and
// verifies the tier 1 promise end to end: blocks rewritten, kz markup and
// asset tags present, assets emitted once, and a clean check pass on the
// second run. The test skips when hugo is not installed, so the corpus
// golden tests remain the portable baseline and this is the live proof.
func TestHugoE2E(t *testing.T) {
	hugoPath, err := exec.LookPath("hugo")
	if err != nil {
		t.Skip("hugo not found in PATH; skipping live end to end test")
	}

	ctx := context.Background()
	publicDir := filepath.Join(t.TempDir(), "public")

	runHugo := func(site, dest string) {
		t.Helper()
		cmd := exec.CommandContext(ctx, hugoPath, "--minify", "--source", site, "--destination", dest)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("hugo build of %s failed: %v\n%s", site, err, out)
		}
	}
	fixtureRoot := filepath.Join("testdata", "hugo-fixture")
	runHugo(filepath.Join(fixtureRoot, "plain-site"), filepath.Join(publicDir, "plain"))
	runHugo(filepath.Join(fixtureRoot, "chroma-site"), filepath.Join(publicDir, "chroma"))
	runHugo(filepath.Join(fixtureRoot, "hook-site"), filepath.Join(publicDir, "hook"))

	p := newProcessor(t, Config{})
	result, err := p.Run(ctx, publicDir)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	rewritten := 0
	for _, fr := range result.Files {
		if fr.Err != nil {
			t.Errorf("%s: %v", fr.Path, fr.Err)
		}
		rewritten += fr.BlocksRewritten
	}
	if rewritten < 7 {
		t.Fatalf("BlocksRewritten total = %d, want at least 7 (one per authored page, two on the collapse page)", rewritten)
	}

	pages := []string{
		filepath.Join(publicDir, "plain", "plain-block", "index.html"),
		filepath.Join(publicDir, "chroma", "chroma-default", "index.html"),
		filepath.Join(publicDir, "chroma", "hl-lines", "index.html"),
		filepath.Join(publicDir, "hook", "title-and-markers", "index.html"),
		filepath.Join(publicDir, "hook", "collapse", "index.html"),
		filepath.Join(publicDir, "hook", "hl-lines-hook", "index.html"),
	}
	for _, page := range pages {
		html, err := os.ReadFile(page)
		if err != nil {
			t.Fatalf("read processed page: %v", err)
		}
		if !bytes.Contains(html, []byte("kz-")) {
			t.Errorf("%s: no kz markup after processing", page)
		}
		if !bytes.Contains(html, []byte(`data-kazari="assets"`)) {
			t.Errorf("%s: no injected asset tags", page)
		}
	}

	// The hl_lines page must carry the translated line markers, visible as
	// marked lines in the rendered output rather than Chroma hl classes.
	hlPage, err := os.ReadFile(pages[2])
	if err != nil {
		t.Fatalf("read hl lines page: %v", err)
	}
	if !bytes.Contains(hlPage, []byte(`class="kz-line highlight mark"`)) {
		t.Errorf("hl-lines page carries no marked line markup; hl_lines translation may have failed")
	}

	// Tier 2: title, mark, ins, and del survive purely through the render
	// hook's data-kz-meta attribute, with no Chroma translation involved.
	titleMarkersPage, err := os.ReadFile(pages[3])
	if err != nil {
		t.Fatalf("read title-and-markers page: %v", err)
	}
	for _, want := range []string{
		`<span class="kz-title">main.go</span>`,
		`class="kz-line highlight mark"`,
		`class="kz-line highlight ins"`,
		`class="kz-line highlight del"`,
	} {
		if !bytes.Contains(titleMarkersPage, []byte(want)) {
			t.Errorf("title-and-markers page missing %q", want)
		}
	}

	// Tier 2 collapse: the default style and the collapsible-start style
	// both come from data-kz-meta, no threshold collapse involved.
	collapsePage, err := os.ReadFile(pages[4])
	if err != nil {
		t.Fatalf("read collapse page: %v", err)
	}
	if !bytes.Contains(collapsePage, []byte(`<details class="kz-section">`)) {
		t.Error("collapse page missing default style collapse markup")
	}
	if !bytes.Contains(collapsePage, []byte(`class="kz-section collapsible-start"`)) {
		t.Error("collapse page missing collapsible-start style markup")
	}

	// Tier 2 hl_lines: the hook template's own string form translation to
	// mark ranges, distinct from the Chroma hl class path covered above.
	hookHLPage, err := os.ReadFile(pages[5])
	if err != nil {
		t.Fatalf("read hook hl-lines page: %v", err)
	}
	if !bytes.Contains(hookHLPage, []byte(`class="kz-line highlight mark"`)) {
		t.Error("hook hl-lines page carries no marked line markup")
	}

	// The mermaid page proves the skip rule survives the hook: the hook
	// still emits a plain language-mermaid block and the processor must
	// leave it untouched rather than wrapping it in kazari markup.
	mermaidPage, err := os.ReadFile(filepath.Join(publicDir, "hook", "mermaid", "index.html"))
	if err != nil {
		t.Fatalf("read mermaid page: %v", err)
	}
	if !bytes.Contains(mermaidPage, []byte("language-mermaid")) {
		t.Error("mermaid page's plain block was altered; expected an untouched skip")
	}
	if bytes.Contains(mermaidPage, []byte(`class="kazari-block`)) {
		t.Error("mermaid block was wrapped in kazari markup; skip rule did not fire through the hook")
	}

	eng := engine(t)
	css, err := os.ReadFile(filepath.Join(publicDir, "kazari.css"))
	if err != nil {
		t.Fatalf("read emitted css: %v", err)
	}
	if string(css) != eng.CSS() {
		t.Error("emitted kazari.css does not match Engine.CSS()")
	}
	// The default CLI engine must ship collapse section styling, since meta
	// collapse ranges reach it through data-kz-meta without WithCollapsible.
	if !bytes.Contains(css, []byte("kz-section")) {
		t.Error("emitted kazari.css is missing collapse section styles")
	}
	js, err := os.ReadFile(filepath.Join(publicDir, "kazari.js"))
	if err != nil {
		t.Fatalf("read emitted js: %v", err)
	}
	if string(js) != eng.JS() {
		t.Error("emitted kazari.js does not match Engine.JS()")
	}

	checker := newProcessor(t, Config{Check: true})
	second, err := checker.Run(ctx, publicDir)
	if err != nil {
		t.Fatalf("check run: %v", err)
	}
	if second.ChangedCount != 0 {
		var changed []string
		for _, fr := range second.Files {
			if fr.Changed {
				changed = append(changed, fr.Path)
			}
		}
		for _, ar := range second.Assets {
			if ar.Action != "unchanged" {
				changed = append(changed, ar.Path+" ("+ar.Action+")")
			}
		}
		t.Fatalf("second run not idempotent, ChangedCount = %d:\n%s",
			second.ChangedCount, strings.Join(changed, "\n"))
	}
}
