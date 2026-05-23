package pipeline

import (
	"github.com/qsocket/logslice/internal/fieldtimestamp"
)

// buildFieldTimestampTransformer constructs a fieldtimestamp.Transformer from
// the pipeline configuration and returns a step function suitable for use in
// the processing chain. It returns nil when no timestamp rules are configured.
func buildFieldTimestampTransformer(rules []fieldtimestamp.Rule) func(map[string]any) (map[string]any, error) {
	if len(rules) == 0 {
		return nil
	}
	tr := fieldtimestamp.New(fieldtimestamp.WithRules(rules))
	return func(entry map[string]any) (map[string]any, error) {
		return tr.Apply(entry)
	}
}
