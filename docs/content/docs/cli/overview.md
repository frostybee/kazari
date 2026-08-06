---
title: "Overview"
description: "Upgrade plain code blocks in built HTML to framed, syntax-highlighted Kazari blocks after any static site generator runs."
sidebar:
  order: 1
---

`kazari process` takes a folder of built HTML and upgrades every plain code block into a framed, syntax-highlighted Kazari block. The site deploys exactly as before; only the code blocks change.

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

## Why post-build

Kazari integrates into a site through two paths, each with different requirements.

The **library path** renders blocks while the page is built. A Go program calls `engine.RenderWithMeta()` directly, or wires the [Goldmark extension](/integrations/goldmark-extension/) into its Markdown pipeline. Kazari sees the fence info string exactly as the author typed it, so every meta string feature works with no extra setup. This path needs a Go Markdown pipeline that is yours to change.

The **post-build path** exists for sites where the Markdown pipeline is not yours to change. It takes the finished HTML, recovers the source text from whatever markup the generator produced, and re-renders each block through the same engine. Nothing about the generator has to change.

The tradeoff is the fence info string. Generators discard it before HTML exists, so features driven by per-block meta cannot survive on their own. A one-file [render hook](/cli/render-hooks/) writes that text into the HTML and closes the gap where the generator supports one.

| Path | Use when |
|---|---|
| Library or Goldmark extension | The site is built by a Go program whose Markdown pipeline you control |
| `kazari process` | The site is built by any other generator, or the pipeline is not yours to change |

Both paths run the same engine and read the same `kazari.config.yaml` format, and recovered source renders byte-identical to the same source going through Goldmark. Moving between them changes the wiring, not the output.

## How it works

The processor walks every `.html` and `.htm` file under the target directory, finds each code block, recovers the original source text from the markup, and re-renders it through the Kazari engine. Everything outside the replaced blocks passes through untouched. Running the command a second time produces zero changes, which makes it safe as an unconditional build step.

`kazari.css` and `kazari.js` are written once at the output root and linked from every page that gained a Kazari block. See [Configuration](/cli/configuration/#assets) for asset options.

:::note
The processor never follows symlinks. Indentation that a generator replaced with non-breaking spaces (U+00A0) is preserved as-is, because normalizing would corrupt intentional non-breaking spaces inside string literals.
:::

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
- Line markers with explicit ranges, including labeled ranges
- Focus lines (`focus="4-6"`)
- Explicit collapse ranges and collapse styles
- Hybrid diff (`diff lang="go"`)
- Per-block theme overrides, including a light and dark pair
- `startLineNumber` and per-block `showLineNumbers`, `wrap`, `preserveIndent`, or `hangingIndent` overrides
- Output panel controls: `withOutput`, `outputLabel`, `outputCollapsed`

:::note
Two features stay out of reach even with the stock Hugo hook: inline text markers (`"text"`) and regex markers (`/regex/`). Hugo parses only `key="value"` pairs inside the fence's brace group and drops bare tokens, so neither can be written in a Hugo fence at all. Extend the template locally to add them.
:::

## Feature availability by generator

| Generator | Zero setup | Render hook |
|---|---|---|
| Hugo | Yes | Yes, one template file |
| Eleventy | Yes | Best-effort, docs-only snippet |
| Jekyll | Yes | No clean hook mechanism today |
| mdBook, Sphinx, Zola, Astro, hand-written | Yes | Not available |

## Dry-run mode

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

## Try it

[Hugo Integration](/cli/hugo-integration/) walks through building a complete example site and shows every feature above rendered, with a before and after comparison.
