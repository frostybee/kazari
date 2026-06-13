package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/frostybee/kazari"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func main() {
	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		log.Fatalf("nuri.New: %v", err)
	}
	defer hl.Close(ctx)

	engine := kazari.New(
		kazari.WithHighlighter(kazarinuri.New(ctx, hl)),
		kazari.WithThemes("github-light", "github-dark"),
		kazari.WithCollapsible(kazari.CollapsibleConfig{LineThreshold: 15, PreviewLines: 8}),
	)

	// Case 1: the filename comment is extracted into the title, so the
	// collapse range {2-4} must cover the import block as numbered after
	// extraction, lines 2 to 4 (blank, import, blank).
	collapseCode := "// main.go\n" +
		"package main\n" +
		"\n" +
		"import \"fmt\"\n" +
		"\n" +
		"func main() {\n" +
		"\tfmt.Println(\"one\")\n" +
		"\tfmt.Println(\"two\")\n" +
		"}"
	block1, err := engine.RenderWithMeta(collapseCode, "go showLineNumbers collapse={2-4}")
	if err != nil {
		log.Fatalf("case 1: %v", err)
	}

	// Case 2: overlapping inline markers. The higher priority ins claims
	// "name string"; the overlapping mark "string) error" must keep its
	// non overlapping tail ") error".
	inlineCode := "func greet(name string) error {\n" +
		"\treturn nil\n" +
		"}"
	block2, err := engine.RenderWithMeta(inlineCode, `go ins="name string" "string) error"`)
	if err != nil {
		log.Fatalf("case 2: %v", err)
	}

	// Case 3: bright background SGR codes 101, 102, 104 next to a standard
	// background 41 for comparison.
	ansiCode := "\x1b[101m bright red bg \x1b[0m \x1b[102m bright green bg \x1b[0m \x1b[104m bright blue bg \x1b[0m\n" +
		"normal \x1b[41m standard red bg \x1b[0m"
	block3, err := engine.RenderWithMeta(ansiCode, "ansi")
	if err != nil {
		log.Fatalf("case 3: %v", err)
	}

	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Kazari Visual Check</title>
<style>%s</style>
<style>body { max-width: 720px; margin: 2rem auto; font-family: system-ui, sans-serif; }</style>
</head>
<body>
<h2>1. Collapse range after filename extraction (collapse={2-4})</h2>
%s
<h2>2. Overlapping inline markers (ins="name string" overlapping mark "string) error")</h2>
%s
<h2>3. ANSI bright backgrounds (SGR 101/102/104)</h2>
%s
<script>%s</script>
</body>
</html>`, engine.CSS(), block1, block2, block3, engine.JS())

	if err := os.WriteFile("demo/nuri/visual-check.html", []byte(page), 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Println("Written: demo/nuri/visual-check.html")
}
