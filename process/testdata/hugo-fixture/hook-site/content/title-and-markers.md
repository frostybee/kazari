---
title: "Title and Markers"
---

A fenced block whose title, mark, ins, and del all come from the render
hook's data-kz-meta attribute rather than translated markup.

```go {title="main.go" mark="3" ins="6-7" del="9"}
package main

import "fmt"

func run() int {
	fmt.Println("run")
	return 42
}
// end
```
