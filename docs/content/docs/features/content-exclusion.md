---
title: "Content Exclusion"
description: "Prevent documentation-site prose styles from leaking into code blocks."
sidebar:
  order: 16
---

Every rendered code block includes the `not-content` class on its wrapper element. This prevents documentation-site prose styles from leaking into Kazari's code block styling.

## How it works

Kazari adds `not-content` to the outermost `<div>` wrapper of every code block:

```html
<div class="kazari-block not-content">
  ...
</div>
```

This class is a convention used by documentation-site frameworks such as Astro Starlight, Tailwind Typography (`prose`), and Docusaurus. These frameworks apply broad typography resets to elements like `pre`, `code`, `figure`, and others inside prose containers. The `not-content` marker tells the host site's CSS to exclude the element and its descendants from those resets.

Kazari itself has no CSS rules for `.not-content`. It is a marker for the host site's stylesheet to interpret. If the host site does not use a prose/typography system, the class has no effect.

## No configuration needed

The `not-content` class is always present on every rendered block. No option is needed to enable it.

`WithContentExclusion` and the `contentExclusion` config file key still exist for backward compatibility but are no-ops. The class is applied regardless of their value.

Code groups also receive the class on the group wrapper:

```html
<div class="kazari-block kz-group not-content">
  ...
</div>
```

## Host-site integration

If the host site uses a prose or typography system that applies broad styles to `pre`, `code`, `figure`, or other HTML elements, the `not-content` class ensures those styles do not apply inside Kazari blocks. Most documentation frameworks that support this convention handle it automatically with no additional configuration from the consumer.

For frameworks that use a different exclusion class, override the host-site CSS to also exclude `.kazari-block` elements, or use Kazari's `WithStyleReset` option to apply `all: revert` on code block internals.
