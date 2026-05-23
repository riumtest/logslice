package pipeline

import (
	"encoding/json"

	"github.com/naturalselectionlabs/logslice/internal/fieldwindow"
)

// buildFieldWindowTransformer constructs a fieldwindow.Transformer from the
// pipeline configuration and returns a step function compatible with the
// pipeline runner. It returns nil when no window rules are configured.
func buildFieldWindowTransformer(rules []fieldwindow.Rule) func(map[string]json.RawMessage) (map[string]json.RawMessage, bool) {
	if len(rules) == 0 {
		return nil
	}
	tr := fieldwindow.New(rules)
	return func(entry map[string]json.RawMessage) (map[string]json.RawMessage, bool) {
		return tr.Transform(entry), true
	}
}
