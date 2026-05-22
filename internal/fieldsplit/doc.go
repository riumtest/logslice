// Package fieldsplit provides a transformer that splits a string field into
// multiple named fields using a configurable delimiter.
//
// Example — split an "addr" field of the form "host:port" into separate
// "host" and "port" fields:
//
//	s := fieldsplit.New(fieldsplit.WithRules([]fieldsplit.Rule{
//		{
//			Source:    "addr",
//			Delimiter: ":",
//			Targets:   []string{"host", "port"},
//		},
//	}))
//
//	out := s.Apply(entry)
//
// The original source field is preserved in the output. If the source field
// is absent or is not a string, the rule is silently skipped. Parts beyond
// the number of declared targets are discarded; if there are fewer parts than
// targets, only the available targets are populated.
package fieldsplit
