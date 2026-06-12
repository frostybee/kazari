package js

import (
	"strings"
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func defaultCfg() *config.Config {
	cfg := config.DefaultConfig()
	return cfg
}

func TestGenerate_CopyButtonAlwaysIncluded(t *testing.T) {
	cfg := defaultCfg()
	out := Generate(cfg)
	if !strings.Contains(out, "clipboard") {
		t.Error("expected copy JS to be included by default")
	}
}

func TestGenerate_FullscreenIncluded(t *testing.T) {
	cfg := defaultCfg()
	out := Generate(cfg)
	if !strings.Contains(out, "fullscreen") {
		t.Error("expected fullscreen JS to be included by default")
	}
}

func TestGenerate_CollapsibleConditional(t *testing.T) {
	cfg := defaultCfg()
	out := Generate(cfg)
	if strings.Contains(out, "kz-collapsed") {
		t.Error("collapsible JS should not be included when Collapsible is nil")
	}

	cfg.Collapsible = &config.CollapsibleConfig{LineThreshold: 10}
	out = Generate(cfg)
	if !strings.Contains(out, "kz-collapsed") {
		t.Error("collapsible JS should be included when Collapsible is set")
	}
}

func TestGenerate_CodeGroupsConditional(t *testing.T) {
	cfg := defaultCfg()
	out := Generate(cfg)
	if strings.Contains(out, "kz-group-tabs") {
		t.Error("code group JS should not be included when CodeGroups is false")
	}

	cfg.CodeGroups = true
	out = Generate(cfg)
	if !strings.Contains(out, "kz-group-tabs") {
		t.Error("code group JS should be included when CodeGroups is true")
	}
}

func TestGenerate_MinifyReducesSize(t *testing.T) {
	cfg := defaultCfg()
	cfg.Minify = false
	unminified := Generate(cfg)

	cfg.Minify = true
	minified := Generate(cfg)

	if len(minified) >= len(unminified) {
		t.Errorf("minified (%d bytes) should be shorter than unminified (%d bytes)", len(minified), len(unminified))
	}
}

func TestGenerate_NothingEnabled(t *testing.T) {
	cfg := defaultCfg()
	cfg.CopyButton = false
	cfg.FullscreenButton = false
	cfg.WrapButton = false
	cfg.Collapsible = nil
	cfg.CodeGroups = false

	out := Generate(cfg)
	if out != "" {
		t.Errorf("expected empty output when no JS features enabled, got %d bytes", len(out))
	}
}
