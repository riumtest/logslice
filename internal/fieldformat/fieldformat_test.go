package fieldformat_test

import (
	"testing"

	"github.com/mikelorant/logslice/internal/entry"
	"github.com/mikelorant/logslice/internal/fieldformat"
)

func entry(pairs ...interface{}) entry.Entry {
	e := make(entry.Entry)
	for i := 0; i+1 < len(pairs); i += 2 {
		e[pairs[i].(string)] = pairs[i+1]
	}
	return e
}

func TestIdentityFormat(t *testing.T) {
	f := fieldformat.New()
	in := entry("msg", "Hello World")
	out := f.Apply(in)
	if out["msg"] != "Hello World" {
		t.Fatalf("expected unchanged value, got %v", out["msg"])
	}
}

func TestUppercase(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "level", Op: fieldformat.OpUppercase},
	})
	out := f.Apply(entry("level", "warn"))
	if out["level"] != "WARN" {
		t.Fatalf("expected WARN, got %v", out["level"])
	}
}

func TestLowercase(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "level", Op: fieldformat.OpLowercase},
	})
	out := f.Apply(entry("level", "ERROR"))
	if out["level"] != "error" {
		t.Fatalf("expected error, got %v", out["level"])
	}
}

func TestTruncate(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "msg", Op: fieldformat.OpTruncate, MaxLen: 5},
	})
	out := f.Apply(entry("msg", "Hello World"))
	if out["msg"] != "Hello" {
		t.Fatalf("expected 'Hello', got %v", out["msg"])
	}
}

func TestTruncateShortString(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "msg", Op: fieldformat.OpTruncate, MaxLen: 20},
	})
	out := f.Apply(entry("msg", "Hi"))
	if out["msg"] != "Hi" {
		t.Fatalf("expected 'Hi' unchanged, got %v", out["msg"])
	}
}

func TestSkipsMissingField(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "missing", Op: fieldformat.OpUppercase},
	})
	out := f.Apply(entry("msg", "hello"))
	if out["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %v", out["msg"])
	}
}

func TestSkipsNonStringField(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "count", Op: fieldformat.OpUppercase},
	})
	out := f.Apply(entry("count", 42))
	if out["count"] != 42 {
		t.Fatalf("expected 42 unchanged, got %v", out["count"])
	}
}

func TestOriginalNotMutated(t *testing.T) {
	f := fieldformat.WithRules([]fieldformat.Rule{
		{Field: "level", Op: fieldformat.OpUppercase},
	})
	in := entry("level", "info")
	_ = f.Apply(in)
	if in["level"] != "info" {
		t.Fatalf("original entry mutated")
	}
}
