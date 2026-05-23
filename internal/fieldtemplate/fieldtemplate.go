// Package fieldtemplate provides a transformer that renders Go text/template
// expressions into a target field using values from the current log entry.
package fieldtemplate

import (
	"bytes"
	"fmt"
	"text/template"
)

// Rule describes a single template rule: the destination field and the
// Go template string that produces its value.
type Rule struct {
	DestField string
	Tmpl      string
}

// Transformer applies template rules to each log entry.
type Transformer struct {
	rules []compiledRule
}

type compiledRule struct {
	dest string
	tmpl *template.Template
}

// WithRules returns a Transformer configured with the given rules.
// Returns an error if any template fails to parse.
func WithRules(rules []Rule) (*Transformer, error) {
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		tmpl, err := template.New(r.DestField).Option("missingkey=zero").Parse(r.Tmpl)
		if err != nil {
			return nil, fmt.Errorf("fieldtemplate: parse rule %q: %w", r.DestField, err)
		}
		compiled = append(compiled, compiledRule{dest: r.DestField, tmpl: tmpl})
	}
	return &Transformer{rules: compiled}, nil
}

// New returns a no-op Transformer (no rules).
func New() *Transformer {
	return &Transformer{}
}

// Apply renders each rule template against the entry and stores the result
// in the destination field. The original entry is not mutated.
func (t *Transformer) Apply(entry map[string]any) map[string]any {
	if len(t.rules) == 0 {
		return entry
	}
	out := make(map[string]any, len(entry)+len(t.rules))
	for k, v := range entry {
		out[k] = v
	}
	var buf bytes.Buffer
	for _, r := range t.rules {
		buf.Reset()
		if err := r.tmpl.Execute(&buf, entry); err != nil {
			continue
		}
		out[r.dest] = buf.String()
	}
	return out
}
