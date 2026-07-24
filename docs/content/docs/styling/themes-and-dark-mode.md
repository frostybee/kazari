---
title: "Themes & Dark Mode"
description: "Configure light and dark themes, dark mode strategies, and per-block theme overrides."
sidebar:
  order: 3
---

Both light and dark token colors are baked into the HTML at build time. Theme switching is pure CSS with no JavaScript, no stylesheet swap, and no re-render.

## Setting themes

Pass light and dark theme names to `WithThemes`. These names are forwarded to the highlighter adapter.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemes("one-light", "one-dark-pro"),
)
```

Default: `"github-light"` and `"github-dark"`.

The Nuri adapter accepts Shiki-compatible theme names. The Chroma adapter requires a style map to translate Kazari theme names to Chroma style names. See [Chroma Adapter](/integrations/chroma-adapter/) for details.

For single-theme rendering, pass an empty string as the dark theme:

```go
kazari.WithThemes("github-light", "")
```

:::tip
Pass an empty string as the dark theme to disable dual-theme rendering entirely. Only light-theme token colors are baked into the HTML.
:::

## Dark mode strategies

Set the dark mode strategy with `kazari.WithDarkMode()`. Three strategies are available:

| Strategy | Constructor | Generated CSS |
|----------|-------------|---------------|
| Selector | `SelectorMode(selector)` | `{root}{selector} { ... }` |
| Media query | `MediaQueryMode()` | `@media (prefers-color-scheme: dark) { ... }` |
| Both | `BothMode(selector)` | Selector block and media query block |

Default: `SelectorMode(".dark")`. The selector parameter accepts any CSS selector that can be appended to `:root`, including class selectors (`.dark`) and attribute selectors (`[data-theme="dark"]`).

```go
// Match a custom data attribute
kazari.WithDarkMode(kazari.SelectorMode("[data-theme='dark']"))

// Respect the OS preference only
kazari.WithDarkMode(kazari.MediaQueryMode())

// Match both a class toggle and the OS preference
kazari.WithDarkMode(kazari.BothMode(".dark"))
```

In a config file, set `darkMode.kind` to `"selector"`, `"mediaQuery"`, or `"both"`, and provide `darkMode.selector` when the kind is `"selector"` or `"both"`.

:::caution[Selector must match your site's toggle]
The `darkMode.selector` value must exactly match how your site activates dark mode. Kazari bakes the selector into the generated CSS. If the selector doesn't match, dark theme styles are scoped to a rule that never applies, and code blocks silently stay light.

Common selectors:

| Site / SSG | Dark mode mechanism | Selector value |
|------------|---------------------|----------------|
| Sarde | `data-theme="dark"` on `<html>` | `'[data-theme="dark"]'` |
| Hugo | `.dark` class on `<html>` (typical) | `".dark"` |
| OS preference only | `prefers-color-scheme` media query | Use `MediaQueryMode()` instead |

Inspect your site's `<html>` element after toggling dark mode to confirm which attribute or class is set.
:::

## Per-block theme override

The `theme=` meta token overrides the theme for a single code block without affecting the engine defaults.

A single name applies that theme for both light and dark contexts:

````
```go theme="dracula"
````

A comma-separated pair sets different themes per context. The first name is the light theme; the second is the dark theme:

````
```go theme="github-light,github-dark-default"
````

Override colors are embedded as inline CSS custom properties on the block wrapper element (`--kz-ovl-*` for light, `--kz-ovd-*` for dark). CSS switching rules remap these to the standard `--kz-*` names without JavaScript. Override results are cached per unique `theme=` string value.

## Theme toggle button

`WithThemeToggle(true)` adds a per-block button that forces the block into light or dark mode, independent of the page-level theme:

```go
kazari.WithThemeToggle(true)
```

Default: `false`. The button sets a `data-kz-theme` attribute on the block, which CSS rules use to force one color direction. Config file key: `themeToggleButton`.

The toggle requires a dark theme to be configured via `WithThemes`. Without a dark theme, there is nothing to toggle to.

Toggle state is persisted to `localStorage` as `kz-theme-<blockId>`, where the block ID is a FNV-1a hash of the code content. Two blocks with identical source code share the same ID and localStorage state. Clicking the toggle a second time clears the attribute, reverting the block to follow the page theme.

## Theme customizer

`WithThemeCustomizer` receives the theme name and extracted colors after adjustments are applied, and returns the modified colors:

```go
kazari.WithThemeCustomizer(func(name string, colors kazari.ThemeInfo) kazari.ThemeInfo {
    if name == "github-dark" {
        colors.BG = "#0d1117"
    }
    return colors
})
```

The customizer runs before CSS generation. It is called at engine construction and again for each per-block `theme=` override. This option is Go API only and has no config file equivalent.

## Theme adjustments (OKLCH)

`WithThemeAdjustments` tints extracted editor chrome colors in the OKLCH color space. This shifts backgrounds, gutters, and selection colors toward a target hue without affecting individual token colors.

`Hue` and `Chroma` are pointer fields. Pass `nil` to leave the corresponding axis unchanged.

```go
floatPtr := func(v float64) *float64 { return &v }

kazari.WithThemeAdjustments(kazari.ThemeAdjustments{
    Hue:    floatPtr(220),  // shift toward blue
    Chroma: floatPtr(0.03), // subtle saturation
})
```

`Targets` is a bitmask controlling which color roles are adjusted:

| Constant | Colors affected |
|----------|----------------|
| `AdjustBackgrounds` (default) | `BG`, `SelectionBG`, `FoldBG` |
| `AdjustForegrounds` | `FG`, `LineNumberFG` |

The zero value of `Targets` is `AdjustBackgrounds`. This option is Go API only and has no config file equivalent.

## Themed scrollbars and selection

`WithThemedScrollbars` makes the scrollbar thumb color adapt to the active theme:

```go
kazari.WithThemedScrollbars(true)
```

Default: `true`. Config file key: `themedScrollbars`.

`WithThemedSelectionColors` applies the theme's `SelectionBG` to text selections inside code blocks:

```go
kazari.WithThemedSelectionColors(true)
```

Default: `false`. Config file key: `themedSelection`.

## Minimum contrast

`WithMinContrast` sets a WCAG contrast ratio floor for token foreground colors against the editor background:

```go
kazari.WithMinContrast(4.5)
```

Default: `5.5`. Valid range: `0` to `21`. Set to `0` to disable contrast adjustment entirely. Config file key: `minContrast`.

## Multi-engine pages

When multiple engines share a single page (for example, a showcase app with different theme configurations), use `engine.ThemeCSS()` for secondary engines instead of `engine.CSS()`:

```go
primaryCSS := primaryEngine.CSS()    // full CSS once
secondaryCSS := otherEngine.ThemeCSS() // theme rules only
```

`ThemeCSS()` emits only theme variables, token-switching rules, and toggle CSS. It does not include the structural CSS (base, toolbar, frame, copy, fullscreen) that the primary engine already provides. Using `CSS()` from a secondary engine duplicates structural CSS and can cause conflicts.

## No-flash guarantee

Both light and dark token colors are baked into the HTML at render time. Theme switching is a CSS operation: no JavaScript, no stylesheet swap, no re-render. `base.css` enforces `transition: none !important` on all code block elements to prevent color interpolation during toggles or page load.

See [Dual-Theme Rendering](/architecture/dual-theme-rendering/) for the full technical explanation.

## Options reference

| Option | Default | Config key | Notes |
|--------|---------|------------|-------|
| `WithThemes(light, dark)` | `"github-light"`, `"github-dark"` | `themes.light`, `themes.dark` | Pass `""` as dark for single-theme |
| `WithDarkMode(dm)` | `SelectorMode(".dark")` | `darkMode.kind`, `darkMode.selector` | Three strategies |
| `WithThemeToggle(bool)` | `false` | `themeToggleButton` | Per-block light/dark toggle button |
| `WithThemeCustomizer(fn)` | `nil` | Go API only | Post-process extracted colors |
| `WithThemeAdjustments(adj)` | `nil` | Go API only | OKLCH hue/chroma tinting of editor chrome |
| `WithThemedScrollbars(bool)` | `true` | `themedScrollbars` | Themed scrollbar colors |
| `WithThemedSelectionColors(bool)` | `false` | `themedSelection` | Themed text selection background |
| `WithMinContrast(ratio)` | `5.5` | `minContrast` | WCAG contrast ratio, 0 to 21 |
