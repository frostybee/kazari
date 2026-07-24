---
title: "Configuration Options"
description: "Complete reference for all kazari.With* functional options."
sidebar:
  order: 2
---

Functional options configure the engine via `kazari.New()`, applied left to right (last writer wins for scalar fields). Most options have a `kazari.config.yaml` equivalent; see [File-Based Config](/reference/file-based-config/) for the file format.

## Highlighter and themes

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithHighlighter(Highlighter)` | `nil` | Go API only | Tokenization backend; without one, render falls back to plain text |
| `WithThemes(light, dark string)` | `"github-light"`, `"github-dark"` | `themes.light`, `themes.dark` | Light and dark theme names forwarded to the highlighter |
| `WithDarkMode(DarkMode)` | `SelectorMode(".dark")` | `darkMode` | Dark mode CSS strategy |
| `WithThemeCustomizer(func(string, ThemeInfo) ThemeInfo)` | `nil` | Go API only | Post-process extracted theme colors; runs after adjustments |
| `WithThemeAdjustments(ThemeAdjustments)` | `nil` | Go API only | Tint editor chrome colors in OKLCH space; runs before customizer |
| `WithThemeCSSRoot(string)` | `":root"` | `themeCSSRoot` | CSS selector scope for theme variable declarations |
| `WithThemedScrollbars(bool)` | `true` | `themedScrollbars` | Apply theme-derived scrollbar colors |
| `WithThemedSelectionColors(bool)` | `false` | `themedSelection` | Apply theme `SelectionBG` to text selection |
| `WithMinContrast(float64)` | `5.5` | `minContrast` | WCAG contrast ratio floor for token colors on marker backgrounds (0 to 21) |

`WithDarkMode` accepts a value from one of three constructors: `SelectorMode(selector)`, `MediaQueryMode()`, or `BothMode(selector)`. The selector parameter accepts any CSS selector appendable to `:root`, such as class selectors (`.dark`) or attribute selectors (`[data-theme="dark"]`). See [Themes & Dark Mode](/styling/themes-and-dark-mode/) for details on each strategy.

### `ThemeAdjustments` fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Hue` | `*float64` | `nil` | Target hue in degrees (0 to 360); nil leaves hue unchanged |
| `Chroma` | `*float64` | `nil` | Target chroma (0 to 0.4 typical); nil leaves chroma unchanged |
| `Targets` | `AdjustTargets` | `AdjustBackgrounds` | Bitmask: `AdjustBackgrounds` (BG, SelectionBG, FoldBG) and/or `AdjustForegrounds` (FG, LineNumberFG) |

## Toolbar

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithCopyButton(bool)` | `true` | `copyButton` | Copy-to-clipboard button |
| `WithFullscreenButton(bool)` | `true` | `fullscreenButton` | Fullscreen expand button with font size controls |
| `WithWrapButton(bool)` | `true` | `wrapButton` | Word wrap toggle button |
| `WithThemeToggle(bool)` | `false` | `themeToggleButton` | Per-block light/dark theme toggle |

## Frames and titles

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithFrameDetection(bool)` | `true` | `frameDetection` | Auto-detect frame type from language and content |
| `WithFileNameExtraction(bool)` | `true` | `fileNameExtraction` | Extract title from first-line filename comment |
| `WithTerminalDotStyle(TerminalDotStyle)` | `DotsColored` | `terminalDotStyle` | Terminal frame window dots: `DotsColored` or `DotsMinimal` |
| `WithTerminalCommentStripping(bool)` | `true` | `terminalCommentStripping` | Strip shell comments from copy text in terminal frames |
| `WithLanguageBadge(bool)` | `true` | `languageBadge` | Language name badge in the toolbar |
| `WithLanguageIconMode(LangIconMode)` | `LangIconNone` | `languageIconMode` | Language icon display: `LangIconNone`, `LangIconOnly`, or `LangIconAndText` |
| `WithFileIcons(bool)` | `true` | `fileIcons` | File extension icon in the title bar |
| `WithFileIconResolver(func(ext string) string)` | `nil` | Go API only | Custom function returning SVG data URL per file extension |

## Content features

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithLineNumbers(bool)` | `false` | `lineNumbers` | Show line numbers by default |
| `WithMermaidPassThrough(bool)` | `true` | `mermaidPassThrough` | Pass `mermaid` blocks through as `<pre class="mermaid">` |
| `WithLinks(bool)` | `false` | `links` | Enable `@[text](url)` inline link annotations |
| `WithDataLineCount(bool)` | `true` | `dataLineCount` | Emit `data-line-count` attribute on the block wrapper |
| `WithContentExclusion(bool)` | `false` | `contentExclusion` | Deprecated no-op. The `not-content` class is always present on rendered blocks regardless of this setting. |

The `data-line-count` attribute enables CSS targeting based on block length (e.g., hiding line numbers on single-line blocks) and JavaScript that adjusts behavior based on line count.

`WithLineNumbers(bool)` is a convenience shortcut that sets `Defaults.LineNumbers`. It writes to the same field as `WithDefaults`; the last call wins.

## CSS and appearance

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithStyleReset(bool)` | `true` | `styleReset` | Apply `all: revert` style isolation on code blocks |
| `WithCascadeLayer(string)` | `"kazari"` | `cascadeLayer` | `@layer` name for generated CSS; empty string disables the layer |
| `WithTabWidth(int)` | `2` | `tabWidth` | Tab-to-space expansion width (>= 1) |
| `WithStyleOverrides(map[string]string)` | `nil` | `styleOverrides` | CSS variable overrides applied to both themes |
| `WithThemedStyleOverrides(map[string]StyleValue)` | `nil` | `styleOverrides` | CSS variable overrides with separate light/dark values |

Both `WithStyleOverrides` and `WithThemedStyleOverrides` map to the `styleOverrides` YAML key. Bare key names (e.g., `radius`) are automatically prefixed to `--kz-radius`; keys starting with `--` are used as-is. See [Style Overrides](/styling/style-overrides/) for per-theme value formats.

## Block defaults

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithDefaults(BlockDefaults)` | (see fields below) | `defaults` | Engine-wide per-block rendering defaults |
| `WithLanguageDefaults(map[string]BlockDefaults)` | `nil` | `languageDefaults` | Per-language rendering defaults |
| `WithLanguageAliases(map[string]string)` | `nil` | `languageAliases` | Map alias names to canonical language names |

### `BlockDefaults` fields

| Field | Type | Default | YAML key | Description |
|-------|------|---------|----------|-------------|
| `Wrap` | `bool` | `false` | `wrap` | Enable word wrap |
| `PreserveIndent` | `bool` | `true` | `preserveIndent` | Preserve leading indentation on wrapped lines |
| `HangingIndent` | `int` | `0` | `hangingIndent` | Extra indentation columns for wrapped continuation lines (>= 0) |
| `LineNumbers` | `bool` | `false` | `lineNumbers` | Show line numbers |
| `Frame` | `Frame` | `FrameAuto` | `frame` | Default frame type: `FrameAuto`, `FrameCode`, `FrameTerminal`, or `FrameNone` |

`WithDefaults` replaces the entire defaults struct. `WithLanguageDefaults` YAML keys accept comma-separated language names (e.g., `"bash, sh, zsh"`); each language receives the same defaults.

## Collapsible

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithCollapsible(CollapsibleConfig)` | `nil` (disabled) | `collapsible` | Enable threshold-based collapsible blocks |

Passing any `CollapsibleConfig` value enables the feature. Fields default to Go zero values (`0`, `false`, empty string) when using the Go API. The YAML config path initializes unset fields to `LineThreshold: 15`, `PreviewLines: 5`, `DefaultCollapsed: true`, `PreserveIndent: true`.

### `CollapsibleConfig` fields

| Field | Type | Go default | YAML default | YAML key | Description |
|-------|------|------------|--------------|----------|-------------|
| `LineThreshold` | `int` | `0` | `15` | `lineThreshold` | Minimum lines to trigger auto-collapse (>= 1) |
| `PreviewLines` | `int` | `0` | `5` | `previewLines` | Visible lines before expand (>= 1) |
| `DefaultCollapsed` | `bool` | `false` | `true` | `defaultCollapsed` | Start collapsed when threshold is met |
| `PreserveIndent` | `bool` | `false` | `true` | `preserveIndent` | Preserve minimum indentation in preview |
| `Style` | `CollapseStyle` | `CollapseGithub` | `CollapseGithub` | `style` | Visual style |
| `ExpandButtonText` | `string` | `""` (locale) | `""` (locale) | `expandButtonText` | Expand button label |
| `CollapseButtonText` | `string` | `""` (locale) | `""` (locale) | `collapseButtonText` | Collapse button label |
| `ExpandedAnnouncement` | `string` | `""` (locale) | `""` (locale) | `expandedAnnouncement` | Screen reader announcement on expand |
| `CollapsedAnnouncement` | `string` | `""` (locale) | `""` (locale) | `collapsedAnnouncement` | Screen reader announcement on collapse |

`CollapseStyle` values: `CollapseGithub`, `CollapseCollapsibleStart`, `CollapseCollapsibleEnd`, `CollapseCollapsibleAuto`. See [Collapsible Sections](/features/collapsible-sections/) for per-block meta string syntax and behavior details.

## Localization

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithLocale(string)` | `"en-US"` | `locale` | BCP 47 locale tag for built-in UI strings |
| `WithUIStrings(map[string]string)` | `nil` | `uiStrings` | Override individual UI strings by key |

`WithUIStrings` keys merge with locale-resolved strings. Override keys take precedence.

## Output panel

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithOutputPanel(bool)` | `false` | `outputPanel` | Enable output panel support (splits code at the separator) |
| `WithOutputCollapsed(bool)` | `false` | `outputDefaultCollapsed` | Start output panels collapsed by default |
| `WithOutputSeparator(string)` | `"---output---"` | `outputSeparator` | Custom separator string between code and output |

See [Output Panel](/features/output-panel/) for per-block meta string options and behavior.

## Post-render callbacks

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithPostRender(func(string, BlockInfo) string)` | `nil` | Go API only | Transform rendered HTML after each block; multiple callbacks chain in order |

Each callback receives the rendered HTML string and a `BlockInfo` struct with block metadata (language, title, frame type, raw code, line count, theme, meta string). Return the modified HTML. Multiple `WithPostRender` calls accumulate; callbacks run in registration order. See [Post-Render Callbacks](/plugins/post-render-callbacks/) for examples.

## Output and loading

| Option | Default | YAML key | Description |
|--------|---------|----------|-------------|
| `WithMinify(bool)` | `true` | `minify` | Collapse whitespace in generated CSS and JS |
| `WithWarningHandler(func(string))` | `nil` (`log.Printf`) | Go API only | Receive non-fatal warnings; pass a no-op to silence |
| `WithConfigDir(string)` | n/a | Go API only | Search a directory for `kazari.config.yaml/.yml/.json` and apply it |

`WithConfigDir` is a functional option like any other. Its position in the `kazari.New()` call determines its priority relative to other options. If no config file is found, it is a silent no-op. Parse errors are reported via `WarningHandler`. See [File-Based Config](/reference/file-based-config/) for the file format and composition rules.
