// Package timerange provides helpers for parsing and matching time-based
// filters against structured log records. It supports ISO-8601 timestamps
// as well as relative durations such as "last 5m" or "last 1h".
package timerange

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Filter holds an optional start and end time used to include only log
// records whose timestamp falls within the range [From, To).
type Filter struct {
	From time.Time
	To   time.Time
}

// IsZero reports whether the filter has no constraints.
func (f Filter) IsZero() bool {
	return f.From.IsZero() && f.To.IsZero()
}

// Match returns true when the record's timestamp field falls within the
// filter range. If the filter is zero it always returns true. The field
// value must be an RFC3339 string; unparseable values are treated as
// non-matching.
func (f Filter) Match(record map[string]json.RawMessage, field string) bool {
	if f.IsZero() {
		return true
	}
	raw, ok := record[field]
	if !ok {
		return false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return false
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return false
	}
	if !f.From.IsZero() && t.Before(f.From) {
		return false
	}
	if !f.To.IsZero() && !t.Before(f.To) {
		return false
	}
	return true
}

// Parse parses a time-range expression and returns a Filter. Accepted forms:
//
//	"last <duration>"   – e.g. "last 5m", "last 2h"
//	"<from>,<to>"       – two RFC3339 timestamps separated by a comma
func Parse(expr string, now time.Time) (Filter, error) {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "last ") {
		durStr := strings.TrimSpace(strings.TrimPrefix(expr, "last "))
		d, err := time.ParseDuration(durStr)
		if err != nil {
			return Filter{}, fmt.Errorf("timerange: invalid duration %q: %w", durStr, err)
		}
		return Filter{From: now.Add(-d), To: now}, nil
	}
	parts := strings.SplitN(expr, ",", 2)
	if len(parts) != 2 {
		return Filter{}, fmt.Errorf("timerange: expected \"last <dur>\" or \"<from>,<to>\", got %q", expr)
	}
	from, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[0]))
	if err != nil {
		return Filter{}, fmt.Errorf("timerange: invalid from timestamp: %w", err)
	}
	to, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[1]))
	if err != nil {
		return Filter{}, fmt.Errorf("timerange: invalid to timestamp: %w", err)
	}
	return Filter{From: from, To: to}, nil
}
