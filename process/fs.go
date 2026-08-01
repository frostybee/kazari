package process

import (
	"io/fs"
	"os"
	"path/filepath"
)

// FileSystem is the minimal surface the processor needs. The default
// implementation is the real disk; tests supply an in-memory fake. WalkFiles
// visits regular files only and never follows symlinks. Its callback error
// is reserved for fatal enumeration failures; per file processing failures
// are captured in FileResult.Err instead and never abort the walk.
type FileSystem interface {
	WalkFiles(root string, fn func(path string) error) error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

type osFileSystem struct{}

func (osFileSystem) WalkFiles(root string, fn func(string) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d == nil || !d.Type().IsRegular() {
			return nil
		}
		return fn(path)
	})
}

func (osFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes to a temp file in the target directory and renames it
// over the destination. A dev server serving the tree mid run never observes
// a truncated file; os.Rename replaces existing files on Windows too.
func (osFileSystem) WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".kazari-tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	_, werr := tmp.Write(data)
	cerr := tmp.Close()
	if werr != nil || cerr != nil {
		os.Remove(tmpName)
		if werr != nil {
			return werr
		}
		return cerr
	}
	mode := os.FileMode(0o644)
	if info, serr := os.Stat(path); serr == nil {
		mode = info.Mode().Perm()
	}
	if cherr := os.Chmod(tmpName, mode); cherr != nil {
		os.Remove(tmpName)
		return cherr
	}
	if rerr := os.Rename(tmpName, path); rerr != nil {
		os.Remove(tmpName)
		return rerr
	}
	return nil
}
