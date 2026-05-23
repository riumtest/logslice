// Package fieldwindow provides a sliding-window rolling-average transformer
// for numeric fields in structured log entries.
//
// Each Rule specifies a source field, a destination field, and a window size.
// As entries are processed the transformer maintains a per-rule circular
// buffer of the most recent Size values and writes their arithmetic mean to
// the destination field.
//
// Example usage:
//
//	tr := fieldwindow.New([]fieldwindow.Rule{
//		{SourceField: "response_ms", DestField: "response_ms_avg3", Size: 3},
//	})
//	for _, entry := range entries {
//		processed := tr.Transform(entry)
//	}
//
// Rules with a Size less than 1 are silently discarded. Entries that are
// missing the source field or whose value cannot be decoded as a float64 are
// passed through without writing the destination field for that rule.
package fieldwindow
