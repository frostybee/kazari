---
title: "Hugo native hl_lines through the hook"
---

Exercises the hl_lines to mark ranges translation inside the hook template
itself (string form), distinct from the Chroma hl class translation path
covered by the chroma site.

```go {hl_lines="2-3"}
package main

func run() int {
	return 42
}
```
