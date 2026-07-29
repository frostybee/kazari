---
title: "Custom Highlighters"
description: "Implement Kazari's Highlighter interface to use any syntax highlighting engine, or call Engine.Tokenize for raw token access."
tags: [highlighter]
sidebar:
  order: 2
---

A custom `Highlighter` lets Kazari work with grammars or engines beyond the built-in [Nuri](/integrations/nuri-adapter/) and [Chroma](/integrations/chroma-adapter/) adapters. Implementing the interface connects any tokenizer to Kazari's rendering pipeline. For consumers who need full control over the output HTML, `Engine.Tokenize` exposes raw tokens without any rendering.

## The Highlighter interface

The interface has three methods:

```go
type Highlighter interface {
    Tokenize(code, lang, theme string) ([][]Token, error)
    GetThemeColors(theme string) (ThemeInfo, error)
    GetLoadedLanguages() []string
}
```

**`Tokenize`** receives source code, a language identifier, and a theme name. It returns a two-dimensional slice of tokens: one inner slice per line, each token carrying `Content`, `Color`, `BgColor`, and `FontStyle` fields. Return an error for genuinely unsupported languages (the engine falls back to plaintext for unknown languages and logs a warning).

**`GetThemeColors`** returns a `ThemeInfo` struct with the theme's background, foreground, selection, line number, and fold colors. Kazari calls this once per theme at engine construction to generate CSS variables.

**`GetLoadedLanguages`** returns all language identifiers the highlighter can tokenize. Kazari uses this to distinguish "unknown language" (plaintext fallback with warning) from "known language that errored" (hard error).

### Minimal implementation

A highlighter that returns unstyled tokens for any language:

```go
type PlaintextHighlighter struct{}

func (p *PlaintextHighlighter) Tokenize(code, lang, theme string) ([][]kazari.Token, error) {
    lines := strings.Split(code, "\n")
    result := make([][]kazari.Token, len(lines))
    for i, line := range lines {
        result[i] = []kazari.Token{{Content: line}}
    }
    return result, nil
}

func (p *PlaintextHighlighter) GetThemeColors(theme string) (kazari.ThemeInfo, error) {
    return kazari.ThemeInfo{
        FG: "#24292f",
        BG: "#ffffff",
    }, nil
}

func (p *PlaintextHighlighter) GetLoadedLanguages() []string {
    return []string{"plaintext"}
}
```

```go
engine := kazari.New(
    kazari.WithHighlighter(&PlaintextHighlighter{}),
)
```

For production implementations, see [`nuri/highlighter.go`](https://github.com/frostybee/kazari/blob/main/nuri/highlighter.go) (full TextMate grammar support with dual-theme optimization) and [`chroma/highlighter.go`](https://github.com/frostybee/kazari/blob/main/chroma/highlighter.go) (Chroma lexer/style adaptation).

## DualThemeTokenizer

When Kazari renders with both a light and dark theme, it needs two token streams with identical boundaries (same line count, same token count per line, same `Content` values). By default, the engine calls `Tokenize` twice and merges the results.

Highlighters that can resolve both themes in a single pass implement the optional `DualThemeTokenizer` interface:

```go
type DualThemeTokenizer interface {
    TokenizeDual(code, lang, lightTheme, darkTheme string) (light, dark [][]Token, err error)
}
```

The engine type-asserts this interface on its `Highlighter` at tokenization time. When present and a dark theme is configured, a single `TokenizeDual` call replaces two `Tokenize` calls. The Nuri adapter implements this interface; the Chroma adapter does not (Chroma's lexer API does not support multi-theme tokenization, so two passes are used).

The two returned streams must have identical token boundaries. TextMate scanning is theme-independent, so any TextMate-based highlighter can satisfy this constraint naturally.

## Raw tokens with Engine.Tokenize

`Engine.Tokenize` exposes raw token data for consumers who build their own HTML instead of using Kazari's renderer:

```go
tokens, err := engine.Tokenize(code, "go")
// tokens is [][]Token: one inner slice per source line
```

Each `Token` carries:

| Field | Type | Description |
|---|---|---|
| `Content` | `string` | The text content of the token |
| `Color` | `string` | Foreground hex color (e.g. `"#cf222e"`) |
| `BgColor` | `string` | Background hex color (usually empty) |
| `FontStyle` | `int` | Bitmask: `FontStyleItalic` (1), `FontStyleBold` (2), `FontStyleUnderline` (4), `FontStyleStrikethrough` (8) |

`Tokenize` uses the engine's light theme only. It resolves language aliases (via `WithLanguageAliases`) before calling the underlying highlighter. A configured highlighter is required; calling `Tokenize` without one returns an error.

## Observing failures

When an unknown language falls back to plaintext, the engine emits a warning through `WithWarningHandler`. For custom highlighters, this lets the consumer log or react to languages the highlighter does not cover:

```go
engine := kazari.New(
    kazari.WithHighlighter(customHL),
    kazari.WithWarningHandler(func(msg string) {
        log.Printf("[kazari] %s", msg)
    }),
)
```

Without a handler, warnings go to the standard logger via `log.Print`.

## Related

- [Types and Interfaces](/reference/types-and-interfaces/) for full field definitions of `Token`, `ThemeInfo`, and font style constants
- [Nuri Adapter](/integrations/nuri-adapter/) and [Chroma Adapter](/integrations/chroma-adapter/) for production adapter implementations
- [Post-Render Callbacks](/plugins/post-render-callbacks/) for modifying rendered HTML after tokenization
