package unwrap

import "testing"

func TestCollapseRanges(t *testing.T) {
	cases := []struct {
		in   []int
		want string
	}{
		{nil, ""},
		{[]int{1}, "1"},
		{[]int{1, 2, 3}, "1-3"},
		{[]int{3, 4, 5, 9}, "3-5,9"},
		{[]int{5, 6, 8, 9, 10}, "5-6,8-10"},
		{[]int{2, 4, 6}, "2,4,6"},
	}
	for _, c := range cases {
		if got := collapseRanges(c.in); got != c.want {
			t.Errorf("collapseRanges(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestBuildMeta(t *testing.T) {
	if got := buildMeta("go", []int{2, 3}); got != "go {2-3}" {
		t.Fatalf("got %q", got)
	}
	if got := buildMeta("go", nil); got != "go" {
		t.Fatalf("got %q", got)
	}
	if got := buildMeta("", nil); got != "" {
		t.Fatalf("got %q", got)
	}
	if got := buildMeta("", []int{1}); got != "{1}" {
		t.Fatalf("got %q", got)
	}
}

func TestCollectHighlightedLines_IndependentOfGutterExclusion(t *testing.T) {
	// Line counting and gutter exclusion are separate walks over separate
	// class vocabularies. The gutter digits must vanish from recovered text
	// while the line counter still sees every line, including gutter
	// carrying ones, at its correct 1 based position.
	fragment := `<span class="line"><span class="ln">1</span><span class="cl">a
</span></span><span class="line hl"><span class="ln">2</span><span class="cl">b
</span></span><span class="line"><span class="ln">3</span><span class="cl">c
</span></span>`
	tokens, err := TokenizeFragment([]byte(fragment))
	if err != nil {
		t.Fatal(err)
	}
	hl := collectHighlightedLines(tokens)
	if len(hl) != 1 || hl[0] != 2 {
		t.Fatalf("highlighted = %v, want [2]", hl)
	}
	if got := recoverText(tokens); got != "a\nb\nc\n" {
		t.Fatalf("recovered %q", got)
	}
}
