// Package pipeline wires together a reader, query filter, and formatter
// into a single processing loop.
package pipeline

import (
	"errors"
	"io"

	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/query"
	"github.com/yourorg/logslice/internal/reader"
)

// Config holds the settings for a pipeline run.
type Config struct {
	Sources    []io.Reader
	FilterExpr string
	Out        io.Writer
	Format     formatter.Format
	TimeKey    string
	LevelKey   string
	MessageKey string
}

// Run reads log entries from all sources, applies the optional filter, and
// writes matching entries to the configured output. It returns the number of
// entries written and any terminal error.
func Run(cfg Config) (int, error) {
	if len(cfg.Sources) == 0 {
		return 0, errors.New("pipeline: no sources provided")
	}

	var filter *query.Filter
	if cfg.FilterExpr != "" {
		var err error
		filter, err = query.Parse(cfg.FilterExpr)
		if err != nil {
			return 0, err
		}
	}

	fopts := []formatter.Option{formatter.WithFormat(cfg.Format)}
	if cfg.TimeKey != "" || cfg.LevelKey != "" || cfg.MessageKey != "" {
		tk := orDefault(cfg.TimeKey, "time")
		lk := orDefault(cfg.LevelKey, "level")
		mk := orDefault(cfg.MessageKey, "msg")
		fopts = append(fopts, formatter.WithKeys(tk, lk, mk))
	}
	fmt := formatter.New(cfg.Out, fopts...)

	r := reader.New(cfg.Sources...)
	written := 0

	for {
		entry, raw, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			// skip unparseable lines
			continue
		}

		if filter != nil && !filter.Match(entry) {
			continue
		}

		if werr := fmt.Write(entry, raw); werr != nil {
			return written, werr
		}
		written++
	}

	return written, nil
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
