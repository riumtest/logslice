package pipeline

import (
	"strings"

	"github.com/your-org/logslice/internal/fieldscope"
)

// buildFieldScopeTransformer parses the --field-scope flag values and returns
// a configured *fieldscope.Transformer, or nil when no rules are provided.
//
// Each flag value is expected in the form "source.path:dest", e.g.
//
//	"http.method:method"
func buildFieldScopeTransformer(specs []string) *fieldscope.Transformer {
	var rules []fieldscope.Rule
	for _, spec := range specs {
		idx := strings.LastIndex(spec, ":")
		if idx <= 0 || idx == len(spec)-1 {
			continue // malformed spec — skip
		}
		rules = append(rules, fieldscope.Rule{
			Source: spec[:idx],
			Dest:   spec[idx+1:],
		})
	}
	if len(rules) == 0 {
		return nil
	}
	return fieldscope.New(rules)
}
