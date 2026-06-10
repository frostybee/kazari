// Package kazari renders framed, syntax-highlighted HTML code blocks
// with full CSS customization via custom properties.
package kazari

import (
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/collapsible"
	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/css"
	"github.com/frostybee/kazari/internal/frame"
	"github.com/frostybee/kazari/internal/js"
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
}

// New creates a new Engine with the given options.
func New(opts ...Option) *Engine {
	b := &engineBuilder{cfg: config.DefaultConfig()}

	for _, opt := range opts {
		opt(b)
	}

	e := &Engine{
		hl:  b.hl,
		cfg: b.cfg,
	}

	// Extract theme colors at construction time for CSS generation.
	if e.hl != nil {
		if info, err := e.hl.GetThemeColors(e.cfg.LightTheme); err == nil {
			e.lightColors = theme.ThemeColors{
				EditorBG:     info.BG,
				EditorFG:     info.FG,
				SelectionBG:  info.SelectionBG,
				LineNumberFG: info.LineNumberFG,
			}
		}
		if e.cfg.DarkTheme != "" {
			if info, err := e.hl.GetThemeColors(e.cfg.DarkTheme); err == nil {
				e.darkColors = theme.ThemeColors{
					EditorBG:     info.BG,
					EditorFG:     info.FG,
					SelectionBG:  info.SelectionBG,
					LineNumberFG: info.LineNumberFG,
				}
			}
		}
	}

	return e
}

// Render renders a code block with structured options.
func (e *Engine) Render(code string, opts Options) (string, error) {
	blockOpts := mapOptionsToBlockOpts(opts)
	lang := e.cfg.ResolveLanguage(opts.Lang)
	resolved := e.cfg.Resolve(lang, blockOpts)
	resolved.LineMarkers = convertLineMarkers(opts.LineMarkers)
	resolved.InlineMarkers = convertInlineMarkers(opts.InlineMarkers)
	resolved.FocusLines = convertRanges(opts.FocusLines)

	var spec *config.CollapseSpec
	if opts.Collapse != nil {
		spec = &config.CollapseSpec{
			Enabled:  opts.Collapse.Enabled,
			Disabled: opts.Collapse.Disabled,
			Ranges:   convertRanges(opts.Collapse.Ranges),
		}
	}
	e.resolveCollapse(code, resolved, spec)

	return e.renderResolved(code, resolved)
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
	e.resolveCollapse(code, resolved, parsed.Collapse)
	return e.renderResolved(code, resolved)
}

func (e *Engine) renderResolved(code string, resolved *config.ResolvedBlock) (string, error) {
	code = e.preprocess(code, resolved)

	lines, err := e.tokenize(code, resolved.Lang)
	if err != nil {
		return "", err
	}

	return render.RenderBlock(lines, resolved, e.cfg), nil
}

func (e *Engine) preprocess(code string, resolved *config.ResolvedBlock) string {
	if e.cfg.FileNameExtraction && resolved.Title == "" {
		title, modified := frame.ExtractFileName(code, resolved.Lang)
		if title != "" {
			resolved.Title = title
			code = modified
		}
	}

	if resolved.Frame == config.FrameAuto {
		if e.cfg.FrameDetection {
			resolved.Frame = frame.DetectFrameType(code, resolved.Lang, resolved.Frame)
		} else {
			resolved.Frame = config.FrameCode
		}
	}

	resolved.RawCode = code
	return code
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
// Called automatically by kazarimd.CodeGroups().
func (e *Engine) EnableCodeGroups() {
	e.cfg.CodeGroups = true
}

func convertLineMarkers(markers []LineMarker) []config.LineMarker {
	if len(markers) == 0 {
		return nil
	}
	out := make([]config.LineMarker, len(markers))
	for i, m := range markers {
		out[i] = config.LineMarker{
			Type:  config.MarkerType(m.Type),
			Lines: convertRanges(m.Lines),
			Label: m.Label,
		}
	}
	return out
}

func convertInlineMarkers(markers []InlineMarker) []config.InlineMarker {
	if len(markers) == 0 {
		return nil
	}
	out := make([]config.InlineMarker, len(markers))
	for i, m := range markers {
		out[i] = config.InlineMarker{
			Type: config.MarkerType(m.Type),
			Text: m.Text,
		}
	}
	return out
}

func convertRanges(ranges []Range) []config.LineRange {
	if len(ranges) == 0 {
		return nil
	}
	out := make([]config.LineRange, len(ranges))
	for i, r := range ranges {
		out[i] = config.LineRange{Start: r.Start, End: r.End}
	}
	return out
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

func (e *Engine) tokenize(code, lang string) ([]render.TokenLine, error) {
	if e.hl == nil {
		return plaintextLines(code), nil
	}

	if e.cfg.TabWidth > 0 {
		code = expandTabs(code, e.cfg.TabWidth)
	}

	lightTokens, err := e.hl.Tokenize(code, lang, e.cfg.LightTheme)
	if err != nil {
		return nil, err
	}

	var darkTokens [][]Token
	if e.cfg.DarkTheme != "" {
		darkTokens, err = e.hl.Tokenize(code, lang, e.cfg.DarkTheme)
		if err != nil {
			return nil, err
		}
	}

	return mergeTokens(lightTokens, darkTokens), nil
}
