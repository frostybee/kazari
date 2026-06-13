package kazarimd

import (
	"strings"

	"github.com/frostybee/kazari"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// codeBlockRenderer renders fenced code blocks through the Kazari engine.
type codeBlockRenderer struct {
	engine *kazari.Engine
}

func newCodeBlockRenderer(engine *kazari.Engine) renderer.NodeRenderer {
	return &codeBlockRenderer{engine: engine}
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *codeBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
}

func (r *codeBlockRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	fcb := node.(*ast.FencedCodeBlock)
	code := extractCode(fcb, source)
	lang, metaRest := extractMeta(fcb, source)

	metaStr := lang
	if metaRest != "" {
		metaStr += " " + metaRest
	}

	html, err := r.engine.RenderWithMeta(code, metaStr)
	if err != nil {
		return ast.WalkContinue, err
	}

	_, _ = w.WriteString(html)
	return ast.WalkSkipChildren, nil
}

// extractCode reads the code content from a fenced code block AST node.
func extractCode(node *ast.FencedCodeBlock, source []byte) string {
	var sb strings.Builder
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)
		sb.Write(seg.Value(source))
	}
	code := sb.String()
	code = strings.TrimRight(code, "\n")
	return code
}

// extractMeta splits the info string into language and meta remainder.
func extractMeta(node *ast.FencedCodeBlock, source []byte) (lang string, metaRest string) {
	info := node.Info
	if info == nil {
		return "", ""
	}
	infoStr := string(info.Text(source))
	parts := strings.SplitN(strings.TrimSpace(infoStr), " ", 2)
	lang = parts[0]
	if len(parts) > 1 {
		metaRest = parts[1]
	}
	return lang, metaRest
}
