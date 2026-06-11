package locale

import "testing"

func TestResolve_EnUS_Defaults(t *testing.T) {
	s := Resolve("en-US", nil)
	if s.CopyLabel != "Copy" {
		t.Errorf("CopyLabel = %q, want %q", s.CopyLabel, "Copy")
	}
	if s.FullscreenLabel != "Fullscreen" {
		t.Errorf("FullscreenLabel = %q, want %q", s.FullscreenLabel, "Fullscreen")
	}
	if s.ExpandButtonText != "Show more" {
		t.Errorf("ExpandButtonText = %q, want %q", s.ExpandButtonText, "Show more")
	}
}

func TestResolve_FrFR(t *testing.T) {
	s := Resolve("fr-FR", nil)
	if s.CopyLabel != "Copier" {
		t.Errorf("CopyLabel = %q, want %q", s.CopyLabel, "Copier")
	}
	if s.CopySuccess != "Copié !" {
		t.Errorf("CopySuccess = %q, want %q", s.CopySuccess, "Copié !")
	}
}

func TestResolve_JaJP(t *testing.T) {
	s := Resolve("ja-JP", nil)
	if s.CopyLabel != "コピー" {
		t.Errorf("CopyLabel = %q, want %q", s.CopyLabel, "コピー")
	}
	if s.FullscreenLabel != "全画面" {
		t.Errorf("FullscreenLabel = %q, want %q", s.FullscreenLabel, "全画面")
	}
}

func TestResolve_UnknownLocaleFallsBack(t *testing.T) {
	s := Resolve("de-DE", nil)
	if s.CopyLabel != "Copy" {
		t.Errorf("unknown locale should fall back to en-US, got CopyLabel=%q", s.CopyLabel)
	}
}

func TestResolve_OverrideSingleKey(t *testing.T) {
	s := Resolve("en-US", map[string]string{"copy.label": "Copy code"})
	if s.CopyLabel != "Copy code" {
		t.Errorf("CopyLabel = %q, want %q", s.CopyLabel, "Copy code")
	}
	if s.CopySuccess != "Copied!" {
		t.Error("non-overridden keys should keep defaults")
	}
}

func TestResolve_OverrideMultipleKeys(t *testing.T) {
	s := Resolve("en-US", map[string]string{
		"copy.label":   "Copy code",
		"copy.success": "Done!",
	})
	if s.CopyLabel != "Copy code" {
		t.Errorf("CopyLabel = %q", s.CopyLabel)
	}
	if s.CopySuccess != "Done!" {
		t.Errorf("CopySuccess = %q", s.CopySuccess)
	}
}

func TestResolve_EmptyOverrideNoOp(t *testing.T) {
	s := Resolve("en-US", map[string]string{})
	if s.CopyLabel != "Copy" {
		t.Error("empty overrides should not change defaults")
	}
}

func TestResolve_OverrideOnLocale(t *testing.T) {
	s := Resolve("fr-FR", map[string]string{"copy.label": "Copier le code"})
	if s.CopyLabel != "Copier le code" {
		t.Errorf("override on fr-FR: CopyLabel = %q", s.CopyLabel)
	}
	if s.FullscreenLabel != "Plein écran" {
		t.Error("non-overridden fr-FR keys should keep French defaults")
	}
}

func TestFormatCollapsedLines_Singular(t *testing.T) {
	s := Resolve("en-US", nil)
	result := FormatCollapsedLines(s, 1)
	if result != "1 collapsed line" {
		t.Errorf("singular = %q, want %q", result, "1 collapsed line")
	}
}

func TestFormatCollapsedLines_Plural(t *testing.T) {
	s := Resolve("en-US", nil)
	result := FormatCollapsedLines(s, 5)
	if result != "5 collapsed lines" {
		t.Errorf("plural = %q, want %q", result, "5 collapsed lines")
	}
}

func TestFormatCollapsedLines_French(t *testing.T) {
	s := Resolve("fr-FR", nil)
	result := FormatCollapsedLines(s, 3)
	if result != "3 lignes masquées" {
		t.Errorf("french plural = %q, want %q", result, "3 lignes masquées")
	}
}
