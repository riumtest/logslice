// Package dedupe provides a filter that suppresses consecutive duplicate
// log entries based on a configurable set of fields.
package dedupe

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
)

// Deduper suppresses repeated log entries that share the same fingerprint.
type Deduper struct {
	fields []string // empty means use entire entry
	last   string
}

// Option is a functional option for Deduper.
type Option func(*Deduper)

// WithFields restricts fingerprinting to the given field names.
func WithFields(fields ...string) Option {
	return func(d *Deduper) {
		copied := make([]string, len(fields))
		copy(copied, fields)
		sort.Strings(copied)
		d.fields = copied
	}
}

// New creates a new Deduper with the provided options.
func New(opts ...Option) *Deduper {
	d := &Deduper{}
	for _, o := range opts {
		o(d)
	}
	return d
}

// IsDuplicate returns true when entry is a duplicate of the previous entry.
// Consecutive duplicates are suppressed; non-consecutive repeats are allowed.
func (d *Deduper) IsDuplicate(entry map[string]any) bool {
	fp := d.fingerprint(entry)
	if fp == d.last {
		return true
	}
	d.last = fp
	return false
}

// Reset clears the stored fingerprint so the next entry is never a duplicate.
func (d *Deduper) Reset() {
	d.last = ""
}

func (d *Deduper) fingerprint(entry map[string]any) string {
	var src any
	if len(d.fields) == 0 {
		src = entry
	} else {
		sub := make(map[string]any, len(d.fields))
		for _, f := range d.fields {
			if v, ok := entry[f]; ok {
				sub[f] = v
			}
		}
		src = sub
	}
	b, err := json.Marshal(src)
	if err != nil {
		return fmt.Sprintf("%v", src)
	}
	sum := sha256.Sum256(b)
	return fmt.Sprintf("%x", sum)
}
