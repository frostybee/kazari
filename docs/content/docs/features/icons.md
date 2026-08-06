---
title: "Icons"
description: "File icons in the title bar and language badges with optional icon slots."
tags: [css-variables]
sidebar:
  order: 15
---

Kazari provides two icon systems: file icons in the title bar and language icon slots in the language badge area. Both use CSS-populated HTML slots; no default icon images are included.

## File icons

When a code block has a title with a file extension, Kazari emits a `<span class="kz-file-icon" data-ext="...">` element before the title text:

````
```go title="main.go"
package main

func main() {}
```
````

```go title="main.go"
package main

func main() {}
```

→ The title bar shows a `kz-file-icon` slot before the `main.go` filename. The slot renders as empty space until styled via CSS.

The extension is derived from the last `.` in the title. A title without an extension or with a trailing dot emits no icon span.

File icons appear only in editor frames. Terminal frames and no-frame blocks do not render the icon slot.

To display an image, target the `data-ext` attribute in your stylesheet:

```css
.kz-file-icon[data-ext="go"] {
    background-image: url("/icons/go.svg");
    background-size: contain;
    background-repeat: no-repeat;
}
```

## Inline SVG helper

Use `kazari.CreateInlineSVGURL` to embed an SVG directly as a data URI:

```go
url := kazari.CreateInlineSVGURL(`<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'>
  <path d='M12 2L2 7l10 5 10-5-10-5z'/>
</svg>`)
```

Then use the returned URL in CSS:

```css
.kz-file-icon[data-ext="go"] {
    background-image: url("data:image/svg+xml,...");
}
```

:::tip
SVG attributes inside `CreateInlineSVGURL()` must use single quotes. Double quotes will corrupt the data URL encoding.
:::

## Custom file icon resolver

`WithFileIconResolver` replaces the default empty span with any HTML string:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithFileIconResolver(func(ext string) string {
        icons := map[string]string{
            "go": "🔵",
            "py": "🐍",
            "js": "🟡",
            "rs": "🦀",
        }
        icon, ok := icons[ext]
        if !ok {
            icon = "📄"
        }
        return fmt.Sprintf(`<span class="kz-file-icon">%s</span>`, icon)
    }),
)
```

The resolver receives the file extension without the dot (e.g., `"go"`, `"ts"`). The returned HTML is inserted verbatim before the title span. The resolver is not called when no title is set or when the title has no extension.

`WithFileIconResolver` is a Go API option only and has no config file equivalent.

## Language badge

The language badge displays the language name in the toolbar. `WithLanguageBadge(false)` disables the entire badge area, including both the text label and any icon slot.

Display names are normalized: `javascript` becomes JavaScript, `typescript` becomes TypeScript, `html` becomes HTML, `css` becomes CSS, `json` becomes JSON, `yaml` becomes YAML, and so on. Other languages get their first letter uppercased.

## Language icon mode

By default, the badge shows text only. `WithLanguageIconMode` adds an icon slot:

| Mode | Constant | Output |
|---|---|---|
| Text only (default) | `kazari.LangIconNone` | `<span class="kz-lang">Go</span>` |
| Icon only | `kazari.LangIconOnly` | `<span class="kz-lang-icon" data-lang="go"></span>` |
| Icon and text | `kazari.LangIconAndText` | Icon span before text span |

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithLanguageIconMode(kazari.LangIconAndText),
)
```

The `kz-lang-icon` span has no default image. Style it via `[data-lang]` selectors:

```css
.kz-lang-icon[data-lang="go"] {
    background-image: url("/icons/go.svg");
    background-size: contain;
    background-repeat: no-repeat;
}
```

For theme-adaptive monochrome icons, use `mask-image` instead of `background-image`. The icon inherits the toolbar text color and adapts to light/dark themes automatically:

```css
.kz-lang-icon[data-lang="go"] {
    -webkit-mask-image: url("/icons/go.svg");
    mask-image: url("/icons/go.svg");
    -webkit-mask-size: contain;
    mask-size: contain;
    -webkit-mask-repeat: no-repeat;
    mask-repeat: no-repeat;
    background-color: currentColor;
}
```

When `LanguageBadge` is false, no icon or text is rendered regardless of `LangIconMode`.

## CSS variables

File icon variables (emitted when `WithFileIcons(true)`, the default):

| Variable | Default | Description |
|---|---|---|
| `--kz-file-icon-size` | `1rem` | Icon width and height |
| `--kz-file-icon-margin` | `0` | Icon margin, on top of the toolbar's own gap |
| `--kz-file-icon-opacity` | `0.8` | Icon opacity |

Language icon variables (emitted when `LangIconMode` is not `LangIconNone`):

| Variable | Default | Description |
|---|---|---|
| `--kz-lang-icon-size` | `1.25rem` | Icon width and height |
| `--kz-lang-icon-margin` | `0` | Icon margin |
| `--kz-lang-icon-opacity` | `0.8` | Icon opacity |

Neither set of variables is emitted unless the corresponding feature is active. See the [CSS Variables](/reference/css-variables/) reference for the complete list.

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `WithFileIcons(bool)` | Go API | `true` | Enable or disable file icon slots in the title bar |
| `WithFileIconResolver(func)` | Go API | `nil` | Custom HTML resolver for file icons (Go-only) |
| `WithLanguageBadge(bool)` | Go API | `true` | Show or hide the language badge area |
| `WithLanguageIconMode(mode)` | Go API | `LangIconNone` | Language badge icon mode |
| `fileIcons` | Config file | `true` | YAML/JSON equivalent for `WithFileIcons` |
| `languageBadge` | Config file | `true` | YAML/JSON equivalent for `WithLanguageBadge` |
| `languageIconMode` | Config file | `"none"` | `"none"`, `"iconOnly"`, or `"iconAndText"` |

## Edge cases

- No default icon images are included. Both systems emit empty HTML slots that the consumer styles.
- Both icon types appear only in editor frames. Terminal frames and frameless blocks do not render icon slots.
- `WithLanguageBadge(false)` suppresses both icon and text regardless of `LangIconMode`.
- The `data-lang` attribute uses the raw lowercase language name (e.g., `go`), not the display-cased name (e.g., Go).
- `FileIconResolver` is Go-only. Functions cannot be serialized to YAML/JSON config files.
- A title without a file extension emits no file icon span. The extension is the text after the last `.` in the title.
