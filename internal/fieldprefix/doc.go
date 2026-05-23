// Package fieldprefix transforms string fields in a log entry by adding or
// stripping a fixed prefix.
//
// Usage:
//
//	rules := []fieldprefix.Rule{
//		{Field: "service", Prefix: "svc-"},           // add
//		{Field: "request_id", Prefix: "req-", Strip: true}, // remove
//	}
//	tr := fieldprefix.New(rules)
//	outEntry := tr.Transform(inEntry)
//
// Fields that are absent or not of type string are silently skipped.
// The original entry map is never modified; a shallow copy is returned.
package fieldprefix
