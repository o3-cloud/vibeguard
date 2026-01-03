package inspector

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestFileCache_FileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	if err := os.WriteFile(filepath.Join(tmpDir, "exists.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	tests := []struct {
		name     string
		file     string
		expected bool
	}{
		{"existing file", "exists.txt", true},
		{"non-existing file", "notexists.txt", false},
		{"directory", "subdir", false}, // FileExists should return false for directories
		{"nested non-existing", "subdir/file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First call - not cached
			result := cache.FileExists(tt.file)
			if result != tt.expected {
				t.Errorf("FileExists(%q) = %v, want %v", tt.file, result, tt.expected)
			}

			// Second call - should use cache
			result = cache.FileExists(tt.file)
			if result != tt.expected {
				t.Errorf("FileExists(%q) cached = %v, want %v", tt.file, result, tt.expected)
			}
		})
	}

	// Verify cache was populated
	existsCount, _ := cache.Stats()
	if existsCount == 0 {
		t.Error("expected cache to be populated")
	}
}

func TestFileCache_DirExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test directory
	if err := os.MkdirAll(filepath.Join(tmpDir, "mydir"), 0755); err != nil {
		t.Fatal(err)
	}
	// Create a file (not a directory)
	if err := os.WriteFile(filepath.Join(tmpDir, "myfile"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	tests := []struct {
		name     string
		dir      string
		expected bool
	}{
		{"existing dir", "mydir", true},
		{"non-existing dir", "notexists", false},
		{"file not dir", "myfile", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cache.DirExists(tt.dir)
			if result != tt.expected {
				t.Errorf("DirExists(%q) = %v, want %v", tt.dir, result, tt.expected)
			}

			// Second call - should use cache
			result = cache.DirExists(tt.dir)
			if result != tt.expected {
				t.Errorf("DirExists(%q) cached = %v, want %v", tt.dir, result, tt.expected)
			}
		})
	}
}

func TestFileCache_ReadFile(t *testing.T) {
	tmpDir := t.TempDir()

	content := "hello world"
	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	// First read
	data := cache.ReadFile("test.txt")
	if string(data) != content {
		t.Errorf("ReadFile() = %q, want %q", string(data), content)
	}

	// Second read (cached)
	data2 := cache.ReadFile("test.txt")
	if string(data2) != content {
		t.Errorf("ReadFile() cached = %q, want %q", string(data2), content)
	}

	// Non-existing file
	data3 := cache.ReadFile("notexists.txt")
	if data3 != nil {
		t.Errorf("ReadFile(notexists) = %v, want nil", data3)
	}

	// Verify content cache was populated
	_, contentCount := cache.Stats()
	if contentCount < 2 { // test.txt and notexists.txt (nil entry)
		t.Errorf("expected content cache to have 2 entries, got %d", contentCount)
	}
}

func TestFileCache_FindFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some test files
	if err := os.WriteFile(filepath.Join(tmpDir, "config.yml"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "config.json"), []byte("b"), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	tests := []struct {
		name     string
		paths    []string
		expected string
	}{
		{"first match", []string{"config.yml", "config.json"}, "config.yml"},
		{"second match", []string{"notexists", "config.json"}, "config.json"},
		{"no match", []string{"notexists1", "notexists2"}, ""},
		{"json first", []string{"config.json", "config.yml"}, "config.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cache.FindFile(tt.paths...)
			if result != tt.expected {
				t.Errorf("FindFile(%v) = %q, want %q", tt.paths, result, tt.expected)
			}
		})
	}
}

func TestFileCache_FileContains(t *testing.T) {
	tmpDir := t.TempDir()

	content := "hello world\nthis is a test\nwith multiple lines"
	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	tests := []struct {
		name     string
		file     string
		substr   string
		expected bool
	}{
		{"contains exact", "test.txt", "hello world", true},
		{"contains partial", "test.txt", "world", true},
		{"contains middle", "test.txt", "is a test", true},
		{"not contains", "test.txt", "foobar", false},
		{"file not exists", "notexists.txt", "anything", false},
		{"empty substr", "test.txt", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cache.FileContains(tt.file, tt.substr)
			if result != tt.expected {
				t.Errorf("FileContains(%q, %q) = %v, want %v", tt.file, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestFileCache_Clear(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	// Populate cache
	cache.FileExists("test.txt")
	cache.ReadFile("test.txt")

	existsCount, contentCount := cache.Stats()
	if existsCount == 0 || contentCount == 0 {
		t.Fatal("cache should be populated before clear")
	}

	// Clear cache
	cache.Clear()

	existsCount, contentCount = cache.Stats()
	if existsCount != 0 || contentCount != 0 {
		t.Errorf("cache not cleared: exists=%d, content=%d", existsCount, contentCount)
	}
}

func TestFileCache_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	for i := 0; i < 10; i++ {
		filename := filepath.Join(tmpDir, filepath.Base(t.TempDir())+".txt")
		if err := os.WriteFile(filename, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "shared.txt"), []byte("shared"), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Concurrent reads should not panic or race
			cache.FileExists("shared.txt")
			cache.ReadFile("shared.txt")
			cache.FileContains("shared.txt", "shared")
			cache.DirExists("nonexistent")
		}()
	}
	wg.Wait()

	// Cache should be populated
	existsCount, contentCount := cache.Stats()
	if existsCount == 0 {
		t.Error("expected exists cache to be populated after concurrent access")
	}
	if contentCount == 0 {
		t.Error("expected content cache to be populated after concurrent access")
	}
}

func TestContainsBytes(t *testing.T) {
	tests := []struct {
		name     string
		b        []byte
		sub      []byte
		expected bool
	}{
		{"empty sub", []byte("hello"), []byte(""), true},
		{"exact match", []byte("hello"), []byte("hello"), true},
		{"prefix match", []byte("hello world"), []byte("hello"), true},
		{"suffix match", []byte("hello world"), []byte("world"), true},
		{"middle match", []byte("hello world"), []byte("lo wo"), true},
		{"no match", []byte("hello"), []byte("world"), false},
		{"sub longer than b", []byte("hi"), []byte("hello"), false},
		{"empty b", []byte(""), []byte("a"), false},
		{"both empty", []byte(""), []byte(""), true},
		{"single char match", []byte("abc"), []byte("b"), true},
		{"single char no match", []byte("abc"), []byte("d"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsBytes(tt.b, tt.sub)
			if result != tt.expected {
				t.Errorf("containsBytes(%q, %q) = %v, want %v", tt.b, tt.sub, result, tt.expected)
			}
		})
	}
}

func TestFileCache_ReadFileMutationSafe(t *testing.T) {
	tmpDir := t.TempDir()

	content := "original"
	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cache := NewFileCache(tmpDir)

	// First read
	data1 := cache.ReadFile("test.txt")

	// Mutate the returned slice
	if len(data1) > 0 {
		data1[0] = 'X'
	}

	// Second read should still return original
	data2 := cache.ReadFile("test.txt")
	if string(data2) != content {
		t.Errorf("cache was mutated: got %q, want %q", string(data2), content)
	}
}
