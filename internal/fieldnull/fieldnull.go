// Package fieldnull provides a transformer that drops or replaces null
// (JSON null / nil) values in log entry fields.
package fieldnull

// Rule describes how to handle a null value for a single field.
type Rule struct {
	// Field is the entry key to inspect.
	Field string
	// Replace is the value written when the field is null or missing.
	// When Replace is nil the field is removed from the entry entirely.
	Replace interface{}
	// DropEntry causes the whole entry to be discarded when the field is null.
	DropEntry bool
}

// Transformer applies null-handling rules to log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the given rules in order.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Transform applies each rule to entry and returns the (possibly modified)
// entry. A nil return value signals that the entry should be dropped.
func (t *Transformer) Transform(entry map[string]interface{}) map[string]interface{} {
	if len(t.rules) == 0 {
		return entry
	}

	out := make(map[string]interface{}, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	for _, r := range t.rules {
		v, exists := out[r.Field]
		if !exists || v == nil {
			if r.DropEntry {
				return nil
			}
			if r.Replace != nil {
				out[r.Field] = r.Replace
			} else {
				delete(out, r.Field)
			}
		}
	}

	return out
}
