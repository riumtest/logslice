// Package fieldparse provides a Transformer that coerces nominated string
// fields inside a log entry into their natural Go types.
//
// Many log pipelines emit every value as a JSON string even when the
// underlying data is numeric, boolean, or a nested object.  fieldparse
// lets callers nominate specific field names that should be re-interpreted:
//
//	tr := fieldparse.New(
//		fieldparse.WithFields("duration_ms", "retries", "success"),
//	)
//	entry := map[string]any{"duration_ms": "142", "success": "true"}
//	parsed := tr.Transform(entry)
//	// parsed["duration_ms"] == float64(142)
//	// parsed["success"]     == true
//
// Fields are tried in order: bool → float64 → JSON object/array → original
// string.  Fields that are absent or already non-string are passed through
// unchanged.  The original entry map is never mutated.
package fieldparse
