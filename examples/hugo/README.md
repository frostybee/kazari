# Kazari with Hugo

A complete Hugo site demonstrating `kazari process`. Six pages, a render
hook, a config file, and a site wide dark mode switch.

## Run it

```sh
cd examples/hugo
hugo --minify
go run github.com/frostybee/kazari/cmd/kazari process --config kazari.config.yaml public
```

With the CLI installed, the second command is `kazari process --config
kazari.config.yaml public`.

Open `public/index.html` by double clicking it; no web server is needed.
To see what the processor contributes, run `hugo --minify` on its own
first and open the same page: the blocks are plain and unstyled, because
this site never links `kazari.css` by hand. The processor writes
`kazari.css` and `kazari.js` into `public/` and injects the tags that
load them.

Browsing the output from disk is the only reason `hugo.toml` sets
`relativeURLs` and `uglyURLs`. Without them Hugo writes `/site.css` and
`/basics/`, which resolve against the drive root under `file://` and
leave the page unstyled with a dead nav. Kazari links its own assets
relative to each page either way, so neither setting belongs in a site
that gets deployed.

`kazari process --check public` reports what a run would change without
writing anything, and exits non zero if anything would change. Run it
after a real pass and it reports nothing, which is how the build stays
safe to run twice.

## What is here

```
hugo.toml                              site config
kazari.config.yaml                     Kazari config, annotated
layouts/_markup/render-codeblock.html  the render hook
layouts/_default/baseof.html           head, header, theme script
layouts/_default/single.html
layouts/index.html
static/site.css                        page chrome and language icons
static/theme-toggle.js                 site wide dark mode
content/                               six pages of code blocks
```

## The render hook

`layouts/_markup/render-codeblock.html` is a verbatim copy of
`integrations/hugo/render-codeblock.html` at the repository root. A test
asserts the two stay identical, so copy the canonical file rather than
this one if they ever differ.

The hook bypasses Chroma and emits a plain block carrying the fence's
options in a `data-kz-meta` attribute. Without it, Kazari still upgrades
every block on the page, but only the language survives the build. Titles,
marked lines, focus ranges, collapse ranges, and the output panel all come
from that attribute.

Hugo lowercases attribute names, so fences here write `showlinenumbers`
and `withoutput` where the Kazari meta grammar spells them
`showLineNumbers` and `withOutput`. The hook translates them back. The
header comment inside the template lists every key it recognises.

Two Kazari features cannot be expressed in a Hugo fence at all: inline
text markers, written as a bare `"text"`, and regex markers, written as
`/pattern/`. Hugo parses only `key="value"` pairs inside the brace group
and drops bare tokens. Extend the hook template locally if a site needs
them.

## Copying this to a real site

1. Copy `integrations/hugo/render-codeblock.html` to
   `layouts/_markup/render-codeblock.html`. On Hugo older than 0.146 use
   `layouts/_default/_markup/render-codeblock.html` instead.
2. Copy `kazari.config.yaml` and trim it to what the site needs.
3. Add `kazari process` after `hugo` in the build.
4. If the site has a dark mode switch, point `darkMode.selector` at
   whatever class or attribute it already sets. Kazari has no page level
   toggle of its own; `static/site.css` and `static/theme-toggle.js` here
   are one way to write one, not part of Kazari.

## Notes on the demo's own CSS

`static/site.css` styles the page around the blocks and stays out of them,
with one exception. `languageIconMode: iconAndText` in the config makes
Kazari open every language badge with an empty
`<span class="kz-lang-icon" data-lang="go">` at the far left of the
toolbar, and Kazari ships no icon art, so the site fills it. The rules at
the bottom of that file set a `background-image` per language, as inline
SVG data URIs, covering only the languages these pages use plus a generic
fallback. They win the cascade without `!important` because Kazari's rules
live in `@layer kazari` and unlayered CSS beats any layer.

The file icon, a second slot Kazari can emit next to the title, is turned
off in the config so each block carries one icon. It is filled the same
way, keyed on `data-ext` instead of `data-lang`.

Reach for a background image rather than `::before` content. Kazari sizes
the slot through `--kz-lang-icon-size` and aligns that box in the toolbar,
so a background scales to it. A glyph injected with `content:` follows the
toolbar font instead of the box, so it overflows the slot and sits off the
badge's baseline.

Transitions in `site.css` are declared on named page chrome selectors, not
through a universal rule. `* { transition: ... }` is the usual shortcut in
theme toggle CSS and it reaches every element inside a code block, fading
the token colours on each switch. Kazari sets `transition: none
!important` on blocks and their descendants to survive that, and naming
the selectors here keeps the reset from having to fire.
