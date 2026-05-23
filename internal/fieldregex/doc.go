// Package fieldregex provides a pipeline stage that extracts fields from
// string values using regular expressions with named capture groups.
//
// Each Rule specifies a source field to read, a pattern to match, and an
// optional prefix applied to all extracted field names. Named capture groups
// in the pattern become new fields in the output entry.
//
// Example:
//
//	rule := fieldregex.Rule{
//		Source:  "message",
//		Pattern: `(?P<level>\w+)\s+(?P<code>\d+)`,
//		Prefix:  "parsed_",
//	}
//	e, err := fieldregex.New(fieldregex.WithRules([]fieldregex.Rule{rule}))
//	if err != nil { ... }
//	out := e.Apply(entry)
//
// If the source field is absent or not a string, the rule is silently skipped.
// If the pattern does not match, no new fields are added.
package fieldregex
