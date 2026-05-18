// Command logslice streams and filters structured JSON logs.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
)

func main() {
	var (
		filter  = flag.String("filter", "", "Query filter expression, e.g. 'level=error'")
		format  = flag.String("format", "text", "Output format: text or json")
		keys    = flag.String("keys", "", "Comma-separated keys to include in text output")
		color   = flag.Bool("color", false, "Colorize output based on log level")
	)
	flag.Usage = usage
	flag.Parse()

	sources := flag.Args()

	var fmtOpts []formatter.Option
	fmtOpts = append(fmtOpts, formatter.WithFormat(*format))
	if *keys != "" {
		fmtOpts = append(fmtOpts, formatter.WithKeys(*keys))
	}

	out := output.New(
		output.WithDestination(os.Stdout),
		output.WithColorize(*color),
	)

	writeFn := func(line string) error {
		return out.Write(line)
	}

	if err := pipeline.Run(sources, *filter, fmtOpts, writeFn); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "logslice: wrote %d lines\n", out.Count())
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: logslice [options] [file ...]

Stream and filter structured JSON logs.

Options:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Examples:
  logslice -filter 'level=error' app.log
  logslice -format json -filter 'status>=500' access.log
  cat app.log | logslice -color -filter 'level=warn'
`)
}
