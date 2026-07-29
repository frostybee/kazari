---
title: "Word Wrap"
description: "Wrap long lines while keeping wrapped continuations aligned with the code's indentation."
tags: [client-side-js]
sidebar:
  order: 7
---

Word wrap keeps long lines fully visible without horizontal scrolling. Wrapped continuations can align with the original indentation or use a fixed hanging indent, making it easy to tell new logical lines apart from continuations.

## Enable word wrap

Add `wrap` to the meta string:

````
```go title="wrap.go" showLineNumbers wrap
func configure(opts *Options) {
    opts.Logger = log.New(os.Stdout, "[kazari] prefix ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}
```
````

```go title="wrap.go" showLineNumbers wrap
func configure(opts *Options) {
    opts.Logger = log.New(os.Stdout, "[kazari] prefix ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}
```

→ Long lines wrap at the block edge. Continuation lines align with the first non-whitespace character on the original line because `preserveIndent` is on by default.

There is no `wrap=false` meta token. The `wrap` keyword is presence-only. To control wrap state at runtime, use the [wrap toggle button](#wrap-toggle-button) in the toolbar. To disable wrap per-block in Go, pass `Options.Wrap` with a `false` pointer value.

The [wrap toggle button](/features/toolbar/) in the toolbar also controls wrap at runtime.

Enable globally:

```go
kazari.WithDefaults(kazari.BlockDefaults{Wrap: true})
```

Or in the config file:

```yaml
defaults:
  wrap: true
```

## Disable indent preservation

With `preserveIndent=false`, wrapped continuations start at the left edge of the code area:

````
```go title="no-preserve.go" showLineNumbers wrap preserveIndent=false
func configure(opts *Options) {
    opts.Logger = log.New(os.Stdout, "[kazari] prefix ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}
```
````

```go title="no-preserve.go" showLineNumbers wrap preserveIndent=false
func configure(opts *Options) {
    opts.Logger = log.New(os.Stdout, "[kazari] prefix ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}
```

→ Wrapped continuations start at the left edge of the code area, ignoring the original indentation.

The bare `preserveIndent` token also exists. It is redundant against the `true` default, but useful to re-enable indent preservation per-block when a language default has set it to `false`.

## Hanging indent

`hangingIndent=N` adds N characters of extra indent to all wrapped continuations:

````
```go title="hanging.go" showLineNumbers wrap preserveIndent=false hangingIndent=4
func configure(opts *Options) {
    opts.Logger = log.New(os.Stdout, "[kazari] prefix ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}
```
````

```go title="hanging.go" showLineNumbers wrap preserveIndent=false hangingIndent=4
func configure(opts *Options) {
    opts.Logger = log.New(os.Stdout, "[kazari] prefix ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}
```

→ Wrapped continuations indent 4 characters from the left edge.

When `preserveIndent` is on, `hangingIndent` adds to the detected whitespace count. The formula is `HangingIndent + len(leadingWhitespace)`. A line indented 2 spaces with `hangingIndent=3` produces `--kz-indent:5ch`.

## Go API

All three wrap fields are available on the `Options` struct for per-block control:

```go
wrap, noPreserve, hanging := true, false, 4
html, err := engine.Render(code, kazari.Options{
    Lang:           "go",
    Title:          "hanging.go",
    Wrap:           &wrap,
    PreserveIndent: &noPreserve,
    HangingIndent:  &hanging,
})
```

All three use pointer semantics. A `nil` value inherits from the engine or language default. A non-nil value overrides it. `Options.Wrap` set to `&falseVal` is the way to force wrap off per-block in Go, even when an engine or language default has wrap enabled.

Engine-wide defaults:

```go
kazari.WithDefaults(kazari.BlockDefaults{
    Wrap:           true,
    PreserveIndent: true,  // default
    HangingIndent:  0,     // default
})
```

See [Configuration Layers](/getting-started/configuration-layers/) for the full cascade and override rules.

## Language defaults

Enable wrap for specific languages without affecting the rest:

```go
kazari.WithLanguageDefaults(map[string]kazari.BlockDefaults{
    "text,markdown,log": {Wrap: true},
})
```

Or in the config file:

```yaml
languageDefaults:
  text,markdown,log:
    wrap: true
```

This is the practical way to scope wrap to prose-heavy languages when the global default has wrap off. To disable wrap for specific languages when the global default is on, set `Wrap: false` in the language defaults for those languages.

## How it works

The renderer extracts leading whitespace from each line into a separate `<span class="indent">` element. That span keeps `white-space: pre` so tabs and spaces render exactly. The remaining code content wraps normally under `pre-wrap`.

A `--kz-indent` CSS variable is set per line based on the whitespace character count plus any `hangingIndent` value. When the computed indent is zero, no `--kz-indent` attribute is emitted. CSS uses a negative `text-indent` combined with matching `padding-inline-start` to create the hanging effect: the first visual line sits at the normal position while continuations are pushed right by `--kz-indent` characters.

```css
.kazari-block pre.wrap .kz-line .kz-code {
  white-space: pre-wrap;
  overflow-wrap: break-word;
  text-indent: calc(var(--kz-indent, 0ch) * -1);
  padding-inline-start: calc(var(--kz-code-padding-inline) + var(--kz-indent, 0ch));
}
```

The wrapping CSS rules live in `base.css` and are always included, regardless of whether `WithWrapButton` is enabled. The `WrapButton` option only controls the toolbar toggle UI.

Tabs are expanded to spaces before the indent count runs (controlled by `WithTabWidth`, default 2), so `--kz-indent` is always a clean character count regardless of source tab usage.

## Wrap toggle button

The wrap toggle button in the [toolbar](/features/toolbar/) gives readers runtime control over word wrap. See the [Toolbar](/features/toolbar/) page for full button details.

Word-wrap-specific behavior:

- The button's initial `aria-pressed` state matches the block's resolved `Wrap` value at render time. A block rendered with `wrap` in the meta string starts with `aria-pressed="true"` and the button in its active state.
- Toggling the button adds or removes the `wrap` class on `<pre>`. It does not recompute `--kz-indent` values. The inline indent styles are baked at render time and become inert when the `wrap` class is removed (no `pre.wrap` ancestor to activate the CSS rules). They become live again when the class is re-added.
- `frame=none` blocks never get a wrap toggle button. Only the floating copy button appears on frameless blocks. Wrap itself still works via meta string or config; there is no runtime toggle for it.
- The button reuses copy button CSS variables (`--kz-copy-fg`, `--kz-copy-bg-hover`, etc.). There are no dedicated `--kz-wrap-*` variables.

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `wrap` | Meta string | off | Enable wrap for this block (presence-only, no `wrap=false`) |
| `preserveIndent` / `preserveIndent=false` | Meta string | `true` | Align or reset wrapped-line indentation |
| `hangingIndent=N` | Meta string | `0` | Extra hanging indent in `ch` units |
| `Options.Wrap` / `PreserveIndent` / `HangingIndent` | Go API (per-block) | nil | Override any of the three per `Render()` call |
| `WithDefaults(BlockDefaults{...})` | Go API (engine) | `Wrap: false`, `PreserveIndent: true`, `HangingIndent: 0` | Engine-wide baseline |
| `WithLanguageDefaults(...)` | Go API (language) | none | Per-language overrides |
| `defaults.wrap` | Config file | `false` | Config file equivalent of `WithDefaults` |
| `wrapButton` | Config file / `WithWrapButton` | `true` | Show/hide the toolbar wrap toggle button |

## CSS variables

| Variable | Default | Description |
|---|---|---|
| `--kz-indent` | `0ch` (per-line) | Computed indent per line, set as an inline custom property by the renderer |
| `--kz-code-padding-inline` | `1.35rem` | Horizontal padding inside the code area (composes into the wrap formula) |

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- Wrap is purely visual. The copy button always copies the original unprocessed code.
- There is no `wrap=false` meta token, but `Options.Wrap` set to a `false` pointer value works as a per-block override in Go.
- `frame=none` blocks get no wrap toggle button, though wrap itself still works via meta string or config.
- `--kz-indent` is also used by collapsible section summary lines for alignment. The two computations are independent and do not interact.
- Markers and line numbers compose with wrap with no special-casing. Whitespace splitting runs before marker and token rendering.
- Tabs are expanded to spaces before the indent count, so `--kz-indent` is always a clean character count regardless of source tab usage.
