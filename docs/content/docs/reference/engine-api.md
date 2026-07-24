---
title: "Engine API"
description: "Reference for all Engine methods and the kazari.New constructor."
sidebar:
  order: 3
---

`kazari.Engine` is the top-level type for rendering code blocks. Construct one with `kazari.New()`, then call its methods to produce HTML, CSS, JS, and raw tokens.

## Constructor

### `kazari.New`

```go
func New(opts ...Option) *Engine
```

Applies options left to right. Scalar fields: last writer wins. Map fields (`styleOverrides`, `languageDefaults`, `languageAliases`, `uiStrings`): keys merge, last value per key wins. Returns a ready-to-use `*Engine` (never nil). Without a highlighter, render falls back to plain text.

See [Configuration Options](/reference/configuration-options/) for the complete list of `With*` options.

## Render methods

### `engine.Render`

```go
func (e *Engine) Render(code string, opts Options) (string, error)
```

Renders a code block using structured per-block options. Returns a self-contained HTML fragment. Returns an error only if the highlighter fails. `Options` fields use pointer types so nil values fall through to engine defaults. See [Types & Interfaces](/reference/types-and-interfaces/) for the `Options` struct.

### `engine.RenderWithMeta`

```go
func (e *Engine) RenderWithMeta(code string, metaStr string) (string, error)
```

Parses `metaStr` as a markdown fence info string and renders the code block. An empty `metaStr` renders with engine defaults. Unknown meta tokens are silently skipped. Returns an error only if the highlighter fails.

This is the method used by the [Goldmark extension](/integrations/goldmark-extension/). See [Meta String Syntax](/reference/meta-string-syntax/) for every supported token.

## Asset output

### `engine.CSS`

```go
func (e *Engine) CSS() string
```

Returns the full stylesheet for this engine configuration: structural rules, CSS custom properties (`--kz-*`), theme variables, token switching rules, and feature-specific styles. Inject the returned string once per page inside `<style>` in `<head>`, or write it to a static `.css` file.

### `engine.ThemeCSS`

```go
func (e *Engine) ThemeCSS() string
```

Returns only theme variables and token switching CSS, without structural rules. Use this when multiple engines share the same structural CSS but need different theme palettes on the same page. Inject the full `CSS()` from one engine, then `ThemeCSS()` from additional engines scoped under different CSS selectors via `WithThemeCSSRoot`.

```go
primary := kazari.New(kazari.WithHighlighter(hl), kazari.WithThemes("github-light", "github-dark"))
secondary := kazari.New(kazari.WithHighlighter(hl), kazari.WithThemes("dracula", "dracula"), kazari.WithThemeCSSRoot(".dracula-scope"))

css := primary.CSS() + "\n" + secondary.ThemeCSS()
```

The primary engine emits full structural + theme CSS. The secondary engine contributes only its theme variables, scoped under `.dracula-scope`. This avoids duplicating structural rules.

### `engine.JS`

```go
func (e *Engine) JS() string
```

Returns the JavaScript module for interactive features: copy button, fullscreen, wrap toggle, collapsible sections, theme toggle, and code group tabs. Inject once per page before `</body>`. The returned string includes only the scripts for features enabled in the engine configuration. Returns an empty string if no interactive features are enabled.

### `engine.Assets`

```go
func (e *Engine) Assets() Assets
```

Returns CSS and JS content paired with content-hashed filenames (FNV-1a 32-bit, 8-char hex). Filenames follow the pattern `kazari-<hash>.css` and `kazari-<hash>.js`. Calls `engine.CSS()` and `engine.JS()` internally. See [Types & Interfaces](/reference/types-and-interfaces/) for the `Assets` and `AssetFile` struct definitions.

## Raw tokenization

### `engine.Tokenize`

```go
func (e *Engine) Tokenize(code string, lang string) ([][]Token, error)
```

Returns raw tokens for custom HTML output. Outer slice = lines, inner slice = tokens per line. Uses the light theme only. Language aliases are resolved before tokenization. Returns an error if no highlighter is configured. See [Types & Interfaces](/reference/types-and-interfaces/) for the `Token` struct.

## Code groups

### `engine.EnableCodeGroups`

```go
func (e *Engine) EnableCodeGroups()
```

Enables code group CSS and JS in the engine's asset output. Called automatically by `kazarimd.CodeGroups()`. **Not goroutine-safe**: call during setup before sharing the engine. See [Goldmark Extension](/integrations/goldmark-extension/) for wiring details.

## Thread safety

All `Engine` methods are safe for concurrent use except `EnableCodeGroups`:

- `Render`, `RenderWithMeta`, `Tokenize`: safe. Per-block theme overrides are cached behind a read-write mutex.
- `CSS`, `ThemeCSS`, `JS`, `Assets`: safe. These read configuration state set at construction.
- `EnableCodeGroups`: **not safe**. Must be called during setup before sharing the engine across goroutines.

:::caution
`EnableCodeGroups` is not goroutine-safe. Call it during setup before sharing the engine across goroutines.
:::

## Other public functions

These package-level functions are documented on their respective pages:

| Function | Page |
|----------|------|
| `ParseConfig(data []byte, format string) (*FileConfig, error)` | [File-Based Config](/reference/file-based-config/) |
| `FileConfigToOptions(fc *FileConfig) ([]Option, error)` | [File-Based Config](/reference/file-based-config/) |
| `LoadConfig(path string) ([]Option, error)` | [File-Based Config](/reference/file-based-config/) |
| `SelectorMode(selector string) DarkMode` | [Themes & Dark Mode](/styling/themes-and-dark-mode/) |
| `MediaQueryMode() DarkMode` | [Themes & Dark Mode](/styling/themes-and-dark-mode/) |
| `BothMode(selector string) DarkMode` | [Themes & Dark Mode](/styling/themes-and-dark-mode/) |
| `CreateInlineSVGURL(svg string) string` | [Icons](/features/icons/) |
