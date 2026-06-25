package theme

import (
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// ThemeToggleCSS generates CSS rules for per-block theme toggling.
// These rules use data-kz-theme="dark"|"light" on .kazari-block to
// override the page-level theme for individual blocks.
func ThemeToggleCSS(cfg *config.Config, light, dark ThemeColors) string {
	if cfg.DarkTheme == "" {
		return ""
	}

	var sb strings.Builder

	writeToggleVars(&sb, ".kazari-block[data-kz-theme=\"dark\"]", blockOverridableVars(dark, cfg), cfg)
	writeToggleVars(&sb, ".kazari-block[data-kz-theme=\"light\"]", blockOverridableVars(light, cfg), cfg)

	sb.WriteString(tokenSwitchRule(".kazari-block[data-kz-theme=\"dark\"]", "--sd", "--sdbg"))
	sb.WriteString(tokenSwitchRule(".kazari-block[data-kz-theme=\"light\"]", "--sl", "--slbg"))

	// kz-themed + force-dark: remap override dark vars onto kz-* names.
	sb.WriteString(themedToggleDarkRule(cfg))
	// kz-themed + force-light: remap override light vars onto kz-* names.
	sb.WriteString(themedToggleLightRule(cfg))

	return sb.String()
}

func themedToggleDarkRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed[data-kz-theme=\"dark\"]", darkVarTemplate, true, cfg)
}

func themedToggleLightRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed[data-kz-theme=\"light\"]", lightVarTemplate, true, cfg)
}

func writeToggleVars(sb *strings.Builder, selector string, vars []struct{ name, value string }, cfg *config.Config) {
	sb.WriteString(selector)
	sb.WriteString(" { ")
	for _, v := range vars {
		if v.value != "" {
			sb.WriteString(fmt.Sprintf("%s: %s; ", v.name, v.value))
		}
	}
	if cfg.Collapsible != nil {
		sb.WriteString(collapseGradientDecl)
	}
	sb.WriteString("}\n")
}
