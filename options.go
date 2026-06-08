package kazari

import "github.com/frostybee/kazari/internal/config"

// Option configures the Engine during construction.
type Option func(*engineBuilder)

// engineBuilder collects configuration during New().
type engineBuilder struct {
	cfg *config.Config
	hl  Highlighter
}

func WithHighlighter(hl Highlighter) Option {
	return func(b *engineBuilder) { b.hl = hl }
}

func WithThemes(light, dark string) Option {
	return func(b *engineBuilder) {
		b.cfg.LightTheme = light
		b.cfg.DarkTheme = dark
	}
}

func WithCopyButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.CopyButton = enabled }
}

func WithFullscreenButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.FullscreenButton = enabled }
}

func WithLineNumbers(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.LineNumbers = enabled }
}

func WithFrameDetection(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.FrameDetection = enabled }
}

func WithFileNameExtraction(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.FileNameExtraction = enabled }
}

func WithLanguageBadge(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.LanguageBadge = enabled }
}

func WithDarkMode(dm DarkMode) Option {
	return func(b *engineBuilder) {
		switch v := dm.(type) {
		case selectorMode:
			b.cfg.DarkMode = config.DarkModeConfig{
				Kind:     config.DarkModeSelectorKind,
				Selector: v.Selector,
			}
		case mediaQueryMode:
			b.cfg.DarkMode = config.DarkModeConfig{
				Kind: config.DarkModeMediaQueryKind,
			}
		case bothMode:
			b.cfg.DarkMode = config.DarkModeConfig{
				Kind:     config.DarkModeBothKind,
				Selector: v.Selector,
			}
		}
	}
}

func WithTabWidth(n int) Option {
	return func(b *engineBuilder) { b.cfg.TabWidth = n }
}

func WithStyleReset(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.StyleReset = enabled }
}

func WithThemedScrollbars(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.ThemedScrollbars = enabled }
}

func WithThemedSelectionColors(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.ThemedSelection = enabled }
}

func WithContentExclusion(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.ContentExclusion = enabled }
}

func WithDefaults(d BlockDefaults) Option {
	return func(b *engineBuilder) {
		b.cfg.Defaults = config.BlockDefaults{
			Wrap:           d.Wrap,
			PreserveIndent: d.PreserveIndent,
			HangingIndent:  d.HangingIndent,
			LineNumbers:    d.LineNumbers,
			Frame:          int(d.Frame),
		}
	}
}

func WithLanguageDefaults(m map[string]BlockDefaults) Option {
	return func(b *engineBuilder) {
		internal := make(map[string]config.BlockDefaults, len(m))
		for key, d := range m {
			internal[key] = config.BlockDefaults{
				Wrap:           d.Wrap,
				PreserveIndent: d.PreserveIndent,
				HangingIndent:  d.HangingIndent,
				LineNumbers:    d.LineNumbers,
				Frame:          int(d.Frame),
			}
		}
		b.cfg.LanguageDefaults = internal
	}
}

func WithCollapsible(c CollapsibleConfig) Option {
	return func(b *engineBuilder) {
		b.cfg.Collapsible = &config.CollapsibleConfig{
			LineThreshold:         c.LineThreshold,
			PreviewLines:          c.PreviewLines,
			DefaultCollapsed:      c.DefaultCollapsed,
			PreserveIndent:        c.PreserveIndent,
			ExpandButtonText:      c.ExpandButtonText,
			CollapseButtonText:    c.CollapseButtonText,
			ExpandedAnnouncement:  c.ExpandedAnnouncement,
			CollapsedAnnouncement: c.CollapsedAnnouncement,
		}
	}
}

func WithMinSyntaxHighlightingColorContrast(ratio float64) Option {
	return func(b *engineBuilder) { b.cfg.MinContrast = ratio }
}

func WithMinify(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.Minify = enabled }
}

func WithCascadeLayer(name string) Option {
	return func(b *engineBuilder) { b.cfg.CascadeLayer = name }
}

func WithLanguageAliases(m map[string]string) Option {
	return func(b *engineBuilder) { b.cfg.LanguageAliases = m }
}

