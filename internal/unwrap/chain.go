package unwrap

// DivWrapperChain is tried for candidate regions rooted at a wrapper element
// such as div.highlight, div.highlighter-rouge, div.highlight-x, or
// figure.highlight. Order matters and first match wins. The Rouge table
// shape runs before the other Rouge shapes because the table nests inside
// them; Chroma runs first and safely so, since its root signature never
// overlaps the Rouge or Pygments roots.
var DivWrapperChain = []Unwrapper{
	chromaDivWrapper{},
	rougeLinenosTable{},
	rougeDivWrapper{},
	rougeHighlightFigure{},
	pygmentsDivWrapper{},
	zolaDivWrapper{},
}

// BarePreChain is tried for candidate regions rooted at a pre element with
// no recognized wrapper. The generic fallback runs last so that every more
// specific shape gets first refusal.
var BarePreChain = []Unwrapper{
	chromaBarePre{},
	prismBarePre{},
	shikiAstroBarePre{},
	zolaBarePre{},
	plainFallback{},
}

// RunChain tries each unwrapper in order and returns the first match along
// with the winning unwrapper's name. When the region carries a data-kz-meta
// attribute, its value replaces the synthesized meta string verbatim; source
// recovery and language detection are unaffected.
func RunChain(chain []Unwrapper, tokens []BufferedToken) (Region, string, bool) {
	for _, u := range chain {
		region, ok := u.Match(tokens)
		if !ok {
			continue
		}
		if meta, found := kzMetaAttr(tokens); found {
			region.Meta = meta
		}
		return region, u.Name(), true
	}
	return Region{}, "", false
}
