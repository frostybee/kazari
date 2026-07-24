---
title: "Output Panel"
description: "Separate command output from syntax-highlighted source code in a toggleable panel."
sidebar:
  order: 10
---

The output panel splits a code block into a syntax-highlighted command section and a plain-text output section below it. Add `withOutput` to the meta string and place a `---output---` separator line between the source code and its output.

## Basic output panel

Add `withOutput` to the fence info string and separate commands from output with `---output---`:

````
```bash withOutput
pwd
ls -la
---output---
/usr/home/boba-tan
total 24
drwxr-xr-x  3 boba boba 4096 Jun 28 14:30 .
```
````

```bash withOutput
pwd
ls -la
---output---
/usr/home/boba-tan
total 24
drwxr-xr-x  3 boba boba 4096 Jun 28 14:30 .
```

→ The terminal frame shows `pwd` and `ls -la` with syntax highlighting. A toggleable panel below displays the command output as plain text.

The copy button copies only the command portion. The separator line and output text are excluded.

## Editor frame with output

Output panels work with editor frames too. Attach program output beneath an editor-framed code block:

````
```go withOutput
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Kazari!")
    fmt.Println("Output panels are here.")
}
---output---
Hello, Kazari!
Output panels are here.
```
````

```go withOutput
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Kazari!")
    fmt.Println("Output panels are here.")
}
---output---
Hello, Kazari!
Output panels are here.
```

→ The editor frame renders the Go source with full syntax highlighting. The output panel below shows the program's printed output.

## Collapsed output

Start the output panel hidden with `outputCollapsed`. The reader clicks the toggle to reveal it:

````
```bash withOutput outputCollapsed
echo "Build complete"
---output---
Build complete
```
````

```bash withOutput outputCollapsed
echo "Build complete"
---output---
Build complete
```

→ The output panel starts hidden. The toggle button reads "Output" and expands the panel on click.

## Per-block overrides

Engine-level defaults apply to all output panels. Per-block meta keys override them for individual blocks.

Set all output panels to start collapsed at the engine level:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithOutputPanel(true),
    kazari.WithOutputCollapsed(true),
)
```

Then expand a specific block by adding `outputCollapsed=false` to its meta string:

````
```bash withOutput outputCollapsed=false
curl -s https://api.example.com/health
---output---
{"status": "ok"}
```
````

```bash withOutput outputCollapsed=false
curl -s https://api.example.com/health
---output---
{"status": "ok"}
```

→ All output panels on the page start collapsed except this one, which starts expanded.

## Custom separator

Change the separator string at the engine level with `WithOutputSeparator`:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithOutputPanel(true),
    kazari.WithOutputSeparator("==="),
)
```

Blocks using this engine split on `===` instead of `---output---`.

## Custom label

Replace the default toggle label with `outputLabel="..."`:

````
```bash withOutput outputLabel="Run result"
node index.js
---output---
Server listening on port 3000
```
````

```bash withOutput outputLabel="Run result"
node index.js
---output---
Server listening on port 3000
```

→ The toggle button displays "Run result" instead of the default locale label.

## Engine configuration

The output panel requires `WithOutputPanel(true)` on the engine. Without it, per-block `withOutput` meta tokens are parsed but have no effect, and no output panel CSS or JS is emitted.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithOutputPanel(true),
    kazari.WithOutputCollapsed(true),
)
```

**Engine options:**

| Option | Default | Description |
|---|---|---|
| `WithOutputPanel(bool)` | `false` | Master enable for the feature and its CSS/JS |
| `WithOutputCollapsed(bool)` | `false` | Default collapsed state for all output panels |
| `WithOutputSeparator(string)` | `"---output---"` | Separator line that splits code from output |

Per-block overrides via the `Options` struct:

```go
html, err := engine.Render(code, kazari.Options{
    Lang:            "bash",
    WithOutput:      &outputOn,
    OutputCollapsed: &collapsed,
    OutputLabel:     "Run result",
})
```

All three engine options also have `kazari.config.yaml` equivalents: `outputPanel`, `outputDefaultCollapsed`, and `outputSeparator`. See [File-Based Config](/reference/file-based-config/#output-panel) for details.

## How it works

The engine scans the code content for the separator line during preprocessing. The first line whose trimmed content exactly matches the separator string splits the block: everything before becomes the syntax-highlighted code, everything after becomes the output panel text. The separator line itself is removed from both sections.

Splitting happens before `RawCode` is set, before line count is computed, and before collapse resolution runs. This means:

- The copy button and `BlockInfo.RawCode` (passed to post-render callbacks) contain only the code portion.
- Line numbers, markers, focus lines, and collapsible sections apply only to the code lines.
- Collapse threshold counts only code lines, not output lines.

Output text is plain HTML-escaped text rendered into a bare `<pre>`. It receives no syntax highlighting, no line numbers, and no `kz-line` span structure.

If no separator line is found in the block, the entire content is treated as code and no output panel is rendered.

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `withOutput` | Meta string | off | Enable the output panel for this block |
| `outputCollapsed` / `outputCollapsed=false` | Meta string | `false` | Start the panel hidden or expanded |
| `outputLabel="..."` | Meta string | locale | Custom toggle button label |
| `Options.WithOutput` / `OutputCollapsed` / `OutputLabel` | Go API (per-block) | nil | Per-block overrides |
| `WithOutputPanel(bool)` | Go API / `outputPanel` | `false` | Master enable for the feature |
| `WithOutputCollapsed(bool)` | Go API / `outputDefaultCollapsed` | `false` | Default collapsed state |
| `WithOutputSeparator(string)` | Go API / `outputSeparator` | `"---output---"` | Separator line |

The default toggle label is locale-aware. English uses "Output", French uses "Sortie", Japanese uses "出力". Override the label globally with `WithUIStrings(map[string]string{"output.label": "..."})` or per-block with `outputLabel="..."`.

## CSS variables

| Variable | Default | Description |
|---|---|---|
| `--kz-output-bg` | `var(--kz-editor-bg)` | Panel background |
| `--kz-output-fg` | `var(--kz-editor-fg)` | Panel text color |
| `--kz-output-opacity` | `0.7` | Panel opacity |
| `--kz-output-border` | `rgba(128,128,128,0.2)` | Header top border color |
| `--kz-output-header-bg` | `var(--kz-toolbar-bg)` | Header background |
| `--kz-output-label-color` | `var(--kz-toolbar-fg)` | Toggle button label color |

These variables are only emitted when `WithOutputPanel(true)` is set. See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- `WithOutputPanel(true)` must be set on the engine. Per-block `withOutput` is inert without it.
- All three engine options have `kazari.config.yaml` equivalents (`outputPanel`, `outputDefaultCollapsed`, `outputSeparator`).
- Output panel fields are not part of `BlockDefaults` or `LanguageDefaults`. Per-language defaults for output panels are not possible.
- The copy button always copies only the code portion. The separator and output text are excluded.
- Output text is never syntax-highlighted, regardless of the block's language.
- Line numbers, markers, focus lines, and collapsible sections apply only to the code portion.
- The wrap toggle button affects only the code `<pre>`, not the output `<pre>`.
- Output panels work inside code group tabs, processed by the same `RenderWithMeta` pipeline.
- Mermaid blocks bypass output splitting entirely.
- Only the first separator line is used. If the separator string appears multiple times, everything after the first occurrence is treated as output.
