---
title: "Overview"
description: "How Kazari's extension points work and why there is no plugin system."
sidebar:
  order: 0
---

Kazari has no plugin system and no AST to hook into. Code blocks are rendered as HTML strings through a fixed pipeline. Extension happens at the boundaries of that pipeline: callbacks that transform the output string, interfaces that supply the token input, and DOM attributes that expose block metadata to client-side code.

## Why no plugin system

Kazari builds HTML directly from token streams with no intermediate tree structure. There is no AST to walk, no node types to extend, and no lifecycle to hook into. The few use cases that would require rendering-level access (TypeScript type overlays, animated code transitions) are complex enough to warrant first-class engine support rather than a generic hook API.

## Extension points

Kazari provides six Go-side hooks and a client-side DOM contract. The table below lists each one and where to find its full documentation.

| Extension point | What it does | Documentation |
|---|---|---|
| `WithPostRender` | Modify rendered HTML per block via callbacks | [Post-Render Callbacks](/plugins/post-render-callbacks/) |
| `Highlighter` / `DualThemeTokenizer` / `Tokenize` | Swap syntax highlighting engines or build fully custom HTML from raw tokens | [Custom Highlighters](/plugins/custom-highlighters/) |
| CSS selectors + `data-*` attributes + JS | Style or script blocks client-side, no Go rebuild | [Client-Side Extensibility](/plugins/client-side-extensibility/) |
| `WithThemeCustomizer` | Rewrite extracted theme colors per theme name | [Themes and Dark Mode](/styling/themes-and-dark-mode/) |
| `WithFileIconResolver` | Return custom HTML markup for file icons by extension | [Icons](/features/icons/) |
| `WithWarningHandler` | Intercept non-fatal warnings (unknown languages, config errors) | [Configuration Options](/reference/configuration-options/) |
