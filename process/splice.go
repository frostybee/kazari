package process

import (
	"bytes"
	"fmt"
	"sort"
)

// edit replaces src[start:end] with replacement. An insertion has
// start == end. All offsets are positions in the original file bytes;
// nothing here ever re-finds an offset in mutated output.
type edit struct {
	start       int
	end         int
	replacement []byte
}

// applyEdits builds the output in a single forward pass over one ascending
// event list, copying verbatim between edits. Insertions sort before a
// replacement starting at the same offset, so a tag injected exactly at a
// region boundary lands outside the replaced span. Overlapping edits are an
// internal invariant violation and abort the file rather than corrupt it.
func applyEdits(src []byte, edits []edit) ([]byte, error) {
	sorted := make([]edit, len(edits))
	copy(sorted, edits)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].start != sorted[j].start {
			return sorted[i].start < sorted[j].start
		}
		return sorted[i].end < sorted[j].end
	})

	var out bytes.Buffer
	pos := 0
	for _, e := range sorted {
		if e.start < pos || e.end < e.start || e.end > len(src) {
			return nil, fmt.Errorf("process: overlapping or out of range edit [%d:%d] at position %d", e.start, e.end, pos)
		}
		out.Write(src[pos:e.start])
		out.Write(e.replacement)
		pos = e.end
	}
	out.Write(src[pos:])
	return out.Bytes(), nil
}
