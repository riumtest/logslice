package pipeline

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/fielddefault"
	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/levelfilter"
	"github.com/yourorg/logslice/internal/query"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/sampler"
	"github.com/yourorg/logslice/internal/timerange"
)

// Source pairs an io.Reader with a display name used in diagnostics.
type Source struct {
	Reader io.Reader
	Name   string
}

// FieldDefault describes a single field-default rule for the pipeline.
type FieldDefault struct {
	Field string
	Value any
}

// Config holds all pipeline configuration.
type Config struct {
	Sources       []Source
	Filter        string
	Format        string
	Keys          []string
	Output        io.Writer
	MinLevel      string
	TimeRange     string
	TimeField     string
	SampleMode    string
	SampleN       int
	SampleRate    float64
	SampleSeed    int64
	FieldDefaults []FieldDefault
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// Run executes the pipeline described by cfg.
func Run(cfg Config) error {
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("no sources provided")
	}

	// Build optional filter.
	var filter *query.Filter
	if cfg.Filter != "" {
		f, err := query.Parse(cfg.Filter)
		if err != nil {
			return fmt.Errorf("invalid filter: %w", err)
		}
		filter = f
	}

	// Build optional level filter.
	var lvlFilter *levelfilter.Filter
	if cfg.MinLevel != "" {
		lvl, err := levelfilter.ParseLevel(cfg.MinLevel)
		if err != nil {
			return fmt.Errorf("invalid level: %w", err)
		}
		lvlFilter = levelfilter.New(lvl, orDefault(cfg.TimeField, "level"))
	}

	// Build optional time-range filter.
	var tr *timerange.Range
	if cfg.TimeRange != "" {
		r, err := timerange.Parse(cfg.TimeRange)
		if err != nil {
			return fmt.Errorf("invalid time range: %w", err)
		}
		tr = r
	}

	// Build optional sampler.
	var samp sampler.Sampler
	if cfg.SampleMode != "" {
		s, err := sampler.New(
			cfg.SampleMode,
			cfg.SampleN,
			cfg.SampleRate,
			sampler.WithSeed(cfg.SampleSeed),
		)
		if err != nil {
			return fmt.Errorf("invalid sampler: %w", err)
		}
		samp = s
	}

	// Build optional field-default transformer.
	var defTransformer *fielddefault.Transformer
	if len(cfg.FieldDefaults) > 0 {
		rules := make([]fielddefault.Rule, len(cfg.FieldDefaults))
		for i, d := range cfg.FieldDefaults {
			rules[i] = fielddefault.Rule{Field: d.Field, Value: d.Value}
		}
		defTransformer = fielddefault.New(rules)
	}

	// Build formatter.
	fopts := []formatter.Option{
		formatter.WithFormat(orDefault(cfg.Format, "text")),
	}
	if len(cfg.Keys) > 0 {
		fopts = append(fopts, formatter.WithKeys(cfg.Keys))
	}
	fmt := formatter.New(fopts...)

	for _, src := range cfg.Sources {
		r := reader.New(src.Reader)
		for {
			entry, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue // skip malformed lines
			}

			// Apply field defaults before any filtering.
			if defTransformer != nil {
				entry = defTransformer.Transform(entry)
			}

			if filter != nil && !filter.Match(entry) {
				continue
			}
			if lvlFilter != nil && !lvlFilter.Allow(entry) {
				continue
			}
			if tr != nil && !tr.Match(entry, orDefault(cfg.TimeField, "time")) {
				continue
			}
			if samp != nil && !samp.Keep(entry) {
				continue
			}

			line, err := fmt.Format(entry)
			if err != nil {
				raw, _ := json.Marshal(entry)
				line = string(raw)
			}
			_, _ = io.WriteString(cfg.Output, line+"\n")
		}
	}
	return nil
}
