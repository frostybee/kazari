package unwrap

import "testing"

func TestRougeDivWrapper_Fixtures(t *testing.T) {
	assertRegion(t, rougeDivWrapper{}, "jekyll-rouge-kramdown")
}

func TestRougeHighlightFigure_Fixtures(t *testing.T) {
	assertRegion(t, rougeHighlightFigure{}, "jekyll-rouge-highlight-figure")
}

func TestRougeLinenosTable_Fixtures(t *testing.T) {
	assertRegion(t, rougeLinenosTable{}, "jekyll-rouge-linenos-table")
}
