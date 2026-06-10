package kazarimd

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari"
	"github.com/frostybee/kazari/internal/frame"
	"github.com/frostybee/kazari/internal/meta"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// --- AST Node ---

// KindCodeGroup is the AST node kind for code group containers.
var KindCodeGroup = ast.NewNodeKind("CodeGroup")

// CodeGroupNode is a container AST node wrapping multiple fenced code blocks.
type CodeGroupNode struct {
	ast.BaseBlock
}

func (n *CodeGroupNode) Kind() ast.NodeKind {
	return KindCodeGroup
}

func (n *CodeGroupNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// --- Block Parser ---

type codeGroupParser struct{}

func (p *codeGroupParser) Trigger() []byte {
	return []byte{':'}
}

func (p *codeGroupParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if !isCodeGroupOpener(line) {
		return nil, parser.NoChildren
	}
	reader.Advance(len(line))
	return &CodeGroupNode{}, parser.HasChildren
}

func (p *codeGroupParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	if isCodeGroupCloser(line) {
		reader.Advance(len(line))
		return parser.Close
	}
	return parser.Continue | parser.HasChildren
}

func (p *codeGroupParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}

func (p *codeGroupParser) CanInterruptParagraph() bool { return false }

func (p *codeGroupParser) CanAcceptIndentedLine() bool { return false }

func isCodeGroupOpener(line []byte) bool {
	return bytes.Equal(bytes.TrimSpace(line), []byte(":::code-group"))
}

func isCodeGroupCloser(line []byte) bool {
	trimmed := bytes.TrimSpace(line)
	return bytes.Equal(trimmed, []byte(":::"))
}

// --- Code Group Renderer ---

type codeGroupRenderer struct {
	engine *kazari.Engine
}

func (r *codeGroupRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCodeGroup, r.renderCodeGroup)
}

func (r *codeGroupRenderer) renderCodeGroup(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	var blocks []*ast.FencedCodeBlock
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if fcb, ok := child.(*ast.FencedCodeBlock); ok {
			blocks = append(blocks, fcb)
		}
	}

	if len(blocks) == 0 {
		return ast.WalkSkipChildren, nil
	}

	labels := make([]string, len(blocks))
	for i, block := range blocks {
		labels[i] = deriveTabLabel(block, source)
	}

	w.WriteString(`<div class="kazari-code kz-group">`)

	// Tab bar.
	w.WriteString(`<div class="kz-group-tabs" role="tablist">`)
	for i, label := range labels {
		selected := "false"
		tabindex := "-1"
		if i == 0 {
			selected = "true"
			tabindex = "0"
		}
		w.WriteString(fmt.Sprintf(
			`<button role="tab" aria-selected="%s" tabindex="%s">%s</button>`,
			selected, tabindex, html.EscapeString(label),
		))
	}
	w.WriteString(`</div>`)

	// Panels.
	w.WriteString(`<div class="kz-group-panels">`)
	for i, block := range blocks {
		code := extractCode(block, source)
		lang, metaRest := extractMeta(block, source)
		metaStr := lang
		if metaRest != "" {
			metaStr += " " + metaRest
		}

		rendered, err := r.engine.RenderWithMeta(code, metaStr)
		if err != nil {
			rendered = fmt.Sprintf("<pre><code>%s</code></pre>", html.EscapeString(code))
		}

		if i == 0 {
			w.WriteString(fmt.Sprintf(`<div role="tabpanel">%s</div>`, rendered))
		} else {
			w.WriteString(fmt.Sprintf(`<div role="tabpanel" hidden>%s</div>`, rendered))
		}
	}
	w.WriteString(`</div>`)

	w.WriteString(`</div>`)

	return ast.WalkSkipChildren, nil
}

func deriveTabLabel(block *ast.FencedCodeBlock, source []byte) string {
	lang, metaRest := extractMeta(block, source)

	// 1. Check for explicit title in meta string.
	if metaRest != "" {
		parsed := meta.Parse(lang + " " + metaRest)
		if parsed.BlockOptions.Title != "" {
			return parsed.BlockOptions.Title
		}
	}

	// 2. Check for file name in code.
	code := extractCode(block, source)
	if title, _ := frame.ExtractFileName(code, lang); title != "" {
		return title
	}

	// 3. Fallback to language name.
	if lang != "" {
		return strings.ToUpper(lang[:1]) + lang[1:]
	}
	return "Code"
}

// --- Extension ---

// CodeGroupExtension is a Goldmark extender for code group containers.
type CodeGroupExtension struct {
	engine *kazari.Engine
}

// CodeGroups creates a Goldmark extender for :::code-group containers.
func CodeGroups(engine *kazari.Engine) goldmark.Extender {
	engine.EnableCodeGroups()
	return &CodeGroupExtension{engine: engine}
}

// Extend implements goldmark.Extender.
func (e *CodeGroupExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&codeGroupParser{}, 50),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&codeGroupRenderer{engine: e.engine}, 100),
		),
	)
}
