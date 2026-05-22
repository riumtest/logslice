// Package fieldconditional provides a transformer that sets a destination
// field to a given value when a source field matches a condition.
package fieldconditional

import (
	"fmt"
	"strings"
)

// Rule describes a single conditional assignment.
type Rule struct {
	// SourceField is the field whose value is inspected.
	SourceField string
	// Op is the comparison operator: "eq", "neq", "contains", "gt", "lt".
	Op string
	// Operand is the value to compare against.
	Operand string
	// DestField is the field to set when the condition is true.
	DestField string
	// DestValue is the value written to DestField.
	DestValue string
}

// Transformer applies conditional field assignments to log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the supplied rules in order.
func New(rules []Rule) (*Transformer, error) {
	for i, r := range rules {
		if r.SourceField == "" || r.DestField == "" {
			return nil, fmt.Errorf("fieldconditional: rule %d: source and dest fields must not be empty", i)
		}
		switch r.Op {
		case "eq", "neq", "contains", "gt", "lt":
		default:
			return nil, fmt.Errorf("fieldconditional: rule %d: unknown op %q", i, r.Op)
		}
	}
	return &Transformer{rules: rules}, nil
}

// Apply evaluates each rule against entry and returns the (possibly modified) copy.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		val, ok := out[r.SourceField]
		if !ok {
			continue
		}
		sv := fmt.Sprintf("%v", val)
		if matches(r.Op, sv, r.Operand) {
			out[r.DestField] = r.DestValue
		}
	}
	return out
}

func matches(op, value, operand string) bool {
	switch op {
	case "eq":
		return value == operand
	case "neq":
		return value != operand
	case "contains":
		return strings.Contains(value, operand)
	case "gt":
		return value > operand
	case "lt":
		return value < operand
	}
	return false
}
