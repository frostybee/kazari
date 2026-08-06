---
title: "Kazari with Hugo"
---

This is a small, complete Hugo site that shows what `kazari process` does
to a built site. Hugo renders plain code blocks. Kazari walks the output
afterwards and upgrades every one of them into a framed, syntax
highlighted block, then writes `kazari.css` and `kazari.js` next to the
HTML and injects the tags that load them.

## Run it

```bash
hugo --minify
kazari process --config kazari.config.yaml public
```

Open `public/index.html` in a browser. To see the difference, run only the
first command and open the same page: the blocks are unstyled, because
this site never links `kazari.css` by hand. The processor adds those tags
while rewriting the blocks.

## What to look at

`layouts/_markup/render-codeblock.html` is the render hook. It is a
verbatim copy of `integrations/hugo/render-codeblock.html` and it is the
one file that carries per block options through the build. Without it,
Kazari still upgrades every block, but only the language survives; titles,
marked lines, collapse ranges, and the rest are lost before the processor
ever sees the page.

`kazari.config.yaml` holds the site wide settings. Every fence in
`content/` puts its options in a single brace group right after the
language, as `key="value"` pairs, which is the only shape Hugo passes
through to the hook. A fence opens with three backticks followed by
`go {title="main.go" mark="3-5" collapse="10-20"}`. Hugo discards any info
string text outside the braces and rejects bare tokens such as `{3-5}`, so
every option needs a key.

## Pages

- **Basics** covers titles, filenames, line numbers, and frames.
- **Annotations** covers marked, inserted, deleted, and focused lines.
- **Collapse** covers the automatic threshold and explicit ranges.
- **Themes** covers the dual theme output and the two kinds of toggle.
- **Output** covers the code and output split.

## Copy list

Four things move to a real site:

1. `layouts/_markup/render-codeblock.html`, unchanged.
2. `kazari.config.yaml`, trimmed to the settings that site wants.
3. The `kazari process` step, after `hugo` in the build.
4. A site wide dark mode switch, if the site has one. `static/site.css`
   and `static/theme-toggle.js` here are an example of that, not part of
   Kazari.
