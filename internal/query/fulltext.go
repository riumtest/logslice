package query

import (
	"encoding/json"
	"strings"
)

// FullTextMatcher checks whether any value in a JSON log record contains
// all of the provided search terms (case-insensitive).
type FullTextMatcher struct {
	terms []string
}

// NewFullTextMatcher returns a FullTextMatcher for the given terms.
// Empty or whitespace-only terms are ignored.
func NewFullTextMatcher(terms []string) *FullTextMatcher {
	filtered := make([]string, 0, len(terms))
	for _, t := range terms {
		t = strings.TrimSpace(t)
		if t != "" {
			filtered = append(filtered, strings.ToLower(t))
		}
	}
	return &FullTextMatcher{terms: filtered}
}

// Match returns true when every term is found somewhere in the flattened
// string representation of the record's values.
func (f *FullTextMatcher) Match(record map[string]interface{}) bool {
	if len(f.terms) == 0 {
		return true
	}

	haystack := flattenValues(record)

	for _, term := range f.terms {
		if !strings.Contains(haystack, term) {
			return false
		}
	}
	return true
}

// HasTerms reports whether the matcher has at least one term.
func (f *FullTextMatcher) HasTerms() bool {
	return len(f.terms) > 0
}

// flattenValues converts all values in the record to a single lower-case
// string so terms can be searched across all fields at once.
func flattenValues(record map[string]interface{}) string {
	var sb strings.Builder
	for _, v := range record {
		switch val := v.(type) {
		case string:
			sb.WriteString(strings.ToLower(val))
		default:
			b, err := json.Marshal(val)
			if err == nil {
				sb.WriteString(strings.ToLower(string(b)))
			}
		}
		sb.WriteByte(' ')
	}
	return sb.String()
}
