package timerange_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timerange"
)

var anchor = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestParseLastDuration(t *testing.T) {
	f, err := timerange.Parse("last 30m", anchor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedFrom := anchor.Add(-30 * time.Minute)
	if !f.From.Equal(expectedFrom) {
		t.Errorf("From: got %v, want %v", f.From, expectedFrom)
	}
	if !f.To.Equal(anchor) {
		t.Errorf("To: got %v, want %v", f.To, anchor)
	}
}

func TestParseExplicitRange(t *testing.T) {
	expr := "2024-06-01T10:00:00Z,2024-06-01T11:00:00Z"
	f, err := timerange.Parse(expr, anchor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.From.Hour() != 10 || f.To.Hour() != 11 {
		t.Errorf("unexpected range: %v – %v", f.From, f.To)
	}
}

func TestParseError(t *testing.T) {
	cases := []string{"bad", "last xyz", "notadate,alsonotadate"}
	for _, c := range cases {
		_, err := timerange.Parse(c, anchor)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func makeRecord(ts string) map[string]json.RawMessage {
	v, _ := json.Marshal(ts)
	return map[string]json.RawMessage{"time": v}
}

func TestMatchWithinRange(t *testing.T) {
	f := timerange.Filter{
		From: anchor.Add(-1 * time.Hour),
		To:   anchor.Add(1 * time.Hour),
	}
	if !f.Match(makeRecord(anchor.Format(time.RFC3339)), "time") {
		t.Error("expected match inside range")
	}
}

func TestMatchOutsideRange(t *testing.T) {
	f := timerange.Filter{
		From: anchor.Add(-1 * time.Hour),
		To:   anchor,
	}
	after := anchor.Add(10 * time.Minute).Format(time.RFC3339)
	if f.Match(makeRecord(after), "time") {
		t.Error("expected no match outside range")
	}
}

func TestMatchZeroFilterAlwaysTrue(t *testing.T) {
	var f timerange.Filter
	if !f.Match(makeRecord(anchor.Format(time.RFC3339)), "time") {
		t.Error("zero filter should always match")
	}
}

func TestMatchMissingField(t *testing.T) {
	f := timerange.Filter{From: anchor.Add(-time.Hour), To: anchor.Add(time.Hour)}
	if f.Match(map[string]json.RawMessage{}, "time") {
		t.Error("expected no match for missing field")
	}
}
