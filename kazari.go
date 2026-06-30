// Package kazari renders framed, syntax-highlighted HTML code blocks
// with full CSS customization via custom properties.
package kazari

import (
	"fmt"
	"html"
	"strings"
	"sync"

	"github.com/frostybee/kazari/internal/ansi"
	"github.com/frostybee/kazari/internal/collapsible"
	"github.com/frostybee/kazari/internal/color"
	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/css"
	"github.com/frostybee/kazari/internal/diff"
	"github.com/frostybee/kazari/internal/frame"
	"github.com/frostybee/kazari/internal/js"
	"github.com/frostybee/kazari/internal/link"
	"github.com/frostybee/kazari/internal/locale"
	"github.com/frostybee/kazari/internal/meta"
	"github.com/frostybee/kazari/internal/render"
	"github.com/frostybee/kazari/internal/theme"
)

// Engine is the main entry point for rendering code blocks.
type Engine struct {
	hl          Highlighter
	cfg         *config.Config
	lightColors theme.ThemeColors
	darkColors  theme.ThemeColors

	// Theme pipeline kept for render-time per-block overrides.
	themeAdjustments *ThemeAdjustments
	themeCustomizer  func(string, ThemeInfo) ThemeInfo

	postRenders []func(string, BlockInfo) string

	overrideMu    sync.RWMutex
	overrideCache map[string]overrideEntry
}

// overrideEntry caches the resolved state for one theme= override string.
type overrideEntry struct {
	style        string
	lightEditorBG string
	darkEditorBG  string
	lightBGs     *config.MarkerEffectiveBGs
	darkBGs      *config.MarkerEffectiveBGs
}

// New creates a new Engine with the given options.
func New(opts ...Option) *Engine {
	b := &engineBuilder{cfg: config.DefaultConfig()}

	for _, opt := range opts {
		opt(b)
	}

	e := &Engine{
		hl:               b.hl,
		cfg:              b.cfg,
		themeAdjustments: b.themeAdjustments,
		themeCustomizer:  b.themeCustomizer,
		postRenders:      b.postRenders,
		overrideCache:    make(map[string]overrideEntry),
	}

	// Extract theme colors at construction time for CSS generation.
	if e.hl != nil {
		if colors, ok := extractThemeColors(e.hl, e.cfg.LightTheme, b.themeAdjustments, b.themeCustomizer); ok {
			e.lightColors = colors
		}
		if e.cfg.DarkTheme != "" {
			if colors, ok := extractThemeColors(e.hl, e.cfg.DarkTheme, b.themeAdjustments, b.themeCustomizer); ok {
				e.darkColors = colors
			}
		}
	}

	if e.cfg.MinContrast > 0 {
		if e.lightColors.EditorBG != "" {
			e.cfg.LightEditorBG = e.lightColors.EditorBG
			e.cfg.LightMarkerBGs = computeMarkerBGs(e.lightColors.EditorBG)
		}
		if e.darkColors.EditorBG != "" {
			e.cfg.DarkEditorBG = e.darkColors.EditorBG
			e.cfg.DarkMarkerBGs = computeMarkerBGs(e.darkColors.EditorBG)
		}
	}

	e.cfg.UIStrings = locale.Resolve(e.cfg.Locale, e.cfg.UIStringOverrides)

	return e
}

// extractThemeColors pulls editor colors from a theme, then applies theme
// adjustments and the theme customizer in that order (customizer gets final say).
func extractThemeColors(hl Highlighter, themeName string, adj *ThemeAdjustments, customizer func(string, ThemeInfo) ThemeInfo) (theme.ThemeColors, bool) {
	info, err := hl.GetThemeColors(themeName)
	if err != nil {
		return theme.ThemeColors{}, false
	}
	ti := applyThemeAdjustments(info, adj)
	if customizer != nil {
		ti = customizer(themeName, ti)
	}
	return theme.ThemeColors{
		EditorBG:     ti.BG,
		EditorFG:     ti.FG,
		SelectionBG:  ti.SelectionBG,
		LineNumberFG: ti.LineNumberFG,
		FoldBG:       ti.FoldBG,
	}, true
}

// applyThemeAdjustments tints the selected extracted colors in OKLCH space,
// preserving each color's lightness and alpha.
func applyThemeAdjustments(ti ThemeInfo, adj *ThemeAdjustments) ThemeInfo {
	if adj == nil || (adj.Hue == nil && adj.Chroma == nil) {
		return ti
	}

	tint := func(hex string) string {
		if hex == "" {
			return hex
		}
		l, c, h, err := color.ToOKLCH(hex)
		if err != nil {
			return hex
		}
		if adj.Hue != nil {
			h = *adj.Hue
		}
		if adj.Chroma != nil {
			c = *adj.Chroma
		}
		out := color.FromOKLCH(l, c, h)
		if _, _, _, a, err := color.ParseHex(hex); err == nil && a < 1 {
			out = color.SetAlpha(out, a)
		}
		return out
	}

	targets := adj.Targets
	if targets == 0 {
		targets = AdjustBackgrounds
	}
	if targets&AdjustBackgrounds != 0 {
		ti.BG = tint(ti.BG)
		ti.SelectionBG = tint(ti.SelectionBG)
		ti.FoldBG = tint(ti.FoldBG)
	}
	if targets&AdjustForegrounds != 0 {
		ti.FG = tint(ti.FG)
		ti.LineNumberFG = tint(ti.LineNumberFG)
	}
	return ti
}

// Render renders a code block with structured options.
func (e *Engine) Render(code string, opts Options) (string, error) {
	blockOpts := mapOptionsToBlockOpts(opts)
	lang := e.cfg.ResolveLanguage(opts.Lang)
	resolved := e.cfg.Resolve(lang, blockOpts)
	resolved.LineMarkers = convertLineMarkers(opts.LineMarkers)
	resolved.InlineMarkers = convertInlineMarkers(opts.InlineMarkers)
	resolved.FocusLines = convertRanges(opts.FocusLines)
	resolved.DiffLang = opts.DiffLang

	var spec *config.CollapseSpec
	if opts.Collapse != nil {
		spec = &config.CollapseSpec{
			Enabled:   opts.Collapse.Enabled,
			Disabled:  opts.Collapse.Disabled,
			Ranges:    convertRanges(opts.Collapse.Ranges),
			Style:     convertCollapseStyle(opts.Collapse.Style),
			Threshold: opts.Collapse.Threshold,
		}
	}

	return e.renderResolved(code, resolved, spec, "")
}

// RenderWithMeta parses a meta string and renders the code block.
func (e *Engine) RenderWithMeta(code string, metaStr string) (string, error) {
	parsed := meta.Parse(metaStr)
	lang := e.cfg.ResolveLanguage(parsed.BlockOptions.Lang)
	parsed.BlockOptions.Lang = lang
	resolved := e.cfg.Resolve(lang, &parsed.BlockOptions)
	resolved.LineMarkers = parsed.LineMarkers
	resolved.InlineMarkers = parsed.InlineMarkers
	resolved.FocusLines = parsed.FocusLines
	resolved.DiffLang = parsed.DiffLang
	return e.renderResolved(code, resolved, parsed.Collapse, metaStr)
}

func (e *Engine) renderResolved(code string, resolved *config.ResolvedBlock, collapseSpec *config.CollapseSpec, rawMeta string) (string, error) {
	if e.cfg.MermaidPassThrough && resolved.Lang == "mermaid" {
		return renderMermaidBlock(code), nil
	}

	code = e.preprocess(code, resolved)

	// Collapse resolution must run on the preprocessed code. Preprocessing can
	// remove a filename comment line, which shifts line counts and ranges.
	e.resolveCollapse(code, resolved, collapseSpec)

	if resolved.Theme != "" && e.hl != nil {
		e.applyThemeOverride(resolved)
	}

	lineCount := strings.Count(code, "\n") + 1

	var rendered string
	if resolved.Lang == "ansi" {
		lines := ansi.Parse(code)
		rendered = render.RenderBlock(lines, resolved, e.cfg)
	} else {
		if resolved.Lang == "diff" && resolved.DiffLang != "" {
			stripped, diffMarkers := diff.ProcessDiffBlock(code)
			code = stripped
			resolved.Lang = resolved.DiffLang
			resolved.LineMarkers = append(resolved.LineMarkers, diffMarkers...)
		}

		lines, err := e.tokenize(code, resolved.Lang, resolved.Theme)
		if err != nil {
			return "", err
		}
		rendered = render.RenderBlock(lines, resolved, e.cfg)
	}

	if len(e.postRenders) > 0 {
		info := BlockInfo{
			Lang:      resolved.Lang,
			Title:     resolved.Title,
			Frame:     Frame(resolved.Frame),
			RawCode:   resolved.RawCode,
			LineCount: lineCount,
			Theme:     resolved.Theme,
			Meta:      rawMeta,
		}
		for _, fn := range e.postRenders {
			rendered = fn(rendered, info)
		}
	}

	return rendered, nil
}

func (e *Engine) preprocess(code string, resolved *config.ResolvedBlock) string {
	if e.cfg.FileNameExtraction && resolved.Title == "" {
		title, modified := frame.ExtractFileName(code, resolved.Lang)
		if title != "" {
			resolved.Title = title
			code = modified
		}
	}

	if e.cfg.Links {
		code, resolved.Links = link.ExtractLinks(code)
	}

	if resolved.Frame == config.FrameAuto {
		if e.cfg.FrameDetection {
			resolved.Frame = frame.DetectFrameType(code, resolved.Lang, resolved.Frame)
		} else {
			resolved.Frame = config.FrameCode
		}
	}

	if resolved.WithOutput && e.cfg.OutputPanel {
		code = e.splitOutputSection(code, resolved)
	}

	resolved.RawCode = code
	if resolved.Frame == config.FrameTerminal && e.cfg.TerminalCommentStripping {
		resolved.RawCode = frame.StripTerminalComments(resolved.RawCode)
	}
	return code
}

func (e *Engine) splitOutputSection(code string, resolved *config.ResolvedBlock) string {
	sep := e.cfg.OutputSeparator
	if sep == "" {
		sep = "---output---"
	}
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == sep {
			resolved.OutputText = strings.Join(lines[i+1:], "\n")
			return strings.Join(lines[:i], "\n")
		}
	}
	return code
}

func renderMermaidBlock(code string) string {
	return fmt.Sprintf("<pre class=\"mermaid\">%s</pre>\n", html.EscapeString(code))
}

// Tokenize returns raw tokens for consumers building custom HTML.
func (e *Engine) Tokenize(code string, lang string) ([][]Token, error) {
	if e.hl == nil {
		return nil, fmt.Errorf("kazari: no highlighter configured")
	}
	lang = e.cfg.ResolveLanguage(lang)
	return e.hl.Tokenize(code, lang, e.cfg.LightTheme)
}

// CSS returns the full stylesheet for this engine configuration.
func (e *Engine) CSS() string {
	return css.Generate(e.cfg, e.lightColors, e.darkColors)
}

// ThemeCSS returns only theme variables and token switching CSS, without
// structural rules. Use this for secondary engines on multi-engine pages
// where one primary engine provides the full CSS via CSS().
func (e *Engine) ThemeCSS() string {
	return css.GenerateThemeOnly(e.cfg, e.lightColors, e.darkColors)
}

// JS returns the JavaScript module for this engine configuration.
func (e *Engine) JS() string {
	return js.Generate(e.cfg)
}

// Assets returns CSS and JS with content-hashed filenames.
func (e *Engine) Assets() Assets {
	cssContent := e.CSS()
	jsContent := e.JS()
	return Assets{
		CSS: makeAssetFile(cssContent, "css"),
		JS:  makeAssetFile(jsContent, "js"),
	}
}

// EnableCodeGroups enables code group CSS/JS in engine output.
// Called automatically by kazarimd.CodeGroups(). It mutates the engine
// configuration without synchronization, so call it during setup before
// the engine is shared with concurrent Render, CSS, or JS callers.
func (e *Engine) EnableCodeGroups() {
	e.cfg.CodeGroups = true
}

func (e *Engine) resolveCollapse(code string, resolved *config.ResolvedBlock, spec *config.CollapseSpec) {
	lineCount := strings.Count(code, "\n") + 1
	result := collapsible.ResolveCollapse(
		lineCount, spec, e.cfg.Collapsible, code,
		resolved.LineMarkers, resolved.FocusLines,
	)
	resolved.CollapseThreshold = result.Threshold
	resolved.CollapseSegments = result.PreviewSegments
	resolved.CollapseBeyondCap = result.BeyondCapCount
	resolved.CollapseRanges = result.Ranges
	resolved.CollapseConfig = e.cfg.Collapsible
}
