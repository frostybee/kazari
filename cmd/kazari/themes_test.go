package main

import (
	"sort"
	"strings"
	"testing"
)

func TestThemesListing(t *testing.T) {
	code, stdout, _ := runCLI(t, "themes")
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	names := strings.Fields(stdout)
	if len(names) < 10 {
		t.Fatalf("only %d themes listed", len(names))
	}
	if !sort.StringsAreSorted(names) {
		t.Fatal("theme names must be sorted")
	}
	set := map[string]bool{}
	for _, n := range names {
		set[n] = true
	}
	for _, want := range []string{"github-dark", "github-light"} {
		if !set[want] {
			t.Fatalf("missing %s in listing", want)
		}
	}
}

func TestNearestName(t *testing.T) {
	names := []string{"github-dark", "github-light", "dracula"}
	if got := nearestName("github-drak", names); got != "github-dark" {
		t.Fatalf("got %q", got)
	}
	if got := nearestName("drakula", names); got != "dracula" {
		t.Fatalf("got %q", got)
	}
	if got := nearestName("zzzzzzzz", names); got != "" {
		t.Fatalf("no suggestion expected, got %q", got)
	}
}

func TestValidateThemeNames(t *testing.T) {
	if err := validateThemeNames("github-light", "github-dark"); err != nil {
		t.Fatalf("valid pair rejected: %v", err)
	}
	err := validateThemeNames("github-light", "github-drak")
	if err == nil {
		t.Fatal("typo must be rejected")
	}
	if !strings.Contains(err.Error(), `did you mean "github-dark"?`) {
		t.Fatalf("suggestion missing: %v", err)
	}
	if !strings.Contains(err.Error(), "kazari themes") {
		t.Fatalf("listing hint missing: %v", err)
	}
}
