---
title: "Collapse"
---

Two fenced blocks: one with the default collapse style, one with an
explicit collapsible-start style, both driven entirely by the hook.

```go {collapse="4-9"}
package main

import "fmt"

func run() int {
	x := 1
	y := 2
	return x + y
}
```

```go {collapsestyle="collapsible-start" collapse="2-5"}
package main

func helper() {
	// step one
	// step two
	// step three
}
```
