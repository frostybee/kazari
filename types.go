package kazari

import (
	"fmt"
	"hash/fnv"
)

// Highlighter abstracts syntax highlighting. The default implementation wraps Nuri,
// but any implementation (including mocks) can be used.
type Highlighter interface {
	Tokenize(code, lang, theme string) ([][]Token, error)
	GetThemeColors(theme string) (ThemeInfo, error)
	GetLoadedLanguages() []string
}

// ThemeInfo holds colors extracted from a syntax theme for CSS variable generation.
type ThemeInfo struct {
	FG           string // default foreground (e.g. "#24292f")
	BG           string // default background (e.g. "#ffffff")
	SelectionBG  string // editor.selectionBackground
	LineNumberFG string // editorLineNumber.foreground
}

// Token represents a single colored token from the highlighter.
type Token struct {
	Content   string
	Color     string // foreground hex color (e.g. "#cf222e")
	BgColor   string // background hex color (usually empty)
	FontStyle int    // bitmask: Italic=1, Bold=2, Underline=4, Strikethrough=8
}

// Font style bitmask constants.
const (
	FontStyleNone          = 0
	FontStyleItalic        = 1
	FontStyleBold          = 2
	FontStyleUnderline     = 4
	FontStyleStrikethrough = 8
)

// Frame determines the code block's visual frame type.
type Frame int

const (
	FrameAuto     Frame = iota // auto-detect from language
	FrameCode                  // editor frame
	FrameTerminal              // terminal frame
	FrameNone                  // no frame
)

// Range represents an inclusive 1-based line range.
type Range struct {
	Start int
	End   int
}

// MarkerType identifies the kind of line or inline marker.
type MarkerType int

const (
	MarkerMark MarkerType = iota // highlight (default)
	MarkerIns                    // inserted
	MarkerDel                    // deleted
)

// LineMarker represents a line-level marker with optional label.
type LineMarker struct {
	Type  MarkerType
	Lines []Range
	Label string // empty for unlabeled
}

// InlineMarker represents an inline text marker.
type InlineMarker struct {
	Type MarkerType
	Text string
}

// Options is the per-block configuration passed to Render or derived from meta string.
type Options struct {
	Lang            string
	Title           string
	Frame           *Frame // nil = use default
	LineNumbers     *bool  // nil = use default
	StartLineNumber int
	Wrap            *bool // nil = use default
	LineMarkers     []LineMarker
	InlineMarkers   []InlineMarker
	FocusLines      []Range
	Collapse        *CollapseOptions
}

// CollapseOptions holds per-block collapse configuration from meta string.
type CollapseOptions struct {
	Enabled  bool    // true = force threshold-based collapse
	Disabled bool    // true = force no collapse (nocollapse)
	Ranges   []Range // specific ranges to collapse
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

// CollapsibleConfig configures threshold-based collapsible code blocks.
type CollapsibleConfig struct {
	LineThreshold        int
	PreviewLines         int
	DefaultCollapsed     bool
	PreserveIndent       bool
	ExpandButtonText     string
	CollapseButtonText   string
	ExpandedAnnouncement string // screen reader
	CollapsedAnnouncement string // screen reader
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
