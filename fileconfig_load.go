package kazari

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var configFileNames = []string{
	"kazari.config.yaml",
	"kazari.config.yml",
	"kazari.config.json",
}

// WithConfigDir searches the given directory for a config file
// (kazari.config.yaml, .yml, or .json) and applies its options.
// If no config file is found, it is a silent no-op. Parse errors
// are reported via WarningHandler.
func WithConfigDir(dir string) Option {
	return func(b *engineBuilder) {
		for _, name := range configFileNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err != nil {
				continue
			}
			opts, err := LoadConfig(path)
			if err != nil {
				msg := fmt.Sprintf("kazari: loading config %s: %v", path, err)
				if b.cfg.WarningHandler != nil {
					b.cfg.WarningHandler(msg)
				} else {
					log.Print(msg)
				}
				return
			}
			for _, opt := range opts {
				opt(b)
			}
			return
		}
	}
}

// LoadConfig reads a config file, detects format from extension, parses,
// validates, and returns functional options.
func LoadConfig(path string) ([]Option, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("kazari: reading config file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var format string
	switch ext {
	case ".yaml", ".yml":
		format = "yaml"
	case ".json":
		format = "json"
	default:
		return nil, fmt.Errorf("kazari: unsupported config file extension %q, must be .yaml, .yml, or .json", ext)
	}

	fc, err := ParseConfig(data, format)
	if err != nil {
		return nil, err
	}

	return FileConfigToOptions(fc)
}
