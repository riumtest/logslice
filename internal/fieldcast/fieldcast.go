// Package fieldcast converts field values to a target type (string, int, float, bool).
package fieldcast

import (
	"fmt"
	"strconv"
)

// Rule describes a single field cast operation.
type Rule struct {
	Field  string
	Target string // "string", "int", "float", "bool"
}

// Caster applies type-cast rules to log entries.
type Caster struct {
	rules []Rule
}

// WithRules returns a new Caster configured with the given rules.
func WithRules(rules []Rule) *Caster {
	return &Caster{rules: rules}
}

// New returns a Caster with no rules (identity transform).
func New() *Caster {
	return &Caster{}
}

// Apply returns a copy of entry with each rule's field cast to the target type.
// Fields that are absent or cannot be converted are left unchanged.
func (c *Caster) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range c.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		casted, err := castValue(v, r.Target)
		if err != nil {
			continue
		}
		out[r.Field] = casted
	}
	return out
}

func castValue(v any, target string) (any, error) {
	raw := fmt.Sprintf("%v", v)
	switch target {
	case "string":
		return raw, nil
	case "int":
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			// try via float truncation
			f, err2 := strconv.ParseFloat(raw, 64)
			if err2 != nil {
				return nil, err
			}
			return int64(f), nil
		}
		return n, nil
	case "float":
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "bool":
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return nil, fmt.Errorf("unknown target type: %s", target)
	}
}
