package fieldtimestamp_test

import (
	"testing"
	"time"

	"github.com/qsocket/logslice/internal/fieldtimestamp"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldtimestamp.New()
	in := entry("ts", "2024-01-02T03:04:05Z")
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ts"] != "2024-01-02T03:04:05Z" {
		t.Fatalf("expected unchanged value, got %v", out["ts"])
	}
}

func TestReformatRFC3339(t *testing.T) {
	rules := []fieldtimestamp.Rule{
		{Field: "ts", InputLayout: time.RFC3339, OutputLayout: "2006/01/02"},
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	in := entry("ts", "2024-06-15T10:00:00Z")
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ts"] != "2024/06/15" {
		t.Fatalf("expected 2024/06/15, got %v", out["ts"])
	}
}

func TestWriteToDestField(t *testing.T) {
	rules := []fieldtimestamp.Rule{
		{Field: "ts", OutputLayout: "2006-01-02", DestField: "date"},
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	in := entry("ts", "2024-03-10T08:30:00Z")
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["date"] != "2024-03-10" {
		t.Fatalf("expected date=2024-03-10, got %v", out["date"])
	}
	// Original field must remain.
	if out["ts"] != "2024-03-10T08:30:00Z" {
		t.Fatalf("original ts should be preserved, got %v", out["ts"])
	}
}

func TestSkipsMissingField(t *testing.T) {
	rules := []fieldtimestamp.Rule{
		{Field: "ts", OutputLayout: "2006-01-02"},
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	in := entry("msg", "hello")
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["ts"]; ok {
		t.Fatal("ts should not be present")
	}
}

func TestSkipsNonStringField(t *testing.T) {
	rules := []fieldtimestamp.Rule{
		{Field: "ts", OutputLayout: "2006-01-02"},
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	in := entry("ts", 12345)
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ts"] != 12345 {
		t.Fatalf("non-string value should be unchanged, got %v", out["ts"])
	}
}

func TestParseError(t *testing.T) {
	rules := []fieldtimestamp.Rule{
		{Field: "ts", InputLayout: time.RFC3339},
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	in := entry("ts", "not-a-timestamp")
	_, err := tr.Apply(in)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	rules := []fieldtimestamp.Rule{
		{Field: "ts", OutputLayout: "2006-01-02"},
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	in := entry("ts", "2024-01-01T00:00:00Z")
	orig := in["ts"]
	_, _ = tr.Apply(in)
	if in["ts"] != orig {
		t.Fatal("original entry was mutated")
	}
}
