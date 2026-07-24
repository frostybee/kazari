---
title: "Dual-Theme Rendering"
description: "How Kazari bakes both light and dark token colors into HTML and switches themes with pure CSS."
sidebar:
  order: 1
---

Kazari embeds both light and dark token colors into the rendered HTML at build time. Theme switching happens entirely through CSS custom properties, with no JavaScript involvement and no additional network requests. The result is that toggling between light and dark mode produces no flash or flicker, even on the initial page load. This strategy mirrors the dual-theme approach used by Shiki.

For configuration options (theme names, dark mode strategies, per-block overrides), see [Themes & Dark Mode](/styling/themes-and-dark-mode/).

## The problem

Most approaches to dual-theme code blocks cause a visible color change when the page loads or the user switches themes. Stylesheet swapping requires fetching a new CSS file, leaving a brief period where no theme is applied. JavaScript re-rendering re-tokenizes or re-colorizes after the DOM is available, producing a flash from an uncolored state. Single-theme output carries only one color set and cannot adapt to the opposite mode without re-rendering.

Kazari avoids all three problems by including both color sets in every token span during the initial render. The browser always has both colors available and only needs to read a different CSS variable to switch between them.

## MergedToken

After tokenization, Kazari merges the light and dark token streams into a slice of `MergedToken` structs per line (`internal/config/config.go`):

```go
type MergedToken struct {
    Content    string
    LightColor string  // e.g. "#cf222e"
    DarkColor  string  // e.g. "#ff7b72"
    LightBG    string
    DarkBG     string
    FontStyle  int     // bitmask: italic=1, bold=2, underline=4, strikethrough=8
}
```

`FontStyle` is always taken from the light theme. Both streams must span identical content: same lines, same character boundaries.

## Token alignment

The light and dark tokenizers do not always produce tokens with identical boundaries. Two separate `Tokenize` calls on the same source may split at different positions depending on the grammar rules activated by each theme. Kazari handles this with two paths:

**Fast path**: when both streams have identical token counts per line, tokens are paired by index. This is the common case when the engine uses a `DualThemeTokenizer` (see below), which produces both streams in a single grammar pass and guarantees identical splits.

**Slow path** (`alignTokens` in `token_merge.go`): when token counts differ, Kazari walks both streams simultaneously by character position. At each step it advances the stream with the shorter remaining token, splitting at that boundary. A merged token is produced for each unique boundary segment. Residual tokens in the light stream beyond the end of the dark stream receive no `DarkColor`.

## The `--sl` and `--sd` inline variables

Each span in the rendered HTML carries both theme colors as inline CSS custom properties:

```html
<span style="--sl:#cf222e;--sd:#ff7b72">func</span>
```

The naming conventions:

| Variable | Meaning |
|----------|---------|
| `--sl` | Light foreground color |
| `--sd` | Dark foreground color |
| `--slbg` | Light background color (only when the token has an explicit background) |
| `--sdbg` | Dark background color (only when the token has an explicit background) |
| `--sfs` | Font style (`italic`) |
| `--sfw` | Font weight (`bold`) |
| `--std` | Text decoration (`underline`, `line-through`, or both) |

When only a single theme is configured (the dark theme is an empty string), only `--sl` is emitted. `--sd` is omitted entirely for single-theme blocks. The `hasDualTheme()` function (`internal/render/render_lines.go`) scans all tokens for at least one non-empty `DarkColor` before gating `--sd` emission.

A span with a background color looks like:

```html
<span style="--sl:#ffffff;--sd:#000000;--slbg:#0550ae;--sdbg:#1f6feb">type</span>
```

A plain text span with no special coloring emits nothing, letting it inherit the editor foreground:

```html
<span>   </span>
```

## The CSS switching rule

Two CSS rules control which color set is visible. The first rule is unconditional and resolves the light theme variables:

```css
.kazari-block .kz-line span[style^="--"] {
    color: var(--sl, inherit);
    background-color: var(--slbg, transparent);
    font-style: var(--sfs, inherit);
    font-weight: var(--sfw, inherit);
    text-decoration: var(--std, inherit);
}
```

The second rule is scoped to the dark mode context. With the default `SelectorMode(".dark")` strategy it looks like:

```css
.dark .kazari-block .kz-line span[style^="--"] {
    color: var(--sd, inherit);
    background-color: var(--sdbg, transparent);
}
```

The attribute selector `span[style^="--"]` targets spans whose `style` attribute starts with `--`. This is a starts-with selector, not a substring selector (`[style*="--"]`). The distinction matters: a starts-with selector matches only spans that carry Kazari token variables, not arbitrary spans that happen to include custom properties somewhere in their style string.

When `--sl` is undefined on a span (for example, plain whitespace with no token color), `var(--sl, inherit)` falls back to `inherit`, which is the correct behavior since those spans should inherit the editor foreground.

## DualThemeTokenizer

The `DualThemeTokenizer` interface allows a highlighter to produce both color streams in a single tokenization pass:

```go
type DualThemeTokenizer interface {
    TokenizeDual(code, lang, lightTheme, darkTheme string) (light, dark [][]Token, err error)
}
```

The Nuri adapter implements this interface. One TextMate grammar scan processes both theme's color rules simultaneously, producing two token streams with guaranteed identical boundaries. This is why the fast path (index-based merging) applies for most Nuri-based rendering.

The Chroma adapter does not implement `DualThemeTokenizer`. When dual-theme rendering is active, the engine type-asserts the highlighter and, on failure, falls back to calling `Tokenize` twice: once for the light theme, once for the dark theme. The merged result is semantically identical; the difference is CPU cost, not output quality. For large documents or high-throughput pipelines, the Nuri adapter halves tokenization work.

## Per-block theme override

The `theme=` meta token triggers a separate color extraction for the specified theme. Override colors are embedded as inline CSS variables directly on the block wrapper element, not on individual token spans:

```html
<figure class="kazari-block kz-themed" style="--kz-ovl-editor-bg:#0d1117;--kz-ovl-editor-fg:#e6edf3;...">
```

The variable naming uses two prefixes: `--kz-ovl-*` for light (override light) values and `--kz-ovd-*` for dark (override dark) values.

CSS rules scoped to `.kz-themed` remap these to the `--kz-*` names that the rest of the block's CSS consumes. The dark fallback chain ensures single-theme overrides work correctly in dark mode:

```css
--kz-editor-bg: var(--kz-ovd-editor-bg, var(--kz-ovl-editor-bg));
```

If `--kz-ovd-editor-bg` is not set (because the `theme=` value was a single name rather than `"light,dark"`), the fallback resolves to `--kz-ovl-editor-bg`, which is always set.

Override results are cached per unique `theme=` string in a `sync.RWMutex`-protected map, so repeated renders of the same override incur no repeated color extraction.

## The no-flash invariant

Three properties together guarantee that code blocks never flash or transition during theme changes.

**Colors in HTML at build time.** Both `--sl` and `--sd` are present in the HTML before any JavaScript runs. The browser has both values on the initial parse. No fetch, no delay, no blank state.

**No CSS transitions.** `base.css` (`internal/css/static/base.css`) sets `transition: none !important` on `.kazari-block` and all its descendants. A CSS transition on a color property would produce a visible animation during theme switches. This rule prevents any such transition, even if a parent element's transition applies to `color` or `background-color`.

**Unconditional `var()` fallbacks.** The `span[style^="--"]` rule always reads `--sl` (light) or `--sd` (dark) through `var()`. There are no conditional rules that match only when a specific string appears in the style attribute. Conditional selectors introduce a window where a span matches neither rule, causing a momentary color drop. The unconditional structure eliminates that window entirely.

These three properties are sufficient. No View Transitions API, no `color-scheme` meta hint, and no JavaScript workaround is needed.
