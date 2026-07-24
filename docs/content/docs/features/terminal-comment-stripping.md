---
title: "Terminal Comment Stripping"
description: "Remove shell comments from the copy payload while keeping them visible in the code."
sidebar:
  order: 17
---

Terminal comment stripping removes full-line `#` comments from the copy-button payload while keeping them visible in the displayed code. Readers see the explanatory comments; pasting produces only the runnable commands.

## How it works

Lines where the first non-whitespace character is `#` are removed from the copy payload. The displayed code is completely unaffected.

````
```bash title="setup"
# Install dependencies
npm install
# Build the project
npm run build
# Start the server
npm start
```
````

```bash title="setup"
# Install dependencies
npm install
# Build the project
npm run build
# Start the server
npm start
```

→ The terminal block displays all six lines. Clicking the copy button copies only the three commands: `npm install`, `npm run build`, `npm start`. The `#` comment lines are excluded from the clipboard.

Stripping rules:

- A line is stripped if its first non-whitespace character is `#`. Indented comments like `  # step 2` are also stripped.
- Inline `#` after content (e.g., `echo hello # note`) is preserved because the line starts with `e`, not `#`.
- Blank lines between commands are preserved.

## When it applies

Terminal comment stripping is gated on the block's resolved frame type, not the language. It applies to any block that resolves to `FrameTerminal`:

- Terminal languages (`bash`, `sh`, `zsh`, `powershell`, `fish`, `console`, etc.) auto-detect as terminal frame when `WithFrameDetection` is active.
- `ansi` blocks also auto-detect as terminal frame.
- Explicit `frame="terminal"` on any language triggers stripping.
- Editor frames and frameless blocks are unaffected regardless of language.

## Disabling

To keep comments in the copy payload, disable stripping at engine construction:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithTerminalCommentStripping(false),
)
```

Or in the config file:

```yaml
terminalCommentStripping: false
```

## Configuration

| Option | Layer | Default | Description |
|---|---|---|---|
| `WithTerminalCommentStripping(bool)` | Go API | `true` | Enable or disable comment stripping from the copy payload |
| `terminalCommentStripping` | Config file | `true` | YAML/JSON equivalent |

Stripping is enabled by default. There is no per-block meta string override. It is an engine-wide setting.

## Edge cases

- Enabled by default. Only affects the copy payload (`data-code` attribute). The displayed code is unchanged.
- Inline `#` after content (e.g., `echo hello # note`) is preserved. Only full-line comments are stripped.
- Indented comments (e.g., `  # step 2`) are also stripped.
- There is no `#!` shebang exception in the stripping logic. However, shebang lines force auto-detection to editor frame, so stripping does not trigger under normal usage. Explicit `frame="terminal"` on code containing `#!/bin/bash` would strip the shebang line from the copy payload.
- Stripping also affects `BlockInfo.RawCode` passed to post-render callbacks registered via `WithPostRender`.
