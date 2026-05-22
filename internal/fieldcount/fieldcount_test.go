package fieldcount_test

import (
	"encoding/json"
	"testing"

	"github.com/yourorg/logslice/internal/fieldcount"
)

func entry(pairs ...string) map[string]json.RawMessage {
	m := make(map[string]json.RawMessage, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = json.RawMessage(`"` + pairs[i+1] + `"`)
	}
	return m
}

func TestDefaultDestField(t *testing.T) {
	tr := fieldcount.New()
	out := tr.Transform(entry("level", "info", "msg", "hello"))

	raw, ok := out["_field_count"]
	if !ok {
		t.Fatal("expected _field_count field to be present")
	}
	var n int
	if err := json.Unmarshal(raw, &n); err != nil {
		t.Fatalf("unmarshal count: %v", err)
	}
	if n != 2 {
		t.Errorf("expected count 2, got %d", n)
	}
}

func TestCustomDestField(t *testing.T) {
	tr := fieldcount.New(fieldcount.WithDestField("num_fields"))
	out := tr.Transform(entry("a", "1", "b", "2", "c", "3"))

	if _, ok := out["_field_count"]; ok {
		t.Error("default field should not be present when custom field is set")
	}
	raw, ok := out["num_fields"]
	if !ok {
		t.Fatal("expected num_fields to be present")
	}
	var n int
	if err := json.Unmarshal(raw, &n); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestEmptyEntryUnchanged(t *testing.T) {
	tr := fieldcount.New()
	out := tr.Transform(map[string]json.RawMessage{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldcount.New()
	orig := entry("x", "1")
	_ = tr.Transform(orig)
	if _, ok := orig["_field_count"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestWithDestFieldIgnoresEmpty(t *testing.T) {
	tr := fieldcount.New(fieldcount.WithDestField(""))
	out := tr.Transform(entry("k", "v"))
	if _, ok := out["_field_count"]; !ok {
		t.Error("expected default field name when empty string passed to WithDestField")
	}
}
