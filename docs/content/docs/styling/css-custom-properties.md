---
title: "CSS Custom Properties"
description: "How Kazari uses CSS custom properties for styling, and how cascade layers and style resets protect code blocks."
sidebar:
  order: 1
---

Kazari's presentation layer is styled through CSS custom properties with the `--kz-` prefix. Override any property in a stylesheet to customize appearance without rebuilding the engine.

## The `--kz-` naming convention

Override any property by redeclaring it in your stylesheet at `:root` or a more specific selector:

```css
:root {
    --kz-radius: 0;
    --kz-font-size: 14px;
}
```

All CSS custom properties use the `--kz-` prefix (from **K**a**z**ari). HTML classes use the `kz-` prefix (for example, `.kz-block`, `.kz-line`). A consumer stylesheet targeting `:root` or `.kazari-block` takes precedence over Kazari's own declarations when no cascade layer is involved.

For the full list of available properties, see [CSS Variables](/reference/css-variables/).

## How overrides work

Kazari generates all `--kz-*` declarations in Go code and writes them into `:root` (or the configured CSS root selector) as part of `engine.CSS()`. The 24 static CSS files that ship with Kazari only *consume* these properties via `var(--kz-*)`. They never declare defaults of their own.

This architecture means consumer overrides always win. Kazari's generated CSS is wrapped in `@layer kazari` by default, and unlayered consumer stylesheets take precedence over layered declarations without needing `!important`.

A few variables are set inline per block during rendering:

- `--kz-ln-width` on the `<code>` element, sized to fit the widest line number (e.g., `3ch` for a 100-line block)
- `--kz-indent` on individual lines, preserving indentation for wrapped text

These inline declarations take highest specificity and cannot be overridden from a stylesheet.

## Cascade layer

`kazari.WithCascadeLayer(name)` wraps all generated Kazari CSS in `@layer <name> { ... }`.

The default layer name is `"kazari"`. CSS is layered by default. Unlayered consumer stylesheets automatically take precedence over layered Kazari styles without requiring `!important`.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithCascadeLayer("kazari"), // default; already applied
)
```

To disable layering, pass an empty string:

```go
kazari.WithCascadeLayer("")
```

With layering disabled, specificity rules apply normally. Consumer overrides need enough specificity to win.

Config file key: `cascadeLayer`

## Style reset

`kazari.WithStyleReset(bool)` controls whether `all: revert` is applied to the `<figure>` element inside `.kazari-block`. The default is `true`.

Without this reset, site-wide CSS (prose margins, padding from a CMS stylesheet, inherited borders) can bleed into code blocks. With it enabled, Kazari reverts the figure to browser defaults and then re-applies only its own properties: `position: relative`, `margin: 0`, `border-radius`, `box-shadow`, and `border`.

```go
kazari.WithStyleReset(false) // disable if your site has no prose bleed
```

Config file key: `styleReset`

## Content exclusion

Every rendered code block includes the `not-content` class on its wrapper element. This prevents documentation-site prose styles from leaking into Kazari's code block styling. The class is always present and requires no configuration.

`WithContentExclusion` and the `contentExclusion` config file key exist for backward compatibility but are no-ops. See [Content Exclusion](/features/content-exclusion/) for details.

## Custom CSS root

`kazari.WithThemeCSSRoot(selector)` changes the CSS selector under which all `--kz-*` variable declarations are placed. The default is `":root"`.

Use this when scoping Kazari variables to a specific container instead of the document root:

```go
kazari.WithThemeCSSRoot(".my-docs")
```

The generated CSS uses that selector for the light theme block. The dark mode selector composes by direct concatenation: `.my-docs` combined with `.dark` produces `.my-docs.dark { ... }`. No space is inserted between the root selector and the dark selector.

Config file key: `themeCSSRoot`

## No-flash contract

`base.css` applies `transition: none !important` to `.kazari-block`, `*`, `*::before`, and `*::after`. This rule covers all descendants. Code blocks never flash or transition colors during theme switches or page load.

:::danger
Do not add CSS transitions to any element inside `.kazari-block`. The `transition: none !important` rule in `base.css` enforces the no-flash contract. Adding transitions causes visible color flashing during theme switches.
:::

See [Themes and Dark Mode](/styling/themes-and-dark-mode/) for the full dual-theme rendering model and the no-flash guarantees it provides.

## Configuration

| Option | Default | Config key | Description |
|--------|---------|------------|-------------|
| `WithCascadeLayer(name)` | `"kazari"` | `cascadeLayer` | Wrap all CSS in `@layer`. Pass `""` to disable. |
| `WithStyleReset(bool)` | `true` | `styleReset` | Apply `all: revert` to the `<figure>` element. |
| `WithThemeCSSRoot(selector)` | `":root"` | `themeCSSRoot` | Custom selector for `--kz-*` declarations. |
