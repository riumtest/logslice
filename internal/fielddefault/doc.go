// Package fielddefault provides a transformer that sets default values for
// fields that are absent or explicitly null in a structured log entry.
//
// Rules are evaluated in order. The first matching rule (field missing or null)
// wins. Fields that already carry a non-null value are left unchanged.
//
// Example:
//
//	tr := fielddefault.New([]fielddefault.Rule{
//		{Field: "env",     Value: "production"},
//		{Field: "retries", Value: 0},
//	})
//	out := tr.Transform(entry)
package fielddefault
