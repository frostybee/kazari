---
title: "CSS Variables"
description: "Complete reference for all --kz-* CSS custom properties."
tags: [css-variables]
sidebar:
  order: 5
---

All Kazari CSS custom properties use the `--kz-` prefix. Override them in the consumer's stylesheet to customize appearance without rebuilding the engine. See [CSS Custom Properties](/styling/css-custom-properties/) for cascade behavior and [Style Overrides](/styling/style-overrides/) for setting variables via the Go API.

## Frame and container

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-radius` | `0.5rem` | Border radius on frames and code group containers |
| `--kz-shadow` | `0 2px 8px rgba(0,0,0,0.15)` | Box shadow on frames |
| `--kz-border` | `1px solid transparent` | Border on frames |
| `--kz-transition` | `150ms ease` | Transition timing for code group tabs |

## Typography

### Code

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-font-family` | `'JetBrains Mono Variable', monospace` | Code font family |
| `--kz-font-size` | `0.875rem` | Code font size |
| `--kz-font-weight` | `500` | Code font weight |
| `--kz-line-height` | `1.6` | Code line height |

### UI elements

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-ui-font-family` | `system-ui, sans-serif` | Font family for toolbar, tabs, labels, and tooltips |
| `--kz-ui-font-size` | `0.9rem` | UI element font size |
| `--kz-ui-font-weight` | `400` | UI element font weight |
| `--kz-ui-line-height` | `1.65` | UI element line height |

## Code area

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-code-padding-block` | `1rem` | Vertical padding of the code element |
| `--kz-code-padding-inline` | `1.35rem` | Horizontal padding of each code line |
| `--kz-title-font-size` | `0.8rem` | Title text font size |
| `--kz-title-padding` | `0.5rem 1rem` | Title padding |
| `--kz-indent` | `0ch` | Per-line indent offset for wrapped-line indentation and collapsed-summary indent preservation (set inline per line) |

## Editor colors

These values are derived from the active syntax theme at engine construction time.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-editor-bg` | Theme-derived | Code area and frame background |
| `--kz-editor-fg` | Theme-derived | Default code text color |

## Toolbar

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-toolbar-padding` | `0.25rem 1rem` | Toolbar strip padding |
| `--kz-toolbar-bg` | light: `rgba(229,231,235,0.15)` / dark: `rgba(39,39,42,0.6)` | Toolbar background |
| `--kz-toolbar-border` | light: `rgba(209,213,219,0.5)` / dark: `rgba(63,63,70,0.4)` | Toolbar bottom border |
| `--kz-lang-font-size` | `0.8rem` | Language badge font size |
| `--kz-lang-font-weight` | `500` | Language badge font weight |
| `--kz-lang-fg` | light: `#4b5563` / dark: `#a1a1aa` | Language badge text color |
| `--kz-separator-color` | `rgba(161,161,170,0.3)` | Separator between language badge and title |

## Buttons

These variables apply to the copy, fullscreen, wrap, and theme toggle buttons.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-copy-radius` | `0.375rem` | Button border radius |
| `--kz-copy-fg` | light: `#4b5563` / dark: `#d4d4d8` | Button icon color |
| `--kz-copy-fg-hover` | light: `#111827` / dark: `#ffffff` | Button icon color on hover |
| `--kz-copy-bg-hover` | light: `rgba(156,163,175,0.2)` / dark: `rgba(63,63,70,0.8)` | Button background on hover |
| `--kz-copy-success-bg` | `rgba(34,197,94,0.9)` | Copy button background on success |
| `--kz-copy-success-fg` | `#ffffff` | Copy button icon color on success |
| `--kz-copy-success-border` | `rgba(34,197,94,0.8)` | Copy button border on success |

## Tooltip

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-tooltip-bg` | `rgba(30,30,30,0.92)` | Tooltip background |
| `--kz-tooltip-fg` | `#ffffff` | Tooltip text color |
| `--kz-tooltip-font-size` | `0.75rem` | Tooltip font size |
| `--kz-tooltip-padding` | `0.35rem 0.75rem` | Tooltip inner padding |
| `--kz-tooltip-radius` | `6px` | Tooltip border radius |
| `--kz-tooltip-offset` | `6px` | Gap between trigger element and tooltip |
| `--kz-tooltip-shadow` | `0 2px 6px rgba(0,0,0,0.25)` | Tooltip box shadow |
| `--kz-tooltip-arrow-size` | `5px` | Tooltip arrow size |

## Focus ring

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-focus-ring` | `rgb(59,130,246)` | Outline color for keyboard-focused interactive elements |

## Line numbers

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-ln-width` | `2ch` | Line number column width (widens dynamically for 3+ digit counts) |
| `--kz-ln-padding-inline` | `2ch` | Inline padding inside each line number cell |
| `--kz-ln-fg` | light: `#6e7781` / dark: `#6e7681` | Line number text color (overridden by theme `LineNumberFG` if present) |
| `--kz-ln-opacity` | `0.9` | Line number opacity |
| `--kz-ln-highlight-fg` | light: `#24292f` / dark: `#e6edf3` | Line number color on highlighted rows |
| `--kz-ln-highlight-opacity` | `0.8` | Line number opacity on highlighted rows |
| `--kz-gutter-border-width` | `1px` | Width of the border between gutter and code |
| `--kz-gutter-border-color` | light: `rgba(0,0,0,0.1)` / dark: `rgba(255,255,255,0.1)` | Gutter border color |

## Line markers

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-mark-bg` | `rgba(255,200,0,0.12)` | Highlight (mark) line background |
| `--kz-mark-border` | `rgba(255,200,0,0.5)` | Highlight line left accent bar color |
| `--kz-mark-border-width` | `3px` | Left accent bar width for all marker types |
| `--kz-mark-accent-margin` | `0rem` | Left accent bar margin |
| `--kz-ins-bg` | `rgba(46,160,67,0.12)` | Inserted line background |
| `--kz-ins-border` | `rgba(46,160,67,0.5)` | Inserted line left accent bar color |
| `--kz-ins-indicator` | `'+'` | Diff indicator character for inserted lines |
| `--kz-ins-indicator-color` | `rgba(46,160,67,0.8)` | Diff indicator color for inserted lines |
| `--kz-del-bg` | `rgba(248,81,73,0.12)` | Deleted line background |
| `--kz-del-border` | `rgba(248,81,73,0.5)` | Deleted line left accent bar color |
| `--kz-del-indicator` | `'-'` | Diff indicator character for deleted lines |
| `--kz-del-indicator-color` | `rgba(248,81,73,0.8)` | Diff indicator color for deleted lines |
| `--kz-diff-indicator-margin` | `0.3rem` | Left margin for diff indicator pseudo-elements |

## Labeled range markers

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-label-mark-bg` | `rgba(255,200,0,0.35)` | Labeled mark range line background |
| `--kz-label-ins-bg` | `rgba(46,160,67,0.35)` | Labeled ins range line background |
| `--kz-label-del-bg` | `rgba(248,81,73,0.35)` | Labeled del range line background |
| `--kz-label-fg` | `#ffffff` | Label badge text color |
| `--kz-label-padding` | `0.1rem 0.3rem` | Label badge padding |
| `--kz-label-font-size` | `0.75rem` | Label badge font size |
| `--kz-label-radius` | `0.2rem` | Label badge border radius |

## Inline markers

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-inline-mark-bg` | `rgba(255,200,0,0.2)` | Inline mark background |
| `--kz-inline-mark-border` | `rgba(255,200,0,0.5)` | Inline mark border color |
| `--kz-inline-mark-radius` | `0.2rem` | Inline mark border radius |
| `--kz-inline-mark-padding` | `0.15rem` | Inline mark padding |
| `--kz-inline-mark-border-width` | `1.5px` | Inline mark border width |
| `--kz-inline-ins-bg` | `rgba(46,160,67,0.2)` | Inline inserted text background |
| `--kz-inline-ins-border` | `rgba(46,160,67,0.5)` | Inline inserted text border color |
| `--kz-inline-del-bg` | `rgba(248,81,73,0.2)` | Inline deleted text background |
| `--kz-inline-del-border` | `rgba(248,81,73,0.5)` | Inline deleted text border color |

## Focus lines

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-focus-dimmed-opacity` | `0.35` | Opacity for non-focused lines when focus is active |

## Terminal frame

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-terminal-bg` | `var(--kz-editor-bg)` | Terminal code area background |
| `--kz-terminal-titlebar-bg` | `var(--kz-toolbar-bg)` | Terminal title bar background |
| `--kz-terminal-header-padding` | `0.5rem 1rem` | Terminal title bar padding |

### Colored dots

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-terminal-dot-red` | `#ff5f57` | Red (close) dot color |
| `--kz-terminal-dot-yellow` | `#febc2e` | Yellow (minimize) dot color |
| `--kz-terminal-dot-green` | `#28c840` | Green (maximize) dot color |

### Minimal dots

Emitted when `WithTerminalDotStyle(DotsMinimal)` is configured.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-terminal-dots-fg` | light: `#24292f` / dark: `#e6edf3` | Monochrome dot color |
| `--kz-terminal-dots-opacity` | `0.15` | Dot opacity |
| `--kz-terminal-icon` | Inline SVG data URL | Three-circle mask image for the dot shape |

## Collapsible sections

Always emitted; range-based collapse can appear on any engine via the meta string.

### Button

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-collapse-btn-bg` | light: `rgba(0,0,0,0.04)` / dark: `rgba(255,255,255,0.1)` | Expand/collapse button background |
| `--kz-collapse-btn-fg` | light: `#4b5563` / dark: `#d4d4d8` | Button text color |
| `--kz-collapse-btn-hover-bg` | light: `rgba(0,0,0,0.08)` / dark: `rgba(255,255,255,0.18)` | Button background on hover |
| `--kz-collapse-btn-border` | light: `rgba(0,0,0,0.15)` / dark: `rgba(255,255,255,0.2)` | Button border color |
| `--kz-collapse-btn-border-hover` | light: `rgba(0,0,0,0.3)` / dark: `rgba(255,255,255,0.35)` | Button border on hover |

### Gradient

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-collapse-gradient-start` | `transparent` | Top color of the fade gradient over collapsed content |
| `--kz-collapse-gradient-end` | `var(--kz-editor-bg)` | Bottom color of the fade gradient |
| `--kz-collapse-transition` | `300ms ease` | Fade gradient transition duration |

### Closed sections

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-collapse-closed-bg` | `rgb(84 174 255 / 20%)` | Collapsed summary row background (overridden by theme fold color) |
| `--kz-collapse-closed-border` | `rgb(84 174 255 / 50%)` | Collapsed section border color (overridden by theme fold color) |
| `--kz-collapse-closed-border-width` | `0` | Collapsed section border width |
| `--kz-collapse-closed-fg` | `currentColor` | Collapsed summary text color |
| `--kz-collapse-closed-padding` | `4px` | Collapsed summary row padding |
| `--kz-collapse-closed-font-family` | `inherit` | Collapsed summary font family |
| `--kz-collapse-closed-font-size` | `inherit` | Collapsed summary font size |
| `--kz-collapse-closed-line-height` | `inherit` | Collapsed summary line height |

### Open sections

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-collapse-open-bg` | `transparent` | Open section background (github style) |
| `--kz-collapse-open-bg-collapsible` | `rgb(84 174 255 / 10%)` | Open section background (collapsible-start/end styles) |
| `--kz-collapse-open-border` | `transparent` | Open section border color |
| `--kz-collapse-open-border-width` | `1px` | Open section border width |

### Icons

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-collapse-expand-icon` | Inline SVG data URL | Mask image for the expand icon on closed sections |
| `--kz-collapse-collapse-icon` | Inline SVG data URL | Mask image for the collapse icon on open sections |

## Code groups

Emitted when `engine.EnableCodeGroups()` is called.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-group-border-width` | `1px` | Code group container border width |
| `--kz-group-border` | light: `rgba(0,0,0,0.1)` / dark: `rgba(255,255,255,0.1)` | Container and tab bar border color |
| `--kz-group-radius` | `var(--kz-radius)` | Container border radius |
| `--kz-group-tab-bg` | `transparent` | Inactive tab background |
| `--kz-group-tab-fg` | `inherit` | Inactive tab text color |
| `--kz-group-tab-active-bg` | light: `rgba(255,255,255,0.8)` / dark: `rgba(255,255,255,0.1)` | Active tab background |
| `--kz-group-tab-active-fg` | light: `#24292f` / dark: `#e6edf3` | Active tab text color |
| `--kz-group-tab-active-border` | `#007acc` | Active tab bottom border (underline) color |
| `--kz-group-tab-padding` | `0.5rem 1rem` | Tab inner padding |

## Output panel

Emitted when `WithOutputPanel(true)` is configured.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-output-bg` | `var(--kz-editor-bg)` | Output panel background |
| `--kz-output-fg` | `var(--kz-editor-fg)` | Output panel text color |
| `--kz-output-opacity` | `0.7` | Output panel content opacity |
| `--kz-output-header-bg` | `var(--kz-toolbar-bg)` | Output panel header background |
| `--kz-output-border` | `rgba(128,128,128,0.2)` | Output panel top border |
| `--kz-output-label-color` | `currentColor` | Output panel toggle label color |

## Scrollbar

Emitted when `WithThemedScrollbars(true)` is configured (default).

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-scrollbar-width` | `5px` | Scrollbar track width |
| `--kz-scrollbar-height` | `5px` | Scrollbar track height |
| `--kz-scrollbar-track` | `transparent` | Scrollbar track background |
| `--kz-scrollbar-thumb` | light: `rgba(0,0,0,0.2)` / dark: `rgba(255,255,255,0.15)` | Scrollbar thumb color |
| `--kz-scrollbar-thumb-hover` | light: `rgba(0,0,0,0.35)` / dark: `rgba(255,255,255,0.3)` | Scrollbar thumb color on hover |

## Text selection

Emitted when `WithThemedSelectionColors(true)` is configured.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-selection-bg` | `rgba(0,122,204,0.3)` (overridden by theme `SelectionBG` if present) | Text selection background |
| `--kz-selection-fg` | `inherit` | Text selection foreground |

## Fullscreen

Emitted when `WithFullscreenButton(true)` is configured (default).

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-fs-font-scale` | `1` | Font scale multiplier in fullscreen mode; adjusted by font size buttons |

## Links

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-link-underline` | `currentColor` | Link underline decoration color |
| `--kz-link-underline-width` | `1px` | Link underline thickness |
| `--kz-link-underline-offset` | `0.15em` | Offset between text baseline and underline |

## File icons

Emitted when `WithFileIcons(true)` is configured (default).

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-file-icon-size` | `1rem` | File icon width and height |
| `--kz-file-icon-margin` | `0` | File icon margin, on top of the toolbar's own gap |
| `--kz-file-icon-opacity` | `0.8` | File icon opacity |

## Language icons

Emitted when `WithLanguageIconMode` is not `LangIconNone`.

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-lang-icon-size` | `1.25rem` | Language icon width and height |
| `--kz-lang-icon-margin` | `0` | Language icon margin |
| `--kz-lang-icon-opacity` | `0.8` | Language icon opacity |

## ANSI colors

Standard 16-color ANSI palette for `lang="ansi"` blocks. See [ANSI Rendering](/features/ansi-rendering/).

| Variable | Default | Description |
|----------|---------|-------------|
| `--kz-ansi-black` | `#000000` | Color 0 (black) |
| `--kz-ansi-red` | `#cc0000` | Color 1 (red) |
| `--kz-ansi-green` | `#4e9a06` | Color 2 (green) |
| `--kz-ansi-yellow` | `#c4a000` | Color 3 (yellow) |
| `--kz-ansi-blue` | `#3465a4` | Color 4 (blue) |
| `--kz-ansi-magenta` | `#75507b` | Color 5 (magenta) |
| `--kz-ansi-cyan` | `#06989a` | Color 6 (cyan) |
| `--kz-ansi-white` | `#d3d7cf` | Color 7 (white) |
| `--kz-ansi-bright-black` | `#555753` | Color 8 (bright black) |
| `--kz-ansi-bright-red` | `#ef2929` | Color 9 (bright red) |
| `--kz-ansi-bright-green` | `#8ae234` | Color 10 (bright green) |
| `--kz-ansi-bright-yellow` | `#fce94f` | Color 11 (bright yellow) |
| `--kz-ansi-bright-blue` | `#729fcf` | Color 12 (bright blue) |
| `--kz-ansi-bright-magenta` | `#ad7fa8` | Color 13 (bright magenta) |
| `--kz-ansi-bright-cyan` | `#34e2e2` | Color 14 (bright cyan) |
| `--kz-ansi-bright-white` | `#eeeeec` | Color 15 (bright white) |
