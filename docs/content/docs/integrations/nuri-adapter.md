---
title: "Nuri Adapter"
description: "Connect Kazari to Nuri for TextMate grammar-based syntax highlighting."
sidebar:
  order: 1
---

The `kazarinuri` package connects Kazari to [Nuri](https://github.com/frostybee/nuri), a pure Go port of Shiki that uses TextMate grammars and VS Code themes. Nuri handles tokenization with VS Code-identical fidelity; Kazari adds the presentation layer.

## Setup

Install both packages:

```bash
go get github.com/frostybee/kazari@latest
go get github.com/frostybee/nuri@latest
```

Wire the adapter at engine construction:

```go
import (
    "context"
    "log"

    "github.com/frostybee/kazari"
    kazarinuri "github.com/frostybee/kazari/nuri"
    "github.com/frostybee/nuri"
    "github.com/frostybee/nuri/bundle/core"
)

func main() {
    ctx := context.Background()

    hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
    if err != nil {
        log.Fatal(err)
    }
    defer hl.Close(ctx)

    nuriHL := kazarinuri.New(ctx, hl)

    engine := kazari.New(
        kazari.WithHighlighter(nuriHL),
        kazari.WithThemes("github-light", "github-dark"),
    )

    // engine.CSS() and engine.JS() return styles and scripts to inject once per page.
    // engine.RenderWithMeta(code, meta) renders each code block.
}
```

`core.FS()` provides 39 bundled language grammars and 65 VS Code themes from the `github.com/frostybee/nuri/bundle/core` package. `nuri.New` compiles the embedded Oniguruma WASM engine and returns an error if initialization fails.

:::caution
Always call `defer hl.Close(ctx)` after creating the Nuri highlighter. This releases the underlying WASM runtime when the program exits. Forgetting this call leaks resources.
:::

The `kazarinuri.New(ctx, hl)` call wraps the Nuri highlighter in an adapter that satisfies Kazari's `Highlighter` interface. Both the context and the Nuri pointer are required.

## Language bundles

Nuri ships two language bundles. Both include the same 65 VS Code themes.

| Bundle | Grammars | Import path | Use case |
|--------|----------|-------------|----------|
| `core` | 39 | `github.com/frostybee/nuri/bundle/core` | Most applications. Covers common languages (Go, JS, TS, Python, Rust, HTML, CSS, JSON, YAML, Bash, and more). |
| `full` | 258 | `github.com/frostybee/nuri/bundle/full` | Broad language coverage. Includes all VS Code built-in grammars. |

Both expose a single `FS() fs.FS` function returning an embedded filesystem with `grammars/` and `themes/` subdirectories. Pass the result to `nuri.WithFS`:

```go
import "github.com/frostybee/nuri/bundle/full"

hl, err := nuri.New(ctx, nuri.WithFS(full.FS()))
```

To register a custom grammar or theme not included in either bundle, use `nuri.WithGrammar` or `nuri.WithTheme` alongside a bundle:

```go
hl, err := nuri.New(ctx,
    nuri.WithFS(core.FS()),
    nuri.WithGrammar("mylang", grammarJSON),
    nuri.WithTheme("my-theme", themeJSON),
)
```

## Dual-theme tokenization

`NuriHighlighter` implements the `DualThemeTokenizer` interface:

```go
type DualThemeTokenizer interface {
    TokenizeDual(code, lang, lightTheme, darkTheme string) (light, dark [][]Token, err error)
}
```

When two themes are configured via `kazari.WithThemes(light, dark)`, the engine detects this interface and calls `TokenizeDual` instead of calling `Tokenize` twice. A single Nuri tokenization pass produces both light and dark token streams with identical structure: same line count, same token count per line, same content, different colors.

This optimization reduces tokenization cost for dual-theme rendering. The [Chroma adapter](/integrations/chroma-adapter/) does not implement `DualThemeTokenizer`; it always requires two separate `Tokenize` calls.

## Theme names

Pass Nuri-compatible theme names directly to `kazari.WithThemes`. Common pairs:

| Light | Dark |
|-------|------|
| `"github-light"` | `"github-dark"` or `"github-dark-default"` |
| `"one-light"` | `"one-dark-pro"` |
| `"solarized-light"` | `"solarized-dark"` |
| `"rose-pine-dawn"` | `"rose-pine"` or `"rose-pine-moon"` |
| `"vitesse-light"` | `"vitesse-dark"` |

Theme names are passed to Nuri without mapping or aliasing. `GetThemeColors` returns an error for unrecognized theme names.

The full list of bundled themes is available in the [Nuri documentation](https://frostybee.github.com/nuri).

## Theme colors

`GetThemeColors` reads VS Code-style color keys from the theme and returns a fully populated `ThemeInfo`:

| Field | Source key |
|-------|-----------|
| `FG` | Theme foreground |
| `BG` | Theme background |
| `SelectionBG` | `editor.selectionBackground` |
| `LineNumberFG` | `editorLineNumber.foreground` |
| `FoldBG` | `editor.foldBackground` |

All five fields are populated for standard VS Code themes. The [Chroma adapter](/integrations/chroma-adapter/) leaves `SelectionBG` and `FoldBG` empty because Chroma themes do not include editor color keys.

## Nuri configuration

These `nuri.New` options control the tokenization engine. They are set on the Nuri highlighter, not on the Kazari engine.

| Option | Default | Description |
|--------|---------|-------------|
| `WithFS(fs.FS)` | required | Grammar and theme filesystem from `core.FS()` or `full.FS()` |
| `WithPoolSize(n)` | `runtime.NumCPU()` | Maximum concurrent WASM instances for tokenization |
| `WithTimeoutMs(ms)` | `0` (no limit) | Per-line timeout in milliseconds; timed-out lines emit unstyled tokens with a diagnostic |
| `WithMaxLineLength(n)` | `0` (no limit) | Lines exceeding this length emit a single unstyled token |
| `WithMinContrast(ratio)` | `5.5` | WCAG contrast floor between token foreground and editor background; `0` disables |
| `WithCompilationCacheDir(dir)` | none | On-disk cache for the AOT-compiled WASM binary; speeds up cold start for CLI and SSG use |
| `WithGrammar(name, data)` | none | Register a custom TextMate grammar JSON at construction |
| `WithTheme(name, data)` | none | Register a custom VS Code theme JSON at construction |
| `WithAlias(alias, target)` | none | Map a language alias to an existing grammar name |

Pool sizing example for a documentation site building many pages concurrently:

```go
hl, err := nuri.New(ctx,
    nuri.WithFS(core.FS()),
    nuri.WithPoolSize(4),
    nuri.WithTimeoutMs(500),
)
```

:::note
Nuri's `WithMinContrast` adjusts token foreground colors at theme load time to meet the specified WCAG contrast ratio against the editor background. Kazari has a separate [`WithMinContrast`](/reference/configuration-options/) option that adjusts marker background contrast. The two are independent and complement each other.
:::

## SVG output

Nuri can render syntax-highlighted code as standalone SVG images, useful for README files, email content, or static exports where HTML and CSS are not available. This feature is called directly on the `*nuri.Highlighter` instance, not through the Kazari adapter.

```go
svg, err := hl.CodeToSVG(ctx, code, nuri.CodeToSVGOptions{
    Lang:  "go",
    Theme: "github-dark",
})
```

### `CodeToSVGOptions` fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Lang` | `string` | required | Language for syntax highlighting |
| `Theme` | `string` | required | Theme name for colors |
| `FontSize` | `float64` | `14` | Font size in pixels |
| `CornerRadius` | `float64` | `8` | Border radius of the SVG container |
| `ShowBackground` | `*bool` | `true` | Set to `false` for a transparent background |

### Examples

Custom font size:

```go
svg, _ := hl.CodeToSVG(ctx, code, nuri.CodeToSVGOptions{
    Lang:     "go",
    Theme:    "github-dark",
    FontSize: 18,
})
```

Transparent background (no container fill):

```go
noBG := false
svg, _ := hl.CodeToSVG(ctx, code, nuri.CodeToSVGOptions{
    Lang:           "go",
    Theme:          "github-dark",
    ShowBackground: &noBG,
    CornerRadius:   0,
})
```

Dual-theme output (generate separate SVGs for light and dark contexts):

```go
lightSVG, _ := hl.CodeToSVG(ctx, code, nuri.CodeToSVGOptions{
    Lang:  "go",
    Theme: "github-light",
})
darkSVG, _ := hl.CodeToSVG(ctx, code, nuri.CodeToSVGOptions{
    Lang:  "go",
    Theme: "github-dark",
})
```

This produces two standalone images that can be toggled with CSS `prefers-color-scheme` media queries or displayed side by side.

## Language loading

Nuri loads language grammars lazily. A language appears in `GetLoadedLanguages()` only after the first `Tokenize` or `TokenizeDual` call for that language. This is expected behavior and does not affect rendering.

When a block specifies an unknown language, the engine falls back to plaintext rendering and emits a warning through the configured warning handler.

The full list of supported languages is available in the [Nuri documentation](https://frostybee.github.com/nuri).

## Context and lifecycle

The Nuri adapter captures the `context.Context` passed to `kazarinuri.New` and reuses it for every subsequent `Tokenize`, `TokenizeDual`, and `GetThemeColors` call. Kazari's `Highlighter` interface does not accept a per-call context, so cancellation must be managed at the Nuri highlighter level.

For long-running servers, pass a context tied to the server's lifecycle:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
if err != nil {
    log.Fatal(err)
}
defer hl.Close(ctx)

nuriHL := kazarinuri.New(ctx, hl)
```

For CLI tools and SSGs that process a batch of files and exit, `context.Background()` is sufficient.
