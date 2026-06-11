package locale

import "fmt"

type UIStrings struct {
	CopyLabel             string
	CopyTitle             string
	CopySuccess           string
	FullscreenLabel       string
	ExpandButtonText      string
	CollapseButtonText    string
	ExpandedAnnouncement  string
	CollapsedAnnouncement string
	CollapsedLineSingular string
	CollapsedLinePlural   string
	CodeGroupFallback     string
}

var locales = map[string]UIStrings{
	"en-US": {
		CopyLabel:             "Copy",
		CopyTitle:             "Copy to clipboard",
		CopySuccess:           "Copied!",
		FullscreenLabel:       "Fullscreen",
		ExpandButtonText:      "Show more",
		CollapseButtonText:    "Show less",
		ExpandedAnnouncement:  "Code block expanded",
		CollapsedAnnouncement: "Code block collapsed",
		CollapsedLineSingular: "1 collapsed line",
		CollapsedLinePlural:   "%d collapsed lines",
		CodeGroupFallback:     "Code",
	},
	"fr-FR": {
		CopyLabel:             "Copier",
		CopyTitle:             "Copier dans le presse-papiers",
		CopySuccess:           "Copié !",
		FullscreenLabel:       "Plein écran",
		ExpandButtonText:      "Afficher plus",
		CollapseButtonText:    "Afficher moins",
		ExpandedAnnouncement:  "Bloc de code déplié",
		CollapsedAnnouncement: "Bloc de code réduit",
		CollapsedLineSingular: "1 ligne masquée",
		CollapsedLinePlural:   "%d lignes masquées",
		CodeGroupFallback:     "Code",
	},
	"ja-JP": {
		CopyLabel:             "コピー",
		CopyTitle:             "クリップボードにコピー",
		CopySuccess:           "コピーしました",
		FullscreenLabel:       "全画面",
		ExpandButtonText:      "もっと見る",
		CollapseButtonText:    "閉じる",
		ExpandedAnnouncement:  "コードブロックを展開しました",
		CollapsedAnnouncement: "コードブロックを折りたたみました",
		CollapsedLineSingular: "1 行を折りたたみ",
		CollapsedLinePlural:   "%d 行を折りたたみ",
		CodeGroupFallback:     "コード",
	},
}

// Resolve returns UI strings for the given locale with optional key overrides.
// Falls back to en-US for unknown locales.
func Resolve(loc string, overrides map[string]string) *UIStrings {
	base, ok := locales[loc]
	if !ok {
		base = locales["en-US"]
	}
	s := base
	for key, val := range overrides {
		applyOverride(&s, key, val)
	}
	return &s
}

func applyOverride(s *UIStrings, key, val string) {
	switch key {
	case "copy.label":
		s.CopyLabel = val
	case "copy.title":
		s.CopyTitle = val
	case "copy.success":
		s.CopySuccess = val
	case "fullscreen.label":
		s.FullscreenLabel = val
	case "collapse.expand":
		s.ExpandButtonText = val
	case "collapse.collapse":
		s.CollapseButtonText = val
	case "collapse.expanded":
		s.ExpandedAnnouncement = val
	case "collapse.collapsed":
		s.CollapsedAnnouncement = val
	case "collapse.summary.singular":
		s.CollapsedLineSingular = val
	case "collapse.summary.plural":
		s.CollapsedLinePlural = val
	case "codegroup.fallback":
		s.CodeGroupFallback = val
	}
}

// FormatCollapsedLines returns the summary text for a collapsed section.
func FormatCollapsedLines(s *UIStrings, count int) string {
	if count == 1 {
		return s.CollapsedLineSingular
	}
	return fmt.Sprintf(s.CollapsedLinePlural, count)
}
