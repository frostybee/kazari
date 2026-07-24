---
title: "Line Markers"
description: "Highlight, insert, and delete lines with colored backgrounds and optional labeled ranges."
sidebar:
  order: 4
---

Line markers highlight entire lines with a colored background and a left-edge accent border. Three marker types are available: mark (neutral highlight), ins (insertion), and del (deletion).

## Mark lines

Add line numbers or ranges in curly braces after the language:

````
```go {2, 4-6}
func main() {
    name := "world"
    greeting := greet(name)
    fmt.Println(greeting)
    fmt.Println("done")
    fmt.Println("bye")
}
```
````

```go {2, 4-6}
func main() {
    name := "world"
    greeting := greet(name)
    fmt.Println(greeting)
    fmt.Println("done")
    fmt.Println("bye")
}
```

→ Lines 2, 4, 5, and 6 have a yellow background with a left accent border.

Lines are 1-based and inclusive. Comma-separated values accept single numbers and ranges. A range like `4-6` covers lines 4, 5, and 6.

## Insertion and deletion markers

Use `ins=` and `del=` to mark lines as added or removed. Insertion lines show a green background with a `+` indicator. Deletion lines show red with a `-` indicator.

````
```go title="diff.go" showLineNumbers {3} del={5-7} ins={9-11}
func process(items []string) ([]string, error) {
    var results []string
    for _, item := range items {
        val := item
        result, err := transform(val)
        if err != nil {
            return nil, err
        }
        cleaned := strings.TrimSpace(item)
        result := convert(cleaned)
        results = append(results, result)
    }
    return results, nil
}
```
````

```go title="diff.go" showLineNumbers {3} del={5-7} ins={9-11}
func process(items []string) ([]string, error) {
    var results []string
    for _, item := range items {
        val := item
        result, err := transform(val)
        if err != nil {
            return nil, err
        }
        cleaned := strings.TrimSpace(item)
        result := convert(cleaned)
        results = append(results, result)
    }
    return results, nil
}
```

→ Line 3 is highlighted (mark), lines 5-7 show red deletion markers, and lines 9-11 show green insertion markers.

The `add=` and `rem=` aliases work identically to `ins=` and `del=`.

## Labeled ranges

Add a label to any marker by placing a quoted string before the colon inside the braces. The label appears as a pill badge on the first line of each range.

````
```go showLineNumbers {"Config":2-4} ins={"Added":7-9}
func NewServer(opts ...Option) *Server {
    cfg := defaultConfig()
    cfg.apply(opts)
    cfg.validate()

    return &Server{
        router: chi.NewRouter(),
        logger: cfg.logger,
        port:   cfg.port,
    }
}
```
````

```go showLineNumbers {"Config":2-4} ins={"Added":7-9}
func NewServer(opts ...Option) *Server {
    cfg := defaultConfig()
    cfg.apply(opts)
    cfg.validate()

    return &Server{
        router: chi.NewRouter(),
        logger: cfg.logger,
        port:   cfg.port,
    }
}
```

→ Lines 2-4 have a "Config" badge, lines 7-9 have an "Added" badge with insertion styling.

Both double and single quotes work for labels. On labeled `ins`/`del` lines, the label badge replaces the `+`/`-` diff indicator.

Compact numeric labels work well for step-by-step annotations:

````
```go title="handler.go" showLineNumbers {"1":2-3} del={"2":5-6} ins={"3":8-9}
func handle(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")

    data, err := fetchByID(id)
    if err != nil { http.Error(w, err.Error(), 500); return }

    result, err := svc.Process(ctx, id)
    if err != nil { http.Error(w, "processing failed", 500); return }

    json.NewEncoder(w).Encode(result)
}
```
````

```go title="handler.go" showLineNumbers {"1":2-3} del={"2":5-6} ins={"3":8-9}
func handle(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")

    data, err := fetchByID(id)
    if err != nil { http.Error(w, err.Error(), 500); return }

    result, err := svc.Process(ctx, id)
    if err != nil { http.Error(w, "processing failed", 500); return }

    json.NewEncoder(w).Encode(result)
}
```

→ Three labeled ranges show a step-by-step code evolution: step 1 (mark), step 2 (deleted), step 3 (inserted).

## Overlap priority

When multiple marker types overlap on the same line, the highest-priority marker wins:

| Priority | Type | Color |
|---|---|---|
| Lowest | `mark` | Yellow / neutral |
| Middle | `del` | Red |
| Highest | `ins` | Green |

For example, `{1-10} ins={5}` marks lines 1-4 and 6-10 as `mark`, but line 5 becomes `ins`. Labels from a lower-priority marker are cleared when a higher-priority marker overwrites the line.

## Go API

Use the `Options.LineMarkers` field to set markers programmatically:

```go
html, err := engine.Render(code, kazari.Options{
    Lang: "go",
    LineMarkers: []kazari.LineMarker{
        {Type: kazari.MarkerMark, Lines: []kazari.Range{{Start: 3, End: 5}}},
        {Type: kazari.MarkerIns, Lines: []kazari.Range{{Start: 8, End: 10}}, Label: "New"},
        {Type: kazari.MarkerDel, Lines: []kazari.Range{{Start: 6, End: 6}}},
    },
})
```

## Syntax reference

| Syntax | Description |
|---|---|
| `{N}` or `{N-M}` | Mark lines (neutral highlight) |
| `{N,M-P}` | Comma-separated numbers and ranges |
| `ins={N-M}` | Mark lines as inserted |
| `del={N-M}` | Mark lines as deleted |
| `add={N-M}` | Alias for `ins=` |
| `rem={N-M}` | Alias for `del=` |
| `{"Label":N-M}` | Labeled marker range |
| `ins={"Label":N-M}` | Labeled insertion range |
| `del={"Label":N-M}` | Labeled deletion range |

## CSS variables

| Variable | Description |
|---|---|
| `--kz-mark-bg` | Mark line background |
| `--kz-mark-border` | Mark line left accent color |
| `--kz-mark-border-width` | Accent border width |
| `--kz-ins-bg` | Insertion line background |
| `--kz-ins-border` | Insertion line left accent color |
| `--kz-ins-indicator` | Insertion gutter indicator (default `"+"`) |
| `--kz-ins-indicator-color` | Insertion indicator color |
| `--kz-del-bg` | Deletion line background |
| `--kz-del-border` | Deletion line left accent color |
| `--kz-del-indicator` | Deletion gutter indicator (default `"-"`) |
| `--kz-del-indicator-color` | Deletion indicator color |
| `--kz-diff-indicator-margin` | Diff indicator left margin |
| `--kz-label-fg` | Label badge text color |
| `--kz-label-font-size` | Label badge font size |
| `--kz-label-radius` | Label badge border radius |
| `--kz-label-padding` | Label badge padding |
| `--kz-label-mark-bg` | Labeled mark line background (stronger than unlabeled) |
| `--kz-label-ins-bg` | Labeled insertion line background |
| `--kz-label-del-bg` | Labeled deletion line background |

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- The left accent border uses `box-shadow: inset` rather than CSS `border-left` to avoid layout shift.
- Labeled lines have a stronger background than unlabeled lines of the same type (`--kz-label-mark-bg` vs `--kz-mark-bg`).
- On highlighted lines, line numbers become more prominent (brighter color, higher opacity). See [Line Numbers](/features/line-numbers/) for details.
- Line markers are purely visual. They do not affect the copied code.
