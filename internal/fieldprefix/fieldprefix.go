package fieldprefix

// Package fieldprefix provides a transformer that adds or strips a string
// prefix from the values of nominated string fields.

import "fmt"

// Rule describes a single prefix operation for one field.
type Rule struct {
	// Field is the JSON key to operate on.
	Field string
	// Prefix is the string to add or remove.
	Prefix string
	// Strip removes the prefix instead of adding it when true.
	Strip bool
}

// Transformer applies prefix rules to log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer configured with the supplied rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Transform applies each rule to a copy of entry and returns the result.
// Fields that are absent or not strings are left unchanged.
func (t *Transformer) Transform(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	for _, r := range t.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if r.Strip {
			if len(s) >= len(r.Prefix) && s[:len(r.Prefix)] == r.Prefix {
				out[r.Field] = s[len(r.Prefix):]
			}
		} else {
			out[r.Field] = fmt.Sprintf("%s%s", r.Prefix, s)
		}
	}

	return out
}
