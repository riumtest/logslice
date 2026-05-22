package fielddefault

import "encoding/json"

// Rule defines a field name and the default value to set when the field is
// absent or null in a log entry.
type Rule struct {
	Field string
	Value any
}

// Transformer sets default values for missing or null fields in log entries.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the given default rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Transform applies default values to the entry. Fields that are already
// present and non-null are left unchanged. A copy of the entry is returned.
func (t *Transformer) Transform(entry map[string]json.RawMessage) map[string]json.RawMessage {
	if len(t.rules) == 0 {
		return entry
	}
	out := make(map[string]json.RawMessage, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		existing, ok := out[r.Field]
		if ok && string(existing) != "null" {
			continue
		}
		raw, err := json.Marshal(r.Value)
		if err != nil {
			continue
		}
		out[r.Field] = json.RawMessage(raw)
	}
	return out
}
