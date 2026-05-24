// Package fieldxpath extracts values from nested JSON objects using
// dot-notation path expressions and writes them to a destination field.
package fieldxpath

import (
	"strings"
)

// Rule defines a single extraction: read from Path (dot-separated) and write
// the resolved value into Dest.
type Rule struct {
	Path string
	Dest string
}

// Transformer extracts nested values via dot-notation paths.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer configured with the provided rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply resolves each rule's Path in entry and stores the result in Dest.
// The original entry is never mutated; a shallow copy is returned.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range t.rules {
		if r.Path == "" || r.Dest == "" {
			continue
		}
		val, ok := resolve(out, strings.Split(r.Path, "."))
		if !ok {
			continue
		}
		out[r.Dest] = val
	}
	return out
}

// resolve walks the parts of a dot-split path through nested maps.
func resolve(obj map[string]any, parts []string) (any, bool) {
	if len(parts) == 0 {
		return nil, false
	}
	val, ok := obj[parts[0]]
	if !ok {
		return nil, false
	}
	if len(parts) == 1 {
		return val, true
	}
	nested, ok := val.(map[string]any)
	if !ok {
		return nil, false
	}
	return resolve(nested, parts[1:])
}

func shallowCopy(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
