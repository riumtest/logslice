package pipeline_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/pipeline"
)

func makeTimedLines(anchor time.Time) string {
	var sb strings.Builder
	offsets := []time.Duration{-90 * time.Minute, -30 * time.Minute, -5 * time.Minute, 10 * time.Minute}
	for i, off := range offsets {
		rec := map[string]interface{}{
			"level":   "info",
			"msg":     fmt.Sprintf("event %d", i),
			"time":    anchor.Add(off).Format(time.RFC3339),
		}
		b, _ := json.Marshal(rec)
		sb.Write(b)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func TestPipelineTimeRangeFilter(t *testing.T) {
	anchor := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	input := makeTimedLines(anchor)

	// keep only records within the last hour relative to anchor
	var out bytes.Buffer
	err := pipeline.Run(pipeline.Options{
		Sources:        []string{"-"},
		Stdin:          strings.NewReader(input),
		Output:         &out,
		Format:         "json",
		TimeField:      "time",
		TimeRangeExpr:  "last 1h",
		TimeRangeNow:   anchor,
	})
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	lines := nonEmptyLines(out.String())
	// offsets -30m and -5m are within last 1h; -90m and +10m are outside
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d:\n%s", len(lines), out.String())
	}
}

func TestPipelineTimeRangeNoField(t *testing.T) {
	anchor := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	input := `{"level":"info","msg":"no ts"}` + "\n"

	var out bytes.Buffer
	err := pipeline.Run(pipeline.Options{
		Sources:       []string{"-"},
		Stdin:         strings.NewReader(input),
		Output:        &out,
		Format:        "json",
		TimeField:     "time",
		TimeRangeExpr: "last 1h",
		TimeRangeNow:  anchor,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if nonEmptyLines(out.String()) != nil {
		t.Error("expected no output when timestamp field is absent")
	}
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
