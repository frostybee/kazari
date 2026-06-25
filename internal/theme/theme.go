package theme

import (
	"fmt"
	"sort"
	"strings"

	"github.com/frostybee/kazari/internal/color"
	"github.com/frostybee/kazari/internal/config"
)

// ThemeColors holds colors extracted from a syntax theme.
type ThemeColors struct {
	EditorBG     string
	EditorFG     string
	SelectionBG  string
	LineNumberFG string
	FoldBG       string
}

func GenerateVars(cfg *config.Config, light, dark ThemeColors) string {
	var sb strings.Builder

	staticVars := buildStaticVars(cfg)

	// Light theme variables.
	lightVars := buildThemeVars(light, cfg)
	// Dark theme variables (only those that differ from light).
	darkVars := buildThemeVars(dark, cfg)
	lightMap := make(map[string]string, len(lightVars))
	for _, v := range lightVars {
		lightMap[v.name] = v.value
	}
	var darkDiffVars []struct{ name, value string }
	for _, v := range darkVars {
		if lightMap[v.name] != v.value {
			darkDiffVars = append(darkDiffVars, v)
		}
	}

	root := cfg.ThemeCSSRoot
	if root == "" {
		root = ":root"
	}

	switch cfg.DarkMode.Kind {
	case config.DarkModeSelectorKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVarLines(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, false)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("%s%s {\n", root, cfg.DarkMode.Selector))
		writeVarLines(&sb, darkDiffVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, true)
		sb.WriteString("}\n")

	case config.DarkModeMediaQueryKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVarLines(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, false)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("@media (prefers-color-scheme: dark) {\n%s {\n", root))
		writeVarLines(&sb, darkDiffVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, true)
		sb.WriteString("}\n}\n")

	case config.DarkModeBothKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVarLines(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, false)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("@media (prefers-color-scheme: dark) {\n%s {\n", root))
		writeVarLines(&sb, darkDiffVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, true)
		sb.WriteString("}\n}\n")
		sb.WriteString(fmt.Sprintf("%s%s {\n", root, cfg.DarkMode.Selector))
		writeVarLines(&sb, darkDiffVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, true)
		sb.WriteString("}\n")
	}

	return sb.String()
}

func buildThemeVars(tc ThemeColors, cfg *config.Config) []struct{ name, value string } {
	vars := blockOverridableVars(tc, cfg)

	if tc.SelectionBG != "" {
		vars = append(vars, nv("--kz-selection-bg", tc.SelectionBG))
	}

	// Collapse section colors derived from the theme's fold background.
	// Emitted after the static defaults, so they win the cascade.
	if cfg.Collapsible != nil && tc.FoldBG != "" {
		vars = append(vars,
			nv("--kz-collapse-closed-bg", color.SetAlpha(tc.FoldBG, 0.2)),
			nv("--kz-collapse-closed-border", color.SetAlpha(tc.FoldBG, 0.5)),
		)
	}

	// Code group tab colors derived from theme luminance.
	if cfg.CodeGroups {
		if color.IsLight(tc.EditorBG) {
			vars = append(vars,
				nv("--kz-group-tab-active-bg", "rgba(255,255,255,0.8)"),
				nv("--kz-group-tab-active-fg", "#24292f"),
				nv("--kz-group-border", "rgba(0,0,0,0.1)"),
			)
		} else {
			vars = append(vars,
				nv("--kz-group-tab-active-bg", "rgba(255,255,255,0.1)"),
				nv("--kz-group-tab-active-fg", "#e6edf3"),
				nv("--kz-group-border", "rgba(255,255,255,0.1)"),
			)
		}
	}

	return vars
}

// blockOverridableVars returns the theme variables that participate in
// per block theme overrides via the theme= meta token. buildThemeVars layers
// the page only extras (selection, fold colors, group tabs) on top of this
// list. BlockOverrideStyle and the kz-themed switching rules derive from the
// same list so the inline emitter and the CSS mappings stay in lockstep.
func blockOverridableVars(tc ThemeColors, cfg *config.Config) []struct{ name, value string } {
	vars := []struct{ name, value string }{
		{"--kz-editor-bg", tc.EditorBG},
		{"--kz-editor-fg", tc.EditorFG},
	}
	if tc.LineNumberFG != "" {
		vars = append(vars, nv("--kz-ln-fg", tc.LineNumberFG))
	}

	if color.IsLight(tc.EditorBG) {
		if tc.LineNumberFG == "" {
			vars = append(vars, nv("--kz-ln-fg", "#6e7781"))
		}
		vars = append(vars,
			nv("--kz-ln-highlight-fg", "#24292f"),
			nv("--kz-gutter-border-color", "rgba(0,0,0,0.1)"),
			nv("--kz-toolbar-bg", "rgba(229, 231, 235, 0.15)"),
			nv("--kz-toolbar-border", "rgba(209, 213, 219, 0.5)"),
			nv("--kz-lang-fg", "#4b5563"),
			nv("--kz-copy-fg", "#4b5563"),
			nv("--kz-copy-fg-hover", "#111827"),
			nv("--kz-copy-bg-hover", "rgba(156, 163, 175, 0.2)"),
		)
		if cfg.Collapsible != nil {
			vars = append(vars,
				nv("--kz-collapse-btn-fg", "#4b5563"),
				nv("--kz-collapse-btn-bg", "rgba(0, 0, 0, 0.04)"),
				nv("--kz-collapse-btn-hover-bg", "rgba(0, 0, 0, 0.08)"),
				nv("--kz-collapse-btn-border", "rgba(0, 0, 0, 0.15)"),
				nv("--kz-collapse-btn-border-hover", "rgba(0, 0, 0, 0.3)"),
			)
		}
		if cfg.ThemedScrollbars {
			vars = append(vars,
				nv("--kz-scrollbar-thumb", "rgba(0, 0, 0, 0.2)"),
				nv("--kz-scrollbar-thumb-hover", "rgba(0, 0, 0, 0.35)"),
			)
		}
	} else {
		if tc.LineNumberFG == "" {
			vars = append(vars, nv("--kz-ln-fg", "#6e7681"))
		}
		vars = append(vars,
			nv("--kz-ln-highlight-fg", "#e6edf3"),
			nv("--kz-gutter-border-color", "rgba(255,255,255,0.1)"),
			nv("--kz-toolbar-bg", "rgba(39, 39, 42, 0.6)"),
			nv("--kz-toolbar-border", "rgba(63, 63, 70, 0.4)"),
			nv("--kz-lang-fg", "#a1a1aa"),
			nv("--kz-copy-fg", "#d4d4d8"),
			nv("--kz-copy-fg-hover", "#ffffff"),
			nv("--kz-copy-bg-hover", "rgba(63, 63, 70, 0.8)"),
		)
		if cfg.Collapsible != nil {
			vars = append(vars,
				nv("--kz-collapse-btn-fg", "#d4d4d8"),
				nv("--kz-collapse-btn-bg", "rgba(255, 255, 255, 0.1)"),
				nv("--kz-collapse-btn-hover-bg", "rgba(255, 255, 255, 0.18)"),
				nv("--kz-collapse-btn-border", "rgba(255, 255, 255, 0.2)"),
				nv("--kz-collapse-btn-border-hover", "rgba(255, 255, 255, 0.35)"),
			)
		}
		if cfg.ThemedScrollbars {
			vars = append(vars,
				nv("--kz-scrollbar-thumb", "rgba(255, 255, 255, 0.15)"),
				nv("--kz-scrollbar-thumb-hover", "rgba(255, 255, 255, 0.3)"),
			)
		}
	}

	// Minimal terminal dots color derived from theme luminance.
	if cfg.TerminalDotStyle == config.DotsMinimal {
		if color.IsLight(tc.EditorBG) {
			vars = append(vars, nv("--kz-terminal-dots-fg", "#24292f"))
		} else {
			vars = append(vars, nv("--kz-terminal-dots-fg", "#e6edf3"))
		}
	}

	return vars
}

// overridableVarNames returns the variable names eligible for per block theme
// overrides. The set depends only on the config, never on theme colors, so a
// placeholder theme is enough to enumerate it.
func overridableVarNames(cfg *config.Config) []string {
	vars := blockOverridableVars(ThemeColors{EditorBG: "#ffffff", EditorFG: "#000000"}, cfg)
	names := make([]string, len(vars))
	for i, v := range vars {
		names[i] = v.name
	}
	return names
}

// BlockOverrideStyle builds the inline style declarations for a per block
// theme override. Light values are emitted as --kz-ovl-* and dark values as
// --kz-ovd-*; the kz-themed rules generated by TokenSwitchingCSS map them
// back onto the regular --kz-* variables. Returns an empty string when the
// light theme colors are unusable.
func BlockOverrideStyle(cfg *config.Config, light, dark ThemeColors) string {
	if light.EditorBG == "" || light.EditorFG == "" {
		return ""
	}
	var sb strings.Builder
	writeOverridePrefixed(&sb, "--kz-ovl-", blockOverridableVars(light, cfg))
	if cfg.DarkTheme != "" && dark.EditorBG != "" && dark.EditorFG != "" {
		writeOverridePrefixed(&sb, "--kz-ovd-", blockOverridableVars(dark, cfg))
	}
	return strings.TrimSuffix(sb.String(), ";")
}

func writeOverridePrefixed(sb *strings.Builder, prefix string, vars []struct{ name, value string }) {
	for _, v := range vars {
		if v.value == "" {
			continue
		}
		sb.WriteString(prefix)
		sb.WriteString(strings.TrimPrefix(v.name, "--kz-"))
		sb.WriteString(":")
		sb.WriteString(v.value)
		sb.WriteString(";")
	}
}

const collapseGradientDecl = "--kz-collapse-gradient-end: var(--kz-editor-bg); "

const (
	lightVarTemplate = "%s: var(--kz-ovl-%s); "
	darkVarTemplate  = "%s: var(--kz-ovd-%s, var(--kz-ovl-%s)); "
)

func nv(name, value string) struct{ name, value string } {
	return struct{ name, value string }{name, value}
}

func writeVarLines(sb *strings.Builder, vars []struct{ name, value string }) {
	for _, v := range vars {
		if v.value != "" {
			sb.WriteString(fmt.Sprintf("  %s: %s;\n", v.name, v.value))
		}
	}
}

func writeStyleOverrides(sb *strings.Builder, overrides map[string]config.StyleValue, isDark bool) {
	if len(overrides) == 0 {
		return
	}
	keys := make([]string, 0, len(overrides))
	for k := range overrides {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sv := overrides[k]
		var v string
		if isDark {
			v = sv.DarkValue()
		} else {
			v = sv.LightValue()
		}
		if v != "" {
			sb.WriteString(fmt.Sprintf("  %s: %s;\n", k, v))
		}
	}
}
