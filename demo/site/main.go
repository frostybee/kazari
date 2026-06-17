package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/frostybee/kazari"
	kazarichroma "github.com/frostybee/kazari/chroma"
	"github.com/frostybee/kazari/demo/internal/showcase"
	"github.com/frostybee/kazari/demo/internal/snippets"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func sourceDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

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
	{Label: "Color Contrast", Href: "color-contrast.html"},
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

	rawHL, err := nuri.New(ctx, nuri.WithFS(core.FS()), nuri.WithMinContrast(0))
	if err != nil {
		log.Fatalf("nuri.New (raw): %v", err)
	}
	defer rawHL.Close(ctx)
	rawNuriHL := kazarinuri.New(ctx, rawHL)

	outDir := sourceDir()

	generateShowcase(outDir, nuriHL, chromaHL, &nuriSVGProvider{hl: hl, ctx: ctx})
	generateNuriVsShiki(outDir, nuriHL)
	generateNuriVsChroma(outDir, nuriHL, chromaHL)
	generateColorContrast(outDir, rawNuriHL, nuriHL)

	log.Printf("All pages written to %s/", outDir)
}

type nuriSVGProvider struct {
	hl  *nuri.Highlighter
	ctx context.Context
}

func (p *nuriSVGProvider) RenderSVG(code string, opts showcase.SVGOptions) (string, error) {
	nuriOpts := nuri.CodeToSVGOptions{
		Lang:         opts.Lang,
		Theme:        opts.Theme,
		FontSize:     opts.FontSize,
		CornerRadius: opts.CornerRadius,
		ShowBackground: opts.ShowBG,
	}
	return p.hl.CodeToSVG(p.ctx, code, nuriOpts)
}

func generateShowcase(outDir string, nuriHL, chromaHL kazari.Highlighter, svgProvider showcase.SVGProvider) {
	configOpts := []kazari.Option{kazari.WithConfigDir(sourceDir())}

	nuriConfig := showcase.Config{
		BackendName:   "Nuri",
		HTMLFile:      "showcase.html",
		OtherName:     "Chroma",
		OtherHref:     "showcase-chroma.html",
		NavLinks:      navLinksWithActive("Showcase"),
		KazariOptions: configOpts,
		SVGProvider:   svgProvider,
	}
	if err := showcase.Generate(outDir, nuriConfig, nuriHL); err != nil {
		log.Fatalf("showcase (nuri): %v", err)
	}
	log.Println("Written: showcase.html, showcase.css, showcase.js")

	chromaConfig := showcase.Config{
		BackendName:   "Chroma",
		HTMLFile:      "showcase-chroma.html",
		OtherName:     "Nuri",
		OtherHref:     "showcase.html",
		NavLinks:      navLinksWithActive("Showcase"),
		KazariOptions: configOpts,
	}
	if err := showcase.Generate(outDir, chromaConfig, chromaHL); err != nil {
		log.Fatalf("showcase (chroma): %v", err)
	}
	log.Println("Written: showcase-chroma.html")
}

func generateNuriVsShiki(outDir string, nuriHL kazari.Highlighter) {
	engine := kazari.New(
		kazari.WithConfigDir(sourceDir()),
		kazari.WithHighlighter(nuriHL),
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
		kazari.WithConfigDir(sourceDir()),
		kazari.WithHighlighter(nuriHL),
	)
	chromaEngine := kazari.New(
		kazari.WithConfigDir(sourceDir()),
		kazari.WithHighlighter(chromaHL),
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

type contrastPair struct {
	light string
	dark  string
	lang  string
	label string
	code  string
}

func generateColorContrast(outDir string, rawHL, correctedHL kazari.Highlighter) {
	pairs := []contrastPair{
		{"solarized-light", "solarized-dark", "go", "Solarized / Go", snippets.GoCode},
		{"rose-pine-dawn", "rose-pine", "javascript", "Rose Pine / JavaScript", snippets.JSCode},
		{"vitesse-light", "vitesse-dark", "python", "Vitesse / Python", snippets.PyCode},
		{"one-light", "one-dark-pro", "typescript", "One Dark Pro / TypeScript", snippets.TSCode},
		{"rose-pine-dawn", "rose-pine-moon", "bash", "Rose Pine Moon / Bash", snippets.BashCode},
		{"solarized-light", "solarized-dark", "css", "Solarized / CSS", snippets.CSSCode},
	}

	noFrame := kazari.FrameNone
	var rows strings.Builder
	var baseCSS, baseJS string

	for _, p := range pairs {
		rawEngine := kazari.New(
			kazari.WithHighlighter(rawHL),
			kazari.WithThemes(p.light, p.dark),
			kazari.WithMinContrast(0),
			kazari.WithMinify(true),
		)
		correctedEngine := kazari.New(
			kazari.WithHighlighter(correctedHL),
			kazari.WithThemes(p.light, p.dark),
			kazari.WithMinify(true),
		)

		if baseCSS == "" {
			baseCSS = rawEngine.CSS()
			baseJS = rawEngine.JS()
		}

		lightInfo, err := rawHL.GetThemeColors(p.light)
		if err != nil {
			log.Fatalf("theme colors %s: %v", p.light, err)
		}
		darkInfo, err := rawHL.GetThemeColors(p.dark)
		if err != nil {
			log.Fatalf("theme colors %s: %v", p.dark, err)
		}

		rawHTML, err := rawEngine.Render(p.code, kazari.Options{
			Lang: p.lang, Frame: &noFrame,
		})
		if err != nil {
			log.Fatalf("contrast raw render %s: %v", p.light, err)
		}
		correctedHTML, err := correctedEngine.Render(p.code, kazari.Options{
			Lang: p.lang, Frame: &noFrame,
		})
		if err != nil {
			log.Fatalf("contrast corrected render %s: %v", p.light, err)
		}

		rows.WriteString(fmt.Sprintf(contrastSectionTmpl,
			p.light, p.label,
			lightInfo.BG, lightInfo.FG, darkInfo.BG, darkInfo.FG,
			rawHTML, correctedHTML,
		))
	}

	page := fmt.Sprintf(contrastPageTmpl,
		baseCSS,
		mustRead("comparison.css"),
		mustRead("nav.css"),
		navHTML("Color Contrast"),
		rows.String(),
		baseJS,
		mustRead("dark-toggle.js"),
	)

	if err := os.WriteFile(outDir+"/color-contrast.html", []byte(page), 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Println("Written: color-contrast.html")
}

const contrastSectionTmpl = `<section class="cmp-row" data-lang="%s" style="--cmp-light-bg:%[3]s;--cmp-light-fg:%[4]s;--cmp-dark-bg:%[5]s;--cmp-dark-fg:%[6]s">
  <h2 class="cmp-lang-heading">%[2]s</h2>
  <div class="cmp-grid">
    <div class="cmp-col">
      <div class="cmp-col-label">Before (raw theme colors)</div>
      %[7]s
    </div>
    <div class="cmp-col">
      <div class="cmp-col-label">After (contrast corrected)</div>
      %[8]s
    </div>
  </div>
</section>`

const contrastPageTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:ital,wght@0,100..800;1,100..800&display=swap" rel="stylesheet">
<title>Kazari: Color Contrast Correction</title>
<style>%s</style>
<style>%s</style>
<style>%s</style>
<style>.cmp-row { --kz-editor-bg: var(--cmp-light-bg); --kz-editor-fg: var(--cmp-light-fg); }
.dark .cmp-row { --kz-editor-bg: var(--cmp-dark-bg); --kz-editor-fg: var(--cmp-dark-fg); }</style>
</head>
<body>
%s
<header class="cmp-header">
  <h1 class="cmp-page-title">Color Contrast Correction</h1>
  <p>Many syntax themes prioritize aesthetics over readability, producing token colors that lack sufficient contrast against the editor background. This is an accessibility concern: the <a href="https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html" target="_blank">WCAG 2.1 Success Criterion 1.4.3 (Contrast Minimum)</a> requires a contrast ratio of at least 4.5:1 for normal text, and 3:1 for large text (Level AA). Code displayed in small monospace fonts falls under the stricter 4.5:1 threshold.</p>
  <p>Kazari can automatically enforce this by adjusting syntax token colors to meet a configurable minimum contrast ratio (default 5.5:1, slightly above WCAG AA) against the theme&#39;s editor background. Correction is applied <strong>per token</strong>: only colors that fall below the threshold are shifted toward a more readable value while staying as close to the original hue as possible. Tokens that already meet the ratio are left untouched.</p>
  <p>As a result, well-designed themes (such as Catppuccin Mocha or Tokyo Night) may show little or no visible difference, while themes known for low-contrast palettes (such as Solarized, Vitesse, or One Dark Pro) will show noticeable adjustments.</p>
  <p>Each row below uses a different theme pair. <strong>Left:</strong> raw theme colors (contrast correction disabled). <strong>Right:</strong> after correction. Toggle dark mode to compare both variants.</p>
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
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:ital,wght@0,100..800;1,100..800&display=swap" rel="stylesheet">
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
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:ital,wght@0,100..800;1,100..800&display=swap" rel="stylesheet">
<title>Kazari: Nuri vs Chroma</title>
<style>%s</style>
<style>%s</style>
<style>%s</style>
</head>
<body>
%s
<header class="cmp-header">
  <p>Both sides are pre-rendered at build time. Nuri uses TextMate grammars (same as Shiki/VS Code). Chroma uses Pygments-based lexers.</p>
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
