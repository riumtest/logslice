package fieldscope

import "encoding/json"

// Rule defines a scope operation: extract a nested JSON field into a top-level key.
type Rule struct {
	// Source is a dot-separated path into the entry, e.g. "meta.request.method".
	Source string
	// Dest is the top-level key to write the extracted value into.
	Dest string
}

// Transformer extracts nested fields from JSON objects and promotes them to
// top-level keys in the entry.
type Transformer struct {
	rules []Rule
}

// New returns a Transformer that applies the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply promotes nested fields to top-level keys according to the configured
// rules. The original entry is not mutated; a shallow copy is returned.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, r := range t.rules {
		val, ok := resolvePath(entry, splitPath(r.Source))
		if !ok {
			continue
		}
		out[r.Dest] = val
	}
	return out
}

// splitPath splits a dot-separated path into segments.
func splitPath(path string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			if i > start {
				parts = append(parts, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		parts = append(parts, path[start:])
	}
	return parts
}

// resolvePath walks into nested maps following the given key segments.
func resolvePath(node any, keys []string) (any, bool) {
	if len(keys) == 0 {
		return node, true
	}
	switch m := node.(type) {
	case map[string]any:
		v, ok := m[keys[0]]
		if !ok {
			return nil, false
		}
		return resolvePath(v, keys[1:])
	case json.RawMessage:
		var sub map[string]any
		if err := json.Unmarshal(m, &sub); err != nil {
			return nil, false
		}
		return resolvePath(sub, keys)
	}
	return nil, false
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
