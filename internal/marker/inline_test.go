package marker

import (
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func tok(content string) config.MergedToken {
	return config.MergedToken{Content: content, LightColor: "#000", DarkColor: "#fff"}
}

func TestProcessInlineMarkers_NoMarkers(t *testing.T) {
	tokens := []config.MergedToken{tok("hello")}
	result := ProcessInlineMarkers(tokens, nil)
	if len(result) != 1 || result[0].Segments[0].Content != "hello" {
		t.Error("expected unchanged token")
	}
	if result[0].Segments[0].Marker != nil {
		t.Error("should have no marker")
	}
}

func TestProcessInlineMarkers_NoMatch(t *testing.T) {
	tokens := []config.MergedToken{tok("hello world")}
	markers := []config.InlineMarker{{Type: config.MarkerMark, Text: "xyz"}}
	result := ProcessInlineMarkers(tokens, markers)
	if len(result) != 1 || len(result[0].Segments) != 1 {
		t.Error("expected single unchanged segment")
	}
}

func TestProcessInlineMarkers_FullTokenMatch(t *testing.T) {
	tokens := []config.MergedToken{tok("foo"), tok(" "), tok("bar")}
	markers := []config.InlineMarker{{Type: config.MarkerMark, Text: "foo"}}
	result := ProcessInlineMarkers(tokens, markers)

	if len(result[0].Segments) != 1 {
		t.Fatalf("expected 1 segment for 'foo', got %d", len(result[0].Segments))
	}
	seg := result[0].Segments[0]
	if seg.Content != "foo" || seg.Marker == nil || seg.Marker.Type != config.MarkerMark {
		t.Error("expected marked 'foo'")
	}
	if seg.Marker.OpenStart || seg.Marker.OpenEnd {
		t.Error("standalone match should have no open flags")
	}
}

func TestProcessInlineMarkers_PartialToken(t *testing.T) {
	tokens := []config.MergedToken{tok("useState")}
	markers := []config.InlineMarker{{Type: config.MarkerMark, Text: "use"}}
	result := ProcessInlineMarkers(tokens, markers)

	segs := result[0].Segments
	if len(segs) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segs))
	}
	if segs[0].Content != "use" || segs[0].Marker == nil {
		t.Error("first segment should be marked 'use'")
	}
	if segs[1].Content != "State" || segs[1].Marker != nil {
		t.Error("second segment should be unmarked 'State'")
	}
}

func TestProcessInlineMarkers_MultiTokenSpan(t *testing.T) {
	// "foo.bar" spans 3 tokens: "foo", ".", "bar"
	tokens := []config.MergedToken{tok("foo"), tok("."), tok("bar")}
	markers := []config.InlineMarker{{Type: config.MarkerMark, Text: "foo.bar"}}
	result := ProcessInlineMarkers(tokens, markers)

	// First token: "foo" with OpenEnd (match continues)
	s0 := result[0].Segments[0]
	if s0.Content != "foo" || s0.Marker == nil {
		t.Error("first token should be marked")
	}
	if s0.Marker.OpenStart || !s0.Marker.OpenEnd {
		t.Errorf("first token: expected OpenStart=false OpenEnd=true, got %v %v", s0.Marker.OpenStart, s0.Marker.OpenEnd)
	}

	// Middle token: "." with both open flags
	s1 := result[1].Segments[0]
	if s1.Content != "." || s1.Marker == nil {
		t.Error("middle token should be marked")
	}
	if !s1.Marker.OpenStart || !s1.Marker.OpenEnd {
		t.Errorf("middle token: expected both open flags, got %v %v", s1.Marker.OpenStart, s1.Marker.OpenEnd)
	}

	// Last token: "bar" with OpenStart (match started before)
	s2 := result[2].Segments[0]
	if s2.Content != "bar" || s2.Marker == nil {
		t.Error("last token should be marked")
	}
	if !s2.Marker.OpenStart || s2.Marker.OpenEnd {
		t.Errorf("last token: expected OpenStart=true OpenEnd=false, got %v %v", s2.Marker.OpenStart, s2.Marker.OpenEnd)
	}
}

func TestProcessInlineMarkers_MultipleOccurrences(t *testing.T) {
	tokens := []config.MergedToken{tok("a b a b")}
	markers := []config.InlineMarker{{Type: config.MarkerMark, Text: "a"}}
	result := ProcessInlineMarkers(tokens, markers)

	segs := result[0].Segments
	// Expected: "a" (marked), " b " (plain), "a" (marked), " b" (plain)
	if len(segs) != 4 {
		t.Fatalf("expected 4 segments, got %d: %+v", len(segs), segs)
	}
	if segs[0].Content != "a" || segs[0].Marker == nil {
		t.Error("segment 0 should be marked 'a'")
	}
	if segs[1].Content != " b " || segs[1].Marker != nil {
		t.Error("segment 1 should be plain ' b '")
	}
	if segs[2].Content != "a" || segs[2].Marker == nil {
		t.Error("segment 2 should be marked 'a'")
	}
	if segs[3].Content != " b" || segs[3].Marker != nil {
		t.Error("segment 3 should be plain ' b'")
	}
}

func TestProcessInlineMarkers_OverlappingPriority(t *testing.T) {
	tokens := []config.MergedToken{tok("abcdef")}
	markers := []config.InlineMarker{
		{Type: config.MarkerMark, Text: "bcde"},  // lower priority
		{Type: config.MarkerIns, Text: "cd"},      // higher priority, overlaps
	}
	result := ProcessInlineMarkers(tokens, markers)

	// Expected segments: "a" (plain), "b" (mark), "cd" (ins), "e" (mark), "f" (plain)
	segs := result[0].Segments
	if len(segs) != 5 {
		t.Fatalf("expected 5 segments, got %d: %+v", len(segs), segs)
	}
	if segs[0].Content != "a" || segs[0].Marker != nil {
		t.Errorf("seg 0: expected plain 'a', got %q marker=%v", segs[0].Content, segs[0].Marker)
	}
	if segs[1].Content != "b" || segs[1].Marker == nil || segs[1].Marker.Type != config.MarkerMark {
		t.Errorf("seg 1: expected mark 'b', got %q marker=%v", segs[1].Content, segs[1].Marker)
	}
	if segs[2].Content != "cd" || segs[2].Marker == nil || segs[2].Marker.Type != config.MarkerIns {
		t.Errorf("seg 2: expected ins 'cd', got %q marker=%v", segs[2].Content, segs[2].Marker)
	}
	if segs[3].Content != "e" || segs[3].Marker == nil || segs[3].Marker.Type != config.MarkerMark {
		t.Errorf("seg 3: expected mark 'e', got %q marker=%v", segs[3].Content, segs[3].Marker)
	}
	if segs[4].Content != "f" || segs[4].Marker != nil {
		t.Errorf("seg 4: expected plain 'f', got %q marker=%v", segs[4].Content, segs[4].Marker)
	}
}

func TestProcessInlineMarkers_InsType(t *testing.T) {
	tokens := []config.MergedToken{tok("hello world")}
	markers := []config.InlineMarker{{Type: config.MarkerIns, Text: "world"}}
	result := ProcessInlineMarkers(tokens, markers)

	segs := result[0].Segments
	if len(segs) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segs))
	}
	if segs[1].Marker == nil || segs[1].Marker.Type != config.MarkerIns {
		t.Error("expected ins marker on 'world'")
	}
}

func TestProcessInlineMarkers_DelType(t *testing.T) {
	tokens := []config.MergedToken{tok("hello world")}
	markers := []config.InlineMarker{{Type: config.MarkerDel, Text: "hello"}}
	result := ProcessInlineMarkers(tokens, markers)

	segs := result[0].Segments
	if segs[0].Marker == nil || segs[0].Marker.Type != config.MarkerDel {
		t.Error("expected del marker on 'hello'")
	}
}
