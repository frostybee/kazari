---
title: "Frames & Titles"
description: "Editor frames, terminal frames, title bars, and auto-detection from language."
tags: [terminal]
sidebar:
  order: 1
---

Frames wrap a code block in editor or terminal chrome, adding a title bar with a language badge and toolbar buttons. Four frame types are available: `code` (editor), `terminal` (shell window), `none` (bare block), and `auto` (detect from language, the default).

## Frame types

### Editor frame

The default frame for most languages. It renders a title bar with an optional file tab, a language badge, and toolbar buttons.

Set an explicit title with `title=`:

````
```go title="main.go"
func main() {
    fmt.Println("Hello, world!")
}
```
````

```go title="main.go"
func main() {
    fmt.Println("Hello, world!")
}
```

→ An editor frame with a `main.go` file tab, language badge, and toolbar buttons.

### Terminal frame

Use a shell language to trigger terminal auto-detection:

````
```bash
curl -s https://api.example.com/status | jq '.health'
```
````

```bash
curl -s https://api.example.com/status | jq '.health'
```

→ A terminal frame with macOS-style colored dots and toolbar buttons.

Add an explicit title to a terminal frame:

````
```powershell title="PowerShell"
Get-Process | Where-Object { $_.CPU -gt 100 }
```
````

```powershell title="PowerShell"
Get-Process | Where-Object { $_.CPU -gt 100 }
```

→ A terminal frame with "PowerShell" centered in the title bar.

The following 19 languages trigger terminal auto-detection:

`ansi` `bash` `bat` `batch` `cmd` `console` `fish` `nu` `nushell` `powershell` `ps` `ps1` `psd1` `psm1` `sh` `shell` `shellscript` `shellsession` `zsh`

Matching is case-insensitive.

### No frame

Remove the frame entirely with `frame="none"`. No title bar, no toolbar. Only the copy button remains, floating in the top-right corner on hover.

````
```go frame="none"
fmt.Println("no frame")
```
````

```go frame="none"
fmt.Println("no frame")
```

→ A bare code block with no title bar or toolbar.

### Forcing a frame type

An explicit `frame=` in the meta string always wins, even if it contradicts the language.

````
```go frame=terminal title="Go in a terminal"
fmt.Println("terminal frame on a Go block")
```
````

```go frame=terminal title="Go in a terminal"
fmt.Println("terminal frame on a Go block")
```

→ A Go block rendered with terminal chrome because `frame=terminal` overrides auto-detection.

## Title bar

### Explicit title

Set a title with `title=` in the meta string. The title appears as a file tab in editor frames and as centered text in terminal frames.

````
```typescript title="server.ts"
const app = express()
app.listen(3000)
```
````

```typescript title="server.ts"
const app = express()
app.listen(3000)
```

→ An editor frame with a `server.ts` file tab in the title bar.

### File name extraction

Place a file path comment on the first line of a block with no `title=` set:

`````title="Markdown source"
````
```javascript
// src/greet.js
export function greet(name) {
    return `Hello, ${name}!`
}
```
````
`````

```javascript
// src/greet.js
export function greet(name) {
    return `Hello, ${name}!`
}
```

→ An editor frame with `src/greet.js` as the title. The comment line is removed from the rendered output.

Kazari scans the first 4 lines for a comment containing a file path. Recognized comment styles:

| Style | Example |
|---|---|
| Line comment (`//`) | `// main.go` |
| Hash comment (`#`) | `# deploy.sh` |
| HTML comment | `<!-- index.html -->` |
| Block comment (`/* */`) | `/* styles.css */` |

Shebangs (`#!/bin/bash`) are not treated as filename comments. Optional label prefixes (`file name:`, `filename:`, `example:`) are stripped automatically.

A comment is accepted as a filename when it has a recognized file extension (~70 supported), matches a special basename (`makefile`, `dockerfile`, `rakefile`, `gemfile`, `procfile`), or is a dotfile (`.env`, `.gitignore`). Comments containing `://` (URLs) or spaces (sentences) are rejected.

An explicit `title=` in the meta string always takes precedence. File name extraction only runs when no title is set.

Disable this behavior with `kazari.WithFileNameExtraction(false)` or `fileNameExtraction: false` in the config file.

### Shebang override

A `#!` shebang in the first 4 lines forces the frame to editor mode, even for terminal languages:

````
```bash
#!/bin/bash
set -euo pipefail
echo "script file, not a terminal session"
```
````

```bash
#!/bin/bash
set -euo pipefail
echo "script file, not a terminal session"
```

→ Despite being `bash`, this block gets an editor frame because the shebang indicates a script file, not a terminal session.

### Pipeline order

File name extraction runs before frame auto-detection. A bash block with a filename comment (like `# deploy.sh`) first gets its title extracted, then the file indicator check assigns an editor frame instead of a terminal frame.

## Terminal dot styles

Two dot styles are available for terminal frames, set at the engine level with `kazari.WithTerminalDotStyle()`:

| Style | Constant | Description |
|---|---|---|
| Colored | `kazari.DotsColored` | Three macOS-style red, yellow, and green dots (default) |
| Minimal | `kazari.DotsMinimal` | Monochrome CSS-only dots via SVG mask |

There is no meta string option for dot style. It applies to all terminal frames in the engine.

```yaml
terminalDotStyle: "minimal"
```

Both CSS rule sets are always included in the generated stylesheet regardless of which style is configured. Only the rendered HTML markup differs.

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `frame=` | Meta string | `"auto"` | Set to `"code"`, `"terminal"`, `"none"`, or `"auto"` |
| `title=` | Meta string | none | Title shown in the title bar |
| `WithFrameDetection(bool)` | Go API | `true` | Enable or disable language-based auto-detection |
| `WithFileNameExtraction(bool)` | Go API | `true` | Extract title from code comments |
| `WithTerminalDotStyle(style)` | Go API | `DotsColored` | Terminal dot appearance |
| `WithTerminalCommentStripping(bool)` | Go API | `true` | Strip `#` comments from terminal copy payload |

Set frames per language via `WithLanguageDefaults` or the config file:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithLanguageDefaults(map[string]kazari.BlockDefaults{
        "bash,sh,zsh": {Frame: kazari.FrameTerminal},
        "go":          {Frame: kazari.FrameCode},
    }),
)
```

```yaml
frameDetection: true
fileNameExtraction: true

languageDefaults:
  "bash,sh,zsh":
    frame: "terminal"
```

## CSS variables

| Variable | Description |
|---|---|
| `--kz-radius` | Border radius of the frame |
| `--kz-border` | Frame border |
| `--kz-shadow` | Frame box shadow |
| `--kz-toolbar-bg` | Toolbar background (theme-derived) |
| `--kz-toolbar-border` | Toolbar bottom border (theme-derived) |
| `--kz-toolbar-padding` | Toolbar padding |
| `--kz-terminal-bg` | Terminal frame background |
| `--kz-terminal-titlebar-bg` | Terminal title bar background |
| `--kz-terminal-header-padding` | Terminal header padding |
| `--kz-terminal-dot-red` | Red dot color (colored style) |
| `--kz-terminal-dot-yellow` | Yellow dot color (colored style) |
| `--kz-terminal-dot-green` | Green dot color (colored style) |
| `--kz-terminal-dots-fg` | Dot color (minimal style, theme-derived) |
| `--kz-terminal-dots-opacity` | Dot opacity (minimal style) |

See the [CSS Variables](/reference/css-variables/) reference for the complete list with default values.

## Edge cases

- An explicit `frame=` in the meta string always wins, even if it contradicts the language.
- The `ansi` language auto-detects as a terminal frame.
- Mermaid blocks (`lang="mermaid"`) bypass framing entirely and render as a bare `<pre class="mermaid">`.
- Terminal frames never show file icons, even when `WithFileIcons(true)` is enabled and a title is set.
- Terminal comment stripping only applies to terminal frames and only affects the copy-button payload, not the displayed code. See [Terminal Comment Stripping](/features/terminal-comment-stripping/) for details.
- When `WithFrameDetection(false)`, all `FrameAuto` blocks default to `FrameCode`. Explicit `frame=terminal` still works.
