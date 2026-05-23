package fieldnull_test

import (
	"testing"

	"github.com/humanlogio/logslice/internal/fieldnull"
)

func entry(pairs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldnull.New(nil)
	in := entry("msg", "hello", "level", "info")
	out := tr.Transform(in)
	if out == nil || out["msg"] != "hello" {
		t.Fatalf("expected unchanged entry, got %v", out)
	}
}

func TestRemovesNullField(t *testing.T) {
	tr := fieldnull.New([]fieldnull.Rule{{Field: "err"}})
	in := entry("msg", "ok", "err", nil)
	out := tr.Transform(in)
	if _, ok := out["err"]; ok {
		t.Fatal("expected err field to be removed")
	}
}

func TestReplacesNullField(t *testing.T) {
	tr := fieldnull.New([]fieldnull.Rule{{Field: "region", Replace: "unknown"}})
	in := entry("msg", "ok", "region", nil)
	out := tr.Transform(in)
	if out["region"] != "unknown" {
		t.Fatalf("expected 'unknown', got %v", out["region"])
	}
}

func TestDropEntryOnNull(t *testing.T) {
	tr := fieldnull.New([]fieldnull.Rule{{Field: "user_id", DropEntry: true}})
	in := entry("msg", "login", "user_id", nil)
	out := tr.Transform(in)
	if out != nil {
		t.Fatalf("expected nil (dropped), got %v", out)
	}
}

func TestNonNullFieldUnchanged(t *testing.T) {
	tr := fieldnull.New([]fieldnull.Rule{{Field: "count", Replace: float64(0)}})
	in := entry("count", float64(42))
	out := tr.Transform(in)
	if out["count"] != float64(42) {
		t.Fatalf("expected 42, got %v", out["count"])
	}
}

func TestMissingFieldTreatedAsNull(t *testing.T) {
	tr := fieldnull.New([]fieldnull.Rule{{Field: "trace_id", Replace: "n/a"}})
	in := entry("msg", "span")
	out := tr.Transform(in)
	if out["trace_id"] != "n/a" {
		t.Fatalf("expected 'n/a' for missing field, got %v", out["trace_id"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldnull.New([]fieldnull.Rule{{Field: "x", Replace: "default"}})
	in := entry("x", nil, "y", "keep")
	_ = tr.Transform(in)
	if in["x"] != nil {
		t.Fatal("original entry was mutated")
	}
}
