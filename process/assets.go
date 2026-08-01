package process

import (
	"bytes"
	"path/filepath"
	"strings"
)

// assetInfo caches the engine's stylesheet and script once per run so every
// page references byte identical content under one hash.
type assetInfo struct {
	cssName    string
	jsName     string
	cssContent []byte
	jsContent  []byte
	cssHash    string
	jsHash     string
}

func (p *Processor) buildAssets() assetInfo {
	a := p.cfg.Engine.Assets()
	info := assetInfo{
		cssContent: []byte(a.CSS.Content),
		jsContent:  []byte(a.JS.Content),
		cssHash:    a.CSS.Hash,
		jsHash:     a.JS.Hash,
		cssName:    "kazari.css",
		jsName:     "kazari.js",
	}
	if p.cfg.HashedAssets {
		info.cssName = a.CSS.Filename
		info.jsName = a.JS.Filename
	}
	return info
}

// writeAssets emits both asset files at the output root, write if different.
// It runs before the concurrent per file pass so no page ever references
// content that is not yet on disk. Check mode reports actions without
// writing.
func (p *Processor) writeAssets(root string) []AssetResult {
	entries := []struct {
		name    string
		content []byte
	}{
		{p.assets.cssName, p.assets.cssContent},
		{p.assets.jsName, p.assets.jsContent},
	}
	var out []AssetResult
	for _, e := range entries {
		path := filepath.Join(root, e.name)
		action := "created"
		if existing, err := p.fs.ReadFile(path); err == nil {
			if bytes.Equal(existing, e.content) {
				action = "unchanged"
			} else {
				action = "updated"
			}
		}
		if action != "unchanged" && !p.cfg.Check {
			if err := p.fs.WriteFile(path, e.content); err != nil {
				p.logf("kazari process: writing %s: %v", path, err)
			}
		}
		out = append(out, AssetResult{Path: path, Action: action})
	}
	return out
}

// assetHref builds the URL a page uses to reference an asset. The default
// is a relative path computed from the page's depth below the output root,
// built with forward slashes explicitly because filepath.Rel would emit
// backslashes on Windows and break the href. Relative paths stay correct
// under subpath deployments like GitHub Pages project sites. AssetsBase
// overrides with a verbatim prefix. Plain names carry a content hash query
// for cache busting; hashed filenames already embed it.
func (p *Processor) assetHref(relPath, name, hash string) string {
	var href string
	if p.cfg.AssetsBase != "" {
		href = strings.TrimSuffix(p.cfg.AssetsBase, "/") + "/" + name
	} else {
		depth := strings.Count(filepath.ToSlash(relPath), "/")
		href = strings.Repeat("../", depth) + name
	}
	if p.cfg.HashedAssets {
		return href
	}
	return href + "?v=" + hash
}

func (p *Processor) linkTag(href string) string {
	return `<link rel="stylesheet" href="` + href + `" data-kazari="assets">`
}

func (p *Processor) scriptTag(href string) string {
	return `<script src="` + href + `" data-kazari="assets"></script>`
}
