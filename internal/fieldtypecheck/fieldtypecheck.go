package fieldtypecheck

import "fmt"

// Rule defines a field and the expected JSON type for that field.
type Rule struct {
	Field    string
	Expected string // "string", "number", "bool", "null", "array", "object"
}

// Checker filters or annotates log entries based on field type expectations.
type Checker struct {
	rules      []Rule
	destField  string
	rejectMode bool
}

// WithDestField sets the field name where a type-error summary is written.
// When empty, entries failing a check are dropped (reject mode).
func WithDestField(f string) func(*Checker) {
	return func(c *Checker) { c.destField = f }
}

// WithRejectMode causes entries that fail any type check to be dropped.
func WithRejectMode() func(*Checker) {
	return func(c *Checker) { c.rejectMode = true }
}

// New creates a Checker with the given rules and options.
func New(rules []Rule, opts ...func(*Checker)) *Checker {
	c := &Checker{rules: rules}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Apply checks each rule against the entry. If a type mismatch is found and
// destField is set, the error summary is written there. If rejectMode is
// true, a mismatched entry is dropped (nil returned).
func (c *Checker) Apply(entry map[string]any) map[string]any {
	if len(c.rules) == 0 {
		return entry
	}
	var errs []string
	for _, r := range c.rules {
		v, ok := entry[r.Field]
		if !ok {
			continue
		}
		if got := jsonType(v); got != r.Expected {
			errs = append(errs, fmt.Sprintf("%s: want %s got %s", r.Field, r.Expected, got))
		}
	}
	if len(errs) == 0 {
		return entry
	}
	if c.rejectMode {
		return nil
	}
	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}
	if c.destField != "" {
		out[c.destField] = errs
	}
	return out
}

func jsonType(v any) string {
	if v == nil {
		return "null"
	}
	switch v.(type) {
	case string:
		return "string"
	case float64, int, int64:
		return "number"
	case bool:
		return "bool"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return "unknown"
	}
}
