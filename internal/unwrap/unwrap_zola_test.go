package unwrap

import "testing"

func TestZolaBarePre_Fixtures(t *testing.T) {
	assertRegion(t, zolaBarePre{}, "zola")
}

func TestZolaBarePre_RequiresBothSignalHalves(t *testing.T) {
	// data-lang alone must not claim a block; the inline styled pre is the
	// second required half of the signature.
	tokens, err := TokenizeFragment([]byte(`<pre><code data-lang="go">x</code></pre>`))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := (zolaBarePre{}).Match(tokens); ok {
		t.Fatal("data-lang without an inline styled pre must not match")
	}
}
