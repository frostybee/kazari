// Command kazari post-processes static site output, upgrading plain code
// blocks to framed, syntax highlighted Kazari blocks. It works on any folder
// of built HTML files regardless of which generator produced them.
package main

import (
	"fmt"
	"io"
	"os"
)

const usage = `kazari upgrades code blocks in built HTML to framed, syntax
highlighted blocks with copy buttons, line numbers, and dual themes.

Usage:
  kazari process [dir] [flags]   upgrade code blocks under dir (default ".")
  kazari themes                  list bundled syntax theme names
  kazari version                 print the kazari version
  kazari help                    show this help

Run "kazari process -h" for the process flags.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "process":
		return runProcess(args[1:], stdout, stderr)
	case "themes":
		return runThemes(stdout, stderr)
	case "version":
		return runVersion(stdout)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "kazari: unknown command %q\n\n", args[0])
		fmt.Fprint(stderr, usage)
		return 2
	}
}
