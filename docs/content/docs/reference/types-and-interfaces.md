---
title: "Types & Interfaces"
description: "Complete reference for all exported types, interfaces, and constants in the kazari package."
sidebar:
  order: 6
---

Every exported type from the `kazari` package, organized by role. For the methods that accept and return these types, see [Engine API](/reference/engine-api/). For the functional options that configure the engine, see [Configuration Options](/reference/configuration-options/).

## Interfaces

### `Highlighter`

```go
type Highlighter interface {
    Tokenize(code, lang, theme string) ([][]Token, error)
    GetThemeColors(theme string) (ThemeInfo, error)
    GetLoadedLanguages() []string
}
```

Abstracts syntax highlighting. The Nuri and Chroma adapters implement this interface. Any implementation works, including mocks for testing. Pass a `Highlighter` to `kazari.WithHighlighter()` at engine construction.

### `DualThemeTokenizer`

```go
type DualThemeTokenizer interface {
    TokenizeDual(code, lang, lightTheme, darkTheme string) (light, dark [][]Token, err error)
}
```

Optional capability for highlighters that resolve two themes from a single tokenization pass. The engine type-asserts this interface on its `Highlighter` at render time and falls back to two separate `Tokenize` calls when it is absent.

The two returned token streams must have identical boundaries: same line count, same per-line token count, and identical `Content` fields. Only colors may differ. The [Nuri adapter](/integrations/nuri-adapter/) implements this interface.

## Core types

### `Token`

```go
type Token struct {
    Content   string
    Color     string // foreground hex color (e.g. "#cf222e")
    BgColor   string // background hex color (usually empty)
    FontStyle int    // bitmask: see constants below
}
```

Returned by `engine.Tokenize()` and consumed by the highlighter interface.

| Constant | Value | Effect |
|----------|-------|--------|
| `FontStyleNone` | `0` | No decoration |
| `FontStyleItalic` | `1` | Italic |
| `FontStyleBold` | `2` | Bold |
| `FontStyleUnderline` | `4` | Underline |
| `FontStyleStrikethrough` | `8` | Strikethrough |

Combine with bitwise OR: `FontStyleBold | FontStyleItalic` produces bold italic.

### `Options`

```go
type Options struct {
    Lang            string
    Title           string
    Theme           string
    Frame           *Frame
    LineNumbers     *bool
    StartLineNumber *int
    Wrap            *bool
    PreserveIndent  *bool
    HangingIndent   *int
    DiffLang        string
    LineMarkers     []LineMarker
    InlineMarkers   []InlineMarker
    FocusLines      []Range
    Collapse        *CollapseOptions
    WithOutput      *bool
    OutputCollapsed *bool
    OutputLabel     string
}
```

Per-block configuration passed to `engine.Render()`. Pointer fields use nil to mean "inherit from engine defaults." Non-nil values override. For the string-based equivalent parsed from markdown fences, see [Meta String Syntax](/reference/meta-string-syntax/).

Output panel fields: `WithOutput` enables the output panel for this block (nil inherits from engine). `OutputCollapsed` starts the panel collapsed (nil uses engine default). `OutputLabel` overrides the toggle label (empty uses locale default).

### `BlockInfo`

```go
type BlockInfo struct {
    Lang      string
    Title     string
    Frame     Frame
    RawCode   string
    LineCount int
    Theme     string
    Meta      string
}
```

Metadata about a rendered code block, passed to every `WithPostRender` callback. `RawCode` is the original source before any processing. `Meta` is the raw meta string from the fence info line. See [Post-Render Callbacks](/plugins/post-render-callbacks/) for usage.

### `CollapseOptions`

```go
type CollapseOptions struct {
    Enabled   bool
    Disabled  bool
    Ranges    []Range
    Style     *CollapseStyle
    Threshold *int
}
```

Per-block collapse configuration within `Options`. `Enabled` forces threshold-based collapse. `Disabled` prevents collapse even when the engine threshold is met. `Ranges` specifies line ranges to collapse. `Style` and `Threshold` override engine defaults for this block.

### `Range`

```go
type Range struct {
    Start int
    End   int
}
```

An inclusive, 1-based line range. A single-line range has `Start == End`.

## Marker types

### `MarkerType`

```go
type MarkerType int

const (
    MarkerMark MarkerType = iota // highlight (default, lowest priority)
    MarkerDel                    // deleted
    MarkerIns                    // inserted (highest priority)
)
```

Higher numeric value wins when markers overlap on the same line.

### `LineMarker`

```go
type LineMarker struct {
    Type  MarkerType
    Lines []Range
    Label string // empty for unlabeled
}
```

A line-level marker with an optional text label. When `Label` is set, the marked lines receive a badge pill in the gutter.

### `InlineMarker`

```go
type InlineMarker struct {
    Type    MarkerType
    Text    string
    IsRegex bool
}
```

An inline text marker. When `IsRegex` is true, `Text` is a regular expression pattern. All occurrences (or matches) within the code are wrapped in `<mark>`, `<ins>`, or `<del>` elements.

## Frame and layout

### `Frame`

```go
type Frame int

const (
    FrameAuto     Frame = iota // auto-detect from language
    FrameCode                  // editor frame
    FrameTerminal              // terminal frame
    FrameNone                  // no frame
)
```

### `BlockDefaults`

```go
type BlockDefaults struct {
    Wrap           bool
    PreserveIndent bool
    HangingIndent  int
    LineNumbers    bool
    Frame          Frame
}
```

Engine-wide per-block rendering defaults. Pass to `kazari.WithDefaults()`. See [Configuration Options](/reference/configuration-options/) for default values.

### `TerminalDotStyle`

```go
type TerminalDotStyle int

const (
    DotsColored TerminalDotStyle = iota // macOS red/yellow/green DOM spans
    DotsMinimal                          // CSS-only monochrome dots via SVG mask
)
```

### `LangIconMode`

```go
type LangIconMode int

const (
    LangIconNone    LangIconMode = iota // text label only
    LangIconOnly                        // icon replaces text
    LangIconAndText                     // icon shown before text label
)
```

## Dark mode

### `DarkMode`

A sealed interface. Only three constructor functions produce valid values:

| Constructor | CSS output |
|-------------|-----------|
| `SelectorMode(selector string)` | `{root}{selector} { ... }` |
| `MediaQueryMode()` | `@media (prefers-color-scheme: dark) { ... }` |
| `BothMode(selector string)` | Selector block and media query block |

The `selector` parameter in `SelectorMode` and `BothMode` accepts any CSS selector that can be appended to `:root`. This includes class selectors (`.dark`), attribute selectors (`[data-theme="dark"]`), and other compound selectors.

Pass the result to `kazari.WithDarkMode()`. See [Themes & Dark Mode](/styling/themes-and-dark-mode/) for usage details.

## Theme customization

### `ThemeInfo`

```go
type ThemeInfo struct {
    FG           string // default foreground (e.g. "#24292f")
    BG           string // default background (e.g. "#ffffff")
    SelectionBG  string // editor.selectionBackground
    LineNumberFG string // editorLineNumber.foreground
    FoldBG       string // editor.foldBackground
}
```

Colors extracted from a syntax theme by `Highlighter.GetThemeColors()`. Used in `WithThemeCustomizer` callbacks and returned by the highlighter for CSS variable generation.

### `ThemeAdjustments`

```go
type ThemeAdjustments struct {
    Hue     *float64      // target hue in degrees (0-360); nil = unchanged
    Chroma  *float64      // target chroma (0-0.4 typical); nil = unchanged
    Targets AdjustTargets // bitmask; zero value = AdjustBackgrounds
}
```

Tints extracted theme colors in OKLCH space. Adjustments apply to editor chrome colors (the `ThemeInfo` fields), not to individual syntax token colors. Pass to `kazari.WithThemeAdjustments()`.

### `AdjustTargets`

```go
type AdjustTargets int

const (
    AdjustBackgrounds AdjustTargets = 1 << iota // BG, SelectionBG, FoldBG
    AdjustForegrounds                           // FG, LineNumberFG
)
```

Combine with bitwise OR to target both: `AdjustBackgrounds | AdjustForegrounds`.

## Collapsible

### `CollapseStyle`

```go
type CollapseStyle int

const (
    CollapseGithub          CollapseStyle = iota // one-way expand, summary disappears
    CollapseCollapsibleStart                     // re-collapsible, summary above content
    CollapseCollapsibleEnd                       // re-collapsible, summary below content
    CollapseCollapsibleAuto                      // auto: end if section reaches last line, else start
)
```

### `CollapsibleConfig`

```go
type CollapsibleConfig struct {
    LineThreshold         int
    PreviewLines          int
    DefaultCollapsed      bool
    PreserveIndent        bool
    Style                 CollapseStyle
    ExpandButtonText      string
    CollapseButtonText    string
    ExpandedAnnouncement  string // screen reader
    CollapsedAnnouncement string // screen reader
}
```

Configures threshold-based collapsible blocks. Pass to `kazari.WithCollapsible()`. See [Configuration Options](/reference/configuration-options/) for default values and YAML equivalents.

## Style

### `StyleValue`

```go
type StyleValue struct {
    Value string // used when both themes share the same value
    Dark  string // dark theme override
    Light string // light theme override
}
```

Represents a CSS variable override that can be universal or per-theme. Pass as values in the map to `kazari.WithThemedStyleOverrides()`.

| Method | Returns |
|--------|---------|
| `IsThemed() bool` | `true` if `Dark` or `Light` is non-empty |
| `LightValue() string` | `Light` if themed, otherwise `Value` |
| `DarkValue() string` | `Dark` if themed, otherwise `Value` |

## Assets

### `Assets`

```go
type Assets struct {
    CSS AssetFile
    JS  AssetFile
}
```

Returned by `engine.Assets()`.

### `AssetFile`

```go
type AssetFile struct {
    Content  string
    Hash     string // 8-char hex (FNV-1a 32-bit)
    Filename string // e.g. "kazari-a1b2c3d4.css"
}
```

## Functional option type

### `Option`

```go
type Option func(*engineBuilder)
```

A functional option passed to `kazari.New()`. The `engineBuilder` type is unexported; create options via the `With*` functions documented in [Configuration Options](/reference/configuration-options/).

## File configuration types

These types are returned by `ParseConfig()` and consumed by `FileConfigToOptions()`. Pointer fields use `nil` to mean "use default" rather than the Go zero value.

### `FileConfig`

The top-level struct deserialized from `kazari.config.yaml` or `kazari.config.json`. Fields map 1:1 to YAML/JSON keys documented in [File-Based Config](/reference/file-based-config/).

| Field | Type | YAML key |
|-------|------|----------|
| `Themes` | `*ThemesFileConfig` | `themes` |
| `DarkMode` | `*DarkModeFileConfig` | `darkMode` |
| `CopyButton` | `*bool` | `copyButton` |
| `FullscreenButton` | `*bool` | `fullscreenButton` |
| `WrapButton` | `*bool` | `wrapButton` |
| `ThemeToggleButton` | `*bool` | `themeToggleButton` |
| `LineNumbers` | `*bool` | `lineNumbers` |
| `FrameDetection` | `*bool` | `frameDetection` |
| `FileNameExtraction` | `*bool` | `fileNameExtraction` |
| `LanguageBadge` | `*bool` | `languageBadge` |
| `ThemedScrollbars` | `*bool` | `themedScrollbars` |
| `ThemedSelection` | `*bool` | `themedSelection` |
| `ContentExclusion` | `*bool` | `contentExclusion` |
| `MermaidPassThrough` | `*bool` | `mermaidPassThrough` |
| `TerminalCommentStripping` | `*bool` | `terminalCommentStripping` |
| `DataLineCount` | `*bool` | `dataLineCount` |
| `FileIcons` | `*bool` | `fileIcons` |
| `Links` | `*bool` | `links` |
| `StyleReset` | `*bool` | `styleReset` |
| `Minify` | `*bool` | `minify` |
| `TabWidth` | `*int` | `tabWidth` |
| `MinContrast` | `*float64` | `minContrast` |
| `CascadeLayer` | `*string` | `cascadeLayer` |
| `ThemeCSSRoot` | `*string` | `themeCSSRoot` |
| `Locale` | `*string` | `locale` |
| `TerminalDotStyle` | `*string` | `terminalDotStyle` |
| `LanguageIconMode` | `*string` | `languageIconMode` |
| `Defaults` | `*BlockDefaultsFileConfig` | `defaults` |
| `LanguageDefaults` | `map[string]*BlockDefaultsFileConfig` | `languageDefaults` |
| `Collapsible` | `*CollapsibleFileConfig` | `collapsible` |
| `LanguageAliases` | `map[string]string` | `languageAliases` |
| `UIStrings` | `map[string]string` | `uiStrings` |
| `StyleOverrides` | `map[string]any` | `styleOverrides` |

### `ThemesFileConfig`

```go
type ThemesFileConfig struct {
    Light string `yaml:"light" json:"light"`
    Dark  string `yaml:"dark"  json:"dark"`
}
```

### `DarkModeFileConfig`

```go
type DarkModeFileConfig struct {
    Kind     string `yaml:"kind"     json:"kind"`
    Selector string `yaml:"selector" json:"selector"`
}
```

### `BlockDefaultsFileConfig`

```go
type BlockDefaultsFileConfig struct {
    Wrap           *bool   `yaml:"wrap"           json:"wrap"`
    PreserveIndent *bool   `yaml:"preserveIndent" json:"preserveIndent"`
    HangingIndent  *int    `yaml:"hangingIndent"  json:"hangingIndent"`
    LineNumbers    *bool   `yaml:"lineNumbers"    json:"lineNumbers"`
    Frame          *string `yaml:"frame"          json:"frame"`
}
```

### `CollapsibleFileConfig`

```go
type CollapsibleFileConfig struct {
    LineThreshold         *int    `yaml:"lineThreshold"         json:"lineThreshold"`
    PreviewLines          *int    `yaml:"previewLines"          json:"previewLines"`
    DefaultCollapsed      *bool   `yaml:"defaultCollapsed"      json:"defaultCollapsed"`
    PreserveIndent        *bool   `yaml:"preserveIndent"        json:"preserveIndent"`
    Style                 *string `yaml:"style"                 json:"style"`
    ExpandButtonText      *string `yaml:"expandButtonText"      json:"expandButtonText"`
    CollapseButtonText    *string `yaml:"collapseButtonText"    json:"collapseButtonText"`
    ExpandedAnnouncement  *string `yaml:"expandedAnnouncement"  json:"expandedAnnouncement"`
    CollapsedAnnouncement *string `yaml:"collapsedAnnouncement" json:"collapsedAnnouncement"`
}
```
