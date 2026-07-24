---
title: "File-Based Config"
description: "Configure Kazari via YAML or JSON files with auto-discovery."
sidebar:
  order: 1
---

A config file sets engine-wide defaults, feature toggles, and style overrides without writing Go code. Kazari loads configuration from `kazari.config.yaml`, `.yml`, or `.json` files.

## Auto-discovery

Pass a directory path to `WithConfigDir`. Kazari searches for config files in this order:

1. `kazari.config.yaml`
2. `kazari.config.yml`
3. `kazari.config.json`

The first file found is loaded; others are ignored even if they exist. If no file is found, the engine is not affected.

```go
engine := kazari.New(
    kazari.WithConfigDir("./config"),
    kazari.WithHighlighter(hl),
    kazari.WithCopyButton(true), // overrides the file config value if set
)
```

If the file exists but fails to parse or validate, `WithConfigDir` logs a warning and continues without applying the config. It does not return an error or panic.

## Loading manually

To load a config file directly by path:

```go
opts, err := kazari.LoadConfig("kazari.config.yaml")
if err != nil {
    log.Fatal(err)
}
engine := kazari.New(append(opts, kazari.WithHighlighter(hl))...)
```

To parse raw bytes (for example, from an embedded file):

```go
fc, err := kazari.ParseConfig(data, "yaml") // "yaml" or "json"
if err != nil {
    log.Fatal(err)
}
opts, err := kazari.FileConfigToOptions(fc)
```

## Themes

```yaml
themes:
  light: "github-light"
  dark: "github-dark"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `themes.light` | `string` | `"github-light"` | Light theme name (passed to the highlighter) |
| `themes.dark` | `string` | `"github-dark"` | Dark theme name (empty string for single-theme) |

## Dark mode

```yaml
darkMode:
  kind: "selector"
  selector: ".dark"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `darkMode.kind` | `string` | `"selector"` | One of: `selector`, `mediaQuery`, `both` |
| `darkMode.selector` | `string` | `".dark"` | CSS selector for dark mode (required for `selector` and `both` kinds) |

`selector` toggles dark mode via a class (e.g., `.dark` on `<html>`). `mediaQuery` uses `@media (prefers-color-scheme: dark)`. `both` emits both so a manual toggle can override OS preference.

## Feature toggles

All feature toggles are optional booleans. Omitted fields keep the engine default.

```yaml
copyButton: true
fullscreenButton: true
wrapButton: true
themeToggleButton: false
lineNumbers: false
frameDetection: true
fileNameExtraction: true
languageBadge: true
themedScrollbars: true
themedSelection: false
contentExclusion: false
mermaidPassThrough: true
terminalCommentStripping: true
dataLineCount: true
fileIcons: true
links: false
styleReset: true
minify: true
```

| Config key | Default | Go API |
|------------|---------|--------|
| `copyButton` | `true` | `WithCopyButton` |
| `fullscreenButton` | `true` | `WithFullscreenButton` |
| `wrapButton` | `true` | `WithWrapButton` |
| `themeToggleButton` | `false` | `WithThemeToggle` |
| `lineNumbers` | `false` | `WithLineNumbers` |
| `frameDetection` | `true` | `WithFrameDetection` |
| `fileNameExtraction` | `true` | `WithFileNameExtraction` |
| `languageBadge` | `true` | `WithLanguageBadge` |
| `themedScrollbars` | `true` | `WithThemedScrollbars` |
| `themedSelection` | `false` | `WithThemedSelectionColors` |
| `contentExclusion` | `false` | `WithContentExclusion` (deprecated no-op) |
| `mermaidPassThrough` | `true` | `WithMermaidPassThrough` |
| `terminalCommentStripping` | `true` | `WithTerminalCommentStripping` |
| `dataLineCount` | `true` | `WithDataLineCount` |
| `fileIcons` | `true` | `WithFileIcons` |
| `links` | `false` | `WithLinks` |
| `styleReset` | `true` | `WithStyleReset` |
| `minify` | `true` | `WithMinify` |

## Engine settings

```yaml
tabWidth: 2
minContrast: 5.5
cascadeLayer: "kazari"
themeCSSRoot: ":root"
terminalDotStyle: "colored"
languageIconMode: "none"
locale: "en-US"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `tabWidth` | `int` | `2` | Spaces per tab for indentation calculation |
| `minContrast` | `float` | `5.5` | Minimum WCAG contrast ratio for syntax colors. Set to `0` to disable |
| `cascadeLayer` | `string` | `"kazari"` | CSS `@layer` name for the generated stylesheet |
| `themeCSSRoot` | `string` | `":root"` | CSS selector for emitting theme variables |
| `terminalDotStyle` | `string` | `"colored"` | Terminal frame dot style: `colored` or `minimal` |
| `languageIconMode` | `string` | `"none"` | Language icon display mode: `none`, `iconOnly`, or `iconAndText` |
| `locale` | `string` | `"en-US"` | Locale for UI strings (copy button labels, etc.) |

## Block defaults

```yaml
defaults:
  wrap: false
  preserveIndent: true
  hangingIndent: 0
  lineNumbers: false
  frame: "auto"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `defaults.wrap` | `bool` | `false` | Enable word wrap |
| `defaults.preserveIndent` | `bool` | `true` | Preserve leading indentation on wrapped lines |
| `defaults.hangingIndent` | `int` | `0` | Extra indent (in `ch`) for continuation lines |
| `defaults.lineNumbers` | `bool` | `false` | Show line numbers |
| `defaults.frame` | `string` | `"auto"` | Frame style: `auto`, `code`, `terminal`, `none` |

## Language defaults

Per-language overrides. Keys accept comma-separated language names to apply the same config to multiple languages:

```yaml
languageDefaults:
  "bash,sh,zsh":
    wrap: true
    frame: "terminal"
  "python":
    lineNumbers: true
  "go":
    lineNumbers: true
    frame: "code"
```

Each language group accepts the same five fields as `defaults`. Language defaults override engine defaults but are overridden by per-block meta strings.

## Collapsible

Setting any `collapsible` field enables collapsible sections. When a `collapsible` section is present in the config file, unset subfields use these defaults:

```yaml
collapsible:
  lineThreshold: 15
  previewLines: 5
  defaultCollapsed: true
  preserveIndent: true
  style: "github"
  expandButtonText: ""
  collapseButtonText: ""
  expandedAnnouncement: ""
  collapsedAnnouncement: ""
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `collapsible.lineThreshold` | `int` | `15` | Blocks with more lines than this are collapsible |
| `collapsible.previewLines` | `int` | `5` | Number of visible lines when collapsed |
| `collapsible.defaultCollapsed` | `bool` | `true` | Whether blocks start collapsed |
| `collapsible.preserveIndent` | `bool` | `true` | Keep indentation in preview lines |
| `collapsible.style` | `string` | `"github"` | Collapse style: `github`, `collapsibleStart`, `collapsibleEnd`, `collapsibleAuto` |
| `collapsible.expandButtonText` | `string` | | Custom expand button label |
| `collapsible.collapseButtonText` | `string` | | Custom collapse button label |
| `collapsible.expandedAnnouncement` | `string` | | Screen reader announcement when expanded |
| `collapsible.collapsedAnnouncement` | `string` | | Screen reader announcement when collapsed |

## Output panel

```yaml
outputPanel: true
outputDefaultCollapsed: false
outputSeparator: "---output---"
```

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `outputPanel` | `bool` | `false` | Enable the output panel feature for separating code from command output |
| `outputDefaultCollapsed` | `bool` | `false` | Whether output panels start collapsed |
| `outputSeparator` | `string` | `"---output---"` | Separator string that splits code from output content in a fenced block |

When enabled, a code block whose content contains the separator string is split into two panels: the code above the separator and the output below it. The separator line itself is not rendered.

## Language aliases

```yaml
languageAliases:
  javascript: "js"
  typescript: "ts"
```

Maps alias names to language names recognized by the highlighter. When a code block uses an alias as its language, the engine resolves it to the target name before tokenization.

## UI string overrides

```yaml
uiStrings:
  "copy.label": "Copy"
  "copy.success": "Copied!"
```

Overrides localized UI strings (button labels, announcements). See the [Localization](/features/localization/) page for the full list of string keys.

## Style overrides

The `styleOverrides` map sets CSS custom property values. Three value formats are accepted.

**Plain string**: same value for both light and dark themes.

```yaml
styleOverrides:
  radius: "0.5rem"
  toolbar-height: "2.5rem"
```

**Array**: dark value first, light value second.

```yaml
styleOverrides:
  editor-bg: ["#1e293b", "#f8fafc"]
  editor-fg: ["#e2e8f0", "#1e293b"]
```

**Map**: explicit `dark` and `light` keys (either or both may be omitted).

```yaml
styleOverrides:
  editor-bg:
    dark: "#1e293b"
    light: "#f8fafc"
```

Setting only one side overrides that theme value while leaving the other from the theme:

```yaml
styleOverrides:
  editor-bg:
    dark: "#1e293b"
```

### Key normalization

Bare names are automatically prefixed with `--kz-`. Keys that already start with `--` are used as-is.

| Config key | Resolved CSS variable |
|------------|----------------------|
| `radius` | `--kz-radius` |
| `editor-bg` | `--kz-editor-bg` |
| `--kz-shadow` | `--kz-shadow` (kept as-is) |
| `--custom-var` | `--custom-var` (kept as-is) |

## Validation

Config files are validated in two passes. The first pass runs during decoding: the YAML decoder rejects unknown keys and type mismatches at parse time. JSON uses `DisallowUnknownFields` for the same behavior. The second pass checks value constraints immediately after decoding.

Both passes must succeed before the config is applied.

### Validated fields

| Field | Rule | Valid | Invalid |
|-------|------|-------|---------|
| `darkMode.kind` | Must be one of the allowed values | `"selector"`, `"mediaQuery"`, `"both"` | `"auto"`, `"css"` |
| `darkMode.selector` | Required when kind is `selector` or `both` | `".dark"`, `"#app"` | omitted with kind `"selector"` |
| `tabWidth` | Integer >= 1 | `2`, `4`, `8` | `0`, `-1` |
| `minContrast` | Float in range [0, 21] | `5.5`, `4.5`, `0` (disabled), `21` (max) | `-1.0`, `22.0` |
| `terminalDotStyle` | Must be one of the allowed values | `"colored"`, `"minimal"` | `"round"` |
| `languageIconMode` | Must be one of the allowed values | `"none"`, `"iconOnly"`, `"iconAndText"` | `"icon"` |
| `defaults.frame` | Must be one of the allowed values | `"auto"`, `"code"`, `"terminal"`, `"none"` | `"editor"` |
| `defaults.hangingIndent` | Integer >= 0 | `0`, `2`, `4` | `-1` |
| `languageDefaults.*.frame` | Same as `defaults.frame` | `"auto"`, `"code"`, `"terminal"`, `"none"` | `"editor"` |
| `languageDefaults.*.hangingIndent` | Integer >= 0 | `0`, `2`, `4` | `-1` |
| `collapsible.lineThreshold` | Integer >= 1 | `10`, `20` | `0`, `-5` |
| `collapsible.previewLines` | Integer >= 1 | `3`, `5` | `0`, `-1` |
| `collapsible.style` | Must be one of the allowed values | `"github"`, `"collapsibleStart"`, `"collapsibleEnd"`, `"collapsibleAuto"` | `"accordion"` |

### Style override validation

Each key in `styleOverrides` must match one of three formats:

| Format | Valid | Invalid |
|--------|-------|---------|
| String | `"0.5rem"`, `"'Fira Code', monospace"` | `42`, `true` (non-string types) |
| Array | `["#1e293b", "#ffffff"]` (exactly 2 string elements) | `["#1e293b"]` (wrong length), `[42, "#fff"]` (non-string) |
| Map | `{light: "#fff", dark: "#000"}`, `{dark: "#000"}` | `{}` (empty), `{light: 42}` (non-string value) |

### Unknown fields

Both YAML and JSON parsers reject unknown fields at decode time:

```yaml
# This produces a parse error:
unknownField: true
```

```
kazari: parsing YAML config: line 5: field unknownField not found in type kazari.FileConfig
```

## Error reporting

How validation errors reach the application depends on the loading path.

**`LoadConfig()` / `ParseConfig()` (direct API)** return errors directly. The caller decides how to handle them:

```go
opts, err := kazari.LoadConfig("kazari.config.yaml")
if err != nil {
    log.Fatalf("bad config: %v", err)
}
engine := kazari.New(append([]kazari.Option{kazari.WithHighlighter(hl)}, opts...)...)
```

**`WithConfigDir()` (discovery-based loading)** does not return an error because it runs inside `kazari.New()` as a functional option. Config errors are routed through the `WarningHandler`. If no handler is set, errors are printed via `log.Print`:

```go
engine := kazari.New(
    kazari.WithHighlighter(hl),
    kazari.WithWarningHandler(func(msg string) {
        log.Printf("kazari warning: %s", msg)
    }),
    kazari.WithConfigDir("."),
)
```

`WithWarningHandler` must appear before `WithConfigDir` in the option list so the handler is registered when the config file is loaded. If no handler is set, the default `log.Print` fallback ensures errors are never silently dropped.

When no config file is found in the directory, `WithConfigDir` is a silent no-op (no warning, no error).

## Composition with the Go API

`WithConfigDir` is an option like any other. Its position in the `New()` call determines priority:

```go
engine := kazari.New(
    kazari.WithCopyButton(false),   // overridden by file config (comes before WithConfigDir)
    kazari.WithConfigDir("./cfg"),  // file config applied here
    kazari.WithHighlighter(hl),     // not in file config, always applied
    kazari.WithMinify(true),        // comes after file config, takes precedence
)
```

Options are applied in order. `WithCopyButton(false)` is set first, then overridden by the file config because `WithConfigDir` comes after it. `WithMinify(true)` comes after the file config, so it takes precedence over any `minify` value in the file.

For scalar fields, the last writer wins. For map fields (`styleOverrides`, `languageDefaults`, `languageAliases`, `uiStrings`), keys from both sources coexist. When the same key appears in both, the later option's value is used.

## Go-API-only options

These options accept Go functions, interfaces, or structs and cannot be expressed in a config file. Pass them directly to `kazari.New()`:

| Option | Purpose |
|--------|---------|
| `WithHighlighter(hl)` | Sets the syntax highlighter (Nuri or Chroma adapter). Required. |
| `WithThemeCustomizer(f)` | Transforms extracted theme colors before CSS generation |
| `WithThemeAdjustments(adj)` | Tints theme colors in OKLCH space before the theme customizer |
| `WithPostRender(f)` | Registers a callback that receives rendered HTML and block metadata |
| `WithWarningHandler(f)` | Sets the function that receives non-fatal warnings |
| `WithFileIconResolver(f)` | Custom function mapping file extensions to icon SVG strings |
All other engine options have a corresponding config file key and can be set in either layer.

## Complete example

```yaml
themes:
  light: "github-light"
  dark: "github-dark"

darkMode:
  kind: "selector"
  selector: ".dark"

copyButton: true
fullscreenButton: true
wrapButton: true
themeToggleButton: true
lineNumbers: true
frameDetection: true
fileNameExtraction: true
languageBadge: true
themedScrollbars: true
links: true
styleReset: true
minify: true
fileIcons: true

outputPanel: true
outputDefaultCollapsed: false

tabWidth: 2
minContrast: 5.5
cascadeLayer: "kazari"
terminalDotStyle: "colored"
languageIconMode: "none"
locale: "en-US"

styleOverrides:
  radius: "0.3rem"
  font-size: "1rem"
  font-family: "'Fira Code', monospace"
  code-padding-block: "1.5rem"
  shadow:
    light: "0 2px 8px rgba(0,0,0,0.1)"
    dark: "none"
  editor-bg:
    light: "#ffffff"
    dark: "#1e293b"

defaults:
  wrap: false
  preserveIndent: true
  lineNumbers: false
  frame: "auto"

languageDefaults:
  "bash,sh,zsh":
    wrap: true
    frame: "terminal"
  "python":
    lineNumbers: true

collapsible:
  lineThreshold: 20
  previewLines: 5
  defaultCollapsed: true
  style: "github"

languageAliases:
  javascript: "js"
  typescript: "ts"

uiStrings:
  "copy.label": "Copy"
  "copy.success": "Copied!"
```
