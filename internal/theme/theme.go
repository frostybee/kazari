package theme

import (
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// ThemeColors holds colors extracted from a syntax theme.
type ThemeColors struct {
	EditorBG     string
	EditorFG     string
	SelectionBG  string
	LineNumberFG string
}

// GenerateVars produces CSS variable declarations for both themes,
// wrapped according to the dark mode strategy.
func GenerateVars(cfg *config.Config, light, dark ThemeColors) string {
	var sb strings.Builder

	// Static defaults (not theme-dependent).
	staticVars := []struct{ name, value string }{
		{"--kz-radius", "0.5rem"},
		{"--kz-shadow", "0 2px 8px rgba(0,0,0,0.15)"},
		{"--kz-border", "1px solid transparent"},
		{"--kz-border-hover", "1px solid var(--kz-accent)"},
		{"--kz-transition", "150ms ease"},
		{"--kz-font-family", "'JetBrains Mono', monospace"},
		{"--kz-font-size", "0.875rem"},
		{"--kz-font-weight", "400"},
		{"--kz-line-height", "1.6"},
		{"--kz-ui-font-family", "system-ui, sans-serif"},
		{"--kz-ui-font-size", "0.9rem"},
		{"--kz-ui-font-weight", "400"},
		{"--kz-ui-line-height", "1.65"},
		{"--kz-code-padding-block", "1rem"},
		{"--kz-code-padding-inline", "1.35rem"},
		{"--kz-title-font-size", "0.8rem"},
		{"--kz-title-padding", "0.5rem 1rem"},
	}

	// Light theme variables.
	lightVars := buildThemeVars(light)
	// Dark theme variables.
	darkVars := buildThemeVars(dark)

	switch cfg.DarkMode.Kind {
	case config.DarkModeSelectorKind:
		sb.WriteString(":root {\n")
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf(":root%s {\n", cfg.DarkMode.Selector))
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n")

	case config.DarkModeMediaQueryKind:
		sb.WriteString(":root {\n")
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		sb.WriteString("}\n")
		sb.WriteString("@media (prefers-color-scheme: dark) {\n:root {\n")
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n}\n")

	case config.DarkModeBothKind:
		sb.WriteString(":root {\n")
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		sb.WriteString("}\n")
		sb.WriteString("@media (prefers-color-scheme: dark) {\n:root {\n")
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n}\n")
		sb.WriteString(fmt.Sprintf(":root%s {\n", cfg.DarkMode.Selector))
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n")
	}

	return sb.String()
}

// TokenSwitchingCSS generates the color switching rules for token spans.
func TokenSwitchingCSS(cfg *config.Config) string {
	var sb strings.Builder

	sb.WriteString(".kazari-code .kz-line span[style] { color: var(--sl); }\n")

	if cfg.DarkTheme == "" {
		return sb.String()
	}

	switch cfg.DarkMode.Kind {
	case config.DarkModeSelectorKind:
		sb.WriteString(fmt.Sprintf("%s .kazari-code .kz-line span[style] { color: var(--sd); }\n", cfg.DarkMode.Selector))

	case config.DarkModeMediaQueryKind:
		sb.WriteString("@media (prefers-color-scheme: dark) {\n")
		sb.WriteString(".kazari-code .kz-line span[style] { color: var(--sd); }\n")
		sb.WriteString("}\n")

	case config.DarkModeBothKind:
		sb.WriteString("@media (prefers-color-scheme: dark) {\n")
		sb.WriteString(".kazari-code .kz-line span[style] { color: var(--sd); }\n")
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("%s .kazari-code .kz-line span[style] { color: var(--sd); }\n", cfg.DarkMode.Selector))
	}

	return sb.String()
}

func buildThemeVars(tc ThemeColors) []struct{ name, value string } {
	vars := []struct{ name, value string }{
		{"--kz-editor-bg", tc.EditorBG},
		{"--kz-editor-fg", tc.EditorFG},
	}
	if tc.SelectionBG != "" {
		vars = append(vars, struct{ name, value string }{"--kz-selection-bg", tc.SelectionBG})
	}
	if tc.LineNumberFG != "" {
		vars = append(vars, struct{ name, value string }{"--kz-ln-fg", tc.LineNumberFG})
	}
	return vars
}

func writeVars(sb *strings.Builder, vars []struct{ name, value string }) {
	for _, v := range vars {
		sb.WriteString(fmt.Sprintf("  %s: %s;\n", v.name, v.value))
	}
}

func writeVarLines(sb *strings.Builder, vars []struct{ name, value string }) {
	for _, v := range vars {
		if v.value != "" {
			sb.WriteString(fmt.Sprintf("  %s: %s;\n", v.name, v.value))
		}
	}
}
