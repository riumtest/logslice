package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pipeline"
)

func levelLines(records []map[string]any) []string {
	var lines []string
	for _, r := range records {
		b, _ := json.Marshal(r)
		lines = append(lines, string(b))
	}
	return lines
}

func TestPipelineLevelFilterWarn(t *testing.T) {
	records := []map[string]any{
		{"level": "debug", "msg": "verbose"},
		{"level": "info", "msg": "startup"},
		{"level": "warn", "msg": "low disk"},
		{"level": "error", "msg": "crash"},
	}
	src := strings.NewReader(strings.Join(levelLines(records), "\n"))
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:   []pipeline.Source{{Reader: src, Name: "test"}},
		MinLevel:  "warn",
		LevelField: "level",
		Output:    &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "verbose") {
		t.Error("debug entry should have been filtered")
	}
	if strings.Contains(out, "startup") {
		t.Error("info entry should have been filtered")
	}
	if !strings.Contains(out, "low disk") {
		t.Error("warn entry should be present")
	}
	if !strings.Contains(out, "crash") {
		t.Error("error entry should be present")
	}
}

func TestPipelineNoLevelFilterPassesAll(t *testing.T) {
	records := []map[string]any{
		{"level": "debug", "msg": "a"},
		{"level": "info", "msg": "b"},
	}
	src := strings.NewReader(strings.Join(levelLines(records), "\n"))
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.Source{{Reader: src, Name: "test"}},
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\"a\"") || !strings.Contains(out, "\"b\"") {
		t.Error("all entries should pass when no level filter is set")
	}
}
