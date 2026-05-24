package fieldencoding_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldencoding"
)

func entry(kvs ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(kvs); i += 2 {
		m[kvs[i].(string)] = kvs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tf := fieldencoding.New()
	in := entry("msg", "hello")
	out, err := tf.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["msg"] != "hello" {
		t.Fatalf("expected hello, got %v", out["msg"])
	}
}

func TestBase64Encode(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "data", Mode: "base64-encode"},
	})
	out, err := tf.Apply(entry("data", "hello"))
	if err != nil {
		t.Fatal(err)
	}
	if out["data"] != "aGVsbG8=" {
		t.Fatalf("unexpected: %v", out["data"])
	}
}

func TestBase64Decode(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "data", Mode: "base64-decode"},
	})
	out, err := tf.Apply(entry("data", "aGVsbG8="))
	if err != nil {
		t.Fatal(err)
	}
	if out["data"] != "hello" {
		t.Fatalf("unexpected: %v", out["data"])
	}
}

func TestHexEncode(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "id", Mode: "hex-encode"},
	})
	out, err := tf.Apply(entry("id", "hi"))
	if err != nil {
		t.Fatal(err)
	}
	if out["id"] != "6869" {
		t.Fatalf("unexpected: %v", out["id"])
	}
}

func TestHexDecode(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "id", Mode: "hex-decode"},
	})
	out, err := tf.Apply(entry("id", "6869"))
	if err != nil {
		t.Fatal(err)
	}
	if out["id"] != "hi" {
		t.Fatalf("unexpected: %v", out["id"])
	}
}

func TestWritesToDestField(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "raw", Mode: "base64-encode", Dest: "encoded"},
	})
	out, err := tf.Apply(entry("raw", "abc"))
	if err != nil {
		t.Fatal(err)
	}
	if out["raw"] != "abc" {
		t.Fatalf("original should be unchanged")
	}
	if out["encoded"] != "YWJj" {
		t.Fatalf("unexpected encoded: %v", out["encoded"])
	}
}

func TestUnknownModeReturnsError(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "data", Mode: "rot13"},
	})
	_, err := tf.Apply(entry("data", "hello"))
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tf := fieldencoding.WithRules([]fieldencoding.Rule{
		{Field: "msg", Mode: "hex-encode"},
	})
	in := entry("msg", "hi")
	_, err := tf.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if in["msg"] != "hi" {
		t.Fatal("original entry was mutated")
	}
}
