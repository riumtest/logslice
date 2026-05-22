// Package fieldround provides a transformer that rounds numeric fields
// in a log entry according to configurable rules.
//
// Each Rule specifies:
//   - Field: the JSON key to target
//   - Mode: one of "round" (default), "floor", or "ceil"
//   - Precision: number of decimal places to keep (0 = integer)
//
// Non-numeric or missing fields are passed through unchanged.
// The original entry map is never mutated.
//
// Example:
//
//	tr := fieldround.New([]fieldround.Rule{
//	    {Field: "latency_ms", Mode: "round", Precision: 1},
//	    {Field: "score",      Mode: "floor", Precision: 0},
//	})
//	out := tr.Apply(entry)
package fieldround
