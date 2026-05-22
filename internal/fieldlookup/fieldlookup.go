// Package fieldlookup provides a transformer that replaces a field's value
// using a static lookup table, optionally writing the result to a separate
// destination field and falling back to a default when no match is found.
package fieldlookup

import "github.com/logslice/logslice/internal/pipeline"

// Option configures a Transformer.
type Option func(*Transformer)

// WithDefault sets the fallback value used when the source field value is not
// found in the lookup table. When not set the entry is left unchanged.
func WithDefault(val string) Option {
	return func(t *Transformer) { t.defaultVal = &val }
}

// WithDestField writes the looked-up value to destField instead of overwriting
// the source field.
func WithDestField(dest string) Option {
	return func(t *Transformer) { t.destField = dest }
}

// Transformer replaces field values via a lookup table.
type Transformer struct {
	srcField   string
	table      map[string]string
	destField  string
	defaultVal *string
}

// New returns a Transformer that maps values in srcField through table.
func New(srcField string, table map[string]string, opts ...Option) *Transformer {
	t := &Transformer{
		srcField: srcField,
		table:    table,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Transform applies the lookup to a single log entry.
func (t *Transformer) Transform(entry pipeline.Entry) pipeline.Entry {
	raw, ok := entry[t.srcField]
	if !ok {
		return entry
	}

	srcStr, ok := raw.(string)
	if !ok {
		return entry
	}

	mapped, found := t.table[srcStr]
	if !found {
		if t.defaultVal == nil {
			return entry
		}
		mapped = *t.defaultVal
	}

	out := make(pipeline.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	dest := t.srcField
	if t.destField != "" {
		dest = t.destField
	}
	out[dest] = mapped
	return out
}
