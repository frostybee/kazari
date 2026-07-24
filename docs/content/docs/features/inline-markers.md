---
title: "Inline Markers"
description: "Mark text and regex matches within lines using highlight, insertion, and deletion styles."
sidebar:
  order: 5
---

Inline markers highlight specific text within code lines. They match by literal string or regex pattern and render as `<mark>`, `<ins>`, or `<del>` elements with a colored background pill.

## Text markers

Place a quoted string in the meta string to highlight all occurrences of that text across every line:

````
```typescript title="cache.ts" showLineNumbers "CacheEntry" ins="factory"
interface CacheEntry {
    key: string
    value: unknown
    ttl: number
}

function factory(key: string): CacheEntry {
    return { key, value: null, ttl: 300 }
}
```
````

```typescript title="cache.ts" showLineNumbers "CacheEntry" ins="factory"
interface CacheEntry {
    key: string
    value: unknown
    ttl: number
}

function factory(key: string): CacheEntry {
    return { key, value: null, ttl: 300 }
}
```

→ Every occurrence of `CacheEntry` shows a neutral highlight pill. `factory` shows a green insertion pill.

Bare quoted strings produce the default `mark` type (neutral highlight). Prefix with `ins=` or `del=` for typed markers. Both double quotes (`"text"`) and single quotes (`'text'`) are accepted.

## Regex markers

Use forward slashes to match a regular expression pattern:

````
```go title="regex-markers.go" showLineNumbers /err\b/ ins=/func\s+\w+/ del=/fmt\.Errorf/
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }
    return &cfg, nil
}
```
````

```go title="regex-markers.go" showLineNumbers /err\b/ ins=/func\s+\w+/ del=/fmt\.Errorf/
func loadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }
    return &cfg, nil
}
```

→ Matches of `/err\b/` show neutral highlight pills. `/func\s+\w+/` shows green insertion pills. `/fmt\.Errorf/` shows red deletion pills.

Bare `/pattern/` produces the `mark` type. Prefix with `ins=` or `del=` for typed matches. Escape forward slashes within the pattern with `\/`. Invalid regex patterns are silently ignored.

## Capture groups

Use `(...)` to highlight only the captured portion of a regex match. Use `(?:...)` for a non-capturing group that highlights the full match. Only the first capturing group is used; additional groups are ignored.

````
```python title="capture_group.py" /ye(s|p)/
result = "yes" if check() else "yep"
print("nope")
```
````

```python title="capture_group.py" /ye(s|p)/
result = "yes" if check() else "yep"
print("nope")
```

→ The regex `/ye(s|p)/` matches `yes` and `yep`, but only `s` and `p` are highlighted (the captured groups). `nope` is not matched.

## Multi-token spanning

When a match crosses syntax token boundaries, Kazari splits the highlight across the affected tokens and joins them visually into a single continuous pill. This works automatically using `open-start` and `open-end` CSS classes that remove borders and border radius at the join point.

## Overlap priority

Inline markers share the same priority system as [line markers](/features/line-markers/): `ins` > `del` > `mark`. When two inline markers overlap on the same characters, the higher-priority marker claims its range. The lower-priority match survives only in the non-overlapping fragments.

For example, `"abcdef" ins="cd"` produces three segments: `ab` (mark), `cd` (ins), `ef` (mark).

## Combined with line markers

Inline markers and line markers are independent. A line can have both a line-level highlight (colored background) and inline text highlights within it:

````
```go title="combined.go" showLineNumbers {4-5} ins={10-12} del={6-8} "db"
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

```go title="combined.go" showLineNumbers {4-5} ins={10-12} del={6-8} "db"
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

→ Lines 4-5 have a yellow mark background with inline `db` highlight pills. Lines 6-8 have red deletion backgrounds. Lines 10-12 have green insertion backgrounds. The `db` text is highlighted on every line where it appears.

## Go API

Use the `Options.InlineMarkers` field to set markers programmatically:

```go
html, err := engine.Render(code, kazari.Options{
    Lang: "go",
    InlineMarkers: []kazari.InlineMarker{
        {Type: kazari.MarkerMark, Text: "CacheEntry"},
        {Type: kazari.MarkerIns, Text: `func\s+\w+`, IsRegex: true},
    },
})
```

## Syntax reference

| Syntax | Description |
|---|---|
| `"text"` | Mark all occurrences (neutral highlight) |
| `'text'` | Same as double quotes |
| `ins="text"` | Mark as inserted |
| `del="text"` | Mark as deleted |
| `/regex/` | Regex match (neutral highlight) |
| `ins=/regex/` | Regex match as inserted |
| `del=/regex/` | Regex match as deleted |
| `add="text"` | Alias for `ins="text"` |
| `rem="text"` | Alias for `del="text"` |
| `add=/regex/` | Alias for `ins=/regex/` |
| `rem=/regex/` | Alias for `del=/regex/` |

## CSS variables

| Variable | Default | Description |
|---|---|---|
| `--kz-inline-mark-bg` | `rgba(255,200,0,0.2)` | Mark highlight background |
| `--kz-inline-mark-border` | `rgba(255,200,0,0.5)` | Mark highlight border color |
| `--kz-inline-mark-border-width` | `1.5px` | Border width |
| `--kz-inline-mark-radius` | `0.2rem` | Border radius (pill shape) |
| `--kz-inline-mark-padding` | `0.15rem` | Padding inside the highlight |
| `--kz-inline-ins-bg` | `rgba(46,160,67,0.2)` | Insertion highlight background |
| `--kz-inline-ins-border` | `rgba(46,160,67,0.5)` | Insertion highlight border |
| `--kz-inline-del-bg` | `rgba(248,81,73,0.2)` | Deletion highlight background |
| `--kz-inline-del-border` | `rgba(248,81,73,0.5)` | Deletion highlight border |

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Edge cases

- Matching is case-sensitive. `"useState"` matches `useState` but not `usestate`.
- All occurrences across all lines are matched, not just the first.
- On lines with a line marker background, all token colors on that line (including text inside inline markers) are contrast-adjusted against the line marker background, not the editor background. This is a per-line effect driven by line-level markers.
- Inline markers are purely visual. They do not affect the copied code.
