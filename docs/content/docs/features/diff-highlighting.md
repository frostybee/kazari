---
title: "Diff Highlighting"
description: "Unified diff with full syntax highlighting in the original language."
tags: [line-annotations, highlighter]
sidebar:
  order: 11
---

Diff highlighting combines unified diff markers with full syntax highlighting. Lines starting with `+` and `-` are marked as insertions and deletions while the code is highlighted in the original language rather than plain diff syntax.

## Hybrid diff syntax

Set the language to `diff` and specify the original language with `lang=`:

````
```diff lang="go" title="hybrid-diff.go" showLineNumbers
 import (
-    "fmt"
+    "log"
     "os"
 )
 
 func main() {
-    fmt.Println("hello")
+    log.Println("hello")
 }
```
````

```diff lang="go" title="hybrid-diff.go" showLineNumbers
 import (
-    "fmt"
+    "log"
     "os"
 )
 
 func main() {
-    fmt.Println("hello")
+    log.Println("hello")
 }
```

→ Lines prefixed with `+` show green insertion markers. Lines prefixed with `-` show red deletion markers. Context lines and the code are highlighted with Go syntax colors.

Any language supported by the configured highlighter can be used with `lang=`:

````
```diff lang="javascript" title="config.js"
 const config = {
-  port: 3000,
+  port: 8080,
   host: "localhost",
 };
```
````

```diff lang="javascript" title="config.js"
 const config = {
-  port: 3000,
+  port: 8080,
   host: "localhost",
 };
```

→ Port 3000 shows a red deletion marker. Port 8080 shows a green insertion marker. The block uses JavaScript syntax highlighting.

## How it works

Two conditions activate hybrid mode: the language must be `diff` AND `lang=` must specify the original language. The preprocessor runs before syntax highlighting:

1. Strips the `+`, `-`, or leading space prefix from each line
2. Generates `ins` line markers for `+` lines and `del` markers for `-` lines
3. Passes the stripped code to the highlighter using the `lang=` language

Lines with no recognized prefix (not starting with `+`, `-`, or space) pass through unchanged. Empty lines are kept as-is. If the prefix is `"+ "` (with a space), both characters are stripped. If it is `"+foo"` (no space), only the `+` is stripped.

After processing, the block effectively becomes the underlying language. The toolbar badge, `data-language` attribute, and language detection all reflect the swapped language.

The `+`/`-` indicators in the rendered output are CSS `::before` pseudo-elements on `.kz-code`, driven by the `--kz-ins-indicator` and `--kz-del-indicator` variables. No diff-specific CSS or JavaScript exists. Diff highlighting reuses the [line marker](/features/line-markers/) infrastructure entirely.

## Plain diff mode

Without `lang=`, a `diff` block renders as plain diff syntax with no Kazari marker classes. The raw `+`/`-` prefixes stay in the code and are tokenized by the highlighter's TextMate diff grammar. The block gets theme-colored text via diff scopes (headers, hunks, added/removed lines) but no green/red marker backgrounds or `+`/`-` CSS indicators.

## Combined with explicit markers

Diff-generated markers are appended to any explicit markers in the meta string. For example, `diff lang="go" {1}` highlights line 1 as `mark` in addition to the diff-generated `ins`/`del` markers. Standard [overlap priority](/features/line-markers/) applies: `mark` < `del` < `ins`.

## Go API

```go
html, err := engine.Render(code, kazari.Options{
    Lang:     "diff",
    DiffLang: "go",
    Title:    "changes.go",
})
```

`DiffLang` is available on `Options` only. It is not part of `BlockDefaults` or `LanguageDefaults` and has no config file equivalent.

## Configuration

| Option | Layer | Syntax | Description |
|---|---|---|---|
| `diff` | Meta string (language) | `` ```diff `` | Set the block language to diff |
| `lang=` | Meta string | `lang="go"` | Original language for syntax highlighting |
| `Options.DiffLang` | Go API (per-block) | `DiffLang: "go"` | Same as `lang=` via Go API |

Both the `diff` language and `lang=` are required for hybrid mode.

## Edge cases

- Without `lang=`, a `diff` block renders as plain diff syntax (TextMate grammar) with no Kazari ins/del markers.
- `DiffLang` is not part of `BlockDefaults` or `LanguageDefaults`. Per-language defaults for diff are not possible.
- No config file support for `DiffLang`.
- Diff-derived markers are added after collapse resolution. Threshold-based collapse preview does not account for diff-derived ins/del lines when computing visible segments.
- Filename comment extraction runs on raw (unstripped) code. A comment on a `+`/`-` line (e.g., `+// file.go`) is not detected as a filename because the prefix blocks the comment-pattern match.
- Diff reuses all marker CSS variables (`--kz-ins-bg`, `--kz-del-bg`, `--kz-ins-indicator`, `--kz-del-indicator`, etc.). No diff-specific CSS or JS exists.
