package main

import (
	"log"

	kazarichroma "github.com/frostybee/kazari/chroma"
	"github.com/frostybee/kazari/demo/internal/showcase"
)

func main() {
	highlighter := kazarichroma.New(kazarichroma.WithStyleMap(map[string]string{
		"github-light": "github",
		"github-dark":  "github-dark",
	}))
	config := showcase.Config{
		BackendName: "Chroma",
		HTMLFile:    "showcase-chroma.html",
		OtherName:   "Nuri",
		OtherHref:   "../nuri/showcase.html",
	}
	if err := showcase.Generate(".", config, highlighter); err != nil {
		log.Fatalf("generate showcase: %v", err)
	}
	log.Println("Written: showcase-chroma.html, showcase.css, showcase.js")
}
