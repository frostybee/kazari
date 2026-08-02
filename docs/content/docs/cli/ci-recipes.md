---
title: "CI Recipes"
description: "Run kazari process as the last step of a static site build."
sidebar:
  order: 4
---

Run `kazari process` after the static site generator writes its output directory and before the deploy step uploads it. The upgraded HTML replaces the generator's output in place; no deploy configuration changes.

## GitHub Actions: Hugo to GitHub Pages

```yaml
- name: Build site
  run: hugo --minify

- name: Enhance code blocks
  run: go run github.com/frostybee/kazari/cmd/kazari@latest process ./public

- uses: actions/upload-pages-artifact@v4
  with:
    path: public
```

The runner needs a Go toolchain; `actions/setup-go` provides one when the image does not.

## General pattern

The recipe is generator-agnostic because the tool only consumes HTML:

1. Build the site with whatever generator the project uses.
2. Run `go run github.com/frostybee/kazari/cmd/kazari@latest process <output-dir>`.
3. Deploy the output directory as before.

For Jekyll the output directory is `_site`, for Eleventy `_site` or `dist`, for mdBook `book`, for Sphinx `_build/html`, for Zola `public`.

## Pin a version

Prefer a version tag over `@latest` in CI so builds stay reproducible:

```bash
go run github.com/frostybee/kazari/cmd/kazari@v1.1.0 process ./public
```

## Verify a render hook installed correctly

In a pull request workflow, a `--check` step fails the build when output would change, which catches a missing or misconfigured [render hook](/cli/render-hooks/) without writing anything:

```yaml
- name: Verify code blocks are processed
  run: |
    hugo --minify
    go run github.com/frostybee/kazari/cmd/kazari@latest process ./public
    go run github.com/frostybee/kazari/cmd/kazari@latest process --check ./public
```

The second `process` run must exit 0: processing already-processed output is a no-op. Exit code 1 means something still changes on every run, which usually points at a config drift or a nondeterministic template.
