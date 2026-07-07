package kazari

import (
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

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
	if fc.ThemeToggleButton != nil {
		opts = append(opts, WithThemeToggle(*fc.ThemeToggleButton))
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
	if fc.Links != nil {
		opts = append(opts, WithLinks(*fc.Links))
	}
	if fc.StyleReset != nil {
		opts = append(opts, WithStyleReset(*fc.StyleReset))
	}
	if fc.Minify != nil {
		opts = append(opts, WithMinify(*fc.Minify))
	}

	if fc.OutputPanel != nil {
		opts = append(opts, WithOutputPanel(*fc.OutputPanel))
	}
	if fc.OutputDefaultCollapsed != nil {
		opts = append(opts, WithOutputCollapsed(*fc.OutputDefaultCollapsed))
	}
	if fc.OutputSeparator != nil {
		opts = append(opts, WithOutputSeparator(*fc.OutputSeparator))
	}

	if fc.TabWidth != nil {
		opts = append(opts, WithTabWidth(*fc.TabWidth))
	}
	if fc.MinContrast != nil {
		opts = append(opts, WithMinContrast(*fc.MinContrast))
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

	if fc.LanguageIconMode != nil {
		mode, err := parseLangIconModeStr(*fc.LanguageIconMode)
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithLanguageIconMode(mode))
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
