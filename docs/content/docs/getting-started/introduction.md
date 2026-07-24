---
title: "Introduction"
description: "Kazari renders framed, syntax-highlighted HTML code blocks in Go with full CSS customization."
sidebar:
  order: 1
---

Kazari is a Go library that renders rich, syntax-highlighted code blocks as standalone HTML. It pairs with a highlighter like Nuri or Chroma for coloring and adds everything else on top: editor and terminal frames, line numbers, copy buttons, collapsible sections, markers, and dual-theme support. All styling is controlled through CSS custom properties, and the output is server-rendered with no framework dependency.

## How it works

Kazari operates in three steps:

1. **Tokenize** the source code via a pluggable highlighter (Nuri or Chroma)
2. **Render** the tokens into an HTML block with frames, markers, line numbers, and interactive controls
3. **Return** separate CSS, JS, and HTML strings for the consumer to place on the page

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemes("github-light", "github-dark"),
    kazari.WithCopyButton(true),
    kazari.WithLineNumbers(true),
)

html, err := engine.RenderWithMeta(code, `go title="main.go" {3-5}`)
css := engine.CSS()  // inject once in <head>
js := engine.JS()    // inject once before </body>
```

`engine.CSS()` and `engine.JS()` return page-wide assets. Inject each once. `engine.Render()` and `engine.RenderWithMeta()` return per-block HTML.

## What Kazari provides

- Editor and terminal frames with title bars, file icons, and language badges
- Toolbar buttons: copy, word wrap, fullscreen, and per-block theme toggle
- Line markers, inline markers, regex markers, focus lines, and labeled ranges
- Diff+syntax hybrid rendering and clickable inline links
- Collapsible sections with threshold-based or range-based collapse
- Dual-theme rendering: light and dark colors baked into HTML, theme switching is pure CSS
- WCAG contrast enforcement, ARIA labels, and `all: revert` style isolation
- 90+ CSS custom properties for visual customization

See [Features](/features/) for the full list with examples and configuration.

## Three configuration layers

Engine defaults, language defaults, and per-block meta strings form a cascade where each layer overrides the previous. See [Configuration Layers](/getting-started/configuration-layers/) for details.

## Kazari is not a highlighter

Kazari handles presentation. Nuri and Chroma handle tokenization. This distinction matters throughout the documentation:

- **Engine** refers to the top-level `kazari.Engine` struct
- **Block** refers to a rendered code block
- **Meta string** refers to the options after the opening triple backticks

Nuri is a pure Go port of Shiki that uses TextMate grammars and VS Code themes. Chroma is a Go syntax highlighter with its own theme format. Kazari accepts either through its `Highlighter` interface and adds the visual layer on top.

## Integrations

Kazari works standalone or through [Goldmark](/integrations/goldmark-extension/) for Markdown rendering. Configuration can also be loaded from a [YAML file](/reference/file-based-config/) instead of Go code.

## Styling

All visual properties are controlled by CSS custom properties with the `--kz-` prefix. Override them in a stylesheet to customize colors, spacing, fonts, borders, and shadows without rebuilding:

```css
.kazari-block {
    --kz-radius: 0.5rem;
    --kz-font-family: 'Fira Code', monospace;
    --kz-editor-bg: #1a1b26;
}
```

No Go code changes are needed for visual adjustments. Inspect any code block in DevTools to discover available variables.
