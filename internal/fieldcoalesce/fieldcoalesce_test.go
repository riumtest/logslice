package fieldcoalesce_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldcoalesce"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr, err := fieldcoalesce.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	in := entry("a", "x")
	out := tr.Transform(in)
	if out["a"] != "x" {
		t.Fatalf("expected x, got %v", out["a"])
	}
}

func TestCoalesceFirstNonEmpty(t *testing.T) {
	tr, _ := fieldcoalesce.New([]fieldcoalesce.Rule{
		{Sources: []string{"a", "b", "c"}, Dest: "result"},
	})
	out := tr.Transform(entry("a", "", "b", "hello", "c", "world"))
	if out["result"] != "hello" {
		t.Fatalf("expected hello, got %v", out["result"])
	}
}

func TestCoalesceSkipsNil(t *testing.T) {
	tr, _ := fieldcoalesce.New([]fieldcoalesce.Rule{
		{Sources: []string{"x", "y"}, Dest: "result"},
	})
	out := tr.Transform(entry("x", nil, "y", 42))
	if out["result"] != 42 {
		t.Fatalf("expected 42, got %v", out["result"])
	}
}

func TestCoalesceNoMatchLeavesDestAbsent(t *testing.T) {
	tr, _ := fieldcoalesce.New([]fieldcoalesce.Rule{
		{Sources: []string{"a", "b"}, Dest: "result"},
	})
	out := tr.Transform(entry("a", "", "b", ""))
	if _, ok := out["result"]; ok {
		t.Fatal("expected result to be absent")
	}
}

func TestCoalesceDoesNotMutateInput(t *testing.T) {
	tr, _ := fieldcoalesce.New([]fieldcoalesce.Rule{
		{Sources: []string{"a"}, Dest: "result"},
	})
	in := entry("a", "val")
	tr.Transform(in)
	if _, ok := in["result"]; ok {
		t.Fatal("input entry was mutated")
	}
}

func TestNewErrorMissingSources(t *testing.T) {
	_, err := fieldcoalesce.New([]fieldcoalesce.Rule{
		{Sources: []string{}, Dest: "result"},
	})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestNewErrorMissingDest(t *testing.T) {
	_, err := fieldcoalesce.New([]fieldcoalesce.Rule{
		{Sources: []string{"a"}, Dest: ""},
	})
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}
