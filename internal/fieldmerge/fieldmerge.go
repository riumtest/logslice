package fieldmerge

// Merger merges values from multiple source fields into a single destination
// field, joining them with a configurable separator.

import (
	"fmt"
	"strings"
)

// Option configures a Merger.
type Option func(*Merger)

// WithSeparator sets the string used to join merged values (default: " ").
func WithSeparator(sep string) Option {
	return func(m *Merger) {
		m.separator = sep
	}
}

// WithSkipMissing controls whether missing source fields are silently skipped
// (default: true). When false, a missing field causes an error placeholder.
func WithSkipMissing(skip bool) Option {
	return func(m *Merger) {
		m.skipMissing = skip
	}
}

// Merger combines values from multiple source fields into a destination field.
type Merger struct {
	sources     []string
	dest        string
	separator   string
	skipMissing bool
}

// New creates a Merger that reads from sources and writes to dest.
func New(dest string, sources []string, opts ...Option) *Merger {
	m := &Merger{
		sources:     sources,
		dest:        dest,
		separator:   " ",
		skipMissing: true,
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Apply merges the configured source fields into the destination field and
// returns a new map. The original entry is not mutated.
func (m *Merger) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry)+1)
	for k, v := range entry {
		out[k] = v
	}

	parts := make([]string, 0, len(m.sources))
	for _, src := range m.sources {
		v, ok := entry[src]
		if !ok {
			if !m.skipMissing {
				parts = append(parts, "<missing:" + src + ">")
			}
			continue
		}
		parts = append(parts, fmt.Sprintf("%v", v))
	}

	out[m.dest] = strings.Join(parts, m.separator)
	return out
}
