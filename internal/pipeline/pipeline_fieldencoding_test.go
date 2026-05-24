package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/logslice/internal/pipeline"
)

func encodingLines(kvs ...any) string {
	m := make(map[string]any)
	for i := 0; i+1 < len(kvs); i += 2 {
		m[kvs[i].(string)] = kvs[i+1]
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func TestPipelineFieldEncodingBase64(t *testing.T) {
	input := encodingLines("msg", "hello", "level", "info")
	src := strings.NewReader(input + "\n")
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources: []pipeline.SourceEntry{{Reader: src}},
		Output:  &buf,
		EncodingRules: []pipeline.EncodingRule{
			{Field: "msg", Mode: "base64-encode", Dest: "msg_b64"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v — %s", err, buf.String())
	}
	if result["msg_b64"] != "aGVsbG8=" {
		t.Fatalf("expected base64 of 'hello', got %v", result["msg_b64"])
	}
	if result["msg"] != "hello" {
		t.Fatalf("original msg should be preserved, got %v", result["msg"])
	}
}

func TestPipelineNoEncodingRulesPassesThrough(t *testing.T) {
	input := encodingLines("msg", "unchanged")
	src := strings.NewReader(input + "\n")
	var buf bytes.Buffer

	err := pipeline.Run(pipeline.Config{
		Sources:       []pipeline.SourceEntry{{Reader: src}},
		Output:        &buf,
		EncodingRules: nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["msg"] != "unchanged" {
		t.Fatalf("expected unchanged, got %v", result["msg"])
	}
}
