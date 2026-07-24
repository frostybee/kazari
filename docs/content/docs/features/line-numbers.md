---
title: "Line Numbers"
description: "Display line numbers with configurable start number and auto-width gutter."
sidebar:
  order: 3
---

Line numbers appear in a gutter column to the left of the code. They are off by default, excluded from text selection and clipboard copy, and marked `aria-hidden` for screen readers.

## Enable line numbers

Add `showLineNumbers` to the meta string:

````
```go title="main.go" showLineNumbers
func main() {
    cfg := loadConfig()
    svc := newService(cfg)
    svc.Run()
}
```
````

```go title="main.go" showLineNumbers
func main() {
    cfg := loadConfig()
    svc := newService(cfg)
    svc.Run()
}
```

→ Line numbers 1 through 5 appear in a gutter column to the left of the code.

To disable line numbers on a specific block when the engine default is on, use `showLineNumbers=false`.

Enable globally with the Go API:

```go
kazari.WithLineNumbers(true)
```

Or in the config file:

```yaml
lineNumbers: true
```

`WithLineNumbers(true)` is shorthand for `WithDefaults(BlockDefaults{LineNumbers: true})`. Both write to the same field.

## Custom start number

Change where numbering begins with `startLineNumber=N`:

````
```go title="calc.go (lines 22-24)" showLineNumbers startLineNumber=22
func add(a, b int) int {
    return a + b
}
```
````

```go title="calc.go (lines 22-24)" showLineNumbers startLineNumber=22
func add(a, b int) int {
    return a + b
}
```

→ Line numbers start at 22, with the gutter auto-sized for two-digit numbers.

Any integer is accepted, including zero and negative values. The first line displays the start number, and each subsequent line increments by 1.

There is no config file equivalent for `startLineNumber`. It is always set per block via the meta string or `Options.StartLineNumber`.

## Gutter width

The gutter auto-sizes based on the digit count of the highest line number. For 1- and 2-digit ranges, the default width of `2ch` applies. For 3+ digits (starting at line 100 or a block with 100+ lines), the width expands automatically via an inline `--kz-ln-width` style on the `<code>` element.

The minus sign counts as one character for negative start numbers. No JavaScript is involved; the width is computed at render time.

## Per-language defaults

Set line numbers on by default for specific languages with `WithLanguageDefaults`:

```go
kazari.WithLanguageDefaults(map[string]kazari.BlockDefaults{
    "go,typescript": {LineNumbers: true},
    "bash,sh,zsh":   {LineNumbers: false},
})
```

```yaml
languageDefaults:
  "go,typescript":
    lineNumbers: true
  "bash,sh,zsh":
    lineNumbers: false
```

Language keys are comma-separated and matched case-insensitively. Per-block meta always overrides language defaults. See [Configuration Layers](/getting-started/configuration-layers/) for the full cascade.

## On highlighted lines

Line numbers on highlighted lines (mark, ins, del) switch to a brighter color and higher opacity to remain readable against the marker background:

```go showLineNumbers {2-3}
func main() {
    cfg := loadConfig()
    svc := newService(cfg)
    svc.Run()
}
```

→ Lines 2 and 3 have brighter line numbers than the unmarked lines. These are controlled by `--kz-ln-highlight-fg` and `--kz-ln-highlight-opacity`.

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `showLineNumbers` | Meta string | off | Enable line numbers for this block |
| `showLineNumbers=false` | Meta string | n/a | Disable line numbers for this block |
| `startLineNumber=N` | Meta string | `1` | Starting line number (any integer) |
| `WithLineNumbers(bool)` | Go API | `false` | Engine-wide default |
| `WithLanguageDefaults(...)` | Go API | none | Per-language defaults |

Config file keys: `lineNumbers` (engine-wide), `defaults.lineNumbers`, `languageDefaults.<lang>.lineNumbers`.

## CSS variables

| Variable | Default | Description |
|---|---|---|
| `--kz-ln-width` | `2ch` | Minimum number width (auto-expanded for 3+ digits) |
| `--kz-ln-padding-inline` | `2ch` | Horizontal padding around numbers |
| `--kz-ln-fg` | Theme-derived | Number text color |
| `--kz-ln-opacity` | `1` | Number opacity |
| `--kz-ln-highlight-fg` | Theme-derived | Number color on highlighted lines |
| `--kz-ln-highlight-opacity` | `0.8` | Number opacity on highlighted lines |
| `--kz-gutter-border-width` | `1px` | Gutter right border width |
| `--kz-gutter-border-color` | Theme-derived | Gutter right border color |

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- Line numbers work in all frame types: editor, terminal, and frameless.
- ANSI blocks (`lang="ansi"`) support line numbers.
- Collapsible sections preserve line numbers; they are never renumbered. Summary and gap rows render empty gutter placeholders to maintain alignment.
- `startLineNumber` is per-block only; there is no engine-wide or config file setting for it.
