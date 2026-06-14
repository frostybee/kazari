package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/frostybee/kazari"
	kazarichroma "github.com/frostybee/kazari/chroma"
	"github.com/frostybee/kazari/demo/internal/showcase"
	"github.com/frostybee/kazari/demo/internal/snippets"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

//go:embed assets/nav.css assets/comparison.css assets/shiki.js assets/dark-toggle.js
var assets embed.FS

func mustRead(name string) string {
	data, err := assets.ReadFile("assets/" + name)
	if err != nil {
		log.Fatalf("read %s: %v", name, err)
	}
	return string(data)
}

var navLinks = []showcase.NavLink{
	{Label: "Showcase", Href: "showcase.html"},
	{Label: "Nuri vs Shiki", Href: "nuri-vs-shiki.html"},
	{Label: "Nuri vs Chroma", Href: "nuri-vs-chroma.html"},
}

func navLinksWithActive(active string) []showcase.NavLink {
	links := make([]showcase.NavLink, len(navLinks))
	copy(links, navLinks)
	for i := range links {
		if links[i].Label == active {
			links[i].Active = true
		}
	}
	return links
}

func main() {
	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		log.Fatalf("nuri.New: %v", err)
	}
	defer hl.Close(ctx)

	nuriHL := kazarinuri.New(ctx, hl)
	chromaHL := kazarichroma.New(kazarichroma.WithStyleMap(map[string]string{
		"github-light": "github",
		"github-dark":  "github-dark",
	}))

	outDir := "demo/site"

	generateShowcase(outDir, nuriHL, chromaHL)
	generateNuriVsShiki(outDir, nuriHL)
	generateNuriVsChroma(outDir, nuriHL, chromaHL)

	log.Printf("All pages written to %s/", outDir)
}

func generateShowcase(outDir string, nuriHL, chromaHL kazari.Highlighter) {
	nuriConfig := showcase.Config{
		BackendName: "Nuri",
		HTMLFile:    "showcase.html",
		OtherName:   "Chroma",
		OtherHref:   "showcase-chroma.html",
		NavLinks:    navLinksWithActive("Showcase"),
	}
	if err := showcase.Generate(outDir, nuriConfig, nuriHL); err != nil {
		log.Fatalf("showcase (nuri): %v", err)
	}
	log.Println("Written: showcase.html, showcase.css, showcase.js")

	chromaConfig := showcase.Config{
		BackendName: "Chroma",
		HTMLFile:    "showcase-chroma.html",
		OtherName:   "Nuri",
		OtherHref:   "showcase.html",
		NavLinks:    navLinksWithActive("Showcase"),
	}
	if err := showcase.Generate(outDir, chromaConfig, chromaHL); err != nil {
		log.Fatalf("showcase (chroma): %v", err)
	}
	log.Println("Written: showcase-chroma.html")
}

func generateNuriVsShiki(outDir string, nuriHL kazari.Highlighter) {
	engine := kazari.New(
		kazari.WithHighlighter(nuriHL),
		kazari.WithThemes("github-light", "github-dark"),
	)

	noFrame := kazari.FrameNone
	var rows strings.Builder
	for _, s := range snippets.All {
		rendered, err := engine.Render(s.Code, kazari.Options{
			Lang:  s.Lang,
			Frame: &noFrame,
		})
		if err != nil {
			log.Fatalf("nuri render %s: %v", s.ID, err)
		}
		rows.WriteString(fmt.Sprintf(shikiSectionTmpl,
			s.ID, s.Label, rendered,
			s.ID, s.Lang,
			s.ID, escapeForScript(s.Code),
		))
	}

	page := fmt.Sprintf(shikiPageTmpl,
		engine.CSS(),
		mustRead("comparison.css"),
		mustRead("nav.css"),
		navHTML("Nuri vs Shiki"),
		rows.String(),
		engine.JS(),
		mustRead("shiki.js"),
	)

	if err := os.WriteFile(outDir+"/nuri-vs-shiki.html", []byte(page), 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Println("Written: nuri-vs-shiki.html")
}

func generateNuriVsChroma(outDir string, nuriHL, chromaHL kazari.Highlighter) {
	nuriEngine := kazari.New(
		kazari.WithHighlighter(nuriHL),
		kazari.WithThemes("github-light", "github-dark"),
	)
	chromaEngine := kazari.New(
		kazari.WithHighlighter(chromaHL),
		kazari.WithThemes("github-light", "github-dark"),
	)

	noFrame := kazari.FrameNone
	var rows strings.Builder
	for _, s := range snippets.All {
		nuriHTML, err := nuriEngine.Render(s.Code, kazari.Options{
			Lang: s.Lang, Frame: &noFrame,
		})
		if err != nil {
			log.Fatalf("nuri render %s: %v", s.ID, err)
		}
		chromaHTML, err := chromaEngine.Render(s.Code, kazari.Options{
			Lang: s.Lang, Frame: &noFrame,
		})
		if err != nil {
			log.Fatalf("chroma render %s: %v", s.ID, err)
		}
		rows.WriteString(fmt.Sprintf(chromaSectionTmpl,
			s.ID, s.Label, nuriHTML, chromaHTML,
		))
	}

	page := fmt.Sprintf(chromaPageTmpl,
		nuriEngine.CSS(),
		mustRead("comparison.css"),
		mustRead("nav.css"),
		navHTML("Nuri vs Chroma"),
		rows.String(),
		nuriEngine.JS(),
		mustRead("dark-toggle.js"),
	)

	if err := os.WriteFile(outDir+"/nuri-vs-chroma.html", []byte(page), 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Println("Written: nuri-vs-chroma.html")
}

func navHTML(active string) string {
	var b strings.Builder
	b.WriteString(`<nav class="site-nav"><div class="site-nav-inner">`)
	b.WriteString(`<a class="site-brand" href="showcase.html">Kazari</a>`)
	b.WriteString(`<div class="site-nav-links">`)
	for _, l := range navLinksWithActive(active) {
		if l.Active {
			b.WriteString(fmt.Sprintf(`<a href="%s" class="active" aria-current="page">%s</a>`, l.Href, l.Label))
		} else {
			b.WriteString(fmt.Sprintf(`<a href="%s">%s</a>`, l.Href, l.Label))
		}
	}
	b.WriteString(`</div>`)
	b.WriteString(`<label class="site-theme-toggle" title="Toggle dark mode"><input type="checkbox" id="dark-toggle"><span class="toggle-track"><span class="toggle-thumb"></span></span></label>`)
	b.WriteString(`</div></nav>`)
	return b.String()
}

func escapeForScript(s string) string {
	return strings.ReplaceAll(s, "</script>", "<\\/script>")
}

const shikiSectionTmpl = `<section class="cmp-row" data-lang="%s">
  <h2 class="cmp-lang-heading">%s</h2>
  <div class="cmp-grid">
    <div class="cmp-col cmp-col--nuri">
      <div class="cmp-col-label">Nuri (Go)</div>
      %s
    </div>
    <div class="cmp-col cmp-col--shiki">
      <div class="cmp-col-label">Shiki (JS / CDN)</div>
      <div class="shiki-target" id="shiki-%s" data-lang="%s"></div>
    </div>
  </div>
  <script type="text/plain" id="src-%s">%s</script>
</section>`

const chromaSectionTmpl = `<section class="cmp-row" data-lang="%s">
  <h2 class="cmp-lang-heading">%s</h2>
  <div class="cmp-grid">
    <div class="cmp-col">
      <div class="cmp-col-label">Nuri</div>
      %s
    </div>
    <div class="cmp-col">
      <div class="cmp-col-label">Chroma</div>
      %s
    </div>
  </div>
</section>`

const shikiPageTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
<title>Kazari: Nuri vs Shiki</title>
<style>%s</style>
<style>%s</style>
<style>%s</style>
</head>
<body>
%s
<header class="cmp-header">
  <span id="shiki-status" class="cmp-status">Loading Shiki from CDN...</span>
  <p>Left: Nuri (Go, pre-rendered at build time). Right: Shiki (JS, loaded from CDN).</p>
</header>
<main>%s</main>
<footer class="site-footer">
  <div class="site-footer-inner">
    <div class="site-footer-about">
      <p class="site-footer-brand">Kazari <span class="site-footer-kanji">飾り</span></p>
      <p>A Go library for rendering framed, syntax-highlighted code blocks with full CSS customization. Powered by <a href="https://github.com/frostybee/nuri">Nuri</a>, a pure Go port of Shiki.</p>
    </div>
    <div class="site-footer-links">
      <a href="https://github.com/frostybee/kazari">GitHub</a>
      <span class="site-footer-sep" aria-hidden="true"></span>
      <a href="https://github.com/frostybee/nuri">Nuri</a>
      <span class="site-footer-sep" aria-hidden="true"></span>
      <a href="https://github.com/frostybee/kazari/blob/main/LICENSE">MIT License</a>
    </div>
  </div>
</footer>
<script>%s</script>
<script type="module">%s</script>
</body>
</html>`

const chromaPageTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
<title>Kazari: Nuri vs Chroma</title>
<style>%s</style>
<style>%s</style>
<style>%s</style>
</head>
<body>
%s
<header class="cmp-header">
  <p>Both sides are pre-rendered at build time. Nuri uses TextMate grammars (same as Shiki/VS Code). Chroma uses its own lexer-based tokenization.</p>
</header>
<main>%s</main>
<footer class="site-footer">
  <div class="site-footer-inner">
    <div class="site-footer-about">
      <p class="site-footer-brand">Kazari <span class="site-footer-kanji">飾り</span></p>
      <p>A Go library for rendering framed, syntax-highlighted code blocks with full CSS customization. Powered by <a href="https://github.com/frostybee/nuri">Nuri</a>, a pure Go port of Shiki.</p>
    </div>
    <div class="site-footer-links">
      <a href="https://github.com/frostybee/kazari">GitHub</a>
      <span class="site-footer-sep" aria-hidden="true"></span>
      <a href="https://github.com/frostybee/nuri">Nuri</a>
      <span class="site-footer-sep" aria-hidden="true"></span>
      <a href="https://github.com/frostybee/kazari/blob/main/LICENSE">MIT License</a>
    </div>
  </div>
</footer>
<script>%s</script>
<script>%s</script>
</body>
</html>`
