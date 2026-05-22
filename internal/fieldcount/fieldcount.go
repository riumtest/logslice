// Package fieldcount provides a transformer that adds a field containing
// the total number of fields present in each log entry.
package fieldcount

import "encoding/json"

// Transformer adds a count field to each log entry reflecting the number
// of keys present in that entry (before the count field is added).
type Transformer struct {
	destField string
}

// Option is a functional option for Transformer.
type Option func(*Transformer)

// WithDestField sets the name of the field written with the count value.
// Defaults to "_field_count".
func WithDestField(name string) Option {
	return func(t *Transformer) {
		if name != "" {
			t.destField = name
		}
	}
}

// New creates a Transformer with the given options.
func New(opts ...Option) *Transformer {
	t := &Transformer{destField: "_field_count"}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Transform returns a copy of entry with the count field injected.
// If entry is nil or empty the original map is returned unchanged.
func (t *Transformer) Transform(entry map[string]json.RawMessage) map[string]json.RawMessage {
	if len(entry) == 0 {
		return entry
	}

	count := len(entry)

	out := make(map[string]json.RawMessage, count+1)
	for k, v := range entry {
		out[k] = v
	}

	raw, err := json.Marshal(count)
	if err != nil {
		return out
	}
	out[t.destField] = json.RawMessage(raw)
	return out
}
