package fieldcopy_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldcopy"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldcopy.New(nil)
	in := entry("level", "info", "msg", "hello")
	out := tr.Apply(in)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestCopyField(t *testing.T) {
	tr := fieldcopy.New([]fieldcopy.Rule{{Src: "msg", Dst: "message"}})
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if out["msg"] != "hello" {
		t.Errorf("source field removed: %v", out)
	}
	if out["message"] != "hello" {
		t.Errorf("destination field not set: %v", out)
	}
}

func TestCopyMissingSourceSkipped(t *testing.T) {
	tr := fieldcopy.New([]fieldcopy.Rule{{Src: "missing", Dst: "dst"}})
	in := entry("level", "warn")
	out := tr.Apply(in)
	if _, ok := out["dst"]; ok {
		t.Errorf("dst should not be set when source is missing: %v", out)
	}
}

func TestCopyMultipleRules(t *testing.T) {
	tr := fieldcopy.New([]fieldcopy.Rule{
		{Src: "level", Dst: "severity"},
		{Src: "msg", Dst: "text"},
	})
	in := entry("level", "error", "msg", "boom")
	out := tr.Apply(in)
	if out["severity"] != "error" {
		t.Errorf("severity not copied: %v", out)
	}
	if out["text"] != "boom" {
		t.Errorf("text not copied: %v", out)
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldcopy.New([]fieldcopy.Rule{{Src: "a", Dst: "b"}})
	in := entry("a", "value")
	_ = tr.Apply(in)
	if _, ok := in["b"]; ok {
		t.Errorf("original entry was mutated")
	}
}

func TestEmptyRuleSkipped(t *testing.T) {
	tr := fieldcopy.New([]fieldcopy.Rule{
		{Src: "", Dst: "dst"},
		{Src: "src", Dst: ""},
		{Src: "level", Dst: "sev"},
	})
	in := entry("level", "debug")
	out := tr.Apply(in)
	if out["sev"] != "debug" {
		t.Errorf("valid rule not applied: %v", out)
	}
}
