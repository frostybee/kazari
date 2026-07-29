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
| `kazari-logo.svg` | `docs/public/images/kazari.svg` | Docs site homepage hero |
| `kazari-favicon.svg` | `docs/public/favicon.svg` | Docs site favicon, via `site.favicon` in `docs/sarde.yaml` |
| `kazari-favicon.svg` | `docs/public/favicon.svg` | Docs site header logo, via `site.logo` in `docs/sarde.yaml`. Same copy as the favicon, no second file. |
| `hero-light.svg` | `docs/public/images/hero-light.svg` | Available to the docs site, currently unreferenced |
| `hero-dark.svg` | `docs/public/images/hero-dark.svg` | Available to the docs site, currently unreferenced |

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
