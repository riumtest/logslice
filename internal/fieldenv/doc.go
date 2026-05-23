// Package fieldenv provides a log entry transformer that injects values from
// environment variables into structured log records.
//
// Each Rule maps one environment variable name to a destination field in the
// entry. If the variable is unset or empty, an optional Default string is used
// instead. When neither the variable nor a default is available, the field is
// left absent.
//
// Existing fields in the entry are never overwritten, making it safe to apply
// this transformer even when some records already carry the destination field.
//
// Example:
//
//	rules := []fieldenv.Rule{
//		{Env: "APP_ENV",  Dest: "environment", Default: "development"},
//		{Env: "AWS_REGION", Dest: "region"},
//	}
//	tr := fieldenv.New(rules)
//	enriched := tr.Transform(entry)
package fieldenv
