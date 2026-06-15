package kazari

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/valiant"
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
	StyleReset               *bool                             `yaml:"styleReset" json:"styleReset"`
	Minify                   *bool                             `yaml:"minify" json:"minify"`
	TabWidth                 *int                              `yaml:"tabWidth" json:"tabWidth"`
	MinContrast              *float64                          `yaml:"minContrast" json:"minContrast"`
	CascadeLayer             *string                           `yaml:"cascadeLayer" json:"cascadeLayer"`
	ThemeCSSRoot             *string                           `yaml:"themeCSSRoot" json:"themeCSSRoot"`
	Locale                   *string                           `yaml:"locale" json:"locale"`
	TerminalDotStyle         *string                           `yaml:"terminalDotStyle" json:"terminalDotStyle"`
	Defaults                 *BlockDefaultsFileConfig          `yaml:"defaults" json:"defaults"`
	LanguageDefaults         map[string]BlockDefaultsFileConfig `yaml:"languageDefaults" json:"languageDefaults"`
	Collapsible              *CollapsibleFileConfig            `yaml:"collapsible" json:"collapsible"`
	LanguageAliases          map[string]string                 `yaml:"languageAliases" json:"languageAliases"`
	UIStrings                map[string]string                 `yaml:"uiStrings" json:"uiStrings"`
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
	v := valiant.New()

	if fc.DarkMode != nil {
		v.Field("darkMode.kind").Value(fc.DarkMode.Kind).OneOf("selector", "mediaQuery", "both")
		v.Field("darkMode.selector").Value(fc.DarkMode.Selector).RequiredIf("darkMode.kind", "selector", "both")
	}

	if fc.TabWidth != nil {
		v.Field("tabWidth").Value(*fc.TabWidth).OmitNil().IntMin(1)
	}

	if fc.MinContrast != nil {
		v.Field("minContrast").Value(*fc.MinContrast).OmitNil().FloatMin(0)
	}

	if fc.TerminalDotStyle != nil {
		v.Field("terminalDotStyle").Value(*fc.TerminalDotStyle).OneOf("colored", "minimal")
	}

	if fc.Defaults != nil && fc.Defaults.Frame != nil {
		v.Field("defaults.frame").Value(*fc.Defaults.Frame).OneOf("auto", "code", "terminal", "none")
	}

	for key, ld := range fc.LanguageDefaults {
		if ld.Frame != nil {
			v.Field(fmt.Sprintf("languageDefaults.%s.frame", key)).Value(*ld.Frame).OneOf("auto", "code", "terminal", "none")
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

// --- FileConfigToOptions ---

// FileConfigToOptions converts a FileConfig into a slice of functional options.
func FileConfigToOptions(fc *FileConfig) ([]Option, error) {
	var opts []Option

	if fc.Themes != nil {
		opts = append(opts, WithThemes(fc.Themes.Light, fc.Themes.Dark))
	}

	if fc.DarkMode != nil {
		dm, err := parseDarkModeConfig(fc.DarkMode)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithDarkMode(dm))
	}

	if fc.CopyButton != nil {
		opts = append(opts, WithCopyButton(*fc.CopyButton))
	}
	if fc.FullscreenButton != nil {
		opts = append(opts, WithFullscreenButton(*fc.FullscreenButton))
	}
	if fc.WrapButton != nil {
		opts = append(opts, WithWrapButton(*fc.WrapButton))
	}
	if fc.LineNumbers != nil {
		opts = append(opts, WithLineNumbers(*fc.LineNumbers))
	}
	if fc.FrameDetection != nil {
		opts = append(opts, WithFrameDetection(*fc.FrameDetection))
	}
	if fc.FileNameExtraction != nil {
		opts = append(opts, WithFileNameExtraction(*fc.FileNameExtraction))
	}
	if fc.LanguageBadge != nil {
		opts = append(opts, WithLanguageBadge(*fc.LanguageBadge))
	}
	if fc.ThemedScrollbars != nil {
		opts = append(opts, WithThemedScrollbars(*fc.ThemedScrollbars))
	}
	if fc.ThemedSelection != nil {
		opts = append(opts, WithThemedSelectionColors(*fc.ThemedSelection))
	}
	if fc.ContentExclusion != nil {
		opts = append(opts, WithContentExclusion(*fc.ContentExclusion))
	}
	if fc.MermaidPassThrough != nil {
		opts = append(opts, WithMermaidPassThrough(*fc.MermaidPassThrough))
	}
	if fc.TerminalCommentStripping != nil {
		opts = append(opts, WithTerminalCommentStripping(*fc.TerminalCommentStripping))
	}
	if fc.DataLineCount != nil {
		opts = append(opts, WithDataLineCount(*fc.DataLineCount))
	}
	if fc.FileIcons != nil {
		opts = append(opts, WithFileIcons(*fc.FileIcons))
	}
	if fc.StyleReset != nil {
		opts = append(opts, WithStyleReset(*fc.StyleReset))
	}
	if fc.Minify != nil {
		opts = append(opts, WithMinify(*fc.Minify))
	}

	if fc.TabWidth != nil {
		opts = append(opts, WithTabWidth(*fc.TabWidth))
	}
	if fc.MinContrast != nil {
		opts = append(opts, WithMinSyntaxHighlightingColorContrast(*fc.MinContrast))
	}
	if fc.CascadeLayer != nil {
		opts = append(opts, WithCascadeLayer(*fc.CascadeLayer))
	}
	if fc.ThemeCSSRoot != nil {
		opts = append(opts, WithThemeCSSRoot(*fc.ThemeCSSRoot))
	}
	if fc.Locale != nil {
		opts = append(opts, WithLocale(*fc.Locale))
	}

	if fc.TerminalDotStyle != nil {
		ds, err := parseTerminalDotStyleStr(*fc.TerminalDotStyle)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithTerminalDotStyle(ds))
	}

	if fc.Defaults != nil {
		d := fc.Defaults
		opts = append(opts, func(b *engineBuilder) {
			if d.Wrap != nil {
				b.cfg.Defaults.Wrap = *d.Wrap
			}
			if d.PreserveIndent != nil {
				b.cfg.Defaults.PreserveIndent = *d.PreserveIndent
			}
			if d.HangingIndent != nil {
				b.cfg.Defaults.HangingIndent = *d.HangingIndent
			}
			if d.LineNumbers != nil {
				b.cfg.Defaults.LineNumbers = *d.LineNumbers
			}
			if d.Frame != nil {
				f, err := parseFrame(*d.Frame)
				if err == nil {
					b.cfg.Defaults.Frame = int(f)
				}
			}
		})
	}

	if len(fc.LanguageDefaults) > 0 {
		for compositeKey, ld := range fc.LanguageDefaults {
			ld := ld
			for _, lang := range splitCommaKey(compositeKey) {
				lang := lang
				opts = append(opts, func(b *engineBuilder) {
					if b.cfg.LanguageDefaults == nil {
						b.cfg.LanguageDefaults = make(map[string]config.BlockDefaults)
					}
					existing := b.cfg.LanguageDefaults[lang]
					if ld.Wrap != nil {
						existing.Wrap = *ld.Wrap
					}
					if ld.PreserveIndent != nil {
						existing.PreserveIndent = *ld.PreserveIndent
					}
					if ld.HangingIndent != nil {
						existing.HangingIndent = *ld.HangingIndent
					}
					if ld.LineNumbers != nil {
						existing.LineNumbers = *ld.LineNumbers
					}
					if ld.Frame != nil {
						f, err := parseFrame(*ld.Frame)
						if err == nil {
							existing.Frame = int(f)
						}
					}
					b.cfg.LanguageDefaults[lang] = existing
				})
			}
		}
	}

	if fc.Collapsible != nil {
		c := fc.Collapsible
		opts = append(opts, func(b *engineBuilder) {
			if b.cfg.Collapsible == nil {
				b.cfg.Collapsible = &config.CollapsibleConfig{
					LineThreshold:    15,
					PreviewLines:     5,
					DefaultCollapsed: true,
					PreserveIndent:   true,
				}
			}
			if c.LineThreshold != nil {
				b.cfg.Collapsible.LineThreshold = *c.LineThreshold
			}
			if c.PreviewLines != nil {
				b.cfg.Collapsible.PreviewLines = *c.PreviewLines
			}
			if c.DefaultCollapsed != nil {
				b.cfg.Collapsible.DefaultCollapsed = *c.DefaultCollapsed
			}
			if c.PreserveIndent != nil {
				b.cfg.Collapsible.PreserveIndent = *c.PreserveIndent
			}
			if c.Style != nil {
				s, err := parseCollapseStyle(*c.Style)
				if err == nil {
					b.cfg.Collapsible.Style = config.CollapseStyle(s)
				}
			}
			if c.ExpandButtonText != nil {
				b.cfg.Collapsible.ExpandButtonText = *c.ExpandButtonText
			}
			if c.CollapseButtonText != nil {
				b.cfg.Collapsible.CollapseButtonText = *c.CollapseButtonText
			}
			if c.ExpandedAnnouncement != nil {
				b.cfg.Collapsible.ExpandedAnnouncement = *c.ExpandedAnnouncement
			}
			if c.CollapsedAnnouncement != nil {
				b.cfg.Collapsible.CollapsedAnnouncement = *c.CollapsedAnnouncement
			}
		})
	}

	if len(fc.LanguageAliases) > 0 {
		opts = append(opts, WithLanguageAliases(fc.LanguageAliases))
	}

	if len(fc.UIStrings) > 0 {
		opts = append(opts, WithUIStrings(fc.UIStrings))
	}

	if len(fc.StyleOverrides) > 0 {
		overrides, err := parseStyleOverrides(fc.StyleOverrides)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithThemedStyleOverrides(overrides))
	}

	return opts, nil
}

func parseStyleOverrides(raw map[string]any) (map[string]StyleValue, error) {
	result := make(map[string]StyleValue, len(raw))
	for k, v := range raw {
		key := normalizeVarName(k)
		switch val := v.(type) {
		case string:
			result[key] = StyleValue{Value: val}
		case []any:
			if len(val) != 2 {
				return nil, fmt.Errorf("kazari: styleOverrides[%q]: array must have exactly 2 elements [dark, light], got %d", k, len(val))
			}
			dark, ok1 := val[0].(string)
			light, ok2 := val[1].(string)
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("kazari: styleOverrides[%q]: array elements must be strings", k)
			}
			result[key] = StyleValue{Dark: dark, Light: light}
		case map[string]any:
			sv := StyleValue{}
			if l, ok := val["light"]; ok {
				s, ok := l.(string)
				if !ok {
					return nil, fmt.Errorf("kazari: styleOverrides[%q].light: must be a string", k)
				}
				sv.Light = s
			}
			if d, ok := val["dark"]; ok {
				s, ok := d.(string)
				if !ok {
					return nil, fmt.Errorf("kazari: styleOverrides[%q].dark: must be a string", k)
				}
				sv.Dark = s
			}
			if sv.Light == "" && sv.Dark == "" {
				return nil, fmt.Errorf("kazari: styleOverrides[%q]: map must have at least one of \"light\" or \"dark\" keys", k)
			}
			result[key] = sv
		default:
			return nil, fmt.Errorf("kazari: styleOverrides[%q]: value must be a string, [dark, light] array, or {light, dark} map", k)
		}
	}
	return result, nil
}

func splitCommaKey(key string) []string {
	parts := strings.Split(key, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// --- WithConfigDir ---

var configFileNames = []string{
	"kazari.config.yaml",
	"kazari.config.yml",
	"kazari.config.json",
}

// WithConfigDir searches the given directory for a config file
// (kazari.config.yaml, .yml, or .json) and applies its options.
// If no config file is found, it is a silent no-op. Parse errors
// are reported via WarningHandler.
func WithConfigDir(dir string) Option {
	return func(b *engineBuilder) {
		for _, name := range configFileNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err != nil {
				continue
			}
			opts, err := LoadConfig(path)
			if err != nil {
				if b.cfg.WarningHandler != nil {
					b.cfg.WarningHandler(fmt.Sprintf("kazari: loading config %s: %v", path, err))
				}
				return
			}
			for _, opt := range opts {
				opt(b)
			}
			return
		}
	}
}

// --- LoadConfig ---

// LoadConfig reads a config file, detects format from extension, parses,
// validates, and returns functional options.
func LoadConfig(path string) ([]Option, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("kazari: reading config file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var format string
	switch ext {
	case ".yaml", ".yml":
		format = "yaml"
	case ".json":
		format = "json"
	default:
		return nil, fmt.Errorf("kazari: unsupported config file extension %q, must be .yaml, .yml, or .json", ext)
	}

	fc, err := ParseConfig(data, format)
	if err != nil {
		return nil, err
	}

	return FileConfigToOptions(fc)
}
