package fieldscope_test

import (
	"testing"

	"github.com/your-org/logslice/internal/fieldscope"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	t.Parallel()
	tr := fieldscope.New(nil)
	in := entry("level", "info", "msg", "hello")
	out := tr.Apply(in)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestPromoteNestedField(t *testing.T) {
	t.Parallel()
	tr := fieldscope.New([]fieldscope.Rule{
		{Source: "meta.user", Dest: "user"},
	})
	in := entry("meta", map[string]any{"user": "alice", "ip": "1.2.3.4"})
	out := tr.Apply(in)
	if out["user"] != "alice" {
		t.Fatalf("expected user=alice, got %v", out["user"])
	}
	// original still intact
	if in["user"] != nil {
		t.Fatal("original entry should not be mutated")
	}
}

func TestPromoteDeeplyNested(t *testing.T) {
	t.Parallel()
	tr := fieldscope.New([]fieldscope.Rule{
		{Source: "a.b.c", Dest: "value"},
	})
	in := entry("a", map[string]any{
		"b": map[string]any{"c": float64(42)},
	})
	out := tr.Apply(in)
	if out["value"] != float64(42) {
		t.Fatalf("expected value=42, got %v", out["value"])
	}
}

func TestMissingPathSkipped(t *testing.T) {
	t.Parallel()
	tr := fieldscope.New([]fieldscope.Rule{
		{Source: "meta.missing", Dest: "x"},
	})
	in := entry("meta", map[string]any{"user": "bob"})
	out := tr.Apply(in)
	if _, exists := out["x"]; exists {
		t.Fatal("dest key should not be set when source path is missing")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	t.Parallel()
	tr := fieldscope.New([]fieldscope.Rule{
		{Source: "req.method", Dest: "method"},
	})
	in := entry("req", map[string]any{"method": "GET"})
	_ = tr.Apply(in)
	if _, exists := in["method"]; exists {
		t.Fatal("original entry must not be mutated")
	}
}

func TestMultipleRules(t *testing.T) {
	t.Parallel()
	tr := fieldscope.New([]fieldscope.Rule{
		{Source: "http.method", Dest: "method"},
		{Source: "http.status", Dest: "status"},
	})
	in := entry("http", map[string]any{"method": "POST", "status": float64(200)})
	out := tr.Apply(in)
	if out["method"] != "POST" {
		t.Fatalf("expected method=POST, got %v", out["method"])
	}
	if out["status"] != float64(200) {
		t.Fatalf("expected status=200, got %v", out["status"])
	}
}
