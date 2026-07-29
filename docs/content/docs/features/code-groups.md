---
title: "Code Groups"
description: "Tabbed code block interface with synced selection across groups."
tags: [markdown, client-side-js, accessibility]
sidebar:
  order: 9
---

Code groups display multiple code blocks as a tabbed interface with synced selection across groups. Each group uses the `:::code-group` container syntax and requires the [Goldmark extension](/integrations/goldmark-extension/).

## Basic code group

Wrap multiple fenced code blocks inside a `:::code-group` container:

```markdown
:::code-group

```go title="main.go"
package main

import "fmt"

func main() {
    fmt.Println("Hello from Go!")
}
```

```python title="main.py"
def main():
    print("Hello from Python!")

if __name__ == "__main__":
    main()
```

```javascript title="index.js"
function main() {
  console.log("Hello from JavaScript!");
}

main();
```

:::
```

→ A tabbed interface with "main.go", "main.py", and "index.js" tabs. The first tab is active by default. Each panel contains a fully rendered Kazari code block with all features.

## Tab labels

Tab labels are derived in priority order:

| Priority | Source | Example |
|---|---|---|
| 1 | `title=` in the fence meta string | `title="main.go"` |
| 2 | File path comment in the first line | `// src/greet.js` |
| 3 | Capitalized language name | `go` becomes "Go" |
| 4 | Fallback | "Code" (when no language is set) |

File path comments follow the same extraction rules as [Frames & Titles](/features/frames-and-titles/). The comment line is removed from the rendered output.

## Tab sync

Add a `sync=` attribute to keep groups in sync. When a tab is selected in one group, the matching tab in all other groups with the same sync key is activated automatically:

```markdown
:::code-group sync="language"

```go
go get github.com/example/pkg
```

```python
pip install example-pkg
```

:::

:::code-group sync="language"

```go
import "github.com/example/pkg"
```

```python
import example_pkg
```

:::
```

→ Selecting "Go" in the first group also selects "Go" in the second group.

Sync behavior:

- Tabs are matched by their exact label text, not by language identifier or tab index. Groups must use identical tab labels for sync to connect them.
- The active tab is persisted in `localStorage` per sync key and restored on page load.
- Groups with different sync keys sync independently.
- A group without `sync=` does not participate in syncing.

## Keyboard navigation

The tab bar follows the ARIA tab panel pattern:

| Key | Action |
|---|---|
| Arrow Right | Next tab (wraps to first) |
| Arrow Left | Previous tab (wraps to last) |
| Home | First tab |
| End | Last tab |

Tabs use the roving tabindex pattern: the active tab has `tabindex="0"`, all others have `tabindex="-1"`. ARIA IDs linking tabs to panels (`aria-controls`, `aria-labelledby`) are assigned client-side by JavaScript on page load.

## Goldmark setup

Register both the base code block renderer and the code group extension:

```go
md := goldmark.New(goldmark.WithExtensions(
    kazarimd.New(engine),
    kazarimd.CodeGroups(engine),
))
```

`kazarimd.CodeGroups(engine)` calls `engine.EnableCodeGroups()` internally. This must be called during setup, before the engine is used for concurrent rendering or before calling `engine.CSS()` and `engine.JS()`, so that the code group styles and scripts are included in the output.

There is no `WithCodeGroups` engine option and no config file equivalent. Code groups are enabled exclusively through the Goldmark extension.

Full setup details are on the [Goldmark Extension](/integrations/goldmark-extension/) page.

## Per-block features

Each tab inside a code group is rendered through the normal `RenderWithMeta` pipeline. All Kazari features work per-tab: line numbers, markers, focus lines, collapsible sections, output panels, frame overrides, and theme overrides. Meta string options are set independently on each fenced code block within the group.

```markdown
:::code-group

```go title="main.go" showLineNumbers {3-5}
package main

func main() {
    fmt.Println("Hello")
}
```

```python title="main.py" ins={2}
def main():
    print("Hello")
```

:::
```

## CSS variables

### Static variables

| Variable | Default | Description |
|---|---|---|
| `--kz-group-tab-bg` | `transparent` | Tab background |
| `--kz-group-tab-fg` | `inherit` | Tab text color |
| `--kz-group-tab-active-border` | `#007acc` | Active tab bottom border |
| `--kz-group-tab-padding` | `0.5rem 1rem` | Tab padding |
| `--kz-group-border-width` | `1px` | Group container border width |
| `--kz-group-radius` | `var(--kz-radius)` | Group container border radius |

### Luminance-derived variables

These values adapt automatically based on each theme's editor background brightness:

| Variable | Light theme | Dark theme |
|---|---|---|
| `--kz-group-tab-active-bg` | `rgba(255,255,255,0.8)` | `rgba(255,255,255,0.1)` |
| `--kz-group-tab-active-fg` | `#24292f` | `#e6edf3` |
| `--kz-group-border` | `rgba(0,0,0,0.1)` | `rgba(255,255,255,0.1)` |

These three variables have no built-in defaults at the `:root` level. They are computed per theme from the editor background luminance and emitted inside each theme's scoped CSS block.

All code group CSS variables are only emitted when `EnableCodeGroups()` has been called. See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- A group with no fenced code children renders nothing. Non-code content (paragraphs, headings) inside `:::code-group` is silently dropped.
- A single-block group renders the full tab/panel structure with one tab.
- Nested `:::code-group` containers are not supported. The inner directive is treated as text.
- Standalone code blocks outside groups render normally, independently.
- `aria-label="Code variants"` on the tab list is hardcoded English and is not localizable via `WithLocale` or `WithUIStrings`.
- Tab sync matches by label text only. If two synced groups use different labels for the same language (e.g., "Go" vs "main.go"), they will not sync with each other.
- There is no programmatic Go API for building code groups. Groups are created exclusively through the `:::code-group` Markdown syntax processed by the Goldmark extension.
