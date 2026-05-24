// Package fieldxpath provides a transformer that extracts values from nested
// JSON log entries using dot-notation path expressions.
//
// # Usage
//
// Create a set of rules, each specifying a dot-separated Path to read from and
// a Dest field to write the resolved value into:
//
//	rules := []fieldxpath.Rule{
//		{Path: "kubernetes.pod.name", Dest: "pod"},
//		{Path: "http.request.method", Dest: "method"},
//	}
//	tr := fieldxpath.New(rules)
//	out := tr.Apply(entry)
//
// If any segment of the path is missing or the intermediate value is not a
// JSON object, the rule is silently skipped and the entry is returned unchanged
// for that rule. The original entry is never mutated.
package fieldxpath
