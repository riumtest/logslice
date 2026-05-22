package fieldclamp_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldclamp"
)

func ptr(f float64) *float64 { return &f }

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityClamp(t *testing.T) {
	c := fieldclamp.New(nil)
	in := entry("x", 5.0)
	out := c.Apply(in)
	if out["x"] != 5.0 {
		t.Fatalf("expected 5.0, got %v", out["x"])
	}
}

func TestClampMin(t *testing.T) {
	c := fieldclamp.New([]fieldclamp.Rule{{Field: "val", Min: ptr(0.0)}})
	out := c.Apply(entry("val", -3.0))
	if out["val"] != 0.0 {
		t.Fatalf("expected 0.0, got %v", out["val"])
	}
}

func TestClampMax(t *testing.T) {
	c := fieldclamp.New([]fieldclamp.Rule{{Field: "score", Max: ptr(100.0)}})
	out := c.Apply(entry("score", 150.0))
	if out["score"] != 100.0 {
		t.Fatalf("expected 100.0, got %v", out["score"])
	}
}

func TestClampWithinBounds(t *testing.T) {
	c := fieldclamp.New([]fieldclamp.Rule{{Field: "n", Min: ptr(1.0), Max: ptr(10.0)}})
	out := c.Apply(entry("n", 5.0))
	if out["n"] != 5.0 {
		t.Fatalf("expected 5.0, got %v", out["n"])
	}
}

func TestClampMissingFieldUnchanged(t *testing.T) {
	c := fieldclamp.New([]fieldclamp.Rule{{Field: "missing", Min: ptr(0.0)}})
	out := c.Apply(entry("other", 42.0))
	if _, ok := out["missing"]; ok {
		t.Fatal("missing field should not be added")
	}
}

func TestClampNonNumericSkipped(t *testing.T) {
	c := fieldclamp.New([]fieldclamp.Rule{{Field: "label", Min: ptr(0.0)}})
	out := c.Apply(entry("label", "hello"))
	if out["label"] != "hello" {
		t.Fatalf("expected original value, got %v", out["label"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	c := fieldclamp.New([]fieldclamp.Rule{{Field: "v", Max: ptr(10.0)}})
	in := entry("v", 99.0)
	c.Apply(in)
	if in["v"] != 99.0 {
		t.Fatal("original entry should not be mutated")
	}
}
