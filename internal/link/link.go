package link

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

var linkRe = regexp.MustCompile(`@\[([^\]]+)\]\(([^)]+)\)`)

func isSafeURL(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	u, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	s := strings.ToLower(u.Scheme)
	if s == "http" || s == "https" || s == "mailto" {
		return true
	}
	// Allow absolute paths (/foo) but reject protocol-relative URLs (//evil.com).
	return s == "" && strings.HasPrefix(trimmed, "/") && !strings.HasPrefix(trimmed, "//")
}

// ExtractLinks finds @[text](url) patterns in the source code, removes the
// syntax leaving only the link text, and returns annotations with character
// offsets into the cleaned text.
func ExtractLinks(code string) (string, [][]config.LinkAnnotation) {
	lines := strings.Split(code, "\n")
	allLinks := make([][]config.LinkAnnotation, len(lines))
	for i, line := range lines {
		cleaned, links := extractLineLinks(line)
		lines[i] = cleaned
		if len(links) > 0 {
			allLinks[i] = links
		}
	}
	return strings.Join(lines, "\n"), allLinks
}

func extractLineLinks(line string) (string, []config.LinkAnnotation) {
	matches := linkRe.FindAllStringSubmatchIndex(line, -1)
	if len(matches) == 0 {
		return line, nil
	}

	var links []config.LinkAnnotation
	var sb strings.Builder
	delta := 0
	prev := 0

	for _, loc := range matches {
		fullStart, fullEnd := loc[0], loc[1]
		textStart, textEnd := loc[2], loc[3]
		urlStart, urlEnd := loc[4], loc[5]

		linkText := line[textStart:textEnd]
		rawURL := line[urlStart:urlEnd]

		if !isSafeURL(rawURL) {
			// Unsafe scheme (e.g. javascript:) — keep original syntax as literal text.
			sb.WriteString(line[prev:fullEnd])
			prev = fullEnd
			continue
		}

		sb.WriteString(line[prev:fullStart])

		cleanStart := fullStart - delta
		cleanEnd := cleanStart + len(linkText)

		links = append(links, config.LinkAnnotation{
			Start: cleanStart,
			End:   cleanEnd,
			URL:   rawURL,
		})

		sb.WriteString(linkText)
		delta += (fullEnd - fullStart) - len(linkText)
		prev = fullEnd
	}

	sb.WriteString(line[prev:])
	return sb.String(), links
}
