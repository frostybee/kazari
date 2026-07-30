# Kazari brand assets

Master copies of the Kazari logo, favicon, and hero illustrations. This folder is
the source of truth. Files served by the documentation site are copies, because
Sarde only serves static assets from `docs/public/` and its public directory is a
hardcoded constant with no configuration hook.

Copying is manual. After editing a master here, copy it to every destination
listed below.

## Files and their destinations

| Master | Copy to | Used by |
|--------|---------|---------|
| `kazari-logo.svg` | no copy needed | `README.md` header image, referenced here directly |
| `kazari-logo.svg` | `docs/public/images/kazari.svg` | Nothing at present. Kept because it is the only full logo the docs site can serve |
| `kazari-favicon.svg` | `docs/public/favicon.svg` | Docs site favicon, via `site.favicon` in `docs/sarde.yaml` |
| `kazari-favicon.svg` | `docs/public/favicon.svg` | Docs site header logo, via `site.logo` in `docs/sarde.yaml`. Same copy as the favicon, no second file. |
| `hero-light.svg` | `docs/public/images/hero-light.svg` | Docs site homepage hero, light variant, via `homepage.hero.image.light` |
| `hero-dark.svg` | `docs/public/images/hero-dark.svg` | Docs site homepage hero, dark variant, via `homepage.hero.image.dark` |
| `kazari-logo.png` | `docs/public/images/kazari-logo.png` | Social card corner mark and watermark, via `social_cards.logo` in `docs/sarde.yaml`. Rasterized from `kazari-logo.svg` at 512 px, since cards composite raster images only; the panel's rounded corners were converted to an explicit path because the offline oksvg rasterizer drops `rx` on `rect` |

`favicon-variants.html` is a standalone comparison page. Open it in a browser to
see every favicon candidate rendered at 16, 32, and 64 pixels against both a light
and a dark browser tab bar. Edit the `variants` array to try new combinations.

## Palette

Taken from `kazari-logo.svg`.

| Colour | Hex | Role |
|--------|-----|------|
| Panel light | `#23243A` | Logo panel gradient start |
| Panel dark | `#171821` | Logo panel gradient end |
| Divider | `#353650` | Logo title bar rule |
| Cyan | `#43C6D9` | Ribbon gradient start, terminal dot |
| Violet | `#7C5CFC` | Ribbon gradient midpoint, primary brand colour |
| Pink | `#FF6B81` | Ribbon gradient end, terminal dot |
| Amber | `#F7C85E` | Terminal dot |
| Near-white | `#F5F3FF` | K stem in the logo |

The shipping favicon uses a different tile and a lightened arm gradient. Both
depart from the logo deliberately, for the reasons in the next section.

## Why the favicon is not the logo

The logo is a 512-unit canvas built for display at 160 pixels and up. Scaling it
to a 16 pixel favicon divides every dimension by 32, which puts most of its
detail below one pixel:

| Element | In the logo | At 16px |
|---------|-------------|---------|
| Terminal dots | `r="12"` | 0.75px across |
| Title bar rule | `stroke-width="12"` | 0.38px |
| Code lines | `stroke-width="16"` | 0.5px |
| K stem | `stroke-width="30"` | 0.94px |
| Panel inset | `x="72"` | 2.25px of dead padding per side |

The favicon therefore drops the dots, the rule, and the code lines, and redraws
the K as three plain strokes instead of the logo's folded ribbon. The ribbon has
an interior notch about 1.5px wide at favicon size; once it closes up the glyph
stops reading as a K.

Current favicon values:

| Part | Value |
|------|-------|
| Tile | `#5B6BA8` to `#3E4A7A` |
| Stem | `#FFFFFF` |
| Arms | `#43C6D9` to `#B3A4FF` to `#FF8FA3` |

The tile is mid-tone rather than the logo's near-black so the icon keeps its
edges against dark browser chrome, and the arm gradient is lightened from the
logo's so it stays legible against that lighter tile.

## Why the header logo is the favicon art

The docs header renders its logo at Sarde's `logo-height` token, 1.75rem or 28px
by default. That is an 18x reduction of the 512-unit canvas, which puts the logo's
code lines at 0.87px, its title bar rule at 0.66px, and its terminal dots at 1.3px
across. Its `x="72"` panel inset alone burns 3.9px of dead padding per side out of
28. Sizing up is not a fix: the logo needs 160px and the header is 64px tall.

The favicon's `stroke-width="56"` renders at about 3.1px at that size, so the K
still reads. `site.logo` therefore points at `/favicon.svg`, the same copy the
favicon uses. The mid slate blue tile that keeps the icon's edges against dark
browser chrome does the same job against both header themes.

A dedicated header asset is only worth authoring if the header should carry more
of the logo's character than a three-stroke K, and that means new artwork tuned
for roughly 24 to 40 pixels, not a rescale of either existing file.

## Why the heroes expand the logo instead of rescaling it

The hero art used to be Sarde's template placeholder, a generic card with two
button shapes and a checkmark, in the theme's default indigo. Nothing in it was
Kazari. The homepage then pointed `image.light` and `image.dark` at
`kazari.svg` instead, which put the same 512-unit square in both slots.

Both are now `hero-light.svg` and `hero-dark.svg`, on a 400 by 300 landscape
canvas. The hero slot is the one place with enough room for detail: Sarde renders
it at `width: 100%` inside a flex column, capped at `20rem` only below the mobile
breakpoint, so the art gets roughly 320 to 480 pixels. That is the opposite of
the favicon problem. Nothing has to be simplified away, so nothing is.

The illustration is therefore the logo's own subject at full size. The logo is a
code window with a ribbon K; the hero is what that window contains once Kazari
renders into it: an editor frame, a title bar with a file name and a copy button,
a line number gutter, tokenized code, and one marked line with the violet gutter
bar. The logo mark then rides the lower right corner, and it is the logo, not the
favicon: the tile carries the logo's panel gradient and its 0.196 corner radius
ratio, and the K is the logo's own `M160 196v152` stem and `M160 270l104-76...`
folded ribbon, referenced verbatim and placed with

    transform="translate(344 236) scale(0.28) translate(-221.5 -272.5)"

The inner translate centres the mark's own bounding box, 145 to 298 by 181 to
364 in logo coordinates. Because the ribbon gradient is `userSpaceOnUse` and is
referenced from inside that group, its logo space coordinates still line up with
the logo space path data, so the gradient needs no adjustment either.

Do not substitute the favicon's three-stroke K here. Its notch argument is
calibrated for 16 pixels; the mark tile renders at roughly 60 to 85, where the
folded ribbon reads without trouble. The stem's `stroke-width="30"` lands at 8.4
units, about 9 pixels.

The two files share their geometry exactly, so the light and dark variants swap
without the hero shifting. What differs is the code block surface and the token
hues. The logo palette is tuned for a near-black panel, and its cyan, amber, and
pink wash out on paper, so the light variant darkens them:

| Role | Dark | Light |
|------|------|-------|
| Keyword | `#7C5CFC` | `#6B4EE6` |
| String | `#43C6D9` | `#1B8B9E` |
| Number | `#F7C85E` | `#A8760C` |
| Identifier | `#C9CCE8` | `#3A3D57` |
| Comment | `#565B7E` | `#8A8EA8` |
| Punctuation | `#4A4E6E` | `#BCBFD6` |

Three elements deliberately do not change between variants: the terminal dots,
the violet line marker, and the logo mark, whose tile stays dark in both because
the logo is dark panelled by nature. Each reads on either surface, and holding
them fixed is what keeps the two files recognisable as one illustration.

On the dark variant the mark tile and the frame share a fill, so a `#353650` rim
and the plate beneath it are what lift the mark off the frame. The light variant
needs neither, and omits the rim.

Depth comes from an offset plate rather than a `filter`, so the art needs no
filter support and stays resolution independent. The token colours are the brand
palette cast into syntax roles; they are not sampled from any shipped Kazari
theme, and they are not meant to track one.

## archive/

Rejected favicon candidates, kept so the reasoning is not lost. Each file
records why it was set aside in a comment at the top.

| File | Why it was rejected |
|------|---------------------|
| `favicon-v1-dark-panel-ribbon.svg` | Reused the logo's ribbon path; the interior notch closes at 16px and the panel is near-black |
| `favicon-v2-violet-solid-k.svg` | Solid white K drops the gradient and the brand colour identity |
| `favicon-v3-slate-dark.svg` | Tile still too dark against dark browser chrome |
| `favicon-v4-pale-tile.svg` | Tile too light |

## Raster exports

`.gitignore` ignores `*.png` repository-wide. A negation for `brand/**/*.png`
keeps exports in this folder tracked. Verify with `git check-ignore -v` before
assuming a new binary asset was committed.
