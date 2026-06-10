package kazarimd

import (
	"github.com/frostybee/kazari"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Extension is a Goldmark extender that renders fenced code blocks through Kazari.
type Extension struct {
	engine *kazari.Engine
}

// New creates a Goldmark extender that renders fenced code blocks through Kazari.
func New(engine *kazari.Engine) goldmark.Extender {
	return &Extension{engine: engine}
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newCodeBlockRenderer(e.engine), 100),
		),
	)
}
