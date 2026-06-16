package locale

import "fmt"

type UIStrings struct {
	CopyLabel             string
	CopyTitle             string
	CopySuccess           string
	FullscreenLabel       string
	FontIncreaseLabel     string
	FontDecreaseLabel     string
	FontResetLabel        string
	FullscreenHint        string
	WrapEnableLabel       string
	WrapDisableLabel      string
	ExpandButtonText      string
	CollapseButtonText    string
	ExpandedAnnouncement  string
	CollapsedAnnouncement string
	CollapsedLineSingular string
	CollapsedLinePlural   string
	CodeGroupFallback        string
	TerminalWindowLabel      string
	ThemeToggleLabel         string
	ThemeToggleAnnouncement  string
}

var locales = map[string]UIStrings{
	"en-US": {
		CopyLabel:             "Copy",
		CopyTitle:             "Copy to clipboard",
		CopySuccess:           "Copied!",
		FullscreenLabel:       "Fullscreen",
		FontIncreaseLabel:     "Increase font size",
		FontDecreaseLabel:     "Decrease font size",
		FontResetLabel:        "Double-click to reset",
		FullscreenHint:        "Press Esc to exit fullscreen",
		WrapEnableLabel:       "Enable word wrap",
		WrapDisableLabel:      "Disable word wrap",
		ExpandButtonText:      "Show more",
		CollapseButtonText:    "Show less",
		ExpandedAnnouncement:  "Code block expanded",
		CollapsedAnnouncement: "Code block collapsed",
		CollapsedLineSingular: "1 collapsed line",
		CollapsedLinePlural:   "%d collapsed lines",
		CodeGroupFallback:        "Code",
		TerminalWindowLabel:      "Terminal window",
		ThemeToggleLabel:         "Toggle theme",
		ThemeToggleAnnouncement:  "Theme toggled",
	},
	"fr-FR": {
		CopyLabel:             "Copier",
		CopyTitle:             "Copier dans le presse-papiers",
		CopySuccess:           "Copié !",
		FullscreenLabel:       "Plein écran",
		FontIncreaseLabel:     "Augmenter la taille",
		FontDecreaseLabel:     "Réduire la taille",
		FontResetLabel:        "Double-cliquez pour réinitialiser",
		FullscreenHint:        "Appuyez sur Échap pour quitter",
		WrapEnableLabel:       "Activer le retour à la ligne",
		WrapDisableLabel:      "Désactiver le retour à la ligne",
		ExpandButtonText:      "Afficher plus",
		CollapseButtonText:    "Afficher moins",
		ExpandedAnnouncement:  "Bloc de code déplié",
		CollapsedAnnouncement: "Bloc de code réduit",
		CollapsedLineSingular: "1 ligne masquée",
		CollapsedLinePlural:   "%d lignes masquées",
		CodeGroupFallback:        "Code",
		TerminalWindowLabel:      "Fenêtre de terminal",
		ThemeToggleLabel:         "Basculer le thème",
		ThemeToggleAnnouncement:  "Thème basculé",
	},
	"ja-JP": {
		CopyLabel:             "コピー",
		CopyTitle:             "クリップボードにコピー",
		CopySuccess:           "コピーしました",
		FullscreenLabel:       "全画面",
		FontIncreaseLabel:     "フォントサイズを拡大",
		FontDecreaseLabel:     "フォントサイズを縮小",
		FontResetLabel:        "ダブルクリックでリセット",
		FullscreenHint:        "Escで全画面を終了",
		WrapEnableLabel:       "折り返しを有効にする",
		WrapDisableLabel:      "折り返しを無効にする",
		ExpandButtonText:      "もっと見る",
		CollapseButtonText:    "閉じる",
		ExpandedAnnouncement:  "コードブロックを展開しました",
		CollapsedAnnouncement: "コードブロックを折りたたみました",
		CollapsedLineSingular: "1 行を折りたたみ",
		CollapsedLinePlural:   "%d 行を折りたたみ",
		CodeGroupFallback:        "コード",
		TerminalWindowLabel:      "ターミナルウィンドウ",
		ThemeToggleLabel:         "テーマ切替",
		ThemeToggleAnnouncement:  "テーマを切り替えました",
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
	case "fullscreen.font.increase":
		s.FontIncreaseLabel = val
	case "fullscreen.font.decrease":
		s.FontDecreaseLabel = val
	case "fullscreen.font.reset":
		s.FontResetLabel = val
	case "fullscreen.hint":
		s.FullscreenHint = val
	case "wrap.enable":
		s.WrapEnableLabel = val
	case "wrap.disable":
		s.WrapDisableLabel = val
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
	case "terminal.label":
		s.TerminalWindowLabel = val
	case "theme.toggle":
		s.ThemeToggleLabel = val
	case "theme.toggle.announcement":
		s.ThemeToggleAnnouncement = val
	}
}

// FormatCollapsedLines returns the summary text for a collapsed section.
func FormatCollapsedLines(s *UIStrings, count int) string {
	if count == 1 {
		return s.CollapsedLineSingular
	}
	return fmt.Sprintf(s.CollapsedLinePlural, count)
}
