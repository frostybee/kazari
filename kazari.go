// Package kazari renders framed, syntax-highlighted HTML code blocks
// with full CSS customization via custom properties.
package kazari

import (
	"fmt"

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
	return e.renderResolved(code, resolved)
}

// RenderWithMeta parses a meta string and renders the code block.
func (e *Engine) RenderWithMeta(code string, metaStr string) (string, error) {
	parsed := meta.Parse(metaStr)
	lang := e.cfg.ResolveLanguage(parsed.BlockOptions.Lang)
	parsed.BlockOptions.Lang = lang
	resolved := e.cfg.Resolve(lang, &parsed.BlockOptions)
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
