package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSite(t *testing.T, pages map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range pages {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const plainBlockPage = `<html><head><title>x</title></head><body>
<pre><code class="language-go">package main</code></pre>
</body></html>`

const unlabeledBlockPage = `<html><head><title>x</title></head><body>
<pre><code>no language</code></pre>
</body></html>`

func TestProcessUnknownThemeFlag(t *testing.T) {
	root := writeSite(t, map[string]string{"index.html": plainBlockPage})
	code, _, stderr := runCLI(t, "process", root, "--theme-dark", "github-drak")
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(stderr, `did you mean "github-dark"?`) {
		t.Fatalf("stderr: %s", stderr)
	}
}

func TestProcessMissingDir(t *testing.T) {
	code, _, stderr := runCLI(t, "process", filepath.Join(t.TempDir(), "nope"))
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(stderr, "is not a directory") {
		t.Fatalf("stderr: %s", stderr)
	}
}

func TestProcessTooManyArgs(t *testing.T) {
	code, _, stderr := runCLI(t, "process", "a", "b")
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(stderr, "too many arguments") {
		t.Fatalf("stderr: %s", stderr)
	}
}

func TestProcessBadFlag(t *testing.T) {
	code, _, _ := runCLI(t, "process", "--no-such-flag")
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
}

func TestProcessMalformedConfig(t *testing.T) {
	root := writeSite(t, map[string]string{"index.html": plainBlockPage})
	cfg := filepath.Join(root, "kazari.config.yaml")
	if err := os.WriteFile(cfg, []byte("notAKey: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderr := runCLI(t, "process", root)
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(stderr, "kazari.config.yaml") {
		t.Fatalf("config path missing from error: %s", stderr)
	}
}

func TestProcessCheckThenApplyThenClean(t *testing.T) {
	root := writeSite(t, map[string]string{"index.html": plainBlockPage})

	code, stdout, stderr := runCLI(t, "process", root, "--check")
	if code != 1 {
		t.Fatalf("dirty check exit = %d, want 1\nstderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "index.html") {
		t.Fatalf("would-change path missing:\n%s", stdout)
	}
	if data, _ := os.ReadFile(filepath.Join(root, "index.html")); string(data) != plainBlockPage {
		t.Fatal("check mode must not write")
	}

	code, stdout, stderr = runCLI(t, "process", root)
	if code != 0 {
		t.Fatalf("apply exit = %d\nstderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "1 blocks upgraded") {
		t.Fatalf("summary: %s", stdout)
	}
	out, _ := os.ReadFile(filepath.Join(root, "index.html"))
	if !strings.Contains(string(out), "kazari-block") {
		t.Fatal("block not upgraded")
	}
	if _, err := os.Stat(filepath.Join(root, "kazari.css")); err != nil {
		t.Fatal("kazari.css missing")
	}

	code, stdout, _ = runCLI(t, "process", root, "--check")
	if code != 0 {
		t.Fatalf("clean check exit = %d, want 0\nstdout: %s", code, stdout)
	}
	if !strings.Contains(stdout, "0 changed") {
		t.Fatalf("summary: %s", stdout)
	}
}

func TestProcessConfigPrecedence(t *testing.T) {
	root := writeSite(t, map[string]string{"index.html": unlabeledBlockPage})
	cfg := filepath.Join(root, "kazari.config.yaml")
	if err := os.WriteFile(cfg, []byte("process:\n  skipUnlabeled: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Config block wins when the flag is absent: the unlabeled block stays.
	code, _, stderr := runCLI(t, "process", root)
	if code != 0 {
		t.Fatalf("exit = %d\nstderr: %s", code, stderr)
	}
	out, _ := os.ReadFile(filepath.Join(root, "index.html"))
	if strings.Contains(string(out), "kazari-block") {
		t.Fatal("config skipUnlabeled must leave the block untouched")
	}

	// An explicitly set flag beats the config block.
	code, _, stderr = runCLI(t, "process", root, "--skip-unlabeled=false")
	if code != 0 {
		t.Fatalf("exit = %d\nstderr: %s", code, stderr)
	}
	out, _ = os.ReadFile(filepath.Join(root, "index.html"))
	if !strings.Contains(string(out), "kazari-block") {
		t.Fatal("explicit flag must override the config block")
	}
}
