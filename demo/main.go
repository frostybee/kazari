package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/frostybee/kazari"
	kazarimd "github.com/frostybee/kazari/goldmark"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
	"github.com/yuin/goldmark"
)

func main() {
	ctx := context.Background()

	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		log.Fatalf("nuri.New: %v", err)
	}
	defer hl.Close(ctx)

	engine := kazari.New(
		kazari.WithHighlighter(kazari.NewNuriHighlighter(ctx, hl)),
		kazari.WithThemes("github-light", "github-dark"),
		kazari.WithMinify(false),
	)

	// Editor frame with title
	goCode := `package main

import "fmt"

func main() {
	name := "Kazari"
	fmt.Printf("Hello, %s!\n", name)
}`

	// Editor frame with file name extraction
	jsCode := `// src/greet.js
const greet = (name) => {
  console.log("Hello, " + name + "!");
  return { greeting: name, time: Date.now() };
};`

	// Terminal frame (auto-detected)
	bashCode := `npm install kazari
go build ./...
echo "Done!"`

	// Terminal frame with title
	psCode := `Write-Output "This one has a title!"`

	// Line numbers enabled
	lnCode := `package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: greet <name>")
		os.Exit(1)
	}
	name := strings.Join(args, " ")
	fmt.Printf("Hello, %s!\n", name)
}`

	// Line numbers with custom start
	lnStartCode := `	name := strings.Join(args, " ")
	fmt.Printf("Hello, %s!\n", name)
}`

	// Frame none
	plainCode := `Just some plain text
with no frame wrapper.`

	goHTML, err := engine.Render(goCode, kazari.Options{Lang: "go", Title: "main.go"})
	if err != nil {
		log.Fatalf("Render go: %v", err)
	}

	jsHTML, err := engine.Render(jsCode, kazari.Options{Lang: "javascript"})
	if err != nil {
		log.Fatalf("Render js: %v", err)
	}

	bashHTML, err := engine.Render(bashCode, kazari.Options{Lang: "bash"})
	if err != nil {
		log.Fatalf("Render bash: %v", err)
	}

	psHTML, err := engine.Render(psCode, kazari.Options{Lang: "powershell", Title: "PowerShell terminal example"})
	if err != nil {
		log.Fatalf("Render ps: %v", err)
	}

	lnEnabled := true
	lnHTML, err := engine.Render(lnCode, kazari.Options{Lang: "go", Title: "main.go", LineNumbers: &lnEnabled})
	if err != nil {
		log.Fatalf("Render ln: %v", err)
	}

	lnStart := 15
	lnStartHTML, err := engine.Render(lnStartCode, kazari.Options{Lang: "go", Title: "main.go (lines 15-17)", LineNumbers: &lnEnabled, StartLineNumber: &lnStart})
	if err != nil {
		log.Fatalf("Render lnStart: %v", err)
	}

	frameNone := kazari.FrameNone
	plainHTML, err := engine.Render(plainCode, kazari.Options{Lang: "text", Frame: &frameNone})
	if err != nil {
		log.Fatalf("Render plain: %v", err)
	}

	// Line markers: mark, ins, del
	markerCode := `package main

import "fmt"

func oldGreet(name string) {
	fmt.Println("Hi,", name)
}

func newGreet(name string) {
	fmt.Printf("Hello, %s! Welcome!\n", name)
}

func main() {
	newGreet("Kazari")
}`

	markerHTML, err := engine.RenderWithMeta(markerCode, `go title="diff.go" showLineNumbers {3} del={5-7} ins={9-11}`)
	if err != nil {
		log.Fatalf("Render markers: %v", err)
	}

	// Labeled range (EC reference example)
	labelCode := `<button
  role="button"
  {...props}

  value={value}
  className={buttonClassName}

  disabled={disabled}
  active={active}
>

  {children &&
    !active &&
    (typeof children === 'string'
      ? <span>{children}</span>
      : children)}
</button>`

	labelHTML, err := engine.RenderWithMeta(labelCode, `jsx title="labeled-line-markers.jsx" showLineNumbers {"1. Provide the value prop here:":4-6} del={"2. Remove the disabled and active states:":7-10} ins={"3. Add this to render the children inside the button:":11-16}`)
	if err != nil {
		log.Fatalf("Render labels: %v", err)
	}

	// Focus lines
	focusCode := `func process(items []string) error {
	for _, item := range items {
		if err := validate(item); err != nil {
			return fmt.Errorf("invalid: %w", err)
		}
		store(item)
	}
	return nil
}`

	focusHTML, err := engine.RenderWithMeta(focusCode, `go title="process.go" showLineNumbers focus={3-5}`)
	if err != nil {
		log.Fatalf("Render focus: %v", err)
	}

	// Inline markers
	inlineCode := `import { useState, useEffect } from 'react';

function App() {
  const [count, setCount] = useState(0);
  useEffect(() => {
    document.title = count;
  }, [count]);
  return <button onClick={() => setCount(count + 1)}>{count}</button>;
}`

	inlineHTML, err := engine.RenderWithMeta(inlineCode, `javascript title="App.jsx" showLineNumbers "useState" ins="useEffect"`)
	if err != nil {
		log.Fatalf("Render inline: %v", err)
	}

	// Single-quote inline markers
	sqCode := `import { useState, useEffect, useRef } from 'react';
const [name, setName] = useState('');
const ref = useRef(null);`

	sqHTML, err := engine.RenderWithMeta(sqCode, `javascript title="single-quotes.jsx" showLineNumbers 'useState' ins='useRef' del='useEffect'`)
	if err != nil {
		log.Fatalf("Render single-quote: %v", err)
	}

	// Combined: line markers + inline + focus
	combinedCode := `func main() {
	db := connect()
	defer db.Close()

	users, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		fmt.Println(u.Name)
	}
}`

	combinedHTML, err := engine.RenderWithMeta(combinedCode, `go title="combined.go" showLineNumbers {4-5} ins={10-12} del={6-8} "db" focus={4-5,10-12}`)
	if err != nil {
		log.Fatalf("Render combined: %v", err)
	}

	// --- Collapsible examples ---

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

	// Threshold-based collapse (long block auto-collapses)
	thresholdCode := `package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"os"
)

type Server struct {
	addr    string
	handler http.Handler
	logger  *log.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		handler: http.DefaultServeMux,
		logger:  log.New(os.Stdout, "[server] ", log.LstdFlags),
	}
}

func (s *Server) Start() error {
	s.logger.Printf("Starting server on %s", s.addr)
	return http.ListenAndServe(s.addr, s.handler)
}`

	thresholdHTML, err := collapseEngine.Render(thresholdCode, kazari.Options{Lang: "go", Title: "server.go"})
	if err != nil {
		log.Fatalf("Render threshold: %v", err)
	}

	// Range-based collapse (specific sections collapsed)
	rangeCode := `package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: calc <expr>")
		os.Exit(1)
	}

	result := evaluate(strings.Join(args, " "))
	fmt.Printf("= %v\n", result)
}

func evaluate(expr string) float64 {
	val, _ := strconv.ParseFloat(expr, 64)
	return val
}`

	rangeHTML, err := collapseEngine.RenderWithMeta(rangeCode, `go title="calc.go" showLineNumbers collapse={3-8}`)
	if err != nil {
		log.Fatalf("Render range: %v", err)
	}

	// Range-based with multiple ranges
	multiRangeCode := `package api

import (
	"encoding/json"
	"net/http"
	"log"
)

type Response struct {
	Status  int
	Message string
	Data    interface{}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	data := fetchData(r.URL.Query())
	json.NewEncoder(w).Encode(Response{
		Status:  200,
		Message: "OK",
		Data:    data,
	})
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)
	result := processData(body)
	json.NewEncoder(w).Encode(Response{
		Status:  201,
		Message: "Created",
		Data:    result,
	})
}`

	multiRangeHTML, err := collapseEngine.RenderWithMeta(multiRangeCode, `go title="api.go" showLineNumbers collapse={3-7,9-13}`)
	if err != nil {
		log.Fatalf("Render multiRange: %v", err)
	}

	// Threshold + markers: gap indicators between non-contiguous segments
	gapCode := `package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"os"
	"strings"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/data", dataHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"status": "success",
		"count":  42,
	}
	json.NewEncoder(w).Encode(data)
}`

	gapHTML, err := collapseEngine.Render(gapCode, kazari.Options{
		Lang:        "go",
		Title:       "server.go",
		LineNumbers: &lnEnabled,
		LineMarkers: []kazari.LineMarker{
			{Type: kazari.MarkerIns, Lines: []kazari.Range{{Start: 10, End: 11}}},
		},
	})
	if err != nil {
		log.Fatalf("Render gap: %v", err)
	}

	// Collapsible-start style (re-collapsible, summary above)
	csStartHTML, err := collapseEngine.RenderWithMeta(rangeCode, `go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-start"`)
	if err != nil {
		log.Fatalf("Render csStart: %v", err)
	}

	// Collapsible-end style (re-collapsible, summary below)
	csEndHTML, err := collapseEngine.RenderWithMeta(rangeCode, `go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-end"`)
	if err != nil {
		log.Fatalf("Render csEnd: %v", err)
	}

	// Collapsible-auto style (auto-selects start/end based on position)
	csAutoHTML, err := collapseEngine.RenderWithMeta(rangeCode, `go title="calc.go" showLineNumbers collapse={3-8,20-24} collapseStyle="collapsible-auto"`)
	if err != nil {
		log.Fatalf("Render csAuto: %v", err)
	}

	// --- Code group example (via Goldmark) ---
	// Use collapseEngine for code groups too (single engine = single CSS/JS output).

	codeGroupMD := []byte(":::code-group\n\n```go title=\"main.go\"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello from Go!\")\n}\n```\n\n```python title=\"main.py\"\ndef main():\n    print(\"Hello from Python!\")\n\nif __name__ == \"__main__\":\n    main()\n```\n\n```javascript title=\"index.js\"\nfunction main() {\n  console.log(\"Hello from JavaScript!\");\n}\n\nmain();\n```\n\n:::\n")

	md := goldmark.New(
		goldmark.WithExtensions(
			kazarimd.New(collapseEngine),
			kazarimd.CodeGroups(collapseEngine),
		),
	)
	var codeGroupBuf bytes.Buffer
	if err := md.Convert(codeGroupMD, &codeGroupBuf); err != nil {
		log.Fatalf("goldmark.Convert: %v", err)
	}
	codeGroupHTML := codeGroupBuf.String()

	css := collapseEngine.CSS()
	js := collapseEngine.JS()

	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Kazari Demo</title>
<style>
%s
body {
  font-family: system-ui, sans-serif;
  max-width: 800px;
  margin: 2rem auto;
  padding: 0 1rem;
  background: #f8f9fa;
  color: #1a1a1a;
}
h1 { font-size: 1.5rem; }
h2 { font-size: 1.1rem; margin-top: 2rem; color: #555; }
.dark body { background: #1a1a2e; color: #e0e0e0; }
.dark h2 { color: #aaa; }
</style>
</head>
<body>
<h1>Kazari Demo</h1>
<p>Dual-theme code blocks with frames, copy buttons, and dark mode.</p>
<button onclick="document.documentElement.classList.toggle('dark')">Toggle Dark Mode</button>

<h2>Editor Frame (explicit title)</h2>
%s

<h2>Editor Frame (file name extracted from comment)</h2>
%s

<h2>Terminal Frame (auto-detected from language)</h2>
%s

<h2>Terminal Frame (with title)</h2>
%s

<h2>Line Numbers</h2>
%s

<h2>Line Numbers (custom start)</h2>
%s

<h2>No Frame</h2>
%s

<h2>Line Markers (mark, ins, del)</h2>
%s

<h2>Labeled Range</h2>
%s

<h2>Focus Lines</h2>
%s

<h2>Inline Markers</h2>
%s

<h2>Inline Markers (single quotes)</h2>
%s

<h2>Combined (markers + inline + focus)</h2>
%s

<h2>Collapsible: Threshold-based (auto-collapses long blocks)</h2>
%s

<h2>Collapsible: Range-based (imports collapsed)</h2>
%s

<h2>Collapsible: Multiple ranges</h2>
%s

<h2>Collapsible: Threshold + markers (gap indicators)</h2>
%s

<h2>Collapsible: collapsible-start (re-collapsible, summary above)</h2>
%s

<h2>Collapsible: collapsible-end (re-collapsible, summary below)</h2>
%s

<h2>Collapsible: collapsible-auto (auto start/end based on position)</h2>
%s

<h2>Code Group (tabbed code blocks via Goldmark)</h2>
%s

<script type="module">
%s
</script>
</body>
</html>`, css, goHTML, jsHTML, bashHTML, psHTML, lnHTML, lnStartHTML, plainHTML,
		markerHTML, labelHTML, focusHTML, inlineHTML, sqHTML, combinedHTML,
		thresholdHTML, rangeHTML, multiRangeHTML, gapHTML,
		csStartHTML, csEndHTML, csAutoHTML, codeGroupHTML, js)

	if err := os.WriteFile("showcase.html", []byte(page), 0644); err != nil {
		log.Fatalf("WriteFile: %v", err)
	}
	fmt.Println("Written: showcase.html")
}
