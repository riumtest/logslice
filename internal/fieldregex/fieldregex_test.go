package fieldregex_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldregex"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	e, err := fieldregex.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("msg", "hello")
	out := e.Apply(in)
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
}

func TestExtractNamedGroups(t *testing.T) {
	e, err := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{
		{Source: "msg", Pattern: `(?P<level>\w+)\s+(?P<code>\d+)`},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := e.Apply(entry("msg", "ERROR 503 upstream"))
	if out["level"] != "ERROR" {
		t.Errorf("expected level=ERROR, got %v", out["level"])
	}
	if out["code"] != "503" {
		t.Errorf("expected code=503, got %v", out["code"])
	}
}

func TestExtractWithPrefix(t *testing.T) {
	e, err := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{
		{Source: "path", Pattern: `(?P<name>[^/]+)$`, Prefix: "file_"},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := e.Apply(entry("path", "/var/log/app.log"))
	if out["file_name"] != "app.log" {
		t.Errorf("expected file_name=app.log, got %v", out["file_name"])
	}
}

func TestNoMatchLeavesEntryUnchanged(t *testing.T) {
	e, err := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{
		{Source: "msg", Pattern: `(?P<ip>\d{1,3}(?:\.\d{1,3}){3})`},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := e.Apply(entry("msg", "no ip here"))
	if _, ok := out["ip"]; ok {
		t.Errorf("expected ip field to be absent")
	}
}

func TestMissingSourceFieldSkipped(t *testing.T) {
	e, err := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{
		{Source: "nonexistent", Pattern: `(?P<x>\w+)`},
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := e.Apply(entry("msg", "hello"))
	if _, ok := out["x"]; ok {
		t.Errorf("expected x field to be absent")
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{
		{Source: "msg", Pattern: `(?P<bad`},
	}))
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	e, _ := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{
		{Source: "msg", Pattern: `(?P<word>\w+)`},
	}))
	in := entry("msg", "hello")
	e.Apply(in)
	if _, ok := in["word"]; ok {
		t.Errorf("original entry should not be mutated")
	}
}
