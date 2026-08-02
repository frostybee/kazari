---
title: "Collapsible Sections"
description: "Collapse long code blocks or specific line ranges with threshold-based and range-based folding."
tags: [client-side-js, accessibility]
sidebar:
  order: 8
---

Collapsible sections reduce the visual footprint of long code blocks. Kazari supports two collapse mechanisms: threshold-based (auto-collapse when a block exceeds a line count) and range-based (collapse specific line ranges).

## Enabling collapsible sections

Range-based collapse works on any engine: a `collapse={N-M}` meta token renders a native `<details>` section, and the collapse stylesheet is always part of `CSS()`. `WithCollapsible` enables threshold-based auto-collapse, its JavaScript module, and the threshold defaults.

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithThemes("github-light", "github-dark"),
    kazari.WithCollapsible(kazari.CollapsibleConfig{
        LineThreshold: 15,
        PreviewLines:  8,
    }),
)
```

Or in the config file:

```yaml
collapsible:
  lineThreshold: 15
  previewLines: 8
  defaultCollapsed: true
```

Setting `WithCollapsible(CollapsibleConfig{})` with zero values enables the feature with internal fallback defaults (`LineThreshold: 15`, `PreviewLines: 8`).

**CollapsibleConfig fields:**

| Field | Type | Default | Description |
|---|---|---|---|
| `LineThreshold` | `int` | `15` | Line count to trigger auto-collapse |
| `PreviewLines` | `int` | `8` | Lines shown in collapsed preview |
| `DefaultCollapsed` | `bool` | `true` | Start in collapsed state |
| `PreserveIndent` | `bool` | `true` | Align summary row with code indentation (range mode) |
| `Style` | `CollapseStyle` | `CollapseGithub` | Default style for range-based collapse |
| `ExpandButtonText` | `string` | locale | Custom expand button text |
| `CollapseButtonText` | `string` | locale | Custom collapse button text |
| `ExpandedAnnouncement` | `string` | locale | Screen reader announcement on expand |
| `CollapsedAnnouncement` | `string` | locale | Screen reader announcement on collapse |

## Threshold-based collapse

When a code block exceeds `LineThreshold` lines, it auto-collapses to a preview of the first `PreviewLines` lines. A gradient overlay fades the bottom, and an expand button reveals the full code.

````
```go title="server.go" collapseThreshold=20
func main() {
    cfg := loadConfig()
    db := connectDB(cfg)
    defer db.Close()

    router := chi.NewRouter()
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)

    router.Get("/health", healthHandler)
    router.Get("/api/users", usersHandler(db))
    router.Post("/api/users", createUserHandler(db))
    router.Get("/api/users/{id}", getUserHandler(db))
    router.Put("/api/users/{id}", updateUserHandler(db))
    router.Delete("/api/users/{id}", deleteUserHandler(db))

    router.Get("/api/posts", postsHandler(db))
    router.Post("/api/posts", createPostHandler(db))
    router.Get("/api/posts/{id}", getPostHandler(db))

    log.Printf("Starting server on :8080")
    http.ListenAndServe(":8080", router)
}
```
````

```go title="server.go" collapseThreshold=20
func main() {
    cfg := loadConfig()
    db := connectDB(cfg)
    defer db.Close()

    router := chi.NewRouter()
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)

    router.Get("/health", healthHandler)
    router.Get("/api/users", usersHandler(db))
    router.Post("/api/users", createUserHandler(db))
    router.Get("/api/users/{id}", getUserHandler(db))
    router.Put("/api/users/{id}", updateUserHandler(db))
    router.Delete("/api/users/{id}", deleteUserHandler(db))

    router.Get("/api/posts", postsHandler(db))
    router.Post("/api/posts", createPostHandler(db))
    router.Get("/api/posts/{id}", getPostHandler(db))

    log.Printf("Starting server on :8080")
    http.ListenAndServe(":8080", router)
}
```

→ The block auto-collapses because it exceeds 20 lines. A gradient overlay and expand button appear at the bottom.

### Per-block threshold

`collapseThreshold=N` overrides the engine `LineThreshold` for a single block. The block auto-collapses only when its line count exceeds N.

### Force and suppress

Use `collapse` to force threshold collapse regardless of line count. Use `nocollapse` to suppress threshold collapse for a specific block.

Priority order: `nocollapse` > `collapse` > `collapseThreshold=N` > engine `LineThreshold`.

The toolbar also shows a chevron toggle button on threshold-collapsed blocks. Both the bottom bar button and the toolbar chevron stay in sync; clicking either one toggles the same collapsed state. See [Toolbar](/features/toolbar/) for details.

## Labeled markers suppress auto-collapse

:::note
When a block contains labeled line markers (e.g., `{"Step 1":3-5}`), auto-threshold collapse is suppressed to keep the labels visible. Only labeled markers trigger this behavior. Plain highlight, insertion, and deletion markers do not affect auto-collapse.

Use the `collapse` flag to override this and force collapse even when labels are present.
:::

See [Line Markers](/features/line-markers/) for labeled marker syntax.

## Marker-aware preview

When a threshold-collapsed block has marked or focused lines beyond the normal preview region, the preview extends to include them. The extension follows these rules:

- **Within 2x cap**: marked or focused lines within `2 * PreviewLines` are shown as additional segments with one line of context on each side. Non-contiguous segments are separated by `⋮` gap indicators.
- **Beyond 2x cap**: marked lines too far away are not shown. Instead, the expand button displays a count: `"Show more (+3 highlighted)"`.

This happens automatically when threshold-based collapse and line markers or focus lines are combined on the same block:

````
```go {10-11} collapse
// 28-line block with markers on lines 10-11
```
````

The preview shows the first few lines, then a gap indicator, then the marked lines with surrounding context. The expand button reveals the full block. If marked lines fall beyond the 2x cap, the button displays a badge like `"Show more (+3 highlighted)"` instead.

## Range-based collapse

Collapse specific line ranges with `collapse={N-M}`. The collapsed range renders as a summary row showing the hidden line count:

````
```go title="calc.go" showLineNumbers collapse={3-8}
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```
````

```go title="calc.go" showLineNumbers collapse={3-8}
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```

→ Lines 3-8 are hidden behind a summary row showing the collapsed line count. Click to expand.

### Multiple ranges

Collapse multiple ranges with comma-separated values. Each range renders its own summary row:

````
```go title="api.go" showLineNumbers collapse={3-7,9-13}
package api

import (
    "encoding/json"
    "net/http"
    "log"
)

type Response struct {
    Status  int
    Message string
    Data    interface{}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    data := fetchData(r.URL.Query())
    json.NewEncoder(w).Encode(Response{Status: 200, Message: "OK", Data: data})
}
```
````

```go title="api.go" showLineNumbers collapse={3-7,9-13}
package api

import (
    "encoding/json"
    "net/http"
    "log"
)

type Response struct {
    Status  int
    Message string
    Data    interface{}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
    data := fetchData(r.URL.Query())
    json.NewEncoder(w).Encode(Response{Status: 200, Message: "OK", Data: data})
}
```

→ The imports (lines 3-7) and type definition (lines 9-13) each collapse independently. Click either summary row to expand that section.

When `PreserveIndent` is `true` (default), the summary text is indented to match the minimum indentation of the collapsed lines. A `--kz-indent` CSS variable is set on the summary, measured in `ch` units.

When line numbers are enabled, the summary row renders an empty gutter placeholder to maintain alignment with surrounding code lines.

## Collapse styles

Set the collapse style per block with `collapseStyle=` or globally via `CollapsibleConfig.Style`. Four styles are available:

| Style | Meta value | Behavior |
|---|---|---|
| GitHub | `collapseStyle="github"` (default) | One-way expand: summary disappears when opened |
| Start | `collapseStyle="collapsible-start"` | Re-collapsible, summary above content |
| End | `collapseStyle="collapsible-end"` | Re-collapsible, summary below content |
| Auto | `collapseStyle="collapsible-auto"` | End if range reaches the last line, Start otherwise |

The GitHub style uses native `<details>` toggling. Once expanded, the summary row is hidden via CSS and there is no UI to re-collapse the section.

The Start and End styles also use `<details>` but keep the summary visible as a toggle. Clicking the summary expands or collapses the content. The End style achieves "summary below content" through CSS `flex-direction: column-reverse`, not DOM reordering.

The Auto style resolves per range. In a block with multiple ranges, one range at the end of the block gets the End style while other ranges get the Start style.

### collapsible-start

The summary appears above the collapsed content. Clicking it toggles the section open and closed:

````
```go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-start"
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```
````

```go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-start"
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```

→ The summary row stays visible after expanding, allowing the section to be collapsed again.

### collapsible-end

The summary appears below the collapsed content:

````
```go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-end"
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```
````

```go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-end"
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```

→ The toggle appears below the hidden lines. Achieved via CSS `flex-direction: column-reverse`, not DOM reordering.

### collapsible-auto

Automatically picks Start or End based on range position. A range that reaches the last line gets End style; other ranges get Start:

````
```go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-auto"
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```
````

```go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-auto"
func Calculate(expr string) (float64, error) {
    tokens := tokenize(expr)
    ast, err := parse(tokens)
    if err != nil {
        return 0, fmt.Errorf("parse: %w", err)
    }
    result := evaluate(ast)
    validate(result)
    return result, nil
}
```

→ Since the range does not reach the last line, Auto resolves to Start style here.

Style cascade: per-block meta `collapseStyle=` > engine config `Style` > `CollapseGithub` default.

## Go API

Control collapse per-block via `Options.Collapse`:

```go
threshold := 30
html, err := engine.Render(code, kazari.Options{
    Lang: "go",
    Collapse: &kazari.CollapseOptions{
        Ranges:    []kazari.Range{{Start: 3, End: 8}},
        Threshold: &threshold,
    },
})
```

**CollapseOptions fields:**

| Field | Type | Description |
|---|---|---|
| `Enabled` | `bool` | Force threshold collapse regardless of line count |
| `Disabled` | `bool` | Suppress threshold collapse |
| `Ranges` | `[]Range` | Line ranges to collapse |
| `Style` | `*CollapseStyle` | Override collapse style (nil = use engine default) |
| `Threshold` | `*int` | Override engine `LineThreshold` for this block |

Or via meta string:

```go
html, err := engine.RenderWithMeta(code, `go collapse={3-8,15-20} collapseStyle="collapsible-start"`)
```

See [Configuration Layers](/getting-started/configuration-layers/) for the full cascade and override rules.

## Configuration

**Meta string options:**

| Option | Syntax | Description |
|---|---|---|
| `collapse` | `collapse` | Force threshold collapse regardless of line count |
| `nocollapse` | `nocollapse` | Suppress threshold collapse for this block |
| `collapse={N-M}` | `collapse={3-8}`, `collapse={3-7,9-13}` | Collapse specific line ranges |
| `collapseStyle=` | `collapseStyle="collapsible-start"` | Override collapse style for ranges |
| `collapseThreshold=N` | `collapseThreshold=20` | Override engine line threshold |

**Config file:**

```yaml
collapsible:
  lineThreshold: 15
  previewLines: 8
  defaultCollapsed: true
  preserveIndent: true
  style: "collapsibleStart"
  expandButtonText: "Show more"
  collapseButtonText: "Show less"
```

Config file style names use camelCase (`collapsibleStart`, `collapsibleEnd`, `collapsibleAuto`, `github`). Meta string uses kebab-case (`collapsible-start`, `collapsible-end`, `collapsible-auto`).

## Range validation

Invalid or conflicting ranges are handled silently:

| Case | Behavior |
|---|---|
| Reversed range (`collapse={8-2}`) | Silently ignored |
| Out of bounds (`collapse={50-60}` on a 20-line block) | Silently ignored |
| Overlapping ranges (`collapse={2-8,5-12}`) | Second range silently ignored |
| Single-line range (`collapse={5-5}`) | Valid, shows "1 collapsed line" |
| Range extends past end (`collapse={18-25}` on a 20-line block) | Clamped to `{18-20}` |

Ranges are sorted by start line before processing. First valid range wins on overlap.

## CSS variables

### Range-based sections

| Variable | Default | Description |
|---|---|---|
| `--kz-collapse-closed-bg` | `rgb(84 174 255 / 20%)` | Summary row background |
| `--kz-collapse-closed-border` | `rgb(84 174 255 / 50%)` | Summary row border color |
| `--kz-collapse-closed-border-width` | `0` | Summary row border width |
| `--kz-collapse-closed-fg` | `currentColor` | Summary text color |
| `--kz-collapse-closed-padding` | `4px` | Summary row vertical padding |
| `--kz-collapse-closed-font-family` | `inherit` | Summary font family |
| `--kz-collapse-closed-font-size` | `inherit` | Summary font size |
| `--kz-collapse-closed-line-height` | `inherit` | Summary line height |
| `--kz-collapse-open-bg` | `transparent` | Expanded background (GitHub style) |
| `--kz-collapse-open-bg-collapsible` | `rgb(84 174 255 / 10%)` | Expanded background (Start/End styles) |
| `--kz-collapse-open-border` | `transparent` | Expanded section border color |
| `--kz-collapse-open-border-width` | `1px` | Expanded section border width |
| `--kz-collapse-expand-icon` | Octicons unfold SVG | Expand icon (inline SVG data URI) |
| `--kz-collapse-collapse-icon` | Octicons fold SVG | Collapse icon (inline SVG data URI) |

### Threshold-based

| Variable | Default | Description |
|---|---|---|
| `--kz-collapse-btn-bg` | `rgba(255,255,255,0.08)` | Button background |
| `--kz-collapse-btn-fg` | `rgba(255,255,255,0.7)` | Button text color |
| `--kz-collapse-btn-hover-bg` | `rgba(255,255,255,0.15)` | Button hover background |
| `--kz-collapse-gradient-start` | `transparent` | Gradient overlay start color |
| `--kz-collapse-gradient-end` | `var(--kz-editor-bg)` | Gradient overlay end color |
| `--kz-collapse-transition` | `300ms ease` | Transition timing for gradient |

These variables are always part of `CSS()`. See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Interaction with other features

- **Line numbers**: summary and gap rows get empty gutter placeholders for alignment. Collapsed lines retain their numbers and resume correctly when expanded.
- **Markers**: marker classes are preserved in hidden DOM elements. When expanded, marked lines display normally with their highlight, insertion, or deletion styles. For threshold mode, marked lines extend the preview as described in [Marker-aware preview](#marker-aware-preview).
- **Focus lines**: included in preview extension, same behavior as markers.
- **Word wrap**: independent. The collapse summary row has its own `--kz-indent` (from the minimum indentation of collapsed lines), separate from wrap's per-line `--kz-indent`.
- **Print**: threshold mode forces all content visible (gradient and button hidden). Range mode shows content and hides summary headers.

## Accessibility

**Range-based sections** use native `<details>/<summary>` elements. Keyboard navigation works out of the box: Tab to focus the summary, Enter or Space to toggle. Screen readers announce the summary text (e.g., "5 collapsed lines"). Expand and collapse icons are decorative CSS pseudo-elements, not announced.

**Threshold-based collapse** uses `aria-expanded` on both the bottom bar button and the toolbar chevron. An `aria-live="polite"` region announces state changes to screen readers. Custom announcement text is configurable via `ExpandedAnnouncement` and `CollapsedAnnouncement` in `CollapsibleConfig`.

## Edge cases

- `WithCollapsible` must be set for threshold-based collapse. Without it, the bare `collapse` meta token has no effect, while range-based `collapse={N-M}` still renders a collapsed section.
- Both modes can coexist: a block can have threshold collapse and range-based collapse simultaneously. Range-based sections take precedence in the render loop for lines that fall within a collapse range.
- `nocollapse` only suppresses threshold-based collapse. Range-based collapse (`collapse={N-M}`) still works on a `nocollapse` block.
- Config file style names use camelCase (`collapsibleStart`), meta string uses kebab-case (`collapsible-start`).
- Unrecognized `collapseStyle` values in meta silently default to `github`. Invalid values in the config file produce a validation error.
- Range-based collapse uses native `<details>/<summary>` toggling and requires no JavaScript. Only threshold-based collapse needs the collapsible JS module.
