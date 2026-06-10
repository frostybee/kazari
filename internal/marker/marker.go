package marker

import "github.com/frostybee/kazari/internal/config"

// ResolvedLine holds the final marker state for a single line after overlap resolution.
type ResolvedLine struct {
	Type    config.MarkerType
	Label   string
	HasMark bool
}

// ResolveLineMarkers flattens all line markers into a per-line map,
// resolving overlaps by priority: mark(0) < del(1) < ins(2).
func ResolveLineMarkers(markers []config.LineMarker) map[int]ResolvedLine {
	if len(markers) == 0 {
		return nil
	}

	result := make(map[int]ResolvedLine)

	for _, m := range markers {
		priority := int(m.Type)

		for _, lr := range m.Lines {
			isFirstLine := true
			for line := lr.Start; line <= lr.End; line++ {
				existing, exists := result[line]

				if !exists || priority >= int(existing.Type) {
					entry := ResolvedLine{
						Type:    m.Type,
						HasMark: true,
					}
					if m.Label != "" && isFirstLine {
						entry.Label = m.Label
					}
					// Same priority: keep existing label if new entry has none.
					if exists && priority == int(existing.Type) && entry.Label == "" && existing.Label != "" {
						entry.Label = existing.Label
					}
					result[line] = entry
				}
				isFirstLine = false
			}
		}
	}

	return result
}

// ResolveFocusSet returns a set of line numbers that are focused.
func ResolveFocusSet(focusLines []config.LineRange) map[int]bool {
	if len(focusLines) == 0 {
		return nil
	}
	set := make(map[int]bool)
	for _, lr := range focusLines {
		for line := lr.Start; line <= lr.End; line++ {
			set[line] = true
		}
	}
	return set
}
