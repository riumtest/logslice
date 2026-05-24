// Package fieldmatch filters or annotates log entries based on whether
// a field's value matches a regular expression, writing a boolean result
// to a configurable destination field.
package fieldmatch

import (
	"regexp"
)

// Rule describes a single match rule: the source field, the compiled
// pattern, and where to write the boolean result.
type Rule struct {
	Field   string
	Pattern *regexp.Regexp
	Dest    string
}

// Transformer applies a set of match rules to each log entry.
type Transformer struct {
	rules []Rule
}

// WithRules returns a new Transformer configured with the given rules.
func WithRules(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// New returns a Transformer with no rules (identity).
func New() *Transformer {
	return &Transformer{}
}

// Apply evaluates each rule against the entry and writes a boolean to the
// destination field. If the source field is absent or not a string the
// destination is set to false. The original entry is never mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range t.rules {
		v, ok := out[r.Field]
		if !ok {
			out[r.Dest] = false
			continue
		}
		s, ok := v.(string)
		if !ok {
			out[r.Dest] = false
			continue
		}
		out[r.Dest] = r.Pattern.MatchString(s)
	}
	return out
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
