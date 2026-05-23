package fieldsanitize_test

import (
	"testing"

	"github.com/celrenheit/logslice/internal/fieldsanitize"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldsanitize.New()
	in := entry("msg", "  hello  ")
	out := tr.Apply(in)
	if out["msg"] != "  hello  " {
		t.Fatalf("expected value unchanged, got %q", out["msg"])
	}
}

func TestTrimWhitespace(t *testing.T) {
	tr := fieldsanitize.WithRules([]fieldsanitize.Rule{
		{Field: "msg"},
	})
	out := tr.Apply(entry("msg", "  hello world  "))
	if out["msg"] != "hello world" {
		t.Fatalf("expected trimmed value, got %q", out["msg"])
	}
}

func TestStripControlCharacters(t *testing.T) {
	tr := fieldsanitize.WithRules([]fieldsanitize.Rule{
		{Field: "msg", StripControl: true},
	})
	// \x01 is a control character; \t should survive
	out := tr.Apply(entry("msg", "hel\x01lo\tworld"))
	if out["msg"] != "hello\tworld" {
		t.Fatalf("expected control chars stripped, got %q", out["msg"])
	}
}

func TestSkipsNonStringField(t *testing.T) {
	tr := fieldsanitize.WithRules([]fieldsanitize.Rule{
		{Field: "count"},
	})
	out := tr.Apply(entry("count", 42))
	if out["count"] != 42 {
		t.Fatalf("expected numeric field unchanged, got %v", out["count"])
	}
}

func TestSkipsMissingField(t *testing.T) {
	tr := fieldsanitize.WithRules([]fieldsanitize.Rule{
		{Field: "missing"},
	})
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if _, ok := out["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldsanitize.WithRules([]fieldsanitize.Rule{
		{Field: "msg"},
	})
	in := entry("msg", "  hi  ")
	_ = tr.Apply(in)
	if in["msg"] != "  hi  " {
		t.Fatal("original entry was mutated")
	}
}

func TestMultipleRules(t *testing.T) {
	tr := fieldsanitize.WithRules([]fieldsanitize.Rule{
		{Field: "a"},
		{Field: "b", StripControl: true},
	})
	out := tr.Apply(entry("a", "  foo  ", "b", "bar\x00baz"))
	if out["a"] != "foo" {
		t.Fatalf("field a: expected %q got %q", "foo", out["a"])
	}
	if out["b"] != "barbaz" {
		t.Fatalf("field b: expected %q got %q", "barbaz", out["b"])
	}
}
