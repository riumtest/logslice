package fieldtransform_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldtransform"
)

func entry(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestIdentityTransform(t *testing.T) {
	tr, err := fieldtransform.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("msg", "hello", "level", "info")
	out := tr.Apply(in)
	if out["msg"] != "hello" || out["level"] != "info" {
		t.Fatalf("identity transform mutated entry: %v", out)
	}
}

func TestRenameField(t *testing.T) {
	tr, err := fieldtransform.New(
		fieldtransform.WithRename("message", "msg"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("message", "hello", "level", "info")
	out := tr.Apply(in)
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
	if _, ok := out["message"]; ok {
		t.Error("old key 'message' should not be present after rename")
	}
}

func TestRenameSkipsConflict(t *testing.T) {
	// If the target key already exists in the entry, the rename is skipped.
	tr, err := fieldtransform.New(
		fieldtransform.WithRename("message", "msg"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("message", "from-message", "msg", "existing")
	out := tr.Apply(in)
	if out["msg"] != "existing" {
		t.Errorf("expected existing msg to be preserved, got %v", out["msg"])
	}
	if out["message"] != "from-message" {
		t.Errorf("expected original message to remain, got %v", out["message"])
	}
}

func TestRenameConflictDetected(t *testing.T) {
	_, err := fieldtransform.New(
		fieldtransform.WithRename("ts", "time"),
		fieldtransform.WithRename("timestamp", "time"),
	)
	if err == nil {
		t.Fatal("expected error for conflicting rename targets, got nil")
	}
}

func TestWithRenameEmptyFieldIgnored(t *testing.T) {
	tr, err := fieldtransform.New(
		fieldtransform.WithRename("", "msg"),
		fieldtransform.WithRename("level", ""),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	in := entry("level", "warn")
	out := tr.Apply(in)
	if out["level"] != "warn" {
		t.Errorf("level should be unchanged, got %v", out["level"])
	}
}
