// Package source provides utilities for resolving log input sources
// from file paths, stdin, or named pipes.
package source

import (
	"fmt"
	"io"
	"os"
)

// Resolver opens one or more named sources and returns their readers.
type Resolver struct {
	stdin io.Reader
}

// Option configures a Resolver.
type Option func(*Resolver)

// WithStdin overrides the default os.Stdin reader (useful for testing).
func WithStdin(r io.Reader) Option {
	return func(rs *Resolver) {
		rs.stdin = r
	}
}

// New creates a new Resolver with the given options.
func New(opts ...Option) *Resolver {
	rs := &Resolver{stdin: os.Stdin}
	for _, o := range opts {
		o(rs)
	}
	return rs
}

// Resolve returns a slice of ReadClosers for the given paths.
// If paths is empty, stdin is returned as the sole source.
// Callers are responsible for closing each returned ReadCloser.
func (rs *Resolver) Resolve(paths []string) ([]io.ReadCloser, error) {
	if len(paths) == 0 {
		return []io.ReadCloser{io.NopCloser(rs.stdin)}, nil
	}

	out := make([]io.ReadCloser, 0, len(paths))
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			// Close any already-opened files before returning.
			for _, rc := range out {
				_ = rc.Close()
			}
			return nil, fmt.Errorf("source: open %q: %w", p, err)
		}
		out = append(out, f)
	}
	return out, nil
}
