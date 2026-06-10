package marker

import (
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func TestResolveLineMarkers_Empty(t *testing.T) {
	result := ResolveLineMarkers(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestResolveLineMarkers_SingleMark(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 3, End: 5}}},
	}
	result := ResolveLineMarkers(markers)

	for _, line := range []int{3, 4, 5} {
		entry, ok := result[line]
		if !ok {
			t.Errorf("line %d should be marked", line)
			continue
		}
		if entry.Type != config.MarkerMark {
			t.Errorf("line %d: expected MarkerMark, got %d", line, entry.Type)
		}
		if !entry.HasMark {
			t.Errorf("line %d: HasMark should be true", line)
		}
	}

	if _, ok := result[2]; ok {
		t.Error("line 2 should not be marked")
	}
	if _, ok := result[6]; ok {
		t.Error("line 6 should not be marked")
	}
}

func TestResolveLineMarkers_MultipleTypes(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 1, End: 3}}},
		{Type: config.MarkerIns, Lines: []config.LineRange{{Start: 5, End: 7}}},
		{Type: config.MarkerDel, Lines: []config.LineRange{{Start: 9, End: 9}}},
	}
	result := ResolveLineMarkers(markers)

	if result[1].Type != config.MarkerMark {
		t.Error("line 1 should be mark")
	}
	if result[5].Type != config.MarkerIns {
		t.Error("line 5 should be ins")
	}
	if result[9].Type != config.MarkerDel {
		t.Error("line 9 should be del")
	}
}

func TestResolveLineMarkers_OverlapHigherPriorityWins(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 10, End: 20}}},
		{Type: config.MarkerDel, Lines: []config.LineRange{{Start: 12, End: 15}}},
	}
	result := ResolveLineMarkers(markers)

	if result[10].Type != config.MarkerMark {
		t.Error("line 10 should be mark (before del range)")
	}
	if result[12].Type != config.MarkerDel {
		t.Error("line 12 should be del (higher priority)")
	}
	if result[15].Type != config.MarkerDel {
		t.Error("line 15 should be del (higher priority)")
	}
	if result[16].Type != config.MarkerMark {
		t.Error("line 16 should be mark (after del range)")
	}
}

func TestResolveLineMarkers_InsOverridesDel(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerDel, Lines: []config.LineRange{{Start: 5, End: 10}}},
		{Type: config.MarkerIns, Lines: []config.LineRange{{Start: 7, End: 8}}},
	}
	result := ResolveLineMarkers(markers)

	if result[5].Type != config.MarkerDel {
		t.Error("line 5 should be del")
	}
	if result[7].Type != config.MarkerIns {
		t.Error("line 7 should be ins (highest priority)")
	}
	if result[8].Type != config.MarkerIns {
		t.Error("line 8 should be ins (highest priority)")
	}
	if result[9].Type != config.MarkerDel {
		t.Error("line 9 should be del")
	}
}

func TestResolveLineMarkers_SameTypeMerge(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 1, End: 5}}},
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 3, End: 8}}},
	}
	result := ResolveLineMarkers(markers)

	for _, line := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		entry, ok := result[line]
		if !ok {
			t.Errorf("line %d should be marked", line)
			continue
		}
		if entry.Type != config.MarkerMark {
			t.Errorf("line %d: expected MarkerMark", line)
		}
	}
}

func TestResolveLineMarkers_LabeledRange(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerIns, Lines: []config.LineRange{{Start: 6, End: 10}}, Label: "A"},
	}
	result := ResolveLineMarkers(markers)

	if result[6].Label != "A" {
		t.Errorf("line 6 (first) should have label 'A', got %q", result[6].Label)
	}
	if result[7].Label != "" {
		t.Errorf("line 7 should have no label, got %q", result[7].Label)
	}
	if result[10].Label != "" {
		t.Errorf("line 10 should have no label, got %q", result[10].Label)
	}
}

func TestResolveLineMarkers_LabelMultipleRanges(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerIns, Lines: []config.LineRange{
			{Start: 2, End: 4},
			{Start: 8, End: 10},
		}, Label: "A"},
	}
	result := ResolveLineMarkers(markers)

	if result[2].Label != "A" {
		t.Errorf("line 2 (first of range 1) should have label")
	}
	if result[3].Label != "" {
		t.Errorf("line 3 should have no label")
	}
	if result[8].Label != "A" {
		t.Errorf("line 8 (first of range 2) should have label")
	}
	if result[9].Label != "" {
		t.Errorf("line 9 should have no label")
	}
}

func TestResolveLineMarkers_LabelOverriddenByHigherPriority(t *testing.T) {
	markers := []config.LineMarker{
		{Type: config.MarkerMark, Lines: []config.LineRange{{Start: 5, End: 10}}, Label: "X"},
		{Type: config.MarkerIns, Lines: []config.LineRange{{Start: 5, End: 5}}},
	}
	result := ResolveLineMarkers(markers)

	if result[5].Type != config.MarkerIns {
		t.Error("line 5 should be ins (higher priority)")
	}
	if result[5].Label != "" {
		t.Errorf("line 5 label should be cleared when overridden, got %q", result[5].Label)
	}
	if result[6].Type != config.MarkerMark {
		t.Error("line 6 should be mark")
	}
}

func TestResolveFocusSet_Empty(t *testing.T) {
	result := ResolveFocusSet(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestResolveFocusSet(t *testing.T) {
	focus := []config.LineRange{
		{Start: 3, End: 5},
		{Start: 10, End: 10},
	}
	result := ResolveFocusSet(focus)

	for _, line := range []int{3, 4, 5, 10} {
		if !result[line] {
			t.Errorf("line %d should be focused", line)
		}
	}
	for _, line := range []int{1, 2, 6, 9, 11} {
		if result[line] {
			t.Errorf("line %d should not be focused", line)
		}
	}
}
