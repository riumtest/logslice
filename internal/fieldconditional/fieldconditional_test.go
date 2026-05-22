package fieldconditional_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldconditional"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr, err := fieldconditional.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	in := entry("level", "info")
	out := tr.Apply(in)
	if out["level"] != "info" {
		t.Fatalf("expected info, got %v", out["level"])
	}
}

func TestEqRuleSetsField(t *testing.T) {
	tr, err := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "level", Op: "eq", Operand: "error", DestField: "alert", DestValue: "true"},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := tr.Apply(entry("level", "error", "msg", "boom"))
	if out["alert"] != "true" {
		t.Fatalf("expected alert=true, got %v", out["alert"])
	}
	out2 := tr.Apply(entry("level", "info"))
	if _, ok := out2["alert"]; ok {
		t.Fatal("alert should not be set for non-error level")
	}
}

func TestContainsRule(t *testing.T) {
	tr, _ := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "msg", Op: "contains", Operand: "timeout", DestField: "category", DestValue: "network"},
	})
	out := tr.Apply(entry("msg", "connection timeout reached"))
	if out["category"] != "network" {
		t.Fatalf("expected network, got %v", out["category"])
	}
}

func TestNeqRule(t *testing.T) {
	tr, _ := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "env", Op: "neq", Operand: "prod", DestField: "debug", DestValue: "1"},
	})
	out := tr.Apply(entry("env", "staging"))
	if out["debug"] != "1" {
		t.Fatalf("expected debug=1, got %v", out["debug"])
	}
	out2 := tr.Apply(entry("env", "prod"))
	if _, ok := out2["debug"]; ok {
		t.Fatal("debug should not be set for prod")
	}
}

func TestMissingSourceFieldSkipped(t *testing.T) {
	tr, _ := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "level", Op: "eq", Operand: "error", DestField: "alert", DestValue: "true"},
	})
	out := tr.Apply(entry("msg", "hello"))
	if _, ok := out["alert"]; ok {
		t.Fatal("alert should not be set when source field is absent")
	}
}

func TestOriginalNotMutated(t *testing.T) {
	tr, _ := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "level", Op: "eq", Operand: "error", DestField: "alert", DestValue: "true"},
	})
	in := entry("level", "error")
	tr.Apply(in)
	if _, ok := in["alert"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}

func TestInvalidOpReturnsError(t *testing.T) {
	_, err := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "level", Op: "regex", Operand: "err.*", DestField: "alert", DestValue: "true"},
	})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestEmptySourceFieldReturnsError(t *testing.T) {
	_, err := fieldconditional.New([]fieldconditional.Rule{
		{SourceField: "", Op: "eq", Operand: "x", DestField: "out", DestValue: "y"},
	})
	if err == nil {
		t.Fatal("expected error for empty source field")
	}
}
