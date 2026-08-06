---
title: "Render Hooks"
description: "Install a one-file hook so per-block titles, markers, and collapse ranges survive the build."
sidebar:
  order: 4
---

A one-file render hook stashes a code block's full meta string in an HTML attribute, so `kazari process` recovers titles, line markers, and collapse ranges instead of falling back to config-only defaults. Everything in the [render hook feature list](/cli/overview/#what-needs-a-render-hook) becomes available.

## Why a hook is needed

Static site generators discard the fence info string before HTML exists. A fence like ` ```go title="main.go" {3-5} ` reaches the built page as a plain Go block; the title and the marked lines are gone, and nothing downstream can recover text that was never written out. The hook runs inside the generator, where the metadata still exists, and writes it into the HTML as a `data-kz-meta` attribute.

## Hugo

Hugo supports render hooks natively, and a single template file is all it takes. Copy the canonical template from the Kazari repository at [`integrations/hugo/render-codeblock.html`](https://github.com/frostybee/kazari/blob/main/integrations/hugo/render-codeblock.html) into the site:

| Hugo version | Install path |
|---|---|
| 0.146 and later | `layouts/_markup/render-codeblock.html` |
| older | `layouts/_default/_markup/render-codeblock.html` |

The hook replaces Chroma for every fence: it emits a plain block carrying the meta string, and `kazari process` does the highlighting afterward.

### Authoring convention

Hugo imposes a specific fence syntax. It discards info string text outside the braces, reads only the first brace group, and rejects bare tokens. Write every Kazari option in one brace group right after the language, as `key="value"` pairs only:

````
```go {title="main.go" mark="3-5" collapse="10-20"}
````

This shape is the only one that survives, because Hugo discards info string text outside the braces, reads only the first brace group, and rejects bare tokens like `{3-5}` as a build error. Hugo also lowercases attribute names before the template sees them; values keep their case.

:::note
Two Kazari features cannot be expressed in a Hugo fence at all: inline text markers (`"text"`) and regex markers (`/regex/`). Hugo parses only `key="value"` pairs inside the brace group and drops bare tokens, so neither reaches the template. Extend it locally to add them.
:::

### Recognized keys

The hook recognizes the keys below and translates each into the format `kazari process` expects. Any lowercase key not listed here, such as `theme`, `frame`, or `lang`, passes through unchanged into the `data-kz-meta` attribute, so future Kazari meta keys work without editing the template.

| Key | Example | Effect |
|---|---|---|
| `title` | `title="main.go"` | Frame title |
| `mark` | `mark="3-5,9"` | Marked line ranges, comma separated |
| `mark` (labeled) | `mark="'Setup':3-5"` | A labeled range. Use single quotes, since the value is already inside double quotes |
| `ins` / `del` | `ins="6-7"` `del="9"` | Inserted or deleted line ranges |
| `add` / `rem` | `add="6-7"` `rem="9"` | Accepted spellings of `ins` and `del` |
| `focus` | `focus="4-6"` | Focused line ranges, dimming everything else |
| `collapse` | `collapse="true"` or `collapse="4-9"` | Collapse the whole block, or the given ranges |
| `nocollapse` | `nocollapse="true"` | Opt this block out of the config wide collapse threshold |
| `collapsestyle` | `collapsestyle="collapsible-start"` | Collapse style: `-start`, `-end`, or `-auto` |
| `collapsethreshold` | `collapsethreshold="20"` | Per-block auto-collapse threshold |
| `showlinenumbers` | `showlinenumbers="true"` or `"false"` | Toggle line numbers |
| `startlinenumber` | `startlinenumber="5"` | First displayed line number |
| `wrap` | `wrap="true"` | Enable word wrap (`wrap="false"` is a no-op; the meta grammar has no per-block switch to turn wrap off) |
| `preserveindent` | `preserveindent="false"` | Whether wrapped lines keep the original indent |
| `hangingindent` | `hangingindent="4"` | Extra indent, in spaces, on wrapped lines |
| `withoutput` | `withoutput="true"` | Split the fence at the output separator into a code panel and an output panel. Requires `outputPanel: true` in the config file; without it the separator stays literal text |
| `outputcollapsed` | `outputcollapsed="true"` or `"false"` | Whether the output panel starts closed |
| `outputlabel` | `outputlabel="Result"` | Label on the output panel's toggle |
| `hl_lines` | `hl_lines="2-4 7"` or `hl_lines=[3,4]` | Hugo's native spelling; both forms translate to marked lines |

### A runnable example

Every key in the table above is exercised by a complete Hugo site shipped with Kazari. [Hugo Integration](/cli/hugo-integration/) has the commands and a page-by-page tour of what each one renders.

## Example: collapse specific lines

Line ranges are 1-based and count the lines of the code block, not the lines of the file the code came from. Given this block:

```
1  package main
2
3  import "fmt"
4
5  func run() int {
6      x := 1
7      y := 2
8      return x + y
9  }
```

Folding the function body takes one attribute:

````
```go {collapse="5-9"}
````

The hook writes that into the built HTML:

```html
<pre><code class="language-go" data-kz-meta="go collapse={5-9}">
```

`kazari process` reads the attribute and renders a collapsed section in place of lines 5 through 9. Multiple ranges accept either separator: `collapse="3-4 8-12"` and `collapse="3-4,8-12"` produce the same result. A single line is `collapse="7"`.

The section opens once and stays open. For a section the reader can close again, add a style:

````
```go {collapse="5-9" collapsestyle="collapsible-start"}
````

`collapsible-start` puts the summary row above the revealed lines, `collapsible-end` below, and `collapsible-auto` picks based on where the range sits in the block.

### Length-based collapse without a hook

Exact ranges need the hook, because the fence text that carries them never reaches the built HTML. Collapsing by length does not: it comes from the config file and applies to every block the processor touches. Put `kazari.config.yaml` at the project root and run the Kazari processor from there, since discovery checks the target directory first and then the working directory:

```yaml
collapsible:
  lineThreshold: 20
  previewLines: 8
  defaultCollapsed: true
```

Every block longer than 20 lines now collapses to an 8-line preview with an expand button, with no changes to any Markdown source. This is the collapse mechanism available to sites with no render hook.

The two mechanisms combine. Keep the threshold for long blocks, then use the hook for per-block exceptions: `nocollapse="true"` leaves one long block fully expanded, and `collapse="5-9"` folds an exact range regardless of the block's length.

:::note
Threshold collapse needs the collapse JavaScript, which ships only when `collapsible` is configured. Range-based sections use native `<details>` elements and work with JavaScript disabled.
:::

## Jekyll

Jekyll has no clean equivalent today. Rouge and kramdown offer no per-fence hook that sees the info string, so Jekyll sites stay on the zero-setup tier unless a custom plugin is written. Rather than a fabricated plugin, target the `data-kz-meta` contract below: any Jekyll plugin able to emit the attribute gets the full feature set.

## Eleventy

Eleventy configures markdown-it directly, and markdown-it exposes the full info string to a `highlight` override. The following snippet is docs-only guidance; it has not been verified against a live Eleventy build, so confirm it in the target project:

```js
const md = require("markdown-it")({
  html: true,
  highlight(code, info) {
    const [lang, ...rest] = info.trim().split(/\s+/);
    const meta = rest.length ? `${lang} ${rest.join(" ")}` : lang;
    const esc = (s) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;");
    return `<pre><code class="language-${lang}" data-kz-meta="${esc(meta)}">${esc(code)}</code></pre>`;
  },
});
```

Unlike Hugo, markdown-it hands the raw info string to the override, so native Kazari meta syntax like ` ```go title="main.go" {3-5} ` passes through unmodified.

## The data-kz-meta contract

`kazari process` looks for the attribute, not for a specific generator. Any pipeline that attaches `data-kz-meta` to the block's `pre` or `code` element gets the full meta string treatment:

- The value is the complete meta string, language first: `go title="main.go" {3-5}`.
- The attribute may sit on the region root or on any `pre` or `code` element inside it.
- Standard HTML attribute escaping applies; both named entities (`&quot;`) and numeric references (`&#34;`) decode correctly.
- Source recovery still runs normally; the attribute replaces only the synthesized meta string.
