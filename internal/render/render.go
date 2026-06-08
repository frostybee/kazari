package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// MergedToken holds both light and dark colors for a single token.
type MergedToken struct {
	Content    string
	LightColor string
	DarkColor  string
	LightBG    string
	DarkBG     string
	FontStyle  int
}

// TokenLine represents one line of merged tokens.
type TokenLine struct {
	Tokens []MergedToken
}

const (
	fontStyleItalic        = 1
	fontStyleBold          = 2
	fontStyleUnderline     = 4
	fontStyleStrikethrough = 8
)

const copySVG = `<svg class="kz-copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>`

const fullscreenSVG = `<svg class="kz-fs-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/></svg>`

// RenderBlock produces the full HTML for a code block.
func RenderBlock(lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config) string {
	var sb strings.Builder
	dualTheme := hasDualTheme(lines)

	sb.WriteString("<div class=\"kazari-code\">\n")

	if resolved.Frame == config.FrameNone {
		renderNoFrame(&sb, lines, resolved, cfg, dualTheme)
	} else {
		renderFramedBlock(&sb, lines, resolved, cfg, dualTheme)
	}

	sb.WriteString("</div>")
	return sb.String()
}

func renderFramedBlock(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	if resolved.Frame == config.FrameTerminal {
		renderTerminalFrame(sb, lines, resolved, cfg, dualTheme)
		return
	}

	classes := "frame"
	if resolved.Title != "" {
		classes += " has-title"
	}
	sb.WriteString(fmt.Sprintf("<figure class=\"%s\" data-lang=\"%s\">", classes, html.EscapeString(resolved.Lang)))

	renderToolbar(sb, resolved, cfg)
	renderPreCode(sb, lines, resolved, dualTheme)

	sb.WriteString("</figure>\n")
}

func renderTerminalFrame(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	classes := "frame is-terminal"
	if resolved.Title != "" {
		classes += " has-title"
	}
	sb.WriteString(fmt.Sprintf("<figure class=\"%s\" data-lang=\"%s\">", classes, html.EscapeString(resolved.Lang)))

	sb.WriteString("<div class=\"kz-terminal-header\">")
	sb.WriteString("<span class=\"kz-terminal-dots\"><span></span><span></span><span></span></span>")
	if resolved.Title != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-title\">%s</span>", html.EscapeString(resolved.Title)))
	}
	if cfg.CopyButton {
		sb.WriteString("<div class=\"kz-terminal-actions\">")
		renderCopyButton(sb, resolved.RawCode)
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")

	renderPreCode(sb, lines, resolved, dualTheme)

	sb.WriteString("</figure>\n")
}

func renderToolbar(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	sb.WriteString("<div class=\"kz-toolbar\">")

	// Left section
	sb.WriteString("<div class=\"kz-toolbar-left\">")
	if resolved.Title != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-title\">%s</span>", html.EscapeString(resolved.Title)))
	} else if cfg.LanguageBadge && resolved.Lang != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-lang\">%s</span>", html.EscapeString(displayLang(resolved.Lang))))
	}
	sb.WriteString("</div>")

	// Right section
	sb.WriteString("<div class=\"kz-toolbar-right\">")
	if resolved.Title != "" && cfg.LanguageBadge && resolved.Lang != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-lang\">%s</span>", html.EscapeString(displayLang(resolved.Lang))))
	}
	if cfg.FullscreenButton {
		sb.WriteString("<button class=\"kz-fs-btn\" aria-label=\"Fullscreen\">")
			sb.WriteString(fullscreenSVG)
			sb.WriteString("</button>")
	}
	if cfg.CopyButton {
		renderCopyButton(sb, resolved.RawCode)
	}
	sb.WriteString("</div>")

	sb.WriteString("</div>")
}

func renderCopyButton(sb *strings.Builder, rawCode string) {
	encoded := encodeForDataCode(rawCode)
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-copy-btn\" title=\"Copy to clipboard\" data-copied=\"Copied!\" data-code=\"%s\">",
		html.EscapeString(encoded),
	))
	sb.WriteString(copySVG)
	sb.WriteString("<span>Copy</span>")
	sb.WriteString("</button>")
}

func renderNoFrame(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	renderPreCode(sb, lines, resolved, dualTheme)
}

func renderPreCode(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, dualTheme bool) {
	if resolved.LineNumbers {
		endNum := resolved.StartLineNumber + len(lines) - 1
		maxDigits := max(digitCount(resolved.StartLineNumber), digitCount(endNum))
		if maxDigits > 2 {
			sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\"><code style=\"--kz-ln-width:%dch\">", html.EscapeString(resolved.Lang), maxDigits))
		} else {
			sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\"><code>", html.EscapeString(resolved.Lang)))
		}
	} else {
		sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\"><code>", html.EscapeString(resolved.Lang)))
	}
	for i, line := range lines {
		lineNum := resolved.StartLineNumber + i
		renderLine(sb, line, lineNum, resolved, dualTheme)
	}
	sb.WriteString("</code></pre>")
}

func encodeForDataCode(code string) string {
	return strings.ReplaceAll(code, "\n", "\x7f")
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

func renderLine(sb *strings.Builder, line TokenLine, lineNum int, resolved *config.ResolvedBlock, dualTheme bool) {
	sb.WriteString("<div class=\"kz-line\">")
	if resolved.LineNumbers {
		sb.WriteString(fmt.Sprintf("<div class=\"gutter\"><div class=\"ln\" aria-hidden=\"true\">%d</div></div>", lineNum))
	}
	sb.WriteString("<div class=\"code\">")
	for _, tok := range line.Tokens {
		if tok.Content == "" {
			continue
		}
		renderToken(sb, tok, dualTheme)
	}
	sb.WriteString("</div></div>")
}

func renderToken(sb *strings.Builder, tok MergedToken, dualTheme bool) {
	style := buildTokenStyle(tok, dualTheme)
	if style != "" {
		sb.WriteString(fmt.Sprintf("<span style=\"%s\">", style))
	} else {
		sb.WriteString("<span>")
	}
	sb.WriteString(html.EscapeString(tok.Content))
	sb.WriteString("</span>")
}

func buildTokenStyle(tok MergedToken, dualTheme bool) string {
	var parts []string

	if tok.LightColor != "" {
		parts = append(parts, fmt.Sprintf("--sl:%s", tok.LightColor))
	}
	if dualTheme && tok.DarkColor != "" {
		parts = append(parts, fmt.Sprintf("--sd:%s", tok.DarkColor))
	}
	if tok.LightBG != "" {
		parts = append(parts, fmt.Sprintf("--slbg:%s", tok.LightBG))
	}
	if dualTheme && tok.DarkBG != "" {
		parts = append(parts, fmt.Sprintf("--sdbg:%s", tok.DarkBG))
	}

	if tok.FontStyle&fontStyleItalic != 0 {
		parts = append(parts, "--sfs:italic")
	}
	if tok.FontStyle&fontStyleBold != 0 {
		parts = append(parts, "--sfw:bold")
	}
	if tok.FontStyle&(fontStyleUnderline|fontStyleStrikethrough) != 0 {
		var decs []string
		if tok.FontStyle&fontStyleUnderline != 0 {
			decs = append(decs, "underline")
		}
		if tok.FontStyle&fontStyleStrikethrough != 0 {
			decs = append(decs, "line-through")
		}
		parts = append(parts, fmt.Sprintf("--std:%s", strings.Join(decs, " ")))
	}

	return strings.Join(parts, ";")
}

func digitCount(n int) int {
	if n < 0 {
		n = -n
	}
	if n == 0 {
		return 1
	}
	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count
}

func hasDualTheme(lines []TokenLine) bool {
	for _, line := range lines {
		for _, tok := range line.Tokens {
			if tok.DarkColor != "" {
				return true
			}
		}
	}
	return false
}
