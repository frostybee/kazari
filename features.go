package kazari

import (
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

func WithCopyButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.CopyButton = enabled }
}

func WithFullscreenButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.FullscreenButton = enabled }
}

func WithWrapButton(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.WrapButton = enabled }
}

func WithLineNumbers(enabled bool) Option {
	return func(b *engineBuilder) {
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

func WithMermaidPassThrough(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.MermaidPassThrough = enabled }
}

func WithTerminalCommentStripping(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.TerminalCommentStripping = enabled }
}

func WithMinSyntaxHighlightingColorContrast(ratio float64) Option {
	return func(b *engineBuilder) { b.cfg.MinContrast = ratio }
}

func WithDataLineCount(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.DataLineCount = enabled }
}

func WithLanguageIconMode(mode LangIconMode) Option {
	return func(b *engineBuilder) { b.cfg.LangIconMode = int(mode) }
}

func WithFileIcons(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.FileIcons = enabled }
}

func WithFileIconResolver(f func(ext string) string) Option {
	return func(b *engineBuilder) { b.cfg.FileIconResolver = f }
}

func normalizeVarName(name string) string {
	if strings.HasPrefix(name, "--") {
		return name
	}
	return "--kz-" + name
}

func WithStyleOverrides(overrides map[string]string) Option {
	return func(b *engineBuilder) {
		if b.cfg.StyleOverrides == nil {
			b.cfg.StyleOverrides = make(map[string]config.StyleValue)
		}
		for k, v := range overrides {
			b.cfg.StyleOverrides[normalizeVarName(k)] = config.StyleValue{Value: v}
		}
	}
}

func WithThemedStyleOverrides(overrides map[string]StyleValue) Option {
	return func(b *engineBuilder) {
		if b.cfg.StyleOverrides == nil {
			b.cfg.StyleOverrides = make(map[string]config.StyleValue)
		}
		for k, v := range overrides {
			b.cfg.StyleOverrides[normalizeVarName(k)] = config.StyleValue{
				Value: v.Value,
				Dark:  v.Dark,
				Light: v.Light,
			}
		}
	}
}
