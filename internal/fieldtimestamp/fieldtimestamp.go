package fieldtimestamp

import (
	"fmt"
	"time"
)

// Rule describes a timestamp parse-and-reformat operation on a single field.
type Rule struct {
	// Field is the JSON key to operate on.
	Field string
	// InputLayout is the Go time layout used to parse the existing value.
	// Use time.RFC3339 if empty.
	InputLayout string
	// OutputLayout is the Go time layout used to format the result.
	// Use time.RFC3339 if empty.
	OutputLayout string
	// DestField, when non-empty, writes the result to a new key instead of
	// overwriting the source field.
	DestField string
}

// Transformer rewrites timestamp fields according to a set of rules.
type Transformer struct {
	rules []Rule
}

// WithRules returns a functional option that sets the transformation rules.
func WithRules(rules []Rule) func(*Transformer) {
	return func(t *Transformer) {
		t.rules = rules
	}
}

// New creates a Transformer with the supplied options.
func New(opts ...func(*Transformer)) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply rewrites timestamp fields in entry according to the configured rules.
// The original map is not mutated; a shallow copy is returned.
func (t *Transformer) Apply(entry map[string]any) (map[string]any, error) {
	if len(t.rules) == 0 {
		return entry, nil
	}
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		raw, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := raw.(string)
		if !ok {
			continue
		}
		inLayout := r.InputLayout
		if inLayout == "" {
			inLayout = time.RFC3339
		}
		outLayout := r.OutputLayout
		if outLayout == "" {
			outLayout = time.RFC3339
		}
		parsed, err := time.Parse(inLayout, s)
		if err != nil {
			return nil, fmt.Errorf("fieldtimestamp: field %q value %q: %w", r.Field, s, err)
		}
		formatted := parsed.UTC().Format(outLayout)
		dest := r.DestField
		if dest == "" {
			dest = r.Field
		}
		out[dest] = formatted
	}
	return out, nil
}
