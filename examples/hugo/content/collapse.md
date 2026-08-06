---
title: "Collapse"
---

## The automatic threshold

`kazari.config.yaml` sets `collapsible.lineThreshold: 18`, so any block
longer than that collapses to a six line preview with a toggle in the
toolbar. This block needs no attributes at all.

```go {title="pipeline.go"}
package pipeline

import (
	"context"
	"errors"
	"sync"
)

type Stage func(context.Context, []byte) ([]byte, error)

type Pipeline struct {
	stages []Stage
	mu     sync.Mutex
}

func New(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

func (p *Pipeline) Run(ctx context.Context, in []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.stages) == 0 {
		return nil, errors.New("pipeline: no stages")
	}
	out := in
	for _, stage := range p.stages {
		next, err := stage(ctx, out)
		if err != nil {
			return nil, err
		}
		out = next
	}
	return out, nil
}
```

## Opting one block out

`nocollapse="true"` keeps a long block fully expanded even though it
crosses the threshold.

```go {title="constants.go" nocollapse="true"}
package config

const (
	DefaultAddr        = ":8080"
	DefaultReadTimeout = 15
	DefaultIdleTimeout = 60
	DefaultMaxHeader   = 1 << 20
	DefaultLogLevel    = "info"
	DefaultCacheTTL    = 300
	DefaultRetries     = 3
	DefaultBackoffMS   = 250
	DefaultPoolSize    = 32
	DefaultQueueDepth  = 128
	DefaultBatchSize   = 64
	DefaultFlushEvery  = 5
	DefaultShutdownSec = 30
	DefaultMetricsPath = "/metrics"
	DefaultHealthPath  = "/healthz"
	DefaultProfilePath = "/debug/pprof"
)
```

## Explicit ranges

`collapse` takes ranges of its own, which folds boilerplate in the middle
of a short block that the threshold would never touch.

```go {title="main.go" collapse="3-9"}
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
```

## Collapse styles

The default style renders each folded region as a summary row. The
`collapsible-start`, `collapsible-end`, and `collapsible-auto` styles
attach the control to the edge of the region instead.

```go {title="collapsible-start" collapsestyle="collapsible-start" collapse="2-6"}
func setup() error {
	if err := loadEnv(); err != nil {
		return err
	}
	if err := openDB(); err != nil {
		return err
	}
	return serve()
}
```

```go {title="collapsible-end" collapsestyle="collapsible-end" collapse="2-6"}
func teardown() error {
	if err := flushQueue(); err != nil {
		return err
	}
	if err := closeDB(); err != nil {
		return err
	}
	return nil
}
```

## A per block threshold

`collapsethreshold` overrides the config wide number for one block. At
`8`, this ten line block collapses even though the site default of `18`
would leave it alone.

```go {title="handlers.go" collapsethreshold="8"}
func routes(mux *http.ServeMux) {
	mux.HandleFunc("/", index)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/search", search)
	mux.HandleFunc("/healthz", health)
	mux.HandleFunc("/metrics", metrics)
	mux.HandleFunc("/version", version)
}
```
