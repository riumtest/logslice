package fieldencoding

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// Rule describes how to encode or decode a single field.
type Rule struct {
	Field  string
	Mode   string // "base64-encode", "base64-decode", "hex-encode", "hex-decode", "url-encode", "url-decode"
	Dest   string // optional; defaults to Field
}

// Transformer applies encoding/decoding rules to log entries.
type Transformer struct {
	rules []Rule
}

// WithRules returns a new Transformer configured with the given rules.
func WithRules(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// New returns a Transformer with no rules (identity).
func New() *Transformer { return &Transformer{} }

// Apply encodes or decodes fields according to the configured rules.
func (t *Transformer) Apply(entry map[string]any) (map[string]any, error) {
	if len(t.rules) == 0 {
		return entry, nil
	}
	out := shallowCopy(entry)
	for _, r := range t.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		result, err := applyMode(r.Mode, s)
		if err != nil {
			return nil, fmt.Errorf("fieldencoding: field %q mode %q: %w", r.Field, r.Mode, err)
		}
		dest := r.Dest
		if dest == "" {
			dest = r.Field
		}
		out[dest] = result
	}
	return out, nil
}

func applyMode(mode, s string) (string, error) {
	switch strings.ToLower(mode) {
	case "base64-encode":
		return base64.StdEncoding.EncodeToString([]byte(s)), nil
	case "base64-decode":
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "hex-encode":
		return hex.EncodeToString([]byte(s)), nil
	case "hex-decode":
		b, err := hex.DecodeString(s)
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("unknown mode %q", mode)
	}
}

func shallowCopy(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
