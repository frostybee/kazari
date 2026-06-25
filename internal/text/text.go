package text

import "strings"

// SplitLines splits code into lines, trimming the trailing empty line
// produced by a final newline.
func SplitLines(code string) []string {
	if code == "" {
		return []string{""}
	}
	lines := strings.Split(code, "\n")
	if strings.HasSuffix(code, "\n") && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
