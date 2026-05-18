package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
)

func TestMainIntegration_NoFilter(t *testing.T) {
	// Write a temp log file and run the pipeline end-to-end.
	lines := []string{
		`{"level":"info","msg":"started"}`,
		`{"level":"error","msg":"failed"}`,
	}

	src := writeTempFile(t, lines)

	var buf bytes.Buffer
	out := output.New(output.WithDestination(&buf))
	writeFn := func(line string) error { return out.Write(line) }

	err := pipeline.Run([]string{src}, "", []formatter.Option{formatter.WithFormat("text")}, writeFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Count() != 2 {
		t.Errorf("expected 2 lines written, got %d", out.Count())
	}
}

func TestMainIntegration_WithFilter(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"ok"}`,
		`{"level":"error","msg":"boom"}`,
	}

	src := writeTempFile(t, lines)

	var buf bytes.Buffer
	out := output.New(output.WithDestination(&buf))
	writeFn := func(line string) error { return out.Write(line) }

	err := pipeline.Run([]string{src}, "level=error", []formatter.Option{formatter.WithFormat("text")}, writeFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Count() != 1 {
		t.Errorf("expected 1 line written, got %d", out.Count())
	}
	if !strings.Contains(buf.String(), "boom") {
		t.Errorf("expected 'boom' in output, got %q", buf.String())
	}
}

// writeTempFile creates a temp file with the given lines and returns its path.
func writeTempFile(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "logslice-test-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	f.Close()
	return f.Name()
}
