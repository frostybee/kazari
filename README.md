# Kazari

A presentation layer for code blocks in Go. Kazari takes syntax-highlighted tokens from [Nuri](https://github.com/frostybee/nuri) or [Chroma](https://github.com/alecthomas/chroma) and renders them into complete code blocks with editor frames, terminal windows, copy buttons, line markers, collapsible sections, and more. The highlighter handles the colors. Kazari handles everything around them.

You call `Render()`, you get back HTML. Styling is yours to control through 60+ `--kz-*` CSS variables.

## Quick Start

```go
import (
    "github.com/frostybee/kazari"
    kazarichroma "github.com/frostybee/kazari/chroma"
)

hl := kazarichroma.New(kazarichroma.WithStyleMap(map[string]string{
    "light": "github",
    "dark":  "github-dark",
}))

engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemes("light", "dark"),
)

html, _ := engine.Render(code, kazari.Options{Lang: "go", Title: "main.go"})
css := engine.CSS()  // inject once in <head>
js := engine.JS()    // inject once before </body>
```

## What You Get

Editor and terminal frames, file name extraction, copy and fullscreen buttons, line numbers, line markers (mark/ins/del), labeled ranges, inline text markers, regex markers, focus lines, collapsible sections, code groups with tab sync, hybrid diff + syntax highlighting, word wrap, mermaid pass-through, ANSI escape rendering, per-block theme overrides, OKLCH theme adjustments, locale/i18n, file icons, contrast enforcement, and a Goldmark extension for markdown-based sites.

All visual customization goes through `--kz-*` CSS variables. No rebuilds needed.

## Usage Pattern

Create one engine per site configuration. Per-block variation (titles, markers, line numbers, collapse, theme overrides) flows through the meta string or `Options` struct, not through separate engines. Call `CSS()` once in `<head>` and `JS()` once before `</body>`. If you genuinely need two engine configurations on the same page, scope the second with `WithThemeCSSRoot()` to avoid CSS variable collisions and use `ThemeCSS()` instead of `CSS()` to emit only its theme variables without duplicating structural rules.

## Installation

```bash
go get github.com/frostybee/kazari
```

Then pick a highlighting backend:

```bash
go get github.com/frostybee/kazari/nuri    # VS Code themes, TextMate grammars
go get github.com/frostybee/kazari/chroma  # Pygments styles, 400+ languages
```

## Documentation

Coming soon.

## Acknowledgments

Kazari is inspired by [Expressive Code](https://expressive-code.com/) and brings its feature set to the Go ecosystem with a simpler architecture. There are no plugin hooks or lifecycle systems. Features are built-in and controlled by config flags. Styling lives entirely in CSS custom properties (`--kz-*`). The HTML structure, CSS styling, and class conventions are derived from Expressive Code by Hippo (MIT license).

## License

MIT
