// Package fieldexist filters or annotates log entries based on the
// presence or absence of specified fields.
package fieldexist

// Rule defines a single existence check.
type Rule struct {
	// Field is the key to check.
	Field string
	// MustExist controls whether the field must be present (true) or absent (false).
	MustExist bool
	// DestField, if non-empty, writes a boolean result instead of filtering.
	DestField string
}

// Transformer checks field existence on each entry.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer configured with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply evaluates each rule against entry.
// If a rule has no DestField, the entry is dropped when the condition is not met.
// If a rule has a DestField, a boolean is written and the entry is always kept.
// Returns nil to signal the entry should be dropped.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}

	out := shallowCopy(entry)

	for _, r := range t.rules {
		_, exists := out[r.Field]

		if r.DestField != "" {
			out[r.DestField] = exists
			continue
		}

		// Filter mode: drop entry when condition not satisfied.
		if r.MustExist && !exists {
			return nil
		}
		if !r.MustExist && exists {
			return nil
		}
	}

	return out
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
