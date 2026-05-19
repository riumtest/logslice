// Package timerange implements time-range filtering for structured log
// records.
//
// # Overview
//
// A [Filter] holds an optional start (From) and end (To) time. When applied
// to a log record it inspects a named timestamp field and returns true only
// when the field's value falls within [From, To).
//
// # Parsing
//
// [Parse] accepts two expression formats:
//
//   - "last <duration>" — a Go duration string prefixed with "last ",
//     e.g. "last 5m" or "last 2h". The range is computed relative to the
//     supplied now value.
//
//   - "<from>,<to>" — two RFC3339 timestamps separated by a comma.
//
// # Usage
//
//	f, err := timerange.Parse("last 1h", time.Now())
//	if err != nil { ... }
//	if f.Match(record, "timestamp") {
//	    // record is within the last hour
//	}
package timerange
