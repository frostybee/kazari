// CSS patterns derived from Expressive Code
// Copyright (c) Hippo (https://github.com/hippotastic)
// MIT License: https://github.com/expressive-code/expressive-code/blob/main/LICENSE

package css

import (
	"fmt"
	"strings"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/minify"
	"github.com/frostybee/kazari/internal/theme"
)

// Generate produces the complete CSS output for the engine configuration.
func Generate(cfg *config.Config, light, dark theme.ThemeColors) string {
	var sb strings.Builder

	// 1. Theme variables
	sb.WriteString(theme.GenerateVars(cfg, light, dark))

	// 2. Token color switching
	sb.WriteString(theme.TokenSwitchingCSS(cfg))

	// 3. Font style application
	sb.WriteString(fontStyleCSS())

	// 4. Base layout
	sb.WriteString(baseLayoutCSS(cfg))

	// 5. Scrollbar theming
	if cfg.ThemedScrollbars {
		sb.WriteString(scrollbarCSS())
	}

	// 6. Selection theming
	if cfg.ThemedSelection {
		sb.WriteString(selectionCSS())
	}

	content := sb.String()

	// Wrap in cascade layer if configured.
	if cfg.CascadeLayer != "" {
		content = fmt.Sprintf("@layer %s {\n%s}\n", cfg.CascadeLayer, content)
	}

	if cfg.Minify {
		return minify.CSS(content)
	}
	return content
}

func fontStyleCSS() string {
	return `.kazari-code .kz-line span[style*="--sfs"] { font-style: var(--sfs); }
.kazari-code .kz-line span[style*="--sfw"] { font-weight: var(--sfw); }
.kazari-code .kz-line span[style*="--std"] { text-decoration: var(--std); }
`
}

func baseLayoutCSS(cfg *config.Config) string {
	var sb strings.Builder

	sb.WriteString(`.kazari-code {
  position: relative;
  margin: 1rem 0;
}
`)

	// Style reset for future figure element (prep for Phase 2).
	if cfg.StyleReset {
		sb.WriteString(`.kazari-code figure {
  all: revert;
  position: relative;
  margin: 0;
  border-radius: var(--kz-radius);
  box-shadow: var(--kz-shadow);
  border: var(--kz-border);
  overflow: hidden;
}
`)
	}

	sb.WriteString(`.kazari-code pre {
  margin: 0;
  padding: 0;
  background: var(--kz-editor-bg);
  color: var(--kz-editor-fg);
  font-family: var(--kz-font-family);
  font-size: var(--kz-font-size);
  font-weight: var(--kz-font-weight);
  line-height: var(--kz-line-height);
  border-radius: var(--kz-radius);
  overflow-x: auto;
}
.kazari-code pre code {
  display: block;
  padding: var(--kz-code-padding-block) 0;
}
.kazari-code .kz-line {
  display: flex;
  min-height: 1lh;
}
.kazari-code .kz-line .code {
  flex: 1;
  padding-inline: var(--kz-code-padding-inline);
  white-space: pre;
}
`)

	// Word wrap support
	sb.WriteString(`.kazari-code pre.wrap .kz-line .code {
  white-space: pre-wrap;
  overflow-wrap: break-word;
}
`)

	return sb.String()
}

func scrollbarCSS() string {
	return `.kazari-code pre::-webkit-scrollbar {
  width: var(--kz-scrollbar-width, 8px);
  height: var(--kz-scrollbar-width, 8px);
}
.kazari-code pre::-webkit-scrollbar-thumb {
  background: var(--kz-scrollbar-thumb, rgba(255,255,255,0.15));
  border-radius: 4px;
}
.kazari-code pre::-webkit-scrollbar-thumb:hover {
  background: var(--kz-scrollbar-thumb-hover, rgba(255,255,255,0.3));
}
.kazari-code pre::-webkit-scrollbar-track {
  background: var(--kz-scrollbar-track, transparent);
}
`
}

func selectionCSS() string {
	return `.kazari-code ::selection {
  background: var(--kz-selection-bg, rgba(0,122,204,0.3));
  color: var(--kz-selection-fg, inherit);
}
`
}
