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
