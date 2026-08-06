---
title: "Output"
---

## Code and its result in one block

`withoutput="true"` splits a fence at a separator line into a code panel
and an output panel. Everything above `---output---` is highlighted as
code. Everything below it is shown verbatim, with no highlighting, because
program output is not source.

This needs `outputPanel: true` in `kazari.config.yaml`. Without it the
separator stays literal text in the middle of the code, which is the whole
of the failure mode: nothing errors, the block only looks wrong.

```go {title="main.go" withoutput="true"}
package main

import "fmt"

func main() {
	for i := 1; i <= 3; i++ {
		fmt.Printf("tick %d\n", i)
	}
}
---output---
tick 1
tick 2
tick 3
```

## Naming the panel

`outputlabel` replaces the default label on the panel's toggle.

```python {title="stats.py" withoutput="true" outputlabel="Console"}
values = [4, 8, 15, 16, 23, 42]
print("count:", len(values))
print("mean:", sum(values) / len(values))
---output---
count: 6
mean: 18.0
```

## Starting collapsed

`outputcollapsed="true"` hides the panel until the reader opens it, which
suits long or noisy output. `outputcollapsed="false"` forces it open on a
site whose config collapses output by default.

```bash {title="build.sh" withoutput="true" outputcollapsed="true" outputlabel="Build log"}
go build ./...
go test ./...
---output---
ok  	github.com/frostybee/kazari	0.412s
ok  	github.com/frostybee/kazari/goldmark	0.188s
ok  	github.com/frostybee/kazari/process	1.307s
```

## Terminal blocks work the same way

A shell block still gets its terminal frame, and the output panel sits
below it.

```bash {withoutput="true"}
kazari process --config kazari.config.yaml public
---output---
7 files, 24 blocks upgraded, 1 skipped, 0 suppressed, 9 changed
```

## A note on the attribute names

The Kazari meta grammar spells these keys `withOutput`, `outputLabel`, and
`outputCollapsed`. Hugo lowercases every attribute name before the render
hook sees it, so the hook translates `withoutput` back to `withOutput` on
the way out. Write them lowercase in Hugo fences and the hook handles the
rest.
