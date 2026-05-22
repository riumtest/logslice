package fieldround_test

import (
	"testing"

	"github.com/humanlogio/logslice/internal/fieldround"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldround.New(nil)
	in := entry("value", 3.7)
	out := tr.Apply(in)
	if out["value"] != 3.7 {
		t.Fatalf("expected 3.7, got %v", out["value"])
	}
}

func TestRoundDefault(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "val", Mode: "round", Precision: 0},
	})
	out := tr.Apply(entry("val", 2.6))
	if out["val"] != float64(3) {
		t.Fatalf("expected 3, got %v", out["val"])
	}
}

func TestFloor(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "val", Mode: "floor", Precision: 0},
	})
	out := tr.Apply(entry("val", 2.9))
	if out["val"] != float64(2) {
		t.Fatalf("expected 2, got %v", out["val"])
	}
}

func TestCeil(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "val", Mode: "ceil", Precision: 0},
	})
	out := tr.Apply(entry("val", 2.1))
	if out["val"] != float64(3) {
		t.Fatalf("expected 3, got %v", out["val"])
	}
}

func TestPrecision(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "val", Mode: "round", Precision: 2},
	})
	out := tr.Apply(entry("val", 1.2349))
	if out["val"] != 1.23 {
		t.Fatalf("expected 1.23, got %v", out["val"])
	}
}

func TestSkipsMissingField(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "missing", Mode: "round", Precision: 0},
	})
	in := entry("other", 5.5)
	out := tr.Apply(in)
	if _, ok := out["missing"]; ok {
		t.Fatal("unexpected key 'missing' in output")
	}
}

func TestSkipsNonNumeric(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "val", Mode: "round", Precision: 0},
	})
	out := tr.Apply(entry("val", "hello"))
	if out["val"] != "hello" {
		t.Fatalf("expected original string, got %v", out["val"])
	}
}

func TestOriginalNotMutated(t *testing.T) {
	tr := fieldround.New([]fieldround.Rule{
		{Field: "val", Mode: "floor", Precision: 0},
	})
	in := entry("val", 9.9)
	_ = tr.Apply(in)
	if in["val"] != 9.9 {
		t.Fatal("original entry was mutated")
	}
}
