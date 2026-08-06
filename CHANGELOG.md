# Changelog

All notable changes to Kazari are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **`kazari process` CLI.** `kazari process [dir]` walks a folder of built HTML files and upgrades plain code blocks in place to framed, syntax-highlighted Kazari blocks with copy buttons, line numbers, and dual themes. It works on the output of any static site generator, including Hugo, Jekyll, Eleventy, mdBook, Sphinx, Zola, and Astro, with no Markdown pipeline integration. `--check` reports pending changes without writing and exits 1 when any exist.
- **The `process` package.** The engine behind the CLI is importable as `github.com/frostybee/kazari/process`, so Go programs can run the same post-build upgrade without shelling out.
- **Goldmark-parity source recovery.** Code recovered from built HTML, whether plain blocks or highlighted markup from Chroma, Rouge, Prism, or Pygments, renders byte-identical to the same source going through the Goldmark path. Hugo `hl_lines` classes translate to Kazari line markers.
- **Hugo render hook.** A `render-codeblock.html` template, shipped at `integrations/hugo/render-codeblock.html`, stashes the full Kazari meta string in a `data-kz-meta` attribute so per-block options survive the build instead of falling back to config defaults. It translates `title`, `mark` (plain and labeled), `ins`, `del`, `add`, `rem`, `focus`, `collapse`, `nocollapse`, `collapsestyle`, `collapsethreshold`, `showlinenumbers`, `startlinenumber`, `wrap`, `preserveindent`, `hangingindent`, `withoutput`, `outputcollapsed`, `outputlabel`, and both forms of Hugo's native `hl_lines`. Keys with no dedicated handler, such as `theme`, `frame`, and `lang`, pass through unchanged.
- **Example Hugo site.** `examples/hugo` is a complete, runnable Hugo site demonstrating the render hook, an annotated config file, and a site-wide dark mode switch, across pages covering frames, annotations, collapse, themes, and the output panel. `cmd/kazari/example_site_test.go` builds it with a real Hugo binary and processes it through the CLI's own entry point.
- **Documented tier limits.** Meta-string-only features (focus lines, labeled ranges, explicit collapse ranges, hybrid diff, per-block overrides, output panel controls) require the render hook; hook-less pipelines still get frames, copy buttons, line numbers, and dual themes with zero configuration. Inline text markers and regex markers cannot be expressed in a Hugo fence at all, since Hugo parses only `key="value"` pairs inside the brace group.

### Fixed

- **Unstyled collapse sections without `WithCollapsible`.** A `collapse={N-M}` meta token renders a `<details>` section on any engine, but the collapse stylesheet and its `--kz-collapse-*` variable defaults, including the expand and collapse icons, were only emitted when `WithCollapsible` was configured. The section rendered with browser default styling and no icons. Both are now always part of `CSS()`. The collapse JavaScript stays conditional, since only threshold-based collapse needs it and threshold markup never renders without the config.

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
