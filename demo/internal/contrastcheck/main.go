package main

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

func luminance(hex string) float64 {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return 0
	}
	var rgb [3]float64
	for i := 0; i < 3; i++ {
		v := 0.0
		for j := 0; j < 2; j++ {
			c := hex[i*2+j]
			switch {
			case c >= '0' && c <= '9':
				v = v*16 + float64(c-'0')
			case c >= 'a' && c <= 'f':
				v = v*16 + float64(c-'a'+10)
			case c >= 'A' && c <= 'F':
				v = v*16 + float64(c-'A'+10)
			}
		}
		v /= 255.0
		if v <= 0.04045 {
			rgb[i] = v / 12.92
		} else {
			rgb[i] = math.Pow((v+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*rgb[0] + 0.7152*rgb[1] + 0.0722*rgb[2]
}

func contrastRatio(hex1, hex2 string) float64 {
	l1, l2 := luminance(hex1), luminance(hex2)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

const goCode = `package main

import "fmt"

func main() {
  name := "world"
  fmt.Printf("Hello, %s!\n", name)
  for i := 0; i < 3; i++ {
    fmt.Println(i)
  }
}`

func main() {
	ctx := context.Background()
	hl, _ := nuri.New(ctx, nuri.WithFS(core.FS()), nuri.WithMinContrast(0))
	defer hl.Close(ctx)
	rawHL := kazarinuri.New(ctx, hl)

	darkThemes := []string{
		"solarized-dark", "nord", "rose-pine", "rose-pine-moon",
		"catppuccin-frappe", "catppuccin-macchiato", "catppuccin-mocha",
		"dracula", "dracula-soft", "material-theme", "material-theme-darker",
		"material-theme-ocean", "material-theme-palenight",
		"night-owl", "one-dark-pro", "tokyo-night", "monokai",
		"everforest-dark", "gruvbox-dark-medium", "gruvbox-dark-soft",
		"ayu-dark", "ayu-mirage", "vitesse-dark", "poimandres",
		"min-dark", "slack-dark", "houston", "laserwave",
		"synthwave-84", "kanagawa-wave", "kanagawa-dragon",
		"github-dark", "github-dark-dimmed",
	}

	type result struct {
		theme       string
		bg          string
		minRatio    float64
		lowCount    int
		totalTokens int
	}
	var results []result

	for _, t := range darkThemes {
		info, err := rawHL.GetThemeColors(t)
		if err != nil {
			continue
		}

		tokens, err := rawHL.Tokenize(goCode, "go", t)
		if err != nil {
			continue
		}

		minR := 21.0
		low := 0
		total := 0
		for _, line := range tokens {
			for _, tok := range line {
				if tok.Color == "" || strings.TrimSpace(tok.Content) == "" {
					continue
				}
				total++
				r := contrastRatio(tok.Color, info.BG)
				if r < minR {
					minR = r
				}
				if r < 4.5 {
					low++
				}
			}
		}
		results = append(results, result{t, info.BG, minR, low, total})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].lowCount > results[j].lowCount })

	fmt.Printf("%-28s %-9s %8s %6s %6s\n", "THEME", "BG", "MIN_CR", "LOW", "TOTAL")
	fmt.Println(strings.Repeat("-", 65))
	for _, r := range results {
		pct := 0.0
		if r.totalTokens > 0 {
			pct = float64(r.lowCount) / float64(r.totalTokens) * 100
		}
		fmt.Printf("%-28s %-9s %8.2f %4d %4d  (%.0f%%)\n",
			r.theme, r.bg, r.minRatio, r.lowCount, r.totalTokens, pct)
	}
}
