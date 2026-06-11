package main

import (
	"context"
	"testing"
	"time"

	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func TestInitTiming(t *testing.T) {
	ctx := context.Background()

	start := time.Now()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	initDur := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("nuri.New: %v", initDur)

	start = time.Now()
	_, err = hl.CodeToHTML(ctx, `fmt.Println("hello")`, nuri.CodeToHTMLOptions{
		Lang:  "go",
		Theme: "github-dark",
	})
	firstRender := time.Since(start)
	t.Logf("First CodeToHTML: %v", firstRender)
	if err != nil {
		t.Fatal(err)
	}

	start = time.Now()
	for i := 0; i < 10; i++ {
		_, _ = hl.CodeToHTML(ctx, `fmt.Println("hello")`, nuri.CodeToHTMLOptions{
			Lang:  "go",
			Theme: "github-dark",
		})
	}
	batchDur := time.Since(start)
	t.Logf("Next 10 CodeToHTML: %v (avg %v)", batchDur, batchDur/10)

	hl.Close(ctx)
}
