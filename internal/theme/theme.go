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

// KnownVarNames returns the deduplicated list of all CSS variable names
// that Kazari can emit for the given config. Used for style override validation.
func KnownVarNames(cfg *config.Config) []string {
	seen := make(map[string]struct{})
	for _, v := range buildStaticVars(cfg) {
		seen[v.name] = struct{}{}
	}
	for _, name := range overridableVarNames(cfg) {
		seen[name] = struct{}{}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	return names
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
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, false)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("%s%s {\n", root, cfg.DarkMode.Selector))
		writeVarLines(&sb, darkDiffVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, true)
		sb.WriteString("}\n")

	case config.DarkModeMediaQueryKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, false)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("@media (prefers-color-scheme: dark) {\n%s {\n", root))
		writeVarLines(&sb, darkDiffVars)
		writeStyleOverrides(&sb, cfg.StyleOverrides, true)
		sb.WriteString("}\n}\n")

	case config.DarkModeBothKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVars(&sb, staticVars)
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

// TokenSwitchingCSS generates the color switching rules for token spans.
func TokenSwitchingCSS(cfg *config.Config) string {
	var sb strings.Builder

	sb.WriteString(tokenSwitchRule(".kazari-block", "--sl", "--slbg"))
	sb.WriteString(themedLightRule(cfg))

	if cfg.DarkTheme == "" {
		return sb.String()
	}

	darkRules := tokenSwitchRule(".kazari-block", "--sd", "--sdbg") +
		themedDarkRule(cfg)

	switch cfg.DarkMode.Kind {
	case config.DarkModeSelectorKind:
		writeScopedRules(&sb, cfg.DarkMode.Selector, darkRules)

	case config.DarkModeMediaQueryKind:
		sb.WriteString("@media (prefers-color-scheme: dark) {\n")
		sb.WriteString(darkRules)
		sb.WriteString("}\n")

	case config.DarkModeBothKind:
		sb.WriteString("@media (prefers-color-scheme: dark) {\n")
		sb.WriteString(darkRules)
		sb.WriteString("}\n")
		writeScopedRules(&sb, cfg.DarkMode.Selector, darkRules)
	}

	return sb.String()
}

// writeScopedRules prefixes each rule line with the dark mode selector.
func writeScopedRules(sb *strings.Builder, selector, rules string) {
	for _, line := range strings.Split(strings.TrimRight(rules, "\n"), "\n") {
		sb.WriteString(selector + " " + line + "\n")
	}
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

func writeThemedRule(selector, varTemplate string, includeGradient bool, cfg *config.Config) string {
	var sb strings.Builder
	sb.WriteString(selector)
	sb.WriteString(" { ")
	for _, name := range overridableVarNames(cfg) {
		suffix := strings.TrimPrefix(name, "--kz-")
		sb.WriteString(fmt.Sprintf(varTemplate, name, suffix, suffix))
	}
	if includeGradient && cfg.Collapsible != nil {
		sb.WriteString(collapseGradientDecl)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func tokenSwitchRule(selector, colorVar, bgVar string) string {
	return fmt.Sprintf(
		"%s .kz-line span[style^=\"--\"] { color: var(%s, inherit); background-color: var(%s, transparent); font-style: var(--sfs, inherit); font-weight: var(--sfw, inherit); text-decoration: var(--std, inherit); }\n",
		selector, colorVar, bgVar,
	)
}

func themedLightRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed", lightVarTemplate, true, cfg)
}

// themedDarkRule must stay on a single line because writeScopedRules prefixes
// each line with the dark mode selector. No gradient re-declaration here
// because it is emitted inside the scoped wrapper.
func themedDarkRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed", darkVarTemplate, false, cfg)
}

func nv(name, value string) struct{ name, value string } {
	return struct{ name, value string }{name, value}
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

// ThemeToggleCSS generates CSS rules for per-block theme toggling.
// These rules use data-kz-theme="dark"|"light" on .kazari-block to
// override the page-level theme for individual blocks.
func ThemeToggleCSS(cfg *config.Config, light, dark ThemeColors) string {
	if cfg.DarkTheme == "" {
		return ""
	}

	var sb strings.Builder

	writeToggleVars(&sb, ".kazari-block[data-kz-theme=\"dark\"]", blockOverridableVars(dark, cfg), cfg)
	writeToggleVars(&sb, ".kazari-block[data-kz-theme=\"light\"]", blockOverridableVars(light, cfg), cfg)

	sb.WriteString(tokenSwitchRule(".kazari-block[data-kz-theme=\"dark\"]", "--sd", "--sdbg"))
	sb.WriteString(tokenSwitchRule(".kazari-block[data-kz-theme=\"light\"]", "--sl", "--slbg"))

	// kz-themed + force-dark: remap override dark vars onto kz-* names.
	sb.WriteString(themedToggleDarkRule(cfg))
	// kz-themed + force-light: remap override light vars onto kz-* names.
	sb.WriteString(themedToggleLightRule(cfg))

	return sb.String()
}

func themedToggleDarkRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed[data-kz-theme=\"dark\"]", darkVarTemplate, true, cfg)
}

func themedToggleLightRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed[data-kz-theme=\"light\"]", lightVarTemplate, true, cfg)
}

func writeToggleVars(sb *strings.Builder, selector string, vars []struct{ name, value string }, cfg *config.Config) {
	sb.WriteString(selector)
	sb.WriteString(" { ")
	for _, v := range vars {
		if v.value != "" {
			sb.WriteString(fmt.Sprintf("%s: %s; ", v.name, v.value))
		}
	}
	if cfg.Collapsible != nil {
		sb.WriteString(collapseGradientDecl)
	}
	sb.WriteString("}\n")
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
		if isDark {
			if !sv.IsThemed() {
				continue
			}
			if v := sv.DarkValue(); v != "" {
				sb.WriteString(fmt.Sprintf("  %s: %s;\n", k, v))
			}
		} else {
			var v string
			if sv.IsThemed() {
				v = sv.LightValue()
			} else {
				v = sv.Value
			}
			if v != "" {
				sb.WriteString(fmt.Sprintf("  %s: %s;\n", k, v))
			}
		}
	}
}
