package fieldpin_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldpin"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldpin.New(nil)
	in := entry("level", "info", "msg", "hello")
	out := tr.Apply(in)
	if len(out) != len(in) {
		t.Fatalf("expected %d keys, got %d", len(in), len(out))
	}
}

func TestPinCreatesSubObject(t *testing.T) {
	tr := fieldpin.New([]fieldpin.Rule{
		{Field: "level", Dest: "_pin"},
	})
	in := entry("level", "warn", "msg", "oops")
	out := tr.Apply(in)

	pin, ok := out["_pin"].(map[string]any)
	if !ok {
		t.Fatal("expected _pin to be a map")
	}
	if pin["level"] != "warn" {
		t.Fatalf("expected level=warn in pin, got %v", pin["level"])
	}
	// Original field still present.
	if out["level"] != "warn" {
		t.Fatal("original field should be preserved")
	}
}

func TestPinMissingFieldSkipped(t *testing.T) {
	tr := fieldpin.New([]fieldpin.Rule{
		{Field: "missing", Dest: "_pin"},
	})
	in := entry("level", "info")
	out := tr.Apply(in)
	if _, ok := out["_pin"]; ok {
		t.Fatal("_pin should not be created when source field is absent")
	}
}

func TestPinMultipleFieldsSameDest(t *testing.T) {
	tr := fieldpin.New([]fieldpin.Rule{
		{Field: "level", Dest: "_pin"},
		{Field: "msg", Dest: "_pin"},
	})
	in := entry("level", "error", "msg", "boom")
	out := tr.Apply(in)

	pin, ok := out["_pin"].(map[string]any)
	if !ok {
		t.Fatal("expected _pin map")
	}
	if pin["level"] != "error" || pin["msg"] != "boom" {
		t.Fatalf("unexpected pin contents: %v", pin)
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldpin.New([]fieldpin.Rule{
		{Field: "level", Dest: "_pin"},
	})
	in := entry("level", "debug")
	_ = tr.Apply(in)
	if _, ok := in["_pin"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}

func TestEmptyRuleFieldOrDestIgnored(t *testing.T) {
	tr := fieldpin.New([]fieldpin.Rule{
		{Field: "", Dest: "_pin"},
		{Field: "level", Dest: ""},
	})
	in := entry("level", "info")
	out := tr.Apply(in)
	if _, ok := out["_pin"]; ok {
		t.Fatal("no pin should be created for empty rules")
	}
}
