package fieldclamp

import "fmt"

// Rule defines a clamp rule for a single field.
type Rule struct {
	Field string
	Min   *float64
	Max   *float64
}

// Clamper clamps numeric field values to optional min/max bounds.
type Clamper struct {
	rules []Rule
}

// New returns a Clamper that applies the given rules.
func New(rules []Rule) *Clamper {
	return &Clamper{rules: rules}
}

// Apply returns a copy of entry with each matching field clamped to its rule bounds.
func (c *Clamper) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range c.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		f, err := toFloat(v)
		if err != nil {
			continue
		}
		if r.Min != nil && f < *r.Min {
			f = *r.Min
		}
		if r.Max != nil && f > *r.Max {
			f = *r.Max
		}
		out[r.Field] = f
	}
	return out
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
	return 0, fmt.Errorf("fieldclamp: unsupported type %T", v)
}
