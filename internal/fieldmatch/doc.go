// Package fieldmatch provides a log-entry transformer that evaluates
// regular-expression rules against string fields and writes a boolean
// match result to a destination field.
//
// # Usage
//
//	tr := fieldmatch.WithRules([]fieldmatch.Rule{
//		{
//			Field:   "url",
//			Pattern: regexp.MustCompile(`^/api/`),
//			Dest:    "is_api_request",
//		},
//	})
//
//	for _, entry := range entries {
//		annotated := tr.Apply(entry)
//		// annotated["is_api_request"] == true or false
//	}
//
// If the source field is absent or is not a string value, the destination
// field is set to false. The original entry map is never modified.
package fieldmatch
