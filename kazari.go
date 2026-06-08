// Package kazari renders framed, syntax-highlighted HTML code blocks
// with full CSS customization via custom properties.
package kazari

import (
	"fmt"
	"hash/fnv"
	"strings"

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

	code = e.preprocess(code, resolved)

	lines, err := e.tokenize(code, resolved.Lang)
	if err != nil {
		return "", err
	}

	return render.RenderBlock(lines, resolved, e.cfg), nil
}

// RenderWithMeta parses a meta string and renders the code block.
func (e *Engine) RenderWithMeta(code string, metaStr string) (string, error) {
	parsed := meta.Parse(metaStr)
	lang := e.cfg.ResolveLanguage(parsed.BlockOptions.Lang)
	parsed.BlockOptions.Lang = lang
	resolved := e.cfg.Resolve(lang, &parsed.BlockOptions)

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

	// Normalize tabs.
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

// mergeTokens pairs light and dark tokens into MergedToken lines.
func mergeTokens(light, dark [][]Token) []render.TokenLine {
	lines := make([]render.TokenLine, len(light))

	for i, lightLine := range light {
		var darkLine []Token
		if dark != nil && i < len(dark) {
			darkLine = dark[i]
		}

		if darkLine == nil || len(lightLine) == len(darkLine) {
			// Fast path: boundaries match (common case).
			tokens := make([]render.MergedToken, len(lightLine))
			for j, lt := range lightLine {
				mt := render.MergedToken{
					Content:    lt.Content,
					LightColor: lt.Color,
					LightBG:    lt.BgColor,
					FontStyle:  lt.FontStyle,
				}
				if darkLine != nil && j < len(darkLine) {
					mt.DarkColor = darkLine[j].Color
					mt.DarkBG = darkLine[j].BgColor
				}
				tokens[j] = mt
			}
			lines[i] = render.TokenLine{Tokens: tokens}
		} else {
			// Slow path: boundaries differ — align by character position.
			lines[i] = render.TokenLine{Tokens: alignTokens(lightLine, darkLine)}
		}
	}

	return lines
}

// alignTokens handles the rare case where light and dark tokens have different boundaries.
func alignTokens(lightLine, darkLine []Token) []render.MergedToken {
	var result []render.MergedToken
	li, di := 0, 0   // token indices
	lo, do := 0, 0   // character offsets within current token

	for li < len(lightLine) && di < len(darkLine) {
		lt := lightLine[li]
		dt := darkLine[di]
		lRemain := len(lt.Content) - lo
		dRemain := len(dt.Content) - do
		take := lRemain
		if dRemain < take {
			take = dRemain
		}

		result = append(result, render.MergedToken{
			Content:    lt.Content[lo : lo+take],
			LightColor: lt.Color,
			DarkColor:  dt.Color,
			LightBG:    lt.BgColor,
			DarkBG:     dt.BgColor,
			FontStyle:  lt.FontStyle,
		})

		lo += take
		do += take
		if lo >= len(lt.Content) {
			li++
			lo = 0
		}
		if do >= len(dt.Content) {
			di++
			do = 0
		}
	}

	// Remaining light tokens (no dark counterpart).
	for li < len(lightLine) {
		lt := lightLine[li]
		content := lt.Content[lo:]
		if content != "" {
			result = append(result, render.MergedToken{
				Content:    content,
				LightColor: lt.Color,
				LightBG:    lt.BgColor,
				FontStyle:  lt.FontStyle,
			})
		}
		li++
		lo = 0
	}

	return result
}

// plaintextLines returns single-token lines for code with no highlighter.
func plaintextLines(code string) []render.TokenLine {
	rawLines := splitLines(code)
	lines := make([]render.TokenLine, len(rawLines))
	for i, content := range rawLines {
		lines[i] = render.TokenLine{
			Tokens: []render.MergedToken{{Content: content}},
		}
	}
	return lines
}

// splitLines splits code into lines, handling the trailing newline correctly.
func splitLines(code string) []string {
	if code == "" {
		return []string{""}
	}
	lines := strings.Split(code, "\n")
	// If code ends with \n, the Split produces a trailing empty string — remove it.
	if strings.HasSuffix(code, "\n") && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func expandTabs(code string, tabWidth int) string {
	if !strings.Contains(code, "\t") {
		return code
	}
	spaces := strings.Repeat(" ", tabWidth)
	return strings.ReplaceAll(code, "\t", spaces)
}

func mapOptionsToBlockOpts(opts Options) *config.BlockOptions {
	bo := &config.BlockOptions{
		Lang:            opts.Lang,
		Title:           opts.Title,
		StartLineNumber: opts.StartLineNumber,
	}
	if opts.Frame != nil {
		f := int(*opts.Frame)
		bo.Frame = &f
	}
	if opts.LineNumbers != nil {
		bo.LineNumbers = opts.LineNumbers
	}
	if opts.Wrap != nil {
		bo.Wrap = opts.Wrap
	}
	return bo
}

func makeAssetFile(content, ext string) AssetFile {
	h := fnv.New32a()
	h.Write([]byte(content))
	hash := fmt.Sprintf("%08x", h.Sum32())
	return AssetFile{
		Content:  content,
		Hash:     hash,
		Filename: fmt.Sprintf("kazari-%s.%s", hash, ext),
	}
}
