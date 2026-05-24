package fieldhasher_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/logslice/logslice/internal/fieldhasher"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func sha256hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func TestIdentityNoRules(t *testing.T) {
	h := fieldhasher.New(nil)
	in := entry("msg", "hello")
	out := h.Apply(in)
	if out["msg"] != "hello" {
		t.Fatalf("unexpected mutation: %v", out)
	}
}

func TestHashSingleField(t *testing.T) {
	h := fieldhasher.New([]fieldhasher.Rule{
		{Fields: []string{"user"}, Dest: "user_hash"},
	})
	in := entry("user", "alice")
	out := h.Apply(in)
	want := sha256hex("alice")
	if out["user_hash"] != want {
		t.Fatalf("got %v, want %v", out["user_hash"], want)
	}
	// original field preserved
	if out["user"] != "alice" {
		t.Fatal("source field was removed")
	}
}

func TestHashMultipleFields(t *testing.T) {
	h := fieldhasher.New([]fieldhasher.Rule{
		{Fields: []string{"a", "b"}, Dest: "ab_hash"},
	})
	in := entry("a", "foo", "b", "bar")
	out := h.Apply(in)
	want := sha256hex("foobar")
	if out["ab_hash"] != want {
		t.Fatalf("got %v, want %v", out["ab_hash"], want)
	}
}

func TestHashMD5Algorithm(t *testing.T) {
	h := fieldhasher.New([]fieldhasher.Rule{
		{Fields: []string{"ip"}, Dest: "ip_md5", Algo: fieldhasher.MD5},
	})
	in := entry("ip", "127.0.0.1")
	out := h.Apply(in)
	if out["ip_md5"] == "" {
		t.Fatal("expected non-empty md5 hash")
	}
}

func TestHashMissingFieldSkipped(t *testing.T) {
	h := fieldhasher.New([]fieldhasher.Rule{
		{Fields: []string{"missing"}, Dest: "h"},
	})
	in := entry("other", "val")
	out := h.Apply(in)
	// dest should still be written (hash of empty input)
	if _, ok := out["h"]; !ok {
		t.Fatal("expected dest field to be present even when source missing")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	h := fieldhasher.New([]fieldhasher.Rule{
		{Fields: []string{"x"}, Dest: "x_hash"},
	})
	in := entry("x", "val")
	h.Apply(in)
	if _, ok := in["x_hash"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestHashNumericValue(t *testing.T) {
	h := fieldhasher.New([]fieldhasher.Rule{
		{Fields: []string{"count"}, Dest: "count_hash"},
	})
	in := entry("count", 42)
	out := h.Apply(in)
	want := sha256hex(fmt.Sprintf("%v", 42))
	if out["count_hash"] != want {
		t.Fatalf("got %v, want %v", out["count_hash"], want)
	}
}
