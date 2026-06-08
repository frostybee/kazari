package meta

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/frostybee/kazari/internal/config"
)

// MarkerType identifies the kind of marker.
type MarkerType int

const (
	MarkerMark MarkerType = iota
	MarkerIns
	MarkerDel
)

// LineRange is an inclusive 1-based line range.
type LineRange struct {
	Start int
	End   int
}

// LineMarker represents a line-level marker with optional label.
type LineMarker struct {
	Type  MarkerType
	Lines []LineRange
	Label string
}

// InlineMarker represents an inline text marker.
type InlineMarker struct {
	Type MarkerType
	Text string
}

// CollapseSpec holds per-block collapse info from meta.
type CollapseSpec struct {
	Enabled  bool
	Disabled bool
	Ranges   []LineRange
}

// ParseResult contains the full parsed meta string output.
type ParseResult struct {
	BlockOptions  config.BlockOptions
	LineMarkers   []LineMarker
	InlineMarkers []InlineMarker
	FocusLines    []LineRange
	Collapse      *CollapseSpec
}

// Parse parses a fence info meta string into structured data.
// The first bare word (no = or {) is treated as the language.
func Parse(meta string) *ParseResult {
	result := &ParseResult{}
	tokens := tokenize(meta)

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		switch {
		case i == 0 && isBareLang(tok):
			result.BlockOptions.Lang = tok

		case tok == "showLineNumbers":
			v := true
			result.BlockOptions.LineNumbers = &v

		case tok == "showLineNumbers=false":
			v := false
			result.BlockOptions.LineNumbers = &v

		case tok == "wrap":
			v := true
			result.BlockOptions.Wrap = &v

		case tok == "collapse":
			if result.Collapse == nil {
				result.Collapse = &CollapseSpec{}
			}
			result.Collapse.Enabled = true

		case tok == "nocollapse":
			if result.Collapse == nil {
				result.Collapse = &CollapseSpec{}
			}
			result.Collapse.Disabled = true

		case strings.HasPrefix(tok, "title="):
			result.BlockOptions.Title = unquote(strings.TrimPrefix(tok, "title="))

		case strings.HasPrefix(tok, "frame="):
			frame := unquote(strings.TrimPrefix(tok, "frame="))
			var f int
			switch frame {
			case "code":
				f = config.FrameCode
			case "terminal":
				f = config.FrameTerminal
			case "none":
				f = config.FrameNone
			default:
				f = config.FrameAuto
			}
			result.BlockOptions.Frame = &f

		case strings.HasPrefix(tok, "startLineNumber="):
			val := strings.TrimPrefix(tok, "startLineNumber=")
			if n, err := strconv.Atoi(val); err == nil {
				result.BlockOptions.StartLineNumber = n
			}

		case strings.HasPrefix(tok, "focus="):
			rangeStr := extractBraces(strings.TrimPrefix(tok, "focus="))
			result.FocusLines = parseRanges(rangeStr)

		case strings.HasPrefix(tok, "ins="):
			remainder := strings.TrimPrefix(tok, "ins=")
			parseMarkerToken(remainder, MarkerIns, result)

		case strings.HasPrefix(tok, "del="):
			remainder := strings.TrimPrefix(tok, "del=")
			parseMarkerToken(remainder, MarkerDel, result)

		case strings.HasPrefix(tok, "add="):
			remainder := strings.TrimPrefix(tok, "add=")
			parseMarkerToken(remainder, MarkerIns, result)

		case strings.HasPrefix(tok, "rem="):
			remainder := strings.TrimPrefix(tok, "rem=")
			parseMarkerToken(remainder, MarkerDel, result)

		case strings.HasPrefix(tok, "collapse="):
			rangeStr := extractBraces(strings.TrimPrefix(tok, "collapse="))
			if result.Collapse == nil {
				result.Collapse = &CollapseSpec{}
			}
			result.Collapse.Ranges = parseRanges(rangeStr)

		case strings.HasPrefix(tok, "{") && strings.HasSuffix(tok, "}"):
			// Bare highlight ranges: {3-5,8}
			rangeStr := tok[1 : len(tok)-1]
			result.LineMarkers = append(result.LineMarkers, LineMarker{
				Type:  MarkerMark,
				Lines: parseRanges(rangeStr),
			})

		case isQuotedString(tok):
			// Inline text marker: "text"
			result.InlineMarkers = append(result.InlineMarkers, InlineMarker{
				Type: MarkerMark,
				Text: unquote(tok),
			})
		}
	}

	return result
}

// parseMarkerToken handles ins="text", ins={lines}, ins={"label":lines}
func parseMarkerToken(remainder string, mtype MarkerType, result *ParseResult) {
	if isQuotedString(remainder) {
		// Inline text marker: ins="text"
		result.InlineMarkers = append(result.InlineMarkers, InlineMarker{
			Type: mtype,
			Text: unquote(remainder),
		})
	} else if strings.HasPrefix(remainder, "{") {
		inner := extractBraces(remainder)
		// Check for labeled range: {"A":6-10}
		if label, ranges, ok := parseLabeledRange(inner); ok {
			result.LineMarkers = append(result.LineMarkers, LineMarker{
				Type:  mtype,
				Lines: ranges,
				Label: label,
			})
		} else {
			result.LineMarkers = append(result.LineMarkers, LineMarker{
				Type:  mtype,
				Lines: parseRanges(inner),
			})
		}
	}
}

// parseLabeledRange handles "A":6-10 inside braces.
func parseLabeledRange(s string) (string, []LineRange, bool) {
	if !strings.HasPrefix(s, "\"") {
		return "", nil, false
	}
	// Find closing quote
	end := strings.Index(s[1:], "\"")
	if end < 0 {
		return "", nil, false
	}
	label := s[1 : end+1]
	rest := s[end+2:]
	// Expect ":"
	if !strings.HasPrefix(rest, ":") {
		return "", nil, false
	}
	ranges := parseRanges(rest[1:])
	return label, ranges, true
}

// parseRanges parses "3-5,8,10-12" into LineRange slices.
func parseRanges(s string) []LineRange {
	var ranges []LineRange
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if dash := strings.Index(part, "-"); dash > 0 {
			start, err1 := strconv.Atoi(strings.TrimSpace(part[:dash]))
			end, err2 := strconv.Atoi(strings.TrimSpace(part[dash+1:]))
			if err1 == nil && err2 == nil {
				ranges = append(ranges, LineRange{Start: start, End: end})
			}
		} else {
			if n, err := strconv.Atoi(part); err == nil {
				ranges = append(ranges, LineRange{Start: n, End: n})
			}
		}
	}
	return ranges
}

// tokenize splits the meta string into logical tokens, preserving quoted strings
// and brace groups as single tokens.
func tokenize(meta string) []string {
	var tokens []string
	runes := []rune(meta)
	i := 0

	for i < len(runes) {
		// Skip whitespace
		if unicode.IsSpace(runes[i]) {
			i++
			continue
		}

		start := i

		// Quoted string at top level (inline text marker)
		if runes[i] == '"' {
			i++
			for i < len(runes) && runes[i] != '"' {
				if runes[i] == '\\' {
					i++
				}
				i++
			}
			if i < len(runes) {
				i++ // closing quote
			}
			tokens = append(tokens, string(runes[start:i]))
			continue
		}

		// Collect until whitespace, but handle quoted values and brace groups inline
		for i < len(runes) && !unicode.IsSpace(runes[i]) {
			if runes[i] == '"' {
				i++
				for i < len(runes) && runes[i] != '"' {
					if runes[i] == '\\' {
						i++
					}
					i++
				}
				if i < len(runes) {
					i++ // closing quote
				}
			} else if runes[i] == '{' {
				depth := 0
				for i < len(runes) {
					if runes[i] == '{' {
						depth++
					} else if runes[i] == '}' {
						depth--
						if depth == 0 {
							i++
							break
						}
					}
					i++
				}
			} else {
				i++
			}
		}

		tokens = append(tokens, string(runes[start:i]))
	}

	return tokens
}

func isBareLang(tok string) bool {
	if tok == "" {
		return false
	}
	if strings.Contains(tok, "=") || strings.HasPrefix(tok, "{") || strings.HasPrefix(tok, "\"") {
		return false
	}
	return true
}

func isQuotedString(tok string) bool {
	return len(tok) >= 2 && tok[0] == '"' && tok[len(tok)-1] == '"'
}

func unquote(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	return strings.ReplaceAll(s, `\"`, `"`)
}

func extractBraces(s string) string {
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		return s[1 : len(s)-1]
	}
	return s
}
