package fieldprefix_test

import (
	"testing"

	"github.com/naturalselectionlabs/logslice/internal/fieldprefix"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldprefix.New(nil)
	in := entry("msg", "hello")
	out := tr.Transform(in)
	if out["msg"] != "hello" {
		t.Fatalf("expected hello, got %v", out["msg"])
	}
}

func TestAddPrefix(t *testing.T) {
	tr := fieldprefix.New([]fieldprefix.Rule{
		{Field: "env", Prefix: "prod-"},
	})
	out := tr.Transform(entry("env", "us-east"))
	if got := out["env"]; got != "prod-us-east" {
		t.Fatalf("expected prod-us-east, got %v", got)
	}
}

func TestStripPrefix(t *testing.T) {
	tr := fieldprefix.New([]fieldprefix.Rule{
		{Field: "env", Prefix: "prod-", Strip: true},
	})
	out := tr.Transform(entry("env", "prod-us-east"))
	if got := out["env"]; got != "us-east" {
		t.Fatalf("expected us-east, got %v", got)
	}
}

func TestStripPrefixNoMatch(t *testing.T) {
	tr := fieldprefix.New([]fieldprefix.Rule{
		{Field: "env", Prefix: "prod-", Strip: true},
	})
	out := tr.Transform(entry("env", "staging"))
	if got := out["env"]; got != "staging" {
		t.Fatalf("expected staging unchanged, got %v", got)
	}
}

func TestSkipsMissingField(t *testing.T) {
	tr := fieldprefix.New([]fieldprefix.Rule{
		{Field: "missing", Prefix: "x-"},
	})
	out := tr.Transform(entry("other", "val"))
	if _, ok := out["missing"]; ok {
		t.Fatal("missing field should not be created")
	}
}

func TestSkipsNonStringField(t *testing.T) {
	tr := fieldprefix.New([]fieldprefix.Rule{
		{Field: "code", Prefix: "err-"},
	})
	out := tr.Transform(entry("code", 42))
	if got := out["code"]; got != 42 {
		t.Fatalf("expected numeric 42 unchanged, got %v", got)
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldprefix.New([]fieldprefix.Rule{
		{Field: "svc", Prefix: "svc-"},
	})
	in := entry("svc", "auth")
	_ = tr.Transform(in)
	if in["svc"] != "auth" {
		t.Fatal("original entry must not be mutated")
	}
}
