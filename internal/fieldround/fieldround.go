package fieldround

import (
	"fmt"
	"math"
)

// Rule defines how a numeric field should be rounded.
type Rule struct {
	Field     string
	Mode      string // "round", "floor", "ceil"
	Precision int    // decimal places (0 = integer)
}

// Transformer rounds numeric fields according to configured rules.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer applying the given rounding rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply returns a new entry with numeric fields rounded per the rules.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		f, err := toFloat(v)
		if err != nil {
			continue
		}
		out[r.Field] = applyMode(f, r.Mode, r.Precision)
	}
	return out
}

func applyMode(f float64, mode string, precision int) float64 {
	mult := math.Pow(10, float64(precision))
	switch mode {
	case "floor":
		return math.Floor(f*mult) / mult
	case "ceil":
		return math.Ceil(f*mult) / mult
	default: // "round"
		return math.Round(f*mult) / mult
	}
}

func toFloat(v any) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case float32:
		return float64(n), nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	}
	return 0, fmt.Errorf("not numeric")
}
