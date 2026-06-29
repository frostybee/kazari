// CSS patterns derived from Expressive Code
// Copyright (c) Hippo (https://github.com/hippotastic)
// MIT License: https://github.com/expressive-code/expressive-code/blob/main/LICENSE

package css

import (
	"embed"
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/minify"
	"github.com/frostybee/kazari/internal/theme"
)

//go:embed static/*.css
var staticFS embed.FS

func readCSS(name string) string {
	data, err := staticFS.ReadFile("static/" + name)
	if err != nil {
		panic("kazari: missing embedded CSS: " + name)
	}
	return string(data)
}

func writeThemeCSS(sb *strings.Builder, cfg *config.Config, light, dark theme.ThemeColors) {
	sb.WriteString(theme.GenerateVars(cfg, light, dark))
	sb.WriteString(theme.TokenSwitchingCSS(cfg))
	if cfg.ThemeToggle && cfg.DarkTheme != "" {
		sb.WriteString(theme.ThemeToggleCSS(cfg, light, dark))
	}
}

// Generate produces the complete CSS output for the engine configuration.
func Generate(cfg *config.Config, light, dark theme.ThemeColors) string {
	var sb strings.Builder

	writeThemeCSS(&sb, cfg, light, dark)

	// Static: always included
	sb.WriteString(readCSS("base.css"))
	sb.WriteString(readCSS("line-numbers.css"))
	sb.WriteString(readCSS("markers.css"))
	sb.WriteString(readCSS("inline-markers.css"))
	sb.WriteString(readCSS("focus.css"))

	// Static: conditional on config
	if cfg.StyleReset {
		sb.WriteString(readCSS("style-reset.css"))
	}
	if cfg.ThemedScrollbars {
		sb.WriteString(readCSS("scrollbar.css"))
	}
	if cfg.ThemedSelection {
		sb.WriteString(readCSS("selection.css"))
	}

	sb.WriteString(readCSS("frame.css"))
	sb.WriteString(readCSS("toolbar.css"))
	sb.WriteString(readCSS("tooltip.css"))
	if cfg.FileIcons {
		sb.WriteString(readCSS("file-icons.css"))
	}
	if cfg.LangIconMode != config.LangIconNone {
		sb.WriteString(readCSS("lang-icons.css"))
	}
	sb.WriteString(readCSS("terminal.css"))

	sb.WriteString(readCSS("terminal-dots-colored.css"))
	sb.WriteString(readCSS("terminal-dots-minimal.css"))

	if cfg.Links {
		sb.WriteString(readCSS("links.css"))
	}
	if cfg.CopyButton {
		sb.WriteString(readCSS("copy.css"))
	}
	if cfg.FullscreenButton {
		sb.WriteString(readCSS("fullscreen.css"))
	}
	if cfg.WrapButton {
		sb.WriteString(readCSS("wrap.css"))
	}
	if cfg.Collapsible != nil {
		sb.WriteString(readCSS("collapsible.css"))
	}
	if cfg.CodeGroups {
		sb.WriteString(readCSS("codegroup.css"))
	}
	if cfg.ThemeToggle {
		sb.WriteString(readCSS("theme-toggle.css"))
	}
	if cfg.OutputPanel {
		sb.WriteString(readCSS("output.css"))
	}

	return finalizeCSS(sb.String(), cfg)
}

// GenerateThemeOnly produces only theme variables and token switching CSS,
// without structural rules. Use this for secondary engines on multi-engine
// pages where one primary engine provides the full CSS via Generate().
func GenerateThemeOnly(cfg *config.Config, light, dark theme.ThemeColors) string {
	var sb strings.Builder
	writeThemeCSS(&sb, cfg, light, dark)
	if cfg.ThemeToggle {
		sb.WriteString(readCSS("theme-toggle.css"))
	}

	return finalizeCSS(sb.String(), cfg)
}

func finalizeCSS(content string, cfg *config.Config) string {
	if cfg.CascadeLayer != "" {
		content = fmt.Sprintf("@layer %s {\n%s}\n", cfg.CascadeLayer, content)
	}
	if cfg.Minify {
		return minify.CSS(content)
	}
	return content
}
