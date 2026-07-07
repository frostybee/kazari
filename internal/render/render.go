package render

import (
	"fmt"
	"hash/fnv"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// RenderBlock produces the full HTML for a code block.
func RenderBlock(lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config) string {
	var sb strings.Builder
	dualTheme := hasDualTheme(lines)

	wrapperClass := "kazari-block"
	if resolved.ThemeOverrideStyle != "" {
		wrapperClass += " kz-themed"
	}
	if resolved.CollapseThreshold && (resolved.CollapseConfig == nil || resolved.CollapseConfig.DefaultCollapsed) {
		wrapperClass += " kz-collapsed"
	}
	wrapperClass += " not-content"
	attrs := fmt.Sprintf(" class=\"%s\"", wrapperClass)
	if cfg.DataLineCount {
		attrs += fmt.Sprintf(" data-lines=\"%d\"", len(lines))
	}
	if cfg.ThemeToggle && cfg.DarkTheme != "" {
		attrs += fmt.Sprintf(" data-kz-id=\"%s\"", blockID(resolved.RawCode))
	}
	if resolved.ThemeOverrideStyle != "" {
		attrs += fmt.Sprintf(" style=\"%s\"", html.EscapeString(resolved.ThemeOverrideStyle))
	}
	sb.WriteString(fmt.Sprintf("<div%s>\n", attrs))

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
	renderCollapseContentStart(sb, resolved)
	renderPreCode(sb, lines, resolved, cfg, dualTheme)
	renderCollapseContentEnd(sb, resolved)

	if resolved.CollapseThreshold {
		renderCollapseBar(sb, resolved, cfg)
	}

	if resolved.OutputText != "" {
		renderOutputPanel(sb, resolved, cfg)
	}

	sb.WriteString("</figure>\n")
}

func renderTerminalFrame(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	classes := "frame is-terminal"
	if resolved.Title != "" {
		classes += " has-title"
	}
	sb.WriteString(fmt.Sprintf("<figure class=\"%s\" data-lang=\"%s\">", classes, html.EscapeString(resolved.Lang)))

	if cfg.TerminalDotStyle == config.DotsMinimal {
		sb.WriteString("<div class=\"kz-terminal-header kz-dots-minimal\">")
	} else {
		sb.WriteString("<div class=\"kz-terminal-header\">")
		sb.WriteString("<span class=\"kz-terminal-dots\" aria-hidden=\"true\"><span></span><span></span><span></span></span>")
	}
	if resolved.Title != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-title\">%s</span>", html.EscapeString(resolved.Title)))
	} else {
		sb.WriteString(fmt.Sprintf("<span class=\"sr-only\">%s</span>", html.EscapeString(cfg.UIStrings.TerminalWindowLabel)))
	}
	if cfg.CopyButton || cfg.WrapButton || cfg.FullscreenButton || (cfg.ThemeToggle && cfg.DarkTheme != "") {
		sb.WriteString("<div class=\"kz-terminal-actions\">")
		renderActionButtons(sb, resolved, cfg)
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")

	renderCollapseContentStart(sb, resolved)
	renderPreCode(sb, lines, resolved, cfg, dualTheme)
	renderCollapseContentEnd(sb, resolved)

	if resolved.CollapseThreshold {
		renderCollapseBar(sb, resolved, cfg)
	}

	if resolved.OutputText != "" {
		renderOutputPanel(sb, resolved, cfg)
	}

	sb.WriteString("</figure>\n")
}

func renderNoFrame(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	renderCollapseContentStart(sb, resolved)
	renderPreCode(sb, lines, resolved, cfg, dualTheme)
	renderCollapseContentEnd(sb, resolved)

	if resolved.CollapseThreshold {
		renderCollapseBar(sb, resolved, cfg)
	}

	if resolved.OutputText != "" {
		renderOutputPanel(sb, resolved, cfg)
	}

	if cfg.CopyButton {
		renderCopyButton(sb, resolved.RawCode, cfg)
	}
}

func blockID(code string) string {
	h := fnv.New32a()
	h.Write([]byte(code))
	return fmt.Sprintf("%08x", h.Sum32())
}
