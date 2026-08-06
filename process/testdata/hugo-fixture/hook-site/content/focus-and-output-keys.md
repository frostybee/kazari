---
title: "Focus and Output"
---

Focus ranges reach the engine through the render hook. Before the hook
learned this key, focus="4-6" passed through quoted and every range was
silently dropped.

```go {focus="4-6"}
package main

import "fmt"

func run() int {
	fmt.Println("run")
	return 42
}
// end
```

The output panel keys are camelCase in the Kazari meta grammar, so they
only survive if the hook translates them. This fixture's engine leaves
outputPanel off, which keeps the separator as literal code; the test
asserts the emitted meta string instead, and the example site under
examples/hugo covers the rendered panel.

```go {withoutput="true" outputlabel="Result" outputcollapsed="false"}
package main

func main() {
	println("hello")
}
---output---
hello
```
