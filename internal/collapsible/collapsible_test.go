package collapsible

import (
	"reflect"
	"testing"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/locale"
)

func intPtr(n int) *int { return &n }

func TestShouldThresholdCollapse(t *testing.T) {
	cfg := &config.CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}

	tests := []struct {
		name      string
		lineCount int
		spec      *config.CollapseSpec
		cfg       *config.CollapsibleConfig
		want      bool
	}{
		{"nil config", 50, nil, nil, false},
		{"below threshold", 10, nil, cfg, false},
		{"at threshold", 15, nil, cfg, false},
		{"above threshold", 16, nil, cfg, true},
		{"force collapse below threshold", 5, &config.CollapseSpec{Enabled: true}, cfg, true},
		{"nocollapse above threshold", 50, &config.CollapseSpec{Disabled: true}, cfg, false},
		{"nocollapse overrides enabled", 50, &config.CollapseSpec{Enabled: true, Disabled: true}, cfg, false},
		{"per-block threshold below", 18, &config.CollapseSpec{Threshold: intPtr(20)}, cfg, false},
		{"per-block threshold above", 21, &config.CollapseSpec{Threshold: intPtr(20)}, cfg, true},
		{"per-block threshold overrides engine", 16, &config.CollapseSpec{Threshold: intPtr(20)}, cfg, false},
		{"per-block threshold ignored with nocollapse", 21, &config.CollapseSpec{Disabled: true, Threshold: intPtr(20)}, cfg, false},
		{"per-block threshold ignored with collapse", 5, &config.CollapseSpec{Enabled: true, Threshold: intPtr(20)}, cfg, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldThresholdCollapse(tt.lineCount, tt.spec, tt.cfg)
			if got != tt.want {
				t.Errorf("ShouldThresholdCollapse(%d) = %v, want %v", tt.lineCount, got, tt.want)
			}
		})
	}
}

func TestValidateRanges(t *testing.T) {
	tests := []struct {
		name      string
		ranges    []config.LineRange
		lineCount int
		want      []config.LineRange
	}{
		{
			"valid ranges",
			[]config.LineRange{{Start: 2, End: 5}, {Start: 10, End: 15}},
			20,
			[]config.LineRange{{Start: 2, End: 5}, {Start: 10, End: 15}},
		},
		{
			"reversed range dropped",
			[]config.LineRange{{Start: 8, End: 2}},
			20,
			nil,
		},
		{
			"out of bounds dropped",
			[]config.LineRange{{Start: 25, End: 30}},
			20,
			nil,
		},
		{
			"clamped to line count",
			[]config.LineRange{{Start: 18, End: 25}},
			20,
			[]config.LineRange{{Start: 18, End: 20}},
		},
		{
			"overlapping second dropped",
			[]config.LineRange{{Start: 2, End: 8}, {Start: 5, End: 12}},
			20,
			[]config.LineRange{{Start: 2, End: 8}},
		},
		{
			"unsorted input sorted",
			[]config.LineRange{{Start: 10, End: 15}, {Start: 2, End: 5}},
			20,
			[]config.LineRange{{Start: 2, End: 5}, {Start: 10, End: 15}},
		},
		{
			"single line range valid",
			[]config.LineRange{{Start: 5, End: 5}},
			20,
			[]config.LineRange{{Start: 5, End: 5}},
		},
		{
			"adjacent ranges valid",
			[]config.LineRange{{Start: 2, End: 5}, {Start: 6, End: 10}},
			20,
			[]config.LineRange{{Start: 2, End: 5}, {Start: 6, End: 10}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateRanges(tt.ranges, tt.lineCount)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateRanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeMinIndent(t *testing.T) {
	code := "line1\n    line2\n        line3\n\n    line4"

	tests := []struct {
		name      string
		start     int
		end       int
		want      int
	}{
		{"first line no indent", 1, 1, 0},
		{"indented lines", 2, 4, 4},
		{"deeply indented", 3, 3, 8},
		{"blank line ignored", 2, 5, 4},
		{"full range", 1, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeMinIndent(code, tt.start, tt.end)
			if got != tt.want {
				t.Errorf("ComputeMinIndent(%d, %d) = %d, want %d", tt.start, tt.end, got, tt.want)
			}
		})
	}
}

func TestSummaryText(t *testing.T) {
	tests := []struct {
		count int
		want  string
	}{
		{1, "1 collapsed line"},
		{5, "5 collapsed lines"},
		{0, "0 collapsed lines"},
	}

	s := locale.Resolve("en-US", nil)
	for _, tt := range tests {
		got := SummaryText(tt.count, s)
		if got != tt.want {
			t.Errorf("SummaryText(%d) = %q, want %q", tt.count, got, tt.want)
		}
	}
}

func TestComputePreviewSegments_FurthestVisibleLine(t *testing.T) {
	tests := []struct {
		name         string
		previewLines int
		lineCount    int
		markers      []config.LineMarker
		focusLines   []config.LineRange
		want         int
	}{
		{"no markers", 8, 50, nil, nil, 8},
		{"preview exceeds line count", 8, 5, nil, nil, 5},
		{
			"extends to marked line with context",
			8, 50,
			[]config.LineMarker{{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 10, End: 10}}}},
			nil,
			11, // line 10 + 1 context
		},
		{
			"capped at 2x",
			8, 50,
			[]config.LineMarker{{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 20, End: 20}}}},
			nil,
			8,
		},
		{
			"extends to focus line with context",
			8, 50,
			nil,
			[]config.LineRange{{Start: 12, End: 12}},
			13, // line 12 + 1 context
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, _ := ComputePreviewSegments(tt.previewLines, tt.lineCount, tt.markers, tt.focusLines)
			if len(segments) == 0 {
				t.Fatal("expected at least one preview segment")
			}
			got := segments[len(segments)-1].End
			if got != tt.want {
				t.Errorf("furthest visible line = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestResolveCollapse(t *testing.T) {
	cfg := &config.CollapsibleConfig{
		LineThreshold:    15,
		PreviewLines:     8,
		DefaultCollapsed: true,
		PreserveIndent:   true,
	}

	t.Run("threshold only", func(t *testing.T) {
		code := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\nline11\nline12\nline13\nline14\nline15\nline16\nline17\nline18\nline19\nline20"
		result := ResolveCollapse(20, nil, cfg, code, nil, nil)
		if !result.Threshold {
			t.Error("expected threshold collapse")
		}
		if len(result.PreviewSegments) != 1 || result.PreviewSegments[0].End != 8 {
			t.Errorf("segments = %v, want [{1 8}]", result.PreviewSegments)
		}
		if len(result.Ranges) != 0 {
			t.Error("expected no ranges")
		}
	})

	t.Run("range only", func(t *testing.T) {
		code := "a\n    b\n    c\n    d\ne"
		spec := &config.CollapseSpec{Ranges: []config.LineRange{{Start: 2, End: 4}}}
		result := ResolveCollapse(5, spec, cfg, code, nil, nil)
		if result.Threshold {
			t.Error("unexpected threshold collapse")
		}
		if len(result.Ranges) != 1 {
			t.Fatalf("got %d ranges, want 1", len(result.Ranges))
		}
		r := result.Ranges[0]
		if r.Start != 2 || r.End != 4 || r.LineCount != 3 || r.MinIndent != 4 {
			t.Errorf("range = %+v, want {2 4 3 4}", r)
		}
	})

	t.Run("both modes", func(t *testing.T) {
		lines := make([]byte, 0, 200)
		for i := 0; i < 20; i++ {
			if i > 0 {
				lines = append(lines, '\n')
			}
			lines = append(lines, "    line"...)
		}
		spec := &config.CollapseSpec{Ranges: []config.LineRange{{Start: 5, End: 8}}}
		result := ResolveCollapse(20, spec, cfg, string(lines), nil, nil)
		if !result.Threshold {
			t.Error("expected threshold collapse")
		}
		if len(result.Ranges) != 1 {
			t.Fatalf("got %d ranges, want 1", len(result.Ranges))
		}
	})
}

func TestComputePreviewSegments_NoMarkers(t *testing.T) {
	segments, beyond := ComputePreviewSegments(8, 50, nil, nil)
	if len(segments) != 1 || segments[0].Start != 1 || segments[0].End != 8 {
		t.Errorf("got %v, want [{1 8}]", segments)
	}
	if beyond != 0 {
		t.Errorf("beyond = %d, want 0", beyond)
	}
}

func TestComputePreviewSegments_MarkerWithinCap(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 12, End: 12}}},
	}
	segments, beyond := ComputePreviewSegments(8, 50, markers, nil)

	if len(segments) != 2 {
		t.Fatalf("got %d segments, want 2: %v", len(segments), segments)
	}
	if segments[0].Start != 1 || segments[0].End != 8 {
		t.Errorf("segment[0] = %v, want {1 8}", segments[0])
	}
	// Marked line 12 with ±1 context = lines 11-13
	if segments[1].Start != 11 || segments[1].End != 13 {
		t.Errorf("segment[1] = %v, want {11 13}", segments[1])
	}
	if beyond != 0 {
		t.Errorf("beyond = %d, want 0", beyond)
	}
}

func TestComputePreviewSegments_MarkerBeyondCap(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 30, End: 32}}},
	}
	segments, beyond := ComputePreviewSegments(8, 50, markers, nil)

	if len(segments) != 1 {
		t.Fatalf("got %d segments, want 1: %v", len(segments), segments)
	}
	if segments[0].Start != 1 || segments[0].End != 8 {
		t.Errorf("segment[0] = %v, want {1 8}", segments[0])
	}
	if beyond != 3 {
		t.Errorf("beyond = %d, want 3", beyond)
	}
}

func TestComputePreviewSegments_AdjacentMerge(t *testing.T) {
	// Lines 10,12 with ±1 context: 9-11 and 11-13.
	// 9 is adjacent to base {1,8} → all merge into {1,13}
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 10, End: 10}, {Start: 12, End: 12}}},
	}
	segments, _ := ComputePreviewSegments(8, 50, markers, nil)

	if len(segments) != 1 {
		t.Fatalf("got %d segments, want 1: %v", len(segments), segments)
	}
	if segments[0].Start != 1 || segments[0].End != 13 {
		t.Errorf("segment[0] = %v, want {1 13}", segments[0])
	}
}

func TestComputePreviewSegments_DisjointSegments(t *testing.T) {
	// Marker at line 14 — context 13-15, gap from base {1,8}
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 14, End: 14}}},
	}
	segments, _ := ComputePreviewSegments(8, 50, markers, nil)

	if len(segments) != 2 {
		t.Fatalf("got %d segments, want 2: %v", len(segments), segments)
	}
	if segments[0].Start != 1 || segments[0].End != 8 {
		t.Errorf("segment[0] = %v, want {1 8}", segments[0])
	}
	if segments[1].Start != 13 || segments[1].End != 15 {
		t.Errorf("segment[1] = %v, want {13 15}", segments[1])
	}
}

func TestComputePreviewSegments_MarkerInBasePreview(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 3, End: 5}}},
	}
	segments, beyond := ComputePreviewSegments(8, 50, markers, nil)

	// All markers within base preview — single segment, no extras
	if len(segments) != 1 {
		t.Fatalf("got %d segments, want 1: %v", len(segments), segments)
	}
	if segments[0].End != 8 {
		t.Errorf("segment[0].End = %d, want 8", segments[0].End)
	}
	if beyond != 0 {
		t.Errorf("beyond = %d, want 0", beyond)
	}
}

func TestComputePreviewSegments_FocusLinesWithinCap(t *testing.T) {
	focus := []config.LineRange{{Start: 14, End: 14}}
	segments, _ := ComputePreviewSegments(8, 50, nil, focus)

	if len(segments) != 2 {
		t.Fatalf("got %d segments, want 2: %v", len(segments), segments)
	}
	if segments[1].Start != 13 || segments[1].End != 15 {
		t.Errorf("segment[1] = %v, want {13 15}", segments[1])
	}
}

func TestComputePreviewSegments_MixedWithinAndBeyond(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{
			{Start: 12, End: 12}, // within 2× cap (16)
			{Start: 25, End: 25}, // beyond 2× cap
		}},
	}
	segments, beyond := ComputePreviewSegments(8, 50, markers, nil)

	if len(segments) != 2 {
		t.Fatalf("got %d segments, want 2: %v", len(segments), segments)
	}
	if beyond != 1 {
		t.Errorf("beyond = %d, want 1", beyond)
	}
}

func TestResolveCollapseStyle(t *testing.T) {
	tests := []struct {
		name      string
		style     config.CollapseStyle
		rangeEnd  int
		lineCount int
		want      config.CollapseStyle
	}{
		{"github passthrough", config.CollapseGithub, 5, 20, config.CollapseGithub},
		{"start passthrough", config.CollapseCollapsibleStart, 5, 20, config.CollapseCollapsibleStart},
		{"end passthrough", config.CollapseCollapsibleEnd, 5, 20, config.CollapseCollapsibleEnd},
		{"auto not at end", config.CollapseCollapsibleAuto, 5, 20, config.CollapseCollapsibleStart},
		{"auto at end", config.CollapseCollapsibleAuto, 20, 20, config.CollapseCollapsibleEnd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveCollapseStyle(tt.style, tt.rangeEnd, tt.lineCount)
			if got != tt.want {
				t.Errorf("ResolveCollapseStyle() = %d, want %d", got, tt.want)
			}
		})
	}
}
