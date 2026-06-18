package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/color"
	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/marker"
)

type MergedToken = config.MergedToken
type TokenLine = config.TokenLine

type lineContext struct {
	resolved         *config.ResolvedBlock
	cfg              *config.Config
	dualTheme        bool
	resolvedMarkers  map[int]marker.ResolvedLine
	focusSet         map[int]bool
	hasFocus         bool
	collapseRangeMap map[int]int
	thresholdVisible map[int]bool
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

		if resolved.CollapseThreshold && lctx.thresholdVisible != nil {
			if !lctx.thresholdVisible[lineNum] {
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

func resolveMarkerClasses(lineNum int, lctx *lineContext) (extraClasses string, markerType *config.MarkerType, labelAttr string) {
	if entry, ok := lctx.resolvedMarkers[lineNum]; ok && entry.HasMark {
		extraClasses = " highlight"
		mt := entry.Type
		markerType = &mt
		switch entry.Type {
		case config.MarkerMark:
			extraClasses += " mark"
		case config.MarkerDel:
			extraClasses += " del"
		case config.MarkerIns:
			extraClasses += " ins"
		}
		if entry.Label != "" {
			extraClasses += " tm-label"
			labelAttr = fmt.Sprintf(" data-label=\"%s\"", html.EscapeString(entry.Label))
		}
	}
	if lctx.hasFocus && lctx.focusSet[lineNum] {
		extraClasses += " focused"
	}
	return
}

func renderLine(sb *strings.Builder, line TokenLine, lineNum int, lctx *lineContext) {
	classes := "kz-line"
	extraClasses, markerType, labelAttr := resolveMarkerClasses(lineNum, lctx)
	classes += extraClasses

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
		sb.WriteString(fmt.Sprintf("<div class=\"kz-gutter\"><div class=\"kz-ln\" aria-hidden=\"true\">%d</div></div>", lineNum))
	}
	sb.WriteString(fmt.Sprintf("<div class=\"kz-code\"%s%s>", labelAttr, indentAttr))
	if indentWS != "" {
		sb.WriteString(fmt.Sprintf("<span class=\"indent\">%s</span>", indentWS))
	}
	lineIdx := lineNum - lctx.resolved.StartLineNumber
	hasInlineMarkers := len(lctx.resolved.InlineMarkers) > 0
	hasLinks := lineIdx >= 0 && lineIdx < len(lctx.resolved.Links) && len(lctx.resolved.Links[lineIdx]) > 0

	if hasInlineMarkers || hasLinks {
		var annotated []marker.TokenWithSegments
		switch {
		case hasInlineMarkers && hasLinks:
			annotated = marker.ProcessInlineMarkersAndLinks(tokens, lctx.resolved.InlineMarkers, lctx.resolved.Links[lineIdx])
		case hasLinks:
			annotated = marker.ProcessLinks(tokens, lctx.resolved.Links[lineIdx])
		default:
			annotated = marker.ProcessInlineMarkers(tokens, lctx.resolved.InlineMarkers)
		}
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
		seg := at.Segments[0]
		isLinkOnly := seg.Marker.Link != "" && seg.Marker.Type == config.MarkerNone
		if isLinkOnly {
			sb.WriteString(fmt.Sprintf("<a class=\"kz-link\" href=\"%s\" rel=\"noopener noreferrer\">", html.EscapeString(seg.Marker.Link)))
			style := buildTokenStyle(at.Token, lctx, markerType)
			if style != "" {
				sb.WriteString(fmt.Sprintf("<span style=\"%s\">", style))
			} else {
				sb.WriteString("<span>")
			}
			sb.WriteString(html.EscapeString(seg.Content))
			sb.WriteString("</span></a>")
			return
		}
		elem := markerElement(seg.Marker.Type)
		var classes []string
		if seg.Marker.OpenStart {
			classes = append(classes, "open-start")
		}
		if seg.Marker.OpenEnd {
			classes = append(classes, "open-end")
		}
		if seg.Marker.Link != "" {
			sb.WriteString(fmt.Sprintf("<a class=\"kz-link\" href=\"%s\" rel=\"noopener noreferrer\">", html.EscapeString(seg.Marker.Link)))
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
		if seg.Marker.Link != "" {
			sb.WriteString("</a>")
		}
		return
	}

	style := buildTokenStyle(at.Token, lctx, markerType)
	if style != "" {
		sb.WriteString(fmt.Sprintf("<span style=\"%s\">", style))
	} else {
		sb.WriteString("<span>")
	}
	for _, seg := range at.Segments {
		if seg.Marker != nil {
			if seg.Marker.Link != "" && seg.Marker.Type == config.MarkerNone {
				sb.WriteString(fmt.Sprintf("<a class=\"kz-link\" href=\"%s\" rel=\"noopener noreferrer\">", html.EscapeString(seg.Marker.Link)))
				sb.WriteString(html.EscapeString(seg.Content))
				sb.WriteString("</a>")
			} else {
				elem := markerElement(seg.Marker.Type)
				if seg.Marker.Link != "" {
					sb.WriteString(fmt.Sprintf("<a class=\"kz-link\" href=\"%s\" rel=\"noopener noreferrer\">", html.EscapeString(seg.Marker.Link)))
				}
				sb.WriteString(fmt.Sprintf("<%s>", elem))
				sb.WriteString(html.EscapeString(seg.Content))
				sb.WriteString(fmt.Sprintf("</%s>", elem))
				if seg.Marker.Link != "" {
					sb.WriteString("</a>")
				}
			}
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

	if tok.FontStyle&config.FontStyleItalic != 0 {
		parts = append(parts, "--sfs:italic")
	}
	if tok.FontStyle&config.FontStyleBold != 0 {
		parts = append(parts, "--sfw:bold")
	}
	if tok.FontStyle&(config.FontStyleUnderline|config.FontStyleStrikethrough) != 0 {
		var decs []string
		if tok.FontStyle&config.FontStyleUnderline != 0 {
			decs = append(decs, "underline")
		}
		if tok.FontStyle&config.FontStyleStrikethrough != 0 {
			decs = append(decs, "line-through")
		}
		parts = append(parts, fmt.Sprintf("--std:%s", strings.Join(decs, " ")))
	}

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
