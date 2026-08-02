package process

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/frostybee/kazari/internal/unwrap"
)

// parityExpectation mirrors the expected.json sidecar of each fixture in
// internal/unwrap/testdata. Those files record the known source, language,
// and meta string each fixture was authored from, which makes them the door
// one truth for the parity test.
type parityExpectation struct {
	Match bool   `json:"match"`
	Lang  string `json:"lang"`
	Meta  string `json:"meta"`
	Code  string `json:"code"`
	Skip  string `json:"skip"`
}

// TestParity renders every unwrap fixture through both front doors and byte
// compares the results. Door one renders the fixture's hand authored known
// Code and Meta directly through RenderWithMeta, exactly as the Goldmark
// path would. Door two recovers a Region from the fixture's input.html using
// the same Discover, SkipReason, chainFor, RunChain sequence processFile
// runs in production, then renders that.
//
// The Phase 3 roadmap entry says "for each corpus fixture"; that is read
// here as the fixture corpus built in Phase 0 (internal/unwrap/testdata),
// the only fixtures carrying known source sidecars. The process corpus under
// testdata/corpus has golden HTML trees only.
//
// Beyond the Phase 0 unit tests, which call one Unwrapper.Match directly,
// this test covers two gaps. First, chain routing: it asserts that Discover
// plus chainFor plus the full ordered RunChain recover the expected Region,
// so a chain ordering regression surfaces here. Second, engine acceptance:
// it is the only place the recovered Code and Meta pairs, including
// synthesized hl_lines ranges, data-kz-meta overrides, NBSP indentation,
// empty code, and empty language, are fed through the real engine. The final
// door one versus door two comparison is largely a corollary of the Region
// equality checks; it stays because it states the tier 1 contract directly:
// processing built HTML yields the same bytes as rendering the source.
func TestParity(t *testing.T) {
	fixtureRoot := filepath.Join("..", "internal", "unwrap", "testdata")
	entries, err := os.ReadDir(fixtureRoot)
	if err != nil {
		t.Fatal(err)
	}
	eng := engine(t)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(fixtureRoot, name)
			src, err := os.ReadFile(filepath.Join(dir, "input.html"))
			if err != nil {
				t.Fatalf("read input: %v", err)
			}
			raw, err := os.ReadFile(filepath.Join(dir, "expected.json"))
			if err != nil {
				t.Fatalf("read expectation: %v", err)
			}
			var want parityExpectation
			if err := json.Unmarshal(raw, &want); err != nil {
				t.Fatalf("parse expectation: %v", err)
			}

			tokens, err := unwrap.TokenizeFragment(src)
			if err != nil {
				t.Fatalf("tokenize: %v", err)
			}
			page := unwrap.Discover(src, tokens)

			// A reserved kz- or kazari- class on the fixture root suppresses
			// discovery before any candidate exists, so those fixtures never
			// reach SkipReason on this path.
			if len(page.Candidates) == 0 {
				if want.Match {
					t.Fatalf("no candidate discovered for a fixture marked match true")
				}
				if want.Skip == unwrap.SkipKzClass && page.Suppressed == 0 {
					t.Fatalf("expected reserved class suppression, got Suppressed = 0")
				}
				t.Skip("discovery yields no candidate; nothing to render")
			}
			if len(page.Candidates) != 1 {
				t.Fatalf("expected exactly 1 candidate, got %d", len(page.Candidates))
			}
			candidate := page.Candidates[0]
			if candidate.Malformed {
				t.Fatal("candidate malformed, fixture shape assumption violated")
			}

			if reason := unwrap.SkipReason(candidate.Tokens); reason != "" {
				if reason != want.Skip {
					t.Fatalf("SkipReason = %q, want %q", reason, want.Skip)
				}
				t.Skipf("skip rule %q fired; nothing to render", reason)
			}

			region, winner, ok := unwrap.RunChain(chainFor(candidate.Kind), candidate.Tokens)
			if !want.Match {
				if ok {
					t.Fatalf("RunChain matched via %s but fixture expects no match", winner)
				}
				t.Skip("no unwrapper claims this fixture")
			}
			if !ok {
				t.Fatal("RunChain found no match for a fixture marked match true; chain routing regression")
			}

			diffParity(t, name+" lang", region.Lang, want.Lang)
			diffParity(t, name+" meta", region.Meta, want.Meta)
			diffParity(t, name+" code", region.Code, want.Code)

			doorOne, err := eng.RenderWithMeta(want.Code, want.Meta)
			if err != nil {
				t.Fatalf("door one render from known source: %v", err)
			}
			doorTwo, err := eng.RenderWithMeta(region.Code, region.Meta)
			if err != nil {
				t.Fatalf("door two render from region recovered by %s: %v", winner, err)
			}
			diffParity(t, name+" door one vs door two", doorTwo, doorOne)
		})
	}
}

// diffParity fails with the byte index of the first mismatch, following the
// diff context pattern of demo/nuri/equivalence_test.go.
func diffParity(t *testing.T, label, got, want string) {
	t.Helper()
	if got == want {
		return
	}
	i := 0
	for i < len(got) && i < len(want) && got[i] == want[i] {
		i++
	}
	lo := i - 60
	if lo < 0 {
		lo = 0
	}
	gotHi, wantHi := i+60, i+60
	if gotHi > len(got) {
		gotHi = len(got)
	}
	if wantHi > len(want) {
		wantHi = len(want)
	}
	t.Errorf("%s differs at byte %d\n got: %q\nwant: %q", label, i, got[lo:gotHi], want[lo:wantHi])
}
