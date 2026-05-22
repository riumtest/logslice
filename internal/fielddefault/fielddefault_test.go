package fielddefault_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logslice/internal/fielddefault"
)

func entry(fields map[string]any) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(fields))
	for k, v := range fields {
		raw, _ := json.Marshal(v)
		out[k] = json.RawMessage(raw)
	}
	return out
}

func TestIdentityNoRules(t *testing.T) {
	tr := fielddefault.New(nil)
	in := entry(map[string]any{"level": "info"})
	out := tr.Transform(in)
	if string(out["level"]) != `"info"` {
		t.Fatalf("expected info, got %s", out["level"])
	}
}

func TestDefaultMissingField(t *testing.T) {
	tr := fielddefault.New([]fielddefault.Rule{
		{Field: "env", Value: "production"},
	})
	in := entry(map[string]any{"level": "info"})
	out := tr.Transform(in)
	if string(out["env"]) != `"production"` {
		t.Fatalf("expected production, got %s", out["env"])
	}
}

func TestDefaultDoesNotOverwriteExisting(t *testing.T) {
	tr := fielddefault.New([]fielddefault.Rule{
		{Field: "env", Value: "production"},
	})
	in := entry(map[string]any{"env": "staging"})
	out := tr.Transform(in)
	if string(out["env"]) != `"staging"` {
		t.Fatalf("expected staging, got %s", out["env"])
	}
}

func TestDefaultReplacesNull(t *testing.T) {
	tr := fielddefault.New([]fielddefault.Rule{
		{Field: "region", Value: "us-east-1"},
	})
	in := entry(map[string]any{"region": nil})
	out := tr.Transform(in)
	if string(out["region"]) != `"us-east-1"` {
		t.Fatalf("expected us-east-1, got %s", out["region"])
	}
}

func TestDefaultNumericValue(t *testing.T) {
	tr := fielddefault.New([]fielddefault.Rule{
		{Field: "retries", Value: 3},
	})
	in := entry(map[string]any{"level": "warn"})
	out := tr.Transform(in)
	if string(out["retries"]) != `3` {
		t.Fatalf("expected 3, got %s", out["retries"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fielddefault.New([]fielddefault.Rule{
		{Field: "env", Value: "production"},
	})
	in := entry(map[string]any{"level": "info"})
	_ = tr.Transform(in)
	if _, ok := in["env"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}
