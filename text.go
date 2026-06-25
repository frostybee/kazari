package kazari

import (
	"strings"

	"github.com/frostybee/kazari/internal/render"
	"github.com/frostybee/kazari/internal/text"
)

func plaintextLines(code string) []render.TokenLine {
	rawLines := text.SplitLines(code)
	lines := make([]render.TokenLine, len(rawLines))
	for i, content := range rawLines {
		lines[i] = render.TokenLine{
			Tokens: []render.MergedToken{{Content: content}},
		}
	}
	return lines
}

func expandTabs(code string, tabWidth int) string {
	if !strings.Contains(code, "\t") {
		return code
	}
	spaces := strings.Repeat(" ", tabWidth)
	return strings.ReplaceAll(code, "\t", spaces)
}
