package fieldsplit_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldsplit"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentitySplit(t *testing.T) {
	s := fieldsplit.New()
	in := entry("msg", "hello")
	out := s.Apply(in)
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", out["msg"])
	}
}

func TestSplitBasic(t *testing.T) {
	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
		{Source: "addr", Delimiter: ":", Targets: []string{"host", "port"}},
	}))
	out := s.Apply(entry("addr", "localhost:8080"))
	if out["host"] != "localhost" {
		t.Errorf("host: got %v", out["host"])
	}
	if out["port"] != "8080" {
		t.Errorf("port: got %v", out["port"])
	}
	if out["addr"] != "localhost:8080" {
		t.Error("source field should be preserved")
	}
}

func TestSplitDefaultDelimiter(t *testing.T) {
	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
		{Source: "pair", Targets: []string{"a", "b"}},
	}))
	out := s.Apply(entry("pair", "foo:bar"))
	if out["a"] != "foo" || out["b"] != "bar" {
		t.Errorf("unexpected: %v", out)
	}
}

func TestSplitFewerPartsThanTargets(t *testing.T) {
	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
		{Source: "x", Delimiter: "/", Targets: []string{"p", "q", "r"}},
	}))
	out := s.Apply(entry("x", "alpha/beta"))
	if out["p"] != "alpha" || out["q"] != "beta" {
		t.Errorf("unexpected: %v", out)
	}
	if _, ok := out["r"]; ok {
		t.Error("r should not be set when part is missing")
	}
}

func TestSplitSkipsMissingSource(t *testing.T) {
	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
		{Source: "missing", Delimiter: ":", Targets: []string{"a"}},
	}))
	out := s.Apply(entry("msg", "hi"))
	if _, ok := out["a"]; ok {
		t.Error("a should not be set when source is absent")
	}
}

func TestSplitSkipsNonStringSource(t *testing.T) {
	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
		{Source: "code", Delimiter: ":", Targets: []string{"x"}},
	}))
	out := s.Apply(entry("code", 42))
	if _, ok := out["x"]; ok {
		t.Error("x should not be set for non-string source")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
		{Source: "kv", Delimiter: "=", Targets: []string{"k", "v"}},
	}))
	in := entry("kv", "foo=bar")
	_ = s.Apply(in)
	if _, ok := in["k"]; ok {
		t.Error("original entry should not be mutated")
	}
}
