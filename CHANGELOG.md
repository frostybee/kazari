# Changelog

All notable changes to Kazari are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-07-24

### Added

- **Output panel.** The `withOutput` meta key plus a `---output---` separator line split a code block into a syntax-highlighted command section and a plain-text output panel below it. The panel is collapsible and carries its own toggle.
- **Output panel keys in file-based config.** `outputPanel`, `outputDefaultCollapsed`, and `outputSeparator` are now recognized in `kazari.config.yaml` and `kazari.config.json`.
- **`WithPostRender` callback.** Register one or more callbacks that run after each block is rendered. Each receives the HTML string and a `BlockInfo` struct describing the block, and returns the HTML to use. Enables wrapping blocks, injecting elements, and adding attributes that built-in features and CSS cannot reach.
- **`BlockInfo` type.** Metadata passed to `WithPostRender` callbacks.
- **External link affordances.** Clickable code annotations created with `@[text](url)` now render an external link icon and open in a new browser tab.

### Changed

- The `not-content` class is now applied to every rendered block unconditionally. It previously required `WithContentExclusion` to be enabled.
- Scrollbars inside code blocks were restyled with a floating thumb, giving them a thinner appearance.
- The output panel chevron-to-label gap increased from `0.35em` to `0.5em`.

### Fixed

- Threshold-based auto-collapse is now suppressed when any line marker carries a label. Labeled ranges exist to guide the reader through specific sections, so collapsing them by default defeated their purpose. An explicit `collapse` meta key still overrides the guard.
- Per-block theme overrides now produce the correct toolbar background color, remap theme variables across all chrome variables, and no longer double-fade line numbers in themes that use alpha colors.
- The per-block theme toggle now overrides terminal background variables, so terminal-framed blocks respond to it correctly.
- An empty language string is treated as plaintext without emitting a warning.

### Documentation

- Added a syntax highlighting page, an ANSI rendering showcase, and cross-links between related pages on the documentation site.

## [1.0.1] - 2026-06-25

### Fixed

- The collapsible section icon color now resolves through `currentColor`, so it inherits the surrounding text color instead of using a fixed value.
- Demo: the no-frame example uses a Go snippet rather than a mismatched sample.
- Demo: the dark theme switched from `github-dark` to `github-dark-default`.

## [1.0.0] - 2026-06-24

Initial release.

### Added

- **Frames and chrome.** Editor and terminal frames with macOS-style dots (colored or minimal), automatic frame detection from the language, title bars with file icons and language badges, and file name extraction from code comments.
- **Toolbar.** Copy to clipboard, word wrap toggle, fullscreen with font size controls, and a per-block theme toggle.
- **Line markers.** Highlight, insertion, and deletion markers with colored backgrounds and gutter accents, range syntax (`{3, 5-8}`), labeled ranges (`{"API":3-7}`), and documented overlap resolution.
- **Inline markers.** Literal text markers, regex markers with capture groups, and inline `ins=` / `del=` variants that span multiple tokens.
- **Focus lines.** `focus={N-M}` dims every line outside the focused range.
- **Line numbers.** Engine-wide, per-language, and per-block control, with `startLineNumber=N` and an auto-width gutter.
- **Word wrap.** `wrap`, `preserveIndent`, and `hangingIndent=N`.
- **Collapsible sections.** Threshold-based auto-collapse with a preview and expand button, plus range-based `collapse={N-M}`, across four collapse styles.
- **Code groups.** `:::code-group` tabbed containers with tab labels, ARIA tab panel semantics, and cross-group tab sync.
- **Diff highlighting.** `diff lang="go"` strips diff prefixes, generates ins/del markers, and applies syntax highlighting to the cleaned source.
- **ANSI rendering.** SGR escape sequences including the standard 16 colors, the 256-color palette, RGB true color, and text attributes.
- **Mermaid pass-through.** `lang="mermaid"` emits a `<pre class="mermaid">` block for client-side rendering.
- **Inline links.** `@[text](url)` renders clickable links inside code, with URL scheme validation.
- **Dual-theme rendering.** Light and dark token colors are baked into the HTML as CSS custom properties. Theme switching is pure CSS with no flash on toggle or page load.
- **Theme configuration.** `WithThemes`, three dark mode strategies (selector, media query, both), per-block `theme=` override, a theme customizer callback, and OKLCH hue and chroma adjustments.
- **Accessibility.** Per-token WCAG contrast enforcement with a configurable minimum ratio, ARIA labels, live regions, and screen reader announcements.
- **Style isolation.** An `all: revert` reset, an optional `@layer kazari` cascade layer, and a `not-content` class for Tailwind Typography.
- **CSS custom properties.** Over 90 `--kz-*` variables covering every visual property, overridable from a stylesheet or through `WithStyleOverrides` and `WithThemedStyleOverrides`.
- **Pluggable highlighter.** A `Highlighter` interface with adapters for Nuri and Chroma, plus an optional `DualThemeTokenizer` interface that halves dual-theme tokenization cost.
- **Goldmark extension.** `kazarimd.New(engine)` for fenced code blocks and `kazarimd.CodeGroups(engine)` for `:::code-group` support.
- **File-based configuration.** Auto-discovery of `kazari.config.yaml`, `.yml`, and `.json` with unknown-field rejection.
- **Localization.** Built-in `en-US`, `fr-FR`, and `ja-JP` locales with per-string overrides.

## Releasing

The changelog is maintained by hand. To cut a release:

1. Add entries under `[Unreleased]` as work lands on `main`.
2. Rename `[Unreleased]` to the new version with today's date, add a fresh empty `[Unreleased]` above it, and update the link references at the bottom of this file.
3. Mirror the new section into `docs/content/docs/changelog.md`.
4. Tag the release: `git tag -a vX.Y.Z -m "vX.Y.Z"` and push it with `git push origin vX.Y.Z`.
5. Publish the GitHub release with the notes from this file: `gh release create vX.Y.Z --notes-file <section>`.

[Unreleased]: https://github.com/frostybee/kazari/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/frostybee/kazari/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/frostybee/kazari/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/frostybee/kazari/releases/tag/v1.0.0
