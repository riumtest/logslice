package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pipeline"
)

func uniqLines(recs []map[string]any) string {
	var sb strings.Builder
	for _, r := range recs {
		sb.WriteString(toJSON(r) + "\n")
	}
	return sb.String()
}

func TestPipelineFieldUniqDeduplicates(t *testing.T) {
	input := uniqLines([]map[string]any{
		{"msg": "first", "tags": []any{"a", "b", "a"}},
		{"msg": "second", "tags": []any{"x", "x", "x"}},
	})

	out, err := runPipeline(t, input, &pipeline.Config{
		UniqFields: []string{"tags"},
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	lines := nonEmptyLines(out)
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}

	if strings.Count(lines[0], `"a"`) != 1 {
		t.Errorf("first line should contain 'a' exactly once: %s", lines[0])
	}
	if strings.Count(lines[1], `"x"`) != 1 {
		t.Errorf("second line should contain 'x' exactly once: %s", lines[1])
	}
}

func TestPipelineNoUniqFieldsPassesThrough(t *testing.T) {
	input := uniqLines([]map[string]any{
		{"msg": "hello", "tags": []any{"a", "a", "b"}},
	})

	out, err := runPipeline(t, input, &pipeline.Config{})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}

	lines := nonEmptyLines(out)
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(lines))
	}
	// With no uniq config the duplicates should still be present.
	if strings.Count(lines[0], `"a"`) != 2 {
		t.Errorf("expected duplicate 'a' values when uniq is disabled: %s", lines[0])
	}
}
