// Package fieldbool normalises boolean-like field values into true JSON booleans.
//
// String values such as "true", "yes", "1", "on" (case-insensitive) are coerced
// to the boolean true; "false", "no", "0", "off" are coerced to false.
// Numeric values follow standard truthiness (non-zero → true).
// Fields that are already booleans are left unchanged.
package fieldbool

import (
	"fmt"
	"strings"
)

// Rule describes a single field-to-boolean coercion.
type Rule struct {
	// Field is the key whose value should be normalised.
	Field string
	// Dest, when non-empty, writes the result to a different key.
	Dest string
}

// Transformer coerces nominated fields to boolean values.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply returns a shallow copy of entry with coerced boolean fields.
// If no rules are configured the original entry is returned unchanged.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range t.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		b, err := coerce(v)
		if err != nil {
			continue
		}
		dest := r.Field
		if r.Dest != "" {
			dest = r.Dest
		}
		out[dest] = b
	}
	return out
}

func coerce(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "true", "yes", "1", "on":
			return true, nil
		case "false", "no", "0", "off":
			return false, nil
		}
		return false, fmt.Errorf("unrecognised boolean string: %q", val)
	case float64:
		return val != 0, nil
	case int:
		return val != 0, nil
	}
	return false, fmt.Errorf("unsupported type %T", v)
}

func shallowCopy(src map[string]any) map[string]any {
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
