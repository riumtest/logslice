package reader_test

import (
	"io"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func TestReaderNext(t *testing.T) {
	input := `{"level":"info","msg":"started"}
{"level":"error","msg":"oops"}
`
	r := reader.New(strings.NewReader(input), "test")

	rec, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec["level"] != "info" {
		t.Errorf("expected level=info, got %v", rec["level"])
	}
	if rec["_source"] != "test" {
		t.Errorf("expected _source=test, got %v", rec["_source"])
	}

	rec, err = r.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec["level"] != "error" {
		t.Errorf("expected level=error, got %v", rec["level"])
	}

	_, err = r.Next()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestReaderSkipsEmptyLines(t *testing.T) {
	input := `{"a":1}

{"a":2}
`
	r := reader.New(strings.NewReader(input), "")
	records, errs := r.ReadAll()
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 records, got %d", len(records))
	}
}

func TestReaderInvalidJSON(t *testing.T) {
	input := `{"ok":true}
not-json
{"ok":false}
`
	r := reader.New(strings.NewReader(input), "src")
	records, errs := r.ReadAll()
	if len(records) != 2 {
		t.Errorf("expected 2 valid records, got %d", len(records))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}

func TestReaderNoSource(t *testing.T) {
	input := `{"x":"y"}
`
	r := reader.New(strings.NewReader(input), "")
	rec, err := r.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := rec["_source"]; ok {
		t.Error("expected no _source key when source is empty")
	}
}
