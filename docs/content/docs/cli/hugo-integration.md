---
title: "Hugo Integration"
description: "Build a real Hugo site, run kazari process over it, and see every per-block feature render."
sidebar:
  order: 2
---

The Kazari repository includes a complete Hugo site at `examples/hugo/` that exercises every feature the processor supports. Building it and processing it takes two commands. The result is a set of pages that can be opened straight from disk, with no server, and compared against the unprocessed build to see exactly what Kazari adds.

This page walks through that example, then covers what to copy when adding Kazari to a real Hugo site.

## The example site

Five content pages cover every processor feature, grouped by concern:

| Page | Covers |
|---|---|
| Basics | Titles, filename comment extraction, line numbers, terminal and editor frame detection, forced frames, word wrap, Mermaid pass-through |
| Annotations | Marked, inserted, and deleted ranges, Hugo's native `hl_lines`, focus lines, hybrid diff, links inside code |
| Collapse | The config threshold, `nocollapse`, explicit ranges, collapse styles, a per-block threshold |
| Themes | Dual themes, per-block theme pinning, a light and dark pair, and the two kinds of toggle |
| Output | The output panel, custom labels, and a panel that starts collapsed |

Focus lines and the output panel are the two worth opening first. Both depend entirely on the render hook carrying their keys through the build, so seeing them render is the clearest proof the hook is installed and working.

The home page lists the two commands and the four files that move to a real site.

## Build and process the example

From a clone of the Kazari repository:

```bash
cd examples/hugo
hugo --minify
```

Open `public/index.html`. Every code block is a plain, unstyled `<pre>`. The site never links `kazari.css` in its templates, so there is nothing to style the blocks.

Now run the processor:

```bash
go run ../../cmd/kazari process --config kazari.config.yaml public
```

With the Kazari processor installed, that becomes `kazari process --config kazari.config.yaml public`.

Open the same `public/index.html` again. Every code block is now framed, syntax-highlighted, and interactive, with a copy button, line numbers, and a theme toggle in the toolbar. The processor wrote `kazari.css` and `kazari.js` into the output root and injected the tags that load them. Nothing else on the page changed.

The processor prints one summary line:

```
6 files, 33 blocks upgraded, 1 skipped, 0 suppressed, 8 changed
```

The skipped block is a Mermaid diagram, which passes through untouched so a client-side renderer can claim it.

Running the processor a second time changes nothing, which is what makes it safe as an unconditional build step. The `--check` flag verifies this without writing:

```bash
go run ../../cmd/kazari process --check --config kazari.config.yaml public
```

`--check` writes nothing and exits `0` when a real run would produce no changes.

:::note
`hugo server` serves pages from memory and never writes to disk, so the processor has nothing to read. Build with `hugo --minify`, process the output, and open the files directly or serve `public/` with a static file server.
:::

## How the example works

The integration relies on three files.

**`layouts/_markup/render-codeblock.html`** is the render hook. It bypasses Chroma and writes each fence's options into a `data-kz-meta` attribute, which is the only way per-block options survive a Hugo build. Without it Kazari still upgrades every block, but only the language survives. See [Render Hooks](/cli/render-hooks/) for the full key table.

**`kazari.config.yaml`** holds the site-wide settings: the theme pair, the dark mode selector, the collapse threshold, and the feature toggles. See [Configuration](/cli/configuration/) for every key.

Everything else is an ordinary Hugo site. The processor reads the built output and needs nothing from the generator beyond the attribute.

## Adding Kazari to a real site

Four steps:

1. Copy `integrations/hugo/render-codeblock.html` to `layouts/_markup/render-codeblock.html`. On Hugo older than 0.146, use `layouts/_default/_markup/render-codeblock.html`.
2. Copy `kazari.config.yaml` into the project root and trim it to the settings the site needs.
3. Add `kazari process` after `hugo` in the build pipeline. See [Build Pipelines](/cli/ci-recipes/).
4. Point `darkMode.selector` at whatever class or attribute the site already uses for dark mode.

### Dark mode

Kazari ships no page-level theme switch. The button in each block's toolbar, enabled with `themeToggleButton`, overrides one block and never touches `<html>`. A site that wants page-wide dark mode writes its own toggle, then tells Kazari which selector marks the dark state:

```yaml
darkMode:
  kind: selector
  selector: ".dark"
```

`kind: selector` ignores the operating system preference and trusts the class. The `selector` key is required whenever `kind` is `selector` or `both`. Kazari sets `transition: none !important` inside code blocks to prevent the usual `* { transition }` shortcut from fading token colors on each switch. An inline script in the document head that applies the stored class before first paint completes the no-flash setup. The example site's `layouts/_default/baseof.html` shows the pattern. See [Themes and Dark Mode](/styling/themes-and-dark-mode/) for the full strategy list.

### Language icons

Every language badge can open with an icon slot at the far left of the toolbar. Set `languageIconMode: iconAndText` in the config and supply one CSS rule per language:

```css
.kz-lang-icon[data-lang="go"] {
  background-image: url("/icons/go.svg");
}
```

Kazari sizes the slot and scales the background to fit. Add a rule without a language selector as a fallback, because the slot renders for every fence and an unmatched language would otherwise hold blank space. The example site covers its languages with inline SVG data URIs at the bottom of `static/site.css`, which is the one part of that file worth copying.

A second slot, the file icon, works the same way for titles with a file extension, keyed on `data-ext` instead of `data-lang`. The example site turns it off with `fileIcons: false` so each block carries a single icon; a site can run either slot, or both. See [Icons](/features/icons/) for both slots in full.

### What not to copy

The example site's `static/site.css` and `static/theme-toggle.js` are not part of Kazari. The theme toggle exists because the example has no theme of its own; a real site already has one. The only part of `site.css` worth taking is the language icon block described above.

## Next

- [Overview](/cli/overview/) for what works without a render hook
- [Render Hooks](/cli/render-hooks/) for every key the hook recognizes
- [Configuration](/cli/configuration/) for flags and config file keys
