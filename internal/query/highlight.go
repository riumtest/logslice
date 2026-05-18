package query

import (
	"strings"
)

// Highlighter holds a compiled set of terms to highlight inside a log line.
type Highlighter struct {
	terms []string
}

// NewHighlighter creates a Highlighter from a slice of search terms.
// Empty terms are ignored.
func NewHighlighter(terms []string) *Highlighter {
	filtered := make([]string, 0, len(terms))
	for _, t := range terms {
		if t != "" {
			filtered = append(filtered, t)
		}
	}
	return &Highlighter{terms: filtered}
}

// Wrap returns the input string with each matching term wrapped between the
// provided prefix and suffix strings (e.g. ANSI colour codes).
// Matching is case-insensitive.
func (h *Highlighter) Wrap(s, prefix, suffix string) string {
	if len(h.terms) == 0 {
		return s
	}
	for _, term := range h.terms {
		s = replaceInsensitive(s, term, prefix+term+suffix)
	}
	return s
}

// HasTerms reports whether the Highlighter has any terms to match.
func (h *Highlighter) HasTerms() bool {
	return len(h.terms) > 0
}

// replaceInsensitive replaces all case-insensitive occurrences of old in s
// with new.
func replaceInsensitive(s, old, newVal string) string {
	if old == "" {
		return s
	}
	lower := strings.ToLower(s)
	lowerOld := strings.ToLower(old)
	var b strings.Builder
	offset := 0
	for {
		idx := strings.Index(lower[offset:], lowerOld)
		if idx < 0 {
			b.WriteString(s[offset:])
			break
		}
		abs := offset + idx
		b.WriteString(s[offset:abs])
		b.WriteString(newVal)
		offset = abs + len(old)
	}
	return b.String()
}
