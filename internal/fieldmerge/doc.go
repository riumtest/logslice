// Package fieldmerge provides a transformer that combines values from multiple
// source fields into a single destination field.
//
// # Usage
//
// Create a Merger with a destination field name and an ordered list of source
// field names. Call Apply on each log entry map to produce a new map with the
// merged field added.
//
//	merger := fieldmerge.New("message", []string{"service", "event"},
//		fieldmerge.WithSeparator(" | "),
//	)
//
//	out := merger.Apply(entry)
//
// # Options
//
//   - WithSeparator — string placed between joined values (default: " ").
//   - WithSkipMissing — when true (default) absent source fields are omitted;
//     when false a "<missing:field>" placeholder is inserted instead.
//
// The original entry map is never mutated.
package fieldmerge
