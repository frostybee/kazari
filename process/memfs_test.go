package process

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// memFS is an in-memory FileSystem for tests that must observe writes (or
// their absence, for the Check mode no-write invariant).
type memFS struct {
	mu     sync.Mutex
	files  map[string][]byte
	writes []string
}

func newMemFS() *memFS {
	return &memFS{files: map[string][]byte{}}
}

func (m *memFS) put(path string, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[path] = data
}

func (m *memFS) WalkFiles(root string, fn func(string) error) error {
	m.mu.Lock()
	prefix := strings.TrimSuffix(root, string(filepath.Separator)) + string(filepath.Separator)
	var paths []string
	for p := range m.files {
		if strings.HasPrefix(p, prefix) {
			paths = append(paths, p)
		}
	}
	m.mu.Unlock()
	sort.Strings(paths)
	for _, p := range paths {
		if err := fn(p); err != nil {
			return err
		}
	}
	return nil
}

func (m *memFS) ReadFile(path string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, ok := m.files[path]
	if !ok {
		return nil, fmt.Errorf("memfs: %s: %w", path, os.ErrNotExist)
	}
	return data, nil
}

func (m *memFS) WriteFile(path string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[path] = data
	m.writes = append(m.writes, path)
	return nil
}

func (m *memFS) writeCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.writes)
}
