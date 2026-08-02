package kazari

import (
	"strings"
	"testing"
)

func TestParseConfigProcessBlockYAML(t *testing.T) {
	fc, err := ParseConfig([]byte(`
process:
  skipUnlabeled: true
  assetsBase: /assets/
  hashedAssets: true
  concurrency: 4
  maxFileBytes: 1048576
`), "yaml")
	if err != nil {
		t.Fatal(err)
	}
	p := fc.Process
	if p == nil {
		t.Fatal("process block not parsed")
	}
	if p.SkipUnlabeled == nil || !*p.SkipUnlabeled {
		t.Fatal("skipUnlabeled")
	}
	if p.AssetsBase == nil || *p.AssetsBase != "/assets/" {
		t.Fatal("assetsBase")
	}
	if p.HashedAssets == nil || !*p.HashedAssets {
		t.Fatal("hashedAssets")
	}
	if p.Concurrency == nil || *p.Concurrency != 4 {
		t.Fatal("concurrency")
	}
	if p.MaxFileBytes == nil || *p.MaxFileBytes != 1048576 {
		t.Fatal("maxFileBytes")
	}
}

func TestParseConfigProcessBlockJSON(t *testing.T) {
	fc, err := ParseConfig([]byte(`{"process":{"skipUnlabeled":true}}`), "json")
	if err != nil {
		t.Fatal(err)
	}
	if fc.Process == nil || fc.Process.SkipUnlabeled == nil || !*fc.Process.SkipUnlabeled {
		t.Fatal("process block not parsed from JSON")
	}
}

func TestParseConfigProcessBlockUnknownKeyStillErrors(t *testing.T) {
	_, err := ParseConfig([]byte("process:\n  selector: \"pre > code\"\n"), "yaml")
	if err == nil {
		t.Fatal("unknown key inside process block must fail strict parsing")
	}
	if !strings.Contains(err.Error(), "selector") {
		t.Fatalf("error should name the unknown key: %v", err)
	}
}

func TestParseConfigProcessValidation(t *testing.T) {
	if _, err := ParseConfig([]byte("process:\n  concurrency: 0\n"), "yaml"); err == nil {
		t.Fatal("concurrency 0 must fail validation")
	}
	if _, err := ParseConfig([]byte("process:\n  maxFileBytes: 0\n"), "yaml"); err == nil {
		t.Fatal("maxFileBytes 0 must fail validation")
	}
}
