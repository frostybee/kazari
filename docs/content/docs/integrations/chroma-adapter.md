---
title: "Chroma Adapter"
description: "Connect Kazari to Chroma for Go-native syntax highlighting."
tags: [highlighter, theming]
sidebar:
  order: 2
---

The `kazarichroma` package connects Kazari to [Chroma](https://github.com/alecthomas/chroma) for pure Go syntax highlighting with zero external dependencies and negligible initialization time. All Kazari presentation features work identically regardless of which adapter is used.

## Setup

Install Kazari and Chroma:

```bash
go get github.com/frostybee/kazari@latest
go get github.com/alecthomas/chroma/v2@latest
```

Wire the adapter at engine construction:

```go
import (
    "github.com/frostybee/kazari"
    kazarichroma "github.com/frostybee/kazari/chroma"
)

chromaHL := kazarichroma.New(kazarichroma.WithStyleMap(map[string]string{
    "github-light":        "github",
    "github-dark-default": "github-dark",
}))

engine := kazari.New(
    kazari.WithHighlighter(chromaHL),
    kazari.WithThemes("github-light", "github-dark-default"),
)
```

No context is required. There is no `Close` method and no deferred cleanup.

## Style map

Kazari theme names and Chroma style names use different naming conventions. `kazarichroma.WithStyleMap` bridges them:

```go
kazarichroma.WithStyleMap(map[string]string{
    "github-light":        "github",
    "github-dark-default": "github-dark",
    "one-dark-pro":        "monokai",
})
```

Keys are the theme names passed to `kazari.WithThemes()`. Values are Chroma style names as returned by Chroma's `styles.Names()` registry.

Without a style map, theme names are passed directly to Chroma's `styles.Get()`.

:::caution
Unknown Chroma style names silently resolve to the Fallback style, which produces monochrome output with no syntax differentiation. Verify style names against Chroma's bundled list before deploying.
:::

## Available styles

Chroma ships with 74+ built-in styles. Common light/dark pairs for dual-theme rendering:

| Light | Dark |
|-------|------|
| `github` | `github-dark` |
| `solarized-light` | `solarized-dark` |
| `gruvbox-light` | `gruvbox` |
| `catppuccin-latte` | `catppuccin-mocha` |
| `tokyonight-day` | `tokyonight-night` |
| `rose-pine-dawn` | `rose-pine` |
| `paraiso-light` | `paraiso-dark` |
| `modus-operandi` | `modus-vivendi` |

For the full list, call `styles.Names()` from `github.com/alecthomas/chroma/v2/styles`.

## Language names

The engine resolves language names via its internal `ResolveLanguage()` before passing them to the adapter. Pass lowercase names to Kazari regardless of what Chroma's registry uses internally.

Chroma supports 400+ languages. Call `lexers.Names(false)` from `github.com/alecthomas/chroma/v2/lexers` for the full list of canonical names.

## Theme colors

`GetThemeColors` populates `ThemeInfo` from the Chroma style:

| Field | Source |
|-------|--------|
| `FG` | Chroma `Text` entry foreground; falls back to `Background` entry foreground |
| `BG` | Chroma `Background` entry background |
| `LineNumberFG` | Chroma `LineNumbers` entry foreground |
| `SelectionBG` | Always empty (Chroma has no selection color) |
| `FoldBG` | Always empty (Chroma has no fold color) |

`GetThemeColors` never returns an error. An unknown style name silently produces a Fallback-based `ThemeInfo` with minimal color data.

## Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `WithStyleMap(map[string]string)` | `nil` | Map Kazari theme names to Chroma style names |

## Differences from Nuri

| Aspect | Nuri | Chroma |
|--------|------|--------|
| Initialization | Requires context and `nuri.New()` | Instant, no context |
| `DualThemeTokenizer` | Yes (one tokenization pass) | No (engine calls `Tokenize` twice) |
| Theme source | VS Code / Shiki themes (TextMate) | Chroma built-in styles |
| Theme name mapping | None needed | `WithStyleMap` required for most setups |
| Token granularity | Hierarchical TextMate scopes | ~92 flat Pygments token types |
| `SelectionBG` in `ThemeInfo` | Populated from theme | Always empty |
| `FoldBG` in `ThemeInfo` | Populated from theme | Always empty |
| `GetThemeColors` error | Returns error for unknown theme | Always returns nil |
| Unknown language | Returns error (engine falls back to plaintext) | Returns error (engine falls back to plaintext) |
| Strikethrough font style | Supported | Not mapped |
| Binary size overhead | ~400KB (embedded WASM) | ~0 (pure Go) |
| Language count | 257 (core: 39, full: 258) | ~400 |

Because `ChromaHighlighter` does not implement `DualThemeTokenizer`, the engine calls `Tokenize` once for the light theme and once for the dark theme when dual-theme rendering is active. For large code blocks this doubles the tokenization work compared to the [Nuri adapter](/integrations/nuri-adapter/). In practice, Chroma's tokenization is fast enough that this has minimal impact on total render time.
