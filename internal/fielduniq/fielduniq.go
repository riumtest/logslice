// Package fielduniq deduplicates values within a repeated string field.
// Given a field whose value is a JSON array of strings, it removes duplicate
// elements while preserving the original order of first occurrence.
package fielduniq

import "encoding/json"

// Transformer removes duplicate values from array-typed string fields.
type Transformer struct {
	fields []string
}

// New returns a Transformer that deduplicates values in the given fields.
// Fields whose values are not JSON arrays of strings are left untouched.
func New(fields []string) *Transformer {
	return &Transformer{fields: fields}
}

// Apply returns a copy of entry with duplicate array elements removed from
// each configured field. The original entry is never mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.fields) == 0 {
		return entry
	}
	out := shallowCopy(entry)
	for _, field := range t.fields {
		val, ok := out[field]
		if !ok {
			continue
		}
		deduped, changed := deduplicateArray(val)
		if changed {
			out[field] = deduped
		}
	}
	return out
}

// deduplicateArray attempts to interpret val as a []any whose elements are
// strings (or json.Number). It returns the deduplicated slice and whether
// any change was made.
func deduplicateArray(val any) ([]any, bool) {
	slice, ok := val.([]any)
	if !ok {
		return nil, false
	}
	seen := make(map[string]struct{}, len(slice))
	result := make([]any, 0, len(slice))
	changed := false
	for _, elem := range slice {
		key := stringify(elem)
		if _, dup := seen[key]; dup {
			changed = true
			continue
		}
		seen[key] = struct{}{}
		result = append(result, elem)
	}
	return result, changed
}

func stringify(v any) string {
	switch s := v.(type) {
	case string:
		return s
	case json.Number:
		return s.String()
	default:
		return ""
	}
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
