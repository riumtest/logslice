// Package fieldenv enriches log entries with values from environment variables.
package fieldenv

import "os"

// Rule maps an environment variable to a destination field in the log entry.
type Rule struct {
	// Env is the name of the environment variable to read.
	Env string
	// Dest is the field name to write the value into.
	Dest string
	// Default is used when the environment variable is unset or empty.
	Default string
}

// Transformer injects environment variable values into log entries.
type Transformer struct {
	rules  []Rule
	lookup func(string) (string, bool)
}

// Option configures a Transformer.
type Option func(*Transformer)

// WithLookup overrides the environment lookup function (useful for testing).
func WithLookup(fn func(string) (string, bool)) Option {
	return func(t *Transformer) {
		t.lookup = fn
	}
}

// New creates a Transformer that applies the given rules.
func New(rules []Rule, opts ...Option) *Transformer {
	t := &Transformer{
		rules:  rules,
		lookup: os.LookupEnv,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Transform applies environment variable rules to a copy of the entry.
// Existing fields are never overwritten.
func (t *Transformer) Transform(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry)+len(t.rules))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		if _, exists := out[r.Dest]; exists {
			continue
		}
		val, ok := t.lookup(r.Env)
		if !ok || val == "" {
			val = r.Default
		}
		if val != "" {
			out[r.Dest] = val
		}
	}
	return out
}
