---
title: "Links"
description: "Add clickable hyperlinks inside code blocks using annotation syntax."
tags: [line-annotations]
sidebar:
  order: 14
---

The `@[text](url)` annotation syntax embeds clickable hyperlinks inside code blocks while preserving token colors and styles. Links are disabled by default; enable them with `kazari.WithLinks(true)`.

## Annotation syntax

Use `@[display text](url)` anywhere in the code. The `@` prefix is required. Plain `[text](url)` is not matched.

```go title="main.go"
import (
    @[fmt](https://pkg.go.dev/fmt)
    @[os](https://pkg.go.dev/os)
)
```

→ The display text (`fmt`, `os`) appears as a clickable hyperlink with a dotted underline. Links inherit the token color.

Go API equivalent:

```go
html, err := engine.RenderWithMeta(code, `go title="main.go"`)
```

## Allowed URL schemes

| URL form | Allowed |
|---|---|
| `https://example.com` | Yes |
| `http://example.com` | Yes |
| `mailto:user@example.com` | Yes |
| `/api/reference` (absolute path) | Yes |
| `//example.com` (protocol-relative) | No |
| `javascript:...` | No |
| `data:...` | No |
| `ftp://...` or other schemes | No |

:::caution
Invalid URLs cause the `@[text](url)` annotation to appear as literal text in the rendered output. No error is raised.
:::

## How it works

Link processing runs inside the preprocessing pipeline before tokenization:

1. Annotation syntax is extracted from each line. The `@[`, `](url)` wrappers are stripped, leaving only the display text.
2. The cleaned code is passed to the highlighter. Byte offsets track where each link text sits in the cleaned output.
3. During HTML rendering, those offsets map back to token spans. The matching tokens are wrapped in `<a class="kz-link" href="..." target="_blank" rel="noopener noreferrer">` elements with an external link icon appended.

The copy button copies the cleaned code, not the original annotation syntax. Link URLs are excluded from the copy payload.

## Links + inline markers

Links compose with inline markers on the same line. Marked text inside a link gets both the marker element and the anchor wrapper:

```jsx title="index.tsx" "createRoot" ins="StrictMode"
const root = @[createRoot](https://react.dev/reference/react-dom/client/createRoot)(
    document.getElementById("root")
)
root.render(<@[StrictMode](https://react.dev/reference/react/StrictMode)><App /></@[StrictMode](https://react.dev/reference/react/StrictMode)>)
```

→ The `"createRoot"` inline marker highlights all occurrences with the default mark style, and `ins="StrictMode"` marks those occurrences as inserted (green). Both markers apply on top of the link underline. The rendered HTML nests `<a class="kz-link">` around the `<mark>` or `<ins>` element.

## CSS variables

| Variable | Default | Description |
|---|---|---|
| `--kz-link-underline` | `currentColor` | Underline color |
| `--kz-link-underline-width` | `1px` | Underline thickness |
| `--kz-link-underline-offset` | `0.15em` | Underline vertical offset |

The link CSS is only included in `engine.CSS()` output when links are enabled. Override any variable on `:root` or a more specific selector:

```css
:root {
    --kz-link-underline: #60a5fa;
    --kz-link-underline-offset: 0.2em;
}
```

See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `WithLinks(bool)` | Go API | `false` | Enable link extraction and rendering |
| `links` | Config file | `false` | YAML/JSON equivalent |

Links are an engine-wide setting. There is no per-block meta string override.

## Edge cases

- Links are disabled by default. Enable with `WithLinks(true)` or `links: true` in the config file.
- Invalid URLs cause the `@[text](url)` annotation to appear as literal text. No error is raised and no link is rendered.
- The copy button copies the cleaned code without annotation syntax or URLs.
- When linked text spans multiple syntax tokens, each token gets its own `<a>` element pointing to the same URL. The external link icon appears only on the last segment.
- Mermaid blocks bypass link extraction entirely (pass-through returns before preprocessing runs).
- The underline style is dotted by default and becomes solid on hover.
