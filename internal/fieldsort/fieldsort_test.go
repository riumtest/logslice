package fieldsort_test

import (
	"testing"

	"github.com/yourusername/logslice/internal/fieldsort"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldsort.New()
	in := entry("a", 1, "b", 2)
	out := tr.Apply(in)
	if len(out) != len(in) {
		t.Fatalf("expected %d keys, got %d", len(in), len(out))
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldsort.New(fieldsort.WithRules([]fieldsort.Rule{{Priority: []string{"z", "a"}}}))
	in := entry("a", 1, "b", 2, "z", 3)
	_ = tr.Apply(in)
	if len(in) != 3 {
		t.Fatal("original entry was mutated")
	}
}

func TestOrderedPriorityFirst(t *testing.T) {
	tr := fieldsort.New(fieldsort.WithRules([]fieldsort.Rule{
		{Priority: []string{"level", "msg", "time"}},
	}))
	in := entry("time", "now", "msg", "hello", "level", "info", "svc", "api")
	keys := tr.Ordered(in)

	if keys[0] != "level" || keys[1] != "msg" || keys[2] != "time" {
		t.Fatalf("unexpected order: %v", keys)
	}
	// remaining key should be present
	if keys[3] != "svc" {
		t.Fatalf("expected svc at index 3, got %s", keys[3])
	}
}

func TestOrderedSkipsMissingPriorityKeys(t *testing.T) {
	tr := fieldsort.New(fieldsort.WithRules([]fieldsort.Rule{
		{Priority: []string{"level", "missing", "msg"}},
	}))
	in := entry("msg", "hi", "level", "warn")
	keys := tr.Ordered(in)

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "level" || keys[1] != "msg" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestOrderedNoRulesReturnsSorted(t *testing.T) {
	tr := fieldsort.New()
	in := entry("z", 1, "a", 2, "m", 3)
	keys := tr.Ordered(in)
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Fatalf("expected alphabetical order, got %v", keys)
	}
}

func TestOrderedEmptyEntry(t *testing.T) {
	tr := fieldsort.New(fieldsort.WithRules([]fieldsort.Rule{{Priority: []string{"a"}}}))  
	keys := tr.Ordered(map[string]any{})
	if len(keys) != 0 {
		t.Fatalf("expected empty slice, got %v", keys)
	}
}
