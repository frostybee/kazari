package ansi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

var standardColors = [16]string{
	"#000000", // 0  black
	"#cc0000", // 1  red
	"#4e9a06", // 2  green
	"#c4a000", // 3  yellow
	"#3465a4", // 4  blue
	"#75507b", // 5  magenta
	"#06989a", // 6  cyan
	"#d3d7cf", // 7  white
	"#555753", // 8  bright black
	"#ef2929", // 9  bright red
	"#8ae234", // 10 bright green
	"#fce94f", // 11 bright yellow
	"#729fcf", // 12 bright blue
	"#ad7fa8", // 13 bright magenta
	"#34e2e2", // 14 bright cyan
	"#eeeeec", // 15 bright white
}

var ansiVarNames = [16]string{
	"--kz-ansi-black", "--kz-ansi-red", "--kz-ansi-green", "--kz-ansi-yellow",
	"--kz-ansi-blue", "--kz-ansi-magenta", "--kz-ansi-cyan", "--kz-ansi-white",
	"--kz-ansi-bright-black", "--kz-ansi-bright-red", "--kz-ansi-bright-green",
	"--kz-ansi-bright-yellow", "--kz-ansi-bright-blue", "--kz-ansi-bright-magenta",
	"--kz-ansi-bright-cyan", "--kz-ansi-bright-white",
}

func standardColorVar(idx int) string {
	return "var(" + ansiVarNames[idx] + ")"
}

type state struct {
	fg        string
	bg        string
	fontStyle int
}

func (s *state) reset() {
	s.fg = ""
	s.bg = ""
	s.fontStyle = 0
}

// Parse converts ANSI-escaped text into token lines for rendering.
// SGR escape codes are parsed into colors and font styles; the escape
// sequences themselves are stripped from the output tokens.
func Parse(code string) []config.TokenLine {
	rawLines := splitLines(code)
	lines := make([]config.TokenLine, len(rawLines))
	var st state

	for i, raw := range rawLines {
		lines[i] = parseLine(raw, &st)
	}
	return lines
}

func parseLine(line string, st *state) config.TokenLine {
	var tokens []config.MergedToken
	var buf strings.Builder

	flush := func() {
		if buf.Len() == 0 {
			return
		}
		tokens = append(tokens, config.MergedToken{
			Content:    buf.String(),
			LightColor: st.fg,
			DarkColor:  st.fg,
			LightBG:    st.bg,
			DarkBG:     st.bg,
			FontStyle:  st.fontStyle,
		})
		buf.Reset()
	}

	pos := 0
	for pos < len(line) {
		if line[pos] == '\x1b' && pos+1 < len(line) && line[pos+1] == '[' {
			end := findSequenceEnd(line, pos+2)
			if end < 0 {
				buf.WriteByte(line[pos])
				pos++
				continue
			}
			flush()
			params := line[pos+2 : end]
			applyParams(params, st)
			pos = end + 1
		} else {
			buf.WriteByte(line[pos])
			pos++
		}
	}
	flush()

	if len(tokens) == 0 {
		tokens = append(tokens, config.MergedToken{Content: ""})
	}
	return config.TokenLine{Tokens: tokens}
}

func findSequenceEnd(line string, start int) int {
	for i := start; i < len(line); i++ {
		b := line[i]
		if b == 'm' {
			return i
		}
		if b != ';' && (b < '0' || b > '9') {
			return -1
		}
	}
	return -1
}

func applyParams(params string, st *state) {
	if params == "" {
		st.reset()
		return
	}
	codes := strings.Split(params, ";")
	for i := 0; i < len(codes); i++ {
		code, err := strconv.Atoi(codes[i])
		if err != nil {
			continue
		}
		switch {
		case code == 0:
			st.reset()
		case code == 1:
			st.fontStyle |= 2 // bold
		case code == 3:
			st.fontStyle |= 1 // italic
		case code == 4:
			st.fontStyle |= 4 // underline
		case code == 9:
			st.fontStyle |= 8 // strikethrough
		case code == 22:
			st.fontStyle &^= 2
		case code == 23:
			st.fontStyle &^= 1
		case code == 24:
			st.fontStyle &^= 4
		case code == 29:
			st.fontStyle &^= 8
		case code >= 30 && code <= 37:
			st.fg = standardColorVar(code - 30)
		case code == 38:
			i = applyExtendedColor(codes, i, &st.fg)
		case code == 39:
			st.fg = ""
		case code >= 40 && code <= 47:
			st.bg = standardColorVar(code - 40)
		case code == 48:
			i = applyExtendedColor(codes, i, &st.bg)
		case code == 49:
			st.bg = ""
		case code >= 90 && code <= 97:
			st.fg = standardColorVar(code - 82)
		case code >= 100 && code <= 107:
			st.bg = standardColorVar(code - 92)
		}
	}
}

func applyExtendedColor(codes []string, i int, target *string) int {
	if i+1 >= len(codes) {
		return i
	}
	mode, err := strconv.Atoi(codes[i+1])
	if err != nil {
		return i
	}
	switch mode {
	case 5:
		if i+2 >= len(codes) {
			return i + 1
		}
		n, err := strconv.Atoi(codes[i+2])
		if err != nil || n < 0 || n > 255 {
			return i + 2
		}
		*target = color256(n)
		return i + 2
	case 2:
		if i+4 >= len(codes) {
			return i + 1
		}
		r, err1 := strconv.Atoi(codes[i+2])
		g, err2 := strconv.Atoi(codes[i+3])
		b, err3 := strconv.Atoi(codes[i+4])
		if err1 != nil || err2 != nil || err3 != nil {
			return i + 4
		}
		*target = fmt.Sprintf("#%02x%02x%02x", clamp(r), clamp(g), clamp(b))
		return i + 4
	}
	return i + 1
}

func color256(n int) string {
	if n < 16 {
		return standardColors[n]
	}
	if n < 232 {
		n -= 16
		b := n % 6
		n /= 6
		g := n % 6
		r := n / 6
		return fmt.Sprintf("#%02x%02x%02x", cubeValue(r), cubeValue(g), cubeValue(b))
	}
	v := 8 + (n-232)*10
	return fmt.Sprintf("#%02x%02x%02x", v, v, v)
}

func cubeValue(i int) int {
	if i == 0 {
		return 0
	}
	return 55 + i*40
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func splitLines(code string) []string {
	if code == "" {
		return []string{""}
	}
	lines := strings.Split(code, "\n")
	if strings.HasSuffix(code, "\n") && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
