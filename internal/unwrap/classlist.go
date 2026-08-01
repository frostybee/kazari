package unwrap

import (
	"strings"

	"golang.org/x/net/html"
)

// attrValue returns the value of the named attribute and whether it exists.
func attrValue(attrs []html.Attribute, key string) (string, bool) {
	for _, a := range attrs {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}

// classList returns the class attribute split into individual tokens.
func classList(attrs []html.Attribute) []string {
	v, ok := attrValue(attrs, "class")
	if !ok {
		return nil
	}
	return strings.Fields(v)
}

// hasClass reports whether the class list contains name as an exact token.
// Substring matching is never acceptable here: "language-css" must not
// satisfy a check for "language-c".
func hasClass(attrs []html.Attribute, name string) bool {
	for _, c := range classList(attrs) {
		if c == name {
			return true
		}
	}
	return false
}

// classWithPrefix returns the first class token carrying the prefix, along
// with whether one was found. The full token is returned, prefix included.
func classWithPrefix(attrs []html.Attribute, prefix string) (string, bool) {
	for _, c := range classList(attrs) {
		if strings.HasPrefix(c, prefix) && len(c) > len(prefix) {
			return c, true
		}
	}
	return "", false
}
