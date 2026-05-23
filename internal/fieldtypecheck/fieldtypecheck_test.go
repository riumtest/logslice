package fieldtypecheck_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldtypecheck"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	c := fieldtypecheck.New(nil)
	e := entry("msg", "hello")
	if got := c.Apply(e); got["msg"] != "hello" {
		t.Fatalf("expected unchanged entry")
	}
}

func TestNoMismatchPassesThrough(t *testing.T) {
	c := fieldtypecheck.New([]fieldtypecheck.Rule{{Field: "count", Expected: "number"}})
	e := entry("count", float64(5))
	out := c.Apply(e)
	if _, ok := out["_type_errors"]; ok {
		t.Fatal("expected no type errors")
	}
}

func TestMismatchWritesToDestField(t *testing.T) {
	c := fieldtypecheck.New(
		[]fieldtypecheck.Rule{{Field: "count", Expected: "number"}},
		fieldtypecheck.WithDestField("_type_errors"),
	)
	e := entry("count", "not-a-number")
	out := c.Apply(e)
	errs, ok := out["_type_errors"]
	if !ok {
		t.Fatal("expected _type_errors field")
	}
	list := errs.([]string)
	if len(list) != 1 {
		t.Fatalf("expected 1 error, got %d", len(list))
	}
}

func TestRejectModeDropsEntry(t *testing.T) {
	c := fieldtypecheck.New(
		[]fieldtypecheck.Rule{{Field: "active", Expected: "bool"}},
		fieldtypecheck.WithRejectMode(),
	)
	e := entry("active", "yes")
	if got := c.Apply(e); got != nil {
		t.Fatal("expected nil for rejected entry")
	}
}

func TestMissingFieldSkipped(t *testing.T) {
	c := fieldtypecheck.New(
		[]fieldtypecheck.Rule{{Field: "score", Expected: "number"}},
		fieldtypecheck.WithDestField("_type_errors"),
	)
	e := entry("msg", "hello")
	out := c.Apply(e)
	if _, ok := out["_type_errors"]; ok {
		t.Fatal("missing field should not trigger error")
	}
}

func TestOriginalNotMutated(t *testing.T) {
	c := fieldtypecheck.New(
		[]fieldtypecheck.Rule{{Field: "val", Expected: "number"}},
		fieldtypecheck.WithDestField("_type_errors"),
	)
	e := entry("val", "oops")
	c.Apply(e)
	if _, ok := e["_type_errors"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}
