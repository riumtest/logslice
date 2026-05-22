package fieldmath_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldmath"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestAddRule(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "a", Right: "b", Dest: "sum", Op: fieldmath.OpAdd},
	})
	out := tf.Apply(entry("a", float64(3), "b", float64(4)))
	if out["sum"] != float64(7) {
		t.Fatalf("expected 7, got %v", out["sum"])
	}
}

func TestSubRule(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "x", Right: "y", Dest: "diff", Op: fieldmath.OpSub},
	})
	out := tf.Apply(entry("x", float64(10), "y", float64(3)))
	if out["diff"] != float64(7) {
		t.Fatalf("expected 7, got %v", out["diff"])
	}
}

func TestMulRule(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "a", Right: "b", Dest: "prod", Op: fieldmath.OpMul},
	})
	out := tf.Apply(entry("a", float64(6), "b", float64(7)))
	if out["prod"] != float64(42) {
		t.Fatalf("expected 42, got %v", out["prod"])
	}
}

func TestDivRule(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "num", Right: "den", Dest: "ratio", Op: fieldmath.OpDiv},
	})
	out := tf.Apply(entry("num", float64(9), "den", float64(3)))
	if out["ratio"] != float64(3) {
		t.Fatalf("expected 3, got %v", out["ratio"])
	}
}

func TestDivByZeroSkipped(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "a", Right: "b", Dest: "r", Op: fieldmath.OpDiv},
	})
	out := tf.Apply(entry("a", float64(5), "b", float64(0)))
	if _, ok := out["r"]; ok {
		t.Fatal("expected dest field to be absent on div-by-zero")
	}
}

func TestMissingFieldSkipped(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "a", Right: "missing", Dest: "r", Op: fieldmath.OpAdd},
	})
	out := tf.Apply(entry("a", float64(1)))
	if _, ok := out["r"]; ok {
		t.Fatal("expected dest field to be absent when operand missing")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tf := fieldmath.New([]fieldmath.Rule{
		{Left: "a", Right: "b", Dest: "sum", Op: fieldmath.OpAdd},
	})
	orig := entry("a", float64(1), "b", float64(2))
	tf.Apply(orig)
	if _, ok := orig["sum"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}
