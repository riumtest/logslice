// Package fieldcopy copies the value of one field into another field
// for each log entry, leaving the source field intact.
package fieldcopy

// Rule describes a single copy operation: read from Src and write to Dst.
type Rule struct {
	Src string
	Dst string
}

// Transformer copies fields according to a set of rules.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the given rules in order.
// Rules with empty Src or Dst are silently skipped.
func New(rules []Rule) *Transformer {
	valid := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Src != "" && r.Dst != "" {
			valid = append(valid, r)
		}
	}
	return &Transformer{rules: valid}
}

// Apply copies fields as configured and returns the modified entry.
// The original map is never mutated; a shallow copy is made first.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range t.rules {
		v, ok := out[r.Src]
		if !ok {
			continue
		}
		out[r.Dst] = v
	}
	return out
}

func shallowCopy(m map[string]any) map[string]any {
	copy := make(map[string]any, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}
