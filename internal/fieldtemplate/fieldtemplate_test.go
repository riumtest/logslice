package fieldtemplate_test

import (
	"testing"

	"github.com/celery/logslice/internal/fieldtemplate"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldtemplate.New()
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", out["msg"])
	}
}

func TestRenderSimpleTemplate(t *testing.T) {
	tr, err := fieldtemplate.WithRules([]fieldtemplate.Rule{
		{DestField: "summary", Tmpl: `{{index . "level"}}: {{index . "msg"}}`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("level", "error", "msg", "disk full")
	out := tr.Apply(in)
	if out["summary"] != "error: disk full" {
		t.Fatalf("expected 'error: disk full', got %q", out["summary"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr, _ := fieldtemplate.WithRules([]fieldtemplate.Rule{
		{DestField: "rendered", Tmpl: `hello`},
	})
	in := entry("x", "1")
	_ = tr.Apply(in)
	if _, ok := in["rendered"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestMissingKeyDoesNotPanic(t *testing.T) {
	tr, err := fieldtemplate.WithRules([]fieldtemplate.Rule{
		{DestField: "out", Tmpl: `{{index . "nonexistent"}}`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("msg", "hi")
	out := tr.Apply(in)
	if _, ok := out["out"]; !ok {
		t.Fatal("expected out field to be set")
	}
}

func TestInvalidTemplateSyntaxReturnsError(t *testing.T) {
	_, err := fieldtemplate.WithRules([]fieldtemplate.Rule{
		{DestField: "bad", Tmpl: `{{.Unclosed`},
	})
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestMultipleRules(t *testing.T) {
	tr, err := fieldtemplate.WithRules([]fieldtemplate.Rule{
		{DestField: "a", Tmpl: `A`},
		{DestField: "b", Tmpl: `B`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := tr.Apply(entry("x", "1"))
	if out["a"] != "A" || out["b"] != "B" {
		t.Fatalf("unexpected output: %v", out)
	}
}
