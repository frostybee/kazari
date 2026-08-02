---
title: "Highlighted Lines"
---

A fenced block with hl_lines so Chroma emits hl classes on lines 3 and 4.
Classes mode is forced per block because inline styles mode carries no hl
class, which makes the highlight unrecoverable by design.

```go {hl_lines=[3,4], noClasses=false}
package main

func run() int {
	return 42
}
```
