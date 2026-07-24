---
title: "Render Pipeline"
description: "The complete call path from Render/RenderWithMeta through meta parsing, config cascade, preprocessing, language routing, tokenization, and HTML output."
sidebar:
  order: 2
---

Every code block Kazari renders follows the same pipeline: parse input, resolve configuration, preprocess the source, route by language, tokenize, and assemble HTML. This page traces that pipeline from the public API entry points through to the final output string.

For the CSS-based theme switching mechanism that displays the rendered tokens, see [Dual-Theme Rendering](/architecture/dual-theme-rendering/).

## Entry points

Both public render methods in `kazari.go` converge on a single internal function, `renderResolved()`.

**`engine.Render(code, opts)`** accepts a structured `Options` value. It converts the Go struct fields into internal block options, resolves the language (applying any aliases from `WithLanguageAliases`), runs the config cascade, and passes the result to `renderResolved()`. Typed marker and collapse fields are copied directly with no string parsing.

**`engine.RenderWithMeta(code, metaStr)`** accepts a raw fence info string. It calls the meta string parser (`internal/meta/meta.go`), which tokenizes the string and extracts the language, block options, line markers, inline markers, focus ranges, and collapse directives. After parsing, the same language resolution and config cascade run. The parsed result then flows into `renderResolved()`.

Both paths produce identical output for equivalent inputs. `RenderWithMeta` is the path used by the [Goldmark extension](/integrations/goldmark-extension/).

## Config cascade

Before rendering begins, block-level configuration is resolved through a three-layer cascade in `Config.Resolve()` (`internal/config/config.go`). The first layer is engine defaults set by `WithDefaults`: frame type, line numbers, wrap, preserve indent, and hanging indent, with `StartLineNumber` defaulting to 1. The second layer is language defaults from `WithLanguageDefaults`, looked up by language name (keys support comma-separated names like `"bash,sh,zsh"`, matched case-insensitively); the first matching key replaces all five default properties. The third and highest-priority layer is per-block overrides from `Options` or the meta string, where each non-nil pointer field overrides its counterpart and a nil pointer means "inherit from the layer above."

The cascade produces a fully resolved block with all layout properties settled. Markers, focus lines, links, and collapse data are assigned separately after the cascade completes.

## Preprocessing

Four steps run in a fixed order on the source code. Each step may mutate the code string and the resolved block.

**1. File name extraction** (`internal/frame/frame.go`). If `WithFileNameExtraction` is enabled and the block has no title, the first line is inspected for a filename comment pattern (e.g., `// main.go` or `# config.toml`). If found, the filename becomes the title and that line is removed from the code.

**2. Link extraction** (`internal/link/link.go`). If `WithLinks` is enabled, `@[text](url)` annotations are scanned and removed from the code. The extracted link positions are stored for re-injection during token rendering.

**3. Frame detection.** If the frame is set to `FrameAuto` and `WithFrameDetection` is enabled, the language and code content determine whether the block gets an editor frame or a terminal frame. Shell languages (`bash`, `sh`, `zsh`, `powershell`, etc.) receive a terminal frame. When detection is disabled, `FrameAuto` falls back to `FrameCode`.

**4. Raw code snapshot.** The preprocessed code is saved as the copy-button source text. If the resolved frame is `FrameTerminal` and `WithTerminalCommentStripping` is enabled, shell comments are stripped from this snapshot (but preserved in the highlighted output).

## Collapse resolution

Collapse resolution runs on the preprocessed code so that line counts are correct after any filename-comment removal.

**Threshold mode.** If the engine has a `CollapsibleConfig` and the line count exceeds the threshold (and the block is not explicitly disabled with `nocollapse`), threshold collapse activates. The resolver computes which lines remain visible in the collapsed state, anchoring preview segments around marked or focused lines so that highlighted content is not hidden behind the fold. A gradient overlay and expand button are added below the visible preview.

**Range mode.** Explicit `collapse={N-M}` ranges are validated and clipped to the line count. Each range becomes a collapsible section in the output. The collapse style (`github`, `collapsible-start`, `collapsible-end`, `collapsible-auto`) determines the HTML structure and interaction model.

The resolver also computes the minimum indentation per range (for indent preservation in summary rows) and, for `collapsible-auto`, decides between start and end placement based on whether the range reaches the final line.

## Per-block theme override

If the meta string includes a `theme=` token and a highlighter is present, the override is resolved before tokenization. A single theme name applies to both modes; a comma-separated pair (e.g., `"dracula,nord"`) sets light and dark independently.

Override results are cached behind a read-write mutex (keyed by the override string). On cache miss, the engine extracts theme colors, runs the theme adjustments and customizer pipeline, and generates inline CSS variable declarations (`--kz-ovl-*` / `--kz-ovd-*`) scoped to the block. If `MinContrast` is active, marker background colors are recomputed against the override theme's editor background.

## Language routing

After preprocessing and collapse resolution, the pipeline branches by language. Four routes are checked in order; the first match wins.

**1. Mermaid pass-through.** If `WithMermaidPassThrough` is enabled and the language is `mermaid`, the code is HTML-escaped and returned as `<pre class="mermaid">`. No highlighting, framing, or other processing occurs.

**2. ANSI.** If the language is `ansi`, the code is parsed for SGR escape sequences. Each escape maps to foreground/background colors and font style bits. The result is a set of token lines with explicit colors from the ANSI palette. No highlighter is called.

**3. Diff.** If the language is `diff` and a `lang=` meta token specifies the original language, the diff preprocessor strips `+`/`-`/space prefixes from each line and generates insertion/deletion line markers. The stripped code is then highlighted as the specified language (e.g., Go syntax in a unified diff). The generated markers are appended to any existing line markers.

**4. Normal.** All other languages are tokenized through the highlighter.

## Tokenization and token merging

For dual-theme rendering, both light and dark token streams are needed. If the highlighter implements `DualThemeTokenizer`, a single call produces both streams. Otherwise, the engine calls `Tokenize` twice (once per theme).

Before tokenization, tabs are expanded to spaces according to `WithTabWidth`.

The two token streams are merged into a single set of `MergedToken` lines (`token_merge.go`). Each merged token carries `Content`, `LightColor`, `DarkColor`, `LightBG`, `DarkBG`, and `FontStyle`. The light and dark arrays must have identical token boundaries (same line count and token count per line) for the merge to produce correct results. The `DualThemeTokenizer` contract guarantees this.

If the language is unknown, the engine logs a warning via `WarningHandler` and falls back to plaintext (no colors).

For the CSS mechanism that switches between light and dark token colors, see [Dual-Theme Rendering](/architecture/dual-theme-rendering/).

## HTML rendering

The `RenderBlock` function (`internal/render/render.go`) assembles the final HTML from the merged token lines and the resolved block configuration.

**Wrapper.** Every block is wrapped in `<div class="kazari-block">`. Optional classes include `kz-themed` (when per-block theme override is active) and `kz-collapsed` (when threshold collapse is active). The `not-content` class is always present. Data attributes include `data-lines` (line count) and `data-kz-id` (block identifier for theme toggle state).

**Frame dispatch.** Three frame types produce different HTML structures. For `FrameCode`, Kazari renders a `<figure class="frame">` containing a toolbar (language badge, title, action buttons), then the `<pre><code>` block. For `FrameTerminal`, Kazari renders a `<figure class="frame is-terminal">` containing a terminal header with window dots (colored or minimal), title, and inline action buttons, then the `<pre><code>` block. For `FrameNone`, the `<pre><code>` block is rendered directly with no surrounding chrome.

The toolbar is built in two halves: a left section (language badge, file icon, title) and a right section (copy, wrap, theme toggle, font controls, fullscreen buttons). Each button is only emitted when its corresponding config flag is enabled.

**Post-render callbacks.** After `RenderBlock` returns the assembled HTML string, any callbacks registered via [`WithPostRender`](/plugins/post-render-callbacks/) run in registration order. Each callback receives the HTML string and a `BlockInfo` struct, and returns a (possibly modified) string. This is the final step before the result is returned to the caller.

## Line rendering

The `<pre><code>` element contains one `<div class="kz-line">` per source line. Each line is rendered through one of three paths:

**Collapse range lines.** Lines within an explicit collapse range are wrapped in a `<details>` element (github style) or a `<div class="kz-section">` with a nested `<details>` (collapsible-start/end styles). A summary row shows the expand/collapse icon and a line count. The actual code lines are rendered normally inside the collapsible container.

**Hidden threshold lines.** In threshold-collapsed blocks, lines beyond the preview are rendered with a `kz-hidden` CSS class. Tokens are still emitted so the copy button captures the full source text. A gap indicator row (`kz-gap`) with a vertical ellipsis separates visible segments.

**Visible lines.** Each visible line contains an optional gutter cell (`kz-gutter` with a line number) and a code cell (`kz-code`). The gutter width is computed from the digit count of the highest line number and set as `--kz-ln-width`. If word wrap is active and the line has leading whitespace, an indent span preserves the indentation and `--kz-indent` is set for hanging indent.

## Token rendering

Each token within a line is rendered as a `<span>` with inline CSS custom properties encoding its colors and styles.

**Style properties.** The `buildTokenStyle` function (`internal/render/render_lines.go`) emits only the properties that differ from defaults:

| Property | Purpose |
|----------|---------|
| `--sl` | Light theme token foreground color |
| `--sd` | Dark theme token foreground color (only in dual-theme mode) |
| `--slbg` | Light theme token background color (rare) |
| `--sdbg` | Dark theme token background color (rare) |
| `--sfs` | Font style (`italic`) |
| `--sfw` | Font weight (`bold`) |
| `--std` | Text decoration (`underline`, `line-through`, or both) |

The stylesheet picks these up via `var(--sl, inherit)` with `inherit` as the fallback for tokens that have no explicit color.

**Contrast enforcement.** When `WithMinContrast` is set, each token's foreground color is checked against the effective background. On a marked line, the marker's opaque composite color is the background; on a normal line, the editor background is used. If the WCAG contrast ratio falls below the threshold, the token color is adjusted. Results are memoized per block.

**Annotated tokens.** When inline markers or links are present on a line, tokens are segmented at marker and link boundaries. Marker segments are wrapped in `<mark>`, `<ins>`, or `<del>` elements. Link segments are wrapped in `<a class="kz-link">`. When a marker spans multiple tokens, intermediate segments receive `open-start` and `open-end` classes for continuous visual styling.

## Data flow

```
(code, meta string or Options)
    |
    v
meta.Parse() or mapOptionsToBlockOpts()
    |
    v
Config.Resolve()  ──  engine defaults → language defaults → per-block
    |
    v
preprocess()
    ├── ExtractFileName   → resolved.Title (line removed from code)
    ├── ExtractLinks      → resolved.Links (annotations removed)
    ├── DetectFrameType   → resolved.Frame
    └── RawCode snapshot  → resolved.RawCode (copy button source)
    |
    v
resolveCollapse()  → segments, ranges, threshold
    |
    v
applyThemeOverride()  → inline CSS override vars + marker BGs
    |
    v
Language routing
    ├── mermaid  → <pre class="mermaid">  (return)
    ├── ansi     → ansi.Parse()           → token lines
    ├── diff     → strip prefixes + markers + re-highlight
    └── normal   → tokenize via highlighter → token lines
              |
              v
         mergeTokens()  → merged token lines (light + dark)
              |
              v
      RenderBlock()
          ├── <div class="kazari-block">
          ├── frame dispatch (toolbar / terminal header / none)
          └── <pre><code>
                └── per-line loop
                      ├── collapse range → <details> / <div.kz-section>
                      ├── hidden line   → <div.kz-line.kz-hidden>
                      └── visible line  → <div.kz-line>
                            ├── kz-gutter (line number)
                            └── kz-code
                                  ├── renderAnnotatedToken (markers/links)
                                  └── renderToken → <span style="--sl:...">
              |
              v
      postRenders[]  → each callback(html, BlockInfo) → modified html
```
