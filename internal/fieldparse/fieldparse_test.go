package fieldparse_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldparse"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityTransform(t *testing.T) {
	tr := fieldparse.New()
	in := entry("msg", "hello", "count", "42")
	out := tr.Transform(in)
	if out["count"] != "42" {
		t.Fatalf("expected string '42', got %v", out["count"])
	}
}

func TestParseInteger(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("count"))
	out := tr.Transform(entry("count", "99"))
	if out["count"] != float64(99) {
		t.Fatalf("expected float64(99), got %v (%T)", out["count"], out["count"])
	}
}

func TestParseFloat(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("ratio"))
	out := tr.Transform(entry("ratio", "3.14"))
	if out["ratio"] != float64(3.14) {
		t.Fatalf("expected 3.14, got %v", out["ratio"])
	}
}

func TestParseBool(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("ok"))
	out := tr.Transform(entry("ok", "true"))
	if out["ok"] != true {
		t.Fatalf("expected true, got %v", out["ok"])
	}
}

func TestParseJSON(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("meta"))
	out := tr.Transform(entry("meta", `{"k":1}`))
	m, ok := out["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", out["meta"])
	}
	if m["k"] != float64(1) {
		t.Fatalf("expected k=1, got %v", m["k"])
	}
}

func TestParseUnparseable(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("tag"))
	out := tr.Transform(entry("tag", "not-a-number"))
	if out["tag"] != "not-a-number" {
		t.Fatalf("expected original string, got %v", out["tag"])
	}
}

func TestSkipsMissingField(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("missing"))
	out := tr.Transform(entry("other", "val"))
	if _, exists := out["missing"]; exists {
		t.Fatal("unexpected key 'missing' in output")
	}
}

func TestSkipsNonStringField(t *testing.T) {
	tr := fieldparse.New(fieldparse.WithFields("count"))
	out := tr.Transform(entry("count", float64(7)))
	if out["count"] != float64(7) {
		t.Fatalf("expected unchanged float64(7), got %v", out["count"])
	}
}
