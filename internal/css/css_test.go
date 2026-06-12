package css

import (
	"strings"
	"testing"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/theme"
)

func defaultCfg() *config.Config {
	return config.DefaultConfig()
}

func colors() (theme.ThemeColors, theme.ThemeColors) {
	light := theme.ThemeColors{EditorBG: "#ffffff", EditorFG: "#24292f"}
	dark := theme.ThemeColors{EditorBG: "#24292f", EditorFG: "#e6edf3"}
	return light, dark
}

func TestGenerate_BaseCSS(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	for _, sel := range []string{".kazari-code", ".kz-line"} {
		if !strings.Contains(out, sel) {
			t.Errorf("expected base selector %q in output", sel)
		}
	}
}

func TestGenerate_LineNumbersCSS(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, ".ln") {
		t.Error("expected line number CSS (.ln) in output")
	}
}

func TestGenerate_MarkersCSS(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, ".kz-line.mark") {
		t.Error("expected markers CSS in output")
	}
	if !strings.Contains(out, ".kz-line.ins") {
		t.Error("expected ins marker CSS in output")
	}
}

func TestGenerate_FocusCSS(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "has-focus") {
		t.Error("expected focus CSS in output")
	}
}

func TestGenerate_CollapsibleConditional(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if strings.Contains(out, "kz-collapse-gradient") {
		t.Error("collapsible CSS should not be included when Collapsible is nil")
	}

	cfg.Collapsible = &config.CollapsibleConfig{LineThreshold: 10}
	out = Generate(cfg, light, dark)
	if !strings.Contains(out, "kz-collapse-gradient") {
		t.Error("collapsible CSS should be included when Collapsible is set")
	}
}

func TestGenerate_CodeGroupsConditional(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if strings.Contains(out, "kz-group-tabs") {
		t.Error("code group CSS should not be included when CodeGroups is false")
	}

	cfg.CodeGroups = true
	out = Generate(cfg, light, dark)
	if !strings.Contains(out, "kz-group-tabs") {
		t.Error("code group CSS should be included when CodeGroups is true")
	}
}

func TestGenerate_FileIconsConditional(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()

	cfg.FileIcons = false
	out := Generate(cfg, light, dark)
	if strings.Contains(out, "kz-file-icon") {
		t.Error("file icon CSS should not be included when FileIcons is false")
	}

	cfg.FileIcons = true
	out = Generate(cfg, light, dark)
	if !strings.Contains(out, "kz-file-icon") {
		t.Error("file icon CSS should be included when FileIcons is true")
	}
}

func TestGenerate_TerminalDotsColored(t *testing.T) {
	cfg := defaultCfg()
	cfg.TerminalDotStyle = config.DotsColored
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "kz-terminal-dots") {
		t.Error("expected terminal dots CSS")
	}
}

func TestGenerate_TerminalDotsMinimal(t *testing.T) {
	cfg := defaultCfg()
	cfg.TerminalDotStyle = config.DotsMinimal
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "kz-dots-minimal") {
		t.Error("expected minimal dots CSS")
	}
}

func TestGenerate_CopyButtonCSS(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "kz-copy-btn") {
		t.Error("expected copy button CSS")
	}
}

func TestGenerate_FullscreenCSS(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "kz-fs-btn") {
		t.Error("expected fullscreen button CSS")
	}
}

func TestGenerate_ThemeVariables(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "--kz-") {
		t.Error("expected --kz- CSS variables in output")
	}
	if !strings.Contains(out, ":root") {
		t.Error("expected :root selector for theme variables")
	}
}

func TestGenerate_CustomCSSRoot(t *testing.T) {
	cfg := defaultCfg()
	cfg.ThemeCSSRoot = ".my-scope"
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, ".my-scope") {
		t.Error("expected custom CSS root .my-scope")
	}
}

func TestGenerate_CascadeLayer(t *testing.T) {
	cfg := defaultCfg()
	cfg.CascadeLayer = "kazari"
	light, dark := colors()
	out := Generate(cfg, light, dark)

	if !strings.Contains(out, "@layer kazari") {
		t.Error("expected @layer kazari wrapper in output")
	}
}

func TestGenerate_MinifyReducesSize(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	cfg.Minify = false
	unminified := Generate(cfg, light, dark)

	cfg.Minify = true
	minified := Generate(cfg, light, dark)

	if len(minified) >= len(unminified) {
		t.Errorf("minified (%d bytes) should be shorter than unminified (%d bytes)", len(minified), len(unminified))
	}
}

func TestGenerate_StyleResetConditional(t *testing.T) {
	light, dark := colors()

	cfgOff := defaultCfg()
	cfgOff.StyleReset = false
	out := Generate(cfgOff, light, dark)
	if strings.Contains(out, "all: revert") {
		t.Error("style reset should not be included when StyleReset is false")
	}

	cfgOn := defaultCfg()
	cfgOn.StyleReset = true
	out = Generate(cfgOn, light, dark)
	if !strings.Contains(out, "all:") {
		t.Errorf("style reset should be included when StyleReset is true, output length: %d", len(out))
	}
}

func TestGenerateThemeOnly_NoStructural(t *testing.T) {
	cfg := defaultCfg()
	light, dark := colors()
	out := GenerateThemeOnly(cfg, light, dark)

	if !strings.Contains(out, "--kz-editor-bg") {
		t.Error("ThemeOnly should contain theme variables")
	}
	if strings.Contains(out, ".kz-toolbar") {
		t.Error("ThemeOnly should not contain structural CSS")
	}
	if strings.Contains(out, ".kz-copy-btn") {
		t.Error("ThemeOnly should not contain structural CSS selectors")
	}
	if strings.Contains(out, "grid-template") {
		t.Error("ThemeOnly should not contain layout CSS")
	}
}
