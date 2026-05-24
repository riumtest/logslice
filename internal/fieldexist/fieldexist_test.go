package fieldexist_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldexist"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldexist.New(nil)
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if out == nil || out["msg"] != "hello" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestFilterDropsWhenFieldMissing(t *testing.T) {
	tr := fieldexist.New([]fieldexist.Rule{
		{Field: "level", MustExist: true},
	})
	out := tr.Apply(entry("msg", "no level here"))
	if out != nil {
		t.Fatalf("expected nil (drop), got %v", out)
	}
}

func TestFilterKeepsWhenFieldPresent(t *testing.T) {
	tr := fieldexist.New([]fieldexist.Rule{
		{Field: "level", MustExist: true},
	})
	out := tr.Apply(entry("level", "info", "msg", "ok"))
	if out == nil {
		t.Fatal("expected entry to be kept")
	}
}

func TestFilterDropsWhenFieldPresent_MustNotExist(t *testing.T) {
	tr := fieldexist.New([]fieldexist.Rule{
		{Field: "debug", MustExist: false},
	})
	out := tr.Apply(entry("debug", true, "msg", "verbose"))
	if out != nil {
		t.Fatalf("expected nil (drop), got %v", out)
	}
}

func TestAnnotateWritesBoolToDestField(t *testing.T) {
	tr := fieldexist.New([]fieldexist.Rule{
		{Field: "error", DestField: "has_error"},
	})

	with := tr.Apply(entry("error", "oops", "msg", "bad"))
	if with == nil || with["has_error"] != true {
		t.Fatalf("expected has_error=true, got %v", with)
	}

	without := tr.Apply(entry("msg", "fine"))
	if without == nil || without["has_error"] != false {
		t.Fatalf("expected has_error=false, got %v", without)
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldexist.New([]fieldexist.Rule{
		{Field: "x", DestField: "has_x"},
	})
	in := entry("x", 1)
	tr.Apply(in)
	if _, ok := in["has_x"]; ok {
		t.Fatal("original entry was mutated")
	}
}
