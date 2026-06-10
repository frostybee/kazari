package config

import "strings"

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

// MarkerType identifies the kind of line or inline marker.
type MarkerType int

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
	Type MarkerType
	Text string
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
	Enabled  bool
	Disabled bool
	Ranges   []LineRange
	Style    *CollapseStyle // nil = use engine default
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

// Config holds the fully resolved engine configuration.
type Config struct {
	LightTheme         string
	DarkTheme          string
	CopyButton         bool
	FullscreenButton   bool
	LineNumbers        bool
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
	Minify             bool
	CascadeLayer       string
	LanguageAliases    map[string]string
	CodeGroups         bool
}

// DefaultConfig returns the engine configuration with all documented defaults.
func DefaultConfig() *Config {
	return &Config{
		LightTheme:         "github-light",
		DarkTheme:          "github-dark",
		CopyButton:         true,
		FullscreenButton:   true,
		LineNumbers:        false,
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
		CascadeLayer:    "kazari",
		LanguageAliases: nil,
	}
}

// BlockOptions represents per-block overrides (from Options or meta string).
type BlockOptions struct {
	Lang            string
	Title           string
	Frame           *int  // nil = use default
	LineNumbers     *bool // nil = use default
	StartLineNumber *int  // nil = use default (1)
	Wrap            *bool // nil = use default
}

// ResolvedBlock is the final merged config for rendering a single code block.
type ResolvedBlock struct {
	Lang            string
	Title           string
	Frame           int
	LineNumbers     bool
	StartLineNumber int
	Wrap            bool
	PreserveIndent  bool
	HangingIndent   int
	RawCode       string // code for copy button (post file-name extraction)
	LineMarkers   []LineMarker
	InlineMarkers []InlineMarker
	FocusLines    []LineRange
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
