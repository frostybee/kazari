package kazari

import (
	"fmt"
	"log"
	"strings"

	"github.com/frostybee/kazari/internal/color"
	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/render"
	"github.com/frostybee/kazari/internal/theme"
)

func (e *Engine) tokenize(code, lang, themeOverride string) ([]render.TokenLine, error) {
	if e.hl == nil {
		return plaintextLines(code), nil
	}

	if e.cfg.TabWidth > 0 {
		code = expandTabs(code, e.cfg.TabWidth)
	}

	lightTheme, darkTheme := e.resolveThemes(themeOverride)

	lines, err := e.tokenizeThemes(code, lang, lightTheme, darkTheme)
	if err != nil {
		if !e.isKnownLanguage(lang) {
			e.warn(fmt.Sprintf("kazari: unknown language %q, falling back to plaintext", lang))
			return plaintextLines(code), nil
		}
		return nil, err
	}
	return lines, nil
}

func (e *Engine) tokenizeThemes(code, lang, lightTheme, darkTheme string) ([]render.TokenLine, error) {
	if darkTheme != "" {
		if dual, ok := e.hl.(DualThemeTokenizer); ok {
			light, dark, err := dual.TokenizeDual(code, lang, lightTheme, darkTheme)
			if err != nil {
				return nil, err
			}
			return mergeTokens(light, dark), nil
		}
	}

	lightTokens, err := e.hl.Tokenize(code, lang, lightTheme)
	if err != nil {
		return nil, err
	}

	var darkTokens [][]Token
	if darkTheme != "" {
		darkTokens, err = e.hl.Tokenize(code, lang, darkTheme)
		if err != nil {
			return nil, err
		}
	}

	return mergeTokens(lightTokens, darkTokens), nil
}

func (e *Engine) isKnownLanguage(lang string) bool {
	for _, loaded := range e.hl.GetLoadedLanguages() {
		if strings.EqualFold(loaded, lang) {
			return true
		}
	}
	return false
}

func (e *Engine) warn(msg string) {
	if e.cfg.WarningHandler != nil {
		e.cfg.WarningHandler(msg)
		return
	}
	log.Print(msg)
}

func (e *Engine) resolveThemes(override string) (light, dark string) {
	light = e.cfg.LightTheme
	dark = e.cfg.DarkTheme
	if override == "" {
		return
	}
	if idx := strings.Index(override, ","); idx >= 0 {
		light = strings.TrimSpace(override[:idx])
		dark = strings.TrimSpace(override[idx+1:])
	} else {
		light = override
		// A single override theme applies to both modes, but only on
		// dual-theme engines so single-theme pages stay single-theme.
		if dark != "" {
			dark = override
		}
	}
	return
}

func (e *Engine) applyThemeOverride(resolved *config.ResolvedBlock) {
	lightName, darkName := e.resolveThemes(resolved.Theme)
	if lightName == e.cfg.LightTheme && darkName == e.cfg.DarkTheme {
		return
	}

	e.overrideMu.RLock()
	entry, ok := e.overrideCache[resolved.Theme]
	e.overrideMu.RUnlock()
	if !ok {
		entry = e.buildOverrideEntry(lightName, darkName)
		e.overrideMu.Lock()
		e.overrideCache[resolved.Theme] = entry
		e.overrideMu.Unlock()
	}

	resolved.ThemeOverrideStyle = entry.style
	resolved.LightEditorBG = entry.lightEditorBG
	resolved.DarkEditorBG = entry.darkEditorBG
	resolved.LightMarkerBGs = entry.lightBGs
	resolved.DarkMarkerBGs = entry.darkBGs
}

func (e *Engine) buildOverrideEntry(lightName, darkName string) overrideEntry {
	light, ok := extractThemeColors(e.hl, lightName, e.themeAdjustments, e.themeCustomizer)
	if !ok {
		e.warn(fmt.Sprintf("kazari: unknown theme %q in per-block override, keeping page colors", lightName))
		return overrideEntry{}
	}

	var dark theme.ThemeColors
	if e.cfg.DarkTheme != "" && darkName != "" {
		if colors, ok := extractThemeColors(e.hl, darkName, e.themeAdjustments, e.themeCustomizer); ok {
			dark = colors
		} else {
			e.warn(fmt.Sprintf("kazari: unknown theme %q in per-block override, keeping page colors for dark mode", darkName))
		}
	}

	entry := overrideEntry{style: theme.BlockOverrideStyle(e.cfg, light, dark)}
	if entry.style == "" {
		return overrideEntry{}
	}
	if e.cfg.MinContrast > 0 {
		if light.EditorBG != "" {
			entry.lightEditorBG = light.EditorBG
			entry.lightBGs = computeMarkerBGs(light.EditorBG)
		}
		if dark.EditorBG != "" {
			entry.darkEditorBG = dark.EditorBG
			entry.darkBGs = computeMarkerBGs(dark.EditorBG)
		}
	}
	return entry
}

func computeMarkerBGs(editorBG string) *config.MarkerEffectiveBGs {
	return &config.MarkerEffectiveBGs{
		Mark: compositeMarkerBG(config.MarkerMark, editorBG),
		Ins:  compositeMarkerBG(config.MarkerIns, editorBG),
		Del:  compositeMarkerBG(config.MarkerDel, editorBG),
	}
}

func compositeMarkerBG(mt config.MarkerType, editorBG string) string {
	rgba := config.MarkerBGColors[mt]
	hex, err := color.RGBAToHex(rgba)
	if err != nil {
		return editorBG
	}
	return color.OnBackground(hex, editorBG)
}
