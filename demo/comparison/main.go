package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/frostybee/kazari"
	"github.com/frostybee/kazari/demo/internal/snippets"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

const sectionTmpl = `<section class="cmp-row" data-lang="%s">
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

const navCSS = `
.site-nav { position: sticky; top: 0; z-index: 30; background: rgba(248,249,250,.96); border-bottom: 1px solid #dfe3e8; backdrop-filter: blur(12px); }
.site-nav-inner { display: flex; align-items: center; gap: 1.25rem; max-width: 1400px; margin: 0 auto; padding: .65rem clamp(1rem,3vw,2.5rem); }
.site-brand { font-size: 1.15rem; font-weight: 700; color: inherit; text-decoration: none; }
.site-nav-links { display: flex; gap: .25rem; }
.site-nav-links a { padding: .4rem .75rem; border-radius: .4rem; color: #5b6470; text-decoration: none; font-size: .88rem; font-weight: 500; transition: background .15s, color .15s; }
.site-nav-links a:hover { background: #f0f1ff; color: #303bc4; }
.site-nav-links a.active { background: #e8eafc; color: #303bc4; font-weight: 650; }
.site-dark-toggle { margin-left: auto; font-size: .88rem; cursor: pointer; user-select: none; }
.dark .site-nav { background: rgba(12,12,23,.96); border-color: #34364b; }
.dark .site-brand { color: #f0f1f5; }
.dark .site-nav-links a { color: #b8bdc9; }
.dark .site-nav-links a:hover { background: #303252; color: #d9dcff; }
.dark .site-nav-links a.active { background: #3c3f67; color: #f0f1ff; }
.dark .site-dark-toggle { color: #e0e0e0; }
`

const pageCSS = `*, *::before, *::after { box-sizing: border-box; }
body { margin: 0; font-family: system-ui, sans-serif; background: #f6f8fa; color: #24292e; }
.dark body { background: #0d1117; color: #c9d1d9; }

.cmp-header {
  max-width: 1400px; margin: 0 auto; padding: 1.5rem clamp(1rem,3vw,2.5rem) .5rem;
  display: flex; align-items: center; gap: 1.5rem; flex-wrap: wrap;
}
.cmp-header p { margin: 0; color: #57606a; font-size: .9rem; flex: 1 1 100%; }
.dark .cmp-header p { color: #8b949e; }
.cmp-status { font-size: .85rem; color: #57606a; }
.cmp-status.ready { color: #1a7f37; }
.cmp-status.error { color: #cf222e; }
.dark .cmp-status { color: #8b949e; }
.dark .cmp-status.ready { color: #3fb950; }

.cmp-row { max-width: 1400px; margin: 0 auto 3rem; padding: 0 clamp(1rem,3vw,2.5rem); position: relative; }

.cmp-lang-heading {
  margin: 0 0 .75rem; font-size: 1rem; font-weight: 600;
  color: #57606a; text-transform: uppercase; letter-spacing: .06em;
}
.dark .cmp-lang-heading { color: #8b949e; }

.cmp-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; }
@media (max-width: 900px) { .cmp-grid { grid-template-columns: 1fr; } }
.cmp-col-label { font-size: .8rem; font-weight: 600; color: #57606a; margin-bottom: .4rem; }
.dark .cmp-col-label { color: #8b949e; }

.shiki-target pre {
  margin: 0; padding: 1rem 1.35rem;
  border-radius: .5rem;
  font-family: 'JetBrains Mono', monospace;
  font-size: .875rem; line-height: 1.6;
  overflow-x: auto;
  tab-size: 4;
  -moz-tab-size: 4;
}
.shiki-target pre code {
  font-family: inherit;
}`

const pageJS = `
const toggle = document.getElementById("dark-toggle");
const status = document.getElementById("shiki-status");

const saved = localStorage.getItem("kz-demo-theme");
if (saved === "dark") {
  document.body.classList.add("dark");
  document.documentElement.classList.add("dark");
  toggle.checked = true;
}

toggle.addEventListener("change", () => {
  const dark = toggle.checked;
  document.body.classList.toggle("dark", dark);
  document.documentElement.classList.toggle("dark", dark);
  localStorage.setItem("kz-demo-theme", dark ? "dark" : "light");
  if (window.__shikiHL) renderAll(window.__shikiHL, dark ? "github-dark" : "github-light");
});

function renderAll(hl, theme) {
  document.querySelectorAll(".shiki-target").forEach(target => {
    const lang = target.dataset.lang;
    const id = target.id.replace("shiki-", "");
    const src = document.getElementById("src-" + id);
    if (!src) return;
    const code = src.textContent.replace(/<\\\/script>/g, "<\/script>");
    target.innerHTML = hl.codeToHtml(code, { lang, theme });
  });
}


async function init() {
  try {
    const { createHighlighter } = await import("https://esm.sh/shiki@3");
    const hl = await createHighlighter({
      themes: ["github-light", "github-dark"],
      langs: ["go", "javascript", "typescript", "python", "bash", "php", "css", "html"],
    });
    window.__shikiHL = hl;
    renderAll(hl, toggle.checked ? "github-dark" : "github-light");
    status.textContent = "Shiki ready.";
    status.className = "cmp-status ready";
  } catch (err) {
    status.textContent = "Shiki CDN load failed: " + err.message;
    status.className = "cmp-status error";
    console.error(err);
  }
}
init();
`

const pageTmpl = `<!DOCTYPE html>
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
<nav class="site-nav">
  <div class="site-nav-inner">
    <a class="site-brand" href="../nuri/showcase.html">Kazari</a>
    <div class="site-nav-links">
      <a href="../nuri/showcase.html">Showcase</a>
      <a href="comparison.html" class="active" aria-current="page">Nuri vs Shiki</a>
      <a href="../comparison-chroma/comparison-chroma.html">Nuri vs Chroma</a>
    </div>
    <label class="site-dark-toggle"><input type="checkbox" id="dark-toggle"> Dark mode</label>
  </div>
</nav>
<header class="cmp-header">
  <span id="shiki-status" class="cmp-status">Loading Shiki from CDN...</span>
  <p>Left: Nuri (Go, pre-rendered at build time). Right: Shiki (JS, loaded from CDN).</p>
</header>
<main>%s</main>
<script>%s</script>
<script type="module">%s</script>
</body>
</html>`

func escapeForScript(s string) string {
	return strings.ReplaceAll(s, "</script>", "<\\/script>")
}

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
	)

	var rows strings.Builder
	for _, s := range snippets.All {
		noFrame := kazari.FrameNone
		rendered, err := engine.Render(s.Code, kazari.Options{
			Lang:  s.Lang,
			Frame: &noFrame,
		})
		if err != nil {
			log.Fatalf("render %s: %v", s.ID, err)
		}
		rows.WriteString(fmt.Sprintf(sectionTmpl,
			s.ID, s.Label, rendered,
			s.ID, s.Lang,
			s.ID, escapeForScript(s.Code),
		))
	}

	page := fmt.Sprintf(pageTmpl,
		engine.CSS(),
		pageCSS,
		navCSS,
		rows.String(),
		engine.JS(),
		pageJS,
	)

	if err := os.WriteFile("demo/comparison/comparison.html", []byte(page), 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Println("Written: demo/comparison/comparison.html")
}
