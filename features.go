package kazari

import "github.com/frostybee/kazari/internal/config"

func WithCopyButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.CopyButton = enabled }
}

func WithFullscreenButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.FullscreenButton = enabled }
}

func WithLineNumbers(enabled bool) Option {
	return func(b *engineBuilder) {
		b.cfg.LineNumbers = enabled
		b.cfg.Defaults.LineNumbers = enabled
	}
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

func WithThemedScrollbars(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.ThemedScrollbars = enabled }
}

func WithThemedSelectionColors(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.ThemedSelection = enabled }
}

func WithContentExclusion(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.ContentExclusion = enabled }
}

func WithCollapsible(c CollapsibleConfig) Option {
	return func(b *engineBuilder) {
		b.cfg.Collapsible = &config.CollapsibleConfig{
			LineThreshold:         c.LineThreshold,
			PreviewLines:          c.PreviewLines,
			DefaultCollapsed:      c.DefaultCollapsed,
			PreserveIndent:        c.PreserveIndent,
			Style:                 config.CollapseStyle(c.Style),
			ExpandButtonText:      c.ExpandButtonText,
			CollapseButtonText:    c.CollapseButtonText,
			ExpandedAnnouncement:  c.ExpandedAnnouncement,
			CollapsedAnnouncement: c.CollapsedAnnouncement,
		}
	}
}

func WithTerminalDotStyle(style TerminalDotStyle) Option {
	return func(b *engineBuilder) { b.cfg.TerminalDotStyle = int(style) }
}

func WithMinSyntaxHighlightingColorContrast(ratio float64) Option {
	return func(b *engineBuilder) { b.cfg.MinContrast = ratio }
}
