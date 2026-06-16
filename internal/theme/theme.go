package theme

import (
	"fmt"
	"sort"
	"strings"

	"github.com/frostybee/kazari/internal/color"
	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/svgutil"
)

// ThemeColors holds colors extracted from a syntax theme.
type ThemeColors struct {
	EditorBG     string
	EditorFG     string
	SelectionBG  string
	LineNumberFG string
	FoldBG       string
}

// GenerateVars produces CSS variable declarations for both themes,
// wrapped according to the dark mode strategy.
func buildStaticVars(cfg *config.Config) []struct{ name, value string } {
	vars := []struct{ name, value string }{
		{"--kz-radius", "0.5rem"},
		{"--kz-shadow", "0 2px 8px rgba(0,0,0,0.15)"},
		{"--kz-border", "1px solid transparent"},
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
		{"--kz-ln-padding-inline", "2ch"},
		{"--kz-gutter-border-width", "1px"},
		{"--kz-mark-bg", "rgba(255,200,0,0.12)"},
		{"--kz-mark-border", "rgba(255,200,0,0.5)"},
		{"--kz-mark-border-width", "3px"},
		{"--kz-ins-bg", "rgba(46,160,67,0.12)"},
		{"--kz-ins-border", "rgba(46,160,67,0.5)"},
		{"--kz-ins-indicator", "'+'"},
		{"--kz-del-bg", "rgba(248,81,73,0.12)"},
		{"--kz-del-border", "rgba(248,81,73,0.5)"},
		{"--kz-del-indicator", "'-'"},
		{"--kz-mark-accent-margin", "0rem"},
		{"--kz-diff-indicator-margin", "0.3rem"},
		{"--kz-ins-indicator-color", "rgba(46,160,67,0.8)"},
		{"--kz-del-indicator-color", "rgba(248,81,73,0.8)"},
		{"--kz-label-mark-bg", "rgba(255,200,0,0.35)"},
		{"--kz-label-ins-bg", "rgba(46,160,67,0.35)"},
		{"--kz-label-del-bg", "rgba(248,81,73,0.35)"},
		{"--kz-label-fg", "#ffffff"},
		{"--kz-label-padding", "0.1rem 0.3rem"},
		{"--kz-label-font-size", "0.75rem"},
		{"--kz-label-radius", "0.2rem"},
		{"--kz-inline-mark-bg", "rgba(255,200,0,0.2)"},
		{"--kz-inline-mark-border", "rgba(255,200,0,0.5)"},
		{"--kz-inline-mark-radius", "0.2rem"},
		{"--kz-inline-mark-padding", "0.15rem"},
		{"--kz-inline-mark-border-width", "1.5px"},
		{"--kz-inline-ins-bg", "rgba(46,160,67,0.2)"},
		{"--kz-inline-ins-border", "rgba(46,160,67,0.5)"},
		{"--kz-inline-del-bg", "rgba(248,81,73,0.2)"},
		{"--kz-inline-del-border", "rgba(248,81,73,0.5)"},
		{"--kz-focus-dimmed-opacity", "0.35"},
		{"--kz-focus-ring", "rgb(59,130,246)"},
		{"--kz-toolbar-padding", "0.25rem 1rem"},
		{"--kz-terminal-header-padding", "0.5rem 1rem"},
		{"--kz-lang-font-size", "0.8rem"},
		{"--kz-lang-font-weight", "500"},
		{"--kz-separator-color", "rgba(161, 161, 170, 0.3)"},
		{"--kz-copy-radius", "0.375rem"},
		{"--kz-copy-success-bg", "rgba(34, 197, 94, 0.9)"},
		{"--kz-copy-success-fg", "#ffffff"},
		{"--kz-copy-success-border", "rgba(34, 197, 94, 0.8)"},
		{"--kz-tooltip-bg", "rgba(30,30,30,0.92)"},
		{"--kz-tooltip-fg", "#ffffff"},
		{"--kz-tooltip-font-size", "0.75rem"},
		{"--kz-tooltip-padding", "0.25rem 0.625rem"},
		{"--kz-tooltip-radius", "999px"},
		{"--kz-tooltip-offset", "6px"},
		{"--kz-tooltip-shadow", "0 2px 6px rgba(0,0,0,0.25)"},
		{"--kz-terminal-bg", "var(--kz-editor-bg)"},
		{"--kz-terminal-titlebar-bg", "var(--kz-toolbar-bg)"},
		{"--kz-terminal-dot-red", "#ff5f57"},
		{"--kz-terminal-dot-yellow", "#febc2e"},
		{"--kz-terminal-dot-green", "#28c840"},
		{"--kz-ln-width", "2ch"},
		{"--kz-ln-opacity", "0.9"},
		{"--kz-ln-highlight-opacity", "0.8"},
	}

	if cfg.Collapsible != nil {
		vars = append(vars,
			struct{ name, value string }{"--kz-collapse-btn-bg", "rgba(255,255,255,0.08)"},
			struct{ name, value string }{"--kz-collapse-btn-fg", "rgba(255,255,255,0.7)"},
			struct{ name, value string }{"--kz-collapse-btn-hover-bg", "rgba(255,255,255,0.15)"},
			struct{ name, value string }{"--kz-collapse-gradient-start", "transparent"},
			struct{ name, value string }{"--kz-collapse-gradient-end", "var(--kz-editor-bg)"},
			struct{ name, value string }{"--kz-collapse-transition", "300ms ease"},
			struct{ name, value string }{"--kz-collapse-closed-bg", "rgb(84 174 255 / 20%)"},
			struct{ name, value string }{"--kz-collapse-closed-border", "rgb(84 174 255 / 50%)"},
			struct{ name, value string }{"--kz-collapse-closed-border-width", "0"},
			struct{ name, value string }{"--kz-collapse-closed-padding", "4px"},
			struct{ name, value string }{"--kz-collapse-open-bg", "transparent"},
			struct{ name, value string }{"--kz-collapse-open-bg-collapsible", "rgb(84 174 255 / 10%)"},
			struct{ name, value string }{"--kz-collapse-open-border", "transparent"},
			struct{ name, value string }{"--kz-collapse-open-border-width", "1px"},
			struct{ name, value string }{"--kz-collapse-closed-fg", "inherit"},
			struct{ name, value string }{"--kz-collapse-closed-font-family", "inherit"},
			struct{ name, value string }{"--kz-collapse-closed-font-size", "inherit"},
			struct{ name, value string }{"--kz-collapse-closed-line-height", "inherit"},
			struct{ name, value string }{"--kz-collapse-expand-icon", `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3E%3Cpath d='m8.177.677 2.896 2.896a.25.25 0 0 1-.177.427H8.75v1.25a.75.75 0 0 1-1.5 0V4H5.104a.25.25 0 0 1-.177-.427L7.823.677a.25.25 0 0 1 .354 0ZM7.25 10.75a.75.75 0 0 1 1.5 0V12h2.146a.25.25 0 0 1 .177.427l-2.896 2.896a.25.25 0 0 1-.354 0l-2.896-2.896A.25.25 0 0 1 5.104 12H7.25v-1.25Zm-5-2a.75.75 0 0 0 0-1.5h-.5a.75.75 0 0 0 0 1.5h.5ZM6 8a.75.75 0 0 1-.75.75h-.5a.75.75 0 0 1 0-1.5h.5A.75.75 0 0 1 6 8Zm2.25.75a.75.75 0 0 0 0-1.5h-.5a.75.75 0 0 0 0 1.5h.5ZM12 8a.75.75 0 0 1-.75.75h-.5a.75.75 0 0 1 0-1.5h.5A.75.75 0 0 1 12 8Zm2.25.75a.75.75 0 0 0 0-1.5h-.5a.75.75 0 0 0 0 1.5h.5Z'/%3E%3C/svg%3E")`},
			struct{ name, value string }{"--kz-collapse-collapse-icon", `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3E%3Cpath d='M10.896 2H8.75V.75a.75.75 0 0 0-1.5 0V2H5.104a.25.25 0 0 0-.177.427l2.896 2.896a.25.25 0 0 0 .354 0l2.896-2.896A.25.25 0 0 0 10.896 2ZM8.75 15.25a.75.75 0 0 1-1.5 0V14H5.104a.25.25 0 0 1-.177-.427l2.896-2.896a.25.25 0 0 1 .354 0l2.896 2.896a.25.25 0 0 1-.177.427H8.75v1.25Zm-6.5-6.5a.75.75 0 0 0 0-1.5h-.5a.75.75 0 0 0 0 1.5h.5ZM6 8a.75.75 0 0 1-.75.75h-.5a.75.75 0 0 1 0-1.5h.5A.75.75 0 0 1 6 8Zm2.25.75a.75.75 0 0 0 0-1.5h-.5a.75.75 0 0 0 0 1.5h.5ZM12 8a.75.75 0 0 1-.75.75h-.5a.75.75 0 0 1 0-1.5h.5A.75.75 0 0 1 12 8Zm2.25.75a.75.75 0 0 0 0-1.5h-.5a.75.75 0 0 0 0 1.5h.5Z'/%3E%3C/svg%3E")`},
		)
	}

	if cfg.CodeGroups {
		vars = append(vars,
			nv("--kz-group-tab-bg", "transparent"),
			nv("--kz-group-tab-fg", "inherit"),
			nv("--kz-group-tab-active-border", "#007acc"),
			nv("--kz-group-tab-padding", "0.5rem 1rem"),
			nv("--kz-group-border-width", "1px"),
			nv("--kz-group-radius", "var(--kz-radius)"),
		)
	}

	if cfg.FileIcons {
		vars = append(vars,
			nv("--kz-file-icon-size", "1rem"),
			nv("--kz-file-icon-margin", "0 0.4rem 0 0"),
			nv("--kz-file-icon-opacity", "0.8"),
		)
	}

	if cfg.LangIconMode != config.LangIconNone {
		vars = append(vars,
			nv("--kz-lang-icon-size", "1.25rem"),
			nv("--kz-lang-icon-margin", "0"),
			nv("--kz-lang-icon-opacity", "0.8"),
		)
	}

	if cfg.ThemedScrollbars {
		vars = append(vars,
			nv("--kz-scrollbar-width", "5px"),
			nv("--kz-scrollbar-height", "5px"),
			nv("--kz-scrollbar-track", "transparent"),
		)
	}

	if cfg.FullscreenButton {
		vars = append(vars,
			nv("--kz-fs-font-scale", "1"),
		)
	}

	if cfg.ThemedSelection {
		vars = append(vars,
			nv("--kz-selection-bg", "rgba(0,122,204,0.3)"),
			nv("--kz-selection-fg", "inherit"),
		)
	}

	if cfg.TerminalDotStyle == config.DotsMinimal {
		dotsSVG := "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 60 16'>" +
			"<circle cx='8' cy='8' r='8'/><circle cx='30' cy='8' r='8'/><circle cx='52' cy='8' r='8'/></svg>"
		vars = append(vars,
			nv("--kz-terminal-dots-opacity", "0.15"),
			nv("--kz-terminal-icon", fmt.Sprintf("url(\"%s\")", svgutil.InlineSVGURL(dotsSVG))),
		)
	}

	return vars
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

	sb.WriteString(".kazari-block .kz-line span[style^=\"--\"] { color: var(--sl, inherit); background-color: var(--slbg, transparent); font-style: var(--sfs, inherit); font-weight: var(--sfw, inherit); text-decoration: var(--std, inherit); }\n")
	sb.WriteString(themedLightRule(cfg))

	if cfg.DarkTheme == "" {
		return sb.String()
	}

	darkRules := ".kazari-block .kz-line span[style^=\"--\"] { color: var(--sd, inherit); background-color: var(--sdbg, transparent); font-style: var(--sfs, inherit); font-weight: var(--sfw, inherit); text-decoration: var(--std, inherit); }\n" +
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

// themedLightRule maps the inline --kz-ovl-* values onto the regular theme
// variables for blocks carrying a per block theme override. The collapse
// gradient end is re declared here because its page level definition
// substitutes var(--kz-editor-bg) at the root, so descendants would inherit
// the page color instead of the override.
func themedLightRule(cfg *config.Config) string {
	var sb strings.Builder
	sb.WriteString(".kazari-block.kz-themed { ")
	for _, name := range overridableVarNames(cfg) {
		suffix := strings.TrimPrefix(name, "--kz-")
		sb.WriteString(fmt.Sprintf("%s: var(--kz-ovl-%s); ", name, suffix))
	}
	if cfg.Collapsible != nil {
		sb.WriteString("--kz-collapse-gradient-end: var(--kz-editor-bg); ")
	}
	sb.WriteString("}\n")
	return sb.String()
}

// themedDarkRule is the dark mode counterpart of themedLightRule. It must
// stay on a single line because writeScopedRules prefixes each line with the
// dark mode selector. Each mapping falls back to the light value so partial
// overrides never leave a variable guaranteed invalid.
func themedDarkRule(cfg *config.Config) string {
	var sb strings.Builder
	sb.WriteString(".kazari-block.kz-themed { ")
	for _, name := range overridableVarNames(cfg) {
		suffix := strings.TrimPrefix(name, "--kz-")
		sb.WriteString(fmt.Sprintf("%s: var(--kz-ovd-%s, var(--kz-ovl-%s)); ", name, suffix, suffix))
	}
	sb.WriteString("}\n")
	return sb.String()
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

	// Force-dark chrome variables on the block.
	sb.WriteString(".kazari-block[data-kz-theme=\"dark\"] { ")
	for _, v := range blockOverridableVars(dark, cfg) {
		if v.value != "" {
			sb.WriteString(fmt.Sprintf("%s: %s; ", v.name, v.value))
		}
	}
	if cfg.Collapsible != nil {
		sb.WriteString("--kz-collapse-gradient-end: var(--kz-editor-bg); ")
	}
	sb.WriteString("}\n")

	// Force-light chrome variables on the block.
	sb.WriteString(".kazari-block[data-kz-theme=\"light\"] { ")
	for _, v := range blockOverridableVars(light, cfg) {
		if v.value != "" {
			sb.WriteString(fmt.Sprintf("%s: %s; ", v.name, v.value))
		}
	}
	if cfg.Collapsible != nil {
		sb.WriteString("--kz-collapse-gradient-end: var(--kz-editor-bg); ")
	}
	sb.WriteString("}\n")

	// Force-dark token switching (overrides page-level light or dark rules).
	sb.WriteString(".kazari-block[data-kz-theme=\"dark\"] .kz-line span[style^=\"--\"] { color: var(--sd, inherit); background-color: var(--sdbg, transparent); font-style: var(--sfs, inherit); font-weight: var(--sfw, inherit); text-decoration: var(--std, inherit); }\n")

	// Force-light token switching (overrides page-level dark rules).
	sb.WriteString(".kazari-block[data-kz-theme=\"light\"] .kz-line span[style^=\"--\"] { color: var(--sl, inherit); background-color: var(--slbg, transparent); font-style: var(--sfs, inherit); font-weight: var(--sfw, inherit); text-decoration: var(--std, inherit); }\n")

	// kz-themed + force-dark: remap override dark vars onto kz-* names.
	sb.WriteString(themedToggleDarkRule(cfg))
	// kz-themed + force-light: remap override light vars onto kz-* names.
	sb.WriteString(themedToggleLightRule(cfg))

	return sb.String()
}

func themedToggleDarkRule(cfg *config.Config) string {
	var sb strings.Builder
	sb.WriteString(".kazari-block.kz-themed[data-kz-theme=\"dark\"] { ")
	for _, name := range overridableVarNames(cfg) {
		suffix := strings.TrimPrefix(name, "--kz-")
		sb.WriteString(fmt.Sprintf("%s: var(--kz-ovd-%s, var(--kz-ovl-%s)); ", name, suffix, suffix))
	}
	if cfg.Collapsible != nil {
		sb.WriteString("--kz-collapse-gradient-end: var(--kz-editor-bg); ")
	}
	sb.WriteString("}\n")
	return sb.String()
}

func themedToggleLightRule(cfg *config.Config) string {
	var sb strings.Builder
	sb.WriteString(".kazari-block.kz-themed[data-kz-theme=\"light\"] { ")
	for _, name := range overridableVarNames(cfg) {
		suffix := strings.TrimPrefix(name, "--kz-")
		sb.WriteString(fmt.Sprintf("%s: var(--kz-ovl-%s); ", name, suffix))
	}
	if cfg.Collapsible != nil {
		sb.WriteString("--kz-collapse-gradient-end: var(--kz-editor-bg); ")
	}
	sb.WriteString("}\n")
	return sb.String()
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
