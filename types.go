package kazari

import (
	"fmt"
	"hash/fnv"

	"github.com/frostybee/kazari/internal/config"
)

// Highlighter abstracts syntax highlighting. The default implementation wraps Nuri,
// but any implementation (including mocks) can be used.
type Highlighter interface {
	Tokenize(code, lang, theme string) ([][]Token, error)
	GetThemeColors(theme string) (ThemeInfo, error)
	GetLoadedLanguages() []string
}

// DualThemeTokenizer is an optional capability for highlighters that can
// resolve two themes from a single tokenization pass. TextMate scanning is
// theme-independent, so a capable highlighter halves dual-theme work.
// The two returned streams MUST have identical token boundaries (same line
// count, same per-line token count and Content). The engine type-asserts
// this interface on its Highlighter and falls back to two Tokenize calls
// when it is absent.
type DualThemeTokenizer interface {
	TokenizeDual(code, lang, lightTheme, darkTheme string) (light, dark [][]Token, err error)
}

// ThemeInfo holds colors extracted from a syntax theme for CSS variable generation.
type ThemeInfo struct {
	FG           string // default foreground (e.g. "#24292f")
	BG           string // default background (e.g. "#ffffff")
	SelectionBG  string // editor.selectionBackground
	LineNumberFG string // editorLineNumber.foreground
	FoldBG       string // editor.foldBackground
}

// ThemeAdjustments tints extracted theme colors in OKLCH space.
// Nil fields leave that channel unchanged. Adjustments apply to the editor
// colors used for CSS variables, not to individual syntax token colors.
type ThemeAdjustments struct {
	Hue     *float64      // target hue in degrees (0-360); nil = unchanged
	Chroma  *float64      // target chroma (0-0.4 typical); nil = unchanged
	Targets AdjustTargets // bitmask; zero value = AdjustBackgrounds
}

// AdjustTargets selects which extracted colors a ThemeAdjustments affects.
type AdjustTargets int

const (
	AdjustBackgrounds AdjustTargets = 1 << iota // BG, SelectionBG, FoldBG
	AdjustForegrounds                           // FG, LineNumberFG
)

// Token represents a single colored token from the highlighter.
type Token struct {
	Content   string
	Color     string // foreground hex color (e.g. "#cf222e")
	BgColor   string // background hex color (usually empty)
	FontStyle int    // bitmask: Italic=1, Bold=2, Underline=4, Strikethrough=8
}

// Font style bitmask constants.
const (
	FontStyleNone          = config.FontStyleNone
	FontStyleItalic        = config.FontStyleItalic
	FontStyleBold          = config.FontStyleBold
	FontStyleUnderline     = config.FontStyleUnderline
	FontStyleStrikethrough = config.FontStyleStrikethrough
)

// Frame determines the code block's visual frame type.
type Frame int

const (
	FrameAuto     Frame = iota // auto-detect from language
	FrameCode                  // editor frame
	FrameTerminal              // terminal frame
	FrameNone                  // no frame
)

// TerminalDotStyle determines how terminal frame dots are rendered.
type TerminalDotStyle int

const (
	DotsColored TerminalDotStyle = iota // macOS red/yellow/green DOM spans (default)
	DotsMinimal                          // CSS-only monochrome dots via SVG mask
)

// LangIconMode controls language icon display in the toolbar badge area.
type LangIconMode int

const (
	LangIconNone    LangIconMode = iota // text label only (default)
	LangIconOnly                        // icon replaces text
	LangIconAndText                     // icon shown before text label
)

// Range represents an inclusive 1-based line range.
type Range struct {
	Start int
	End   int
}

// MarkerType identifies the kind of line or inline marker.
type MarkerType int

// Priority order: mark(0) < del(1) < ins(2). Higher value wins in overlap resolution.
const (
	MarkerMark MarkerType = iota // highlight (default, lowest priority)
	MarkerDel                    // deleted
	MarkerIns                    // inserted (highest priority)
)

// LineMarker represents a line-level marker with optional label.
type LineMarker struct {
	Type  MarkerType
	Lines []Range
	Label string // empty for unlabeled
}

// InlineMarker represents an inline text marker.
type InlineMarker struct {
	Type    MarkerType
	Text    string
	IsRegex bool
}

// Options is the per-block configuration passed to Render or derived from meta string.
type Options struct {
	Lang            string
	Title           string
	Theme           string
	Frame           *Frame // nil = use default
	LineNumbers     *bool  // nil = use default
	StartLineNumber *int   // nil = use default (1)
	Wrap            *bool  // nil = use default
	PreserveIndent  *bool  // nil = use default
	HangingIndent   *int   // nil = use default
	DiffLang        string // original language for diff+syntax hybrid (e.g., "go")
	LineMarkers     []LineMarker
	InlineMarkers   []InlineMarker
	FocusLines      []Range
	Collapse        *CollapseOptions
}

// CollapseOptions holds per-block collapse configuration.
type CollapseOptions struct {
	Enabled   bool           // true = force threshold-based collapse
	Disabled  bool           // true = force no collapse (nocollapse)
	Ranges    []Range        // specific ranges to collapse
	Style     *CollapseStyle // nil = use engine default
	Threshold *int           // nil = use engine default; overrides engine LineThreshold for this block
}

// DarkMode controls how dark mode CSS is generated.
type DarkMode interface {
	darkMode() // sealed — only Kazari's implementations allowed
}

type selectorMode struct{ Selector string }
type mediaQueryMode struct{}
type bothMode struct{ Selector string }

func (selectorMode) darkMode()  {}
func (mediaQueryMode) darkMode() {}
func (bothMode) darkMode()      {}

// SelectorMode uses a CSS class selector for dark mode switching.
// Example: SelectorMode(".dark") produces :root.dark { ... }
func SelectorMode(selector string) DarkMode { return selectorMode{Selector: selector} }

// MediaQueryMode uses prefers-color-scheme media query for dark mode.
func MediaQueryMode() DarkMode { return mediaQueryMode{} }

// BothMode combines media query with a class toggle override.
func BothMode(selector string) DarkMode { return bothMode{Selector: selector} }

// BlockDefaults holds default block rendering properties.
type BlockDefaults struct {
	Wrap           bool
	PreserveIndent bool
	HangingIndent  int
	LineNumbers    bool
	Frame          Frame
}

// CollapseStyle identifies the visual style for range-based collapsible sections.
type CollapseStyle int

const (
	CollapseGithub          CollapseStyle = iota // one-way expand, summary disappears
	CollapseCollapsibleStart                     // re-collapsible, summary above content
	CollapseCollapsibleEnd                       // re-collapsible, summary below content
	CollapseCollapsibleAuto                      // auto: end if section reaches last line, else start
)

// CollapsibleConfig configures threshold-based collapsible code blocks.
type CollapsibleConfig struct {
	LineThreshold        int
	PreviewLines         int
	DefaultCollapsed     bool
	PreserveIndent       bool
	Style                CollapseStyle
	ExpandButtonText     string
	CollapseButtonText   string
	ExpandedAnnouncement string // screen reader
	CollapsedAnnouncement string // screen reader
}

// StyleValue represents a CSS variable override that can be universal or per-theme.
type StyleValue struct {
	Value string // used when both themes share the same value
	Dark  string // dark theme override
	Light string // light theme override
}

// IsThemed reports whether this value has per-theme overrides.
func (sv StyleValue) IsThemed() bool {
	return sv.Dark != "" || sv.Light != ""
}

// LightValue returns the value to use for the light theme.
func (sv StyleValue) LightValue() string {
	if sv.IsThemed() {
		return sv.Light
	}
	return sv.Value
}

// DarkValue returns the value to use for the dark theme.
func (sv StyleValue) DarkValue() string {
	if sv.IsThemed() {
		return sv.Dark
	}
	return sv.Value
}

// Assets holds CSS and JS output with content-hashed filenames.
type Assets struct {
	CSS AssetFile
	JS  AssetFile
}

// AssetFile holds content and its hashed filename.
type AssetFile struct {
	Content  string
	Hash     string // 8-char hex (FNV-1a)
	Filename string // e.g. "kazari-a1b2c3d4.css"
}

func makeAssetFile(content, ext string) AssetFile {
	h := fnv.New32a()
	h.Write([]byte(content))
	hash := fmt.Sprintf("%08x", h.Sum32())
	return AssetFile{
		Content:  content,
		Hash:     hash,
		Filename: fmt.Sprintf("kazari-%s.%s", hash, ext),
	}
}
