// Package fieldtemplate renders Go text/template expressions into new or
// existing fields of a log entry.
//
// Each Rule maps a destination field name to a template string. The template
// is executed with the current log entry (map[string]any) as its data object,
// so any field value can be referenced via the built-in index function:
//
//	{{index . "level"}}: {{index . "msg"}}
//
// Rules are compiled once at construction time via WithRules; a parse error
// in any template causes WithRules to return an error immediately.
//
// The Apply method does not mutate the original entry — it always returns a
// shallow copy with the rendered fields added or overwritten.
package fieldtemplate
