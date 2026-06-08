package js

import "github.com/frostybee/kazari/internal/config"

// Generate produces the JavaScript module output for the engine configuration.
// Phase 1: returns empty string (no interactive features yet).
// Later phases add copy button, collapsible, and fullscreen handlers.
func Generate(cfg *config.Config) string {
	_ = cfg
	return ""
}
