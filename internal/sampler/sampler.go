// Package sampler provides rate-based and head/tail sampling for log streams.
package sampler

import (
	"math/rand"
)

// Mode controls which records are kept.
type Mode int

const (
	// ModeRate keeps approximately 1-in-N records.
	ModeRate Mode = iota
	// ModeHead keeps the first N records.
	ModeHead
	// ModeTail keeps the last N records.
	ModeTail
)

// Sampler decides which log records to pass through.
type Sampler struct {
	mode  Mode
	n     int
	count int
	buf   []map[string]any
	rng   *rand.Rand
}

// Option configures a Sampler.
type Option func(*Sampler)

// WithSeed sets the random seed used for rate sampling.
func WithSeed(seed int64) Option {
	return func(s *Sampler) {
		s.rng = rand.New(rand.NewSource(seed)) //nolint:gosec
	}
}

// New creates a Sampler with the given mode and N value.
func New(mode Mode, n int, opts ...Option) *Sampler {
	if n < 1 {
		n = 1
	}
	s := &Sampler{
		mode: mode,
		n:    n,
		rng:  rand.New(rand.NewSource(42)), //nolint:gosec
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Feed processes a single record and returns (record, true) if it should be
// emitted immediately, or (nil, false) if it is buffered or dropped.
func (s *Sampler) Feed(record map[string]any) (map[string]any, bool) {
	switch s.mode {
	case ModeRate:
		s.count++
		if s.rng.Intn(s.n) == 0 {
			return record, true
		}
		return nil, false
	case ModeHead:
		s.count++
		if s.count <= s.n {
			return record, true
		}
		return nil, false
	case ModeTail:
		s.buf = append(s.buf, record)
		if len(s.buf) > s.n {
			s.buf = s.buf[1:]
		}
		return nil, false
	}
	return record, true
}

// Flush returns any buffered records (only meaningful for ModeTail).
func (s *Sampler) Flush() []map[string]any {
	out := s.buf
	s.buf = nil
	return out
}
