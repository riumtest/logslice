package source_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/source"
)

func TestExpandGlobsNoWildcard(t *testing.T) {
	result, err := source.ExpandGlobs([]string{"/some/file.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "/some/file.log" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestExpandGlobsWithWildcard(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "a.log"), []byte{}, 0600)
	_ = os.WriteFile(filepath.Join(dir, "b.log"), []byte{}, 0600)
	_ = os.WriteFile(filepath.Join(dir, "c.txt"), []byte{}, 0600)

	pattern := filepath.Join(dir, "*.log")
	result, err := source.ExpandGlobs([]string{pattern})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 matches, got %d: %v", len(result), result)
	}
}

func TestExpandGlobsDeduplicate(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "x.log")
	_ = os.WriteFile(file, []byte{}, 0600)

	result, err := source.ExpandGlobs([]string{file, file})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 unique path, got %d", len(result))
	}
}

func TestExpandGlobsNoMatchPreservesPattern(t *testing.T) {
	result, err := source.ExpandGlobs([]string{"/no/match/*.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should preserve the pattern so downstream gives a clear open error.
	if len(result) != 1 {
		t.Errorf("expected pattern preserved, got %v", result)
	}
}
