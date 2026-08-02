---
title: "Overview"
description: "Upgrade plain code blocks in built HTML to framed, syntax-highlighted Kazari blocks after any static site generator runs."
sidebar:
  order: 1
---

`kazari process ./public` walks a folder of built HTML files and upgrades every plain code block to a framed, syntax-highlighted Kazari block with a copy button, line numbers, and dual light and dark themes. The site deploys exactly as before; only the code blocks change.

## Install

Requires Go 1.25 or later.

```bash
go install github.com/frostybee/kazari/cmd/kazari@latest
```

Or run it without installing:

```bash
go run github.com/frostybee/kazari/cmd/kazari@latest process ./public
```

## Usage

```
kazari upgrades code blocks in built HTML to framed, syntax
highlighted blocks with copy buttons, line numbers, and dual themes.

Usage:
  kazari process [dir] [flags]   upgrade code blocks under dir (default ".")
  kazari themes                  list bundled syntax theme names
  kazari version                 print the kazari version
  kazari help                    show this help
```

| Subcommand | Purpose |
|---|---|
| `process [dir]` | Upgrade code blocks in every HTML file under `dir` (default `.`) |
| `themes` | List every bundled syntax theme name, one per line |
| `version` | Print the kazari version |
| `help` | Show usage |

## How it works

The processor visits `.html` and `.htm` regular files only and never follows symlinks. In each file it finds candidate code blocks (`pre > code` pairs, plus `div` wrappers with highlight classes for the Pygments family), recovers the original source text from the markup, re-renders it through the same engine the Go library uses, and splices the result back byte for byte. Everything outside the replaced blocks passes through untouched, so minified output stays minified. Running the command a second time reports zero changes.

Because the input is a folder of HTML, the generator that produced it does not matter: Hugo, Jekyll, Eleventy, mdBook, Sphinx, Zola, Astro, and hand-written sites all work. The tool upgrades the blocks already present; it never adds content.

`kazari.css` and `kazari.js` are written once at the output root and linked from every page that gained a Kazari block. See [Configuration](/cli/configuration/#assets) for asset options.

## What works with zero setup

These features need no configuration and no changes to the source site, because their triggers survive the build: engine config, the code text itself, or generator markup that gets translated.

- Editor and terminal frames with automatic detection
- Copy, fullscreen, and wrap buttons
- Theme toggle and dual light and dark themes
- Language badge and file icons
- Line numbers and word wrap via config defaults
- Threshold-based collapsible sections
- Terminal comment stripping, links, and localization
- Filename comment extraction (a `// main.go` first line becomes the title)
- Output panel via the `---output---` separator line
- Hugo and Chroma `hl_lines` classes translated to marked lines
- Plain diff rendering for `language-diff` blocks
- Mermaid blocks detected and left untouched

## What needs a render hook

Features driven by the fence meta string cannot survive a normal build: generators discard that text before HTML exists. A one-file [render hook](/cli/render-hooks/) preserves them.

- Explicit titles (`title="main.go"`)
- Line markers with explicit ranges, labeled ranges, focus lines
- Explicit collapse ranges and collapse styles
- Hybrid diff (`diff lang="go"`)
- `startLineNumber` and per-block `showLineNumbers` or `wrap` overrides
- Output panel via `withOutput`

> [!NOTE]
> Inline text and regex markers are not supported by the stock Hugo hook template even in tier 2. The meta grammar accepts them, but the shipped template does not emit them; extend the template locally to add them.

## Feature availability by generator

| Generator | Zero setup | Render hook |
|---|---|---|
| Hugo | Yes | Yes, one template file |
| Eleventy | Yes | Best-effort, docs-only snippet |
| Jekyll | Yes | No clean hook mechanism today |
| mdBook, Sphinx, Zola, Astro, hand-written | Yes | Not available |

## Check before writing

`--check` computes every change without writing anything. It prints the path of each file and asset that would change, one per line, then a summary:

```
public/docs/index.html
public/kazari.css
public/kazari.js
14 files, 22 blocks upgraded, 1 skipped, 0 suppressed, 3 changed
```

| Exit code | Meaning |
|---|---|
| `0` | Nothing to do, or processing succeeded |
| `1` | `--check` found pending changes |
| `2` | Errors: bad arguments, config parse failure, unknown theme, or any per-file error |

The same summary line prints on every run: files scanned, blocks upgraded, blocks skipped, blocks suppressed inside existing Kazari markup, and files changed.
