package main

import (
	"bytes"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestNoArgsShowsUsage(t *testing.T) {
	code, _, stderr := runCLI(t)
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(stderr, "Usage:") {
		t.Fatal("usage missing from stderr")
	}
}

func TestUnknownCommand(t *testing.T) {
	code, _, stderr := runCLI(t, "frobnicate")
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if !strings.Contains(stderr, `unknown command "frobnicate"`) {
		t.Fatalf("stderr: %s", stderr)
	}
}

func TestHelp(t *testing.T) {
	code, stdout, _ := runCLI(t, "help")
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "kazari process") {
		t.Fatal("usage missing from stdout")
	}
}

func TestVersion(t *testing.T) {
	code, stdout, _ := runCLI(t, "version")
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.HasPrefix(stdout, "kazari ") {
		t.Fatalf("stdout: %q", stdout)
	}
}
