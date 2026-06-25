package kazari

import "github.com/frostybee/kazari/internal/config"

func convertLineMarkers(markers []LineMarker) []config.LineMarker {
	if len(markers) == 0 {
		return nil
	}
	out := make([]config.LineMarker, len(markers))
	for i, m := range markers {
		out[i] = config.LineMarker{
			Type:  config.MarkerType(m.Type),
			Lines: convertRanges(m.Lines),
			Label: m.Label,
		}
	}
	return out
}

func convertInlineMarkers(markers []InlineMarker) []config.InlineMarker {
	if len(markers) == 0 {
		return nil
	}
	out := make([]config.InlineMarker, len(markers))
	for i, m := range markers {
		out[i] = config.InlineMarker{
			Type:    config.MarkerType(m.Type),
			Text:    m.Text,
			IsRegex: m.IsRegex,
		}
	}
	return out
}

func convertRanges(ranges []Range) []config.LineRange {
	if len(ranges) == 0 {
		return nil
	}
	out := make([]config.LineRange, len(ranges))
	for i, r := range ranges {
		out[i] = config.LineRange{Start: r.Start, End: r.End}
	}
	return out
}

func convertCollapseStyle(s *CollapseStyle) *config.CollapseStyle {
	if s == nil {
		return nil
	}
	cs := config.CollapseStyle(*s)
	return &cs
}

func mapOptionsToBlockOpts(opts Options) *config.BlockOptions {
	bo := &config.BlockOptions{
		Lang:            opts.Lang,
		Title:           opts.Title,
		Theme:           opts.Theme,
		StartLineNumber: opts.StartLineNumber,
	}
	if opts.Frame != nil {
		f := int(*opts.Frame)
		bo.Frame = &f
	}
	if opts.LineNumbers != nil {
		bo.LineNumbers = opts.LineNumbers
	}
	if opts.Wrap != nil {
		bo.Wrap = opts.Wrap
	}
	if opts.PreserveIndent != nil {
		bo.PreserveIndent = opts.PreserveIndent
	}
	if opts.HangingIndent != nil {
		bo.HangingIndent = opts.HangingIndent
	}
	return bo
}
