---
title: "Goldmark Extension"
description: "Wire Kazari into a Goldmark Markdown pipeline for fenced code blocks and code groups."
tags: [markdown]
sidebar:
  order: 3
---

The `kazarimd` package integrates Kazari into a [Goldmark](https://github.com/yuin/goldmark) Markdown pipeline, replacing the default fenced code block renderer and adding `:::code-group` container support. All meta string features work automatically through the fence info string.

Import path: `github.com/frostybee/kazari/goldmark`

## Install

```bash
go get github.com/frostybee/kazari/goldmark@latest
```

## Fenced code blocks

Register `kazarimd.New(engine)` as a Goldmark extension to route all fenced code blocks through Kazari:

```go
import (
    "github.com/frostybee/kazari"
    kazarimd "github.com/frostybee/kazari/goldmark"
    "github.com/yuin/goldmark"
)

md := goldmark.New(
    goldmark.WithExtensions(kazarimd.New(engine)),
)

var buf bytes.Buffer
if err := md.Convert([]byte(markdownSource), &buf); err != nil {
    log.Fatal(err)
}
```

The extension reads the fence info string (everything after the opening backticks), splits it into language and meta, and calls `engine.RenderWithMeta(code, meta)`. The language is the first word; everything after the first space is the meta string.

````markdown
```go title="main.go" {3-5} showLineNumbers
````

In this example, `go` is the language and `title="main.go" {3-5} showLineNumbers` is the meta string passed to the engine. All meta string features work automatically: `title=`, line markers, `showLineNumbers`, `focus=`, `collapse`, and the rest.

## Code groups

Register `kazarimd.CodeGroups(engine)` to add `:::code-group` container support:

```go
md := goldmark.New(
    goldmark.WithExtensions(
        kazarimd.New(engine),
        kazarimd.CodeGroups(engine),
    ),
)
```

`CodeGroups(engine)` calls `engine.EnableCodeGroups()` internally, which enables code group CSS and JS in the engine output. No separate setup step is needed.

Register `CodeGroups(engine)` during setup, before the engine is used for rendering. `EnableCodeGroups()` mutates the engine configuration without synchronization, so calling it concurrently with active `Convert()` calls is not safe.

Tab labels are derived in priority order: explicit `title=` meta attribute, then extracted filename comment from the first line of code, then capitalized language name, then `"Code"` as a fallback. See the [Code Groups](/features/code-groups/) page for the full `:::code-group` syntax, tab behavior, `sync=` attribute, and keyboard navigation.

## CSS and JS injection

:::important
The Goldmark extensions produce per-block HTML only. Call `engine.CSS()` and `engine.JS()` separately and inject them once per rendered page. Without this step, code blocks render unstyled and without interactive features.
:::

`CSS()` and `JS()` are config-derived and idempotent. The output is the same regardless of which blocks were rendered. Call once per site build, not per page.

When `CodeGroups(engine)` is registered before `engine.CSS()` is called, the CSS output includes tab bar styles and the JS output includes the tab-switching and sync handler.

### Inline assets

```go
var buf bytes.Buffer
md.Convert(markdownSource, &buf)

fmt.Fprintf(w, "<style>%s</style>", engine.CSS())
// ... converted markdown output ...
fmt.Fprintf(w, "<script>%s</script>", engine.JS())
```

### External files (SSG pattern)

Write assets to disk once during the build, then reference them from page templates:

```go
os.WriteFile("assets/kazari.css", []byte(engine.CSS()), 0644)
os.WriteFile("assets/kazari.js", []byte(engine.JS()), 0644)
```

```html
<link rel="stylesheet" href="/assets/kazari.css">
<script type="module" src="/assets/kazari.js"></script>
```

### Hashed assets

`engine.Assets()` returns CSS and JS with content-hashed filenames for cache busting:

```go
assets := engine.Assets()
os.WriteFile("assets/"+assets.CSS.Filename, []byte(assets.CSS.Content), 0644)
os.WriteFile("assets/"+assets.JS.Filename, []byte(assets.JS.Content), 0644)
```

Each `AssetFile` includes `Content`, `Hash` (8-char hex FNV-1a), and `Filename` (e.g., `kazari-a1b2c3d4.css`).

## Full example

```go
package main

import (
    "bytes"
    "context"
    "fmt"
    "log"

    "github.com/frostybee/kazari"
    kazarinuri "github.com/frostybee/kazari/nuri"
    kazarimd "github.com/frostybee/kazari/goldmark"
    "github.com/frostybee/nuri"
    "github.com/frostybee/nuri/bundle/core"
    "github.com/yuin/goldmark"
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
        kazari.WithCopyButton(true),
    )

    md := goldmark.New(
        goldmark.WithExtensions(
            kazarimd.New(engine),
            kazarimd.CodeGroups(engine),
        ),
    )

    src := []byte("# Hello\n\n```go title=\"main.go\"\nfmt.Println(\"hello\")\n```\n")
    var buf bytes.Buffer
    if err := md.Convert(src, &buf); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("<html><head><style>%s</style></head><body>%s<script>%s</script></body></html>",
        engine.CSS(), buf.String(), engine.JS())
}
```

## Single engine pattern

When using the Goldmark extension alongside direct `engine.Render()` or `engine.RenderWithMeta()` calls, use a single engine instance for everything. Each engine produces its own CSS and JS output. Including output from two separate engines in the same page causes duplicate function declarations.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemes("github-light", "github-dark"),
)

md := goldmark.New(goldmark.WithExtensions(
    kazarimd.New(engine),
    kazarimd.CodeGroups(engine),
))
var buf bytes.Buffer
md.Convert(markdownSource, &buf)

customBlock, _ := engine.Render(code, kazari.Options{Lang: "go", Title: "example.go"})

css := engine.CSS()
js := engine.JS()
```

One `engine.CSS()` and `engine.JS()` call covers both the Goldmark-rendered blocks and the directly rendered blocks.

## Edge cases

| Case | Behavior |
|------|----------|
| No code blocks in markdown | Zero Kazari overhead; extension is dormant |
| Code block with no language | Rendered as plaintext with editor frame |
| Code block with unknown language | Plaintext fallback with a warning |
| Empty code block | Renders empty figure with frame |
| Code group with single block | Renders as a group with one tab |
| Code group with no code blocks | Container omitted entirely |
| Nested `:::code-group` | Not supported; inner directive treated as text |
| Code block outside code group | Rendered as standalone block |
| Mixed standalone and grouped | Both render correctly, independently |
