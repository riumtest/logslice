// Package fieldtruncate truncates string fields to a maximum byte or rune length.
package fieldtruncate

import "unicode/utf8"

// Rule describes a single truncation rule.
type Rule struct {
	// Field is the JSON key to truncate.
	Field string
	// MaxLen is the maximum number of runes allowed. Values longer than this
	// are truncated and a suffix is appended.
	MaxLen int
	// Suffix is appended when the value is truncated (default: "...").
	Suffix string
}

// Transformer truncates string fields according to a set of rules.
type Transformer struct {
	rules []Rule
}

// WithRules returns a new Transformer configured with the given rules.
func WithRules(rules []Rule) *Transformer {
	normalized := make([]Rule, len(rules))
	for i, r := range rules {
		if r.Suffix == "" {
			r.Suffix = "..."
		}
		if r.MaxLen < 0 {
			r.MaxLen = 0
		}
		normalized[i] = r
	}
	return &Transformer{rules: normalized}
}

// New returns a Transformer with no rules (identity transform).
func New() *Transformer {
	return &Transformer{}
}

// Apply applies all truncation rules to the entry and returns a new map.
// The original entry is never mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		val, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if utf8.RuneCountInString(s) > r.MaxLen {
			out[r.Field] = truncateRunes(s, r.MaxLen) + r.Suffix
		}
	}
	return out
}

// truncateRunes returns the first n runes of s.
func truncateRunes(s string, n int) string {
	count := 0
	for i := range s {
		if count == n {
			return s[:i]
		}
		count++
	}
	return s
}
