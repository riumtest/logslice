package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/fieldenv"
	"github.com/yourorg/logslice/internal/pipeline"
)

func envLines() []string {
	return []string{
		`{"msg":"hello"}`,
		`{"msg":"world","environment":"existing"}`,
	}
}

func TestPipelineFieldEnvInjects(t *testing.T) {
	fakeEnv := func(key string) (string, bool) {
		if key == "APP_ENV" {
			return "production", true
		}
		return "", false
	}
	rules := []fieldenv.Rule{{Env: "APP_ENV", Dest: "environment", Default: "dev"}}
	tr := fieldenv.New(rules, fieldenv.WithLookup(fakeEnv))

	lines := envLines()
	var buf strings.Builder
	err := pipeline.Run(pipeline.Config{
		Sources:          readerFromLines(lines),
		Output:           &buf,
		FieldEnvTransform: tr,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected 'production' in output, got: %s", out)
	}
	// second record already had the field; should remain unchanged
	if strings.Count(out, "existing") != 1 {
		t.Errorf("expected 'existing' to appear once, got: %s", out)
	}
}

func TestPipelineNoFieldEnvPassesThrough(t *testing.T) {
	lines := envLines()
	var buf strings.Builder
	err := pipeline.Run(pipeline.Config{
		Sources: readerFromLines(lines),
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}
