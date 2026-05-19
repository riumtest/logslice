package dedupe_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/dedupe"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestNoDuplicates(t *testing.T) {
	d := dedupe.New()
	a := entry("msg", "hello", "level", "info")
	b := entry("msg", "world", "level", "info")
	if d.IsDuplicate(a) {
		t.Fatal("first entry should never be a duplicate")
	}
	if d.IsDuplicate(b) {
		t.Fatal("different entry should not be a duplicate")
	}
}

func TestConsecutiveDuplicate(t *testing.T) {
	d := dedupe.New()
	a := entry("msg", "hello", "level", "error")
	if d.IsDuplicate(a) {
		t.Fatal("first entry should not be a duplicate")
	}
	if !d.IsDuplicate(a) {
		t.Fatal("identical consecutive entry should be a duplicate")
	}
}

func TestNonConsecutiveRepeatAllowed(t *testing.T) {
	d := dedupe.New()
	a := entry("msg", "hello")
	b := entry("msg", "world")
	d.IsDuplicate(a) // prime
	d.IsDuplicate(b) // different
	if d.IsDuplicate(a) {
		t.Fatal("non-consecutive repeat should not be treated as duplicate")
	}
}

func TestWithFieldsSubset(t *testing.T) {
	d := dedupe.New(dedupe.WithFields("msg"))
	a := entry("msg", "hello", "ts", "2024-01-01T00:00:00Z")
	b := entry("msg", "hello", "ts", "2024-01-02T00:00:00Z") // same msg, different ts
	if d.IsDuplicate(a) {
		t.Fatal("first entry should not be duplicate")
	}
	if !d.IsDuplicate(b) {
		t.Fatal("entry with same msg field should be duplicate when only msg is tracked")
	}
}

func TestReset(t *testing.T) {
	d := dedupe.New()
	a := entry("msg", "hello")
	d.IsDuplicate(a)
	d.Reset()
	if d.IsDuplicate(a) {
		t.Fatal("entry after Reset should not be a duplicate")
	}
}

func TestEmptyEntry(t *testing.T) {
	d := dedupe.New()
	e := entry()
	if d.IsDuplicate(e) {
		t.Fatal("first empty entry should not be a duplicate")
	}
	if !d.IsDuplicate(e) {
		t.Fatal("second empty entry should be a duplicate")
	}
}
