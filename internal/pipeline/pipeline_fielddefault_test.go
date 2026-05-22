package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pipeline"
)

func defaultLines() string {
	return strings.Join([]string{
		`{"level":"info","msg":"started"}`,
		`{"level":"warn","msg":"retry","env":"staging"}`,
		`{"level":"error","msg":"failed","env":null}`,
	}, "\n")
}

func TestPipelineFieldDefault(t *testing.T) {
	src := strings.NewReader(defaultLines())
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.Source{{Reader: src, Name: "test"}},
		Output:  &buf,
		FieldDefaults: []pipeline.FieldDefault{
			{Field: "env", Value: "production"},
		},
		Format: "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// First line has no env field — default should be applied.
	if !strings.Contains(lines[0], `"production"`) {
		t.Errorf("line 0: expected default env=production, got %s", lines[0])
	}
	// Second line already has env=staging — should be unchanged.
	if !strings.Contains(lines[1], `"staging"`) {
		t.Errorf("line 1: expected env=staging to be preserved, got %s", lines[1])
	}
	// Third line has env=null — default should replace it.
	if !strings.Contains(lines[2], `"production"`) {
		t.Errorf("line 2: expected null env replaced with production, got %s", lines[2])
	}
}

func TestPipelineNoFieldDefaults(t *testing.T) {
	src := strings.NewReader(`{"level":"info","msg":"ok"}`)
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.Source{{Reader: src, Name: "test"}},
		Output:  &buf,
		Format:  "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"ok"`) {
		t.Errorf("expected msg=ok in output, got %s", buf.String())
	}
}
