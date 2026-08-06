---
title: "Build Pipelines"
description: "Add kazari process to a static site build pipeline."
sidebar:
  order: 5
---

Run `kazari process` after the static site generator writes its output directory and before the deploy step uploads it. The upgraded HTML replaces the generator's output in place; no deploy configuration changes.

## GitHub Actions: Hugo to GitHub Pages

A single step between the Hugo build and the GitHub Pages upload is enough. The processor reads from `public/`, upgrades every code block, and writes the assets alongside the HTML.

```yaml
- name: Build site
  run: hugo --minify

- name: Enhance code blocks
  run: >-
    go run github.com/frostybee/kazari/cmd/kazari@latest
    process --config kazari.config.yaml ./public

- uses: actions/upload-pages-artifact@v4
  with:
    path: public
```

The runner needs a Go toolchain; `actions/setup-go` provides one when the image does not.

## General pattern

The processor only needs a directory of HTML files, so the same three-step shape works regardless of the generator or CI platform.

1. Build the site with whatever generator the project uses.
2. Run `kazari process --config kazari.config.yaml <output-dir>`.
3. Deploy the output directory as before.

For Jekyll the output directory is `_site`, for Eleventy `_site` or `dist`, for mdBook `book`, for Sphinx `_build/html`, for Zola `public`.

## Pin a version

`@latest` pulls whatever version exists at build time, which can break a pipeline without a code change. Pinning to a version tag keeps builds reproducible.

```bash
go run github.com/frostybee/kazari/cmd/kazari@v1.1.0 process --config kazari.config.yaml ./public
```

## Verify idempotency

Processing the same output twice should produce no changes. Adding a `--check` step after the real run catches config drift or nondeterministic templates that would otherwise go unnoticed.

```yaml
- name: Process code blocks
  run: >-
    go run github.com/frostybee/kazari/cmd/kazari@latest
    process --config kazari.config.yaml ./public

- name: Verify processing is stable
  run: >-
    go run github.com/frostybee/kazari/cmd/kazari@latest
    process --check --config kazari.config.yaml ./public
```

The second run must exit `0`. Exit code `1` means something still changes on every run, which usually points at a config drift or a nondeterministic template.
