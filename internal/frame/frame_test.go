package frame

import (
	"testing"

	"github.com/frostybee/kazari/internal/config"
)

func TestStripTerminalComments(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "empty input",
			code: "",
			want: "",
		},
		{
			name: "no comments",
			code: "npm install\nnpm start",
			want: "npm install\nnpm start",
		},
		{
			name: "all comments",
			code: "# step 1\n# step 2",
			want: "",
		},
		{
			name: "mixed comments and commands",
			code: "# Install deps\nnpm install\n# Start server\nnpm start",
			want: "npm install\nnpm start",
		},
		{
			name: "indented comments stripped",
			code: "  # indented comment\nnpm install",
			want: "npm install",
		},
		{
			name: "inline hash preserved",
			code: "echo \"hello\" # not a comment line",
			want: "echo \"hello\" # not a comment line",
		},
		{
			name: "trailing blank lines cleaned",
			code: "npm install\n# trailing comment\n",
			want: "npm install",
		},
		{
			name: "blank lines between commands preserved",
			code: "npm install\n\nnpm start",
			want: "npm install\n\nnpm start",
		},
		{
			name: "bare hash line stripped",
			code: "#\nnpm install",
			want: "npm install",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripTerminalComments(tt.code)
			if got != tt.want {
				t.Errorf("StripTerminalComments():\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

// --- IsTerminalLanguage ---

func TestIsTerminalLanguage(t *testing.T) {
	terminals := []string{"bash", "sh", "zsh", "powershell", "fish", "console", "ansi", "nu"}
	for _, lang := range terminals {
		if !isTerminalLanguage(lang) {
			t.Errorf("isTerminalLanguage(%q) = false, want true", lang)
		}
	}
}

func TestIsTerminalLanguage_NonTerminal(t *testing.T) {
	nonTerminals := []string{"go", "python", "javascript", "rust", ""}
	for _, lang := range nonTerminals {
		if isTerminalLanguage(lang) {
			t.Errorf("isTerminalLanguage(%q) = true, want false", lang)
		}
	}
}

func TestIsTerminalLanguage_CaseInsensitive(t *testing.T) {
	for _, lang := range []string{"BASH", "Bash", "BaSh"} {
		if !isTerminalLanguage(lang) {
			t.Errorf("isTerminalLanguage(%q) = false, want true (case-insensitive)", lang)
		}
	}
}

// --- DetectFrameType ---

func TestDetectFrameType_ExplicitDefault(t *testing.T) {
	got := DetectFrameType("echo hello", "bash", config.FrameCode)
	if got != config.FrameCode {
		t.Errorf("explicit FrameCode should be returned as-is, got %d", got)
	}

	got = DetectFrameType("echo hello", "bash", config.FrameNone)
	if got != config.FrameNone {
		t.Errorf("explicit FrameNone should be returned as-is, got %d", got)
	}
}

func TestDetectFrameType_NonTerminalLanguage(t *testing.T) {
	got := DetectFrameType("func main() {}", "go", config.FrameAuto)
	if got != config.FrameCode {
		t.Errorf("non-terminal language should be FrameCode, got %d", got)
	}
}

func TestDetectFrameType_TerminalNoIndicator(t *testing.T) {
	got := DetectFrameType("npm install\nnpm start", "bash", config.FrameAuto)
	if got != config.FrameTerminal {
		t.Errorf("terminal lang without file indicator should be FrameTerminal, got %d", got)
	}
}

func TestDetectFrameType_TerminalWithShebang(t *testing.T) {
	got := DetectFrameType("#!/bin/bash\necho hello", "bash", config.FrameAuto)
	if got != config.FrameCode {
		t.Errorf("terminal lang with shebang should be FrameCode (script), got %d", got)
	}
}

func TestDetectFrameType_TerminalWithFileComment(t *testing.T) {
	got := DetectFrameType("# deploy.sh\necho hello", "bash", config.FrameAuto)
	if got != config.FrameCode {
		t.Errorf("terminal lang with file comment should be FrameCode (script), got %d", got)
	}
}

// --- ExtractFileName ---

func TestExtractFileName_SlashSlash(t *testing.T) {
	title, modified := ExtractFileName("// main.go\npackage main", "go")
	if title != "main.go" {
		t.Errorf("title = %q, want %q", title, "main.go")
	}
	if modified == "// main.go\npackage main" {
		t.Error("modified code should have comment removed")
	}
}

func TestExtractFileName_Hash(t *testing.T) {
	title, _ := ExtractFileName("# deploy.sh\necho hello", "bash")
	if title != "deploy.sh" {
		t.Errorf("title = %q, want %q", title, "deploy.sh")
	}
}

func TestExtractFileName_HTML(t *testing.T) {
	title, _ := ExtractFileName("<!-- index.html -->\n<div>hello</div>", "html")
	if title != "index.html" {
		t.Errorf("title = %q, want %q", title, "index.html")
	}
}

func TestExtractFileName_BlockComment(t *testing.T) {
	title, _ := ExtractFileName("/* styles.css */\nbody {}", "css")
	if title != "styles.css" {
		t.Errorf("title = %q, want %q", title, "styles.css")
	}
}

func TestExtractFileName_WithPrefix(t *testing.T) {
	title, _ := ExtractFileName("// file name: main.go\npackage main", "go")
	if title != "main.go" {
		t.Errorf("title = %q, want %q", title, "main.go")
	}
}

func TestExtractFileName_NotFound(t *testing.T) {
	title, code := ExtractFileName("package main\nfunc main() {}", "go")
	if title != "" {
		t.Errorf("title = %q, want empty", title)
	}
	if code != "package main\nfunc main() {}" {
		t.Error("code should be unchanged when no file name found")
	}
}

func TestExtractFileName_BeyondLine4(t *testing.T) {
	code := "line1\nline2\nline3\nline4\n// main.go\ncode"
	title, _ := ExtractFileName(code, "go")
	if title != "" {
		t.Errorf("file name beyond line 4 should not be found, got %q", title)
	}
}

func TestExtractFileName_BlankLineAfterRemoval(t *testing.T) {
	_, modified := ExtractFileName("// main.go\n\npackage main", "go")
	if modified != "package main" {
		t.Errorf("blank line after file comment should be skipped, got %q", modified)
	}
}

// --- isFilePath (tested via extractFromComment) ---

func TestExtractFromComment_KnownExtensions(t *testing.T) {
	for _, path := range []string{"main.go", "app.js", "script.py"} {
		got := extractFromComment("// " + path)
		if got != path {
			t.Errorf("extractFromComment(%q) = %q, want %q", "// "+path, got, path)
		}
	}
}

func TestExtractFromComment_Dotfiles(t *testing.T) {
	for _, path := range []string{".env", ".gitignore"} {
		got := extractFromComment("// " + path)
		if got != path {
			t.Errorf("extractFromComment(%q) = %q, want %q", "// "+path, got, path)
		}
	}
}

func TestExtractFromComment_SpecialNames(t *testing.T) {
	for _, name := range []string{"Makefile", "Dockerfile"} {
		got := extractFromComment("// " + name)
		if got != name {
			t.Errorf("extractFromComment(%q) = %q, want %q", "// "+name, got, name)
		}
	}
}

func TestExtractFromComment_URLRejected(t *testing.T) {
	got := extractFromComment("// https://example.com")
	if got != "" {
		t.Errorf("URLs should be rejected, got %q", got)
	}
}

func TestExtractFromComment_SpacesRejected(t *testing.T) {
	got := extractFromComment("// this is a sentence")
	if got != "" {
		t.Errorf("strings with spaces should be rejected, got %q", got)
	}
}

func TestExtractFromComment_UnknownExtension(t *testing.T) {
	got := extractFromComment("// file.xyz")
	if got != "" {
		t.Errorf("unknown extension should be rejected, got %q", got)
	}
}

func TestExtractFromComment_ShebangExcluded(t *testing.T) {
	got := extractFromComment("#!/bin/bash")
	if got != "" {
		t.Errorf("shebang should be excluded from comment extraction, got %q", got)
	}
}

func TestExtractFromComment_EmptyContent(t *testing.T) {
	got := extractFromComment("// ")
	if got != "" {
		t.Errorf("empty content after comment should return empty, got %q", got)
	}
}

// --- Frontmatter delimiter cleanup ---

func TestExtractFileName_FrontmatterCleanup(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantTitle string
		wantCode  string
	}{
		{
			name:      "YAML frontmatter emptied",
			code:      "---\n# filename: doc.md\n---\n\ncontent",
			wantTitle: "doc.md",
			wantCode:  "content",
		},
		{
			name:      "TOML frontmatter emptied",
			code:      "+++\n# filename: config.toml\n+++\ncontent",
			wantTitle: "config.toml",
			wantCode:  "content",
		},
		{
			name:      "non-empty frontmatter untouched",
			code:      "---\ntitle: hello\n# filename: doc.md\n---\ncontent",
			wantTitle: "doc.md",
			wantCode:  "---\ntitle: hello\n---\ncontent",
		},
		{
			name:      "no frontmatter untouched",
			code:      "// main.go\npackage main",
			wantTitle: "main.go",
			wantCode:  "package main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, modified := ExtractFileName(tt.code, "md")
			if title != tt.wantTitle {
				t.Errorf("title: got %q, want %q", title, tt.wantTitle)
			}
			if modified != tt.wantCode {
				t.Errorf("code:\ngot:  %q\nwant: %q", modified, tt.wantCode)
			}
		})
	}
}

func TestRemoveEmptyFrontmatter_MismatchedDelimiters(t *testing.T) {
	code := "---\n+++\ncontent"
	if got := removeEmptyFrontmatter(code); got != code {
		t.Errorf("mismatched delimiters should be untouched, got %q", got)
	}
}
