package fieldmask_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldmask"
)

func entry() map[string]any {
	return map[string]any{
		"level":   "info",
		"msg":     "hello",
		"service": "api",
		"ts":      "2024-01-01T00:00:00Z",
	}
}

func TestIdentityMask(t *testing.T) {
	m := fieldmask.New()
	if !m.IsIdentity() {
		t.Fatal("expected identity mask")
	}
	out := m.Apply(entry())
	if len(out) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(out))
	}
}

func TestIncludeFields(t *testing.T) {
	m := fieldmask.New(fieldmask.WithInclude([]string{"level", "msg"}))
	if m.IsIdentity() {
		t.Fatal("should not be identity")
	}
	out := m.Apply(entry())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if _, ok := out["level"]; !ok {
		t.Error("expected 'level' in output")
	}
	if _, ok := out["msg"]; !ok {
		t.Error("expected 'msg' in output")
	}
}

func TestExcludeFields(t *testing.T) {
	m := fieldmask.New(fieldmask.WithExclude([]string{"ts", "service"}))
	out := m.Apply(entry())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if _, ok := out["ts"]; ok {
		t.Error("'ts' should be excluded")
	}
	if _, ok := out["service"]; ok {
		t.Error("'service' should be excluded")
	}
}

func TestIncludeOverridesExclude(t *testing.T) {
	// include wins: only 'level' and 'msg' kept; exclude of 'level' still removes it
	m := fieldmask.New(
		fieldmask.WithInclude([]string{"level", "msg"}),
		fieldmask.WithExclude([]string{"level"}),
	)
	out := m.Apply(entry())
	if len(out) != 1 {
		t.Fatalf("expected 1 field, got %d", len(out))
	}
	if _, ok := out["msg"]; !ok {
		t.Error("expected 'msg' in output")
	}
}

func TestEmptyIncludeSliceIsNoop(t *testing.T) {
	m := fieldmask.New(fieldmask.WithInclude([]string{}))
	if !m.IsIdentity() {
		t.Fatal("empty include should be identity")
	}
}

func TestBlankFieldNamesIgnored(t *testing.T) {
	m := fieldmask.New(fieldmask.WithExclude([]string{"  ", ""}))
	if !m.IsIdentity() {
		t.Fatal("blank field names should be ignored")
	}
}
