# Examples

Runnable sites that show what `kazari process` does to real static site
output. These are demos meant to be browsed and copied from, unlike the
minimal fixtures under `process/testdata/`, which exist to be asserted
against and are deliberately unstyled.

| Example | Generator | What it shows |
|---|---|---|
| [`hugo/`](hugo/) | Hugo | The tier 2 render hook, per block options, dual themes, collapse, and the output panel |

## The contract these share

`kazari process` walks built HTML and upgrades every code block it finds.
That alone needs nothing from the generator: plain `<pre><code>` blocks are
recognised and rewritten, and this is the tier 1 path that works with any
tool.

Per block options are the part that needs help. A generator writes its own
markup and drops the fence info string on the way, so `title="main.go"` or
`mark="3-5"` never reaches the built page. Kazari's answer is a single
attribute:

```html
<pre><code class="language-go" data-kz-meta='go title="main.go" {3-5}'>...</code></pre>
```

Any generator that can emit `data-kz-meta` carrying the fence's Kazari meta
string gets the full feature set. The processor reads the attribute,
passes the value verbatim to the engine, and removes it while rewriting
the block. This is the tier 2 path.

The Hugo example does this with a render hook. Other generators need their
own bridge, which is why the folder is split by generator rather than
holding one site. Jekyll and Eleventy have no equivalent hook mechanism
today, so neither has an example yet.

## Running any of them

Each example has its own README with the exact commands. All of them
follow the same two steps: build the site with its generator, then run
`kazari process` over the output directory.
