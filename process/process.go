// Package process post-processes a static site generator's built HTML
// output. It finds code blocks in each page, recovers their source through
// the unwrap layer, re-renders them with a shared kazari Engine, splices the
// results back while preserving every other byte of the page, and emits the
// engine's CSS and JS assets once with injected link and script tags. It is
// the library behind the kazari process command and is importable by Go
// programs directly.
package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/frostybee/kazari"
	"github.com/frostybee/kazari/internal/unwrap"
	"golang.org/x/sync/errgroup"
)

const defaultMaxFileBytes = 32 << 20

// Config configures a Processor. Engine is required and shared across
// workers; construct it with kazari.WithWarningHandler wired to the same
// sink as Logger so plaintext fallbacks for unknown languages stay visible.
type Config struct {
	Engine *kazari.Engine

	// Check reports what would change without writing anything.
	Check bool

	// SkipUnlabeled leaves blocks with no detectable language untouched.
	// The default renders them as plain text frames.
	SkipUnlabeled bool

	// AssetsBase, when set, prefixes asset URLs verbatim instead of the
	// default per file relative path.
	AssetsBase string

	// HashedAssets switches from kazari.css and kazari.js to content
	// hashed filenames. Old hashed files from previous runs are not
	// cleaned up.
	HashedAssets bool

	// Concurrency bounds the worker pool; zero means GOMAXPROCS.
	Concurrency int

	// MaxFileBytes skips files larger than this; zero means 32 MiB.
	// Token dense highlighted pages inflate far beyond their byte size
	// in memory, so the guard protects CI runners from pathological
	// inputs.
	MaxFileBytes int64

	// Logger receives progress and warnings; nil is silent.
	Logger func(format string, args ...any)

	// FS abstracts the disk; nil uses the real filesystem with atomic
	// writes.
	FS FileSystem
}

// AssetResult records what happened to one emitted asset file. Action is
// created, updated, or unchanged; Check mode reports the same actions
// without writing.
type AssetResult struct {
	Path   string
	Action string
}

// FileResult records the outcome for one HTML file.
type FileResult struct {
	Path            string
	BlocksFound     int
	BlocksRewritten int
	BlocksSkipped   []string
	Suppressed      int
	Changed         bool
	Err             error
}

// Result is the outcome of one Run.
type Result struct {
	Files        []FileResult
	Assets       []AssetResult
	Suppressed   int
	ChangedCount int
}

// Processor rewrites code blocks across an output directory.
type Processor struct {
	cfg    Config
	fs     FileSystem
	assets assetInfo
}

// New validates the configuration and returns a Processor.
func New(cfg Config) (*Processor, error) {
	if cfg.Engine == nil {
		return nil, errors.New("process: Config.Engine is required")
	}
	if cfg.FS == nil {
		cfg.FS = osFileSystem{}
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = runtime.GOMAXPROCS(0)
	}
	if cfg.MaxFileBytes <= 0 {
		cfg.MaxFileBytes = defaultMaxFileBytes
	}
	return &Processor{cfg: cfg, fs: cfg.FS}, nil
}

func (p *Processor) logf(format string, args ...any) {
	if p.cfg.Logger != nil {
		p.cfg.Logger(format, args...)
	}
}

// Run processes every HTML file under root. Assets are written first so no
// page ever references content that is not yet on disk. Files are processed
// concurrently, each goroutine writing only its own result slot; paths are
// sorted so output ordering is deterministic. Cancellation returns the
// partial Result together with the context error.
func (p *Processor) Run(ctx context.Context, root string) (Result, error) {
	p.assets = p.buildAssets()

	var result Result
	result.Assets = p.writeAssets(root)
	for _, a := range result.Assets {
		if a.Action != "unchanged" {
			result.ChangedCount++
		}
	}

	var paths []string
	err := p.fs.WalkFiles(root, func(path string) error {
		if p.shouldProcess(path) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return result, fmt.Errorf("process: walking %s: %w", root, err)
	}
	sort.Strings(paths)

	results := make([]FileResult, len(paths))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(p.cfg.Concurrency)
	for i, path := range paths {
		g.Go(func() error {
			select {
			case <-gctx.Done():
				results[i] = FileResult{Path: path, Err: gctx.Err()}
				return gctx.Err()
			default:
			}
			results[i] = p.processFile(root, path)
			return nil
		})
	}
	werr := g.Wait()

	result.Files = results
	for i := range results {
		if results[i].Changed {
			result.ChangedCount++
		}
		result.Suppressed += results[i].Suppressed
	}
	return result, werr
}

// shouldProcess filters the walk to HTML pages and keeps the emitted asset
// files out of the pipeline.
func (p *Processor) shouldProcess(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".html" && ext != ".htm" {
		return false
	}
	base := filepath.Base(path)
	return base != p.assets.cssName && base != p.assets.jsName
}

func chainFor(kind unwrap.CandidateKind) []unwrap.Unwrapper {
	if kind == unwrap.CandidateWrapper {
		return unwrap.DivWrapperChain
	}
	return unwrap.BarePreChain
}

func (p *Processor) processFile(root, path string) FileResult {
	fr := FileResult{Path: path}

	src, err := p.fs.ReadFile(path)
	if err != nil {
		fr.Err = err
		return fr
	}
	if int64(len(src)) > p.cfg.MaxFileBytes {
		fr.Err = fmt.Errorf("process: %s is %d bytes, over the %d byte limit", path, len(src), p.cfg.MaxFileBytes)
		p.logf("kazari process: skipping oversized file %s", path)
		return fr
	}
	tokens, err := unwrap.TokenizeFragment(src)
	if err != nil {
		fr.Err = fmt.Errorf("process: tokenizing %s: %w", path, err)
		p.logf("kazari process: %s left untouched: %v", path, err)
		return fr
	}

	page := unwrap.Discover(src, tokens)
	fr.Suppressed = page.Suppressed
	if page.SuppressedUnclosed {
		p.logf("kazari process: %s: an unclosed kazari block scope suppressed the rest of the file", path)
	}

	var edits []edit
	var handleCandidate func(c unwrap.Candidate, allowReoffer bool)
	handleCandidate = func(c unwrap.Candidate, allowReoffer bool) {
		fr.BlocksFound++
		if c.Malformed {
			fr.BlocksSkipped = append(fr.BlocksSkipped, "malformed")
			p.logf("kazari process: %s: unclosed element at byte %d left untouched", path, c.ByteStart)
			return
		}
		if reason := unwrap.SkipReason(c.Tokens); reason != "" {
			fr.BlocksSkipped = append(fr.BlocksSkipped, reason)
			return
		}
		region, _, ok := unwrap.RunChain(chainFor(c.Kind), c.Tokens)
		if !ok {
			// A site component using a trigger class must not swallow a
			// genuine code block nested inside it: re-offer the contents
			// once, one level deep only.
			if c.Kind == unwrap.CandidateWrapper && allowReoffer && len(c.Tokens) > 2 {
				inner := unwrap.DiscoverWithin(c.Tokens[1 : len(c.Tokens)-1])
				fr.Suppressed += inner.Suppressed
				if len(inner.Candidates) > 0 {
					for _, ic := range inner.Candidates {
						handleCandidate(ic, false)
					}
					return
				}
			}
			fr.BlocksSkipped = append(fr.BlocksSkipped, "unrecognized-shape")
			return
		}
		if region.Lang == "" && p.cfg.SkipUnlabeled {
			fr.BlocksSkipped = append(fr.BlocksSkipped, "unlabeled")
			return
		}
		rendered, rerr := p.cfg.Engine.RenderWithMeta(region.Code, region.Meta)
		if rerr != nil {
			fr.BlocksSkipped = append(fr.BlocksSkipped, "render-error")
			p.logf("kazari process: %s: render error, block left untouched: %v", path, rerr)
			return
		}
		edits = append(edits, edit{start: c.ByteStart, end: c.ByteEnd, replacement: []byte(rendered)})
		fr.BlocksRewritten++
	}
	for _, c := range page.Candidates {
		handleCandidate(c, true)
	}

	rel, rerr := filepath.Rel(root, path)
	if rerr != nil {
		rel = filepath.Base(path)
	}
	cssHref := p.assetHref(rel, p.assets.cssName, p.assets.cssHash)
	jsHref := p.assetHref(rel, p.assets.jsName, p.assets.jsHash)
	linkTag, scriptTag := p.linkTag(cssHref), p.scriptTag(jsHref)

	// Existing tags are always reconciled to canonical form, whether or not
	// any block in this file was rewritten; a config only change must still
	// refresh stale cache busting hashes. Only injecting a brand new tag is
	// gated on a rewritten block.
	haveLink, haveScript := false, false
	for _, tag := range page.AssetTags {
		var canonical string
		switch tag.Kind {
		case "link":
			canonical = linkTag
			haveLink = true
		case "script":
			canonical = scriptTag
			haveScript = true
		default:
			continue
		}
		if !bytes.Equal(src[tag.ByteStart:tag.ByteEnd], []byte(canonical)) {
			edits = append(edits, edit{start: tag.ByteStart, end: tag.ByteEnd, replacement: []byte(canonical)})
		}
	}
	if fr.BlocksRewritten > 0 {
		if !haveLink {
			edits = append(edits, edit{start: page.LinkInsert, end: page.LinkInsert, replacement: []byte(linkTag)})
		}
		if !haveScript {
			edits = append(edits, edit{start: page.ScriptInsert, end: page.ScriptInsert, replacement: []byte(scriptTag)})
		}
	}

	if len(edits) == 0 {
		return fr
	}
	out, aerr := applyEdits(src, edits)
	if aerr != nil {
		fr.Err = aerr
		p.logf("kazari process: %s left untouched: %v", path, aerr)
		return fr
	}
	if bytes.Equal(out, src) {
		return fr
	}
	fr.Changed = true
	if p.cfg.Check {
		return fr
	}
	if werr := p.fs.WriteFile(path, out); werr != nil {
		fr.Err = werr
	}
	return fr
}
