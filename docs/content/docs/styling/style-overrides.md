---
title: "Style Overrides"
description: "Programmatically set CSS custom property values from Go code."
tags: [css-variables, theming]
sidebar:
  order: 2
---

Style overrides inject CSS custom property values into the generated stylesheet. Two functions are available: `WithStyleOverrides` for values shared across both themes, and `WithThemedStyleOverrides` for values that differ between light and dark.

## Theme-neutral overrides

`kazari.WithStyleOverrides` sets the same value for both themes:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithStyleOverrides(map[string]string{
        "radius":    "0",
        "font-size": "14px",
    }),
)
```

Key normalization: bare names like `radius` are automatically prefixed to `--kz-radius`. Keys that already start with `--` are used as-is, whether they begin with `--kz-` or not.

| Key in map | Key in CSS output |
|---|---|
| `"radius"` | `--kz-radius` |
| `"--kz-shadow"` | `--kz-shadow` |
| `"--my-custom-var"` | `--my-custom-var` |

Universal values appear in both the light and dark CSS blocks. This ensures the value is never overridden by a theme-derived default that would otherwise win in the dark block.

Multiple `WithStyleOverrides` calls merge into the same internal map. When the same key appears in two calls, the later call's value is used.

## Themed overrides

`kazari.WithThemedStyleOverrides` sets separate values per theme:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemedStyleOverrides(map[string]kazari.StyleValue{
        "--kz-editor-bg": {Dark: "#1e293b", Light: "#ffffff"},
        "--kz-shadow":    {Dark: "none", Light: "0 2px 8px rgba(0,0,0,0.1)"},
    }),
)
```

The same key normalization rules apply. Light values appear in the `:root` block; dark values appear in the dark mode block. Setting only `Dark` or only `Light` is valid: the missing theme emits no line for that variable.

Both functions write to the same internal map and can be combined freely:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithStyleOverrides(map[string]string{
        "radius": "0",
    }),
    kazari.WithThemedStyleOverrides(map[string]kazari.StyleValue{
        "--kz-editor-bg": {Dark: "#1e293b", Light: "#f8fafc"},
    }),
)
```

## The StyleValue type

`StyleValue` carries one universal value or separate per-theme values:

```go
type StyleValue struct {
    Value string // both themes (theme-neutral)
    Dark  string // dark theme only
    Light string // light theme only
}
```

Helper methods:

| Method | Returns |
|---|---|
| `IsThemed()` | `true` if `Dark` or `Light` is non-empty |
| `LightValue()` | `Light` if themed, otherwise `Value` |
| `DarkValue()` | `Dark` if themed, otherwise `Value` |

An empty string for any field means no CSS line is emitted for that slot. This makes it safe to set only one theme without producing a blank property in the other.

## How overrides appear in CSS

Overrides are sorted alphabetically by key and emitted inside the root and dark blocks:

```css
/* theme-neutral override: appears in both blocks */
:root {
    --kz-font-size: 14px;
    --kz-radius: 0;
}

[data-theme="dark"] {
    --kz-font-size: 14px;  /* universal value repeated */
    --kz-radius: 0;
}
```

Themed overrides split between the two blocks:

```css
:root {
    --kz-editor-bg: #f8fafc;
    --kz-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

[data-theme="dark"] {
    --kz-editor-bg: #1e293b;
    --kz-shadow: none;
}
```

Alphabetical sorting makes the CSS output deterministic across builds.

## File config

Style overrides in `kazari.config.yaml` use the `styleOverrides` key with three value formats:

**Plain string** sets the same value for both themes:

```yaml
styleOverrides:
  radius: "0.5rem"
  toolbar-height: "2.5rem"
```

**Array** sets dark first, light second (index 0 is dark):

```yaml
styleOverrides:
  editor-bg: ["#1e293b", "#f8fafc"]
  editor-fg: ["#e2e8f0", "#1e293b"]
```

**Map** uses explicit `dark` and `light` keys; either may be omitted:

```yaml
styleOverrides:
  shadow:
    dark: "none"
    light: "0 2px 8px rgba(0,0,0,0.1)"
```

Key normalization applies the same rules as the Go API. See [File-Based Config](/reference/file-based-config/) for the full config file format.

## Precedence

Style overrides are placed after Kazari's theme-derived defaults in the generated CSS, so overrides always win within the stylesheet. When `WithCascadeLayer` is active (default: `@layer kazari`), unlayered consumer CSS automatically takes precedence over layered overrides without `!important`. When layering is disabled, consumer overrides need enough specificity to win.

Go API and file config overrides write to the same internal map. The last call wins for the same key: a `WithStyleOverrides` call placed after `WithConfigDir` in `kazari.New()` overrides the file config value for that key.

## Configuration

| Option | Default | Description |
|---|---|---|
| `WithStyleOverrides(map[string]string)` | `nil` | Theme-neutral CSS variable overrides |
| `WithThemedStyleOverrides(map[string]StyleValue)` | `nil` | Per-theme CSS variable overrides |
| `styleOverrides` (config file) | `nil` | YAML/JSON equivalent; accepts plain string, array, or map |

For the full list of `--kz-*` variables that can be overridden, see the [CSS Variables](/reference/css-variables/) reference.
