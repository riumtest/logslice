package fieldmerge_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldmerge"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestMergeBasic(t *testing.T) {
	m := fieldmerge.New("full", []string{"first", "last"})
	out := m.Apply(entry("first", "Jane", "last", "Doe"))
	if got := out["full"]; got != "Jane Doe" {
		t.Fatalf("expected 'Jane Doe', got %q", got)
	}
}

func TestMergeCustomSeparator(t *testing.T) {
	m := fieldmerge.New("path", []string{"dir", "file"}, fieldmerge.WithSeparator("/"))
	out := m.Apply(entry("dir", "home", "file", "log.txt"))
	if got := out["path"]; got != "home/log.txt" {
		t.Fatalf("expected 'home/log.txt', got %q", got)
	}
}

func TestMergeSkipsMissingByDefault(t *testing.T) {
	m := fieldmerge.New("full", []string{"first", "middle", "last"})
	out := m.Apply(entry("first", "Jane", "last", "Doe"))
	if got := out["full"]; got != "Jane Doe" {
		t.Fatalf("expected 'Jane Doe', got %q", got)
	}
}

func TestMergeShowsMissingPlaceholder(t *testing.T) {
	m := fieldmerge.New("full", []string{"first", "last"}, fieldmerge.WithSkipMissing(false))
	out := m.Apply(entry("first", "Jane"))
	if got := out["full"]; got != "Jane <missing:last>" {
		t.Fatalf("unexpected value: %q", got)
	}
}

func TestMergeDoesNotMutateOriginal(t *testing.T) {
	orig := entry("a", "hello", "b", "world")
	m := fieldmerge.New("merged", []string{"a", "b"})
	_ = m.Apply(orig)
	if _, ok := orig["merged"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestMergePreservesExistingFields(t *testing.T) {
	m := fieldmerge.New("merged", []string{"x", "y"})
	out := m.Apply(entry("x", "1", "y", "2", "z", "keep"))
	if got := out["z"]; got != "keep" {
		t.Fatalf("expected 'keep', got %q", got)
	}
}

func TestMergeNumericValues(t *testing.T) {
	m := fieldmerge.New("label", []string{"code", "msg"}, fieldmerge.WithSeparator(":"))
	out := m.Apply(entry("code", 404, "msg", "not found"))
	if got := out["label"]; got != "404:not found" {
		t.Fatalf("unexpected value: %q", got)
	}
}
