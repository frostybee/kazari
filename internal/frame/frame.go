package frame

import (
	"path/filepath"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

var terminalLanguages = map[string]bool{
	"ansi": true, "bash": true, "bat": true, "batch": true, "cmd": true,
	"console": true, "fish": true, "nu": true, "nushell": true,
	"powershell": true, "ps": true, "ps1": true, "psd1": true, "psm1": true,
	"sh": true, "shell": true, "shellscript": true, "shellsession": true, "zsh": true,
}

var knownExtensions = map[string]bool{
	".go": true, ".js": true, ".ts": true, ".jsx": true, ".tsx": true,
	".py": true, ".rb": true, ".rs": true, ".c": true, ".h": true,
	".cpp": true, ".hpp": true, ".java": true, ".kt": true, ".swift": true,
	".cs": true, ".fs": true, ".php": true, ".lua": true, ".zig": true,
	".css": true, ".scss": true, ".sass": true, ".less": true,
	".html": true, ".htm": true, ".xml": true, ".svg": true,
	".json": true, ".jsonc": true, ".yaml": true, ".yml": true, ".toml": true,
	".md": true, ".mdx": true, ".txt": true, ".csv": true,
	".sh": true, ".bash": true, ".zsh": true, ".fish": true, ".ps1": true, ".bat": true, ".cmd": true,
	".sql": true, ".graphql": true, ".gql": true, ".proto": true,
	".Dockerfile": true, ".dockerfile": true,
	".env": true, ".ini": true, ".conf": true, ".cfg": true,
	".vue": true, ".svelte": true, ".astro": true,
	".wasm": true, ".tf": true, ".hcl": true, ".nix": true,
	".r": true, ".R": true, ".jl": true, ".ex": true, ".exs": true,
	".dart": true, ".groovy": true, ".scala": true, ".clj": true,
	".ml": true, ".mli": true, ".hs": true, ".elm": true,
	".makefile": true, ".mk": true,
}

func isTerminalLanguage(lang string) bool {
	return terminalLanguages[strings.ToLower(lang)]
}

// DetectFrameType determines the frame type for a code block.
func DetectFrameType(code, lang string, frameDefault int) int {
	if frameDefault != config.FrameAuto {
		return frameDefault
	}
	if !isTerminalLanguage(lang) {
		return config.FrameCode
	}
	if hasFileIndicator(code) {
		return config.FrameCode
	}
	return config.FrameTerminal
}

// ExtractFileName scans the first 4 lines for a comment containing a file path.
// Returns the extracted title and modified code (with comment line removed).
// If no file name is found, returns ("", original code).
func ExtractFileName(code, lang string) (string, string) {
	lines := strings.Split(code, "\n")
	limit := 4
	if len(lines) < limit {
		limit = len(lines)
	}

	for i := 0; i < limit; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if title := extractFromComment(trimmed); title != "" {
			modified := removeEmptyFrontmatter(removeLineFromCode(lines, i))
			return title, modified
		}
	}

	return "", code
}

// removeEmptyFrontmatter strips a frontmatter block that became empty after
// the file name comment was removed: two adjacent delimiter lines at the top
// of the code, plus one following blank line.
func removeEmptyFrontmatter(code string) string {
	lines := strings.Split(code, "\n")
	if len(lines) < 2 {
		return code
	}
	delim := strings.TrimSpace(lines[0])
	if delim != "---" && delim != "+++" {
		return code
	}
	if strings.TrimSpace(lines[1]) != delim {
		return code
	}
	rest := lines[2:]
	if len(rest) > 0 && strings.TrimSpace(rest[0]) == "" {
		rest = rest[1:]
	}
	return strings.Join(rest, "\n")
}

func hasFileIndicator(code string) bool {
	lines := strings.Split(code, "\n")
	limit := 4
	if len(lines) < limit {
		limit = len(lines)
	}
	for i := 0; i < limit; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "#!") {
			return true
		}
		if extractFromComment(trimmed) != "" {
			return true
		}
	}
	return false
}

func extractFromComment(line string) string {
	var content string

	switch {
	case strings.HasPrefix(line, "//"):
		content = strings.TrimSpace(line[2:])
	case strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "#!"):
		content = strings.TrimSpace(line[1:])
	case strings.HasPrefix(line, "<!--"):
		content = strings.TrimSpace(line[4:])
		content = strings.TrimSuffix(content, "-->")
		content = strings.TrimSpace(content)
	case strings.HasPrefix(line, "/*"):
		content = strings.TrimSpace(line[2:])
		content = strings.TrimSuffix(content, "*/")
		content = strings.TrimSpace(content)
	default:
		return ""
	}

	if content == "" {
		return ""
	}

	content = stripOptionalPrefix(content)

	if isFilePath(content) {
		return content
	}
	return ""
}

func stripOptionalPrefix(content string) string {
	lc := strings.ToLower(content)
	for _, prefix := range []string{"file name:", "filename:", "example:"} {
		if strings.HasPrefix(lc, prefix) {
			return strings.TrimSpace(content[len(prefix):])
		}
	}
	return content
}

func isFilePath(s string) bool {
	if s == "" {
		return false
	}
	// Reject URLs
	if strings.Contains(s, "://") {
		return false
	}
	// Reject strings with spaces (likely a sentence, not a path)
	if strings.Contains(s, " ") {
		return false
	}

	// Dotfiles (.gitignore, .env, etc.)
	if strings.HasPrefix(s, ".") && !strings.Contains(s, " ") && len(s) > 1 {
		return true
	}

	// Check for known extension
	ext := filepath.Ext(s)
	if ext != "" {
		if knownExtensions[ext] || knownExtensions[strings.ToLower(ext)] {
			return true
		}
		// Also match Makefile, Dockerfile etc. by name
		base := filepath.Base(s)
		if knownExtensions["."+strings.ToLower(base)] {
			return true
		}
	}

	// Files like "Makefile", "Dockerfile" without extension
	base := filepath.Base(s)
	switch strings.ToLower(base) {
	case "makefile", "dockerfile", "rakefile", "gemfile", "procfile", "cmakelists.txt":
		return true
	}

	return false
}

// StripTerminalComments removes lines whose first non-whitespace character is '#'.
func StripTerminalComments(code string) string {
	lines := strings.Split(code, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || !strings.HasPrefix(trimmed, "#") {
			kept = append(kept, line)
		}
	}
	result := strings.Join(kept, "\n")
	return strings.TrimRight(result, "\n")
}

func removeLineFromCode(lines []string, index int) string {
	result := make([]string, 0, len(lines)-1)
	for i, line := range lines {
		if i == index {
			continue
		}
		// Skip leading empty line after removal
		if i == index+1 && strings.TrimSpace(line) == "" {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
