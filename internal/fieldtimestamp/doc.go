// Package fieldtimestamp provides a Transformer that parses timestamp strings
// in JSON log entries and reformats them according to configurable Go time
// layouts.
//
// Each Rule specifies:
//   - Field      – the source JSON key containing the timestamp string.
//   - InputLayout  – the Go time layout used to parse the raw value
//     (defaults to time.RFC3339).
//   - OutputLayout – the Go time layout used when writing the result
//     (defaults to time.RFC3339).
//   - DestField  – optional target key; when empty the source field is
//     overwritten in-place.
//
// Entries that do not contain the named field, or whose value is not a string,
// are passed through unchanged. A parse failure returns an error so the
// pipeline can surface malformed timestamps.
package fieldtimestamp
