// Package fielduniq provides a transformer that removes duplicate elements
// from array-typed fields in a log entry.
//
// Given a field whose JSON value is an array of strings (or numbers), the
// transformer retains only the first occurrence of each element, preserving
// the original order.
//
// Example configuration:
//
//	tr := fielduniq.New([]string{"tags", "labels"})
//	out := tr.Apply(entry)
//
// Fields that are absent, null, or not an array are passed through unchanged.
package fielduniq
