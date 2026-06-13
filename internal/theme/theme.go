package theme

import (
	"fmt"
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
func GenerateVars(cfg *config.Config, light, dark ThemeColors) string {
	var sb strings.Builder

	// Static defaults (not theme-dependent).
	staticVars := []struct{ name, value string }{
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
		// Line marker defaults
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
		// Label defaults
		{"--kz-label-mark-bg", "rgba(255,200,0,0.35)"},
		{"--kz-label-ins-bg", "rgba(46,160,67,0.35)"},
		{"--kz-label-del-bg", "rgba(248,81,73,0.35)"},
		{"--kz-label-fg", "#ffffff"},
		{"--kz-label-padding", "0.1rem 0.3rem"},
		{"--kz-label-font-size", "0.75rem"},
		{"--kz-label-radius", "0.2rem"},
		// Inline marker defaults
		{"--kz-inline-mark-bg", "rgba(255,200,0,0.2)"},
		{"--kz-inline-mark-border", "rgba(255,200,0,0.5)"},
		{"--kz-inline-mark-radius", "0.2rem"},
		{"--kz-inline-mark-padding", "0.15rem"},
		{"--kz-inline-mark-border-width", "1.5px"},
		{"--kz-inline-ins-bg", "rgba(46,160,67,0.2)"},
		{"--kz-inline-ins-border", "rgba(46,160,67,0.5)"},
		{"--kz-inline-del-bg", "rgba(248,81,73,0.2)"},
		{"--kz-inline-del-border", "rgba(248,81,73,0.5)"},
		// Focus defaults
		{"--kz-focus-dimmed-opacity", "0.35"},
	}


	// Collapsible defaults (conditional)
	if cfg.Collapsible != nil {
		staticVars = append(staticVars,
			// Threshold button
			struct{ name, value string }{"--kz-collapse-btn-bg", "rgba(255,255,255,0.08)"},
			struct{ name, value string }{"--kz-collapse-btn-fg", "rgba(255,255,255,0.7)"},
			struct{ name, value string }{"--kz-collapse-btn-hover-bg", "rgba(255,255,255,0.15)"},
			struct{ name, value string }{"--kz-collapse-gradient-start", "transparent"},
			struct{ name, value string }{"--kz-collapse-gradient-end", "var(--kz-editor-bg)"},
			struct{ name, value string }{"--kz-collapse-transition", "300ms ease"},
			// Range section (EC-matching defaults)
			struct{ name, value string }{"--kz-collapse-closed-bg", "rgb(84 174 255 / 20%)"},
			struct{ name, value string }{"--kz-collapse-closed-border", "rgb(84 174 255 / 50%)"},
			struct{ name, value string }{"--kz-collapse-closed-border-width", "0"},
			struct{ name, value string }{"--kz-collapse-closed-padding", "4px"},
			struct{ name, value string }{"--kz-collapse-open-bg", "transparent"},
			struct{ name, value string }{"--kz-collapse-open-bg-collapsible", "rgb(84 174 255 / 10%)"},
			struct{ name, value string }{"--kz-collapse-open-border", "transparent"},
			struct{ name, value string }{"--kz-collapse-open-border-width", "1px"},
		)
	}

	// Code group defaults (conditional)
	if cfg.CodeGroups {
		staticVars = append(staticVars,
			nv("--kz-group-tab-bg", "transparent"),
			nv("--kz-group-tab-fg", "inherit"),
			nv("--kz-group-tab-active-border", "#007acc"),
			nv("--kz-group-tab-padding", "0.5rem 1rem"),
			nv("--kz-group-border-width", "1px"),
			nv("--kz-group-radius", "var(--kz-radius)"),
		)
	}

	// File icon defaults (conditional)
	if cfg.FileIcons {
		staticVars = append(staticVars,
			nv("--kz-file-icon-size", "1rem"),
			nv("--kz-file-icon-margin", "0 0.4rem 0 0"),
			nv("--kz-file-icon-opacity", "0.8"),
		)
	}

	// Minimal terminal dots (conditional)
	if cfg.TerminalDotStyle == config.DotsMinimal {
		dotsSVG := "<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 60 16'>" +
			"<circle cx='8' cy='8' r='8'/><circle cx='30' cy='8' r='8'/><circle cx='52' cy='8' r='8'/></svg>"
		staticVars = append(staticVars,
			nv("--kz-terminal-dots-opacity", "0.15"),
			nv("--kz-terminal-icon", fmt.Sprintf("url(\"%s\")", svgutil.InlineSVGURL(dotsSVG))),
		)
	}

	// Light theme variables.
	lightVars := buildThemeVars(light, cfg)
	// Dark theme variables.
	darkVars := buildThemeVars(dark, cfg)

	root := cfg.ThemeCSSRoot
	if root == "" {
		root = ":root"
	}

	switch cfg.DarkMode.Kind {
	case config.DarkModeSelectorKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("%s%s {\n", root, cfg.DarkMode.Selector))
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n")

	case config.DarkModeMediaQueryKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("@media (prefers-color-scheme: dark) {\n%s {\n", root))
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n}\n")

	case config.DarkModeBothKind:
		sb.WriteString(fmt.Sprintf("%s {\n", root))
		writeVars(&sb, staticVars)
		writeVarLines(&sb, lightVars)
		sb.WriteString("}\n")
		sb.WriteString(fmt.Sprintf("@media (prefers-color-scheme: dark) {\n%s {\n", root))
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n}\n")
		sb.WriteString(fmt.Sprintf("%s%s {\n", root, cfg.DarkMode.Selector))
		writeVarLines(&sb, darkVars)
		sb.WriteString("}\n")
	}

	return sb.String()
}

// TokenSwitchingCSS generates the color switching rules for token spans.
func TokenSwitchingCSS(cfg *config.Config) string {
	var sb strings.Builder

	sb.WriteString(".kazari-code .kz-line span[style] { color: var(--sl); }\n")
	sb.WriteString(themedLightRule(cfg))

	if cfg.DarkTheme == "" {
		return sb.String()
	}

	darkRules := ".kazari-code .kz-line span[style] { color: var(--sd); }\n" +
		".kazari-code .kz-line span[style*=\"--sdbg\"] { background-color: var(--sdbg); }\n" +
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
	sb.WriteString(".kazari-code.kz-themed { ")
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
	sb.WriteString(".kazari-code.kz-themed { ")
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
