package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/pipeline"
)

func geoLines(ips ...string) string {
	var sb strings.Builder
	for _, ip := range ips {
		b, _ := json.Marshal(map[string]any{"ip": ip, "msg": "req"})
		sb.Write(b)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func TestPipelineFieldGeoAnnotates(t *testing.T) {
	t.Parallel()
	input := geoLines("10.0.0.1", "1.2.3.4")
	var out bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Reader: strings.NewReader(input),
		Writer: &out,
		GeoRules: []pipeline.GeoRule{
			{
				SrcField:  "ip",
				DestField: "region",
				Regions:   map[string]string{"10.0.0.0/8": "private"},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmptyLines(out.String())
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var first map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if first["region"] != "private" {
		t.Errorf("expected region=private, got %v", first["region"])
	}
	var second map[string]any
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := second["region"]; ok {
		t.Error("unmatched IP should not have region field")
	}
}

func TestPipelineNoGeoRulesPassesThrough(t *testing.T) {
	t.Parallel()
	input := geoLines("10.0.0.1")
	var out bytes.Buffer
	err := pipeline.Run(pipeline.Config{
		Reader: strings.NewReader(input),
		Writer: &out,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := nonEmptyLines(out.String())
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["region"]; ok {
		t.Error("no geo rules: region field must be absent")
	}
}
