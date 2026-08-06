---
title: "Annotations"
---

## Marked lines

`mark` takes comma separated line numbers and ranges.

```go {title="cache.go" mark="5-7,11"}
package cache

import "sync"

type Store struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.items[key]
	return v, ok
}
```

## Hugo's own spelling

`hl_lines` is Hugo's native attribute for highlighted lines. The hook
translates it to the same marked lines, so a site migrating from Chroma
keeps its existing fences working unchanged.

```go {hl_lines="3-4"}
func sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}
```

## Inserted and deleted lines

`ins` and `del` mark lines as added or removed without turning the block
into a diff, so the code stays syntax highlighted. `add` and `rem` are
accepted spellings of the same two.

```go {title="retry.go" ins="6-8" del="4"}
package client

func send(req *Request) error {
	return transport.Do(req)
	for attempt := 0; attempt < 3; attempt++ {
		if err := transport.Do(req); err == nil {
			return nil
		}
	}
	return ErrExhausted
}
```

## Focused lines

`focus` dims everything outside the range, which points at one part of a
long block without hiding the rest.

```go {title="middleware.go" focus="6-8"}
package server

import "net/http"

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
```

## Hybrid diff

A `diff` fence with `lang` set highlights the code in that language and
keeps the diff gutter, instead of colouring every line as diff syntax.

```diff {lang="go" title="config.go"}
 type Config struct {
 	Addr string
-	TLS  bool
+	TLS      bool
+	CertFile string
 }
```

## Links inside code

`@[label](url)` becomes a real link with an external icon and the syntax
itself disappears from the rendered code. This needs `links: true` in the
config.

```go {title="doc.go" showlinenumbers="false"}
// Package client implements the @[HTTP API](https://example.org/api).
// Rate limits are described in the @[quota guide](https://example.org/quota).
package client
```

## What cannot come through a Hugo fence

Kazari also supports inline text markers and regex markers, written as a
bare `"text"` or `/pattern/` in the meta string. Hugo parses only
`key="value"` pairs inside the brace group and drops bare tokens, so
neither reaches the hook. A site that needs them has to extend the hook
template locally.
