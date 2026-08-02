package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frostybee/kazari"
	kazarinuri "github.com/frostybee/kazari/nuri"
	"github.com/frostybee/kazari/process"
	"github.com/frostybee/nuri"
	"github.com/frostybee/nuri/bundle/core"
)

const (
	defaultLightTheme = "github-light"
	defaultDarkTheme  = "github-dark"
)

func runProcess(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("process", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: kazari process [dir] [flags]\n\nFlags:\n")
		fs.PrintDefaults()
	}
	check := fs.Bool("check", false, "report would-be changes without writing; exit 1 if any")
	configPath := fs.String("config", "", "config file path (default: auto-discover kazari.config.yaml|.yml|.json in dir, then the working directory)")
	themeLight := fs.String("theme-light", "", "light syntax theme name (overrides config)")
	themeDark := fs.String("theme-dark", "", "dark syntax theme name (overrides config)")
	assetsBase := fs.String("assets-base", "", "fixed asset URL prefix instead of per-file relative paths")
	hashedAssets := fs.Bool("hashed-assets", false, "use content hashed asset filenames")
	skipUnlabeled := fs.Bool("skip-unlabeled", false, "leave blocks without a detectable language untouched")
	concurrency := fs.Int("concurrency", 0, "max files processed concurrently (default: number of CPUs)")
	verbose := fs.Bool("verbose", false, "log per-file progress to stderr")

	// The stdlib flag package stops at the first positional argument, so
	// re-parse the remainder to accept flags on either side of the dir.
	dir := ""
	rest := args
	for {
		if err := fs.Parse(rest); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return 0
			}
			return 2
		}
		if fs.NArg() == 0 {
			break
		}
		if dir != "" {
			fmt.Fprintf(stderr, "kazari: too many arguments: %s\n", strings.Join(fs.Args(), " "))
			return 2
		}
		dir = fs.Arg(0)
		rest = fs.Args()[1:]
	}
	if dir == "" {
		dir = "."
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "kazari: %q is not a directory\n", dir)
		return 2
	}

	fc, fcPath, err := loadFileConfig(*configPath, dir)
	if err != nil {
		fmt.Fprintf(stderr, "kazari: %v\n", err)
		return 2
	}
	if fc != nil && *verbose {
		fmt.Fprintf(stderr, "kazari: using config %s\n", fcPath)
	}

	light, dark := defaultLightTheme, defaultDarkTheme
	if fc != nil && fc.Themes != nil {
		if fc.Themes.Light != "" {
			light = fc.Themes.Light
		}
		if fc.Themes.Dark != "" {
			dark = fc.Themes.Dark
		}
	}
	if *themeLight != "" {
		light = *themeLight
	}
	if *themeDark != "" {
		dark = *themeDark
	}
	if err := validateThemeNames(light, dark); err != nil {
		fmt.Fprintf(stderr, "kazari: %v\n", err)
		return 2
	}

	ctx := context.Background()
	hl, err := nuri.New(ctx, nuri.WithFS(core.FS()))
	if err != nil {
		fmt.Fprintf(stderr, "kazari: initializing highlighter: %v\n", err)
		return 2
	}
	var opts []kazari.Option
	if fc != nil {
		fcOpts, oerr := kazari.FileConfigToOptions(fc)
		if oerr != nil {
			fmt.Fprintf(stderr, "kazari: %v\n", oerr)
			return 2
		}
		opts = append(opts, fcOpts...)
	}
	opts = append(opts,
		kazari.WithHighlighter(kazarinuri.New(ctx, hl)),
		kazari.WithThemes(light, dark),
		kazari.WithWarningHandler(func(msg string) { fmt.Fprintln(stderr, msg) }),
	)
	engine := kazari.New(opts...)

	pcfg := process.Config{Engine: engine, Check: *check}
	if fc != nil && fc.Process != nil {
		pc := fc.Process
		if pc.SkipUnlabeled != nil {
			pcfg.SkipUnlabeled = *pc.SkipUnlabeled
		}
		if pc.AssetsBase != nil {
			pcfg.AssetsBase = *pc.AssetsBase
		}
		if pc.HashedAssets != nil {
			pcfg.HashedAssets = *pc.HashedAssets
		}
		if pc.Concurrency != nil {
			pcfg.Concurrency = *pc.Concurrency
		}
		if pc.MaxFileBytes != nil {
			pcfg.MaxFileBytes = *pc.MaxFileBytes
		}
	}
	explicit := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { explicit[f.Name] = true })
	if explicit["skip-unlabeled"] {
		pcfg.SkipUnlabeled = *skipUnlabeled
	}
	if explicit["assets-base"] {
		pcfg.AssetsBase = *assetsBase
	}
	if explicit["hashed-assets"] {
		pcfg.HashedAssets = *hashedAssets
	}
	if explicit["concurrency"] {
		pcfg.Concurrency = *concurrency
	}
	if *verbose {
		pcfg.Logger = func(format string, logArgs ...any) {
			fmt.Fprintf(stderr, format+"\n", logArgs...)
		}
	}

	proc, err := process.New(pcfg)
	if err != nil {
		fmt.Fprintf(stderr, "kazari: %v\n", err)
		return 2
	}
	result, err := proc.Run(ctx, dir)
	if err != nil {
		fmt.Fprintf(stderr, "kazari: %v\n", err)
		return 2
	}

	failed := false
	for _, f := range result.Files {
		if f.Err != nil {
			failed = true
			fmt.Fprintf(stderr, "kazari: %s: %v\n", f.Path, f.Err)
		}
	}

	if *check {
		for _, a := range result.Assets {
			if a.Action != "unchanged" {
				fmt.Fprintln(stdout, a.Path)
			}
		}
		for _, f := range result.Files {
			if f.Changed {
				fmt.Fprintln(stdout, f.Path)
			}
		}
	}

	found, rewritten, skipped := 0, 0, 0
	for _, f := range result.Files {
		found += f.BlocksFound
		rewritten += f.BlocksRewritten
		skipped += len(f.BlocksSkipped)
	}
	fmt.Fprintf(stdout, "%d files, %d blocks upgraded, %d skipped, %d suppressed, %d changed\n",
		len(result.Files), rewritten, skipped, result.Suppressed, result.ChangedCount)

	switch {
	case failed:
		return 2
	case *check && result.ChangedCount > 0:
		return 1
	default:
		return 0
	}
}

// loadFileConfig resolves and parses the config file. An explicit path must
// exist and parse; auto discovery probes the target dir then the working
// directory and finding nothing is fine. Parse errors are always hard
// errors here, unlike WithConfigDir, which demotes them to warnings and is
// therefore wrong for a CLI.
func loadFileConfig(explicitPath, dir string) (*kazari.FileConfig, string, error) {
	if explicitPath != "" {
		fc, err := parseConfigFile(explicitPath)
		return fc, explicitPath, err
	}
	for _, d := range []string{dir, "."} {
		for _, name := range []string{"kazari.config.yaml", "kazari.config.yml", "kazari.config.json"} {
			path := filepath.Join(d, name)
			if _, err := os.Stat(path); err != nil {
				continue
			}
			fc, err := parseConfigFile(path)
			return fc, path, err
		}
	}
	return nil, "", nil
}

func parseConfigFile(path string) (*kazari.FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var format string
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		format = "yaml"
	case ".json":
		format = "json"
	default:
		return nil, fmt.Errorf("config %s: unsupported extension, use .yaml, .yml, or .json", path)
	}
	fc, err := kazari.ParseConfig(data, format)
	if err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	return fc, nil
}
