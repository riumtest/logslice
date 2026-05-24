package fieldmatch_test

import (
	"regexp"
	"testing"

	"github.com/your-org/logslice/internal/fieldmatch"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldmatch.New()
	in := entry("msg", "hello")
	out := tr.Apply(in)
	if out["msg"] != "hello" || len(out) != 1 {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestMatchTrue(t *testing.T) {
	tr := fieldmatch.WithRules([]fieldmatch.Rule{
		{Field: "url", Pattern: regexp.MustCompile(`^/api/`), Dest: "is_api"},
	})
	out := tr.Apply(entry("url", "/api/v1/users"))
	if out["is_api"] != true {
		t.Fatalf("expected true, got %v", out["is_api"])
	}
}

func TestMatchFalse(t *testing.T) {
	tr := fieldmatch.WithRules([]fieldmatch.Rule{
		{Field: "url", Pattern: regexp.MustCompile(`^/api/`), Dest: "is_api"},
	})
	out := tr.Apply(entry("url", "/static/logo.png"))
	if out["is_api"] != false {
		t.Fatalf("expected false, got %v", out["is_api"])
	}
}

func TestMatchMissingFieldSetsFalse(t *testing.T) {
	tr := fieldmatch.WithRules([]fieldmatch.Rule{
		{Field: "url", Pattern: regexp.MustCompile(`.*`), Dest: "matched"},
	})
	out := tr.Apply(entry("msg", "no url here"))
	if out["matched"] != false {
		t.Fatalf("expected false for missing field, got %v", out["matched"])
	}
}

func TestMatchNonStringFieldSetsFalse(t *testing.T) {
	tr := fieldmatch.WithRules([]fieldmatch.Rule{
		{Field: "code", Pattern: regexp.MustCompile(`200`), Dest: "is200"},
	})
	out := tr.Apply(entry("code", 200))
	if out["is200"] != false {
		t.Fatalf("expected false for non-string field, got %v", out["is200"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldmatch.WithRules([]fieldmatch.Rule{
		{Field: "msg", Pattern: regexp.MustCompile(`error`), Dest: "is_error"},
	})
	in := entry("msg", "some error occurred")
	_ = tr.Apply(in)
	if _, ok := in["is_error"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestMultipleRules(t *testing.T) {
	tr := fieldmatch.WithRules([]fieldmatch.Rule{
		{Field: "msg", Pattern: regexp.MustCompile(`error`), Dest: "is_error"},
		{Field: "msg", Pattern: regexp.MustCompile(`timeout`), Dest: "is_timeout"},
	})
	out := tr.Apply(entry("msg", "connection timeout error"))
	if out["is_error"] != true || out["is_timeout"] != true {
		t.Fatalf("unexpected results: %v", out)
	}
}
