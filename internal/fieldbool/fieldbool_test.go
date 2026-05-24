package fieldbool_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldbool"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldbool.New(nil)
	in := entry("active", "yes")
	out := tr.Apply(in)
	if out["active"] != "yes" {
		t.Fatalf("expected unchanged value, got %v", out["active"])
	}
}

func TestCoerceTruthy(t *testing.T) {
	cases := []any{"true", "yes", "1", "on", "TRUE", "YES"}
	for _, c := range cases {
		tr := fieldbool.New([]fieldbool.Rule{{Field: "f"}})
		out := tr.Apply(entry("f", c))
		if out["f"] != true {
			t.Errorf("input %v: expected true, got %v", c, out["f"])
		}
	}
}

func TestCoerceFalsy(t *testing.T) {
	cases := []any{"false", "no", "0", "off", "FALSE"}
	for _, c := range cases {
		tr := fieldbool.New([]fieldbool.Rule{{Field: "f"}})
		out := tr.Apply(entry("f", c))
		if out["f"] != false {
			t.Errorf("input %v: expected false, got %v", c, out["f"])
		}
	}
}

func TestCoerceBoolPassthrough(t *testing.T) {
	tr := fieldbool.New([]fieldbool.Rule{{Field: "ok"}})
	out := tr.Apply(entry("ok", true))
	if out["ok"] != true {
		t.Fatalf("expected true, got %v", out["ok"])
	}
}

func TestCoerceNumeric(t *testing.T) {
	tr := fieldbool.New([]fieldbool.Rule{{Field: "n"}})
	out := tr.Apply(entry("n", float64(42)))
	if out["n"] != true {
		t.Fatalf("expected true for non-zero float, got %v", out["n"])
	}
	out2 := tr.Apply(entry("n", float64(0)))
	if out2["n"] != false {
		t.Fatalf("expected false for zero float, got %v", out2["n"])
	}
}

func TestCoerceWritesToDest(t *testing.T) {
	tr := fieldbool.New([]fieldbool.Rule{{Field: "raw", Dest: "enabled"}})
	out := tr.Apply(entry("raw", "yes"))
	if out["enabled"] != true {
		t.Fatalf("expected true in dest field, got %v", out["enabled"])
	}
	if _, exists := out["raw"]; !exists {
		t.Fatal("source field should still be present")
	}
}

func TestUnrecognisedStringSkipped(t *testing.T) {
	tr := fieldbool.New([]fieldbool.Rule{{Field: "x"}})
	in := entry("x", "maybe")
	out := tr.Apply(in)
	// unrecognised string: field should be unchanged
	if out["x"] != "maybe" {
		t.Fatalf("expected original value preserved, got %v", out["x"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldbool.New([]fieldbool.Rule{{Field: "flag"}})
	in := entry("flag", "on")
	_ = tr.Apply(in)
	if in["flag"] != "on" {
		t.Fatal("original entry must not be mutated")
	}
}
