package fielduniq_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fielduniq"
)

func entry(fields map[string]any) map[string]any { return fields }

func TestIdentityNoRules(t *testing.T) {
	tr := fielduniq.New(nil)
	in := entry(map[string]any{"tags": []any{"a", "b", "a"}})
	out := tr.Apply(in)
	if len(out["tags"].([]any)) != 3 {
		t.Fatal("expected no dedup when no fields configured")
	}
}

func TestDeduplicatesStrings(t *testing.T) {
	tr := fielduniq.New([]string{"tags"})
	in := entry(map[string]any{"tags": []any{"a", "b", "a", "c", "b"}})
	out := tr.Apply(in)
	got := out["tags"].([]any)
	want := []any{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("index %d: got %v want %v", i, got[i], v)
		}
	}
}

func TestNoDuplicatesUnchanged(t *testing.T) {
	tr := fielduniq.New([]string{"tags"})
	in := entry(map[string]any{"tags": []any{"x", "y", "z"}})
	out := tr.Apply(in)
	got := out["tags"].([]any)
	if len(got) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(got))
	}
}

func TestNonArrayFieldSkipped(t *testing.T) {
	tr := fielduniq.New([]string{"level"})
	in := entry(map[string]any{"level": "info"})
	out := tr.Apply(in)
	if out["level"] != "info" {
		t.Fatal("non-array field should be left untouched")
	}
}

func TestMissingFieldSkipped(t *testing.T) {
	tr := fielduniq.New([]string{"tags"})
	in := entry(map[string]any{"msg": "hello"})
	out := tr.Apply(in)
	if _, ok := out["tags"]; ok {
		t.Fatal("missing field should not be added")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fielduniq.New([]string{"tags"})
	orig := []any{"a", "a", "b"}
	in := entry(map[string]any{"tags": orig})
	tr.Apply(in)
	if len(orig) != 3 {
		t.Fatal("original slice should not be mutated")
	}
}

func TestMultipleFields(t *testing.T) {
	tr := fielduniq.New([]string{"tags", "labels"})
	in := entry(map[string]any{
		"tags":   []any{"a", "a"},
		"labels": []any{"x", "y", "x"},
	})
	out := tr.Apply(in)
	if len(out["tags"].([]any)) != 1 {
		t.Error("tags not deduped")
	}
	if len(out["labels"].([]any)) != 2 {
		t.Error("labels not deduped")
	}
}
