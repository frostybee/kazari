package diff

import (
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func TestProcessDiffBlock_Mixed(t *testing.T) {
	code := " const x = 1;\n+ const y = 2;\n- const z = 3;"
	stripped, markers := ProcessDiffBlock(code)

	if stripped != "const x = 1;\nconst y = 2;\nconst z = 3;" {
		t.Errorf("stripped = %q", stripped)
	}

	insCount, delCount := 0, 0
	for _, m := range markers {
		if m.Type == config.MarkerIns {
			insCount += len(m.Lines)
		}
		if m.Type == config.MarkerDel {
			delCount += len(m.Lines)
		}
	}
	if insCount != 1 {
		t.Errorf("expected 1 ins line, got %d", insCount)
	}
	if delCount != 1 {
		t.Errorf("expected 1 del line, got %d", delCount)
	}
}

func TestProcessDiffBlock_AllAdditions(t *testing.T) {
	code := "+ line1\n+ line2\n+ line3"
	stripped, markers := ProcessDiffBlock(code)

	if stripped != "line1\nline2\nline3" {
		t.Errorf("stripped = %q", stripped)
	}

	for _, m := range markers {
		if m.Type == config.MarkerIns && len(m.Lines) != 3 {
			t.Errorf("expected 3 ins lines, got %d", len(m.Lines))
		}
	}
}

func TestProcessDiffBlock_AllDeletions(t *testing.T) {
	code := "- line1\n- line2"
	stripped, markers := ProcessDiffBlock(code)

	if stripped != "line1\nline2" {
		t.Errorf("stripped = %q", stripped)
	}

	for _, m := range markers {
		if m.Type == config.MarkerDel && len(m.Lines) != 2 {
			t.Errorf("expected 2 del lines, got %d", len(m.Lines))
		}
	}
}

func TestProcessDiffBlock_NoPrefix(t *testing.T) {
	code := "plain line\nanother"
	stripped, markers := ProcessDiffBlock(code)

	if stripped != code {
		t.Errorf("lines without prefix should be unchanged, got %q", stripped)
	}
	if len(markers) != 0 {
		t.Errorf("expected no markers, got %d", len(markers))
	}
}

func TestProcessDiffBlock_Empty(t *testing.T) {
	stripped, markers := ProcessDiffBlock("")

	if stripped != "" {
		t.Errorf("empty input should return empty, got %q", stripped)
	}
	if len(markers) != 0 {
		t.Errorf("expected no markers, got %d", len(markers))
	}
}

func TestProcessDiffBlock_ContextLines(t *testing.T) {
	code := " unchanged\n+ added\n also unchanged"
	stripped, markers := ProcessDiffBlock(code)

	if stripped != "unchanged\nadded\nalso unchanged" {
		t.Errorf("stripped = %q", stripped)
	}

	total := 0
	for _, m := range markers {
		total += len(m.Lines)
	}
	if total != 1 {
		t.Errorf("context lines should have no markers, total marked = %d", total)
	}
}

func TestProcessDiffBlock_LineNumbers(t *testing.T) {
	code := " ctx\n+ add\n- del"
	_, markers := ProcessDiffBlock(code)

	for _, m := range markers {
		if m.Type == config.MarkerIns {
			if m.Lines[0].Start != 2 {
				t.Errorf("ins should be line 2, got %d", m.Lines[0].Start)
			}
		}
		if m.Type == config.MarkerDel {
			if m.Lines[0].Start != 3 {
				t.Errorf("del should be line 3, got %d", m.Lines[0].Start)
			}
		}
	}
}
