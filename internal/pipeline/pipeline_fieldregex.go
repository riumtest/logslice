package pipeline

import (
	"github.com/user/logslice/internal/fieldregex"
)

// FieldRegexRule is the pipeline-level configuration for a regex extraction rule.
type FieldRegexRule struct {
	Source  string
	Pattern string
	Prefix  string
}

// buildFieldRegexExtractor constructs a fieldregex.Extractor from pipeline config,
// returning nil if no rules are defined.
func buildFieldRegexExtractor(rules []FieldRegexRule) (*fieldregex.Extractor, error) {
	if len(rules) == 0 {
		return nil, nil
	}
	inner := make([]fieldregex.Rule, len(rules))
	for i, r := range rules {
		inner[i] = fieldregex.Rule{
			Source:  r.Source,
			Pattern: r.Pattern,
			Prefix:  r.Prefix,
		}
	}
	return fieldregex.New(fieldregex.WithRules(inner))
}
