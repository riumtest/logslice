package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/logslice/internal/pipeline"
)

func typedLines(kv ...any) []string {
	var lines []string
	for i := 0; i+1 < len(kv); i += 2 {
		m := map[string]any{kv[i].(string): kv[i+1]}
		b, _ := json.Marshal(m)
		lines = append(lines, string(b))
	}
	return lines
}

func TestPipelineTypeCheckAnnotates(t *testing.T) {
	lines := []string{
		`{"count":"bad"}`,
		`{"count":3}`,
	}
	src := strings.NewReader(strings.Join(lines, "\n"))
	var buf bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.SourceConfig{{Reader: src}},
		TypeCheckRules: []pipeline.TypeCheckRule{{Field: "count", Expected: "number"}},
		TypeCheckDestField: "_type_errors",
		Output: &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "_type_errors") {
		t.Error("expected _type_errors annotation in output")
	}
	if strings.Count(out, "\n") != 2 {
		t.Errorf("expected 2 output lines, got: %q", out)
	}
}

func TestPipelineTypeCheckRejectMode(t *testing.T) {
	lines := []string{
		`{"active":"yes"}`,
		`{"active":true}`,
	}
	src := strings.NewReader(strings.Join(lines, "\n"))
	var buf bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.SourceConfig{{Reader: src}},
		TypeCheckRules: []pipeline.TypeCheckRule{{Field: "active", Expected: "bool"}},
		TypeCheckRejectMode: true,
		Output: &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Count(out, "\n") != 1 {
		t.Errorf("expected 1 output line after rejection, got: %q", out)
	}
}

func TestPipelineNoTypeCheckPassesAll(t *testing.T) {
	lines := []string{`{"msg":"a"}`, `{"msg":"b"}`}
	src := strings.NewReader(strings.Join(lines, "\n"))
	var buf bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.SourceConfig{{Reader: src}},
		Output: &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Count(buf.String(), "\n") != 2 {
		t.Error("expected all lines to pass through")
	}
}
