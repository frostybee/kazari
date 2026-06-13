package kazarichroma

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/frostybee/kazari"
)

// ChromaHighlighter adapts Chroma's lexer/style system to satisfy
// Kazari's Highlighter interface.
type ChromaHighlighter struct {
	styleMap map[string]string
}

// Option configures a ChromaHighlighter.
type Option func(*ChromaHighlighter)

// WithStyleMap maps Kazari theme names to Chroma style names.
// Keys are the theme names passed to kazari.WithThemes(); values
// are Chroma style names (as returned by styles.Names()).
func WithStyleMap(m map[string]string) Option {
	return func(h *ChromaHighlighter) {
		h.styleMap = m
	}
}

// New creates a ChromaHighlighter.
// Without a style map, Kazari theme names are passed directly to
// styles.Get() (which returns Fallback for unknown names).
func New(opts ...Option) *ChromaHighlighter {
	h := &ChromaHighlighter{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (c *ChromaHighlighter) resolveStyle(theme string) *chroma.Style {
	name := theme
	if c.styleMap != nil {
		if mapped, ok := c.styleMap[theme]; ok {
			name = mapped
		}
	}
	return styles.Get(name)
}

func (c *ChromaHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
	lexer := lexers.Get(lang)
	if lexer == nil {
		return nil, fmt.Errorf("kazarichroma: unknown language %q", lang)
	}
	lexer = chroma.Coalesce(lexer)

	style := c.resolveStyle(theme)
	themeBG := style.Get(chroma.Background).Background

	iter, err := lexer.Tokenise(nil, code)
	if err != nil {
		return nil, fmt.Errorf("kazarichroma: tokenise %q: %w", lang, err)
	}

	chromaLines := chroma.SplitTokensIntoLines(iter.Tokens())

	if len(chromaLines) == 0 {
		return [][]kazari.Token{{{Content: ""}}}, nil
	}

	lines := make([][]kazari.Token, len(chromaLines))
	for i, chromaLine := range chromaLines {
		tokens := make([]kazari.Token, 0, len(chromaLine))
		for _, ct := range chromaLine {
			content := strings.TrimRight(ct.Value, "\n")
			if content == "" && len(chromaLine) > 1 {
				continue
			}

			entry := style.Get(ct.Type)
			tok := kazari.Token{Content: content}

			if entry.Colour.IsSet() {
				tok.Color = entry.Colour.String()
			}
			if entry.Background.IsSet() && entry.Background != themeBG {
				tok.BgColor = entry.Background.String()
			}

			var fs int
			if entry.Italic == chroma.Yes {
				fs |= kazari.FontStyleItalic
			}
			if entry.Bold == chroma.Yes {
				fs |= kazari.FontStyleBold
			}
			if entry.Underline == chroma.Yes {
				fs |= kazari.FontStyleUnderline
			}
			tok.FontStyle = fs

			tokens = append(tokens, tok)
		}
		if len(tokens) == 0 {
			tokens = []kazari.Token{{Content: ""}}
		}
		lines[i] = tokens
	}

	return lines, nil
}

func (c *ChromaHighlighter) GetThemeColors(theme string) (kazari.ThemeInfo, error) {
	style := c.resolveStyle(theme)

	bgEntry := style.Get(chroma.Background)
	textEntry := style.Get(chroma.Text)

	info := kazari.ThemeInfo{}

	if textEntry.Colour.IsSet() {
		info.FG = textEntry.Colour.String()
	} else if bgEntry.Colour.IsSet() {
		info.FG = bgEntry.Colour.String()
	}

	if bgEntry.Background.IsSet() {
		info.BG = bgEntry.Background.String()
	}

	lnEntry := style.Get(chroma.LineNumbers)
	if lnEntry.Colour.IsSet() {
		info.LineNumberFG = lnEntry.Colour.String()
	}

	return info, nil
}

func (c *ChromaHighlighter) GetLoadedLanguages() []string {
	return lexers.Names(false)
}
