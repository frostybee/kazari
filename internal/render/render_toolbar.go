package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

const copySVG = `<svg class="kz-copy-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>`

const fullscreenSVG = `<svg class="kz-fs-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/></svg>`

const fullscreenExitSVG = `<svg class="kz-fs-exit-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 4v5H4m5 0L4 4M15 4v5h5m-5 0l5-5M9 20v-5H4m5 0l-5 5M15 20v-5h5m-5 0l5 5"/></svg>`

const chevronSVG = `<svg class="kz-collapse-toggle-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>`

const fontIncreaseSVG = `<svg class="kz-font-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6v12m-6-6h12"/></svg>`

const fontDecreaseSVG = `<svg class="kz-font-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 12h12"/></svg>`

const wrapSVG = `<svg class="kz-wrap-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6h18M3 12h15a3 3 0 110 6h-4m0 0l2-2m-2 2l2 2"/></svg>`

const wrapOffSVG = `<svg class="kz-wrap-off-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6h18M3 12h18M3 18h18"/></svg>`

const themeToggleLightSVG = `<svg class="kz-theme-toggle-light-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><circle cx="12" cy="12" r="5" stroke-width="2"/><path stroke-linecap="round" stroke-width="2" d="M12 1v2m0 18v2M4.22 4.22l1.42 1.42m12.72 12.72l1.42 1.42M1 12h2m18 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>`

const themeToggleDarkSVG = `<svg class="kz-theme-toggle-dark-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/></svg>`

func renderToolbar(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	sb.WriteString("<div class=\"kz-toolbar\">")

	sb.WriteString("<div class=\"kz-toolbar-left\">")
	if cfg.LanguageBadge && resolved.Lang != "" {
		renderLangBadge(sb, resolved.Lang, cfg)
	}
	if resolved.Title != "" {
		if cfg.FileIcons {
			ext := fileExt(resolved.Title)
			if ext != "" {
				if cfg.FileIconResolver != nil {
					sb.WriteString(cfg.FileIconResolver(ext))
				} else {
					sb.WriteString(fmt.Sprintf(`<span class="kz-file-icon" data-ext="%s"></span>`, html.EscapeString(ext)))
				}
			}
		}
		sb.WriteString(fmt.Sprintf("<span class=\"kz-title\">%s</span>", html.EscapeString(resolved.Title)))
	}
	sb.WriteString("</div>")

	sb.WriteString("<div class=\"kz-toolbar-right\">")
	renderActionButtons(sb, resolved, cfg)
	if resolved.CollapseThreshold {
		initiallyCollapsed := resolved.CollapseConfig == nil || resolved.CollapseConfig.DefaultCollapsed
		expanded := "false"
		tooltipText := cfg.UIStrings.ExpandButtonText
		if !initiallyCollapsed {
			expanded = "true"
			tooltipText = cfg.UIStrings.CollapseButtonText
		}
		sb.WriteString(fmt.Sprintf(
			"<button class=\"kz-collapse-toggle\" aria-expanded=\"%s\" aria-label=\"%s\" data-tooltip=\"%s\" data-expand=\"%s\" data-collapse=\"%s\">",
			expanded,
			html.EscapeString(tooltipText),
			html.EscapeString(tooltipText),
			html.EscapeString(cfg.UIStrings.ExpandButtonText),
			html.EscapeString(cfg.UIStrings.CollapseButtonText),
		))
		sb.WriteString(chevronSVG)
		sb.WriteString("</button>")
	}
	sb.WriteString("</div>")

	sb.WriteString("</div>")
}

func renderActionButtons(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	if cfg.CopyButton {
		renderCopyButton(sb, resolved.RawCode, cfg)
	}
	if cfg.WrapButton {
		renderWrapButton(sb, resolved, cfg)
	}
	if cfg.ThemeToggle && cfg.DarkTheme != "" {
		renderThemeToggleButton(sb, cfg)
	}
	if cfg.FullscreenButton {
		renderFontControls(sb, cfg)
		renderFullscreenButton(sb, cfg)
	}
}

func renderCopyButton(sb *strings.Builder, rawCode string, cfg *config.Config) {
	encoded := encodeForDataCode(rawCode)
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-copy-btn\" aria-label=\"%s\" data-tooltip=\"%s\" data-copied=\"%s\" data-code=\"%s\">",
		html.EscapeString(cfg.UIStrings.CopyLabel),
		html.EscapeString(cfg.UIStrings.CopyLabel),
		html.EscapeString(cfg.UIStrings.CopySuccess),
		html.EscapeString(encoded),
	))
	sb.WriteString(copySVG)
	sb.WriteString("</button>")
	sb.WriteString(`<span class="kz-sr-announce" aria-live="polite"></span>`)
}

func renderWrapButton(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	pressed := "false"
	title := cfg.UIStrings.WrapEnableLabel
	if resolved.Wrap {
		pressed = "true"
		title = cfg.UIStrings.WrapDisableLabel
	}
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-wrap-btn\" aria-pressed=\"%s\" aria-label=\"%s\" data-tooltip=\"%s\" data-enable=\"%s\" data-disable=\"%s\">",
		pressed,
		html.EscapeString(title),
		html.EscapeString(title),
		html.EscapeString(cfg.UIStrings.WrapEnableLabel),
		html.EscapeString(cfg.UIStrings.WrapDisableLabel),
	))
	sb.WriteString(wrapSVG)
	sb.WriteString(wrapOffSVG)
	sb.WriteString("</button>")
}

func renderFullscreenButton(sb *strings.Builder, cfg *config.Config) {
	sb.WriteString(fmt.Sprintf("<button class=\"kz-fs-btn\" aria-label=\"%s\" data-tooltip=\"%s\" aria-expanded=\"false\">",
		html.EscapeString(cfg.UIStrings.FullscreenLabel),
		html.EscapeString(cfg.UIStrings.FullscreenLabel)))
	sb.WriteString(fullscreenSVG)
	sb.WriteString(fullscreenExitSVG)
	sb.WriteString("</button>")
}

func renderFontControls(sb *strings.Builder, cfg *config.Config) {
	sb.WriteString("<div class=\"kz-font-controls\">")
	sb.WriteString(fmt.Sprintf("<button class=\"kz-font-dec\" aria-label=\"%s\" data-tooltip=\"%s\">",
		html.EscapeString(cfg.UIStrings.FontDecreaseLabel),
		html.EscapeString(cfg.UIStrings.FontDecreaseLabel)))
	sb.WriteString(fontDecreaseSVG)
	sb.WriteString("</button>")
	sb.WriteString(fmt.Sprintf("<button class=\"kz-font-inc\" aria-label=\"%s\" data-tooltip=\"%s\">",
		html.EscapeString(cfg.UIStrings.FontIncreaseLabel),
		html.EscapeString(cfg.UIStrings.FontIncreaseLabel)))
	sb.WriteString(fontIncreaseSVG)
	sb.WriteString("</button>")
	sb.WriteString("</div>")
}

func renderThemeToggleButton(sb *strings.Builder, cfg *config.Config) {
	darkModeStr := "selector"
	switch cfg.DarkMode.Kind {
	case config.DarkModeMediaQueryKind:
		darkModeStr = "media"
	case config.DarkModeBothKind:
		darkModeStr = "both"
	}
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-theme-toggle-btn\" aria-pressed=\"false\" aria-label=\"%s\" data-tooltip=\"%s\" data-label=\"%s\" data-toggled=\"%s\" data-announcement=\"%s\" data-kz-dark-selector=\"%s\" data-kz-dark-mode=\"%s\">",
		html.EscapeString(cfg.UIStrings.ThemeToggleLabel),
		html.EscapeString(cfg.UIStrings.ThemeToggleLabel),
		html.EscapeString(cfg.UIStrings.ThemeToggleLabel),
		html.EscapeString(cfg.UIStrings.ThemeToggleLabel),
		html.EscapeString(cfg.UIStrings.ThemeToggleAnnouncement),
		html.EscapeString(cfg.DarkMode.Selector),
		darkModeStr,
	))
	sb.WriteString(themeToggleLightSVG)
	sb.WriteString(themeToggleDarkSVG)
	sb.WriteString("</button>")
}

func renderLangBadge(sb *strings.Builder, lang string, cfg *config.Config) {
	mode := cfg.LangIconMode
	if mode == config.LangIconOnly || mode == config.LangIconAndText {
		sb.WriteString(fmt.Sprintf(`<span class="kz-lang-icon" data-lang="%s"></span>`, html.EscapeString(lang)))
	}
	if mode != config.LangIconOnly {
		sb.WriteString(fmt.Sprintf(`<span class="kz-lang">%s</span>`, html.EscapeString(displayLang(lang))))
	}
}

func displayLang(lang string) string {
	upper := map[string]string{
		"javascript": "JavaScript", "typescript": "TypeScript",
		"css": "CSS", "html": "HTML", "json": "JSON", "yaml": "YAML",
		"sql": "SQL", "php": "PHP", "xml": "XML", "svg": "SVG",
		"jsx": "JSX", "tsx": "TSX", "graphql": "GraphQL",
	}
	if display, ok := upper[strings.ToLower(lang)]; ok {
		return display
	}
	if len(lang) > 0 {
		return strings.ToUpper(lang[:1]) + lang[1:]
	}
	return lang
}

func fileExt(title string) string {
	idx := strings.LastIndex(title, ".")
	if idx < 0 || idx == len(title)-1 {
		return ""
	}
	return title[idx+1:]
}

func encodeForDataCode(code string) string {
	return strings.ReplaceAll(code, "\n", "\x7f")
}
