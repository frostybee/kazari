package unwrap

import (
	"bytes"

	"golang.org/x/net/html"
)

// CandidateKind distinguishes the two discovery entry points, which map to
// the two unwrapper chains.
type CandidateKind int

const (
	CandidateWrapper CandidateKind = iota
	CandidateBarePre
)

// Candidate is one potential code block region found in a page. Tokens spans
// the root start tag through its matching end tag, preserving the convention
// every Unwrapper and SkipReason relies on: the root tag sits at index 0.
// Malformed candidates (unclosed root element) carry only the root token and
// must never be spliced; their bytes stay untouched.
type Candidate struct {
	Kind      CandidateKind
	Tokens    []BufferedToken
	ByteStart int
	ByteEnd   int
	Malformed bool
}

// AssetTag is an existing link or script tag injected by a previous run,
// identified by its data-kazari="assets" marker.
type AssetTag struct {
	ByteStart int
	ByteEnd   int
	Kind      string // link or script
	Ref       string // href or src value as found
}

// Page is the result of discovering one HTML file.
type Page struct {
	Candidates []Candidate
	AssetTags  []AssetTag
	Generator  string // content of meta name="generator", empty if absent

	// LinkInsert is the byte offset where a stylesheet link belongs: before
	// the head close tag, else right after the body start tag, else the
	// start of content (byte 3 when the file opens with a UTF-8 BOM).
	LinkInsert int
	// ScriptInsert is the byte offset where the script tag belongs: before
	// the body close tag, else before the html close tag, else EOF.
	ScriptInsert int

	// Suppressed counts candidates that were not emitted because they sit
	// inside an element carrying a kz- or kazari- prefixed class. That
	// namespace is reserved for Kazari's own output; the count keeps the
	// suppression visible instead of silently omitting content.
	Suppressed int
	// SuppressedUnclosed reports that a reserved class scope was still open
	// at the end of the file, so suppression ran to EOF (fail closed). A
	// scope also ends when any ancestor element closes, mirroring browser
	// recovery for unclosed elements, in which case this stays false.
	SuppressedUnclosed bool
}

// Discover walks a tokenized HTML file once and returns its candidates,
// injection offsets, generator identity, and existing asset tags. All byte
// offsets are positions in the ORIGINAL src bytes; splicing merges them into
// a single ascending edit list and never re-finds offsets in mutated output.
func Discover(src []byte, tokens []BufferedToken) Page {
	var page Page
	headEnd, bodyAfter, bodyEnd, htmlEnd := -1, -1, -1, -1

	depth := 0
	suppressAt := -1 // element depth at which a reserved class scope opened

	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		switch tok.Type {
		case html.StartTagToken, html.SelfClosingTagToken:
			void := voidElements[tok.Data] || tok.Type == html.SelfClosingTagToken

			switch tok.Data {
			case "meta":
				if n, _ := attrValue(tok.Attrs, "name"); n == "generator" && page.Generator == "" {
					page.Generator, _ = attrValue(tok.Attrs, "content")
				}
			case "body":
				if bodyAfter < 0 && tok.Type == html.StartTagToken {
					bodyAfter = tok.Offset + len(tok.Raw)
				}
			case "link":
				if v, _ := attrValue(tok.Attrs, "data-kazari"); v == "assets" {
					ref, _ := attrValue(tok.Attrs, "href")
					page.AssetTags = append(page.AssetTags, AssetTag{
						ByteStart: tok.Offset,
						ByteEnd:   tok.Offset + len(tok.Raw),
						Kind:      "link",
						Ref:       ref,
					})
				}
			case "script":
				if v, _ := attrValue(tok.Attrs, "data-kazari"); v == "assets" {
					end := tok.Offset + len(tok.Raw)
					if j, closed := findMatchingEnd(tokens, i); closed {
						end = tokens[j].Offset + len(tokens[j].Raw)
					}
					ref, _ := attrValue(tok.Attrs, "src")
					page.AssetTags = append(page.AssetTags, AssetTag{
						ByteStart: tok.Offset,
						ByteEnd:   end,
						Kind:      "script",
						Ref:       ref,
					})
				}
			}

			if suppressAt < 0 && !void && hasReservedClass(tok.Attrs) {
				suppressAt = depth
			}

			if suppressAt >= 0 {
				if isWrapperTrigger(tok) || (tok.Data == "pre" && !void) {
					page.Suppressed++
				}
				if !void {
					depth++
				}
				i++
				continue
			}

			if !void && isWrapperTrigger(tok) {
				cand, next, closed := makeCandidate(tokens, i, CandidateWrapper)
				page.Candidates = append(page.Candidates, cand)
				if closed {
					// The subtree is balanced, so depth is unchanged and
					// nested candidates are consumed: outer wins.
					i = next
					continue
				}
				// Unclosed root: resume at the very next token so one
				// malformed wrapper never blinds the rest of the file.
				depth++
				i++
				continue
			}
			if !void && tok.Data == "pre" {
				cand, next, closed := makeCandidate(tokens, i, CandidateBarePre)
				page.Candidates = append(page.Candidates, cand)
				if closed {
					i = next
					continue
				}
				depth++
				i++
				continue
			}

			if !void {
				depth++
			}
			i++

		case html.EndTagToken:
			switch tok.Data {
			case "head":
				if headEnd < 0 {
					headEnd = tok.Offset
				}
			case "body":
				if bodyEnd < 0 {
					bodyEnd = tok.Offset
				}
			case "html":
				if htmlEnd < 0 {
					htmlEnd = tok.Offset
				}
			}
			if !voidElements[tok.Data] {
				if depth > 0 {
					depth--
				}
				if suppressAt >= 0 && depth <= suppressAt {
					suppressAt = -1
				}
			}
			i++

		default:
			i++
		}
	}
	if suppressAt >= 0 {
		page.SuppressedUnclosed = true
	}

	contentStart := 0
	if bytes.HasPrefix(src, []byte{0xEF, 0xBB, 0xBF}) {
		contentStart = 3
	}
	page.LinkInsert = headEnd
	if page.LinkInsert < 0 {
		page.LinkInsert = bodyAfter
	}
	if page.LinkInsert < 0 {
		page.LinkInsert = contentStart
	}
	page.ScriptInsert = bodyEnd
	if page.ScriptInsert < 0 {
		page.ScriptInsert = htmlEnd
	}
	if page.ScriptInsert < 0 {
		page.ScriptInsert = len(src)
	}
	return page
}

// DiscoverWithin re-runs candidate discovery on the inner tokens of a
// wrapper candidate that no unwrapper matched. A site component that happens
// to use a trigger class must not swallow a genuine code block nested inside
// it. Only Candidates and Suppressed are meaningful in the returned Page;
// the caller applies this one level deep and never recurses further.
func DiscoverWithin(inner []BufferedToken) Page {
	return Discover(nil, inner)
}

// isWrapperTrigger reports whether a start tag opens a wrapper candidate:
// a div or figure whose class list carries the highlight token, the
// highlighter-rouge token, or a highlight- prefixed token. The bare chroma
// token is deliberately not a trigger; no generator roots a block there and
// unrelated widget libraries use the same class name.
func isWrapperTrigger(tok BufferedToken) bool {
	if tok.Type != html.StartTagToken || (tok.Data != "div" && tok.Data != "figure") {
		return false
	}
	if hasClass(tok.Attrs, "highlight") || hasClass(tok.Attrs, "highlighter-rouge") {
		return true
	}
	_, ok := classWithPrefix(tok.Attrs, "highlight-")
	return ok
}

// makeCandidate builds a candidate for the element opened at index i. The
// returned index is where the outer walk resumes, and closed reports whether
// the element's end tag was found. Unlike subtree, an unclosed element is
// detected explicitly rather than yielding the rest of the window, because a
// same named nested element could otherwise fake a balanced close.
func makeCandidate(tokens []BufferedToken, i int, kind CandidateKind) (Candidate, int, bool) {
	root := tokens[i]
	j, closed := findMatchingEnd(tokens, i)
	if !closed {
		return Candidate{
			Kind:      kind,
			Tokens:    tokens[i : i+1],
			ByteStart: root.Offset,
			ByteEnd:   root.Offset + len(root.Raw),
			Malformed: true,
		}, i + 1, false
	}
	end := tokens[j]
	return Candidate{
		Kind:      kind,
		Tokens:    tokens[i : j+1],
		ByteStart: root.Offset,
		ByteEnd:   end.Offset + len(end.Raw),
	}, j + 1, true
}

// findMatchingEnd locates the end tag matching the element opened at index
// i, tracking depth on that tag name only.
func findMatchingEnd(tokens []BufferedToken, i int) (int, bool) {
	tag := tokens[i].Data
	depth := 1
	for j := i + 1; j < len(tokens); j++ {
		switch tokens[j].Type {
		case html.StartTagToken:
			if tokens[j].Data == tag {
				depth++
			}
		case html.EndTagToken:
			if tokens[j].Data == tag {
				depth--
				if depth == 0 {
					return j, true
				}
			}
		}
	}
	return len(tokens), false
}
