---
title: "Configuration"
description: "Flags, the process config key, theme selection, and dark mode setup for kazari process."
sidebar:
  order: 2
---

Flags override the config file; the config file overrides built-in defaults. Most projects set the dark mode selector once and run with defaults for everything else.

## Flags

| Flag | Default | Description |
|---|---|---|
| `--check` | off | Report would-be changes without writing; exit 1 if any |
| `--config` | auto-discover | Config file path. Without it, `kazari.config.yaml`, `.yml`, or `.json` is discovered in the target directory, then the working directory |
| `--theme-light` | `github-light` | Light syntax theme name (overrides config) |
| `--theme-dark` | `github-dark` | Dark syntax theme name (overrides config) |
| `--assets-base` | relative paths | Fixed asset URL prefix instead of per-file relative paths |
| `--hashed-assets` | `false` | Use content hashed asset filenames |
| `--skip-unlabeled` | `false` | Leave blocks without a detectable language untouched |
| `--concurrency` | `0` (CPU count) | Max files processed concurrently |
| `--verbose` | `false` | Log per-file progress to stderr |

The positional `[dir]` argument defaults to `.` and may appear before or after the flags. A flag overrides a config file value only when it is explicitly passed; an unset flag never clobbers a config value with its zero default.

## Unknown theme names

Theme names are validated against the bundled set before any work starts. A typo produces a suggestion:

```
kazari: unknown theme "github-drak", did you mean "github-dark"? Run "kazari themes" to list all bundled themes.
```

`kazari themes` prints every bundled theme name, sorted, one per line.

## The process config key

The `process:` block in `kazari.config.yaml` configures the CLI. All keys are optional.

```yaml
process:
  skipUnlabeled: true
  assetsBase: "/assets/"
  hashedAssets: false
  concurrency: 4
  maxFileBytes: 33554432
```

| Field | Type | Default | Description |
|---|---|---|---|
| `skipUnlabeled` | `bool` | `false` | Leave blocks without a detectable language untouched instead of rendering a plain text frame |
| `assetsBase` | `string` | `""` | Fixed asset URL prefix; empty means per-file relative paths |
| `hashedAssets` | `bool` | `false` | Emit `kazari-<hash>.css` and `kazari-<hash>.js` instead of `kazari.css` and `kazari.js` |
| `concurrency` | `int` | `0` | Max files processed concurrently; zero means the CPU count |
| `maxFileBytes` | `int64` | `33554432` | Files larger than this (32 MiB) are skipped and reported as errors |

> [!NOTE]
> The `process:` key affects only the `kazari process` command. The Goldmark extension and the library API ignore it.

The same config file also carries the engine-level keys, including `themes:` and `darkMode:`, which the CLI applies when building its engine. See [File-Based Config](/reference/file-based-config/) for the full engine key reference.

## Dark mode

Dark mode is the one setting to configure by hand. Kazari bakes both light and dark token colors into every block; the `darkMode` key decides which CSS mechanism switches between them, and it must match the mechanism the site's theme actually uses. There is no auto-detection.

```yaml
darkMode:
  kind: selector
  selector: ".dark"
```

| Mechanism | `kind` | `selector` | Typical on |
|---|---|---|---|
| OS preference via `prefers-color-scheme` | `mediaQuery` | none | Sites without a manual theme toggle |
| `dark` class on the root element | `selector` | `.dark` | Tailwind and Docusaurus style themes |
| `data-theme` attribute | `selector` | `[data-theme="dark"]` | Many Hugo and Jekyll themes |
| Theme-specific classes | `both` or a custom selector | site specific | mdBook and other theme pickers |

These are typical defaults; third-party themes vary. To find the right value, open the site in a browser, toggle its theme switcher, and watch in devtools which class or attribute changes on the `<html>` element. The engine default is `selector` with `.dark`.

## Collapsing long blocks

Engine keys in the same config file shape what the processor renders. The most useful one for sites without a render hook is `collapsible`, which folds any block over a line count down to a preview:

```yaml
collapsible:
  lineThreshold: 20
  previewLines: 8
  defaultCollapsed: true
```

Adding the `collapsible` key at all is what enables the feature. Unset subfields fall back to these values:

| Field | Type | Default | Description |
|---|---|---|---|
| `lineThreshold` | `int` | `15` | Blocks with more lines than this become collapsible |
| `previewLines` | `int` | `5` | Lines kept visible when collapsed |
| `defaultCollapsed` | `bool` | `true` | Render collapsed on page load |
| `style` | `string` | `"github"` | Section style for range-based collapse |

Four more fields cover indentation and screen reader labels. See [File-Based Config](/reference/file-based-config/#collapsible) for the complete list.

This applies to every block the processor touches, with no changes to any Markdown source. Collapsing an exact line range instead, or exempting a single block, needs a [render hook](/cli/render-hooks/#example-collapse-specific-lines).

## Assets

`kazari.css` and `kazari.js` are written once at the output root, only when their content differs from what is already on disk. Pages reference them with relative `../` paths built from each page's depth, which works unchanged under subpath hosting such as GitHub Pages project sites. A `?v=` query carrying an 8 character content hash busts caches when the stylesheet changes.

`--assets-base` replaces the relative paths with a fixed prefix; a trailing slash on the prefix is normalized. `--hashed-assets` switches to content hashed filenames and drops the `?v=` query, since the hash lives in the name.

Injected link and script tags carry a `data-kazari="assets"` attribute. Re-runs find those tags and update them in place, including after `--assets-base` or `--hashed-assets` changed, so tags never duplicate.

> [!WARNING]
> `--hashed-assets` does not clean up hashed files left behind by earlier runs. Stale `kazari-<oldhash>.css` files accumulate until removed manually.

## Known limitation

> [!NOTE]
> Indentation that a generator replaced with non-breaking spaces decodes to U+00A0, not a regular space, and is preserved as-is. Normalizing would corrupt intentional non-breaking spaces inside string literals, so the processor does not attempt it.
