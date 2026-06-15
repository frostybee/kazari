package showcase

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/frostybee/kazari"
	kazarimd "github.com/frostybee/kazari/goldmark"
	"github.com/yuin/goldmark"
)

type catalogBuilder struct {
	highlighter   kazari.Highlighter
	kazariOptions []kazari.Option
	err           error
}

func (b *catalogBuilder) engine(options ...kazari.Option) *kazari.Engine {
	all := make([]kazari.Option, 0, len(b.kazariOptions)+2+len(options))
	all = append(all, b.kazariOptions...)
	all = append(all,
		kazari.WithHighlighter(b.highlighter),
		kazari.WithMinify(false),
	)
	all = append(all, options...)
	return kazari.New(all...)
}

func (b *catalogBuilder) render(engine *kazari.Engine, code string, options kazari.Options) template.HTML {
	if b.err != nil {
		return ""
	}
	html, err := engine.Render(code, options)
	if err != nil {
		b.err = err
		return ""
	}
	return template.HTML(html)
}

func (b *catalogBuilder) renderMeta(engine *kazari.Engine, code, meta string) template.HTML {
	if b.err != nil {
		return ""
	}
	html, err := engine.RenderWithMeta(code, meta)
	if err != nil {
		b.err = err
		return ""
	}
	return template.HTML(html)
}

func (b *catalogBuilder) renderMarkdown(markdown goldmark.Markdown, source string) template.HTML {
	if b.err != nil {
		return ""
	}
	var output bytes.Buffer
	if err := markdown.Convert([]byte(source), &output); err != nil {
		b.err = err
		return ""
	}
	return template.HTML(output.String())
}

func recipe(label, code string) Recipe {
	return Recipe{Label: label, Code: code}
}

var exampleDescriptions = map[string]string{
	"editor-explicit-title": "Adds familiar editor chrome and an explicit file name above the highlighted code.",
	"editor-comment-title":  "Extracts the file name from a leading source comment when no title is provided.",
	"terminal-auto":         "Automatically presents shell languages in a terminal-style frame instead of an editor frame.",
	"terminal-title":        "Adds a descriptive session or command label to the terminal title bar.",
	"terminal-minimal":      "Uses a compact engine-level dot style for a quieter terminal title bar.",
	"no-frame":              "Removes the surrounding window chrome while preserving syntax highlighting and code structure.",
	"line-numbers":          "Adds a numbered gutter so readers can reference specific lines precisely.",
	"line-numbers-start":    "Continues numbering from the source file's original location when showing an excerpt.",
	"word-wrap":             "Wraps long lines within the block while preserving indentation for readable continuations.",
	"word-wrap-no-preserve": "Wrapped continuations start at the left edge instead of aligning with the original indentation.",
	"word-wrap-hanging":     "Adds a fixed hanging indent to wrapped continuations so new logical lines are easy to spot.",
	"line-markers":          "Distinguishes highlighted, inserted, and deleted lines with clear full-line treatments.",
	"labeled-range":         "Attaches explanatory labels to marked line ranges so each change can carry context.",
	"labeled-range-no-ln":   "Labeled ranges without line numbers place the badge flush at the left edge of the block.",
	"labeled-range-numbers": "Short numeric labels act as compact reference badges on the first line of each range.",
	"focus-lines":           "Keeps selected lines prominent while dimming the surrounding code for emphasis.",
	"inline-markers":        "Highlights matching text inside a line without losing the underlying syntax colors.",
	"inline-markers-single": "Uses single-quoted marker expressions when the highlighted text or meta string needs simpler escaping.",
	"combined-markers":      "Layers line markers, inline matches, and focused ranges in the same code block.",
	"collapse-threshold":    "Automatically collapses long blocks when they exceed the engine's configured line threshold.",
	"collapse-range":        "Hides a selected line range behind an expandable summary to keep the main logic visible.",
	"collapse-multiple":     "Collapses multiple independent ranges so distant supporting sections stay compact.",
	"collapse-gaps":         "Shows omitted regions as gap indicators while preserving important marked lines around them.",
	"collapse-start":        "Places the expansion summary above a range that can be collapsed again after opening.",
	"collapse-end":          "Places the expansion summary below a range that can be collapsed again after opening.",
	"collapse-auto":         "Chooses the summary position from the collapsed range's location within the code block.",
	"mermaid":               "Passes Mermaid source through unchanged so a diagram renderer can process it later.",
	"regex-markers":         "Marks text that matches regular expressions, including inserted and deleted match styles.",
	"regex-capture":         "Highlights only the captured subgroup of a regular-expression match instead of the full match.",
	"hybrid-diff":           "Combines diff prefixes with syntax highlighting from the underlying source language.",
	"code-group":            "Presents related language examples as an accessible tabbed group generated from Markdown.",
	"ansi":                  "Converts ANSI SGR escape sequences into styled terminal colors and text treatments.",
	"code-group-sync":       "Synchronizes matching tabs across separate code groups that share the same key.",
	"theme-override":        "Selects an alternate theme for one block without changing the rest of the page.",
	"theme-override-dual":   "Gives one block its own light and dark themes, switched by the page's dark mode toggle. Inverted here on purpose: dracula in light mode, github-light in dark mode.",
	"theme-customizer":      "Adjusts resolved theme colors through a callback before Kazari generates the CSS variables.",
	"theme-adjustments":     "Applies an OKLCH tint to generated theme colors while preserving their visual relationships.",
	"scoped-css":            "Emits theme variables beneath a custom selector so Kazari styles stay inside a chosen container.",
	"locale-french":         "Localizes built-in copy and fullscreen controls through the engine's locale setting.",
	"file-icons":            "Resolves a custom icon from each title's file extension and places it in the frame toolbar.",
}

func joinHTML(parts ...template.HTML) template.HTML {
	var output strings.Builder
	for _, part := range parts {
		output.WriteString(string(part))
	}
	return template.HTML(output.String())
}

func buildCatalog(highlighter kazari.Highlighter, kazariOptions []kazari.Option) ([]Category, string, string, error) {
	b := &catalogBuilder{highlighter: highlighter, kazariOptions: kazariOptions}
	engine := b.engine()

	goCode := `package main

import "fmt"

func main() {
	name := "Kazari"
	fmt.Printf("Hello, %s!\n", name)
}`
	jsCode := `// src/greet.js
const greet = (name) => {
  console.log("Hello, " + name + "!");
  return { greeting: name, time: Date.now() };
};`
	bashCode := `npm install kazari
go build ./...
echo "Done!"`
	psCode := `Get-ChildItem -Path ./dist -Recurse | Measure-Object -Property Length -Sum`
	plainCode := `Just some plain text
with no frame wrapper.`

	dotsEngine := b.engine(kazari.WithTerminalDotStyle(kazari.DotsMinimal))
	frameNone := kazari.FrameNone

	frames := Category{
		ID:          "frames",
		Title:       "Frames",
		Description: "Editor, terminal, and unframed presentation styles.",
		Examples: []Example{
			{
				ID:       "editor-explicit-title",
				Title:    "Editor Frame (explicit title)",
				NavTitle: "Editor title",
				HTML:     b.render(engine, goCode, kazari.Options{Lang: "go", Title: "main.go"}),
				Recipes: []Recipe{
					recipe("Meta", `go title="main.go"`),
					recipe("Go", `html, err := engine.Render(code, kazari.Options{
	Lang:  "go",
	Title: "main.go",
})`),
				},
			},
			{
				ID:       "editor-comment-title",
				Title:    "Editor Frame (file name extracted from comment)",
				NavTitle: "Extracted file name",
				HTML:     b.render(engine, jsCode, kazari.Options{Lang: "javascript"}),
				Recipes: []Recipe{
					recipe("Meta", `javascript`),
					recipe("Go", `html, err := engine.Render(code, kazari.Options{Lang: "javascript"})`),
				},
			},
			{
				ID:       "terminal-auto",
				Title:    "Terminal Frame (auto-detected from language)",
				NavTitle: "Terminal detection",
				HTML:     b.render(engine, bashCode, kazari.Options{Lang: "bash"}),
				Recipes: []Recipe{
					recipe("Meta", `bash`),
					recipe("Go", `html, err := engine.Render(code, kazari.Options{Lang: "bash"})`),
				},
			},
			{
				ID:       "terminal-title",
				Title:    "Terminal Frame (with title)",
				NavTitle: "Terminal title",
				HTML:     b.render(engine, psCode, kazari.Options{Lang: "powershell", Title: "PowerShell terminal example"}),
				Recipes: []Recipe{
					recipe("Meta", `powershell title="PowerShell terminal example"`),
					recipe("Go", `html, err := engine.Render(code, kazari.Options{
	Lang:  "powershell",
	Title: "PowerShell terminal example",
})`),
				},
			},
			{
				ID:          "terminal-minimal",
				Title:       "Terminal Frame (minimal dots)",
				NavTitle:    "Minimal dots",
				Description: "Terminal dot style is configured at engine level.",
				HTML:        b.render(dotsEngine, bashCode, kazari.Options{Lang: "bash", Title: "Minimal dots"}),
				Recipes: []Recipe{recipe("Go", `engine := kazari.New(
	kazari.WithHighlighter(highlighter),
	kazari.WithTerminalDotStyle(kazari.DotsMinimal),
)
html, err := engine.Render(code, kazari.Options{Lang: "bash", Title: "Minimal dots"})`)},
			},
			{
				ID:       "no-frame",
				Title:    "No Frame",
				NavTitle: "No frame",
				HTML:     b.render(engine, plainCode, kazari.Options{Lang: "text", Frame: &frameNone}),
				Recipes: []Recipe{
					recipe("Meta", `text frame="none"`),
					recipe("Go", `frame := kazari.FrameNone
html, err := engine.Render(code, kazari.Options{Lang: "text", Frame: &frame})`),
				},
			},
		},
	}

	lineNumbers := true
	startLine := 22
	lineCode := `package main

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
	lineStartCode := "\tresult := evaluate(strings.Join(args, \" \"))\n\tfmt.Printf(\"= %v\\n\", result)\n}"
	wrapCode := `func configure(opts *Options) {
	opts.Logger = log.New(os.Stdout, "[kazari] a deliberately long prefix string that forces this line to wrap inside the demo container", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	opts.Description = "Word wrap keeps long lines visible without horizontal scrolling, and preserved indentation keeps wrapped continuations aligned with the code structure."
}`

	layout := Category{
		ID:          "layout",
		Title:       "Layout",
		Description: "Line numbering and wrapping behavior for different code shapes.",
		Examples: []Example{
			{
				ID:       "line-numbers",
				Title:    "Line Numbers",
				NavTitle: "Line numbers",
				HTML:     b.render(engine, lineCode, kazari.Options{Lang: "go", Title: "main.go", LineNumbers: &lineNumbers}),
				Recipes: []Recipe{
					recipe("Meta", `go title="main.go" showLineNumbers`),
					recipe("Go", `enabled := true
html, err := engine.Render(code, kazari.Options{
	Lang: "go", Title: "main.go", LineNumbers: &enabled,
})`),
				},
			},
			{
				ID:       "line-numbers-start",
				Title:    "Line Numbers (custom start)",
				NavTitle: "Custom start",
				HTML:     b.render(engine, lineStartCode, kazari.Options{Lang: "go", Title: "calc.go (lines 22-24)", LineNumbers: &lineNumbers, StartLineNumber: &startLine}),
				Recipes: []Recipe{
					recipe("Meta", `go title="calc.go (lines 22-24)" showLineNumbers startLineNumber=22`),
					recipe("Go", `enabled, start := true, 22
html, err := engine.Render(code, kazari.Options{
	Lang: "go", Title: "calc.go (lines 22-24)",
	LineNumbers: &enabled, StartLineNumber: &start,
})`),
				},
			},
			{
				ID:       "word-wrap",
				Title:    "Word Wrap (long lines wrap, indent preserved)",
				NavTitle: "Word wrap",
				HTML:     b.renderMeta(engine, wrapCode, `go title="wrap.go" showLineNumbers wrap`),
				Recipes: []Recipe{
					recipe("Meta", `go title="wrap.go" showLineNumbers wrap`),
					recipe("Go", `enabled, wrap := true, true
html, err := engine.Render(code, kazari.Options{
	Lang: "go", Title: "wrap.go", LineNumbers: &enabled, Wrap: &wrap,
})`),
				},
			},
			{
				ID:          "word-wrap-no-preserve",
				Title:       "Word Wrap (preserveIndent=false)",
				NavTitle:    "No preserve indent",
				Description: "Wrapped continuations start at the left edge instead of aligning with the original indentation.",
				HTML:        b.renderMeta(engine, wrapCode, `go title="no-preserve.go" showLineNumbers wrap preserveIndent=false`),
				Recipes: []Recipe{
					recipe("Meta", `go title="no-preserve.go" showLineNumbers wrap preserveIndent=false`),
				},
			},
			{
				ID:          "word-wrap-hanging",
				Title:       "Word Wrap (hangingIndent=4)",
				NavTitle:    "Hanging indent",
				Description: "Adds a fixed 4-character hanging indent to all wrapped continuations, making it easy to spot where a new logical line begins.",
				HTML:        b.renderMeta(engine, wrapCode, `go title="hanging.go" showLineNumbers wrap preserveIndent=false hangingIndent=4`),
				Recipes: []Recipe{
					recipe("Meta", `go title="hanging.go" showLineNumbers wrap preserveIndent=false hangingIndent=4`),
				},
			},
		},
	}

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
	labelCode := `class UserController extends Controller
{
    private UserRepository $users;
    private LoggerInterface $logger;

    public function __construct(
        UserRepository $users,
        LoggerInterface $logger
    ) {
        $this->users = $users;
        $this->logger = $logger;
    }

    public function show(int $id): Response
    {

        $user = $this->users->find($id);
        if ($user === null) {
            throw new NotFoundHttpException();
        }

        return $this->json($user);
    }
}`
	focusCode := `func process(items []string) error {
	for _, item := range items {
		if err := validate(item); err != nil {
			return fmt.Errorf("invalid: %w", err)
		}
		store(item)
	}
	return nil
}`
	inlineCode := `interface CacheEntry<T> {
  key: string;
  value: T;
  expiresAt: number;
}

function getOrSet<T>(cache: Map<string, CacheEntry<T>>, key: string, factory: () => T): T {
  const entry = cache.get(key);
  if (entry && entry.expiresAt > Date.now()) {
    return entry.value;
  }
  const value = factory();
  cache.set(key, { key, value, expiresAt: Date.now() + 3600_000 });
  return value;
}`
	singleQuoteCode := `var query = context.Users
    .Where(u => u.IsActive)
    .OrderBy(u => u.CreatedAt)
    .Select(u => new UserDto(u.Name, u.Email));`
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

	markers := Category{
		ID:          "markers",
		Title:       "Markers and Focus",
		Description: "Call attention to lines, ranges, and inline text without losing syntax highlighting.",
		Examples: []Example{
			metaGoExample(b, engine, "line-markers", "Line Markers (mark, ins, del)", "Line markers", markerCode,
				`go title="diff.go" showLineNumbers {3} del={5-7} ins={9-11}`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`go title="diff.go" showLineNumbers {3} del={5-7} ins={9-11}`+"`"+`)`),
			metaGoExample(b, engine, "labeled-range", "Labeled Range", "Labeled range", labelCode,
				`php title="UserController.php" showLineNumbers {"1. Inject dependencies via constructor:":5-12} del={"2. Remove inline lookup logic:":16-20} ins={"3. Return a JSON response:":21-22}`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, engine, "labeled-range-no-ln", "Labeled Range (no line numbers)", "No line numbers", labelCode,
				`php title="UserController.php" {"1. Inject dependencies via constructor:":5-12} del={"2. Remove inline lookup logic:":16-20} ins={"3. Return a JSON response:":21-22}`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, engine, "labeled-range-numbers", "Labeled Range (numbered)", "Numbered labels", labelCode,
				`php title="UserController.php" {"1":6-9} del={"2":17-19} ins={"3":21-22}`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, engine, "focus-lines", "Focus Lines", "Focus lines", focusCode,
				`go title="process.go" showLineNumbers focus={3-5}`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`go title="process.go" showLineNumbers focus={3-5}`+"`"+`)`),
			metaGoExample(b, engine, "inline-markers", "Inline Markers", "Inline markers", inlineCode,
				`typescript title="cache.ts" showLineNumbers "CacheEntry" ins="factory"`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`typescript title="cache.ts" showLineNumbers "CacheEntry" ins="factory"`+"`"+`)`),
			metaGoExample(b, engine, "inline-markers-single", "Inline Markers (single quotes)", "Single quotes", singleQuoteCode,
				`csharp title="UserQuery.cs" showLineNumbers 'context' ins='OrderBy' del='Select'`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, engine, "combined-markers", "Combined (markers + inline + focus)", "Combined markers", combinedCode,
				`go title="combined.go" showLineNumbers {4-5} ins={10-12} del={6-8} "db" focus={4-5,10-12}`,
				`html, err := engine.RenderWithMeta(code, meta)`),
		},
	}

	collapseEngine := b.engine(
		kazari.WithCollapsible(kazari.CollapsibleConfig{
			LineThreshold:    12,
			PreviewLines:     6,
			DefaultCollapsed: true,
			PreserveIndent:   true,
		}),
	)
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
	json.NewEncoder(w).Encode(Response{Status: 200, Message: "OK", Data: data})
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)
	result := processData(body)
	json.NewEncoder(w).Encode(Response{Status: 201, Message: "Created", Data: result})
}`
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
	data := map[string]interface{}{"status": "success", "count": 42}
	json.NewEncoder(w).Encode(data)
}`
	gapHTML := b.render(collapseEngine, gapCode, kazari.Options{
		Lang:        "go",
		Title:       "server.go",
		LineNumbers: &lineNumbers,
		LineMarkers: []kazari.LineMarker{{Type: kazari.MarkerIns, Lines: []kazari.Range{{Start: 10, End: 11}}}},
	})

	collapsible := Category{
		ID:          "collapsible",
		Title:       "Collapsible Sections",
		Description: "Threshold and range-based strategies for keeping long examples compact.",
		Examples: []Example{
			{
				ID:          "collapse-threshold",
				Title:       "Threshold-based (auto-collapses long blocks)",
				NavTitle:    "Threshold",
				Description: "Threshold behavior is configured at engine level.",
				HTML:        b.render(collapseEngine, thresholdCode, kazari.Options{Lang: "go", Title: "server.go"}),
				Recipes: []Recipe{recipe("Go", `engine := kazari.New(kazari.WithCollapsible(kazari.CollapsibleConfig{
	LineThreshold: 12, PreviewLines: 6, DefaultCollapsed: true,
}))
html, err := engine.Render(code, kazari.Options{Lang: "go", Title: "server.go"})`)},
			},
			metaGoExample(b, collapseEngine, "collapse-range", "Range-based (imports collapsed)", "Range", rangeCode,
				`go title="calc.go" showLineNumbers collapse={3-8}`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`go title="calc.go" showLineNumbers collapse={3-8}`+"`"+`)`),
			metaGoExample(b, collapseEngine, "collapse-multiple", "Multiple ranges", "Multiple ranges", multiRangeCode,
				`go title="api.go" showLineNumbers collapse={3-7,9-13}`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			{
				ID:          "collapse-gaps",
				Title:       "Threshold + markers (gap indicators)",
				NavTitle:    "Gap indicators",
				Description: "Structured options combine threshold collapsing with highlighted lines.",
				HTML:        gapHTML,
				Recipes: []Recipe{recipe("Go", `enabled := true
html, err := engine.Render(code, kazari.Options{
	Lang: "go", Title: "server.go", LineNumbers: &enabled,
	LineMarkers: []kazari.LineMarker{{
		Type: kazari.MarkerIns,
		Lines: []kazari.Range{{Start: 10, End: 11}},
	}},
})`)},
			},
			metaGoExample(b, collapseEngine, "collapse-start", "collapsible-start (re-collapsible, summary above)", "Collapsible start", rangeCode,
				`go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-start"`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, collapseEngine, "collapse-end", "collapsible-end (re-collapsible, summary below)", "Collapsible end", rangeCode,
				`go title="calc.go" showLineNumbers collapse={3-8} collapseStyle="collapsible-end"`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, collapseEngine, "collapse-auto", "collapsible-auto (auto start/end based on position)", "Collapsible auto", rangeCode,
				`go title="calc.go" showLineNumbers collapse={3-8,20-24} collapseStyle="collapsible-auto"`,
				`html, err := engine.RenderWithMeta(code, meta)`),
		},
	}

	mermaidCode := `graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Do something]
    B -->|No| D[Do something else]
    C --> E[End]
    D --> E`
	regexCode := `func fetchUsers() ([]User, error) {
	resp, err := http.Get("/api/users")
	if err != nil {
		return nil, fmt.Errorf("fetchUsers failed: %w", err)
	}
	defer resp.Body.Close()
	var users []User
	json.NewDecoder(resp.Body).Decode(&users)
	return users, nil
}`
	captureCode := "haystack = \"yes\"\nconfirm = \"yep\"\nreject = \"nope\""
	diffCode := " import (\n-\t\"fmt\"\n+\t\"log\"\n \t\"os\"\n )\n \n func main() {\n-\tfmt.Println(\"hello\")\n+\tlog.Println(\"hello\")\n }"
	codeGroupMarkdown := ":::code-group\n\n```go title=\"main.go\"\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello from Go!\")\n}\n```\n\n```python title=\"main.py\"\ndef main():\n    print(\"Hello from Python!\")\n\nif __name__ == \"__main__\":\n    main()\n```\n\n```javascript title=\"index.js\"\nfunction main() {\n  console.log(\"Hello from JavaScript!\");\n}\n\nmain();\n```\n\n:::\n"
	ansiCode := "\x1b[1;34mINFO\x1b[0m  Server started on \x1b[32m:8080\x1b[0m\n" +
		"\x1b[1;33mWARN\x1b[0m  Cache miss for key \x1b[36m\"user:42\"\x1b[0m\n" +
		"\x1b[1;31mERROR\x1b[0m Connection refused: \x1b[4mdb.example.com:5432\x1b[0m\n" +
		"\x1b[90m2024-01-15 10:30:45\x1b[0m \x1b[38;5;208mDEBUG\x1b[0m Retrying in \x1b[1m3s\x1b[0m..."
	syncMarkdown := ":::code-group sync=\"language\"\n\n```go\ngo get github.com/example/pkg\n```\n\n```python\npip install example-pkg\n```\n\n```javascript\nnpm install example-pkg\n```\n\n:::\n\n<p>Select a language above and the group below syncs automatically.</p>\n\n:::code-group sync=\"language\"\n\n```go\nimport \"github.com/example/pkg\"\n```\n\n```python\nimport example_pkg\n```\n\n```javascript\nconst pkg = require('example-pkg');\n```\n\n:::\n\n<p>This group uses a different sync key (<code>sync=\"platform\"</code>) and syncs independently.</p>\n\n:::code-group sync=\"platform\"\n\n```bash title=\"Linux\"\nsudo apt install build-essential\n```\n\n```powershell title=\"Windows\"\nwinget install Microsoft.VisualStudio.BuildTools\n```\n\n```bash title=\"macOS\"\nbrew install gcc\n```\n\n:::\n"

	markdownRenderer := goldmark.New(goldmark.WithExtensions(
		kazarimd.New(collapseEngine),
		kazarimd.CodeGroups(collapseEngine),
	))
	codeGroupHTML := b.renderMarkdown(markdownRenderer, codeGroupMarkdown)
	syncHTML := b.renderMarkdown(markdownRenderer, syncMarkdown)

	formats := Category{
		ID:          "formats",
		Title:       "Formats and Groups",
		Description: "Specialized formats, regex matching, diffs, ANSI output, and tabbed groups.",
		Examples: []Example{
			metaGoExample(b, engine, "mermaid", "Mermaid Pass-Through (raw code for Mermaid.js)", "Mermaid", mermaidCode,
				`mermaid`, `html, err := engine.Render(code, kazari.Options{Lang: "mermaid"})`),
			metaGoExample(b, engine, "regex-markers", "Regex Markers", "Regex markers", regexCode,
				`go title="regex-markers.go" showLineNumbers /err\b/ ins=/func\s+\w+/ del=/fmt\.Errorf/`,
				`html, err := engine.RenderWithMeta(code, meta)`),
			metaGoExample(b, engine, "regex-capture", `Regex Capture Group (/ye(s|p)/ marks only "s" or "p")`, "Regex capture", captureCode,
				`python title="capture_group.py" /ye(s|p)/`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`python title="capture_group.py" /ye(s|p)/`+"`"+`)`),
			metaGoExample(b, engine, "hybrid-diff", `Hybrid Diff + Syntax Highlighting (diff lang="go")`, "Hybrid diff", diffCode,
				`diff lang="go" title="hybrid-diff.go" showLineNumbers`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`diff lang="go" title="hybrid-diff.go" showLineNumbers`+"`"+`)`),
			{
				ID:       "code-group",
				Title:    "Code Group (tabbed code blocks via Goldmark)",
				NavTitle: "Code group",
				HTML:     codeGroupHTML,
				Recipes:  []Recipe{recipe("Markdown", codeGroupMarkdown)},
			},
			metaGoExample(b, engine, "ansi", "ANSI Escape Sequences (parsed SGR codes)", "ANSI", ansiCode,
				`ansi title="server.log" showLineNumbers`,
				`html, err := engine.RenderWithMeta(code, `+"`"+`ansi title="server.log" showLineNumbers`+"`"+`)`),
			{
				ID:       "code-group-sync",
				Title:    "Code Group Tab Sync (tabs synced across groups)",
				NavTitle: "Tab sync",
				HTML:     syncHTML,
				Recipes:  []Recipe{recipe("Markdown", syncMarkdown)},
			},
		},
	}

	themeCode := "func main() {\n\tfmt.Println(\"Same code, different theme\")\n}"
	themeDefaultHTML := b.renderMeta(engine, themeCode, `go title="default theme" showLineNumbers`)
	themeDraculaHTML := b.renderMeta(engine, themeCode, `go title="theme=dracula" showLineNumbers theme="dracula"`)
	themeDualHTML := b.renderMeta(engine, themeCode, `go title="theme=dracula,github-light" showLineNumbers theme="dracula,github-light"`)
	customizerEngine := b.engine(
		kazari.WithThemeCSSRoot(".kazari-customizer"),
		kazari.WithThemeCustomizer(func(name string, colors kazari.ThemeInfo) kazari.ThemeInfo {
			if name == "github-dark" {
				colors.BG = "#1a1b26"
			}
			return colors
		}),
	)
	customizerHTML := b.render(customizerEngine, `fmt.Println("Custom dark BG: #1a1b26")`, kazari.Options{Lang: "go", Title: "customized-theme.go"})
	tintHue, tintChroma := 195.0, 0.04
	tintedEngine := b.engine(
		kazari.WithThemeCSSRoot(".kazari-tinted"),
		kazari.WithThemeAdjustments(kazari.ThemeAdjustments{Hue: &tintHue, Chroma: &tintChroma}),
	)
	tintedHTML := b.render(tintedEngine, `fmt.Println("Backgrounds tinted toward teal in OKLCH space")`, kazari.Options{Lang: "go", Title: "tinted-theme.go"})
	scopedEngine := b.engine(kazari.WithThemeCSSRoot(".kazari-scoped"))
	scopedHTML := b.render(scopedEngine, `fmt.Println("CSS vars scoped to .kazari-scoped")`, kazari.Options{Lang: "go", Title: "scoped.go"})

	themes := Category{
		ID:          "themes",
		Title:       "Themes and Customization",
		Description: "Per-block themes, generated adjustments, and scoped CSS output.",
		Examples: []Example{
			{
				ID:       "theme-override",
				Title:    "Per-Block Theme Override (default vs dracula)",
				NavTitle: "Per-block override",
				HTML:     joinHTML(themeDefaultHTML, themeDraculaHTML),
				Recipes: []Recipe{
					recipe("Meta", "go title=\"default theme\" showLineNumbers\ngo title=\"theme=dracula\" showLineNumbers theme=\"dracula\""),
					recipe("Go", `defaultHTML, err := engine.RenderWithMeta(code, `+"`"+`go title="default theme" showLineNumbers`+"`"+`)
draculaHTML, err := engine.RenderWithMeta(code, `+"`"+`go title="theme=dracula" showLineNumbers theme="dracula"`+"`"+`)`),
				},
			},
			{
				ID:       "theme-override-dual",
				Title:    "Per-Block Theme Override (dual: dracula + github-light)",
				NavTitle: "Dual override",
				HTML:     themeDualHTML,
				Recipes: []Recipe{
					recipe("Meta", `go title="theme=dracula,github-light" showLineNumbers theme="dracula,github-light"`),
					recipe("Go", `html, err := engine.RenderWithMeta(code, `+"`"+`go showLineNumbers theme="dracula,github-light"`+"`"+`)`),
				},
			},
			{
				ID:           "theme-customizer",
				Title:        "Theme Customizer (dark BG changed to #1a1b26)",
				NavTitle:     "Theme customizer",
				HTML:         customizerHTML,
				WrapperClass: "kazari-customizer",
				Recipes: []Recipe{recipe("Go", `engine := kazari.New(
	kazari.WithThemeCSSRoot(".kazari-customizer"),
	kazari.WithThemeCustomizer(func(name string, colors kazari.ThemeInfo) kazari.ThemeInfo {
		if name == "github-dark" { colors.BG = "#1a1b26" }
		return colors
	}),
)`)},
			},
			{
				ID:           "theme-adjustments",
				Title:        "Theme Adjustments (OKLCH teal tint)",
				NavTitle:     "Theme adjustments",
				Description:  "Toggle dark mode to see the background tint.",
				HTML:         tintedHTML,
				WrapperClass: "kazari-tinted",
				Recipes: []Recipe{recipe("Go", `hue, chroma := 195.0, 0.04
engine := kazari.New(
	kazari.WithThemeCSSRoot(".kazari-tinted"),
	kazari.WithThemeAdjustments(kazari.ThemeAdjustments{
		Hue: &hue, Chroma: &chroma,
	}),
)`)},
			},
			{
				ID:           "scoped-css",
				Title:        "Scoped CSS Root (.kazari-scoped)",
				NavTitle:     "Scoped CSS",
				HTML:         scopedHTML,
				WrapperClass: "kazari-scoped",
				Recipes: []Recipe{recipe("Go", `engine := kazari.New(kazari.WithThemeCSSRoot(".kazari-scoped"))
css := engine.CSS()`)},
			},
		},
	}

	frenchEngine := b.engine(kazari.WithLocale("fr-FR"))
	frenchHTML := b.render(frenchEngine, `fmt.Println("Bonjour le monde !")`, kazari.Options{Lang: "go", Title: "locale-fr.go"})
	iconEngine := b.engine(kazari.WithFileIconResolver(func(extension string) string {
		icons := map[string]string{"go": "🔵", "py": "🐍", "js": "🟡", "rs": "🦀", "css": "🎨"}
		icon := "📄"
		if resolved, ok := icons[extension]; ok {
			icon = resolved
		}
		return fmt.Sprintf(`<span class="kz-file-icon">%s</span>`, icon)
	}))
	iconHTML := joinHTML(
		b.render(iconEngine, `fmt.Println("Go")`, kazari.Options{Lang: "go", Title: "main.go"}),
		b.render(iconEngine, `print("Python")`, kazari.Options{Lang: "python", Title: "app.py"}),
		b.render(iconEngine, `console.log("JS")`, kazari.Options{Lang: "javascript", Title: "index.js"}),
		b.render(iconEngine, `fn main() {}`, kazari.Options{Lang: "rust", Title: "main.rs"}),
	)

	localization := Category{
		ID:          "localization",
		Title:       "Localization and Assets",
		Description: "Localized controls and consumer-provided file icon resolution.",
		Examples: []Example{
			{
				ID:          "locale-french",
				Title:       `Locale: French (WithLocale("fr-FR"))`,
				NavTitle:    "French locale",
				Description: `Copy and fullscreen controls use French labels.`,
				HTML:        frenchHTML,
				Recipes: []Recipe{recipe("Go", `engine := kazari.New(kazari.WithLocale("fr-FR"))
html, err := engine.Render(code, kazari.Options{Lang: "go", Title: "locale-fr.go"})`)},
			},
			{
				ID:          "file-icons",
				Title:       "File Icons (custom resolver)",
				NavTitle:    "File icons",
				Description: "WithFileIconResolver injects an icon based on the title extension.",
				HTML:        iconHTML,
				Recipes: []Recipe{recipe("Go", `engine := kazari.New(kazari.WithFileIconResolver(func(ext string) string {
	icons := map[string]string{"go": "🔵", "py": "🐍", "js": "🟡", "rs": "🦀"}
	return fmt.Sprintf("<span class=\"kz-file-icon\">%s</span>", icons[ext])
}))`)},
			},
		},
	}

	if b.err != nil {
		return nil, "", "", fmt.Errorf("render showcase example: %w", b.err)
	}

	css := strings.Join([]string{
		collapseEngine.CSS(),
		dotsEngine.ThemeCSS(),
		customizerEngine.ThemeCSS(),
		tintedEngine.ThemeCSS(),
		scopedEngine.ThemeCSS(),
		`.kazari-code .kz-file-icon { font-size: 1rem; margin-right: .4rem; }`,
	}, "\n")
	return []Category{frames, layout, markers, collapsible, formats, themes, localization}, css, collapseEngine.JS(), nil
}

func metaGoExample(b *catalogBuilder, engine *kazari.Engine, id, title, navTitle, code, meta, goRecipe string) Example {
	return Example{
		ID:       id,
		Title:    title,
		NavTitle: navTitle,
		HTML:     b.renderMeta(engine, code, meta),
		Recipes: []Recipe{
			recipe("Meta", meta),
			recipe("Go", goRecipe),
		},
	}
}
