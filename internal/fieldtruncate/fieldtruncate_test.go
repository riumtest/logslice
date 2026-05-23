package fieldtruncate_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldtruncate"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldtruncate.New()
	in := entry("msg", "hello world")
	out := tr.Apply(in)
	if out["msg"] != "hello world" {
		t.Fatalf("expected unchanged, got %v", out["msg"])
	}
}

func TestTruncateShortValueUnchanged(t *testing.T) {
	tr := fieldtruncate.WithRules([]fieldtruncate.Rule{{Field: "msg", MaxLen: 20}})
	in := entry("msg", "hi")
	out := tr.Apply(in)
	if out["msg"] != "hi" {
		t.Fatalf("expected 'hi', got %v", out["msg"])
	}
}

func TestTruncateLongValue(t *testing.T) {
	tr := fieldtruncate.WithRules([]fieldtruncate.Rule{{Field: "msg", MaxLen: 5}})
	in := entry("msg", "hello world")
	out := tr.Apply(in)
	if out["msg"] != "hello..." {
		t.Fatalf("expected 'hello...', got %v", out["msg"])
	}
}

func TestTruncateCustomSuffix(t *testing.T) {
	tr := fieldtruncate.WithRules([]fieldtruncate.Rule{{Field: "msg", MaxLen: 4, Suffix: "--"}})
	in := entry("msg", "abcdefgh")
	out := tr.Apply(in)
	if out["msg"] != "abcd--" {
		t.Fatalf("expected 'abcd--', got %v", out["msg"])
	}
}

func TestTruncateNonStringSkipped(t *testing.T) {
	tr := fieldtruncate.WithRules([]fieldtruncate.Rule{{Field: "count", MaxLen: 2}})
	in := entry("count", float64(12345))
	out := tr.Apply(in)
	if out["count"] != float64(12345) {
		t.Fatalf("expected numeric unchanged, got %v", out["count"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldtruncate.WithRules([]fieldtruncate.Rule{{Field: "msg", MaxLen: 3}})
	in := entry("msg", "hello world")
	original := in["msg"]
	tr.Apply(in)
	if in["msg"] != original {
		t.Fatal("original entry was mutated")
	}
}

func TestTruncateMultibyteRunes(t *testing.T) {
	// Each emoji is multiple bytes but one rune.
	tr := fieldtruncate.WithRules([]fieldtruncate.Rule{{Field: "msg", MaxLen: 3}})
	in := entry("msg", "😀😁😂😃😄")
	out := tr.Apply(in)
	if out["msg"] != "😀😁😂..." {
		t.Fatalf("expected 3 emoji + suffix, got %v", out["msg"])
	}
}
