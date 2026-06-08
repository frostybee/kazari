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

// BlockDefaults holds default block rendering properties.
type BlockDefaults struct {
	Wrap           bool
	PreserveIndent bool
	HangingIndent  int
	LineNumbers    bool
	Frame          int // FrameAuto, FrameCode, FrameTerminal, FrameNone
}

// CollapsibleConfig configures threshold-based collapsible code blocks.
type CollapsibleConfig struct {
	LineThreshold         int
	PreviewLines          int
	DefaultCollapsed      bool
	PreserveIndent        bool
	ExpandButtonText      string
	CollapseButtonText    string
	ExpandedAnnouncement  string
	CollapsedAnnouncement string
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
	StartLineNumber int
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
	RawCode         string // code for copy button (post file-name extraction)
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
		if blockOpts.StartLineNumber > 0 {
			resolved.StartLineNumber = blockOpts.StartLineNumber
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
