---
title: "Syntax Highlighting"
description: "Language-aware code coloring powered by a pluggable highlighter."
tags: [highlighter]
sidebar:
  order: 0
---

Kazari renders fenced code blocks with language-aware syntax coloring. A pluggable highlighter ([Nuri](/integrations/nuri-adapter/) or [Chroma](/integrations/chroma-adapter/)) tokenizes the source code, and Kazari handles all presentation: frames, line numbers, markers, copy buttons, and theme switching.

## Basic usage

Set the language identifier as the first word after the opening backticks:

````
```go
func main() {
    fmt.Println("Hello, world!")
}
```
````

The language identifier controls both syntax coloring and frame auto-detection. Additional per-block options (title, line numbers, markers) follow the language in the same line. See [Meta String Syntax](/reference/meta-string-syntax/) for the full list.

````
```go title="main.go" showLineNumbers {3}
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```
````

Pass code through the Go API with the same options:

```go
html, err := engine.Render(code, kazari.Options{
    Lang: "go",
})
```

Or use a meta string directly:

```go
html, err := engine.RenderWithMeta(code, `go title="main.go" showLineNumbers {3}`)
```

## Supported languages

Kazari delegates tokenization to whichever `Highlighter` is configured on the engine. Two adapters are provided:

| Adapter | Languages | Grammar type |
|---|---|---|
| [Nuri](/integrations/nuri-adapter/) | 39 (core bundle) / 258 (full bundle) | VS Code TextMate grammars |
| [Chroma](/integrations/chroma-adapter/) | 400+ | Regex-based lexers |

Nuri uses the same TextMate grammars as Shiki. See the [full list of supported languages and themes](https://github.com/shikijs/textmate-grammars-themes/blob/main/packages/tm-grammars/README.md) on GitHub. See each adapter page for setup instructions and configuration.

## Language aliases

Map shorthand names to canonical language identifiers so code fences resolve correctly. Aliases are case-insensitive.

```yaml
languageAliases:
  ts: typescript
  js: javascript
  py: python
  rs: rust
```

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithLanguageAliases(map[string]string{
        "ts": "typescript",
        "js": "javascript",
    }),
)
```

When a code block uses `ts` as its language, the engine resolves it to `typescript` before passing it to the highlighter. See [Configuration Options](/reference/configuration-options/) and [File-Based Config](/reference/file-based-config/) for the full option reference.

## Unknown languages

When the language identifier is not recognized by the configured highlighter, Kazari falls back to plaintext rendering (no coloring) and emits a warning through the engine's warning handler. The block still renders with all presentation features (frame, copy button, line numbers) intact.

## Terminal output

For rendering terminal output with ANSI escape sequences, set the language to `ansi`. Kazari bypasses the external highlighter and parses SGR escape codes directly into colored, styled tokens. See [ANSI Rendering](/features/ansi-rendering/) for supported sequences, color modes, and the visual showcase.
