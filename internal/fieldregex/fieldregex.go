package fieldregex

import (
	"fmt"
	"regexp"
)

// Rule defines a regex extraction rule: read Source, apply Pattern,
// and write named capture groups as new fields (optionally prefixed).
type Rule struct {
	Source  string
	Pattern string
	Prefix  string
	re      *regexp.Regexp
}

// Extractor applies regex extraction rules to log entries.
type Extractor struct {
	rules []Rule
}

// WithRules returns an option function that sets the extraction rules.
func WithRules(rules []Rule) func(*Extractor) error {
	return func(e *Extractor) error {
		for i, r := range rules {
			if r.Source == "" {
				return fmt.Errorf("rule %d: source field must not be empty", i)
			}
			if r.Pattern == "" {
				return fmt.Errorf("rule %d: pattern must not be empty", i)
			}
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return fmt.Errorf("rule %d: invalid pattern: %w", i, err)
			}
			rules[i].re = re
		}
		e.rules = rules
		return nil
	}
}

// New creates a new Extractor with the provided options.
func New(opts ...func(*Extractor) error) (*Extractor, error) {
	e := &Extractor{}
	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

// Apply runs all extraction rules against entry, returning a new map
// with extracted fields merged in. The original entry is not mutated.
func (e *Extractor) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, rule := range e.rules {
		src, ok := entry[rule.Source]
		if !ok {
			continue
		}
		str, ok := src.(string)
		if !ok {
			continue
		}
		match := rule.re.FindStringSubmatch(str)
		if match == nil {
			continue
		}
		for i, name := range rule.re.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			key := rule.Prefix + name
			out[key] = match[i]
		}
	}
	return out
}
