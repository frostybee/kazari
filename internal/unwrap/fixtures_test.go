package unwrap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// expectation mirrors expected.json. The match flag distinguishes a region
// no unwrapper claims from zero valued fields, and skip records which pre
// chain rule fired so the corpus stays reusable by the processor phase.
type expectation struct {
	Match bool   `json:"match"`
	Lang  string `json:"lang"`
	Meta  string `json:"meta"`
	Code  string `json:"code"`
	Skip  string `json:"skip"`
}

// loadFixture reads one fixture directory. os.ReadFile returns the exact
// file bytes; nothing here may normalize line endings, because the CRLF
// fixture pins that raw carriage returns reach the tokenizer.
func loadFixture(t *testing.T, dir string) ([]BufferedToken, expectation) {
	t.Helper()
	src, err := os.ReadFile(filepath.Join("testdata", dir, "input.html"))
	if err != nil {
		t.Fatalf("read input: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join("testdata", dir, "expected.json"))
	if err != nil {
		t.Fatalf("read expectation: %v", err)
	}
	var want expectation
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatalf("parse expectation: %v", err)
	}
	tokens, err := TokenizeFragment(src)
	if err != nil {
		t.Fatalf("tokenize: %v", err)
	}
	return tokens, want
}

// diffStrings fails with the byte index of the first mismatch, following the
// diff context pattern of demo/nuri/equivalence_test.go.
func diffStrings(t *testing.T, label, got, want string) {
	t.Helper()
	if got == want {
		return
	}
	i := 0
	for i < len(got) && i < len(want) && got[i] == want[i] {
		i++
	}
	t.Errorf("%s differs at byte %d\n got: %q\nwant: %q", label, i, got, want)
}

// assertRegion runs one unwrapper against one fixture and checks every
// field byte exactly.
func assertRegion(t *testing.T, u Unwrapper, dir string) {
	t.Helper()
	tokens, want := loadFixture(t, dir)
	region, ok := u.Match(tokens)
	if ok != want.Match {
		t.Fatalf("%s: Match = %v, want %v", dir, ok, want.Match)
	}
	if !ok {
		return
	}
	diffStrings(t, dir+" lang", region.Lang, want.Lang)
	diffStrings(t, dir+" meta", region.Meta, want.Meta)
	diffStrings(t, dir+" code", region.Code, want.Code)
}
