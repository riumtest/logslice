package fieldredact_test

import (
	"encoding/json"
	"testing"

	"github.com/yourusername/logslice/internal/fieldredact"
)

func entry(kv ...string) map[string]json.RawMessage {
	m := make(map[string]json.RawMessage, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		v, _ := json.Marshal(kv[i+1])
		m[kv[i]] = v
	}
	return m
}

func str(raw json.RawMessage) string {
	var s string
	_ = json.Unmarshal(raw, &s)
	return s
}

func TestIdentityRedact(t *testing.T) {
	r := fieldredact.New(nil)
	e := entry("msg", "hello", "user", "alice")
	out := r.Apply(e)
	if str(out["user"]) != "alice" {
		t.Fatalf("expected alice, got %s", str(out["user"]))
	}
}

func TestRedactSingleField(t *testing.T) {
	r := fieldredact.New([]string{"password"})
	e := entry("msg", "login", "password", "s3cr3t")
	out := r.Apply(e)
	if str(out["password"]) != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %s", str(out["password"]))
	}
	if str(out["msg"]) != "login" {
		t.Fatalf("msg should be unchanged")
	}
}

func TestRedactMultipleFields(t *testing.T) {
	r := fieldredact.New([]string{"token", "secret"})
	e := entry("token", "abc", "secret", "xyz", "level", "info")
	out := r.Apply(e)
	if str(out["token"]) != "[REDACTED]" {
		t.Fatalf("token not redacted")
	}
	if str(out["secret"]) != "[REDACTED]" {
		t.Fatalf("secret not redacted")
	}
	if str(out["level"]) != "info" {
		t.Fatalf("level should be unchanged")
	}
}

func TestRedactMissingFieldIgnored(t *testing.T) {
	r := fieldredact.New([]string{"password"})
	e := entry("msg", "hello")
	out := r.Apply(e)
	if _, ok := out["password"]; ok {
		t.Fatal("missing field should not be added")
	}
}

func TestRedactCustomMask(t *testing.T) {
	r := fieldredact.New([]string{"ssn"}, fieldredact.WithMask("***"))
	e := entry("ssn", "123-45-6789")
	out := r.Apply(e)
	if str(out["ssn"]) != "***" {
		t.Fatalf("expected ***, got %s", str(out["ssn"]))
	}
}

func TestRedactEmptyMaskIgnored(t *testing.T) {
	r := fieldredact.New([]string{"key"}, fieldredact.WithMask(""))
	e := entry("key", "value")
	out := r.Apply(e)
	if str(out["key"]) != "[REDACTED]" {
		t.Fatalf("empty mask should fall back to default, got %s", str(out["key"]))
	}
}
