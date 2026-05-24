package fieldabsent_test

import (
	"testing"

	"github.com/humanlogio/logslice/internal/fieldabsent"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldabsent.New(nil)
	in := entry("level", "info", "msg", "hello")
	out := tr.Transform(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
}

func TestRemovesSingleField(t *testing.T) {
	tr := fieldabsent.New([]fieldabsent.Rule{{Field: "secret"}})
	in := entry("level", "info", "secret", "token123", "msg", "hi")
	out := tr.Transform(in)
	if _, ok := out["secret"]; ok {
		t.Fatal("expected 'secret' to be removed")
	}
	if out["msg"] != "hi" {
		t.Fatalf("expected msg=hi, got %v", out["msg"])
	}
}

func TestRemovesMultipleFields(t *testing.T) {
	tr := fieldabsent.New([]fieldabsent.Rule{
		{Field: "password"},
		{Field: "token"},
	})
	in := entry("user", "alice", "password", "s3cr3t", "token", "abc", "action", "login")
	out := tr.Transform(in)
	if _, ok := out["password"]; ok {
		t.Error("expected 'password' to be removed")
	}
	if _, ok := out["token"]; ok {
		t.Error("expected 'token' to be removed")
	}
	if out["user"] != "alice" {
		t.Errorf("expected user=alice, got %v", out["user"])
	}
}

func TestMissingFieldIsNoop(t *testing.T) {
	tr := fieldabsent.New([]fieldabsent.Rule{{Field: "nonexistent"}})
	in := entry("level", "debug", "msg", "ok")
	out := tr.Transform(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldabsent.New([]fieldabsent.Rule{{Field: "drop_me"}})
	in := entry("keep", "yes", "drop_me", "bye")
	_ = tr.Transform(in)
	if _, ok := in["drop_me"]; !ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestEmptyFieldNameSkipped(t *testing.T) {
	tr := fieldabsent.New([]fieldabsent.Rule{{Field: ""}})
	in := entry("level", "info", "msg", "test")
	out := tr.Transform(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
}
