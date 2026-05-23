package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/naturalselectionlabs/logslice/internal/fieldwindow"
	"github.com/naturalselectionlabs/logslice/internal/pipeline"
)

func windowLines(vals []float64) []string {
	lines := make([]string, len(vals))
	for i, v := range vals {
		b, _ := json.Marshal(map[string]any{"v": v, "level": "info"})
		lines[i] = string(b)
	}
	return lines
}

func TestPipelineFieldWindowRollingAvg(t *testing.T) {
	lines := windowLines([]float64{10, 20, 30, 40})
	src := strings.NewReader(strings.Join(lines, "\n"))
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:      []interface{ Read() (map[string]json.RawMessage, error) }{},
		WindowRules:  []fieldwindow.Rule{{SourceField: "v", DestField: "v_avg", Size: 3}},
		RawReader:    src,
		Output:       &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outLines := nonEmptyLines(buf.String())
	if len(outLines) != 4 {
		t.Fatalf("expected 4 output lines, got %d", len(outLines))
	}

	// Last line window = [20,30,40], avg = 30
	var last map[string]json.RawMessage
	if err := json.Unmarshal([]byte(outLines[3]), &last); err != nil {
		t.Fatalf("unmarshal last line: %v", err)
	}
	var avg float64
	if err := json.Unmarshal(last["v_avg"], &avg); err != nil {
		t.Fatalf("unmarshal v_avg: %v", err)
	}
	if avg != 30.0 {
		t.Fatalf("want avg=30, got %v", avg)
	}
}

func TestPipelineNoWindowRulesPassesThrough(t *testing.T) {
	lines := windowLines([]float64{5, 10})
	src := strings.NewReader(strings.Join(lines, "\n"))
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:   []interface{ Read() (map[string]json.RawMessage, error) }{},
		RawReader: src,
		Output:    &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outLines := nonEmptyLines(buf.String())
	if len(outLines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(outLines))
	}
	var first map[string]json.RawMessage
	if err := json.Unmarshal([]byte(outLines[0]), &first); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := first["v_avg"]; ok {
		t.Fatal("expected v_avg to be absent when no window rules configured")
	}
}
