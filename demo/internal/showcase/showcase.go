package showcase

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/frostybee/kazari"
)

//go:embed templates/page.html assets/showcase.css assets/showcase.js
var sourceFiles embed.FS

type SVGProvider interface {
	RenderSVG(code string, opts SVGOptions) (string, error)
}

type SVGOptions struct {
	Lang         string
	Theme        string
	FontSize     float64
	CornerRadius float64
	ShowBG       *bool
}

type Config struct {
	BackendName    string
	HTMLFile       string
	OtherName      string
	OtherHref      string
	NavLinks       []NavLink
	KazariOptions  []kazari.Option
	SVGProvider    SVGProvider
}

type NavLink struct {
	Label  string
	Href   string
	Active bool
}

type BackendLink struct {
	Name   string
	Href   string
	Active bool
}

type Recipe struct {
	Label string
	Code  string
}

type Example struct {
	ID           string
	Title        string
	NavTitle     string
	Description  template.HTML
	HTML         template.HTML
	WrapperClass string
	Recipes      []Recipe
	SearchText   string
}

func (e Example) HasRecipeTabs() bool {
	return len(e.Recipes) > 1
}

type Category struct {
	ID          string
	Title       string
	Description string
	Examples    []Example
}

type Page struct {
	Title        string
	Subtitle     string
	Backends     []BackendLink
	NavLinks     []NavLink
	Categories   []Category
	ExampleCount int
}

type Output struct {
	Page Page
	HTML string
	CSS  string
	JS   string
}

func Generate(dir string, cfg Config, highlighter kazari.Highlighter) error {
	output, err := Build(cfg, highlighter)
	if err != nil {
		return err
	}

	files := map[string]string{
		cfg.HTMLFile:   output.HTML,
		"showcase.css": output.CSS,
		"showcase.js":  output.JS,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

func Build(cfg Config, highlighter kazari.Highlighter) (Output, error) {
	if cfg.BackendName == "" || cfg.HTMLFile == "" || cfg.OtherName == "" || cfg.OtherHref == "" {
		return Output{}, fmt.Errorf("showcase: incomplete configuration")
	}
	if highlighter == nil {
		return Output{}, fmt.Errorf("showcase: nil highlighter")
	}

	catalog, generatedCSS, generatedJS, err := buildCatalog(highlighter, cfg.KazariOptions, cfg.SVGProvider)
	if err != nil {
		return Output{}, err
	}
	prepareCatalog(catalog)

	page := Page{
		Title:        "Kazari Showcase with " + cfg.BackendName,
		Subtitle:     "Code blocks highlighted by " + cfg.BackendName + ".",
		Categories:   catalog,
		ExampleCount: countExamples(catalog),
		NavLinks:     cfg.NavLinks,
		Backends: []BackendLink{
			{Name: cfg.BackendName, Active: true},
			{Name: cfg.OtherName, Href: cfg.OtherHref},
		},
	}

	pageTemplate, err := template.ParseFS(sourceFiles, "templates/page.html")
	if err != nil {
		return Output{}, fmt.Errorf("parse page template: %w", err)
	}
	var html bytes.Buffer
	if err := pageTemplate.Execute(&html, page); err != nil {
		return Output{}, fmt.Errorf("render page template: %w", err)
	}

	pageCSS, err := sourceFiles.ReadFile("assets/showcase.css")
	if err != nil {
		return Output{}, fmt.Errorf("read showcase CSS: %w", err)
	}
	pageJS, err := sourceFiles.ReadFile("assets/showcase.js")
	if err != nil {
		return Output{}, fmt.Errorf("read showcase JS: %w", err)
	}

	return Output{
		Page: page,
		HTML: trimTrailingWhitespace(html.String()),
		CSS:  generatedCSS + "\n" + string(pageCSS),
		JS:   generatedJS + "\n" + string(pageJS),
	}, nil
}

func trimTrailingWhitespace(content string) string {
	lines := strings.Split(content, "\n")
	for index := range lines {
		lines[index] = strings.TrimRight(lines[index], " \t")
	}
	return strings.Join(lines, "\n")
}

func prepareCatalog(categories []Category) {
	for categoryIndex := range categories {
		category := &categories[categoryIndex]
		for exampleIndex := range category.Examples {
			example := &category.Examples[exampleIndex]
			if description, ok := exampleDescriptions[example.ID]; ok {
				example.Description = template.HTML(description)
			}
			if example.NavTitle == "" {
				example.NavTitle = example.Title
			}
			var search strings.Builder
			search.WriteString(category.Title)
			search.WriteByte(' ')
			search.WriteString(example.Title)
			search.WriteByte(' ')
			search.WriteString(string(example.Description))
			for _, recipe := range example.Recipes {
				search.WriteByte(' ')
				search.WriteString(recipe.Label)
				search.WriteByte(' ')
				search.WriteString(recipe.Code)
			}
			example.SearchText = strings.ToLower(search.String())
		}
	}
}

func countExamples(categories []Category) int {
	total := 0
	for _, category := range categories {
		total += len(category.Examples)
	}
	return total
}
