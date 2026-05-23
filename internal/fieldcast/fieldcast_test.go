package fieldcast_test

import (
	"testing"

	"github.com/yourusername/logslice/internal/fieldcast"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	c := fieldcast.New()
	in := entry("level", "info", "msg", "hello")
	out := c.Apply(in)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestCastToString(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "code", Target: "string"}})
	out := c.Apply(entry("code", 404))
	if out["code"] != "404" {
		t.Fatalf("expected \"404\", got %v", out["code"])
	}
}

func TestCastToInt(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "count", Target: "int"}})
	out := c.Apply(entry("count", "42"))
	if out["count"] != int64(42) {
		t.Fatalf("expected int64(42), got %v (%T)", out["count"], out["count"])
	}
}

func TestCastToIntFromFloat(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "val", Target: "int"}})
	out := c.Apply(entry("val", "3.9"))
	if out["val"] != int64(3) {
		t.Fatalf("expected int64(3), got %v", out["val"])
	}
}

func TestCastToFloat(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "ratio", Target: "float"}})
	out := c.Apply(entry("ratio", "1.5"))
	if out["ratio"] != 1.5 {
		t.Fatalf("expected 1.5, got %v", out["ratio"])
	}
}

func TestCastToBool(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "ok", Target: "bool"}})
	out := c.Apply(entry("ok", "true"))
	if out["ok"] != true {
		t.Fatalf("expected true, got %v", out["ok"])
	}
}

func TestCastMissingFieldSkipped(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "missing", Target: "int"}})
	in := entry("level", "info")
	out := c.Apply(in)
	if _, exists := out["missing"]; exists {
		t.Fatal("missing field should not appear in output")
	}
}

func TestCastFailurePreservesOriginal(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "val", Target: "int"}})
	out := c.Apply(entry("val", "not-a-number"))
	if out["val"] != "not-a-number" {
		t.Fatalf("expected original value preserved, got %v", out["val"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	c := fieldcast.WithRules([]fieldcast.Rule{{Field: "count", Target: "int"}})
	in := entry("count", "7")
	_ = c.Apply(in)
	if in["count"] != "7" {
		t.Fatal("original entry was mutated")
	}
}
