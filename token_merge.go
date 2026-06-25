package kazari

import "github.com/frostybee/kazari/internal/render"

// mergeTokens pairs light and dark tokens into MergedToken lines.
func mergeTokens(light, dark [][]Token) []render.TokenLine {
	lines := make([]render.TokenLine, len(light))

	for i, lightLine := range light {
		var darkLine []Token
		if dark != nil && i < len(dark) {
			darkLine = dark[i]
		}

		if darkLine == nil || len(lightLine) == len(darkLine) {
			// Fast path: boundaries match (common case).
			tokens := make([]render.MergedToken, len(lightLine))
			for j, lt := range lightLine {
				mt := render.MergedToken{
					Content:    lt.Content,
					LightColor: lt.Color,
					LightBG:    lt.BgColor,
					FontStyle:  lt.FontStyle,
				}
				if darkLine != nil && j < len(darkLine) {
					mt.DarkColor = darkLine[j].Color
					mt.DarkBG = darkLine[j].BgColor
				}
				tokens[j] = mt
			}
			lines[i] = render.TokenLine{Tokens: tokens}
		} else {
			// Slow path: boundaries differ — align by character position.
			lines[i] = render.TokenLine{Tokens: alignTokens(lightLine, darkLine)}
		}
	}

	return lines
}

// alignTokens handles the rare case where light and dark tokens have different boundaries.
func alignTokens(lightLine, darkLine []Token) []render.MergedToken {
	var result []render.MergedToken
	li, di := 0, 0
	lo, do := 0, 0

	for li < len(lightLine) && di < len(darkLine) {
		lt := lightLine[li]
		dt := darkLine[di]
		lRemain := len(lt.Content) - lo
		dRemain := len(dt.Content) - do
		take := lRemain
		if dRemain < take {
			take = dRemain
		}

		result = append(result, render.MergedToken{
			Content:    lt.Content[lo : lo+take],
			LightColor: lt.Color,
			DarkColor:  dt.Color,
			LightBG:    lt.BgColor,
			DarkBG:     dt.BgColor,
			FontStyle:  lt.FontStyle,
		})

		lo += take
		do += take
		if lo >= len(lt.Content) {
			li++
			lo = 0
		}
		if do >= len(dt.Content) {
			di++
			do = 0
		}
	}

	for li < len(lightLine) {
		lt := lightLine[li]
		content := lt.Content[lo:]
		if content != "" {
			result = append(result, render.MergedToken{
				Content:    content,
				LightColor: lt.Color,
				LightBG:    lt.BgColor,
				FontStyle:  lt.FontStyle,
			})
		}
		li++
		lo = 0
	}

	return result
}
