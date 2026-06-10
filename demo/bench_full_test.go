package main

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/frostybee/kazari"
	kazarimd "github.com/frostybee/kazari/goldmark"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
	"github.com/yuin/goldmark"
)

func TestFullPipeline(t *testing.T) {
	ctx := context.Background()

	start := time.Now()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	t.Logf("nuri.New: %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}
	defer hl.Close(ctx)

	start = time.Now()
	engine := kazari.New(
		kazari.WithHighlighter(kazari.NewNuriHighlighter(ctx, hl)),
		kazari.WithThemes("github-light", "github-dark"),
		kazari.WithMinify(false),
	)
	t.Logf("kazari.New: %v", time.Since(start))

	goCode := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}"

	start = time.Now()
	_, err = engine.Render(goCode, kazari.Options{Lang: "go", Title: "main.go"})
	t.Logf("First Render (go): %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	jsCode := "const x = 1;\nconsole.log(x);"
	start = time.Now()
	_, err = engine.Render(jsCode, kazari.Options{Lang: "javascript"})
	t.Logf("Second Render (js): %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	bashCode := "echo hello\nls -la"
	start = time.Now()
	_, err = engine.Render(bashCode, kazari.Options{Lang: "bash"})
	t.Logf("Third Render (bash): %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	start = time.Now()
	_, err = engine.Render(goCode, kazari.Options{Lang: "go", Title: "cached.go"})
	t.Logf("Fourth Render (go again): %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	start = time.Now()
	collapseEngine := kazari.New(
		kazari.WithHighlighter(kazari.NewNuriHighlighter(ctx, hl)),
		kazari.WithThemes("github-light", "github-dark"),
		kazari.WithMinify(false),
		kazari.WithCollapsible(kazari.CollapsibleConfig{
			LineThreshold:    12,
			PreviewLines:     6,
			DefaultCollapsed: true,
			PreserveIndent:   true,
		}),
	)
	t.Logf("kazari.New (collapse): %v", time.Since(start))

	longCode := "package main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n\t\"log\"\n\t\"encoding/json\"\n\t\"os\"\n)\n\ntype Server struct {\n\taddr    string\n\thandler http.Handler\n\tlogger  *log.Logger\n}\n\nfunc NewServer(addr string) *Server {\n\treturn &Server{\n\t\taddr:    addr,\n\t\thandler: http.DefaultServeMux,\n\t\tlogger:  log.New(os.Stdout, \"[server] \", log.LstdFlags),\n\t}\n}\n"
	start = time.Now()
	_, err = collapseEngine.Render(longCode, kazari.Options{Lang: "go", Title: "server.go"})
	t.Logf("Collapse Render: %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	start = time.Now()
	md := goldmark.New(
		goldmark.WithExtensions(
			kazarimd.New(collapseEngine),
			kazarimd.CodeGroups(collapseEngine),
		),
	)
	mdSrc := []byte(":::code-group\n\n```go title=\"main.go\"\npackage main\n```\n\n```python title=\"main.py\"\nprint(\"hello\")\n```\n\n:::\n")
	var buf bytes.Buffer
	err = md.Convert(mdSrc, &buf)
	t.Logf("Goldmark convert: %v", time.Since(start))
	if err != nil {
		t.Fatal(err)
	}

	start = time.Now()
	_ = collapseEngine.CSS()
	t.Logf("CSS(): %v", time.Since(start))

	start = time.Now()
	_ = collapseEngine.JS()
	t.Logf("JS(): %v", time.Since(start))
}
