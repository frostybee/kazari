package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestExampleHugoSite builds examples/hugo with a real hugo binary and runs
// the CLI over the output, proving the shipped demo actually works rather
// than only looking right in review. It goes through run() rather than the
// process package so that flag parsing and the example's own config file
// are exercised: a config that fails to parse or validate would otherwise
// pass unnoticed, since a default engine would still render every block.
//
// The assertions target the markup the example is there to demonstrate,
// weighted toward the features the render hook only recently stopped
// dropping. The test skips when hugo is not installed.
func TestExampleHugoSite(t *testing.T) {
	hugoPath, err := exec.LookPath("hugo")
	if err != nil {
		t.Skip("hugo not found in PATH; skipping example site build")
	}

	site := filepath.Join("..", "..", "examples", "hugo")
	configPath := filepath.Join(site, "kazari.config.yaml")
	public := filepath.Join(t.TempDir(), "public")

	cmd := exec.CommandContext(context.Background(), hugoPath,
		"--minify", "--source", site, "--destination", public)
	if out, herr := cmd.CombinedOutput(); herr != nil {
		t.Fatalf("hugo build of %s failed: %v\n%s", site, herr, out)
	}

	code, stdout, stderr := runCLI(t, "process", "--config", configPath, public)
	if code != 0 {
		t.Fatalf("process exit = %d, want 0\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}

	pages := map[string][]string{
		"basics": {
			`<span class="kz-title">server.go</span>`,
			// languageIconMode: iconAndText opens the badge with the icon
			// slot, and fileIcons: false keeps the title bare, so this is
			// the only icon span on the page.
			`<span class="kz-lang-icon" data-lang="go"></span>`,
			`<figure class="frame is-terminal" data-lang="bash"`,
			// frame="none" renders no figure at all, so the wrapper is the
			// only evidence the block was upgraded.
			`class="kazari-block`,
			// The mermaid fence must survive untouched for a diagram
			// renderer to claim it.
			`language-mermaid`,
		},
		"annotations": {
			`class="kz-line highlight mark"`,
			`class="kz-line highlight ins"`,
			`class="kz-line highlight del"`,
			// focus reached the engine only after the hook learned to emit
			// it as a brace group.
			`class="kz-line focused"`,
			`class="kz-link" href="https://example.org/api"`,
		},
		"collapse": {
			`<details class="kz-section">`,
			`class="kz-section collapsible-start"`,
			`class="kz-section collapsible-end"`,
			`class="kz-collapse-toggle"`,
		},
		"themes": {
			`class="kazari-block kz-themed`,
			`class="kz-theme-toggle-btn"`,
		},
		"output": {
			// The whole output panel was unreachable through Hugo until the
			// hook translated the camelCase keys.
			`<div class="kz-output">`,
			`<div class="kz-output-header">`,
			`<pre class="kz-output-pre">`,
			`>Console</button>`,
			`<div class="kz-output kz-output-hidden">`,
		},
	}
	for page, wants := range pages {
		// uglyURLs is on, so pages land beside the home page rather than in
		// a directory of their own.
		html := readPage(t, filepath.Join(public, page+".html"))
		for _, want := range wants {
			if !strings.Contains(html, want) {
				t.Errorf("%s page missing %q", page, want)
			}
		}
	}

	// The example deliberately does not link kazari.css by hand, so the
	// injected tags and the emitted files are the only way the styles reach
	// the page.
	home := readPage(t, filepath.Join(public, "index.html"))
	for _, want := range []string{`data-kazari="assets"`, "kazari.css", "kazari.js"} {
		if !strings.Contains(home, want) {
			t.Errorf("home page missing %q", want)
		}
	}
	for _, asset := range []string{"kazari.css", "kazari.js"} {
		if _, serr := os.Stat(filepath.Join(public, asset)); serr != nil {
			t.Errorf("processor emitted no %s: %v", asset, serr)
		}
	}

	// Running the processor twice must be safe, which is what makes it
	// usable as an unconditional build step.
	code, stdout, stderr = runCLI(t, "process", "--check", "--config", configPath, public)
	if code != 0 {
		t.Fatalf("check exit = %d, want 0 (second run not idempotent)\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
}

func readPage(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read processed page: %v", err)
	}
	return string(data)
}
