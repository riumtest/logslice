// Package fieldabsent removes specified fields from log entries.
package fieldabsent

// Rule describes a field to remove from a log entry.
type Rule struct {
	// Field is the key to remove.
	Field string
}

// Transformer removes configured fields from each entry.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that will drop the given fields from every entry.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Transform returns a shallow copy of entry with the configured fields removed.
// If no rules are configured the original entry is returned unchanged.
func (t *Transformer) Transform(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}

	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	for _, r := range t.rules {
		if r.Field == "" {
			continue
		}
		delete(out, r.Field)
	}

	return out
}
