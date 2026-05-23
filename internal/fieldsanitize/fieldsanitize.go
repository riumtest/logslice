// Package fieldsanitize trims whitespace and optionally removes control
// characters from string fields in a log entry.
package fieldsanitize

import (
	"strings"
	"unicode"
)

// Rule describes how a single field should be sanitized.
type Rule struct {
	// Field is the JSON key to sanitize.
	Field string
	// StripControl removes non-printable / control characters when true.
	StripControl bool
}

// Transformer sanitizes string fields according to a set of rules.
type Transformer struct {
	rules []Rule
}

// WithRules returns a new Transformer configured with the given rules.
func WithRules(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// New returns a Transformer with no rules (identity transform).
func New() *Transformer {
	return &Transformer{}
}

// Apply sanitizes each configured field in the entry and returns a new map.
// Fields that are not strings, or not present, are left untouched.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}

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
		s = strings.TrimSpace(s)
		if r.StripControl {
			s = stripControl(s)
		}
		out[r.Field] = s
	}

	return out
}

// stripControl removes runes that are control characters (except common
// whitespace that TrimSpace already handles).
func stripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return -1
		}
		return r
	}, s)
}
