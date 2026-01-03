// Package inspector provides project inspection and detection capabilities
// for the AI agent-assisted setup feature.
package inspector

import (
	"os"
	"path/filepath"
	"sync"
)

// FileCache provides caching for file system operations to improve performance.
// It caches file existence checks and file contents for repeated lookups.
type FileCache struct {
	root         string
	existsCache  map[string]bool
	contentCache map[string][]byte
	mu           sync.RWMutex
}

// NewFileCache creates a new FileCache for the given root directory.
func NewFileCache(root string) *FileCache {
	return &FileCache{
		root:         root,
		existsCache:  make(map[string]bool),
		contentCache: make(map[string][]byte),
	}
}

// FileExists checks if a file exists, using cache when available.
func (c *FileCache) FileExists(name string) bool {
	path := filepath.Join(c.root, name)

	c.mu.RLock()
	exists, cached := c.existsCache[path]
	c.mu.RUnlock()

	if cached {
		return exists
	}

	info, err := os.Stat(path)
	exists = err == nil && !info.IsDir()

	c.mu.Lock()
	c.existsCache[path] = exists
	c.mu.Unlock()

	return exists
}

// DirExists checks if a directory exists, using cache when available.
func (c *FileCache) DirExists(name string) bool {
	path := filepath.Join(c.root, name)

	// Use a different cache key prefix for directories
	cacheKey := "dir:" + path

	c.mu.RLock()
	exists, cached := c.existsCache[cacheKey]
	c.mu.RUnlock()

	if cached {
		return exists
	}

	info, err := os.Stat(path)
	exists = err == nil && info.IsDir()

	c.mu.Lock()
	c.existsCache[cacheKey] = exists
	c.mu.Unlock()

	return exists
}

// ReadFile reads a file's contents, using cache when available.
// Returns nil if the file doesn't exist or can't be read.
// The returned slice is a copy and safe to modify.
func (c *FileCache) ReadFile(name string) []byte {
	path := filepath.Join(c.root, name)

	c.mu.RLock()
	content, cached := c.contentCache[path]
	c.mu.RUnlock()

	if cached {
		// Return a copy to prevent mutation
		if content == nil {
			return nil
		}
		result := make([]byte, len(content))
		copy(result, content)
		return result
	}

	data, err := os.ReadFile(path)
	if err != nil {
		c.mu.Lock()
		c.contentCache[path] = nil
		c.mu.Unlock()
		return nil
	}

	// Store a copy in the cache
	cacheCopy := make([]byte, len(data))
	copy(cacheCopy, data)

	c.mu.Lock()
	c.contentCache[path] = cacheCopy
	c.mu.Unlock()

	// Return another copy for the caller
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

// FindFile returns the first existing file from the given paths.
// Uses cached file existence checks.
func (c *FileCache) FindFile(paths ...string) string {
	for _, p := range paths {
		if c.FileExists(p) {
			return p
		}
	}
	return ""
}

// FileContains checks if a file exists and contains the given substring.
// Uses cached file contents.
func (c *FileCache) FileContains(name, substr string) bool {
	content := c.ReadFile(name)
	if content == nil {
		return false
	}
	return containsBytes(content, []byte(substr))
}

// containsBytes checks if b contains sub.
func containsBytes(b, sub []byte) bool {
	if len(sub) == 0 {
		return true
	}
	if len(b) < len(sub) {
		return false
	}
	for i := 0; i <= len(b)-len(sub); i++ {
		if b[i] == sub[0] {
			match := true
			for j := 1; j < len(sub); j++ {
				if b[i+j] != sub[j] {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

// Clear clears all cached data.
func (c *FileCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.existsCache = make(map[string]bool)
	c.contentCache = make(map[string][]byte)
}

// Stats returns cache statistics.
func (c *FileCache) Stats() (existsCount, contentCount int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.existsCache), len(c.contentCache)
}
