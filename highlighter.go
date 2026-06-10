package kazari

import (
	"context"

	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/ast"
)

// NuriHighlighter adapts a Nuri Highlighter to satisfy Kazari's Highlighter interface.
type NuriHighlighter struct {
	hl  *nuri.Highlighter
	ctx context.Context
}

// NewNuriHighlighter creates an adapter from a Nuri Highlighter instance.
func NewNuriHighlighter(ctx context.Context, hl *nuri.Highlighter) *NuriHighlighter {
	return &NuriHighlighter{hl: hl, ctx: ctx}
}

func (n *NuriHighlighter) Tokenize(code, lang, themeName string) ([][]Token, error) {
	result, err := n.hl.CodeToTokens(n.ctx, code, ast.CodeToTokensOptions{
		Lang:  lang,
		Theme: themeName,
	})
	if err != nil {
		return nil, err
	}

	lines := make([][]Token, len(result.Tokens))
	for i, nuriLine := range result.Tokens {
		tokens := make([]Token, len(nuriLine))
		for j, nt := range nuriLine {
			tokens[j] = Token{
				Content:   nt.Content,
				Color:     nt.Color,
				BgColor:   nt.BgColor,
				FontStyle: int(nt.FontStyle),
			}
		}
		lines[i] = tokens
	}
	return lines, nil
}

// TokenizeDual implements DualThemeTokenizer: one tokenization pass resolves
// both themes via Nuri's multi-theme mode, halving dual-theme cost. Both
// returned streams share token boundaries by construction.
func (n *NuriHighlighter) TokenizeDual(code, lang, lightTheme, darkTheme string) ([][]Token, [][]Token, error) {
	result, err := n.hl.CodeToTokens(n.ctx, code, ast.CodeToTokensOptions{
		Lang: lang,
		Themes: map[string]string{
			"light": lightTheme,
			"dark":  darkTheme,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	// Nuri sorts theme keys, so "dark" is the default theme: its style fills
	// each token's own Color/BgColor/FontStyle fields, while "light" lives in
	// ThemeStyles["light"]. A missing ThemeStyles entry falls back to the
	// default-theme style rather than emitting an uncolored token.
	light := make([][]Token, len(result.Tokens))
	dark := make([][]Token, len(result.Tokens))
	for i, line := range result.Tokens {
		lightLine := make([]Token, len(line))
		darkLine := make([]Token, len(line))
		for j, nt := range line {
			darkLine[j] = Token{
				Content:   nt.Content,
				Color:     nt.Color,
				BgColor:   nt.BgColor,
				FontStyle: int(nt.FontStyle),
			}
			lightLine[j] = darkLine[j]
			if ls, ok := nt.ThemeStyles["light"]; ok {
				lightLine[j].Color = ls.Color
				lightLine[j].BgColor = ls.BgColor
				lightLine[j].FontStyle = int(ls.FontStyle)
			}
		}
		light[i] = lightLine
		dark[i] = darkLine
	}
	return light, dark, nil
}

func (n *NuriHighlighter) GetThemeColors(themeName string) (ThemeInfo, error) {
	result, err := n.hl.CodeToTokens(n.ctx, "", ast.CodeToTokensOptions{
		Lang:  "text",
		Theme: themeName,
	})
	if err != nil {
		return ThemeInfo{}, err
	}

	return ThemeInfo{
		FG: result.FG,
		BG: result.BG,
	}, nil
}

func (n *NuriHighlighter) GetLoadedLanguages() []string {
	return n.hl.LoadedLanguages()
}
