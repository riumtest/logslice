// Package fieldparse provides a transformer that parses string fields
// into typed JSON values (numbers, booleans, or nested JSON objects).
package fieldparse

import (
	"encoding/json"
	"strconv"
)

// Transformer parses nominated string fields into their inferred types.
type Transformer struct {
	fields map[string]struct{}
}

// Option configures a Transformer.
type Option func(*Transformer)

// WithFields specifies which field names should be parsed.
func WithFields(names ...string) Option {
	return func(t *Transformer) {
		for _, n := range names {
			t.fields[n] = struct{}{}
		}
	}
}

// New creates a Transformer with the supplied options.
func New(opts ...Option) *Transformer {
	t := &Transformer{fields: make(map[string]struct{})}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Transform returns a copy of entry with targeted string fields parsed into
// their natural types. Fields that cannot be parsed are left unchanged.
func (t *Transformer) Transform(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	for k := range t.fields {
		raw, ok := out[k]
		if !ok {
			continue
		}
		s, ok := raw.(string)
		if !ok {
			continue
		}
		out[k] = parseValue(s)
	}
	return out
}

// parseValue attempts to decode s as bool, float64, or a JSON object/array;
// falls back to the original string on failure.
func parseValue(s string) any {
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if len(s) > 0 && (s[0] == '{' || s[0] == '[') {
		var v any
		if err := json.Unmarshal([]byte(s), &v); err == nil {
			return v
		}
	}
	return s
}
