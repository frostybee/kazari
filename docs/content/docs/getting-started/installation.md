---
title: "Quick Start"
description: "Get started with Kazari: install, render a code block, and view it in the browser."
tags: [markdown]
sidebar:
  order: 2
---

Render a syntax-highlighted code block and write it to an HTML file in under 30 lines of Go. By the end of this page, a working `output.html` opens in the browser with a framed, dual-theme code block.

## Install

Kazari requires Go 1.25 or later and a syntax highlighter. Install Kazari and the Nuri adapter (recommended):

```bash
go get github.com/frostybee/kazari@latest
go get github.com/frostybee/nuri@latest
```

Chroma is also supported as an alternative highlighter. See the [Chroma Adapter](/integrations/chroma-adapter/) page for setup.

## Create the engine

Set up a Nuri highlighter, wrap it with the Kazari adapter, and pass it to `kazari.New()`:

```go
ctx := context.Background()
hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
if err != nil {
    log.Fatal(err)
}
defer hl.Close(ctx)

engine := kazari.New(
    kazari.WithHighlighter(kazarinuri.New(ctx, hl)),
    kazari.WithThemes("github-light", "github-dark"),
)
```

`WithThemes` sets both a light and dark theme. Kazari bakes both sets of token colors into the HTML at build time via CSS custom properties. Theme switching is pure CSS with no JavaScript and no flash on toggle or page load. Dark mode defaults to a `.dark` class selector on the root element. See the [Themes & Dark Mode](/styling/themes-and-dark-mode/) page for other strategies.

## Render a code block

Pass source code and a meta string to `engine.RenderWithMeta()`:

```go
code := `func main() {
    cfg := loadConfig()
    svc := newService(cfg)
    svc.Run()
}`

html, err := engine.RenderWithMeta(code, `go title="main.go" showLineNumbers {2}`)
```

The meta string controls per-block behavior. In this example:

| Segment | Effect |
|---|---|
| `go` | Language for syntax highlighting |
| `title="main.go"` | Title displayed in the frame header |
| `showLineNumbers` | Show line numbers in the gutter |
| `{2}` | Highlight line 2 with a colored background |

## Assemble the page

Kazari returns three separate outputs. Combine them into a complete HTML file:

```go
package main

import (
    "context"
    "log"
    "os"

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

    engine := kazari.New(
        kazari.WithHighlighter(kazarinuri.New(ctx, hl)),
        kazari.WithThemes("github-light", "github-dark"),
    )

    code := `func main() {
    cfg := loadConfig()
    svc := newService(cfg)
    svc.Run()
}`
    block, err := engine.RenderWithMeta(code, `go title="main.go" showLineNumbers {2}`)
    if err != nil {
        log.Fatal(err)
    }

    page := "<!DOCTYPE html><html><head><style>" +
        engine.CSS() +
        "</style></head><body>" +
        block +
        "<script>" + engine.JS() + "</script></body></html>"

    os.WriteFile("output.html", []byte(page), 0644)
}
```

Run the program and open `output.html` in a browser:

```bash
go run main.go
```

→ A framed code block appears with the title "main.go", line numbers, and line 2 highlighted in yellow.

## Output structure

Each method returns a string. The consumer decides where to place them.

| Method | Returns | Inject |
|---|---|---|
| `engine.CSS()` | Full stylesheet (structural rules + theme variables) | Once in `<head>` |
| `engine.JS()` | Scripts for interactive features (copy, fullscreen, collapse) | Once before `</body>` |
| `engine.Render()` | Per-block HTML | Once per code block |
| `engine.RenderWithMeta()` | Per-block HTML (parses a meta string) | Once per code block |

`engine.Assets()` returns the same CSS and JS content with content-hashed filenames, suitable for production builds with cache busting. See [Engine API](/reference/engine-api/) for details.

## Add features

Enable common features with functional options passed to `kazari.New()`:

```go
engine := kazari.New(
    kazari.WithHighlighter(kazarinuri.New(ctx, hl)),
    kazari.WithThemes("github-light", "github-dark"),
    kazari.WithCopyButton(true),
    kazari.WithLineNumbers(true),
    kazari.WithCollapsible(kazari.CollapsibleConfig{
        LineThreshold: 20,
        PreviewLines:  8,
    }),
)
```

Per-block settings override engine defaults via the meta string:

````
```go title="server.go" showLineNumbers {3-5} collapse
````

See [Features](/features/) for the full list of available options.

## With Goldmark

Render fenced code blocks in Markdown through Goldmark:

```go
import (
    kazarimd "github.com/frostybee/kazari/goldmark"
    "github.com/yuin/goldmark"
)

md := goldmark.New(
    goldmark.WithExtensions(kazarimd.New(engine)),
)
```

Add code group support for tabbed containers:

```go
md := goldmark.New(
    goldmark.WithExtensions(
        kazarimd.New(engine),
        kazarimd.CodeGroups(engine),
    ),
)
```

See the [Goldmark Extension](/integrations/goldmark-extension/) page for configuration details.

## Next steps

- [Configuration Layers](/getting-started/configuration-layers/) explains the three-layer cascade: engine defaults, language defaults, and per-block meta strings
- [Features](/features/) covers every feature with meta string syntax, Go API, and CSS variables
- [CSS Custom Properties](/styling/css-custom-properties/) shows how to customize colors, spacing, and fonts with `--kz-*` CSS variables
- [File-Based Config](/reference/file-based-config/) sets engine options from a `kazari.config.yaml` file without writing Go code
