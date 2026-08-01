package process

import (
	"bytes"
	"testing"
)

func TestApplyEdits_NoEditsPassthrough(t *testing.T) {
	src := []byte(`<html><body>untouched</body></html>`)
	out, err := applyEdits(src, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, src) {
		t.Fatal("no edits must reproduce the input exactly")
	}
}

func TestApplyEdits_MultipleRegions(t *testing.T) {
	src := []byte("aaaBBBcccDDDeee")
	out, err := applyEdits(src, []edit{
		{start: 9, end: 12, replacement: []byte("Y")},
		{start: 3, end: 6, replacement: []byte("X")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "aaaXcccYeee" {
		t.Fatalf("got %q", out)
	}
}

func TestApplyEdits_InsertionBeforeReplacementAtSameOffset(t *testing.T) {
	// A tag injected exactly at a region boundary must land outside the
	// replaced span; insertions sort before replacements at equal start.
	src := []byte("aaaBBBccc")
	out, err := applyEdits(src, []edit{
		{start: 3, end: 6, replacement: []byte("R")},
		{start: 3, end: 3, replacement: []byte("<ins>")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "aaa<ins>Rccc" {
		t.Fatalf("got %q", out)
	}
}

func TestApplyEdits_OverlapRejected(t *testing.T) {
	src := []byte("aaaBBBccc")
	if _, err := applyEdits(src, []edit{
		{start: 2, end: 5, replacement: []byte("X")},
		{start: 4, end: 7, replacement: []byte("Y")},
	}); err == nil {
		t.Fatal("overlapping edits must abort the file, not corrupt it")
	}
}

func TestApplyEdits_InsertionAtEOF(t *testing.T) {
	src := []byte("abc")
	out, err := applyEdits(src, []edit{{start: 3, end: 3, replacement: []byte("<script>")}})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "abc<script>" {
		t.Fatalf("got %q", out)
	}
}
