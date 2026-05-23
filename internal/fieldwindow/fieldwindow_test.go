package fieldwindow_test

import (
	"encoding/json"
	"testing"

	"github.com/naturalselectionlabs/logslice/internal/fieldwindow"
)

func entry(fields map[string]any) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(fields))
	for k, v := range fields {
		b, _ := json.Marshal(v)
		out[k] = b
	}
	return out
}

func floatVal(t *testing.T, e map[string]json.RawMessage, key string) float64 {
	t.Helper()
	var f float64
	if err := json.Unmarshal(e[key], &f); err != nil {
		t.Fatalf("key %q: %v", key, err)
	}
	return f
}

func TestIdentityNoRules(t *testing.T) {
	tr := fieldwindow.New(nil)
	in := entry(map[string]any{"latency": 10.0})
	out := tr.Transform(in)
	if _, ok := out["latency"]; !ok {
		t.Fatal("expected latency field to be preserved")
	}
}

func TestRollingAverageSingleValue(t *testing.T) {
	tr := fieldwindow.New([]fieldwindow.Rule{
		{SourceField: "latency", DestField: "latency_avg", Size: 3},
	})
	out := tr.Transform(entry(map[string]any{"latency": 30.0}))
	if got := floatVal(t, out, "latency_avg"); got != 30.0 {
		t.Fatalf("want 30, got %v", got)
	}
}

func TestRollingAverageAccumulates(t *testing.T) {
	tr := fieldwindow.New([]fieldwindow.Rule{
		{SourceField: "v", DestField: "v_avg", Size: 3},
	})
	for _, val := range []float64{10, 20, 30} {
		tr.Transform(entry(map[string]any{"v": val}))
	}
	out := tr.Transform(entry(map[string]any{"v": 40.0}))
	// window holds [20,30,40], avg = 30
	if got := floatVal(t, out, "v_avg"); got != 30.0 {
		t.Fatalf("want 30, got %v", got)
	}
}

func TestSkipsMissingSourceField(t *testing.T) {
	tr := fieldwindow.New([]fieldwindow.Rule{
		{SourceField: "missing", DestField: "out", Size: 2},
	})
	out := tr.Transform(entry(map[string]any{"other": "x"}))
	if _, ok := out["out"]; ok {
		t.Fatal("expected dest field to be absent when source is missing")
	}
}

func TestSkipsNonNumericField(t *testing.T) {
	tr := fieldwindow.New([]fieldwindow.Rule{
		{SourceField: "v", DestField: "v_avg", Size: 2},
	})
	out := tr.Transform(entry(map[string]any{"v": "not-a-number"}))
	if _, ok := out["v_avg"]; ok {
		t.Fatal("expected dest field to be absent for non-numeric source")
	}
}

func TestZeroSizeRuleIgnored(t *testing.T) {
	tr := fieldwindow.New([]fieldwindow.Rule{
		{SourceField: "v", DestField: "v_avg", Size: 0},
	})
	out := tr.Transform(entry(map[string]any{"v": 5.0}))
	if _, ok := out["v_avg"]; ok {
		t.Fatal("expected zero-size rule to be ignored")
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldwindow.New([]fieldwindow.Rule{
		{SourceField: "v", DestField: "v_avg", Size: 2},
	})
	in := entry(map[string]any{"v": 7.0})
	tr.Transform(in)
	if _, ok := in["v_avg"]; ok {
		t.Fatal("original entry was mutated")
	}
}
