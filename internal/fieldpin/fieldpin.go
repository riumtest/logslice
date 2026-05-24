// Package fieldpin copies a set of top-level field values into a nested
// sub-object so that they are preserved after downstream transforms that may
// rename or drop the originals.
package fieldpin

import "fmt"

// Rule describes a single pin operation: the source field whose value should
// be pinned and the destination sub-object key under which it will be stored.
type Rule struct {
	Field string // source field name
	Dest  string // destination sub-object name (created if absent)
}

// Transformer pins field values into a snapshot sub-object.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the given rules.
// Rules with empty Field or Dest are silently ignored.
func New(rules []Rule) *Transformer {
	filtered := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Field != "" && r.Dest != "" {
			filtered = append(filtered, r)
		}
	}
	return &Transformer{rules: filtered}
}

// Apply copies source field values into their respective destination
// sub-objects. The original entry is never mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}

	out := shallowCopy(entry)

	for _, r := range t.rules {
		val, ok := out[r.Field]
		if !ok {
			continue
		}

		// Retrieve or create the destination sub-object.
		var dest map[string]any
		if existing, ok := out[r.Dest]; ok {
			if m, ok := existing.(map[string]any); ok {
				dest = shallowCopy(m)
			} else {
				// Destination exists but is not a map – skip.
				continue
			}
		} else {
			dest = make(map[string]any)
		}

		dest[fmt.Sprintf("%s", r.Field)] = val
		out[r.Dest] = dest
	}

	return out
}

func shallowCopy(m map[string]any) map[string]any {
	copied := make(map[string]any, len(m))
	for k, v := range m {
		copied[k] = v
	}
	return copied
}
