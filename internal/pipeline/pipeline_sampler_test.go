package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/sampler"
)

func jsonLines(records []map[string]any) string {
	var sb strings.Builder
	for _, r := range records {
		b, _ := json.Marshal(r)
		sb.Write(b)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func TestPipelineHeadSampling(t *testing.T) {
	records := make([]map[string]any, 10)
	for i := range records {
		records[i] = map[string]any{"level": "info", "message": "msg", "ts": "2024-01-01T00:00:00Z"}
	}
	src := strings.NewReader(jsonLines(records))

	var buf bytes.Buffer
	w := output.New(output.WithDestination(&buf))
	f := formatter.New()
	s := sampler.New(sampler.ModeHead, 3)

	err := pipeline.Run(pipeline.Config{
		Sources:   []interface{ Read() (map[string]any, error) }{mustReader(src)},
		Filter:    "",
		Formatter: f,
		Writer:    w,
		Sampler:   s,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 sampled lines, got %d: %q", len(lines), buf.String())
	}
}

func TestPipelineTailSampling(t *testing.T) {
	records := make([]map[string]any, 8)
	for i := range records {
		records[i] = map[string]any{"level": "info", "message": "msg", "ts": "2024-01-01T00:00:00Z"}
	}
	src := strings.NewReader(jsonLines(records))

	var buf bytes.Buffer
	w := output.New(output.WithDestination(&buf))
	f := formatter.New()
	s := sampler.New(sampler.ModeTail, 4)

	err := pipeline.Run(pipeline.Config{
		Sources:   []interface{ Read() (map[string]any, error) }{mustReader(src)},
		Filter:    "",
		Formatter: f,
		Writer:    w,
		Sampler:   s,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 tail lines, got %d: %q", len(lines), buf.String())
	}
}
