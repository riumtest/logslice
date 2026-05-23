package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/logslice/internal/pipeline"
)

func scopeLines(records []map[string]any) []string {
	var lines []string
	for _, r := range records {
		b, _ := json.Marshal(r)
		lines = append(lines, string(b))
	}
	return lines
}

func TestPipelineFieldScopePromotes(t *testing.T) {
	t.Parallel()
	records := []map[string]any{
		{"http": map[string]any{"method": "GET", "status": float64(200)}},
		{"http": map[string]any{"method": "POST", "status": float64(201)}},
	}
	input := strings.NewReader(strings.Join(scopeLines(records), "\n"))
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:    []string{"-"},
		FieldScope: []string{"http.method:method", "http.status:status"},
	}, input, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &out); err != nil {
		t.Fatalf("invalid JSON in output: %v", err)
	}
	if out["method"] != "GET" {
		t.Errorf("expected method=GET, got %v", out["method"])
	}
	if out["status"] != float64(200) {
		t.Errorf("expected status=200, got %v", out["status"])
	}
}

func TestPipelineNoFieldScopePassesThrough(t *testing.T) {
	t.Parallel()
	records := []map[string]any{
		{"level": "info", "msg": "ok"},
	}
	input := strings.NewReader(strings.Join(scopeLines(records), "\n"))
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:    []string{"-"},
		FieldScope: nil,
	}, input, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(lines))
	}
}
