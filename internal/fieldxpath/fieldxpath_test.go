package fieldxpath_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldxpath"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldxpath.New(nil)
	in := entry("a", "b")
	out := tr.Apply(in)
	if out["a"] != "b" {
		t.Fatalf("expected a=b, got %v", out["a"])
	}
}

func TestExtractTopLevel(t *testing.T) {
	tr := fieldxpath.New([]fieldxpath.Rule{{Path: "msg", Dest: "message"}})
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if out["message"] != "hello" {
		t.Fatalf("expected message=hello, got %v", out["message"])
	}
}

func TestExtractNestedField(t *testing.T) {
	tr := fieldxpath.New([]fieldxpath.Rule{{Path: "meta.region", Dest: "region"}})
	in := entry("meta", map[string]any{"region": "us-east-1"})
	out := tr.Apply(in)
	if out["region"] != "us-east-1" {
		t.Fatalf("expected region=us-east-1, got %v", out["region"])
	}
}

func TestExtractDeeplyNested(t *testing.T) {
	tr := fieldxpath.New([]fieldxpath.Rule{{Path: "a.b.c", Dest: "deep"}})
	in := entry("a", map[string]any{"b": map[string]any{"c": 42.0}})
	out := tr.Apply(in)
	if out["deep"] != 42.0 {
		t.Fatalf("expected deep=42, got %v", out["deep"])
	}
}

func TestMissingPathSkipped(t *testing.T) {
	tr := fieldxpath.New([]fieldxpath.Rule{{Path: "x.y", Dest: "result"}})
	in := entry("foo", "bar")
	out := tr.Apply(in)
	if _, ok := out["result"]; ok {
		t.Fatal("expected result to be absent")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldxpath.New([]fieldxpath.Rule{{Path: "meta.host", Dest: "host"}})
	in := entry("meta", map[string]any{"host": "srv1"})
	tr.Apply(in)
	if _, ok := in["host"]; ok {
		t.Fatal("original entry should not have been mutated")
	}
}

func TestIntermediateNonMapSkipped(t *testing.T) {
	tr := fieldxpath.New([]fieldxpath.Rule{{Path: "a.b", Dest: "out"}})
	in := entry("a", "not-a-map")
	out := tr.Apply(in)
	if _, ok := out["out"]; ok {
		t.Fatal("expected out to be absent when path traversal fails")
	}
}
