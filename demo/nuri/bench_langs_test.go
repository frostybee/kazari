package main

import (
	"context"
	"testing"
	"time"

	"github.com/frostybee/kazari"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func TestPerLangCost(t *testing.T) {
	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		t.Fatal(err)
	}
	defer hl.Close(ctx)

	engine := kazari.New(
		kazari.WithHighlighter(kazarinuri.New(ctx, hl)),
		kazari.WithThemes("github-light", "github-dark"),
		kazari.WithMinify(false),
	)

	langs := []struct {
		name string
		code string
	}{
		{"go", "package main\nfunc main() {}"},
		{"javascript", "const x = 1;"},
		{"bash", "echo hello"},
		{"powershell", "Write-Output 'hi'"},
		{"jsx", "<div>{x}</div>"},
		{"python", "print('hello')"},
		{"text", "plain text"},
	}

	for _, l := range langs {
		start := time.Now()
		_, err := engine.Render(l.code, kazari.Options{Lang: l.name, Title: "test"})
		dur := time.Since(start)
		if err != nil {
			t.Errorf("%s: %v", l.name, err)
			continue
		}
		t.Logf("%-12s first render: %v", l.name, dur)

		start = time.Now()
		_, _ = engine.Render(l.code, kazari.Options{Lang: l.name, Title: "test2"})
		t.Logf("%-12s second render: %v", l.name, time.Since(start))
	}
}
