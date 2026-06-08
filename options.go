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

func WithMinify(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.Minify = enabled }
}

func WithCascadeLayer(name string) Option {
	return func(b *engineBuilder) { b.cfg.CascadeLayer = name }
}

func WithLanguageAliases(m map[string]string) Option {
	return func(b *engineBuilder) { b.cfg.LanguageAliases = m }
}

func mapOptionsToBlockOpts(opts Options) *config.BlockOptions {
	bo := &config.BlockOptions{
		Lang:            opts.Lang,
		Title:           opts.Title,
		StartLineNumber: opts.StartLineNumber,
	}
	if opts.Frame != nil {
		f := int(*opts.Frame)
		bo.Frame = &f
	}
	if opts.LineNumbers != nil {
		bo.LineNumbers = opts.LineNumbers
	}
	if opts.Wrap != nil {
		bo.Wrap = opts.Wrap
	}
	return bo
}

