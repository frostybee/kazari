---
title: "Themes"
---

## Two themes in one build

`kazari.config.yaml` names `github-light` and `dracula`. Kazari tokenizes
every block against both and bakes both sets of colours into the HTML as
inline custom properties. Switching themes changes which set applies and
nothing else. No JavaScript runs, no tokens are recomputed, and the block
never flashes an intermediate state.

The config also declares which selector marks the dark state:

```yaml {title="kazari.config.yaml" mark="2-3"}
darkMode:
  kind: selector
  selector: ".dark"
```

`kind: selector` tells Kazari to ignore the operating system preference
entirely and trust the class. The `selector` key is required whenever
`kind` is `selector` or `both`.

## The two toggles do different things

The **Dark** button in the page header belongs to this site, not to
Kazari. It lives in `static/theme-toggle.js`, flips the `.dark` class on
`<html>`, and stores the choice in `localStorage`. Kazari ships no page
level switch, so any site that wants one writes it.

The small sun and moon button in each block toolbar is Kazari's own, from
`themeToggleButton: true`. It sets `data-kz-theme` on that one block and
leaves the rest of the page alone. Use it to read one block in the other
theme without changing the whole page.

Toggle the header button and this block follows the page. Toggle the
button in its toolbar and it stops following until the page toggle catches
up with it.

```go {title="themes.go" mark="4"}
package theme

func Resolve(pref string) string {
	if pref == "" {
		return "system"
	}
	return pref
}
```

## Pinning one block to one theme

`theme` fixes a block to a named theme in both page states, which is
useful when a screenshot or a brand colour has to stay put.

```go {title="always-dracula.go" theme="dracula"}
package brand

var Accent = "#bd93f9"
```

```go {title="always-light.go" theme="github-light"}
package brand

var Paper = "#ffffff"
```

## Two themes on one block

`theme` also accepts a light and dark pair, comma separated, which
overrides the site pair for that block only. This one uses Nord in dark
mode instead of Dracula.

```go {title="nord-at-night.go" theme="github-light,nord"}
package brand

var Frost = "#88c0d0"
```

Run `kazari themes` to list every bundled theme name.

## No flash, by construction

Two things make the switch instant. Kazari sets `transition: none` with
`!important` on a code block and every element inside it, so nothing
animates between colour sets. That reset exists because the usual theme
toggle shortcut is a universal rule such as `* { transition: ... }`, which
would otherwise fade every token on every switch. This site declares its
own transitions on named page chrome selectors in `static/site.css`, so
the reset never has to fire.

The second rule is the inline script in `layouts/_default/baseof.html`. It
applies the stored class before the first paint, so a reload in dark mode
never shows a light frame first.
