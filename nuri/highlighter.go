package kazarinuri

import (
	"context"

	"github.com/frostybee/kazari"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/ast"
)

// NuriHighlighter adapts a Nuri Highlighter to satisfy Kazari's Highlighter interface.
type NuriHighlighter struct {
	hl  *nuri.Highlighter
	ctx context.Context
}

// New creates an adapter from a Nuri Highlighter instance.
func New(ctx context.Context, hl *nuri.Highlighter) *NuriHighlighter {
	return &NuriHighlighter{hl: hl, ctx: ctx}
}

func (n *NuriHighlighter) Tokenize(code, lang, themeName string) ([][]kazari.Token, error) {
	result, err := n.hl.CodeToTokens(n.ctx, code, ast.CodeToTokensOptions{
		Lang:  lang,
		Theme: themeName,
	})
	if err != nil {
		return nil, err
	}

	lines := make([][]kazari.Token, len(result.Tokens))
	for i, nuriLine := range result.Tokens {
		tokens := make([]kazari.Token, len(nuriLine))
		for j, nt := range nuriLine {
			tokens[j] = kazari.Token{
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
func (n *NuriHighlighter) TokenizeDual(code, lang, lightTheme, darkTheme string) ([][]kazari.Token, [][]kazari.Token, error) {
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

	light := make([][]kazari.Token, len(result.Tokens))
	dark := make([][]kazari.Token, len(result.Tokens))
	for i, line := range result.Tokens {
		lightLine := make([]kazari.Token, len(line))
		darkLine := make([]kazari.Token, len(line))
		for j, nt := range line {
			darkLine[j] = kazari.Token{
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

func (n *NuriHighlighter) GetThemeColors(themeName string) (kazari.ThemeInfo, error) {
	tc, err := n.hl.GetThemeColors(themeName)
	if err != nil {
		return kazari.ThemeInfo{}, err
	}

	return kazari.ThemeInfo{
		FG:           tc.Foreground,
		BG:           tc.Background,
		SelectionBG:  tc.SelectionBackground,
		LineNumberFG: tc.Colors["editorLineNumber.foreground"],
		FoldBG:       tc.Colors["editor.foldBackground"],
	}, nil
}

func (n *NuriHighlighter) GetLoadedLanguages() []string {
	return n.hl.LoadedLanguages()
}
