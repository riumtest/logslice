// Package fieldformat provides string formatting transformations for fields
// in structured log entries.
//
// Supported operations:
//
//   - OpUppercase — converts the field value to upper-case.
//   - OpLowercase — converts the field value to lower-case.
//   - OpTruncate  — truncates the field value to at most MaxLen characters.
//
// Only fields whose values are strings are affected; numeric or boolean values
// are left unchanged. Missing fields are silently skipped.
//
// Example:
//
//	f := fieldformat.WithRules([]fieldformat.Rule{
//		{Field: "level", Op: fieldformat.OpUppercase},
//		{Field: "msg",   Op: fieldformat.OpTruncate, MaxLen: 80},
//	})
//	out := f.Apply(entry)
package fieldformat
