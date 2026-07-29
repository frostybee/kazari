---
title: "Focus Lines"
description: "Dim all lines except the ones that matter."
tags: [line-annotations]
sidebar:
  order: 6
---

Focus lines dim everything except the specified lines, drawing attention to the code that matters. Non-focused lines fade to reduced opacity while focused lines remain fully visible.

## Focus a range

Add `focus={N-M}` to the meta string:

````
```go title="process.go" showLineNumbers focus={3-5}
func process(items []string) error {
    for _, item := range items {
        if err := validate(item); err != nil {
            return fmt.Errorf("invalid: %w", err)
        }
    }
    return nil
}
```
````

```go title="process.go" showLineNumbers focus={3-5}
func process(items []string) error {
    for _, item := range items {
        if err := validate(item); err != nil {
            return fmt.Errorf("invalid: %w", err)
        }
    }
    return nil
}
```

→ Lines 3-5 render at full opacity. All other lines fade to reduced opacity.

The syntax accepts comma-separated line numbers and ranges, the same format as [line markers](/features/line-markers/): `focus={1-3,7,10-12}`.

## Multiple ranges

Focus non-contiguous sections of code:

````
```go title="server.go" showLineNumbers focus={2-3,7-8}
func main() {
    cfg := loadConfig()
    db := connectDB(cfg)
    defer db.Close()

    router := chi.NewRouter()
    router.Get("/health", healthHandler)
    router.Get("/api/users", usersHandler(db))

    http.ListenAndServe(":8080", router)
}
```
````

```go title="server.go" showLineNumbers focus={2-3,7-8}
func main() {
    cfg := loadConfig()
    db := connectDB(cfg)
    defer db.Close()

    router := chi.NewRouter()
    router.Get("/health", healthHandler)
    router.Get("/api/users", usersHandler(db))

    http.ListenAndServe(":8080", router)
}
```

→ Lines 2-3 and 7-8 render at full opacity. All other lines are dimmed.

## Combined with markers

Focus and line markers are independent. Both can appear on the same block. Non-focused lines are always dimmed, including lines with markers.

A marked line in the focus set retains full opacity. A marked line outside the focus set is dimmed along with all other non-focused lines.

````
```go title="combined.go" showLineNumbers {4-5} ins={10-12} del={6-8} focus={4-5,10-12}
func main() {
    ctx := context.Background()

    db, err := sql.Open("postgres", dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()
    cache := newCache()
    cache.Warm(ctx)

    svc := newService(db, cache)
    srv := newServer(svc, db)
    srv.Run(ctx)
}
```
````

```go title="combined.go" showLineNumbers {4-5} ins={10-12} del={6-8} focus={4-5,10-12}
func main() {
    ctx := context.Background()

    db, err := sql.Open("postgres", dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()
    cache := newCache()
    cache.Warm(ctx)

    svc := newService(db, cache)
    srv := newServer(svc, db)
    srv.Run(ctx)
}
```

→ Lines 4-5 show highlighted backgrounds at full opacity. Lines 6-8 have deletion markers but are dimmed (outside the focus set). Lines 10-12 have insertion markers at full opacity.

## Configuration

| Option | Syntax | Description |
|---|---|---|
| `focus={N-M}` | `focus={3-5}`, `focus={1-3,7}` | Focus the specified lines |

Go API equivalent via the `Options.FocusLines` field:

```go
html, err := engine.Render(code, kazari.Options{
    Lang:       "go",
    FocusLines: []kazari.Range{{Start: 3, End: 5}},
})
```

There is no engine-level default for focus. It is always a per-block setting.

## CSS variables

| Variable | Default | Description |
|---|---|---|
| `--kz-focus-dimmed-opacity` | `0.35` | Opacity of non-focused lines |

Lower values create more contrast between focused and unfocused lines.

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- Focus is purely visual. It does not affect the copied code.
- Focus can be combined with inline markers. Inline highlights on dimmed lines are dimmed along with the rest of the line.
- Line numbers on dimmed lines inherit the reduced opacity.
