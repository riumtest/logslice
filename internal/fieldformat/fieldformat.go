// Package fieldformat applies formatting transformations to string fields
// in a log entry, such as upper-casing, lower-casing, or truncating values.
package fieldformat

import (
	"strings"

	"github.com/mikelorant/logslice/internal/entry"
)

// Op represents a formatting operation to apply to a field value.
type Op int

const (
	OpUppercase Op = iota
	OpLowercase
	OpTruncate
)

// Rule describes a single field formatting rule.
type Rule struct {
	Field  string
	Op     Op
	MaxLen int // used only for OpTruncate
}

// Formatter applies a set of Rules to log entries.
type Formatter struct {
	rules []Rule
}

// WithRules returns a new Formatter configured with the provided rules.
func WithRules(rules []Rule) *Formatter {
	return &Formatter{rules: rules}
}

// New returns a Formatter with no rules (identity transform).
func New() *Formatter {
	return &Formatter{}
}

// Apply returns a new entry with the formatting rules applied.
// Fields not matched by any rule are left unchanged.
func (f *Formatter) Apply(e entry.Entry) entry.Entry {
	out := make(entry.Entry, len(e))
	for k, v := range e {
		out[k] = v
	}

	for _, r := range f.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		switch r.Op {
		case OpUppercase:
			out[r.Field] = strings.ToUpper(s)
		case OpLowercase:
			out[r.Field] = strings.ToLower(s)
		case OpTruncate:
			if r.MaxLen > 0 && len(s) > r.MaxLen {
				out[r.Field] = s[:r.MaxLen]
			}
		}
	}
	return out
}
