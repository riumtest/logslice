package fieldlookup_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldlookup"
	"github.com/logslice/logslice/internal/pipeline"
)

func entry(kv ...any) pipeline.Entry {
	e := make(pipeline.Entry)
	for i := 0; i+1 < len(kv); i += 2 {
		e[kv[i].(string)] = kv[i+1]
	}
	return e
}

func TestIdentityLookup(t *testing.T) {
	tr := fieldlookup.New("env", map[string]string{})
	in := entry("env", "prod", "msg", "hello")
	out := tr.Transform(in)
	if out["env"] != "prod" {
		t.Fatalf("expected prod, got %v", out["env"])
	}
}

func TestLookupReplaces(t *testing.T) {
	tr := fieldlookup.New("env", map[string]string{"prod": "production", "dev": "development"})
	out := tr.Transform(entry("env", "prod"))
	if out["env"] != "production" {
		t.Fatalf("expected production, got %v", out["env"])
	}
}

func TestLookupMissingKeyNoDefault(t *testing.T) {
	tr := fieldlookup.New("env", map[string]string{"prod": "production"})
	out := tr.Transform(entry("env", "staging"))
	if out["env"] != "staging" {
		t.Fatalf("expected original value staging, got %v", out["env"])
	}
}

func TestLookupMissingKeyWithDefault(t *testing.T) {
	tr := fieldlookup.New("env", map[string]string{"prod": "production"},
		fieldlookup.WithDefault("unknown"))
	out := tr.Transform(entry("env", "staging"))
	if out["env"] != "unknown" {
		t.Fatalf("expected unknown, got %v", out["env"])
	}
}

func TestLookupDestField(t *testing.T) {
	tr := fieldlookup.New("code", map[string]string{"200": "OK"},
		fieldlookup.WithDestField("status_text"))
	out := tr.Transform(entry("code", "200"))
	if out["code"] != "200" {
		t.Fatalf("source field should be unchanged, got %v", out["code"])
	}
	if out["status_text"] != "OK" {
		t.Fatalf("expected OK in status_text, got %v", out["status_text"])
	}
}

func TestLookupNonStringFieldUnchanged(t *testing.T) {
	tr := fieldlookup.New("code", map[string]string{"200": "OK"})
	out := tr.Transform(entry("code", 200))
	if out["code"] != 200 {
		t.Fatalf("non-string field should be unchanged, got %v", out["code"])
	}
}

func TestLookupMissingFieldUnchanged(t *testing.T) {
	tr := fieldlookup.New("env", map[string]string{"prod": "production"})
	out := tr.Transform(entry("msg", "hello"))
	if _, ok := out["env"]; ok {
		t.Fatal("env field should not be added when missing from entry")
	}
}

func TestLookupDoesNotMutateOriginal(t *testing.T) {
	tr := fieldlookup.New("env", map[string]string{"prod": "production"})
	in := entry("env", "prod")
	_ = tr.Transform(in)
	if in["env"] != "prod" {
		t.Fatal("original entry must not be mutated")
	}
}

// TestLookupDestFieldWithDefault verifies that when a dest field is configured
// alongside a default value, the default is written to the dest field rather
// than the source field when the lookup key is not found.
func TestLookupDestFieldWithDefault(t *testing.T) {
	tr := fieldlookup.New("code", map[string]string{"200": "OK"},
		fieldlookup.WithDestField("status_text"),
		fieldlookup.WithDefault("UNKNOWN"))
	out := tr.Transform(entry("code", "404"))
	if out["code"] != "404" {
		t.Fatalf("source field should be unchanged, got %v", out["code"])
	}
	if out["status_text"] != "UNKNOWN" {
		t.Fatalf("expected UNKNOWN in status_text, got %v", out["status_text"])
	}
}
