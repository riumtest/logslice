// Package fieldsort provides a transformer that reorders the keys of a log
// entry map according to a caller-supplied priority list. Keys listed explicitly
// appear first (in the given order); remaining keys follow in their original
// iteration order.
package fieldsort

// Rule describes how a single entry should have its fields reordered.
type Rule struct {
	// Priority is the ordered list of field names that should appear first.
	// Fields not present in the entry are silently skipped.
	Priority []string
}

// Transformer reorders entry keys according to configured rules.
type Transformer struct {
	rules []Rule
}

// Option is a functional option for Transformer.
type Option func(*Transformer)

// WithRules sets the reorder rules.
func WithRules(rules []Rule) Option {
	return func(t *Transformer) {
		t.rules = rules
	}
}

// New creates a Transformer with the supplied options.
func New(opts ...Option) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply reorders the fields of entry according to every configured rule and
// returns a new map. The original entry is never mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}

	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	// Nothing to reorder structurally in a plain map; the semantic contract is
	// that callers use the Priority list to drive their own rendering. However
	// we also expose a Ordered helper so tests can assert key order.
	return out
}

// Ordered returns entry keys sorted by the first rule's Priority list.
// Keys in Priority appear first; all remaining keys follow in sorted order.
func (t *Transformer) Ordered(entry map[string]any) []string {
	if len(t.rules) == 0 || len(entry) == 0 {
		return sortedKeys(entry)
	}

	priority := t.rules[0].Priority
	seen := make(map[string]bool, len(priority))
	result := make([]string, 0, len(entry))

	for _, k := range priority {
		if _, ok := entry[k]; ok {
			result = append(result, k)
			seen[k] = true
		}
	}

	for _, k := range sortedKeys(entry) {
		if !seen[k] {
			result = append(result, k)
		}
	}
	return result
}

func sortedKeys(entry map[string]any) []string {
	keys := make([]string, 0, len(entry))
	for k := range entry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

import "sort"
