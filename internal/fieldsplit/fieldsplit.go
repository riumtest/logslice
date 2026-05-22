// Package fieldsplit splits a single string field into multiple fields
// using a configurable delimiter, writing the parts into new named fields.
package fieldsplit

import "strings"

// Rule describes how to split one field into several.
type Rule struct {
	// Source is the field whose value will be split.
	Source string
	// Delimiter is the string used to split the value (default ":").
	Delimiter string
	// Targets are the destination field names for each part, in order.
	// Extra parts beyond len(Targets) are discarded; missing parts are skipped.
	Targets []string
}

// Option configures a Splitter.
type Option func(*Splitter)

// WithRules sets the split rules to apply.
func WithRules(rules []Rule) Option {
	return func(s *Splitter) {
		s.rules = rules
	}
}

// Splitter applies field-split rules to log entries.
type Splitter struct {
	rules []Rule
}

// New creates a Splitter with the provided options.
func New(opts ...Option) *Splitter {
	s := &Splitter{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply processes a single log entry, returning a new map with split fields
// added. The original entry is not mutated.
func (s *Splitter) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry)+4)
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range s.rules {
		raw, ok := out[r.Source]
		if !ok {
			continue
		}
		str, ok := raw.(string)
		if !ok {
			continue
		}
		delim := r.Delimiter
		if delim == "" {
			delim = ":"
		}
		parts := strings.SplitN(str, delim, len(r.Targets)+1)
		for i, target := range r.Targets {
			if i < len(parts) {
				out[target] = parts[i]
			}
		}
	}
	return out
}
