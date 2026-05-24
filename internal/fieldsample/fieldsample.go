package fieldsample

import "math/rand"

// Rule defines a sampling rule: keep only every Nth entry where the
// given field matches the provided value. If Value is empty, the rule
// applies to all entries regardless of field content.
type Rule struct {
	Field string
	Value string
	N     int // keep 1 in every N entries
}

// Transformer probabilistically drops entries based on field-value
// sampling rules. Each rule maintains its own counter so that
// different field values are sampled independently.
type Transformer struct {
	rules   []Rule
	counts  []int
	randFn  func(n int) int
}

// Option configures a Transformer.
type Option func(*Transformer)

// WithRandFn replaces the default random source (useful for testing).
func WithRandFn(fn func(n int) int) Option {
	return func(t *Transformer) { t.randFn = fn }
}

// New creates a Transformer from the given rules.
// Rules with N < 1 are normalised to 1 (keep all).
func New(rules []Rule, opts ...Option) *Transformer {
	normalised := make([]Rule, len(rules))
	for i, r := range rules {
		if r.N < 1 {
			r.N = 1
		}
		normalised[i] = r
	}
	t := &Transformer{
		rules:  normalised,
		counts: make([]int, len(normalised)),
		randFn: rand.Intn,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply returns nil when the entry should be dropped, otherwise a
// shallow copy of entry is returned unchanged.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	for i, r := range t.rules {
		if r.Field != "" {
			v, ok := entry[r.Field]
			if !ok {
				continue
			}
			if r.Value != "" {
				s, _ := v.(string)
				if s != r.Value {
					continue
				}
			}
		}
		// Rule matches — apply sampling.
		if r.N == 1 {
			continue // keep all
		}
		t.counts[i]++
		if t.randFn(r.N) != 0 {
			return nil
		}
	}
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	return out
}
