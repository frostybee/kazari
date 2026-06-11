package svgutil

import (
	"strings"
	"testing"
)

func TestInlineSVGURL_Encoding(t *testing.T) {
	got := InlineSVGURL(`<svg viewBox="0 0 1 1">{x}</svg>`)
	want := "data:image/svg+xml,%3Csvg viewBox=%220 0 1 1%22%3E%7Bx%7D%3C/svg%3E"
	if got != want {
		t.Errorf("InlineSVGURL():\ngot:  %s\nwant: %s", got, want)
	}
}

func TestInlineSVGURL_PercentEncodedFirst(t *testing.T) {
	got := InlineSVGURL("<svg width='100%'/>")
	if !strings.Contains(got, "100%25") {
		t.Errorf("percent sign should encode to %%25 without double-encoding, got %s", got)
	}
	if strings.Contains(got, "%2525") {
		t.Errorf("percent sign was double-encoded: %s", got)
	}
}

func TestInlineSVGURL_SingleQuotesPreserved(t *testing.T) {
	got := InlineSVGURL("<svg xmlns='http://www.w3.org/2000/svg'/>")
	if !strings.Contains(got, "xmlns='http://www.w3.org/2000/svg'") {
		t.Errorf("single quotes should pass through unchanged, got %s", got)
	}
}

func TestInlineSVGURL_MatchesTerminalIcon(t *testing.T) {
	svg := "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 60 16'>" +
		"<circle cx='8' cy='8' r='8'/><circle cx='30' cy='8' r='8'/><circle cx='52' cy='8' r='8'/></svg>"
	// The exact URL previously hardcoded in theme.go for the minimal dots icon.
	want := "data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 60 16'%3E%3Ccircle cx='8' cy='8' r='8'/%3E%3Ccircle cx='30' cy='8' r='8'/%3E%3Ccircle cx='52' cy='8' r='8'/%3E%3C/svg%3E"
	if got := InlineSVGURL(svg); got != want {
		t.Errorf("terminal icon URL changed:\ngot:  %s\nwant: %s", got, want)
	}
}
