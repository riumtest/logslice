// Package fieldmask provides filtering of JSON log entry fields,
// allowing callers to include or exclude specific keys from output.
package fieldmask

import "strings"

// Mask holds the configuration for field inclusion/exclusion.
type Mask struct {
	include map[string]struct{}
	exclude map[string]struct{}
}

// Option configures a Mask.
type Option func(*Mask)

// WithInclude restricts output to only the specified fields.
// An empty slice is a no-op.
func WithInclude(fields []string) Option {
	return func(m *Mask) {
		for _, f := range fields {
			if f = strings.TrimSpace(f); f != "" {
				m.include[f] = struct{}{}
			}
		}
	}
}

// WithExclude removes the specified fields from output.
// An empty slice is a no-op.
func WithExclude(fields []string) Option {
	return func(m *Mask) {
		for _, f := range fields {
			if f = strings.TrimSpace(f); f != "" {
				m.exclude[f] = struct{}{}
			}
		}
	}
}

// New creates a Mask with the given options.
func New(opts ...Option) *Mask {
	m := &Mask{
		include: make(map[string]struct{}),
		exclude: make(map[string]struct{}),
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Apply returns a copy of entry with the mask applied.
// Include takes precedence: if any include fields are set, only those keys
// are kept (minus any also listed in exclude).
func (m *Mask) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))

	for k, v := range entry {
		if len(m.include) > 0 {
			if _, ok := m.include[k]; !ok {
				continue
			}
		}
		if _, ok := m.exclude[k]; ok {
			continue
		}
		out[k] = v
	}
	return out
}

// IsIdentity returns true when the mask will not alter any entry.
func (m *Mask) IsIdentity() bool {
	return len(m.include) == 0 && len(m.exclude) == 0
}
