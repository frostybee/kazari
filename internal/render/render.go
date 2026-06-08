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
	LightColor string // --sl value
	DarkColor  string // --sd value (empty in single-theme mode)
	LightBG    string // --slbg (rare)
	DarkBG     string // --sdbg (rare)
	FontStyle  int    // bitmask: Italic=1, Bold=2, Underline=4, Strikethrough=8
}

// TokenLine represents one line of merged tokens.
type TokenLine struct {
	Tokens []MergedToken
}

// RenderBlock produces the full HTML for a code block.
func RenderBlock(lines []TokenLine, resolved *config.ResolvedBlock) string {
	var sb strings.Builder
	dualTheme := hasDualTheme(lines)

	sb.WriteString("<div class=\"kazari-code\">\n")
	sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\"><code>", html.EscapeString(resolved.Lang)))

	for _, line := range lines {
		renderLine(&sb, line, dualTheme)
	}

	sb.WriteString("</code></pre>\n")
	sb.WriteString("</div>")

	return sb.String()
}

func renderLine(sb *strings.Builder, line TokenLine, dualTheme bool) {
	sb.WriteString("<div class=\"kz-line\"><div class=\"code\">")
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

	// Font style bitmask
	if tok.FontStyle&1 != 0 { // Italic
		parts = append(parts, "--sfs:italic")
	}
	if tok.FontStyle&2 != 0 { // Bold
		parts = append(parts, "--sfw:bold")
	}
	if tok.FontStyle&(4|8) != 0 { // Underline and/or Strikethrough
		var decs []string
		if tok.FontStyle&4 != 0 {
			decs = append(decs, "underline")
		}
		if tok.FontStyle&8 != 0 {
			decs = append(decs, "line-through")
		}
		parts = append(parts, fmt.Sprintf("--std:%s", strings.Join(decs, " ")))
	}

	return strings.Join(parts, ";")
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
