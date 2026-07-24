---
title: "Post-Render Callbacks"
description: "Modify rendered HTML with WithPostRender callbacks that receive block metadata."
sidebar:
  order: 1
---

Post-render callbacks run after each code block is rendered, receiving the HTML string and a `BlockInfo` struct with metadata about the block. They enable HTML modifications that built-in features and CSS cannot achieve: injecting new elements, wrapping blocks in custom containers, or adding attributes based on block context.

## Basic callback

Register a callback with `WithPostRender`. This example wraps every titled block in a `<figure>` with a `<figcaption>`:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithPostRender(func(html string, info kazari.BlockInfo) string {
        if info.Title == "" {
            return html
        }
        return fmt.Sprintf(
            "<figure>%s<figcaption>%s</figcaption></figure>",
            html, info.Title,
        )
    }),
)
```

→ Blocks with a title render inside a `<figure>` element. Blocks without a title pass through unchanged.

The callback receives the complete HTML of the rendered `.kazari-block` div. Return the modified string. Return the input unchanged to skip modification.

## TODO/FIXME badges

Inject badge elements next to comment annotations. This is distinct from regex markers, which style existing text but cannot insert new HTML:

```go
todoRe := regexp.MustCompile(`(>[^<]*?)(TODO:)`)
fixmeRe := regexp.MustCompile(`(>[^<]*?)(FIXME:)`)

engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithPostRender(func(html string, info kazari.BlockInfo) string {
        html = todoRe.ReplaceAllString(html,
            `${1}<span class="kz-todo-badge">TODO</span> `)
        html = fixmeRe.ReplaceAllString(html,
            `${1}<span class="kz-fixme-badge">FIXME</span> `)
        return html
    }),
)
```

→ `TODO:` and `FIXME:` in comments are replaced with styled badge spans. The regex pattern `(>[^<]*?)` targets only visible text content, not attribute values like `data-code`.

Style the badges with CSS in the consumer's stylesheet:

```css
.kz-todo-badge {
    display: inline-block;
    font-size: 0.7em;
    font-weight: 700;
    padding: 0.1em 0.45em;
    border-radius: 3px;
    background: rgba(234, 179, 8, 0.2);
    color: #b45309;
}
.kz-fixme-badge {
    display: inline-block;
    font-size: 0.7em;
    font-weight: 700;
    padding: 0.1em 0.45em;
    border-radius: 3px;
    background: rgba(239, 68, 68, 0.2);
    color: #dc2626;
}
```

## BlockInfo fields

Every callback receives a `BlockInfo` struct with metadata about the rendered block:

| Field | Type | Description |
|---|---|---|
| `Lang` | `string` | Resolved canonical language name |
| `Title` | `string` | Block title (after filename extraction) |
| `Frame` | `Frame` | Frame type: `FrameCode`, `FrameTerminal`, or `FrameNone` |
| `RawCode` | `string` | Preprocessed source code (what the copy button copies) |
| `LineCount` | `int` | Number of rendered lines |
| `Theme` | `string` | Per-block theme override (empty when using engine defaults) |
| `Meta` | `string` | Raw meta string from `RenderWithMeta`; empty from `Render` |

Use `BlockInfo` fields to conditionally apply modifications:

```go
kazari.WithPostRender(func(html string, info kazari.BlockInfo) string {
    if info.Lang != "go" {
        return html
    }
    return html + `<a class="playground-link" href="#">Run in Playground</a>`
})
```

→ A playground link appears below Go blocks only. Other languages pass through unchanged.

## Multiple callbacks

Call `WithPostRender` multiple times to register multiple callbacks. They run in registration order, each receiving the output of the previous one:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithPostRender(addFigCaption),
    kazari.WithPostRender(addTodoBadges),
)
```

→ `addFigCaption` runs first, then `addTodoBadges` receives the already-wrapped HTML. The same `BlockInfo` is passed to every callback in the chain.

## Configuration

| Option | Layer | Description |
|---|---|---|
| `WithPostRender(func)` | Go API | Register a post-render callback. Multiple calls append; callbacks run in order. |

Post-render callbacks are Go API only. They cannot be configured via meta strings or config files.

Mermaid pass-through blocks bypass post-render callbacks. They return a bare `<pre class="mermaid">` element before the post-render step runs.

Post-render callbacks execute at the very end of the [render pipeline](/architecture/render-pipeline/), after `RenderBlock` assembles the full HTML string.

## Related

- [Custom Highlighters](/plugins/custom-highlighters/) for swapping tokenization engines or accessing raw tokens
- [Client-Side Extensibility](/plugins/client-side-extensibility/) for extending blocks with CSS and JavaScript instead of Go callbacks
