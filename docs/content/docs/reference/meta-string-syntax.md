---
title: "Meta String Syntax"
description: "Complete reference for every token in the markdown fence info string."
tags: [line-annotations]
sidebar:
  order: 4
---

The meta string is the text after the opening triple backticks in a fenced code block. Kazari parses it into per-block options, markers, focus ranges, and collapse directives. Tokens are separated by whitespace.

````
```go title="main.go" {3-5} showLineNumbers
````

A Go code block with a `main.go` title tab, lines 3-5 highlighted, and line numbers enabled.

## Block options

| Token | Description | Example |
|-------|-------------|---------|
| *language* | First bare word sets the language | `go`, `typescript`, `bash` |
| `title="..."` | Title shown in the frame title bar | `title="main.go"` |
| `frame=` | Frame type: `code`, `terminal`, `none`, `auto` | `frame=terminal` |
| `theme="..."` | Override the theme for this block | `theme="dracula"` |
| `showLineNumbers` | Enable line numbers | `showLineNumbers` |
| `showLineNumbers=false` | Disable line numbers (overrides engine default) | `showLineNumbers=false` |
| `startLineNumber=N` | First displayed line number | `startLineNumber=10` |
| `wrap` | Enable word wrap | `wrap` |
| `preserveIndent` | Preserve indentation on wrapped lines | `preserveIndent` |
| `preserveIndent=false` | Disable indent preservation | `preserveIndent=false` |
| `hangingIndent=N` | Extra indent columns for wrapped continuation lines | `hangingIndent=2` |
| `lang="..."` | Language for diff+syntax hybrid rendering | `lang="go"` |

The language must be the first token and must not contain `=` or start with `{`, `"`, or `'`. The `lang=` token specifies the original language when using `diff` as the fence language. See [Diff Highlighting](/features/diff-highlighting/).

Values after `=` accept double quotes, single quotes, or bare words: `title="My Title"`, `title='My Title'`, and `title=MyTitle` are all valid.

## Line markers

| Token | Type | Description | Example |
|-------|------|-------------|---------|
| `{N-M,...}` | mark | Highlight lines with the default style | `{3-5}`, `{2,4-6}` |
| `{"Label":N-M}` | mark | Highlight lines with a labeled badge | `{"Added":3-5}` |
| `ins={N-M,...}` | ins | Mark lines as inserted (green) | `ins={10-12}` |
| `ins={"Label":N-M}` | ins | Mark lines as inserted with a label | `ins={"New":1-3}` |
| `del={N-M,...}` | del | Mark lines as deleted (red) | `del={7}` |
| `del={"Label":N-M}` | del | Mark lines as deleted with a label | `del={"Removed":4-6}` |
| `add={...}` | ins | Alias for `ins=` | `add={1-3}` |
| `add={"Label":N-M}` | ins | Alias for `ins={"Label":N-M}` | `add={"New":1-3}` |
| `rem={...}` | del | Alias for `del=` | `rem={7}` |
| `rem={"Label":N-M}` | del | Alias for `del={"Label":N-M}` | `rem={"Old":4-6}` |

Multiple line marker tokens in the same meta string each produce a separate marker entry. When markers overlap on the same line, higher-priority types win: mark (lowest) < del < ins (highest). See [Line Markers](/features/line-markers/) for visual behavior.

## Inline markers

| Token | Type | Description | Example |
|-------|------|-------------|---------|
| `"text"` | mark | Highlight all occurrences of text | `"useState"` |
| `'text'` | mark | Same, with single quotes | `'myFunction'` |
| `/regex/` | mark | Highlight all regex matches | `/func\s+\w+/` |
| `ins="text"` | ins | Mark text occurrences as inserted | `ins="added"` |
| `ins=/regex/` | ins | Mark regex matches as inserted | `ins=/new\w+/` |
| `del="text"` | del | Mark text occurrences as deleted | `del="removed"` |
| `del=/regex/` | del | Mark regex matches as deleted | `del=/old\w+/` |
| `add="text"` | ins | Alias for `ins="text"` | `add="added"` |
| `add=/regex/` | ins | Alias for `ins=/regex/` | `add=/new\w+/` |
| `rem="text"` | del | Alias for `del="text"` | `rem="removed"` |
| `rem=/regex/` | del | Alias for `del=/regex/` | `rem=/old\w+/` |

Bare quoted strings (`"text"` or `'text'`) at any position are parsed as inline markers, never as the language. Use `\/` to escape literal slashes inside regex patterns: `/\/path\//` matches `/path/`. See [Inline Markers](/features/inline-markers/) for multi-token spanning and open-start/open-end behavior.

## Focus

| Token | Description | Example |
|-------|-------------|---------|
| `focus={N-M,...}` | Dim all lines except the specified ranges | `focus={1-3,7}` |

Non-focused lines receive reduced opacity via `--kz-focus-dimmed-opacity`. See [Focus Lines](/features/focus-lines/).

## Collapsible

| Token | Description | Example |
|-------|-------------|---------|
| `collapse` | Enable threshold-based collapse for this block | `collapse` |
| `nocollapse` | Prevent collapse even when engine threshold is met | `nocollapse` |
| `collapse={N-M,...}` | Collapse specific line ranges | `collapse={5-15}` |
| `collapseStyle=` | Visual style: `github`, `collapsible-start`, `collapsible-end`, `collapsible-auto` | `collapseStyle=collapsible-start` |
| `collapseThreshold=N` | Override engine threshold for this block (must be > 0) | `collapseThreshold=30` |

Multiple `collapse={...}` tokens accumulate: each appends to the list of collapsed ranges. `collapseStyle=` and `collapseThreshold=` create the collapse configuration even without a `collapse` or `nocollapse` token. Unrecognized style values fall back to `github`. See [Collapsible Sections](/features/collapsible-sections/) for threshold and range behavior.

## Output panel

| Token | Description | Example |
|-------|-------------|---------|
| `withOutput` | Enable the output panel for this block | `withOutput` |
| `outputCollapsed` | Start the output panel collapsed | `outputCollapsed` |
| `outputCollapsed=false` | Start the output panel expanded (overrides engine default) | `outputCollapsed=false` |
| `outputLabel="..."` | Custom toggle label for the output panel | `outputLabel="Run result"` |

The output panel splits a code block at the `---output---` separator (configurable via `WithOutputSeparator`). Code above the separator is syntax-highlighted; content below renders in a plain panel beneath the block. Requires `WithOutputPanel(true)` at engine level. See [Output Panel](/features/output-panel/).

## Range syntax

Line markers, focus, and collapse tokens all share the same range format:

- Comma-separated: `3-5,8,10-12`
- Single line: `8` is equivalent to `8-8`
- All line numbers are 1-based and inclusive
- Whitespace around commas and numbers is trimmed

## Quoting rules

Values after `=` accept three formats:

- **Double-quoted:** `title="My Title"`. Use `\"` to include a literal double quote.
- **Single-quoted:** `title='My Title'`. Use `\'` to include a literal single quote.
- **Bare word:** `frame=terminal`. No quoting needed when the value contains no spaces.

Inline text markers use the same quoting rules. Regex patterns are delimited by `/` and support `\/` for literal slashes.
