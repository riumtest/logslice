package fieldmath

import (
	"fmt"
	"math"
)

// Op represents a binary math operation.
type Op string

const (
	OpAdd Op = "add"
	OpSub Op = "sub"
	OpMul Op = "mul"
	OpDiv Op = "div"
)

// Rule defines a math operation: Dest = Left <Op> Right.
type Rule struct {
	Left  string
	Right string
	Dest  string
	Op    Op
}

// Transformer applies arithmetic rules to log entries.
type Transformer struct {
	rules []Rule
}

// New creates a Transformer with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply applies all rules to a copy of entry and returns the result.
// Rules whose operand fields are missing or non-numeric are skipped.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		lv, lok := toFloat(out[r.Left])
		rv, rok := toFloat(out[r.Right])
		if !lok || !rok {
			continue
		}
		result, err := compute(lv, rv, r.Op)
		if err != nil {
			continue
		}
		out[r.Dest] = result
	}
	return out
}

func compute(l, r float64, op Op) (float64, error) {
	switch op {
	case OpAdd:
		return l + r, nil
	case OpSub:
		return l - r, nil
	case OpMul:
		return l * r, nil
	case OpDiv:
		if r == 0 {
			return math.NaN(), fmt.Errorf("division by zero")
		}
		return l / r, nil
	default:
		return 0, fmt.Errorf("unknown op: %s", op)
	}
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}
