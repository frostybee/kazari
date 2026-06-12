package main

import (
	"context"
	"log"

	"github.com/frostybee/kazari/demo/internal/showcase"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func main() {
	ctx := context.Background()
	highlighter, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		log.Fatalf("nuri.New: %v", err)
	}
	defer highlighter.Close(ctx)

	config := showcase.Config{
		BackendName: "Nuri",
		HTMLFile:    "showcase.html",
		OtherName:   "Chroma",
		OtherHref:   "../chroma/showcase-chroma.html",
	}
	if err := showcase.Generate(".", config, kazarinuri.New(ctx, highlighter)); err != nil {
		log.Fatalf("generate showcase: %v", err)
	}
	log.Println("Written: showcase.html, showcase.css, showcase.js")
}
