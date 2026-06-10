package js

import (
	"embed"
	"strings"

	"github.com/frostybee/kazari/internal/config"
	"github.com/frostybee/kazari/internal/minify"
)

//go:embed static/*.js
var staticFS embed.FS

func readJS(name string) string {
	data, err := staticFS.ReadFile("static/" + name)
	if err != nil {
		panic("kazari: missing embedded JS: " + name)
	}
	return string(data)
}

// Generate produces the JavaScript module output for the engine configuration.
func Generate(cfg *config.Config) string {
	var sb strings.Builder

	if cfg.CopyButton {
		sb.WriteString(readJS("copy.js"))
	}
	if cfg.FullscreenButton {
		sb.WriteString(readJS("fullscreen.js"))
	}
	if cfg.Collapsible != nil {
		sb.WriteString(readJS("collapsible.js"))
	}
	if cfg.CodeGroups {
		sb.WriteString(readJS("codegroup.js"))
	}

	content := sb.String()
	if content == "" {
		return ""
	}

	if cfg.Minify {
		return minify.JS(content)
	}
	return content
}
