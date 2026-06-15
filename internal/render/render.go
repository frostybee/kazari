package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/collapsible"
	"github.com/frostybee/kazari/internal/color"
	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/marker"
)

// Type aliases for backward compatibility within the package.
type MergedToken = config.MergedToken
type TokenLine = config.TokenLine

const (
	fontStyleItalic        = 1
	fontStyleBold          = 2
	fontStyleUnderline     = 4
	fontStyleStrikethrough = 8
)

const copySVG = `<svg class="kz-copy-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>`

const fullscreenSVG = `<svg class="kz-fs-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/></svg>`

const fullscreenExitSVG = `<svg class="kz-fs-exit-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 4v5H4m5 0L4 4M15 4v5h5m-5 0l5-5M9 20v-5H4m5 0l-5 5M15 20v-5h5m-5 0l5 5"/></svg>`

const chevronSVG = `<svg class="kz-collapse-toggle-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>`

const fontIncreaseSVG = `<svg class="kz-font-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6v12m-6-6h12"/></svg>`

const fontDecreaseSVG = `<svg class="kz-font-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 12h12"/></svg>`

const wrapSVG = `<svg class="kz-wrap-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6h18M3 12h15a3 3 0 110 6h-4m0 0l2-2m-2 2l2 2"/></svg>`

const wrapOffSVG = `<svg class="kz-wrap-off-icon" aria-hidden="true" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6h18M3 12h18M3 18h18"/></svg>`

// RenderBlock produces the full HTML for a code block.
func RenderBlock(lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config) string {
	var sb strings.Builder
	dualTheme := hasDualTheme(lines)

	wrapperClass := "kazari-code"
	if resolved.ThemeOverrideStyle != "" {
		wrapperClass += " kz-themed"
	}
	if resolved.CollapseThreshold && (resolved.CollapseConfig == nil || resolved.CollapseConfig.DefaultCollapsed) {
		wrapperClass += " kz-collapsed"
	}
	attrs := fmt.Sprintf(" class=\"%s\"", wrapperClass)
	if cfg.DataLineCount {
		attrs += fmt.Sprintf(" data-lines=\"%d\"", len(lines))
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
	renderFullscreenHint(sb, cfg)
	renderCollapseContentStart(sb, resolved)
	renderPreCode(sb, lines, resolved, cfg, dualTheme)
	renderCollapseContentEnd(sb, resolved)

	if resolved.CollapseThreshold {
		renderCollapseBar(sb, resolved, cfg)
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
	if cfg.CopyButton || cfg.WrapButton || cfg.FullscreenButton {
		sb.WriteString("<div class=\"kz-terminal-actions\">")
		if cfg.CopyButton {
			renderCopyButton(sb, resolved.RawCode, cfg)
		}
		if cfg.WrapButton {
			renderWrapButton(sb, resolved, cfg)
		}
		if cfg.FullscreenButton {
			renderFontControls(sb, cfg)
			renderFullscreenButton(sb, cfg)
		}
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")

	renderFullscreenHint(sb, cfg)
	renderCollapseContentStart(sb, resolved)
	renderPreCode(sb, lines, resolved, cfg, dualTheme)
	renderCollapseContentEnd(sb, resolved)

	if resolved.CollapseThreshold {
		renderCollapseBar(sb, resolved, cfg)
	}

	sb.WriteString("</figure>\n")
}

func renderToolbar(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	sb.WriteString("<div class=\"kz-toolbar\">")

	// Left section: language label, then separator + title (if present)
	sb.WriteString("<div class=\"kz-toolbar-left\">")
	if cfg.LanguageBadge && resolved.Lang != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-lang\">%s</span>", html.EscapeString(displayLang(resolved.Lang))))
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

	// Right section: action buttons only
	sb.WriteString("<div class=\"kz-toolbar-right\">")
	if cfg.CopyButton {
		renderCopyButton(sb, resolved.RawCode, cfg)
	}
	if cfg.WrapButton {
		renderWrapButton(sb, resolved, cfg)
	}
	if cfg.FullscreenButton {
		renderFontControls(sb, cfg)
		renderFullscreenButton(sb, cfg)
	}
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

func renderFullscreenHint(sb *strings.Builder, cfg *config.Config) {
}

func renderNoFrame(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	renderCollapseContentStart(sb, resolved)
	renderPreCode(sb, lines, resolved, cfg, dualTheme)
	renderCollapseContentEnd(sb, resolved)

	if resolved.CollapseThreshold {
		renderCollapseBar(sb, resolved, cfg)
	}

	if cfg.CopyButton {
		renderCopyButton(sb, resolved.RawCode, cfg)
	}
}

type lineContext struct {
	resolved         *config.ResolvedBlock
	cfg              *config.Config
	dualTheme        bool
	resolvedMarkers  map[int]marker.ResolvedLine
	focusSet         map[int]bool
	hasFocus         bool
	collapseRangeMap map[int]int  // lineNum -> range index
	thresholdVisible map[int]bool // lines visible in threshold preview
	contrastCache    map[string]string
}

func renderPreCode(sb *strings.Builder, lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config, dualTheme bool) {
	lctx := &lineContext{
		resolved:         resolved,
		cfg:              cfg,
		dualTheme:        dualTheme,
		resolvedMarkers:  marker.ResolveLineMarkers(resolved.LineMarkers),
		focusSet:         marker.ResolveFocusSet(resolved.FocusLines),
		hasFocus:         len(resolved.FocusLines) > 0,
		collapseRangeMap: buildCollapseRangeMap(resolved.CollapseRanges),
		thresholdVisible: buildThresholdVisibleSet(resolved.CollapseSegments),
		contrastCache:    make(map[string]string),
	}

	preClass := ""
	if resolved.Wrap {
		preClass = " class=\"wrap\""
	}
	sb.WriteString(fmt.Sprintf("<pre%s data-language=\"%s\">", preClass, html.EscapeString(resolved.Lang)))

	lnWidth := 0
	if resolved.LineNumbers {
		endNum := resolved.StartLineNumber + len(lines) - 1
		maxDigits := max(digitCount(resolved.StartLineNumber), digitCount(endNum))
		if maxDigits > 2 {
			lnWidth = maxDigits
		}
	}
	renderCodeOpen(sb, lctx, lnWidth)

	for i, line := range lines {
		lineNum := resolved.StartLineNumber + i

		// Range-based collapse: wrap ranges in <details> or collapsible div
		if rangeIdx, inRange := lctx.inCollapseRange(lineNum); inRange {
			cr := resolved.CollapseRanges[rangeIdx]
			if lineNum == cr.Start {
				renderCollapseRangeOpen(sb, resolved, cr, cfg)
			}
			renderLine(sb, line, lineNum, lctx)
			if lineNum == cr.End {
				renderCollapseRangeClose(sb, cr)
			}
			continue
		}

		// Threshold-based collapse: segment-based visibility with gap indicators
		if resolved.CollapseThreshold && lctx.thresholdVisible != nil {
			if !lctx.thresholdVisible[lineNum] {
				// Emit gap indicator between non-contiguous segments (not at trailing edge)
				if len(resolved.CollapseSegments) > 1 {
					prevVisible := lctx.thresholdVisible[lineNum-1] || lctx.inCollapseRangeEnd(lineNum-1)
					nextVisibleExists := lctx.hasVisibleLineAfter(lineNum)
					if prevVisible && nextVisibleExists {
						renderGapIndicator(sb, resolved)
					}
				}
				renderHiddenLine(sb, line, lineNum, lctx)
				continue
			}
		}

		renderLine(sb, line, lineNum, lctx)
	}

	sb.WriteString("</code></pre>")
}

func renderCodeOpen(sb *strings.Builder, lctx *lineContext, lnWidth int) {
	classes := ""
	if lctx.hasFocus {
		classes = " class=\"has-focus\""
	}
	if lnWidth > 0 {
		sb.WriteString(fmt.Sprintf("<code%s style=\"--kz-ln-width:%dch\">", classes, lnWidth))
	} else {
		sb.WriteString(fmt.Sprintf("<code%s>", classes))
	}
}

func encodeForDataCode(code string) string {
	return strings.ReplaceAll(code, "\n", "\x7f")
}

func fileExt(title string) string {
	idx := strings.LastIndex(title, ".")
	if idx < 0 || idx == len(title)-1 {
		return ""
	}
	return title[idx+1:]
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

func renderLine(sb *strings.Builder, line TokenLine, lineNum int, lctx *lineContext) {
	classes := "kz-line"
	labelAttr := ""
	var markerType *config.MarkerType

	if entry, ok := lctx.resolvedMarkers[lineNum]; ok && entry.HasMark {
		classes += " highlight"
		mt := entry.Type
		markerType = &mt
		switch entry.Type {
		case config.MarkerMark:
			classes += " mark"
		case config.MarkerDel:
			classes += " del"
		case config.MarkerIns:
			classes += " ins"
		}
		if entry.Label != "" {
			classes += " tm-label"
			labelAttr = fmt.Sprintf(" data-label=\"%s\"", html.EscapeString(entry.Label))
		}
	}

	if lctx.hasFocus && lctx.focusSet[lineNum] {
		classes += " focused"
	}

	tokens := line.Tokens
	indentAttr := ""
	indentWS := ""
	if lctx.resolved.Wrap {
		ws, rest := splitLeadingWhitespace(line.Tokens)
		indent := lctx.resolved.HangingIndent
		if lctx.resolved.PreserveIndent {
			indent += len(ws)
		}
		if indent > 0 {
			indentAttr = fmt.Sprintf(" style=\"--kz-indent:%dch\"", indent)
			indentWS = ws
			tokens = rest
		}
	}

	sb.WriteString(fmt.Sprintf("<div class=\"%s\">", classes))
	if lctx.resolved.LineNumbers {
		sb.WriteString(fmt.Sprintf("<div class=\"gutter\"><div class=\"ln\" aria-hidden=\"true\">%d</div></div>", lineNum))
	}
	sb.WriteString(fmt.Sprintf("<div class=\"code\"%s%s>", labelAttr, indentAttr))
	if indentWS != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"indent\">%s</span>", indentWS))
	}
	if len(lctx.resolved.InlineMarkers) > 0 {
		annotated := marker.ProcessInlineMarkers(tokens, lctx.resolved.InlineMarkers)
		for _, at := range annotated {
			if at.Token.Content == "" {
				continue
			}
			renderAnnotatedToken(sb, at, lctx, markerType)
		}
	} else {
		for _, tok := range tokens {
			if tok.Content == "" {
				continue
			}
			renderToken(sb, tok, lctx, markerType)
		}
	}
	sb.WriteString("</div></div>")
}

// splitLeadingWhitespace separates the leading whitespace run from a token
// stream. The returned rest slice is a copy with the first token trimmed,
// so callers can render it without mutating the input line.
func splitLeadingWhitespace(tokens []MergedToken) (string, []MergedToken) {
	var ws strings.Builder
	for i, tok := range tokens {
		trimmed := strings.TrimLeft(tok.Content, " \t")
		ws.WriteString(tok.Content[:len(tok.Content)-len(trimmed)])
		if trimmed != "" {
			rest := make([]MergedToken, len(tokens)-i)
			copy(rest, tokens[i:])
			rest[0].Content = trimmed
			return ws.String(), rest
		}
	}
	return ws.String(), nil
}

func markerElement(mtype config.MarkerType) string {
	switch mtype {
	case config.MarkerIns:
		return "ins"
	case config.MarkerDel:
		return "del"
	default:
		return "mark"
	}
}

func renderAnnotatedToken(sb *strings.Builder, at marker.TokenWithSegments, lctx *lineContext, markerType *config.MarkerType) {
	hasInlineMarker := false
	for _, seg := range at.Segments {
		if seg.Marker != nil {
			hasInlineMarker = true
			break
		}
	}

	if !hasInlineMarker {
		renderToken(sb, at.Token, lctx, markerType)
		return
	}

	allOneSegmentSpanning := len(at.Segments) == 1 && at.Segments[0].Marker != nil &&
		(at.Segments[0].Marker.OpenStart || at.Segments[0].Marker.OpenEnd)

	if allOneSegmentSpanning {
		// Multi-token span: mark wraps the span.
		seg := at.Segments[0]
		elem := markerElement(seg.Marker.Type)
		var classes []string
		if seg.Marker.OpenStart {
			classes = append(classes, "open-start")
		}
		if seg.Marker.OpenEnd {
			classes = append(classes, "open-end")
		}
		sb.WriteString(fmt.Sprintf("<%s class=\"%s\">", elem, strings.Join(classes, " ")))
		style := buildTokenStyle(at.Token, lctx, markerType)
		if style != "" {
			sb.WriteString(fmt.Sprintf("<span style=\"%s\">", style))
		} else {
			sb.WriteString("<span>")
		}
		sb.WriteString(html.EscapeString(seg.Content))
		sb.WriteString("</span>")
		sb.WriteString(fmt.Sprintf("</%s>", elem))
		return
	}

	// Partial match or standalone: mark goes inside the span.
	style := buildTokenStyle(at.Token, lctx, markerType)
	if style != "" {
		sb.WriteString(fmt.Sprintf("<span style=\"%s\">", style))
	} else {
		sb.WriteString("<span>")
	}
	for _, seg := range at.Segments {
		if seg.Marker != nil {
			elem := markerElement(seg.Marker.Type)
			sb.WriteString(fmt.Sprintf("<%s>", elem))
			sb.WriteString(html.EscapeString(seg.Content))
			sb.WriteString(fmt.Sprintf("</%s>", elem))
		} else {
			sb.WriteString(html.EscapeString(seg.Content))
		}
	}
	sb.WriteString("</span>")
}

func renderToken(sb *strings.Builder, tok MergedToken, lctx *lineContext, markerType *config.MarkerType) {
	style := buildTokenStyle(tok, lctx, markerType)
	if style != "" {
		sb.WriteString(fmt.Sprintf("<span style=\"%s\">", style))
	} else {
		sb.WriteString("<span>")
	}
	sb.WriteString(html.EscapeString(tok.Content))
	sb.WriteString("</span>")
}

func buildTokenStyle(tok MergedToken, lctx *lineContext, markerType *config.MarkerType) string {
	var parts []string

	lightColor := tok.LightColor
	darkColor := tok.DarkColor

	if lctx.cfg.MinContrast > 0 {
		if markerType != nil {
			// Per-block theme overrides carry their own marker backgrounds so
			// contrast is computed against the override canvas, not the page's.
			lightBGs := lctx.cfg.LightMarkerBGs
			if lctx.resolved.LightMarkerBGs != nil {
				lightBGs = lctx.resolved.LightMarkerBGs
			}
			darkBGs := lctx.cfg.DarkMarkerBGs
			if lctx.resolved.DarkMarkerBGs != nil {
				darkBGs = lctx.resolved.DarkMarkerBGs
			}
			if lightColor != "" && lightBGs != nil {
				lightColor = adjustContrast(lightColor, lightBGs.BG(*markerType), lctx.cfg.MinContrast, lctx.contrastCache)
			}
			if darkColor != "" && darkBGs != nil {
				darkColor = adjustContrast(darkColor, darkBGs.BG(*markerType), lctx.cfg.MinContrast, lctx.contrastCache)
			}
		} else {
			lightBG := lctx.cfg.LightEditorBG
			if lctx.resolved.LightEditorBG != "" {
				lightBG = lctx.resolved.LightEditorBG
			}
			darkBG := lctx.cfg.DarkEditorBG
			if lctx.resolved.DarkEditorBG != "" {
				darkBG = lctx.resolved.DarkEditorBG
			}
			if lightColor != "" && lightBG != "" {
				lightColor = adjustContrast(lightColor, lightBG, lctx.cfg.MinContrast, lctx.contrastCache)
			}
			if darkColor != "" && darkBG != "" {
				darkColor = adjustContrast(darkColor, darkBG, lctx.cfg.MinContrast, lctx.contrastCache)
			}
		}
	}

	if lightColor != "" {
		parts = append(parts, fmt.Sprintf("--sl:%s", lightColor))
	}
	if lctx.dualTheme && darkColor != "" {
		parts = append(parts, fmt.Sprintf("--sd:%s", darkColor))
	}
	if tok.LightBG != "" {
		parts = append(parts, fmt.Sprintf("--slbg:%s", tok.LightBG))
	}
	if lctx.dualTheme && tok.DarkBG != "" {
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

	// Escape the assembled style so values from custom Highlighter
	// implementations cannot break out of the style attribute.
	return html.EscapeString(strings.Join(parts, ";"))
}

func adjustContrast(tokenColor, effectiveBG string, minContrast float64, cache map[string]string) string {
	key := tokenColor + "|" + effectiveBG
	if adjusted, ok := cache[key]; ok {
		return adjusted
	}
	adjusted := color.EnsureContrastOnBackground(tokenColor, effectiveBG, minContrast)
	cache[key] = adjusted
	return adjusted
}

func digitCount(n int) int {
	extra := 0
	if n < 0 {
		extra = 1
		n = -n
	}
	if n == 0 {
		return 1 + extra
	}
	count := 0
	for n > 0 {
		count++
		n /= 10
	}
	return count + extra
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

// --- Collapsible rendering helpers ---

func buildCollapseRangeMap(ranges []config.CollapseRange) map[int]int {
	if len(ranges) == 0 {
		return nil
	}
	m := make(map[int]int)
	for i, cr := range ranges {
		for line := cr.Start; line <= cr.End; line++ {
			m[line] = i
		}
	}
	return m
}

func (lctx *lineContext) inCollapseRange(lineNum int) (int, bool) {
	if lctx.collapseRangeMap == nil {
		return 0, false
	}
	idx, ok := lctx.collapseRangeMap[lineNum]
	return idx, ok
}

func (lctx *lineContext) hasVisibleLineAfter(lineNum int) bool {
	for _, seg := range lctx.resolved.CollapseSegments {
		if seg.Start > lineNum {
			return true
		}
	}
	return false
}

func (lctx *lineContext) inCollapseRangeEnd(lineNum int) bool {
	idx, ok := lctx.collapseRangeMap[lineNum]
	if !ok {
		return false
	}
	return lineNum == lctx.resolved.CollapseRanges[idx].End
}

func buildThresholdVisibleSet(segments []config.PreviewSegment) map[int]bool {
	if len(segments) == 0 {
		return nil
	}
	m := make(map[int]bool)
	for _, seg := range segments {
		for line := seg.Start; line <= seg.End; line++ {
			m[line] = true
		}
	}
	return m
}

func renderSummaryLine(sb *strings.Builder, resolved *config.ResolvedBlock, cr config.CollapseRange, cfg *config.Config) {
	sb.WriteString("<summary>")
	sb.WriteString("<div class=\"kz-line\">")
	if resolved.LineNumbers {
		sb.WriteString("<div class=\"gutter\"><div class=\"ln\"></div></div>")
	}

	indentStyle := ""
	if cr.MinIndent > 0 && resolved.CollapseConfig != nil && resolved.CollapseConfig.PreserveIndent {
		indentStyle = fmt.Sprintf(" style=\"--kz-indent:%dch\"", cr.MinIndent)
	}

	sb.WriteString(fmt.Sprintf("<div class=\"code\"%s>", indentStyle))
	sb.WriteString("<span class=\"expand\" aria-hidden=\"true\"></span>")
	sb.WriteString("<span class=\"collapse\" aria-hidden=\"true\"></span>")
	sb.WriteString(fmt.Sprintf("<span class=\"text\">%s</span>", collapsible.SummaryText(cr.LineCount, cfg.UIStrings)))
	sb.WriteString("</div>")
	sb.WriteString("</div>")
	sb.WriteString("</summary>")
}

func renderCollapseRangeOpen(sb *strings.Builder, resolved *config.ResolvedBlock, cr config.CollapseRange, cfg *config.Config) {
	switch cr.Style {
	case config.CollapseCollapsibleStart:
		sb.WriteString("<div class=\"kz-section collapsible-start\">")
		sb.WriteString("<details>")
		renderSummaryLine(sb, resolved, cr, cfg)
		sb.WriteString("</details>")
		sb.WriteString("<div class=\"content-lines\">")

	case config.CollapseCollapsibleEnd:
		sb.WriteString("<div class=\"kz-section collapsible-end\">")
		sb.WriteString("<details>")
		renderSummaryLine(sb, resolved, cr, cfg)
		sb.WriteString("</details>")
		sb.WriteString("<div class=\"content-lines\">")

	default: // github
		sb.WriteString("<details class=\"kz-section\">")
		renderSummaryLine(sb, resolved, cr, cfg)
	}
}

func renderCollapseRangeClose(sb *strings.Builder, cr config.CollapseRange) {
	switch cr.Style {
	case config.CollapseCollapsibleStart, config.CollapseCollapsibleEnd:
		sb.WriteString("</div></div>") // close content-lines + kz-section wrapper
	default: // github
		sb.WriteString("</details>")
	}
}

func renderGapIndicator(sb *strings.Builder, resolved *config.ResolvedBlock) {
	sb.WriteString("<div class=\"kz-line kz-gap\">")
	if resolved.LineNumbers {
		sb.WriteString("<div class=\"gutter\"><div class=\"ln\"></div></div>")
	}
	sb.WriteString("<div class=\"code\"><span class=\"kz-gap-indicator\" aria-hidden=\"true\">⋮</span><span class=\"sr-only\">Lines hidden</span></div>")
	sb.WriteString("</div>")
}

func renderHiddenLine(sb *strings.Builder, line TokenLine, lineNum int, lctx *lineContext) {
	classes := "kz-line kz-hidden"
	var markerType *config.MarkerType

	if entry, ok := lctx.resolvedMarkers[lineNum]; ok && entry.HasMark {
		classes += " highlight"
		mt := entry.Type
		markerType = &mt
		switch entry.Type {
		case config.MarkerMark:
			classes += " mark"
		case config.MarkerDel:
			classes += " del"
		case config.MarkerIns:
			classes += " ins"
		}
	}

	if lctx.hasFocus && lctx.focusSet[lineNum] {
		classes += " focused"
	}

	sb.WriteString(fmt.Sprintf("<div class=\"%s\">", classes))
	if lctx.resolved.LineNumbers {
		sb.WriteString(fmt.Sprintf("<div class=\"gutter\"><div class=\"ln\" aria-hidden=\"true\">%d</div></div>", lineNum))
	}
	sb.WriteString("<div class=\"code\">")
	for _, tok := range line.Tokens {
		if tok.Content == "" {
			continue
		}
		renderToken(sb, tok, lctx, markerType)
	}
	sb.WriteString("</div></div>")
}

func renderCollapseContentStart(sb *strings.Builder, resolved *config.ResolvedBlock) {
	if resolved.CollapseThreshold {
		sb.WriteString("<div class=\"kz-collapse-content\">")
	}
}

func renderCollapseContentEnd(sb *strings.Builder, resolved *config.ResolvedBlock) {
	if resolved.CollapseThreshold {
		sb.WriteString("<div class=\"kz-collapse-gradient\"></div>")
		sb.WriteString("</div>")
	}
}

func renderCollapseBar(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	expandText := cfg.UIStrings.ExpandButtonText
	collapseText := cfg.UIStrings.CollapseButtonText
	expandedAnnouncement := cfg.UIStrings.ExpandedAnnouncement
	collapsedAnnouncement := cfg.UIStrings.CollapsedAnnouncement
	if resolved.CollapseConfig != nil {
		if resolved.CollapseConfig.ExpandButtonText != "" {
			expandText = resolved.CollapseConfig.ExpandButtonText
		}
		if resolved.CollapseConfig.CollapseButtonText != "" {
			collapseText = resolved.CollapseConfig.CollapseButtonText
		}
		if resolved.CollapseConfig.ExpandedAnnouncement != "" {
			expandedAnnouncement = resolved.CollapseConfig.ExpandedAnnouncement
		}
		if resolved.CollapseConfig.CollapsedAnnouncement != "" {
			collapsedAnnouncement = resolved.CollapseConfig.CollapsedAnnouncement
		}
	}

	if resolved.CollapseBeyondCap > 0 {
		expandText = fmt.Sprintf("%s (+%d highlighted)", expandText, resolved.CollapseBeyondCap)
	}

	sb.WriteString("<div class=\"kz-collapse-bar\">")
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-collapse-btn\" aria-expanded=\"false\" data-expand=\"%s\" data-collapse=\"%s\" data-expanded-msg=\"%s\" data-collapsed-msg=\"%s\">%s</button>",
		html.EscapeString(expandText),
		html.EscapeString(collapseText),
		html.EscapeString(expandedAnnouncement),
		html.EscapeString(collapsedAnnouncement),
		html.EscapeString(expandText),
	))
	sb.WriteString("</div>")
	sb.WriteString("<div class=\"kz-sr-announce\" aria-live=\"polite\"></div>")
}
