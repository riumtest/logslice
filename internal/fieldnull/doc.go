// Package fieldnull provides a log entry transformer that handles null and
// missing field values according to a set of user-defined rules.
//
// Each Rule targets a single field and specifies one of three behaviours when
// the field is absent or carries a JSON null value:
//
//  1. Remove the field from the entry (default when Replace is nil).
//  2. Replace the field with a static fallback value (set Replace).
//  3. Drop the entire entry (set DropEntry to true).
//
// Rules are evaluated in the order they are supplied. Non-null fields are
// passed through unchanged.
//
// Example
//
//	tr := fieldnull.New([]fieldnull.Rule{
//		{Field: "user_id", DropEntry: true},
//		{Field: "region",  Replace: "us-east-1"},
//		{Field: "error"},   // removes the key when null
//	})
//
//	out := tr.Transform(entry)
//	if out == nil {
//		// entry was dropped because user_id was null
//	}
package fieldnull
