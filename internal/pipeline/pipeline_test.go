package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/pipeline"
)

func TestPipelineNoFilter(t *testing.T) {
	src := strings.NewReader(
		`{"level":"info","msg":"hello"}` + "\n" +
			`{"level":"error","msg":"oops"}` + "\n",
	)

	var out bytes.Buffer
	n, err := pipeline.Run(pipeline.Config{
		Sources: []interface{ Read([]byte) (int, error) }{src},
		Out:     &out,
		Format:  formatter.FormatJSON,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 entries, got %d", n)
	}
}

func TestPipelineWithFilter(t *testing.T) {
	src := strings.NewReader(
		`{"level":"info","msg":"ok"}` + "\n" +
			`{"level":"error","msg":"fail"}` + "\n" +
			`{"level":"warn","msg":"maybe"}` + "\n",
	)

	var out bytes.Buffer
	n, err := pipeline.Run(pipeline.Config{
		Sources:    []interface{ Read([]byte) (int, error) }{src},
		FilterExpr: `level = "error"`,
		Out:        &out,
		Format:     formatter.FormatJSON,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 entry, got %d", n)
	}
	if !strings.Contains(out.String(), "fail") {
		t.Errorf("expected filtered entry in output, got: %s", out.String())
	}
}

func TestPipelineNoSources(t *testing.T) {
	var out bytes.Buffer
	_, err := pipeline.Run(pipeline.Config{Out: &out})
	if err == nil {
		t.Error("expected error for empty sources")
	}
}

func TestPipelineBadFilter(t *testing.T) {
	src := strings.NewReader(`{"level":"info"}` + "\n")
	var out bytes.Buffer
	_, err := pipeline.Run(pipeline.Config{
		Sources:    []interface{ Read([]byte) (int, error) }{src},
		FilterExpr: "!!!",
		Out:        &out,
	})
	if err == nil {
		t.Error("expected parse error for bad filter")
	}
}

func TestPipelineSkipsInvalidJSON(t *testing.T) {
	src := strings.NewReader(
		"not json\n" +
			`{"level":"info","msg":"valid"}` + "\n",
	)
	var out bytes.Buffer
	n, err := pipeline.Run(pipeline.Config{
		Sources: []interface{ Read([]byte) (int, error) }{src},
		Out:     &out,
		Format:  formatter.FormatJSON,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 valid entry, got %d", n)
	}
}
