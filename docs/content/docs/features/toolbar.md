---
title: "Toolbar"
description: "Copy, fullscreen, wrap toggle, and theme toggle buttons."
tags: [client-side-js, accessibility]
sidebar:
  order: 2
---

The toolbar sits in the frame's title bar and holds action buttons for interacting with the code block. Five components are available: copy button, fullscreen button, wrap toggle, theme toggle, and language badge.

All toolbar buttons are engine-level options set via `kazari.New()`. There are no meta string toggles for individual buttons.

## Copy button

Copies the raw code to the clipboard. Enabled by default.

```go
kazari.WithCopyButton(false)
```

```yaml
copyButton: false
```

After a successful copy, a checkmark icon and "Copied!" tooltip appear for 2 seconds, then reset. Screen readers receive an `aria-live` announcement.

For terminal frames, `WithTerminalCommentStripping(true)` (the default) excludes lines starting with `#` from the copied text, so the reader can paste commands directly. See [Terminal Comment Stripping](/features/terminal-comment-stripping/) for details.

On frameless blocks (`frame=none`), the copy button floats in the top-right corner, visible on hover and keyboard focus. On touch devices, it is always visible. It is hidden in print.

## Fullscreen button

Enters native browser fullscreen on the code block wrapper. Enabled by default.

```go
kazari.WithFullscreenButton(false)
```

```yaml
fullscreenButton: false
```

In fullscreen mode, two font size controls appear in the toolbar:

- **Increase**: +0.1 per click, up to a maximum of 5.0
- **Decrease**: -0.1 per click, down to a minimum of 0.6
- **Double-click the decrease button** to reset the font scale to 1.0

The font scale persists in `localStorage` and is restored on subsequent page loads. The font controls are always tied to the fullscreen button; there is no separate option to enable them independently.

## Wrap toggle button

Toggles word wrap on and off at runtime. Enabled by default.

```go
kazari.WithWrapButton(false)
```

```yaml
wrapButton: false
```

The button updates `aria-pressed` to reflect the current state. The initial state matches the `wrap` meta option or engine default. Toggling does not persist across page loads. Wrap layout options (`preserveIndent`, `hangingIndent`) are documented on the [Word Wrap](/features/word-wrap/) page.

## Theme toggle button

Toggles the block between light and dark themes independently of the page-level theme. Disabled by default.

```go
kazari.WithThemeToggle(true)
```

```yaml
themeToggleButton: true
```

The theme toggle button requires a dark theme configured via `WithThemes(light, dark)`. It is silently suppressed when no dark theme is set.

The toggle state persists per block in `localStorage`, keyed by an 8-character hash of the code content. It works with all three dark mode strategies (`SelectorMode`, `MediaQueryMode`, `BothMode`).

A per-block theme override can also be set directly in the meta string with `theme=`:

````
```go title="dracula.go" theme="dracula"
func main() {
    fmt.Println("different theme")
}
```
````

```go title="dracula.go" theme="dracula"
func main() {
    fmt.Println("different theme")
}
```

→ The code block renders with the Dracula theme colors regardless of the page-level theme.

Full theme configuration is documented on the [Themes & Dark Mode](/styling/themes-and-dark-mode/) page.

## Language badge

Shows the language name in the toolbar area. Enabled by default.

```go
kazari.WithLanguageBadge(false)
```

```yaml
languageBadge: false
```

Display names follow proper casing for known languages: `javascript` renders as "JavaScript", `typescript` as "TypeScript", `css` as "CSS", `html` as "HTML", `json` as "JSON", `yaml` as "YAML", `sql` as "SQL", `graphql` as "GraphQL". Other languages are title-cased.

Language icons (shown alongside or instead of the text label) are configured with `WithLanguageIconMode()`. See the [Icons](/features/icons/) page for details.

## Collapse chevron

When collapsible sections are configured via `WithCollapsible`, a chevron button appears in the toolbar for blocks that exceed the line threshold. It syncs with the bottom expand/collapse bar.

The chevron is only rendered in editor frames. Terminal frames show only the bottom bar. See [Collapsible Sections](/features/collapsible-sections/) for configuration details.

## Terminal frames

Terminal frames share the same action buttons (copy, wrap, theme toggle, fullscreen) but display them differently. The buttons appear in a `kz-terminal-actions` container on the right side of the header. Terminal frames do not show the language badge, file icons, or the collapse chevron.

## Configuration

| Option | Default | Description |
|---|---|---|
| `WithCopyButton(bool)` | `true` | Show or hide the copy button |
| `WithFullscreenButton(bool)` | `true` | Show or hide the fullscreen button (includes font controls) |
| `WithWrapButton(bool)` | `true` | Show or hide the wrap toggle button |
| `WithThemeToggle(bool)` | `false` | Show or hide the per-block theme toggle |
| `WithLanguageBadge(bool)` | `true` | Show or hide the language label |
| `WithTerminalCommentStripping(bool)` | `true` | Strip `#` lines from copied text in terminal frames |

Config file equivalents:

```yaml
copyButton: true
fullscreenButton: true
wrapButton: true
themeToggleButton: false
languageBadge: true
terminalCommentStripping: true
```

## CSS variables

| Variable | Description |
|---|---|
| `--kz-toolbar-bg` | Toolbar background color |
| `--kz-toolbar-border` | Toolbar bottom border |
| `--kz-toolbar-padding` | Toolbar padding |
| `--kz-copy-fg` | Copy button icon color |
| `--kz-copy-fg-hover` | Copy button icon color on hover |
| `--kz-copy-bg-hover` | Copy button background on hover |
| `--kz-copy-success-bg` | Copy button background on success |
| `--kz-copy-success-fg` | Copy button icon color on success |
| `--kz-copy-success-border` | Copy button border on success |
| `--kz-copy-radius` | Copy button border radius |
| `--kz-tooltip-bg` | Tooltip background |
| `--kz-tooltip-fg` | Tooltip text color |
| `--kz-tooltip-radius` | Tooltip border radius |
| `--kz-lang-fg` | Language badge text color |
| `--kz-lang-font-size` | Language badge font size |
| `--kz-separator-color` | Divider between language badge and title |
| `--kz-fs-font-scale` | Font scale factor in fullscreen mode |

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- Font controls are hidden outside of fullscreen mode. They appear only after entering fullscreen.
- Tooltips are hidden on touch devices (`hover: none`) and in print media.
- All toolbar buttons include ARIA labels and keyboard focus support.
