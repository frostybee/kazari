package minify

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

var m *minify.M

func init() {
	m = minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)
}

// CSS minifies CSS content. Returns input unchanged on error.
func CSS(input string) string {
	if input == "" {
		return ""
	}
	result, err := m.String("text/css", input)
	if err != nil {
		return input
	}
	return result
}

// JS minifies JavaScript content. Returns input unchanged on error.
func JS(input string) string {
	if input == "" {
		return ""
	}
	result, err := m.String("application/javascript", input)
	if err != nil {
		return input
	}
	return result
}
