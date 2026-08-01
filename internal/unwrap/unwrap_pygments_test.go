package unwrap

import "testing"

func TestPygmentsDivWrapper_Fixtures(t *testing.T) {
	dirs := []string{
		"sphinx-pygments-plain",
		"sphinx-pygments-linenos",
	}
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) { assertRegion(t, pygmentsDivWrapper{}, dir) })
	}
}
