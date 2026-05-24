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
//   - WithOverwrite — when true (default) an existing destination field is
//     overwritten; when false the entry is returned unchanged if the
//     destination field already exists.
//
// # Behaviour
//
// Source field values are converted to strings via fmt.Sprintf("%v", v).
// Empty string values are treated as present and included in the join.
// The original entry map is never mutated.
package fieldmerge
