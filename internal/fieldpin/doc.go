// Package fieldpin provides a transformer that pins (snapshots) selected
// top-level field values into a named sub-object within each log entry.
//
// This is useful when downstream pipeline stages may rename, redact, or drop
// fields and you want to retain the original values for auditing or debugging.
//
// Example
//
//	rules := []fieldpin.Rule{
//		{Field: "level", Dest: "_original"},
//		{Field: "user_id", Dest: "_original"},
//	}
//	tr := fieldpin.New(rules)
//	out := tr.Apply(entry)
//	// out["_original"] == map[string]any{"level": "warn", "user_id": "u42"}
//
// Rules with an empty Field or Dest are silently ignored.
// If the destination key already holds a map, new pins are merged into it;
// if it holds a non-map value the rule is skipped.
package fieldpin
