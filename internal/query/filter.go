// Package query provides a simple DSL for filtering structured JSON log entries.
package query

import (
	"fmt"
	"strconv"
	"strings"
)

// Op represents a comparison operator.
type Op string

const (
	OpEq  Op = "="
	OpNeq Op = "!="
	OpGt  Op = ">"
	OpLt  Op = "<"
	OpGte Op = ">="
	OpLte Op = "<="
	OpContains Op = "~"
)

// Filter represents a single field comparison expression.
type Filter struct {
	Field string
	Op    Op
	Value string
}

// Match reports whether the given log entry (as a flat map) satisfies the filter.
func (f *Filter) Match(entry map[string]any) bool {
	raw, ok := entry[f.Field]
	if !ok {
		return false
	}

	switch f.Op {
	case OpEq:
		return fmt.Sprintf("%v", raw) == f.Value
	case OpNeq:
		return fmt.Sprintf("%v", raw) != f.Value
	case OpContains:
		return strings.Contains(fmt.Sprintf("%v", raw), f.Value)
	case OpGt, OpLt, OpGte, OpLte:
		return compareNumeric(raw, f.Op, f.Value)
	}
	return false
}

func compareNumeric(raw any, op Op, value string) bool {
	var left float64
	switch v := raw.(type) {
	case float64:
		left = v
	case int:
		left = float64(v)
	case string:
		var err error
		left, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
	default:
		return false
	}

	right, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}

	switch op {
	case OpGt:
		return left > right
	case OpLt:
		return left < right
	case OpGte:
		return left >= right
	case OpLte:
		return left <= right
	}
	return false
}
