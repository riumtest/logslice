// Package fieldflip swaps the values of two fields within a log entry.
package fieldflip

// Rule describes a pair of fields whose values should be exchanged.
type Rule struct {
	A string // source field
	B string // destination field
}

// Flipper swaps field values according to a set of rules.
type Flipper struct {
	rules []Rule
}

// New returns a Flipper that applies the given swap rules.
func New(rules []Rule) *Flipper {
	return &Flipper{rules: rules}
}

// Apply exchanges the values of each rule's A and B fields in a copy of the
// entry. If either field is absent the entry is returned unchanged for that
// rule. The original map is never mutated.
func (f *Flipper) Apply(entry map[string]any) map[string]any {
	if len(f.rules) == 0 {
		return entry
	}

	out := shallowCopy(entry)
	for _, r := range f.rules {
		valA, okA := out[r.A]
		valB, okB := out[r.B]
		if !okA || !okB {
			continue
		}
		out[r.A] = valB
		out[r.B] = valA
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
