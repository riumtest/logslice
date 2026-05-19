// Package pipeline wires together the reader, filter, sampler, time-range
// filter, formatter, and output writer into a single streaming pipeline.
package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/query"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/sampler"
	"github.com/yourorg/logslice/internal/source"
	"github.com/yourorg/logslice/internal/timerange"
)

// Options configures a pipeline run.
type Options struct {
	Sources       []string
	Stdin         io.Reader
	Output        io.Writer
	Format        string
	Keys          []string
	FilterExpr    string
	Colorize      bool
	SampleMode    string // "head", "tail", "rate"
	SampleN       int
	SampleRate    float64
	SampleSeed    int64
	TimeField     string
	TimeRangeExpr string
	TimeRangeNow  time.Time
}

// Run executes the pipeline, writing formatted log lines to opts.Output.
func Run(opts Options) error {
	srcs, err := source.New(opts.Sources, source.WithStdin(opts.Stdin))
	if err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}
	defer func() {
		for _, s := range srcs {
			s.Close()
		}
	}()

	var filt *query.Filter
	if opts.FilterExpr != "" {
		f, err := query.Parse(opts.FilterExpr)
		if err != nil {
			return fmt.Errorf("pipeline: bad filter: %w", err)
		}
		filt = f
	}

	var trFilter timerange.Filter
	if opts.TimeRangeExpr != "" {
		now := opts.TimeRangeNow
		if now.IsZero() {
			now = time.Now()
		}
		trFilter, err = timerange.Parse(opts.TimeRangeExpr, now)
		if err != nil {
			return fmt.Errorf("pipeline: bad time range: %w", err)
		}
	}

	var records []map[string]json.RawMessage
	for _, src := range srcs {
		r := reader.New(src)
		for {
			rec, err := r.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				continue
			}
			if filt != nil && !filt.Match(rec) {
				continue
			}
			tf := orDefault(opts.TimeField, "time")
			if !trFilter.IsZero() && !trFilter.Match(rec, tf) {
				continue
			}
			records = append(records, rec)
		}
	}

	if opts.SampleMode != "" {
		var samplerOpts []sampler.Option
		if opts.SampleSeed != 0 {
			samplerOpts = append(samplerOpts, sampler.WithSeed(opts.SampleSeed))
		}
		s := sampler.New(opts.SampleMode, opts.SampleN, opts.SampleRate, samplerOpts...)
		records = s.Sample(records)
	}

	fmt_ := orDefault(opts.Format, "text")
	keys := opts.Keys
	if len(keys) == 0 {
		keys = []string{"time", "level", "msg"}
	}
	fmtr := formatter.New(
		formatter.WithFormat(fmt_),
		formatter.WithKeys(keys),
	)

	w := output.New(
		output.WithDestination(opts.Output),
		output.WithColorize(opts.Colorize),
	)

	for _, rec := range records {
		line, err := fmtr.Format(rec)
		if err != nil {
			continue
		}
		if err := w.Write(strings.TrimRight(line, "\n")); err != nil {
			return err
		}
	}
	return nil
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
