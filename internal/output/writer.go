// Package output provides writers for directing formatted log output
// to various destinations such as stdout, files, or buffers.
package output

import (
	"fmt"
	"io"
	"os"
)

// Writer wraps an io.Writer with optional colorization and line counting.
type Writer struct {
	dst      io.Writer
	colorize bool
	count    int
}

// Option configures a Writer.
type Option func(*Writer)

// WithColorize enables ANSI color codes in output.
func WithColorize(c bool) Option {
	return func(w *Writer) {
		w.colorize = c
	}
}

// WithDestination sets the underlying io.Writer.
func WithDestination(dst io.Writer) Option {
	return func(w *Writer) {
		w.dst = dst
	}
}

// New creates a Writer with the given options.
// Defaults to os.Stdout with no colorization.
func New(opts ...Option) *Writer {
	w := &Writer{dst: os.Stdout}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Write writes a formatted line to the destination.
// It appends a newline if the line does not already end with one.
func (w *Writer) Write(line string) error {
	if w.colorize {
		line = colorize(line)
	}
	if len(line) == 0 || line[len(line)-1] != '\n' {
		line += "\n"
	}
	_, err := fmt.Fprint(w.dst, line)
	if err != nil {
		return err
	}
	w.count++
	return nil
}

// Count returns the number of lines successfully written.
func (w *Writer) Count() int {
	return w.count
}

// colorize applies simple ANSI coloring based on log level keywords.
func colorize(line string) string {
	const (
		red    = "\033[31m"
		yellow = "\033[33m"
		cyan   = "\033[36m"
		reset  = "\033[0m"
	)
	switch {
	case contains(line, "error") || contains(line, "ERROR"):
		return red + line + reset
	case contains(line, "warn") || contains(line, "WARN"):
		return yellow + line + reset
	case contains(line, "info") || contains(line, "INFO"):
		return cyan + line + reset
	default:
		return line
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
