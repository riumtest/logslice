// Package fieldsample provides a transformer that probabilistically
// drops log entries based on field-value sampling rules.
//
// Each Rule targets a specific field (and optionally a specific value
// within that field). When an entry matches a rule, it is kept with
// probability 1/N and dropped otherwise. Multiple rules are evaluated
// in order; the first matching rule that decides to drop the entry
// wins.
//
// Example usage:
//
//	rules := []fieldsample.Rule{
//		{Field: "level", Value: "debug", N: 100}, // keep 1% of debug
//		{Field: "level", Value: "info",  N: 10},  // keep 10% of info
//	}
//	tr := fieldsample.New(rules)
//	out := tr.Apply(entry)
//	if out == nil {
//		// entry was dropped
//	}
package fieldsample
