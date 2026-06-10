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

const copySVG = `<svg class="kz-copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>`

const fullscreenSVG = `<svg class="kz-fs-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/></svg>`

// RenderBlock produces the full HTML for a code block.
func RenderBlock(lines []TokenLine, resolved *config.ResolvedBlock, cfg *config.Config) string {
	var sb strings.Builder
	dualTheme := hasDualTheme(lines)

	wrapperClass := "kazari-code"
	if resolved.CollapseThreshold && (resolved.CollapseConfig == nil || resolved.CollapseConfig.DefaultCollapsed) {
		wrapperClass += " kz-collapsed"
	}
	sb.WriteString(fmt.Sprintf("<div class=\"%s\">\n", wrapperClass))

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
	renderPreCode(sb, lines, resolved, cfg, dualTheme)

	if resolved.CollapseThreshold {
		renderThresholdOverlay(sb, resolved)
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
		sb.WriteString("<span class=\"kz-terminal-dots\"><span></span><span></span><span></span></span>")
	}
	if resolved.Title != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"kz-title\">%s</span>", html.EscapeString(resolved.Title)))
	}
	if cfg.CopyButton {
		sb.WriteString("<div class=\"kz-terminal-actions\">")
		renderCopyButton(sb, resolved.RawCode)
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")

	renderPreCode(sb, lines, resolved, cfg, dualTheme)

	if resolved.CollapseThreshold {
		renderThresholdOverlay(sb, resolved)
	}

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
	renderPreCode(sb, lines, resolved, cfg, dualTheme)

	if resolved.CollapseThreshold {
		renderThresholdOverlay(sb, resolved)
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

	if resolved.LineNumbers {
		endNum := resolved.StartLineNumber + len(lines) - 1
		maxDigits := max(digitCount(resolved.StartLineNumber), digitCount(endNum))
		if maxDigits > 2 {
			sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\">", html.EscapeString(resolved.Lang)))
			renderCodeOpen(sb, lctx, maxDigits)
		} else {
			sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\">", html.EscapeString(resolved.Lang)))
			renderCodeOpen(sb, lctx, 0)
		}
	} else {
		sb.WriteString(fmt.Sprintf("<pre data-language=\"%s\">", html.EscapeString(resolved.Lang)))
		renderCodeOpen(sb, lctx, 0)
	}

	for i, line := range lines {
		lineNum := resolved.StartLineNumber + i

		// Range-based collapse: wrap ranges in <details> or collapsible div
		if rangeIdx, inRange := lctx.inCollapseRange(lineNum); inRange {
			cr := resolved.CollapseRanges[rangeIdx]
			if lineNum == cr.Start {
				renderCollapseRangeOpen(sb, resolved, cr)
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

	sb.WriteString(fmt.Sprintf("<div class=\"%s\">", classes))
	if lctx.resolved.LineNumbers {
		sb.WriteString(fmt.Sprintf("<div class=\"gutter\"><div class=\"ln\" aria-hidden=\"true\">%d</div></div>", lineNum))
	}
	sb.WriteString(fmt.Sprintf("<div class=\"code\"%s>", labelAttr))
	if len(lctx.resolved.InlineMarkers) > 0 {
		annotated := marker.ProcessInlineMarkers(line.Tokens, lctx.resolved.InlineMarkers)
		for _, at := range annotated {
			if at.Token.Content == "" {
				continue
			}
			renderAnnotatedToken(sb, at, lctx, markerType)
		}
	} else {
		for _, tok := range line.Tokens {
			if tok.Content == "" {
				continue
			}
			renderToken(sb, tok, lctx, markerType)
		}
	}
	sb.WriteString("</div></div>")
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

	if markerType != nil && lctx.cfg.MinContrast > 0 {
		if lightColor != "" && lctx.cfg.LightMarkerBGs != nil {
			lightColor = adjustContrast(lightColor, lctx.cfg.LightMarkerBGs.BG(*markerType), lctx.cfg.MinContrast, lctx.contrastCache)
		}
		if darkColor != "" && lctx.cfg.DarkMarkerBGs != nil {
			darkColor = adjustContrast(darkColor, lctx.cfg.DarkMarkerBGs.BG(*markerType), lctx.cfg.MinContrast, lctx.contrastCache)
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

	return strings.Join(parts, ";")
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

func renderSummaryLine(sb *strings.Builder, resolved *config.ResolvedBlock, cr config.CollapseRange) {
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
	sb.WriteString("<span class=\"expand\"></span>")
	sb.WriteString("<span class=\"collapse\"></span>")
	sb.WriteString(fmt.Sprintf("<span class=\"text\">%s</span>", collapsible.SummaryText(cr.LineCount)))
	sb.WriteString("</div>")
	sb.WriteString("</div>")
	sb.WriteString("</summary>")
}

func renderCollapseRangeOpen(sb *strings.Builder, resolved *config.ResolvedBlock, cr config.CollapseRange) {
	switch cr.Style {
	case config.CollapseCollapsibleStart:
		sb.WriteString("<div class=\"kz-section collapsible-start\">")
		sb.WriteString("<details>")
		renderSummaryLine(sb, resolved, cr)
		sb.WriteString("</details>")
		sb.WriteString("<div class=\"content-lines\">")

	case config.CollapseCollapsibleEnd:
		sb.WriteString("<div class=\"kz-section collapsible-end\">")
		sb.WriteString("<details>")
		renderSummaryLine(sb, resolved, cr)
		sb.WriteString("</details>")
		sb.WriteString("<div class=\"content-lines\">")

	default: // github
		sb.WriteString("<details class=\"kz-section\">")
		renderSummaryLine(sb, resolved, cr)
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
	sb.WriteString("<div class=\"code\"><span class=\"kz-gap-indicator\">⋮</span></div>")
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

func renderThresholdOverlay(sb *strings.Builder, resolved *config.ResolvedBlock) {
	expandText := "Show more"
	collapseText := "Show less"
	expandedAnnouncement := "Code block expanded"
	collapsedAnnouncement := "Code block collapsed"
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

	// Badge fallback: append highlighted line count when markers are beyond 2× cap
	if resolved.CollapseBeyondCap > 0 {
		expandText = fmt.Sprintf("%s (+%d highlighted)", expandText, resolved.CollapseBeyondCap)
	}

	sb.WriteString("<div class=\"kz-collapse-gradient\"></div>")
	sb.WriteString(fmt.Sprintf(
		"<button class=\"kz-collapse-btn\" aria-expanded=\"false\" data-expand=\"%s\" data-collapse=\"%s\" data-expanded-msg=\"%s\" data-collapsed-msg=\"%s\">%s</button>",
		html.EscapeString(expandText),
		html.EscapeString(collapseText),
		html.EscapeString(expandedAnnouncement),
		html.EscapeString(collapsedAnnouncement),
		html.EscapeString(expandText),
	))
	sb.WriteString("<div class=\"kz-sr-announce\" aria-live=\"polite\"></div>")
}
