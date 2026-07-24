---
title: "Configuration Layers"
description: "Understand how Kazari resolves block settings through its three-layer cascade: engine defaults, language defaults, and per-block overrides."
sidebar:
  order: 3
---

Every code block Kazari renders gets its final settings from a three-layer cascade. Engine defaults establish the baseline, language defaults refine behavior for specific languages, and per-block overrides give each block the final say. Higher layers override lower ones, and omitted values fall through to the layer below.

## The cascade

```
Engine defaults
  └─ Language defaults
       └─ Per-block overrides (meta string or Options struct)
```

Kazari resolves five block properties through this cascade:

| Property | Meta string key | Options field | Description |
|---|---|---|---|
| Frame | `frame=code` | `Frame` | Visual frame style: `auto`, `code`, `terminal`, `none` |
| Line numbers | `showLineNumbers` | `LineNumbers` | Show or hide the line number gutter |
| Word wrap | `wrap` | `Wrap` | Enable word wrapping for long lines |
| Preserve indent | `preserveIndent` | `PreserveIndent` | Keep leading indentation on wrapped lines |
| Hanging indent | `hangingIndent=2` | `HangingIndent` | Extra indent (in `ch`) for continuation lines |

Feature toggles like `WithCopyButton` and `WithCollapsible` are engine-level only. They do not participate in the cascade and cannot be overridden per block.

## Layer 1: Engine defaults

Set during engine construction via functional options passed to `kazari.New()`. These apply to every code block the engine renders.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemes("github-light", "github-dark"),
    kazari.WithLineNumbers(true),
    kazari.WithDefaults(kazari.BlockDefaults{
        Wrap:  true,
        Frame: kazari.FrameCode,
    }),
)
```

`WithDefaults` sets block-level properties for all languages. Individual `With*` functions like `WithLineNumbers` modify the same defaults. Both approaches write to the same internal structure, so the last one in the option list wins.

### Built-in defaults

When no options are provided, the engine starts with these values:

| Property | Default |
|---|---|
| Frame | `FrameAuto` (auto-detect from language) |
| Line numbers | `false` |
| Word wrap | `false` |
| Preserve indent | `true` |
| Hanging indent | `0` (disabled) |

These come from `DefaultConfig()` and apply before any functional options run.

### File-based configuration

Engine defaults can also come from a `kazari.config.yaml` (or `.yml` or `.json`) file instead of Go code:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithConfigDir("."),
)
```

The config file maps directly to functional options. `WithConfigDir` searches for a config file and applies its settings as options during engine construction. Options passed after `WithConfigDir` override the file:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithConfigDir("."),        // loads config file
    kazari.WithCopyButton(false),     // overrides config file's copyButton
)
```

See [File-Based Config](/reference/file-based-config/) for the full schema and validation rules.

## Layer 2: Language defaults

Apply block properties to all blocks of a specific language. Set via `WithLanguageDefaults` or the `languageDefaults` key in a config file.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithLanguageDefaults(map[string]kazari.BlockDefaults{
        "bash,sh,zsh": {
            Wrap:  true,
            Frame: kazari.FrameTerminal,
        },
        "go": {
            LineNumbers: true,
        },
    }),
)
```

Comma-separated keys apply the same defaults to multiple languages. `"bash,sh,zsh"` sets all three to use terminal frames with word wrap.

The equivalent in a config file:

```yaml
languageDefaults:
  "bash,sh,zsh":
    wrap: true
    frame: "terminal"
  "go":
    lineNumbers: true
```

Language defaults override engine defaults for matching blocks. When a code block's language matches a key, all five properties from that entry replace the engine defaults for that block. Non-matching blocks fall through to engine defaults unchanged.

When multiple comma-separated keys could match the same language, keys are sorted alphabetically and the first match wins.

## Layer 3: Per-block overrides

Override any cascaded property for a single code block. Two interfaces provide per-block control: the meta string (for Markdown workflows) and the `Options` struct (for the Go API).

### Meta string

The meta string follows the language identifier in a fenced code block:

````
```go title="main.go" showLineNumbers frame=code {3-5}
func main() {
    cfg := loadConfig()
    svc := newService(cfg)
    svc.Run()
}
```
````

Cascade-relevant meta string keys:

| Key | Effect |
|---|---|
| `showLineNumbers` | Enable line numbers |
| `showLineNumbers=false` | Disable line numbers |
| `wrap` | Enable word wrap |
| `preserveIndent` | Enable preserved indentation on wrapped lines |
| `preserveIndent=false` | Disable preserved indentation |
| `hangingIndent=2` | Set hanging indent to 2ch |
| `frame=code` | Force editor frame |
| `frame=terminal` | Force terminal frame |
| `frame=none` | Remove frame |

The meta string also carries non-cascade properties like `title`, `theme`, markers (`{3-5}`), focus lines, and collapse directives. These are per-block only and do not participate in the cascade.

### Options struct

When calling `engine.Render()` directly, use the `Options` struct:

```go
ln := true
html, err := engine.Render(code, kazari.Options{
    Lang:        "go",
    Title:       "main.go",
    LineNumbers: &ln,
    LineMarkers: []kazari.LineMarker{
        {Type: kazari.MarkerMark, Lines: []kazari.Range{{Start: 3, End: 5}}},
    },
})
```

Cascade properties in `Options` use pointer types (`*bool`, `*int`, `*Frame`). A `nil` pointer means "use the cascaded value." A non-nil pointer overrides it. This makes the distinction between "not specified" and "explicitly set to the zero value" unambiguous.

| Field | Type | nil behavior |
|---|---|---|
| `Frame` | `*Frame` | Use language or engine default |
| `LineNumbers` | `*bool` | Use language or engine default |
| `Wrap` | `*bool` | Use language or engine default |
| `PreserveIndent` | `*bool` | Use language or engine default |
| `HangingIndent` | `*int` | Use language or engine default |
| `StartLineNumber` | `*int` | Default: 1 |

## Resolution example

Given this engine setup:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithDefaults(kazari.BlockDefaults{
        LineNumbers: true,
        Frame:       kazari.FrameCode,
    }),
    kazari.WithLanguageDefaults(map[string]kazari.BlockDefaults{
        "bash,sh": {
            Frame: kazari.FrameTerminal,
            Wrap:  true,
        },
    }),
)
```

Three code blocks resolve differently:

**Block A** (`` ```go ``)

| Property | Value | Source |
|---|---|---|
| Frame | `FrameCode` | Engine default |
| Line numbers | `true` | Engine default |
| Wrap | `false` | Built-in default |

**Block B** (`` ```bash ``)

| Property | Value | Source |
|---|---|---|
| Frame | `FrameTerminal` | Language default |
| Line numbers | `false` | Language default (replaces engine) |
| Wrap | `true` | Language default |

**Block C** (`` ```bash showLineNumbers frame=none ``)

| Property | Value | Source |
|---|---|---|
| Frame | `FrameNone` | Meta string |
| Line numbers | `true` | Meta string |
| Wrap | `true` | Language default |

Block C shows all three layers in action: `Wrap` comes from the language default, `Frame` and `LineNumbers` come from the meta string, and the engine defaults are fully overridden.

## What does not cascade

Engine-level features are set once during construction and apply uniformly. They cannot be overridden per language or per block.

| Feature | Option |
|---|---|
| Copy button | `WithCopyButton` |
| Fullscreen button | `WithFullscreenButton` |
| Wrap toggle button | `WithWrapButton` |
| Theme toggle | `WithThemeToggle` |
| Language badge | `WithLanguageBadge` |
| Frame detection | `WithFrameDetection` |
| File name extraction | `WithFileNameExtraction` |
| Style reset | `WithStyleReset` |
| Themed scrollbars | `WithThemedScrollbars` |
| Content exclusion | `WithContentExclusion` |
| Minimum contrast | `WithMinContrast` |
| Collapsible config | `WithCollapsible` |
| Links | `WithLinks` |
| File icons | `WithFileIcons` |
| Mermaid pass-through | `WithMermaidPassThrough` |
| Terminal comment stripping | `WithTerminalCommentStripping` |
| Output panel | `WithOutputPanel` |

## Quick reference

This page explains how the cascade works. For a complete list of every option, config file field, and meta string key in one place, see the Reference section:

- [Configuration Options](/reference/configuration-options/) lists every `With*` functional option and its corresponding config file field, type, and default value
- [Meta String Syntax](/reference/meta-string-syntax/) lists every key recognized in the fence info string with usage examples
- [File-Based Config](/reference/file-based-config/) covers the full YAML/JSON schema and validation rules
