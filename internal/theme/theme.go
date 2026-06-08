package theme

import (
	"fmt"
	"math"
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

	// Toolbar, badge, and copy button colors derived from theme luminance.
	if isLightTheme(tc.EditorBG) {
		vars = append(vars,
			nv("--kz-toolbar-bg", "rgba(229, 231, 235, 0.15)"),
			nv("--kz-toolbar-border", "rgba(209, 213, 219, 0.5)"),
			nv("--kz-lang-fg", "#4b5563"),
			nv("--kz-lang-bg", "rgba(209, 213, 219, 0.4)"),
			nv("--kz-copy-fg", "#4b5563"),
			nv("--kz-copy-bg", "rgba(209, 213, 219, 0.4)"),
			nv("--kz-copy-border", "rgba(156, 163, 175, 0.2)"),
			nv("--kz-copy-fg-hover", "#111827"),
			nv("--kz-copy-bg-hover", "rgba(156, 163, 175, 0.2)"),
		)
	} else {
		vars = append(vars,
			nv("--kz-toolbar-bg", "rgba(39, 39, 42, 0.6)"),
			nv("--kz-toolbar-border", "rgba(63, 63, 70, 0.4)"),
			nv("--kz-lang-fg", "#a1a1aa"),
			nv("--kz-lang-bg", "rgba(51, 51, 55, 0.8)"),
			nv("--kz-copy-fg", "#d4d4d8"),
			nv("--kz-copy-bg", "rgba(51, 51, 55, 0.9)"),
			nv("--kz-copy-border", "rgba(63, 63, 70, 0.5)"),
			nv("--kz-copy-fg-hover", "#ffffff"),
			nv("--kz-copy-bg-hover", "rgba(63, 63, 70, 0.8)"),
		)
	}

	return vars
}

func nv(name, value string) struct{ name, value string } {
	return struct{ name, value string }{name, value}
}

func isLightTheme(bg string) bool {
	if bg == "" {
		return false
	}
	r, g, b, _, err := parseHexSimple(bg)
	if err {
		return false
	}
	lum := 0.2126*linearize(r) + 0.7152*linearize(g) + 0.0722*linearize(b)
	return lum > 0.5
}

func linearize(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func parseHexSimple(hex string) (r, g, b, a float64, fail bool) {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	a = 1.0
	switch len(hex) {
	case 3:
		var ri, gi, bi uint8
		_, err := fmt.Sscanf(hex, "%1x%1x%1x", &ri, &gi, &bi)
		if err != nil {
			return 0, 0, 0, 0, true
		}
		r = float64(ri*17) / 255
		g = float64(gi*17) / 255
		b = float64(bi*17) / 255
	case 6:
		var ri, gi, bi uint8
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &ri, &gi, &bi)
		if err != nil {
			return 0, 0, 0, 0, true
		}
		r = float64(ri) / 255
		g = float64(gi) / 255
		b = float64(bi) / 255
	default:
		return 0, 0, 0, 0, true
	}
	return r, g, b, a, false
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
