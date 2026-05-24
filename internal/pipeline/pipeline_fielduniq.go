package pipeline

import "github.com/yourorg/logslice/internal/fielduniq"

// buildFieldUniqTransformer constructs a fielduniq.Transformer from the
// pipeline configuration and registers it as a transform stage. If no
// fields are configured the stage is skipped entirely.
func buildFieldUniqTransformer(cfg *Config) stage {
	if len(cfg.UniqFields) == 0 {
		return nil
	}
	tr := fielduniq.New(cfg.UniqFields)
	return func(entry map[string]any) (map[string]any, bool) {
		return tr.Apply(entry), true
	}
}
