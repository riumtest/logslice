package source_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/source"
)

func TestResolveStdin(t *testing.T) {
	fakeStdin := strings.NewReader("hello\n")
	rs := source.New(source.WithStdin(fakeStdin))

	sources, err := rs.Resolve(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(sources))
	}
	defer sources[0].Close()

	data, _ := io.ReadAll(sources[0])
	if string(data) != "hello\n" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestResolveFiles(t *testing.T) {
	dir := t.TempDir()

	file1 := filepath.Join(dir, "a.log")
	file2 := filepath.Join(dir, "b.log")
	_ = os.WriteFile(file1, []byte("line1\n"), 0600)
	_ = os.WriteFile(file2, []byte("line2\n"), 0600)

	rs := source.New()
	sources, err := rs.Resolve([]string{file1, file2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(sources))
	}
	for _, s := range sources {
		s.Close()
	}
}

func TestResolveMissingFile(t *testing.T) {
	rs := source.New()
	_, err := rs.Resolve([]string{"/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestResolvePartialFailureClosesOpened(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.log")
	_ = os.WriteFile(good, []byte("ok\n"), 0600)

	rs := source.New()
	_, err := rs.Resolve([]string{good, "/no/such/file.log"})
	if err == nil {
		t.Fatal("expected error when second file is missing")
	}
}
