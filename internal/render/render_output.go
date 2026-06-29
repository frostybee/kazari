package render

import (
	"fmt"
	"html"
	"strings"

	"github.com/frostybee/kazari/internal/config"
)

func renderOutputPanel(sb *strings.Builder, resolved *config.ResolvedBlock, cfg *config.Config) {
	panelClass := "kz-output"
	expanded := "true"
	if resolved.OutputCollapsed {
		panelClass += " kz-output-hidden"
		expanded = "false"
	}
	sb.WriteString(fmt.Sprintf("<div class=\"%s\">", panelClass))

	label := resolved.OutputLabel
	if label == "" {
		label = cfg.UIStrings.OutputLabel
	}
	sb.WriteString("<div class=\"kz-output-header\">")
	sb.WriteString(fmt.Sprintf("<button class=\"kz-output-toggle\" aria-expanded=\"%s\">", expanded))
	sb.WriteString(html.EscapeString(label))
	sb.WriteString("</button>")
	sb.WriteString("</div>")

	sb.WriteString("<pre class=\"kz-output-pre\">")
	sb.WriteString(html.EscapeString(resolved.OutputText))
	sb.WriteString("</pre>")
	sb.WriteString("</div>")
}
