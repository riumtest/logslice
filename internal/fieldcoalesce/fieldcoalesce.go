// Package fieldcoalesce provides a transformer that sets a destination field
// to the value of the first non-null, non-empty source field found.
package fieldcoalesce

import "fmt"

// Rule describes a single coalesce operation: check each field in Sources in
// order and write the first non-empty value to Dest.
type Rule struct {
	Sources []string
	Dest    string
}

// Transformer applies coalesce rules to log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer configured with the provided rules.
func New(rules []Rule) (*Transformer, error) {
	for i, r := range rules {
		if len(r.Sources) == 0 {
			return nil, fmt.Errorf("rule %d: sources must not be empty", i)
		}
		if r.Dest == "" {
			return nil, fmt.Errorf("rule %d: dest must not be empty", i)
		}
	}
	return &Transformer{rules: rules}, nil
}

// Transform applies each rule to entry, returning a new map with dest fields
// populated from the first non-empty source value found.
func (t *Transformer) Transform(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		for _, src := range r.Sources {
			v, ok := out[src]
			if !ok || v == nil {
				continue
			}
			if s, ok := v.(string); ok && s == "" {
				continue
			}
			out[r.Dest] = v
			break
		}
	}
	return out
}
