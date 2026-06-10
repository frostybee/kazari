package meta

import (
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func TestParse_SingleQuote_BareInlineMarker(t *testing.T) {
	result := Parse("go 'useState'")
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Text != "useState" {
		t.Errorf("expected text 'useState', got %q", result.InlineMarkers[0].Text)
	}
	if result.InlineMarkers[0].Type != config.MarkerMark {
		t.Errorf("expected MarkerMark, got %d", result.InlineMarkers[0].Type)
	}
}

func TestParse_SingleQuote_InsInlineMarker(t *testing.T) {
	result := Parse("go ins='added'")
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Text != "added" {
		t.Errorf("expected text 'added', got %q", result.InlineMarkers[0].Text)
	}
	if result.InlineMarkers[0].Type != config.MarkerIns {
		t.Errorf("expected MarkerIns, got %d", result.InlineMarkers[0].Type)
	}
}

func TestParse_SingleQuote_DelInlineMarker(t *testing.T) {
	result := Parse("go del='removed'")
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Text != "removed" {
		t.Errorf("expected text 'removed', got %q", result.InlineMarkers[0].Text)
	}
	if result.InlineMarkers[0].Type != config.MarkerDel {
		t.Errorf("expected MarkerDel, got %d", result.InlineMarkers[0].Type)
	}
}

func TestParse_SingleQuote_EscapedQuote(t *testing.T) {
	result := Parse(`go 'it\'s'`)
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Text != "it's" {
		t.Errorf("expected text \"it's\", got %q", result.InlineMarkers[0].Text)
	}
}

func TestParse_SingleQuote_MixedWithDoubleQuotes(t *testing.T) {
	result := Parse(`go "double" 'single'`)
	if len(result.InlineMarkers) != 2 {
		t.Fatalf("expected 2 inline markers, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Text != "double" {
		t.Errorf("expected first marker text 'double', got %q", result.InlineMarkers[0].Text)
	}
	if result.InlineMarkers[1].Text != "single" {
		t.Errorf("expected second marker text 'single', got %q", result.InlineMarkers[1].Text)
	}
}

func TestParse_SingleQuote_LabeledRange(t *testing.T) {
	result := Parse("go {'A':3-5}")
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
	if result.LineMarkers[0].Label != "A" {
		t.Errorf("expected label 'A', got %q", result.LineMarkers[0].Label)
	}
	if len(result.LineMarkers[0].Lines) != 1 || result.LineMarkers[0].Lines[0].Start != 3 || result.LineMarkers[0].Lines[0].End != 5 {
		t.Errorf("expected range 3-5, got %+v", result.LineMarkers[0].Lines)
	}
}

func TestParse_SingleQuote_NotMistakenForLang(t *testing.T) {
	result := Parse("'text'")
	if result.BlockOptions.Lang != "" {
		t.Errorf("single-quoted string should not be parsed as language, got %q", result.BlockOptions.Lang)
	}
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
}

func TestParse_SingleQuote_InsLabeledRange(t *testing.T) {
	result := Parse("go ins={'Added':6-10}")
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
	if result.LineMarkers[0].Label != "Added" {
		t.Errorf("expected label 'Added', got %q", result.LineMarkers[0].Label)
	}
	if result.LineMarkers[0].Type != config.MarkerIns {
		t.Errorf("expected MarkerIns, got %d", result.LineMarkers[0].Type)
	}
}
