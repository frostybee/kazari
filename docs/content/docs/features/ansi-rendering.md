---
title: "ANSI Rendering"
description: "Render terminal output with ANSI escape sequences as colored, styled code blocks."
tags: [terminal, highlighter]
sidebar:
  order: 12
---

ANSI rendering converts terminal output containing SGR escape sequences into colored, styled code blocks. Set the language to `ansi` and Kazari routes the code through a built-in parser instead of the external highlighter.

## Basic usage

Set the language to `ansi` in the fence meta string:

````
```ansi title="test output"
...terminal output with escape sequences...
```
````

:::caution
ANSI escape sequences contain raw ESC bytes (0x1B) that are invisible in most text editors. When using the Go API directly, pass strings with `\x1b[...m` escape literals to `engine.RenderWithMeta()`. When using the Goldmark extension, raw ESC bytes can be embedded in Markdown fence blocks and are processed at build time.
:::

Pass ANSI-escaped strings via the Go API:

```go
output := "\x1b[32m✔\x1b[0m Tests passed  \x1b[33m3 skipped\x1b[0m\n" +
    "\x1b[31m✘\x1b[0m \x1b[1mBuildError:\x1b[0m missing return statement"

html, err := engine.RenderWithMeta(output, `ansi title="test output"`)
```

→ The first line renders a green check mark followed by normal text and a yellow "3 skipped". The second line renders a red cross, then "BuildError:" in bold, then normal text.

All standard meta string options work on ANSI blocks: `title=`, `frame=`, `showLineNumbers`, `{N-M}` line markers, `ins={}`, `del={}`, `focus={}`, and inline markers. Because `ansi` is in the terminal language list, blocks auto-detect as terminal frame when `WithFrameDetection` is active.

```go
html, err := engine.RenderWithMeta(output, `ansi title="build log" showLineNumbers {3}`)
```

→ A terminal-framed block titled "build log" with numbered lines and line 3 highlighted.

ANSI rendering bypasses the external highlighter entirely. It works with no highlighter configured on the engine.

## Visual showcase

The block below demonstrates every color mode and font style that Kazari's ANSI parser supports.

```ansi
[1;4mStandard ANSI colors:[0m
- Foreground:  [30mBlack[0m [31mRed[0m [32mGreen[0m [33mYellow[0m [34mBlue[0m [35mMagenta[0m [36mCyan[0m [37mWhite[0m
- Background:  [41m Red [0m [42m Green [0m [43m Yellow [0m [44m Blue [0m [45m Magenta [0m [46m Cyan [0m [47m White [0m

[1;4mBright colors:[0m
- Foreground:  [90mBlack[0m [91mRed[0m [92mGreen[0m [93mYellow[0m [94mBlue[0m [95mMagenta[0m [96mCyan[0m [97mWhite[0m
- Background:  [101m Red [0m [102m Green [0m [103m Yellow [0m [104m Blue [0m [105m Magenta [0m [106m Cyan [0m [107m White [0m

[1;4m8-bit colors (sample from 256-color palette):[0m
- Foreground:  [38;5;196m196[0m [38;5;208m208[0m [38;5;220m220[0m [38;5;46m46[0m [38;5;33m33[0m [38;5;129m129[0m
- Background:  [48;5;196m 196 [0m [48;5;208m 208 [0m [48;5;220m 220 [0m [48;5;46m 46 [0m [48;5;33m 33 [0m [48;5;129m 129 [0m
- Grayscale:   [38;5;232m██[0m [38;5;236m██[0m [38;5;240m██[0m [38;5;244m██[0m [38;5;248m██[0m [38;5;252m██[0m [38;5;255m██[0m

[1;4m24-bit true colors (full RGB):[0m
- Foreground:  [38;2;34;139;34mForestGreen[0m [38;2;255;165;0mOrange[0m [38;2;102;51;153mRebeccaPurple[0m [38;2;0;191;255mDeepSkyBlue[0m
- Background:  [48;2;34;139;34m ForestGreen [0m [48;2;255;165;0m Orange [0m [48;2;102;51;153m RebeccaPurple [0m [48;2;0;191;255m DeepSkyBlue [0m

[1;4mFont styles:[0m
- Default
- [1mBold[0m
- [3mItalic[0m
- [4mUnderline[0m
- [9mStrikethrough[0m
- [1;3mBold + Italic[0m
- [1;4;32mBold + Underline + Green[0m
```

## Supported sequences

| SGR code | Effect |
|---|---|
| `0` (or empty params) | Reset all (clear color, background, and font style) |
| `1` | Bold |
| `3` | Italic |
| `4` | Underline |
| `9` | Strikethrough |
| `22` | Remove bold |
| `23` | Remove italic |
| `24` | Remove underline |
| `29` | Remove strikethrough |
| `30`-`37` | Standard foreground color (8 colors) |
| `38;5;N` | 256-color foreground (N = 0-255) |
| `38;2;R;G;B` | 24-bit true-color foreground |
| `39` | Default foreground (clears color) |
| `40`-`47` | Standard background color (8 colors) |
| `48;5;N` | 256-color background |
| `48;2;R;G;B` | 24-bit true-color background |
| `49` | Default background (clears background) |
| `90`-`97` | Bright foreground (8 colors) |
| `100`-`107` | Bright background (8 colors) |

SGR state persists across lines within a single block, matching real terminal behavior. A color set on one line remains active on the next line unless explicitly reset. Malformed sequences (those containing non-digit, non-semicolon characters before `m`) pass through as literal text.

The following SGR codes are not supported: dim (`2`), blink (`5`/`6`), reverse video (`7`), and hidden/conceal (`8`). Non-SGR escape sequences (CSI cursor movement, OSC sequences) are also not handled.

## Color modes

Standard 16 colors and bright colors (SGR 30-37, 40-47, 90-97, 100-107) resolve to CSS variable references (`var(--kz-ansi-red)`, etc.) rather than hardcoded hex values. This allows retheming by overriding the `--kz-ansi-*` variables without touching the HTML.

256-color mode (SGR `38;5;N` / `48;5;N`) resolves differently by index range:

- **Index 0-15**: hardcoded hex values from the standard palette (not CSS variable references)
- **Index 16-231**: computed from a 6x6x6 RGB color cube
- **Index 232-255**: grayscale ramp from `#080808` to `#eeeeee`

24-bit true-color (SGR `38;2;R;G;B` / `48;2;R;G;B`) produces direct `#rrggbb` hex values.

Standard 16 colors are themeable via the CSS variables. 256-color and 24-bit true-color values use literal hex and are not themeable through CSS variables.

## CSS variables

The `--kz-ansi-*` variables define the standard 16-color palette. They are always emitted in `engine.CSS()` output regardless of configuration.

| Variable | Default |
|---|---|
| `--kz-ansi-black` | `#000000` |
| `--kz-ansi-red` | `#cc0000` |
| `--kz-ansi-green` | `#4e9a06` |
| `--kz-ansi-yellow` | `#c4a000` |
| `--kz-ansi-blue` | `#3465a4` |
| `--kz-ansi-magenta` | `#75507b` |
| `--kz-ansi-cyan` | `#06989a` |
| `--kz-ansi-white` | `#d3d7cf` |
| `--kz-ansi-bright-black` | `#555753` |
| `--kz-ansi-bright-red` | `#ef2929` |
| `--kz-ansi-bright-green` | `#8ae234` |
| `--kz-ansi-bright-yellow` | `#fce94f` |
| `--kz-ansi-bright-blue` | `#729fcf` |
| `--kz-ansi-bright-magenta` | `#ad7fa8` |
| `--kz-ansi-bright-cyan` | `#34e2e2` |
| `--kz-ansi-bright-white` | `#eeeeec` |

The defaults match the GNOME Terminal / VGA palette. Override any variable by redeclaring it on `:root` or a more specific selector in your stylesheet:

```css
:root {
    --kz-ansi-red: #e06c75;
    --kz-ansi-green: #98c379;
}
```

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- ANSI rendering bypasses the external highlighter entirely. It works with no highlighter configured on the engine.
- The copy button copies raw text including escape codes, not the display text.
- [Terminal comment stripping](/features/terminal-comment-stripping/) applies to ANSI blocks by default (terminal frame). Lines starting with `#` are stripped from the copy payload only, not from the display.
- ANSI escape sequences contain raw ESC bytes (0x1B) that are invisible in most text editors. The Go API uses `\x1b[...m` escape literals; the Goldmark extension processes raw ESC bytes embedded in Markdown fence blocks at build time.
- 256-color index 0-15 uses hardcoded hex values, not the `--kz-ansi-*` CSS variables. Only the standard 16 SGR colors (30-37, 40-47, 90-97, 100-107) use CSS variables.
- The `--kz-ansi-*` variables are theme-independent. They are emitted once in the static CSS block and do not change between light and dark themes.
