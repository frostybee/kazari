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

// --- Block option tokens ---

func TestParse_ShowLineNumbers(t *testing.T) {
	result := Parse("go showLineNumbers")
	if result.BlockOptions.LineNumbers == nil || !*result.BlockOptions.LineNumbers {
		t.Error("showLineNumbers should set LineNumbers to true")
	}
}

func TestParse_ShowLineNumbersFalse(t *testing.T) {
	result := Parse("go showLineNumbers=false")
	if result.BlockOptions.LineNumbers == nil || *result.BlockOptions.LineNumbers {
		t.Error("showLineNumbers=false should set LineNumbers to false")
	}
}

func TestParse_Wrap(t *testing.T) {
	result := Parse("go wrap")
	if result.BlockOptions.Wrap == nil || !*result.BlockOptions.Wrap {
		t.Error("wrap should set Wrap to true")
	}
}

func TestParse_PreserveIndent(t *testing.T) {
	result := Parse("go wrap preserveIndent")
	if result.BlockOptions.PreserveIndent == nil || !*result.BlockOptions.PreserveIndent {
		t.Error("preserveIndent should set PreserveIndent to true")
	}
}

func TestParse_PreserveIndentFalse(t *testing.T) {
	result := Parse("go wrap preserveIndent=false")
	if result.BlockOptions.PreserveIndent == nil || *result.BlockOptions.PreserveIndent {
		t.Error("preserveIndent=false should set PreserveIndent to false")
	}
}

func TestParse_HangingIndent(t *testing.T) {
	result := Parse("go wrap hangingIndent=2")
	if result.BlockOptions.HangingIndent == nil || *result.BlockOptions.HangingIndent != 2 {
		t.Error("hangingIndent=2 should set HangingIndent to 2")
	}
}

func TestParse_WrapCombo(t *testing.T) {
	result := Parse("go wrap preserveIndent=false hangingIndent=4")
	if result.BlockOptions.Wrap == nil || !*result.BlockOptions.Wrap {
		t.Error("expected Wrap=true")
	}
	if result.BlockOptions.PreserveIndent == nil || *result.BlockOptions.PreserveIndent {
		t.Error("expected PreserveIndent=false")
	}
	if result.BlockOptions.HangingIndent == nil || *result.BlockOptions.HangingIndent != 4 {
		t.Error("expected HangingIndent=4")
	}
}

func TestParse_Title(t *testing.T) {
	result := Parse(`go title="My Title"`)
	if result.BlockOptions.Title != "My Title" {
		t.Errorf("Title = %q, want %q", result.BlockOptions.Title, "My Title")
	}
}

func TestParse_FrameTerminal(t *testing.T) {
	result := Parse(`bash frame=terminal`)
	if result.BlockOptions.Frame == nil || *result.BlockOptions.Frame != config.FrameTerminal {
		t.Error("frame=terminal should set Frame to FrameTerminal")
	}
}

func TestParse_StartLineNumber(t *testing.T) {
	result := Parse("go startLineNumber=5")
	if result.BlockOptions.StartLineNumber == nil || *result.BlockOptions.StartLineNumber != 5 {
		t.Errorf("StartLineNumber = %v, want 5", result.BlockOptions.StartLineNumber)
	}
}

// --- Line markers ---

func TestParse_BareMarkerRanges(t *testing.T) {
	result := Parse("go {3-5,8}")
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
	m := result.LineMarkers[0]
	if m.Type != config.MarkerMark {
		t.Errorf("Type = %d, want MarkerMark", m.Type)
	}
	if len(m.Lines) != 2 {
		t.Fatalf("expected 2 ranges, got %d", len(m.Lines))
	}
	if m.Lines[0].Start != 3 || m.Lines[0].End != 5 {
		t.Errorf("first range = %d-%d, want 3-5", m.Lines[0].Start, m.Lines[0].End)
	}
	if m.Lines[1].Start != 8 || m.Lines[1].End != 8 {
		t.Errorf("second range = %d-%d, want 8-8", m.Lines[1].Start, m.Lines[1].End)
	}
}

func TestParse_InsLineMarker(t *testing.T) {
	result := Parse("go ins={10-12}")
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
	if result.LineMarkers[0].Type != config.MarkerIns {
		t.Errorf("Type = %d, want MarkerIns", result.LineMarkers[0].Type)
	}
}

func TestParse_DelLineMarker(t *testing.T) {
	result := Parse("go del={7}")
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
	if result.LineMarkers[0].Type != config.MarkerDel {
		t.Errorf("Type = %d, want MarkerDel", result.LineMarkers[0].Type)
	}
}

func TestParse_AddAlias(t *testing.T) {
	result := Parse("go add={1-3}")
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
	if result.LineMarkers[0].Type != config.MarkerIns {
		t.Errorf("add= should map to MarkerIns, got %d", result.LineMarkers[0].Type)
	}
}

// --- Focus lines ---

func TestParse_Focus(t *testing.T) {
	result := Parse("go focus={1-3,7}")
	if len(result.FocusLines) != 2 {
		t.Fatalf("expected 2 focus ranges, got %d", len(result.FocusLines))
	}
	if result.FocusLines[0].Start != 1 || result.FocusLines[0].End != 3 {
		t.Errorf("first focus = %d-%d, want 1-3", result.FocusLines[0].Start, result.FocusLines[0].End)
	}
	if result.FocusLines[1].Start != 7 || result.FocusLines[1].End != 7 {
		t.Errorf("second focus = %d-%d, want 7-7", result.FocusLines[1].Start, result.FocusLines[1].End)
	}
}

func TestParse_FocusSingleLine(t *testing.T) {
	result := Parse("go focus={5}")
	if len(result.FocusLines) != 1 {
		t.Fatalf("expected 1 focus range, got %d", len(result.FocusLines))
	}
	if result.FocusLines[0].Start != 5 || result.FocusLines[0].End != 5 {
		t.Errorf("focus = %d-%d, want 5-5", result.FocusLines[0].Start, result.FocusLines[0].End)
	}
}

// --- Collapse tokens ---

func TestParse_Collapse(t *testing.T) {
	result := Parse("go collapse")
	if result.Collapse == nil || !result.Collapse.Enabled {
		t.Error("collapse should set Enabled=true")
	}
}

func TestParse_NoCollapse(t *testing.T) {
	result := Parse("go nocollapse")
	if result.Collapse == nil || !result.Collapse.Disabled {
		t.Error("nocollapse should set Disabled=true")
	}
}

func TestParse_CollapseRanges(t *testing.T) {
	result := Parse("go collapse={2-5}")
	if result.Collapse == nil {
		t.Fatal("Collapse should not be nil")
	}
	if len(result.Collapse.Ranges) != 1 {
		t.Fatalf("expected 1 collapse range, got %d", len(result.Collapse.Ranges))
	}
	if result.Collapse.Ranges[0].Start != 2 || result.Collapse.Ranges[0].End != 5 {
		t.Errorf("range = %d-%d, want 2-5", result.Collapse.Ranges[0].Start, result.Collapse.Ranges[0].End)
	}
}

func TestParse_CollapseStyle(t *testing.T) {
	result := Parse("go collapseStyle=collapsible-start")
	if result.Collapse == nil || result.Collapse.Style == nil {
		t.Fatal("Collapse.Style should not be nil")
	}
	if *result.Collapse.Style != config.CollapseCollapsibleStart {
		t.Errorf("Style = %d, want CollapseCollapsibleStart", *result.Collapse.Style)
	}
}

func TestParse_CollapseThreshold(t *testing.T) {
	result := Parse("go collapseThreshold=20")
	if result.Collapse == nil || result.Collapse.Threshold == nil {
		t.Fatal("Collapse.Threshold should not be nil")
	}
	if *result.Collapse.Threshold != 20 {
		t.Errorf("Threshold = %d, want 20", *result.Collapse.Threshold)
	}
}

func TestParse_CollapseThreshold_Invalid(t *testing.T) {
	for _, meta := range []string{"go collapseThreshold=0", "go collapseThreshold=-5", "go collapseThreshold=abc"} {
		result := Parse(meta)
		if result.Collapse != nil && result.Collapse.Threshold != nil {
			t.Errorf("Parse(%q): Threshold should be nil for invalid value", meta)
		}
	}
}

func TestParse_CollapseThreshold_WithStyle(t *testing.T) {
	result := Parse("go collapseThreshold=20 collapseStyle=collapsible-start")
	if result.Collapse == nil {
		t.Fatal("Collapse should not be nil")
	}
	if result.Collapse.Threshold == nil || *result.Collapse.Threshold != 20 {
		t.Error("Threshold should be 20")
	}
	if result.Collapse.Style == nil || *result.Collapse.Style != config.CollapseCollapsibleStart {
		t.Error("Style should be collapsible-start")
	}
}

// --- Combined meta strings ---

func TestParse_Combined_AllOptions(t *testing.T) {
	result := Parse(`go title="main.go" {3-5} showLineNumbers`)
	if result.BlockOptions.Lang != "go" {
		t.Errorf("Lang = %q, want %q", result.BlockOptions.Lang, "go")
	}
	if result.BlockOptions.Title != "main.go" {
		t.Errorf("Title = %q, want %q", result.BlockOptions.Title, "main.go")
	}
	if result.BlockOptions.LineNumbers == nil || !*result.BlockOptions.LineNumbers {
		t.Error("LineNumbers should be true")
	}
	if len(result.LineMarkers) != 1 {
		t.Fatalf("expected 1 line marker, got %d", len(result.LineMarkers))
	}
}

func TestParse_Combined_MarkersAndFocus(t *testing.T) {
	result := Parse("bash ins={1-3} del={5} focus={8-10}")
	if len(result.LineMarkers) != 2 {
		t.Fatalf("expected 2 line markers, got %d", len(result.LineMarkers))
	}
	if len(result.FocusLines) != 1 {
		t.Fatalf("expected 1 focus range, got %d", len(result.FocusLines))
	}
}

func TestParse_Combined_MultipleOptions(t *testing.T) {
	result := Parse(`typescript wrap frame=none startLineNumber=10`)
	if result.BlockOptions.Lang != "typescript" {
		t.Errorf("Lang = %q, want %q", result.BlockOptions.Lang, "typescript")
	}
	if result.BlockOptions.Wrap == nil || !*result.BlockOptions.Wrap {
		t.Error("Wrap should be true")
	}
	if result.BlockOptions.Frame == nil || *result.BlockOptions.Frame != config.FrameNone {
		t.Error("Frame should be FrameNone")
	}
	if result.BlockOptions.StartLineNumber == nil || *result.BlockOptions.StartLineNumber != 10 {
		t.Error("StartLineNumber should be 10")
	}
}

// --- Edge cases ---

func TestParse_LanguageOnly(t *testing.T) {
	result := Parse("go")
	if result.BlockOptions.Lang != "go" {
		t.Errorf("Lang = %q, want %q", result.BlockOptions.Lang, "go")
	}
	if len(result.LineMarkers) != 0 {
		t.Errorf("expected no markers, got %d", len(result.LineMarkers))
	}
}

func TestParse_EmptyString(t *testing.T) {
	result := Parse("")
	if result.BlockOptions.Lang != "" {
		t.Errorf("empty string should produce empty Lang, got %q", result.BlockOptions.Lang)
	}
}

func TestParse_ExtraWhitespace(t *testing.T) {
	result := Parse("  go   showLineNumbers  ")
	if result.BlockOptions.Lang != "go" {
		t.Errorf("Lang = %q, want %q", result.BlockOptions.Lang, "go")
	}
	if result.BlockOptions.LineNumbers == nil || !*result.BlockOptions.LineNumbers {
		t.Error("showLineNumbers should be parsed despite extra whitespace")
	}
}

func TestParse_FrameDefaultUnknown(t *testing.T) {
	result := Parse("go frame=unknown")
	if result.BlockOptions.Frame == nil || *result.BlockOptions.Frame != config.FrameAuto {
		t.Error("unknown frame value should default to FrameAuto")
	}
}

// --- Regex markers ---

func TestParse_RegexPattern(t *testing.T) {
	result := Parse(`go /func\s+\w+/`)
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	m := result.InlineMarkers[0]
	if !m.IsRegex {
		t.Error("IsRegex should be true")
	}
	if m.Text != `func\s+\w+` {
		t.Errorf("Text = %q, want %q", m.Text, `func\s+\w+`)
	}
	if m.Type != config.MarkerMark {
		t.Errorf("Type = %d, want MarkerMark", m.Type)
	}
}

func TestParse_InsRegex(t *testing.T) {
	result := Parse(`go ins=/added\s+/`)
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if !result.InlineMarkers[0].IsRegex {
		t.Error("IsRegex should be true")
	}
	if result.InlineMarkers[0].Type != config.MarkerIns {
		t.Errorf("Type = %d, want MarkerIns", result.InlineMarkers[0].Type)
	}
}

func TestParse_DelRegex(t *testing.T) {
	result := Parse(`go del=/old\w+/`)
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Type != config.MarkerDel {
		t.Errorf("Type = %d, want MarkerDel", result.InlineMarkers[0].Type)
	}
}

// --- Theme and lang tokens ---

func TestParse_ThemeOverride(t *testing.T) {
	result := Parse(`go theme="dracula"`)
	if result.BlockOptions.Theme != "dracula" {
		t.Errorf("Theme = %q, want %q", result.BlockOptions.Theme, "dracula")
	}
}

func TestParse_ThemeDual(t *testing.T) {
	result := Parse(`go theme="dracula,nord"`)
	if result.BlockOptions.Theme != "dracula,nord" {
		t.Errorf("Theme = %q, want %q", result.BlockOptions.Theme, "dracula,nord")
	}
}

func TestParse_DiffLang(t *testing.T) {
	result := Parse(`diff lang="javascript"`)
	if result.BlockOptions.Lang != "diff" {
		t.Errorf("Lang = %q, want %q", result.BlockOptions.Lang, "diff")
	}
	if result.DiffLang != "javascript" {
		t.Errorf("DiffLang = %q, want %q", result.DiffLang, "javascript")
	}
}

func TestParse_RegexEscapedSlash(t *testing.T) {
	result := Parse(`go /\/path\//`)
	if len(result.InlineMarkers) != 1 {
		t.Fatalf("expected 1 inline marker, got %d", len(result.InlineMarkers))
	}
	if result.InlineMarkers[0].Text != `/path/` {
		t.Errorf("Text = %q, want %q", result.InlineMarkers[0].Text, `/path/`)
	}
}
