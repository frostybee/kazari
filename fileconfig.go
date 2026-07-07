package kazari

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/frostybee/edict"
	"gopkg.in/yaml.v3"
)

// FileConfig represents a Kazari configuration file (YAML or JSON).
// All fields use pointer types so omitted keys are distinguishable from zero values.
type FileConfig struct {
	Themes                   *ThemesFileConfig                  `yaml:"themes" json:"themes"`
	DarkMode                 *DarkModeFileConfig                `yaml:"darkMode" json:"darkMode"`
	CopyButton               *bool                             `yaml:"copyButton" json:"copyButton"`
	FullscreenButton         *bool                             `yaml:"fullscreenButton" json:"fullscreenButton"`
	WrapButton               *bool                             `yaml:"wrapButton" json:"wrapButton"`
	ThemeToggleButton        *bool                             `yaml:"themeToggleButton" json:"themeToggleButton"`
	LineNumbers              *bool                             `yaml:"lineNumbers" json:"lineNumbers"`
	FrameDetection           *bool                             `yaml:"frameDetection" json:"frameDetection"`
	FileNameExtraction       *bool                             `yaml:"fileNameExtraction" json:"fileNameExtraction"`
	LanguageBadge            *bool                             `yaml:"languageBadge" json:"languageBadge"`
	ThemedScrollbars         *bool                             `yaml:"themedScrollbars" json:"themedScrollbars"`
	ThemedSelection          *bool                             `yaml:"themedSelection" json:"themedSelection"`
	ContentExclusion         *bool                             `yaml:"contentExclusion" json:"contentExclusion"`
	MermaidPassThrough       *bool                             `yaml:"mermaidPassThrough" json:"mermaidPassThrough"`
	TerminalCommentStripping *bool                             `yaml:"terminalCommentStripping" json:"terminalCommentStripping"`
	DataLineCount            *bool                             `yaml:"dataLineCount" json:"dataLineCount"`
	FileIcons                *bool                             `yaml:"fileIcons" json:"fileIcons"`
	Links                    *bool                             `yaml:"links" json:"links"`
	StyleReset               *bool                             `yaml:"styleReset" json:"styleReset"`
	Minify                   *bool                             `yaml:"minify" json:"minify"`
	TabWidth                 *int                              `yaml:"tabWidth" json:"tabWidth"`
	MinContrast              *float64                          `yaml:"minContrast" json:"minContrast"`
	CascadeLayer             *string                           `yaml:"cascadeLayer" json:"cascadeLayer"`
	ThemeCSSRoot             *string                           `yaml:"themeCSSRoot" json:"themeCSSRoot"`
	Locale                   *string                           `yaml:"locale" json:"locale"`
	TerminalDotStyle         *string                           `yaml:"terminalDotStyle" json:"terminalDotStyle"`
	LanguageIconMode         *string                           `yaml:"languageIconMode" json:"languageIconMode"`
	Defaults                 *BlockDefaultsFileConfig          `yaml:"defaults" json:"defaults"`
	LanguageDefaults         map[string]BlockDefaultsFileConfig `yaml:"languageDefaults" json:"languageDefaults"`
	Collapsible              *CollapsibleFileConfig            `yaml:"collapsible" json:"collapsible"`
	LanguageAliases          map[string]string                 `yaml:"languageAliases" json:"languageAliases"`
	UIStrings                map[string]string                 `yaml:"uiStrings" json:"uiStrings"`
	OutputPanel              *bool                             `yaml:"outputPanel" json:"outputPanel"`
	OutputDefaultCollapsed   *bool                             `yaml:"outputDefaultCollapsed" json:"outputDefaultCollapsed"`
	OutputSeparator          *string                           `yaml:"outputSeparator" json:"outputSeparator"`
	StyleOverrides           map[string]any                    `yaml:"styleOverrides" json:"styleOverrides"`
}

type ThemesFileConfig struct {
	Light string `yaml:"light" json:"light"`
	Dark  string `yaml:"dark" json:"dark"`
}

type DarkModeFileConfig struct {
	Kind     string `yaml:"kind" json:"kind"`
	Selector string `yaml:"selector" json:"selector"`
}

type BlockDefaultsFileConfig struct {
	Wrap           *bool   `yaml:"wrap" json:"wrap"`
	PreserveIndent *bool   `yaml:"preserveIndent" json:"preserveIndent"`
	HangingIndent  *int    `yaml:"hangingIndent" json:"hangingIndent"`
	LineNumbers    *bool   `yaml:"lineNumbers" json:"lineNumbers"`
	Frame          *string `yaml:"frame" json:"frame"`
}

type CollapsibleFileConfig struct {
	LineThreshold         *int    `yaml:"lineThreshold" json:"lineThreshold"`
	PreviewLines          *int    `yaml:"previewLines" json:"previewLines"`
	DefaultCollapsed      *bool   `yaml:"defaultCollapsed" json:"defaultCollapsed"`
	PreserveIndent        *bool   `yaml:"preserveIndent" json:"preserveIndent"`
	Style                 *string `yaml:"style" json:"style"`
	ExpandButtonText      *string `yaml:"expandButtonText" json:"expandButtonText"`
	CollapseButtonText    *string `yaml:"collapseButtonText" json:"collapseButtonText"`
	ExpandedAnnouncement  *string `yaml:"expandedAnnouncement" json:"expandedAnnouncement"`
	CollapsedAnnouncement *string `yaml:"collapsedAnnouncement" json:"collapsedAnnouncement"`
}

// --- Enum parsers ---

func parseFrame(s string) (Frame, error) {
	switch s {
	case "auto":
		return FrameAuto, nil
	case "code":
		return FrameCode, nil
	case "terminal":
		return FrameTerminal, nil
	case "none":
		return FrameNone, nil
	default:
		return 0, fmt.Errorf("kazari: invalid frame %q, must be one of: auto, code, terminal, none", s)
	}
}

func parseCollapseStyle(s string) (CollapseStyle, error) {
	switch s {
	case "github":
		return CollapseGithub, nil
	case "collapsibleStart":
		return CollapseCollapsibleStart, nil
	case "collapsibleEnd":
		return CollapseCollapsibleEnd, nil
	case "collapsibleAuto":
		return CollapseCollapsibleAuto, nil
	default:
		return 0, fmt.Errorf("kazari: invalid collapsible style %q, must be one of: github, collapsibleStart, collapsibleEnd, collapsibleAuto", s)
	}
}

func parseTerminalDotStyleStr(s string) (TerminalDotStyle, error) {
	switch s {
	case "colored":
		return DotsColored, nil
	case "minimal":
		return DotsMinimal, nil
	default:
		return 0, fmt.Errorf("kazari: invalid terminalDotStyle %q, must be one of: colored, minimal", s)
	}
}

func parseLangIconModeStr(s string) (LangIconMode, error) {
	switch s {
	case "none":
		return LangIconNone, nil
	case "iconOnly":
		return LangIconOnly, nil
	case "iconAndText":
		return LangIconAndText, nil
	default:
		return 0, fmt.Errorf("kazari: invalid languageIconMode %q, must be one of: none, iconOnly, iconAndText", s)
	}
}

func parseDarkModeConfig(fc *DarkModeFileConfig) (DarkMode, error) {
	switch fc.Kind {
	case "selector":
		return SelectorMode(fc.Selector), nil
	case "mediaQuery":
		return MediaQueryMode(), nil
	case "both":
		return BothMode(fc.Selector), nil
	default:
		return nil, fmt.Errorf("kazari: invalid darkMode.kind %q, must be one of: selector, mediaQuery, both", fc.Kind)
	}
}

// --- Validation ---

func validateFileConfig(fc *FileConfig) error {
	v := edict.New()

	if fc.DarkMode != nil {
		v.Field("darkMode.kind").Value(fc.DarkMode.Kind).OneOf("selector", "mediaQuery", "both")
		v.Field("darkMode.selector").Value(fc.DarkMode.Selector).RequiredIf("darkMode.kind", "selector", "both")
	}

	if fc.TabWidth != nil {
		v.Field("tabWidth").Value(*fc.TabWidth).OmitNil().IntMin(1)
	}

	if fc.MinContrast != nil {
		v.Field("minContrast").Value(*fc.MinContrast).OmitNil().FloatRange(0, 21)
	}

	if fc.TerminalDotStyle != nil {
		v.Field("terminalDotStyle").Value(*fc.TerminalDotStyle).OneOf("colored", "minimal")
	}

	if fc.LanguageIconMode != nil {
		v.Field("languageIconMode").Value(*fc.LanguageIconMode).OneOf("none", "iconOnly", "iconAndText")
	}

	if fc.Defaults != nil {
		if fc.Defaults.Frame != nil {
			v.Field("defaults.frame").Value(*fc.Defaults.Frame).OneOf("auto", "code", "terminal", "none")
		}
		if fc.Defaults.HangingIndent != nil {
			v.Field("defaults.hangingIndent").Value(*fc.Defaults.HangingIndent).OmitNil().IntMin(0)
		}
	}

	for key, ld := range fc.LanguageDefaults {
		if ld.Frame != nil {
			v.Field(fmt.Sprintf("languageDefaults.%s.frame", key)).Value(*ld.Frame).OneOf("auto", "code", "terminal", "none")
		}
		if ld.HangingIndent != nil {
			v.Field(fmt.Sprintf("languageDefaults.%s.hangingIndent", key)).Value(*ld.HangingIndent).OmitNil().IntMin(0)
		}
	}

	if fc.Collapsible != nil {
		if fc.Collapsible.LineThreshold != nil {
			v.Field("collapsible.lineThreshold").Value(*fc.Collapsible.LineThreshold).OmitNil().IntMin(1)
		}
		if fc.Collapsible.PreviewLines != nil {
			v.Field("collapsible.previewLines").Value(*fc.Collapsible.PreviewLines).OmitNil().IntMin(1)
		}
		if fc.Collapsible.Style != nil {
			v.Field("collapsible.style").Value(*fc.Collapsible.Style).OneOf("github", "collapsibleStart", "collapsibleEnd", "collapsibleAuto")
		}
	}

	return v.Validate()
}

// --- ParseConfig ---

// ParseConfig parses config bytes in the given format ("yaml" or "json")
// and validates the result.
func ParseConfig(data []byte, format string) (*FileConfig, error) {
	fc := &FileConfig{}

	switch format {
	case "yaml":
		dec := yaml.NewDecoder(bytes.NewReader(data))
		dec.KnownFields(true)
		if err := dec.Decode(fc); err != nil {
			return nil, fmt.Errorf("kazari: parsing YAML config: %w", err)
		}
	case "json":
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields()
		if err := dec.Decode(fc); err != nil {
			return nil, fmt.Errorf("kazari: parsing JSON config: %w", err)
		}
	default:
		return nil, fmt.Errorf("kazari: unsupported config format %q, must be yaml or json", format)
	}

	if err := validateFileConfig(fc); err != nil {
		return nil, err
	}

	return fc, nil
}
