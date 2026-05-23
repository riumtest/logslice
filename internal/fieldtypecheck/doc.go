// Package fieldtypecheck provides a transformer that validates the JSON types
// of named fields in a log entry against expected types.
//
// When a field's runtime type does not match the declared expectation the
// checker can either:
//
//   - Annotate the entry with a list of type-error descriptions written to a
//     configurable destination field (default behaviour when WithDestField is
//     used).
//   - Drop the entry entirely when WithRejectMode is applied.
//
// Supported type names match JSON primitives: "string", "number", "bool",
// "null", "array", and "object".
package fieldtypecheck
