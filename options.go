package kazari

import "github.com/frostybee/kazari/internal/config"

// Option configures the Engine during construction.
type Option func(*engineBuilder)

// engineBuilder collects configuration during New().
type engineBuilder struct {
	cfg              *config.Config
	hl               Highlighter
	themeCustomizer  func(string, ThemeInfo) ThemeInfo
	themeAdjustments *ThemeAdjustments
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
		b.cfg.Defaults = toConfigBlockDefaults(d)
	}
}

func WithLanguageDefaults(m map[string]BlockDefaults) Option {
	return func(b *engineBuilder) {
		if b.cfg.LanguageDefaults == nil {
			b.cfg.LanguageDefaults = make(map[string]config.BlockDefaults)
		}
		for key, d := range m {
			b.cfg.LanguageDefaults[key] = toConfigBlockDefaults(d)
		}
	}
}

func toConfigBlockDefaults(d BlockDefaults) config.BlockDefaults {
	return config.BlockDefaults{
		Wrap:           d.Wrap,
		PreserveIndent: d.PreserveIndent,
		HangingIndent:  d.HangingIndent,
		LineNumbers:    d.LineNumbers,
		Frame:          int(d.Frame),
	}
}

func WithMinify(enabled bool) Option {
	return func(b *engineBuilder) { b.cfg.Minify = enabled }
}

func WithCascadeLayer(name string) Option {
	return func(b *engineBuilder) { b.cfg.CascadeLayer = name }
}

func WithLanguageAliases(m map[string]string) Option {
	return func(b *engineBuilder) {
		if b.cfg.LanguageAliases == nil {
			b.cfg.LanguageAliases = make(map[string]string)
		}
		for k, v := range m {
			b.cfg.LanguageAliases[k] = v
		}
	}
}

func WithThemeCSSRoot(selector string) Option {
	return func(b *engineBuilder) { b.cfg.ThemeCSSRoot = selector }
}

func WithThemeCustomizer(f func(themeName string, colors ThemeInfo) ThemeInfo) Option {
	return func(b *engineBuilder) { b.themeCustomizer = f }
}

// WithThemeAdjustments tints the extracted theme colors in OKLCH space,
// applied to both themes before the theme customizer runs.
func WithThemeAdjustments(adj ThemeAdjustments) Option {
	return func(b *engineBuilder) { b.themeAdjustments = &adj }
}

func WithLocale(loc string) Option {
	return func(b *engineBuilder) { b.cfg.Locale = loc }
}

func WithUIStrings(overrides map[string]string) Option {
	return func(b *engineBuilder) {
		if b.cfg.UIStringOverrides == nil {
			b.cfg.UIStringOverrides = make(map[string]string)
		}
		for k, v := range overrides {
			b.cfg.UIStringOverrides[k] = v
		}
	}
}

// WithWarningHandler sets the function that receives non-fatal warnings,
// such as unknown language fallbacks. When unset, warnings go to log.Printf.
// Pass a no-op function to silence warnings entirely.
func WithWarningHandler(f func(string)) Option {
	return func(b *engineBuilder) { b.cfg.WarningHandler = f }
}

func mapOptionsToBlockOpts(opts Options) *config.BlockOptions {
	bo := &config.BlockOptions{
		Lang:            opts.Lang,
		Title:           opts.Title,
		Theme:           opts.Theme,
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
	if opts.PreserveIndent != nil {
		bo.PreserveIndent = opts.PreserveIndent
	}
	if opts.HangingIndent != nil {
		bo.HangingIndent = opts.HangingIndent
	}
	return bo
}

