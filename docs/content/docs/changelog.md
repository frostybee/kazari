---
title: "Changelog"
description: "Release notes and version history for Kazari."
sidebar:
  order: 10
---

Every released version of Kazari and what changed in it. Entries follow [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and version numbers follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The canonical copy lives in [`CHANGELOG.md`](https://github.com/frostybee/kazari/blob/main/CHANGELOG.md) at the repository root. Downloadable archives and release notes are on the [GitHub releases page](https://github.com/frostybee/kazari/releases).

## v1.1.0

Released 2026-07-24.

### Added

- **Output panel.** The `withOutput` meta key plus a `---output---` separator line split a code block into a syntax-highlighted command section and a plain-text output panel below it. The panel is collapsible and carries its own toggle. See [Output Panel](/docs/features/output-panel/).
- **Output panel keys in file-based config.** `outputPanel`, `outputDefaultCollapsed`, and `outputSeparator` are now recognized in `kazari.config.yaml` and `kazari.config.json`. See [File-Based Config](/docs/reference/file-based-config/#output-panel).
- **`WithPostRender` callback.** Register one or more callbacks that run after each block is rendered. Each receives the HTML string and a `BlockInfo` struct describing the block, and returns the HTML to use. Enables wrapping blocks, injecting elements, and adding attributes that built-in features and CSS cannot reach. See [Post-Render Callbacks](/docs/plugins/post-render-callbacks/).
- **`BlockInfo` type.** Metadata passed to `WithPostRender` callbacks.
- **External link affordances.** Clickable code annotations created with `@[text](url)` now render an external link icon and open in a new browser tab. See [Links](/docs/features/links/).

### Changed

- The `not-content` class is now applied to every rendered block unconditionally. It previously required `WithContentExclusion` to be enabled. See [Content Exclusion](/docs/features/content-exclusion/).
- Scrollbars inside code blocks were restyled with a floating thumb, giving them a thinner appearance.
- The output panel chevron-to-label gap increased from `0.35em` to `0.5em`.

### Fixed

- Threshold-based auto-collapse is now suppressed when any line marker carries a label. Labeled ranges exist to guide the reader through specific sections, so collapsing them by default defeated their purpose. An explicit `collapse` meta key still overrides the guard. See [Collapsible Sections](/docs/features/collapsible-sections/).
- Per-block theme overrides now produce the correct toolbar background color, remap theme variables across all chrome variables, and no longer double-fade line numbers in themes that use alpha colors.
- The per-block theme toggle now overrides terminal background variables, so terminal-framed blocks respond to it correctly.
- An empty language string is treated as plaintext without emitting a warning.

### Documentation

- Added a [Syntax Highlighting](/docs/features/syntax-highlighting/) page, an ANSI rendering showcase, and cross-links between related pages.

## v1.0.1

Released 2026-06-25.

### Fixed

- The collapsible section icon color now resolves through `currentColor`, so it inherits the surrounding text color instead of using a fixed value.
- Demo: the no-frame example uses a Go snippet rather than a mismatched sample.
- Demo: the dark theme switched from `github-dark` to `github-dark-default`.

## v1.0.0

Released 2026-06-24. Initial release.

### Added

- **Frames and chrome.** Editor and terminal frames with macOS-style dots (colored or minimal), automatic frame detection from the language, title bars with file icons and language badges, and file name extraction from code comments. See [Frames and Titles](/docs/features/frames-and-titles/) and [Icons](/docs/features/icons/).
- **Toolbar.** Copy to clipboard, word wrap toggle, fullscreen with font size controls, and a per-block theme toggle. See [Toolbar](/docs/features/toolbar/).
- **Line markers.** Highlight, insertion, and deletion markers with colored backgrounds and gutter accents, range syntax (`{3, 5-8}`), labeled ranges (`{"API":3-7}`), and documented overlap resolution. See [Line Markers](/docs/features/line-markers/).
- **Inline markers.** Literal text markers, regex markers with capture groups, and inline `ins=` and `del=` variants that span multiple tokens. See [Inline Markers](/docs/features/inline-markers/).
- **Focus lines.** `focus={N-M}` dims every line outside the focused range. See [Focus Lines](/docs/features/focus-lines/).
- **Line numbers.** Engine-wide, per-language, and per-block control, with `startLineNumber=N` and an auto-width gutter. See [Line Numbers](/docs/features/line-numbers/).
- **Word wrap.** `wrap`, `preserveIndent`, and `hangingIndent=N`. See [Word Wrap](/docs/features/word-wrap/).
- **Collapsible sections.** Threshold-based auto-collapse with a preview and expand button, plus range-based `collapse={N-M}`, across four collapse styles. See [Collapsible Sections](/docs/features/collapsible-sections/).
- **Code groups.** `:::code-group` tabbed containers with tab labels, ARIA tab panel semantics, and cross-group tab sync. See [Code Groups](/docs/features/code-groups/).
- **Diff highlighting.** `diff lang="go"` strips diff prefixes, generates ins and del markers, and applies syntax highlighting to the cleaned source. See [Diff Highlighting](/docs/features/diff-highlighting/).
- **ANSI rendering.** SGR escape sequences including the standard 16 colors, the 256-color palette, RGB true color, and text attributes. See [ANSI Rendering](/docs/features/ansi-rendering/).
- **Mermaid pass-through.** `lang="mermaid"` emits a `<pre class="mermaid">` block for client-side rendering. See [Mermaid Pass-Through](/docs/features/mermaid-pass-through/).
- **Inline links.** `@[text](url)` renders clickable links inside code, with URL scheme validation. See [Links](/docs/features/links/).
- **Dual-theme rendering.** Light and dark token colors are baked into the HTML as CSS custom properties. Theme switching is pure CSS with no flash on toggle or page load. See [Dual-Theme Rendering](/docs/architecture/dual-theme-rendering/).
- **Theme configuration.** `WithThemes`, three dark mode strategies (selector, media query, both), per-block `theme=` override, a theme customizer callback, and OKLCH hue and chroma adjustments. See [Themes and Dark Mode](/docs/styling/themes-and-dark-mode/).
- **Accessibility.** Per-token WCAG contrast enforcement with a configurable minimum ratio, ARIA labels, live regions, and screen reader announcements.
- **Style isolation.** An `all: revert` reset, an optional `@layer kazari` cascade layer, and a `not-content` class for Tailwind Typography. See [CSS Custom Properties](/docs/styling/css-custom-properties/).
- **CSS custom properties.** Over 90 `--kz-*` variables covering every visual property, overridable from a stylesheet or through `WithStyleOverrides` and `WithThemedStyleOverrides`. See [CSS Variables](/docs/reference/css-variables/).
- **Pluggable highlighter.** A `Highlighter` interface with adapters for Nuri and Chroma, plus an optional `DualThemeTokenizer` interface that halves dual-theme tokenization cost. See [Custom Highlighters](/docs/plugins/custom-highlighters/).
- **Goldmark extension.** `kazarimd.New(engine)` for fenced code blocks and `kazarimd.CodeGroups(engine)` for `:::code-group` support. See [Goldmark Extension](/docs/integrations/goldmark-extension/).
- **File-based configuration.** Auto-discovery of `kazari.config.yaml`, `.yml`, and `.json` with unknown-field rejection. See [File-Based Config](/docs/reference/file-based-config/).
- **Localization.** Built-in `en-US`, `fr-FR`, and `ja-JP` locales with per-string overrides. See [Localization](/docs/features/localization/).
