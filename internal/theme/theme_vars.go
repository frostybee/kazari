package theme

import (
	"fmt"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/svgutil"
)

func buildStaticVars(cfg *config.Config) []struct{ name, value string } {
	vars := []struct{ name, value string }{
		{"--kz-radius", "0.5rem"},
		{"--kz-shadow", "0 2px 8px rgba(0,0,0,0.15)"},
		{"--kz-border", "1px solid transparent"},
		{"--kz-transition", "150ms ease"},
		{"--kz-font-family", "'JetBrains Mono Variable', monospace"},
		{"--kz-font-size", "0.875rem"},
		{"--kz-font-weight", "500"},
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
		{"--kz-tooltip-padding", "0.35rem 0.75rem"},
		{"--kz-tooltip-radius", "6px"},
		{"--kz-tooltip-offset", "6px"},
		{"--kz-tooltip-shadow", "0 2px 6px rgba(0,0,0,0.25)"},
		{"--kz-tooltip-arrow-size", "5px"},
		{"--kz-terminal-bg", "var(--kz-editor-bg)"},
		{"--kz-terminal-titlebar-bg", "var(--kz-toolbar-bg)"},
		{"--kz-terminal-dot-red", "#ff5f57"},
		{"--kz-terminal-dot-yellow", "#febc2e"},
		{"--kz-terminal-dot-green", "#28c840"},
		{"--kz-ln-width", "2ch"},
		{"--kz-ln-opacity", "1"},
		{"--kz-ln-highlight-opacity", "0.8"},
		{"--kz-ansi-black", "#000000"},
		{"--kz-ansi-red", "#cc0000"},
		{"--kz-ansi-green", "#4e9a06"},
		{"--kz-ansi-yellow", "#c4a000"},
		{"--kz-ansi-blue", "#3465a4"},
		{"--kz-ansi-magenta", "#75507b"},
		{"--kz-ansi-cyan", "#06989a"},
		{"--kz-ansi-white", "#d3d7cf"},
		{"--kz-ansi-bright-black", "#555753"},
		{"--kz-ansi-bright-red", "#ef2929"},
		{"--kz-ansi-bright-green", "#8ae234"},
		{"--kz-ansi-bright-yellow", "#fce94f"},
		{"--kz-ansi-bright-blue", "#729fcf"},
		{"--kz-ansi-bright-magenta", "#ad7fa8"},
		{"--kz-ansi-bright-cyan", "#34e2e2"},
		{"--kz-ansi-bright-white", "#eeeeec"},
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
			struct{ name, value string }{"--kz-collapse-closed-fg", "currentColor"},
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
			nv("--kz-scrollbar-width", "8px"),
			nv("--kz-scrollbar-height", "8px"),
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
