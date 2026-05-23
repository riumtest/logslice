package pipeline

import (
	"fmt"

	"github.com/user/logslice/internal/fieldregex"
)

// FieldRegexRule is the pipeline-level configuration for a regex extraction rule.
type FieldRegexRule struct {
	Source  string
	Pattern string
	Prefix  string
}

// validate checks that the rule has the required fields set.
func (r FieldRegexRule) validate() error {
	if r.Source == "" {
		return fmt.Errorf("fieldregex rule missing required field: source")
	}
	if r.Pattern == "" {
		return fmt.Errorf("fieldregex rule for source %q missing required field: pattern", r.Source)
	}
	return nil
}

// buildFieldRegexExtractor constructs a fieldregex.Extractor from pipeline config,
// returning nil if no rules are defined.
func buildFieldRegexExtractor(rules []FieldRegexRule) (*fieldregex.Extractor, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	inner := make([]fieldregex.Rule, len(rules))
	for i, r := range rules {
		if err := r.validate(); err != nil {
			return nil, fmt.Errorf("invalid fieldregex rule at index %d: %w", i, err)
		}
		inner[i] = fieldregex.Rule{
			Source:  r.Source,
			Pattern: r.Pattern,
			Prefix:  r.Prefix,
		}
	}
	return fieldregex.New(fieldregex.WithRules(inner))
}
