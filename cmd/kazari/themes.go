package main

import (
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"

	"github.com/frostybee/nuri/bundle/core"
)

// bundledThemes enumerates every theme name in Nuri's embedded registry.
// The bundle FS presents themes/<name>.json.gz as virtual <name>.json
// entries, so the theme name is the entry name minus the .json suffix.
// There is no exported list helper in Nuri, and LoadedThemes only reflects
// lazily loaded themes, so walking the FS is the one reliable source.
func bundledThemes() ([]string, error) {
	entries, err := fs.ReadDir(core.FS(), "themes")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		if name != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

func runThemes(stdout, stderr io.Writer) int {
	names, err := bundledThemes()
	if err != nil {
		fmt.Fprintf(stderr, "kazari: listing themes: %v\n", err)
		return 2
	}
	for _, n := range names {
		fmt.Fprintln(stdout, n)
	}
	return 0
}

// validateThemeNames checks the resolved theme pair against the bundled set
// before any engine work happens. The engine never fails on a bad theme
// name (construction and CSS generation silently degrade, and render errors
// are masked when the language is also unknown), so this early check is the
// only place a typo produces a clear message.
func validateThemeNames(light, dark string) error {
	names, err := bundledThemes()
	if err != nil {
		return fmt.Errorf("listing themes: %w", err)
	}
	valid := make(map[string]bool, len(names))
	for _, n := range names {
		valid[n] = true
	}
	for _, name := range []string{light, dark} {
		if valid[name] {
			continue
		}
		msg := fmt.Sprintf("unknown theme %q", name)
		if s := nearestName(name, names); s != "" {
			msg += fmt.Sprintf(", did you mean %q?", s)
		}
		return fmt.Errorf("%s Run \"kazari themes\" to list all bundled themes.", msg)
	}
	return nil
}

// nearestName suggests the closest bundled name within an edit distance of
// two, which covers the typo cases worth suggesting without ever proposing
// something wild.
func nearestName(name string, names []string) string {
	best, bestDist := "", 3
	for _, n := range names {
		if d := editDistance(name, n, bestDist); d < bestDist {
			best, bestDist = n, d
		}
	}
	return best
}

// editDistance computes Levenshtein distance, giving up early once the
// distance can no longer beat the current bound.
func editDistance(a, b string, bound int) int {
	if abs(len(a)-len(b)) >= bound {
		return bound
	}
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		rowMin := cur[0]
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = minInt(prev[j]+1, minInt(cur[j-1]+1, prev[j-1]+cost))
			if cur[j] < rowMin {
				rowMin = cur[j]
			}
		}
		if rowMin >= bound {
			return bound
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
