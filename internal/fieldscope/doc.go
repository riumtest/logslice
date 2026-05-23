// Package fieldscope provides a transformer that promotes nested JSON fields
// to top-level keys in a log entry.
//
// Many structured loggers emit deeply nested objects (e.g. an "http" block
// containing method, path, and status). fieldscope lets you extract specific
// leaf values and surface them as flat, top-level fields so that downstream
// filters, formatters, and aggregators can reference them directly.
//
// Example:
//
//	rules := []fieldscope.Rule{
//		{Source: "http.request.method", Dest: "method"},
//		{Source: "http.response.status", Dest: "status"},
//	}
//	tr := fieldscope.New(rules)
//	out := tr.Apply(entry)
//	// out["method"] == "GET"
//	// out["status"] == 200
package fieldscope
