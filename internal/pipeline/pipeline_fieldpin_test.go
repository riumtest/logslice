package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pipeline"
)

func pinLines() string {
	lines := []map[string]any{
		{"level": "info", "msg": "started"},
		{"level": "error", "msg": "failed"},
		{"level": "warn", "msg": "retrying"},
	}
	var sb strings.Builder
	for _, l := range lines {
		b, _ := json.Marshal(l)
		sb.Write(b)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func TestPipelineFieldPinAnnotates(t *testing.T) {
	src := strings.NewReader(pinLines())
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.Source{{Reader: src}},
		Output:  &buf,
		FieldPinRules: []pipeline.FieldPinRule{
			{Field: "level", Dest: "_pin"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 3 {
		t.Fatalf("expected 3 output lines, got %d", len(lines))
	}

	for _, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		pin, ok := entry["_pin"].(map[string]any)
		if !ok {
			t.Fatalf("expected _pin map in: %s", line)
		}
		if _, ok := pin["level"]; !ok {
			t.Fatalf("expected level in _pin for: %s", line)
		}
	}
}

func TestPipelineNoFieldPinPassesThrough(t *testing.T) {
	src := strings.NewReader(pinLines())
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:       []pipeline.Source{{Reader: src}},
		Output:        &buf,
		FieldPinRules: nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	for _, line := range lines {
		if strings.Contains(line, "_pin") {
			t.Fatalf("unexpected _pin field in: %s", line)
		}
	}
}
