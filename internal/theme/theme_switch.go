package theme

import (
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// TokenSwitchingCSS generates the color switching rules for token spans.
func TokenSwitchingCSS(cfg *config.Config) string {
	var sb strings.Builder

	sb.WriteString(tokenSwitchRule(".kazari-block", "--sl", "--slbg"))
	sb.WriteString(themedLightRule(cfg))

	if cfg.DarkTheme == "" {
		return sb.String()
	}

	darkRules := tokenSwitchRule(".kazari-block", "--sd", "--sdbg") +
		themedDarkRule(cfg)

	switch cfg.DarkMode.Kind {
	case config.DarkModeSelectorKind:
		writeScopedRules(&sb, cfg.DarkMode.Selector, darkRules)

	case config.DarkModeMediaQueryKind:
		sb.WriteString("@media (prefers-color-scheme: dark) {\n")
		sb.WriteString(darkRules)
		sb.WriteString("}\n")

	case config.DarkModeBothKind:
		sb.WriteString("@media (prefers-color-scheme: dark) {\n")
		sb.WriteString(darkRules)
		sb.WriteString("}\n")
		writeScopedRules(&sb, cfg.DarkMode.Selector, darkRules)
	}

	return sb.String()
}

// writeScopedRules prefixes each rule line with the dark mode selector.
func writeScopedRules(sb *strings.Builder, selector, rules string) {
	for _, line := range strings.Split(strings.TrimRight(rules, "\n"), "\n") {
		sb.WriteString(selector + " " + line + "\n")
	}
}

func tokenSwitchRule(selector, colorVar, bgVar string) string {
	return fmt.Sprintf(
		"%s .kz-line span[style^=\"--\"] { color: var(%s, inherit); background-color: var(%s, transparent); font-style: var(--sfs, inherit); font-weight: var(--sfw, inherit); text-decoration: var(--std, inherit); }\n",
		selector, colorVar, bgVar,
	)
}

func writeThemedRule(selector, varTemplate string, includeGradient bool, cfg *config.Config) string {
	var sb strings.Builder
	sb.WriteString(selector)
	sb.WriteString(" { ")
	for _, name := range overridableVarNames(cfg) {
		suffix := strings.TrimPrefix(name, "--kz-")
		sb.WriteString(fmt.Sprintf(varTemplate, name, suffix, suffix))
	}
	if includeGradient && cfg.Collapsible != nil {
		sb.WriteString(collapseGradientDecl)
	}
	sb.WriteString("}\n")
	return sb.String()
}

func themedLightRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed", lightVarTemplate, true, cfg)
}

// themedDarkRule must stay on a single line because writeScopedRules prefixes
// each line with the dark mode selector. No gradient re-declaration here
// because it is emitted inside the scoped wrapper.
func themedDarkRule(cfg *config.Config) string {
	return writeThemedRule(".kazari-block.kz-themed", darkVarTemplate, false, cfg)
}
