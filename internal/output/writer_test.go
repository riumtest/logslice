package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriterWritesLine(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf))

	if err := w.Write("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := buf.String(); got != "hello world\n" {
		t.Errorf("expected 'hello world\\n', got %q", got)
	}
}

func TestWriterCount(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf))

	for i := 0; i < 5; i++ {
		_ = w.Write("line")
	}

	if w.Count() != 5 {
		t.Errorf("expected count 5, got %d", w.Count())
	}
}

func TestWriterNoDoubleNewline(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf))

	_ = w.Write("already newline\n")

	if strings.Count(buf.String(), "\n") != 1 {
		t.Errorf("expected exactly one newline, got %q", buf.String())
	}
}

func TestWriterColorizeError(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf), WithColorize(true))

	_ = w.Write("level=error msg=something failed")

	if !strings.Contains(buf.String(), "\033[31m") {
		t.Error("expected red ANSI code for error level")
	}
}

func TestWriterColorizeWarn(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf), WithColorize(true))

	_ = w.Write("level=warn msg=low disk")

	if !strings.Contains(buf.String(), "\033[33m") {
		t.Error("expected yellow ANSI code for warn level")
	}
}

func TestWriterColorizeInfo(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf), WithColorize(true))

	_ = w.Write("level=info msg=started")

	if !strings.Contains(buf.String(), "\033[36m") {
		t.Error("expected cyan ANSI code for info level")
	}
}

func TestWriterNoColor(t *testing.T) {
	var buf bytes.Buffer
	w := New(WithDestination(&buf), WithColorize(false))

	_ = w.Write("level=error msg=oops")

	if strings.Contains(buf.String(), "\033[") {
		t.Error("expected no ANSI codes when colorize is false")
	}
}
