package pipeline

import (
	"github.com/user/logslice/internal/fieldencoding"
)

// EncodingRule mirrors fieldencoding.Rule for pipeline configuration.
type EncodingRule struct {
	Field string
	Mode  string
	Dest  string
}

func buildFieldEncodingTransformer(rules []EncodingRule) *fieldencoding.Transformer {
	if len(rules) == 0 {
		return nil
	}
	inner := make([]fieldencoding.Rule, len(rules))
	for i, r := range rules {
		inner[i] = fieldencoding.Rule{
			Field: r.Field,
			Mode:  r.Mode,
			Dest:  r.Dest,
		}
	}
	return fieldencoding.WithRules(inner)
}
