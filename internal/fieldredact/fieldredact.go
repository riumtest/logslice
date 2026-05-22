// Package fieldredact masks or removes sensitive field values from log entries.
package fieldredact

import "encoding/json"

const defaultMask = "[REDACTED]"

// Redactor replaces the values of named fields with a mask string.
type Redactor struct {
	fields map[string]struct{}
	mask   string
}

// Option configures a Redactor.
type Option func(*Redactor)

// WithMask overrides the default mask string.
func WithMask(mask string) Option {
	return func(r *Redactor) {
		if mask != "" {
			r.mask = mask
		}
	}
}

// New returns a Redactor that masks the supplied field names.
// If fields is empty the Redactor is a no-op.
func New(fields []string, opts ...Option) *Redactor {
	r := &Redactor{
		fields: make(map[string]struct{}, len(fields)),
		mask:   defaultMask,
	}
	for _, f := range fields {
		if f != "" {
			r.fields[f] = struct{}{}
		}
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Apply returns a copy of entry with sensitive fields replaced by the mask.
// Fields not present in the entry are silently ignored.
func (r *Redactor) Apply(entry map[string]json.RawMessage) map[string]json.RawMessage {
	if len(r.fields) == 0 {
		return entry
	}
	out := make(map[string]json.RawMessage, len(entry))
	for k, v := range entry {
		if _, redact := r.fields[k]; redact {
			masked, _ := json.Marshal(r.mask)
			out[k] = masked
		} else {
			out[k] = v
		}
	}
	return out
}
