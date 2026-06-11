package diff

import (
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

// ProcessDiffBlock strips diff prefixes (+/-/space) from each line and returns
// the cleaned code along with line markers derived from the prefixes.
func ProcessDiffBlock(code string) (string, []config.LineMarker) {
	lines := strings.Split(code, "\n")
	stripped := make([]string, len(lines))

	var insLines, delLines []config.LineRange

	for i, line := range lines {
		lineNum := i + 1
		if len(line) == 0 {
			stripped[i] = ""
			continue
		}

		switch line[0] {
		case '+':
			stripped[i] = stripDiffPrefix(line)
			insLines = append(insLines, config.LineRange{Start: lineNum, End: lineNum})
		case '-':
			stripped[i] = stripDiffPrefix(line)
			delLines = append(delLines, config.LineRange{Start: lineNum, End: lineNum})
		case ' ':
			stripped[i] = line[1:]
		default:
			stripped[i] = line
		}
	}

	var markers []config.LineMarker
	if len(insLines) > 0 {
		markers = append(markers, config.LineMarker{Type: config.MarkerIns, Lines: insLines})
	}
	if len(delLines) > 0 {
		markers = append(markers, config.LineMarker{Type: config.MarkerDel, Lines: delLines})
	}

	return strings.Join(stripped, "\n"), markers
}

func stripDiffPrefix(line string) string {
	if len(line) <= 1 {
		return ""
	}
	if line[1] == ' ' {
		return line[2:]
	}
	return line[1:]
}
