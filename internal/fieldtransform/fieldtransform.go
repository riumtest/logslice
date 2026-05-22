// Package fieldtransform provides a Transformer that renames or remaps
// fields in a log entry map before formatting.
package fieldtransform

import "fmt"

// Option configures a Transformer.
type Option func(*Transformer)

// Transformer renames fields in a log entry according to a mapping.
type Transformer struct {
	renames map[string]string // oldKey -> newKey
}

// WithRename adds a rename rule: the field named `from` will be output as `to`.
// If `to` is empty the call is a no-op.
func WithRename(from, to string) Option {
	return func(t *Transformer) {
		if from != "" && to != "" {
			t.renames[from] = to
		}
	}
}

// New creates a Transformer with the supplied options.
func New(opts ...Option) (*Transformer, error) {
	t := &Transformer{
		renames: make(map[string]string),
	}
	for _, o := range opts {
		o(t)
	}
	// Detect collisions: two source fields mapped to the same target.
	seen := make(map[string]string)
	for from, to := range t.renames {
		if prev, ok := seen[to]; ok {
			return nil, fmt.Errorf("fieldtransform: rename conflict: both %q and %q map to %q", prev, from, to)
		}
		seen[to] = from
	}
	return t, nil
}

// Apply returns a shallow copy of entry with fields renamed according to the
// configured mapping. Fields whose target name already exists in entry are
// skipped to avoid silent overwrites.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.renames) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		if newKey, ok := t.renames[k]; ok {
			// Only rename when the target key is not already present.
			if _, exists := entry[newKey]; !exists {
				out[newKey] = v
				continue
			}
		}
		out[k] = v
	}
	return out
}
