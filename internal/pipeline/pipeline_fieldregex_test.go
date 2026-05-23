package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/pipeline"
)

func regexLines() string {
	return strings.Join([]string{
		`{"msg":"ERROR 503 upstream timeout"}`,
		`{"msg":"INFO 200 ok"}`,
		`{"msg":"no match here"}`,
	}, "\n")
}

func TestPipelineFieldRegex(t *testing.T) {
	src := strings.NewReader(regexLines())
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.NamedReader{{Name: "test", Reader: src}},
		Output:  &buf,
		FieldRegexRules: []pipeline.FieldRegexRule{
			{Source: "msg", Pattern: `(?P<level>\w+)\s+(?P<code>\d+)`},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 3 {
		t.Fatalf("expected 3 output lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "ERROR") {
		t.Errorf("expected first line to contain ERROR, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "503") {
		t.Errorf("expected first line to contain 503, got: %s", lines[0])
	}
}

func TestPipelineNoFieldRegex(t *testing.T) {
	src := strings.NewReader(`{"msg":"ERROR 503"}` + "\n")
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.NamedReader{{Name: "test", Reader: src}},
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := nonEmptyLines(buf.String())
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(lines))
	}
	if strings.Contains(lines[0], "level") {
		t.Errorf("expected no extracted fields, got: %s", lines[0])
	}
}
