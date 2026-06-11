package main

import (
	"context"
	"testing"
	"time"

	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func TestSingleVsDualTheme(t *testing.T) {
	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		t.Fatal(err)
	}
	defer hl.Close(ctx)

	jsCode := "const x = 1;\nconsole.log(x);"

	// warmup
	hl.CodeToHTML(ctx, jsCode, nuri.CodeToHTMLOptions{Lang: "javascript", Theme: "github-dark"})

	start := time.Now()
	for i := 0; i < 5; i++ {
		hl.CodeToHTML(ctx, jsCode, nuri.CodeToHTMLOptions{Lang: "javascript", Theme: "github-dark"})
	}
	t.Logf("Single theme 5x js: %v (avg %v)", time.Since(start), time.Since(start)/5)

	start = time.Now()
	for i := 0; i < 5; i++ {
		hl.CodeToHTML(ctx, jsCode, nuri.CodeToHTMLOptions{Lang: "javascript", Theme: "github-light"})
		hl.CodeToHTML(ctx, jsCode, nuri.CodeToHTMLOptions{Lang: "javascript", Theme: "github-dark"})
	}
	t.Logf("Dual two-pass 5x js:    %v (avg %v) — legacy: tokenize twice", time.Since(start), time.Since(start)/5)

	start = time.Now()
	for i := 0; i < 5; i++ {
		hl.CodeToHTML(ctx, jsCode, nuri.CodeToHTMLOptions{Lang: "javascript", Themes: map[string]string{
			"light": "github-light",
			"dark":  "github-dark",
		}})
	}
	t.Logf("Dual single-pass 5x js: %v (avg %v) — tokenize once, resolve both (Kazari's path)", time.Since(start), time.Since(start)/5)
}
