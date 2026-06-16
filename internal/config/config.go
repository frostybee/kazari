package config

import (
	"strings"

	"github.com/frostybee/kazari/internal/locale"
)

// DarkModeKind identifies the dark mode CSS strategy.
type DarkModeKind int

const (
	DarkModeSelectorKind   DarkModeKind = iota // CSS class selector
	DarkModeMediaQueryKind                     // prefers-color-scheme
	DarkModeBothKind                           // media query + class toggle
)

// DarkModeConfig holds the dark mode strategy settings.
type DarkModeConfig struct {
	Kind     DarkModeKind
	Selector string // used by Selector and Both modes (e.g. ".dark")
}

// Frame constants matching the public Frame enum.
const (
	FrameAuto     = 0
	FrameCode     = 1
	FrameTerminal = 2
	FrameNone     = 3
)

// Terminal dot style constants matching the public TerminalDotStyle enum.
const (
	DotsColored = 0
	DotsMinimal = 1
)

// Language icon mode constants matching the public LangIconMode enum.
const (
	LangIconNone    = 0
	LangIconOnly    = 1
	LangIconAndText = 2
)

// MarkerType identifies the kind of line or inline marker.
type MarkerType int

// MarkerNone indicates a segment with no line/inline marker (used for link-only annotations).
const MarkerNone MarkerType = -1

// Priority order: mark(0) < del(1) < ins(2). Higher value wins in overlap resolution.
const (
	MarkerMark MarkerType = iota // highlight (default, lowest priority)
	MarkerDel                    // deleted
	MarkerIns                    // inserted (highest priority)
)

// LineRange is an inclusive 1-based line range.
type LineRange struct {
	Start int
	End   int
}

// LineMarker represents a line-level marker with optional label.
type LineMarker struct {
	Type  MarkerType
	Lines []LineRange
	Label string
}

// InlineMarker represents an inline text marker.
type InlineMarker struct {
	Type    MarkerType
	Text    string
	IsRegex bool
}

// LinkAnnotation records the position and URL of an extracted inline link.
type LinkAnnotation struct {
	Start int    // character offset in cleaned text (inclusive)
	End   int    // character offset in cleaned text (exclusive)
	URL   string
}

// MarkerBGColors maps marker types to their CSS rgba background values.
var MarkerBGColors = map[MarkerType]string{
	MarkerMark: "rgba(255,200,0,0.12)",
	MarkerIns:  "rgba(46,160,67,0.12)",
	MarkerDel:  "rgba(248,81,73,0.12)",
}

// MarkerEffectiveBGs holds opaque hex backgrounds after compositing marker RGBA on editor BG.
type MarkerEffectiveBGs struct {
	Mark string
	Ins  string
	Del  string
}

// BG returns the effective background for a given marker type.
func (m *MarkerEffectiveBGs) BG(mt MarkerType) string {
	switch mt {
	case MarkerIns:
		return m.Ins
	case MarkerDel:
		return m.Del
	default:
		return m.Mark
	}
}

// MergedToken holds both light and dark colors for a single token.
type MergedToken struct {
	Content    string
	LightColor string
	DarkColor  string
	LightBG    string
	DarkBG     string
	FontStyle  int
}

// TokenLine represents one line of merged tokens.
type TokenLine struct {
	Tokens []MergedToken
}

// BlockDefaults holds default block rendering properties.
type BlockDefaults struct {
	Wrap           bool
	PreserveIndent bool
	HangingIndent  int
	LineNumbers    bool
	Frame          int // FrameAuto, FrameCode, FrameTerminal, FrameNone
}

// CollapseStyle identifies the visual style for range-based collapsible sections.
type CollapseStyle int

const (
	CollapseGithub          CollapseStyle = iota // default — one-way expand, summary disappears
	CollapseCollapsibleStart                     // re-collapsible, summary above content
	CollapseCollapsibleEnd                       // re-collapsible, summary below content
	CollapseCollapsibleAuto                      // auto: end if section reaches last line, else start
)

// CollapsibleConfig configures threshold-based collapsible code blocks.
type CollapsibleConfig struct {
	LineThreshold         int
	PreviewLines          int
	DefaultCollapsed      bool
	PreserveIndent        bool
	Style                 CollapseStyle
	ExpandButtonText      string
	CollapseButtonText    string
	ExpandedAnnouncement  string
	CollapsedAnnouncement string
}

// CollapseSpec holds per-block collapse directives from meta string or Options.
type CollapseSpec struct {
	Enabled   bool
	Disabled  bool
	Ranges    []LineRange
	Style     *CollapseStyle // nil = use engine default
	Threshold *int           // nil = use engine default; overrides CollapsibleConfig.LineThreshold
}

// CollapseRange is a validated, render-ready collapse range with pre-computed metadata.
type CollapseRange struct {
	Start     int           // 1-based inclusive
	End       int           // 1-based inclusive
	LineCount int           // End - Start + 1
	MinIndent int           // minimum indentation in spaces (for --kz-indent)
	Style     CollapseStyle // resolved style (auto already resolved)
}

// PreviewSegment represents a contiguous range of visible lines in threshold preview.
type PreviewSegment struct {
	Start int // 1-based inclusive
	End   int // 1-based inclusive
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

// Config holds the fully resolved engine configuration.
type Config struct {
	LightTheme         string
	DarkTheme          string
	CopyButton         bool
	FullscreenButton   bool
	WrapButton         bool
	ThemeToggle        bool
	FrameDetection     bool
	FileNameExtraction bool
	LanguageBadge      bool
	StyleReset         bool
	ThemedScrollbars   bool
	ThemedSelection    bool
	ContentExclusion   bool
	Collapsible        *CollapsibleConfig
	Defaults           BlockDefaults
	LanguageDefaults   map[string]BlockDefaults
	DarkMode           DarkModeConfig
	TabWidth           int
	MinContrast        float64
	LightEditorBG      string
	DarkEditorBG       string
	LightMarkerBGs     *MarkerEffectiveBGs
	DarkMarkerBGs      *MarkerEffectiveBGs
	Minify             bool
	CascadeLayer       string
	LanguageAliases    map[string]string
	CodeGroups         bool
	TerminalDotStyle           int // DotsColored or DotsMinimal
	TerminalCommentStripping   bool
	MermaidPassThrough         bool
	DataLineCount              bool
	ThemeCSSRoot               string
	StyleOverrides             map[string]StyleValue
	Locale                     string
	UIStringOverrides          map[string]string
	UIStrings                  *locale.UIStrings
	LangIconMode               int
	Links                      bool
	FileIcons                  bool
	FileIconResolver           func(string) string
	WarningHandler             func(string)
}

// DefaultConfig returns the engine configuration with all documented defaults.
func DefaultConfig() *Config {
	return &Config{
		LightTheme:         "github-light",
		DarkTheme:          "github-dark",
		CopyButton:         true,
		FullscreenButton:   true,
		WrapButton:         true,
		FrameDetection:     true,
		FileNameExtraction: true,
		LanguageBadge:      true,
		StyleReset:         true,
		ThemedScrollbars:   true,
		ThemedSelection:    false,
		ContentExclusion:   false,
		Collapsible:        nil,
		Defaults: BlockDefaults{
			Wrap:           false,
			PreserveIndent: true,
			HangingIndent:  0,
			LineNumbers:    false,
			Frame:          FrameAuto,
		},
		LanguageDefaults: nil,
		DarkMode: DarkModeConfig{
			Kind:     DarkModeSelectorKind,
			Selector: ".dark",
		},
		TabWidth:        2,
		MinContrast:     5.5,
		Minify:          true,
		CascadeLayer:     "kazari",
		LanguageAliases:  nil,
		TerminalDotStyle:         DotsColored,
		TerminalCommentStripping: true,
		MermaidPassThrough:       true,
		DataLineCount:            true,
		ThemeCSSRoot:             ":root",
		Locale:                   "en-US",
		FileIcons:                true,
	}
}

// BlockOptions represents per-block overrides (from Options or meta string).
type BlockOptions struct {
	Lang            string
	Title           string
	Theme           string
	Frame           *int  // nil = use default
	LineNumbers     *bool // nil = use default
	StartLineNumber *int  // nil = use default (1)
	Wrap            *bool // nil = use default
	PreserveIndent  *bool // nil = use default
	HangingIndent   *int  // nil = use default
}

// ResolvedBlock is the final merged config for rendering a single code block.
type ResolvedBlock struct {
	Lang            string
	Title           string
	Theme           string
	DiffLang        string
	Frame           int
	LineNumbers     bool
	StartLineNumber int
	Wrap            bool
	PreserveIndent  bool
	HangingIndent   int
	RawCode       string // code for copy button (post file-name extraction)
	LineMarkers   []LineMarker
	InlineMarkers []InlineMarker
	Links         [][]LinkAnnotation // per-line link annotations (indexed by line)
	FocusLines    []LineRange
	// Per-block theme override state (populated when Theme is set)
	ThemeOverrideStyle string              // inline --kz-ovl-*/--kz-ovd-* declarations for the wrapper
	LightEditorBG      string              // editor BG for unmarked-line contrast (override light)
	DarkEditorBG       string              // editor BG for unmarked-line contrast (override dark)
	LightMarkerBGs     *MarkerEffectiveBGs // marker compositing against the override light BG
	DarkMarkerBGs      *MarkerEffectiveBGs // marker compositing against the override dark BG
	// Collapse state (populated by collapsible.ResolveCollapse)
	CollapseThreshold  bool
	CollapseRanges     []CollapseRange
	CollapseSegments   []PreviewSegment
	CollapseBeyondCap  int // marked lines beyond 2× preview cap (for badge)
	CollapseConfig     *CollapsibleConfig
}

// Resolve applies the config cascade for a specific block:
// engine defaults -> Defaults -> LanguageDefaults -> per-block options.
func (c *Config) Resolve(lang string, blockOpts *BlockOptions) *ResolvedBlock {
	resolved := &ResolvedBlock{
		Lang:            lang,
		Frame:           c.Defaults.Frame,
		LineNumbers:     c.Defaults.LineNumbers,
		StartLineNumber: 1,
		Wrap:            c.Defaults.Wrap,
		PreserveIndent:  c.Defaults.PreserveIndent,
		HangingIndent:   c.Defaults.HangingIndent,
	}

	// Apply language defaults (comma-separated key matching).
	if c.LanguageDefaults != nil {
		for key, langDef := range c.LanguageDefaults {
			if matchesLanguageKey(key, lang) {
				resolved.Frame = langDef.Frame
				resolved.LineNumbers = langDef.LineNumbers
				resolved.Wrap = langDef.Wrap
				resolved.PreserveIndent = langDef.PreserveIndent
				resolved.HangingIndent = langDef.HangingIndent
				break
			}
		}
	}

	// Apply per-block overrides.
	if blockOpts != nil {
		if blockOpts.Lang != "" {
			resolved.Lang = blockOpts.Lang
		}
		if blockOpts.Title != "" {
			resolved.Title = blockOpts.Title
		}
		if blockOpts.Theme != "" {
			resolved.Theme = blockOpts.Theme
		}
		if blockOpts.Frame != nil {
			resolved.Frame = *blockOpts.Frame
		}
		if blockOpts.LineNumbers != nil {
			resolved.LineNumbers = *blockOpts.LineNumbers
		}
		if blockOpts.StartLineNumber != nil {
			resolved.StartLineNumber = *blockOpts.StartLineNumber
		}
		if blockOpts.Wrap != nil {
			resolved.Wrap = *blockOpts.Wrap
		}
		if blockOpts.PreserveIndent != nil {
			resolved.PreserveIndent = *blockOpts.PreserveIndent
		}
		if blockOpts.HangingIndent != nil {
			resolved.HangingIndent = *blockOpts.HangingIndent
		}
	}

	return resolved
}

// ResolveLanguage applies language aliases and returns the canonical language name.
func (c *Config) ResolveLanguage(lang string) string {
	lower := strings.ToLower(lang)
	if c.LanguageAliases != nil {
		if canonical, ok := c.LanguageAliases[lower]; ok {
			return canonical
		}
	}
	return lower
}

// matchesLanguageKey checks if lang matches a comma-separated key (e.g. "bash,sh,zsh").
func matchesLanguageKey(key, lang string) bool {
	lower := strings.ToLower(lang)
	for _, k := range strings.Split(key, ",") {
		if strings.TrimSpace(strings.ToLower(k)) == lower {
			return true
		}
	}
	return false
}
