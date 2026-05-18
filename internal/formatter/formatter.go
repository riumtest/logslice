// Package formatter provides human-readable output formatting for structured JSON log entries.
package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Format controls the output style.
type Format int

const (
	// FormatText renders logs as human-readable text.
	FormatText Format = iota
	// FormatJSON renders logs as raw JSON lines.
	FormatJSON
)

// Formatter writes log entries to an output stream.
type Formatter struct {
	out    io.Writer
	format Format
	timeKey   string
	levelKey  string
	messageKey string
}

// Option is a functional option for Formatter.
type Option func(*Formatter)

// WithFormat sets the output format.
func WithFormat(f Format) Option {
	return func(fmt *Formatter) { fmt.format = f }
}

// WithKeys overrides the field names used for time, level, and message.
func WithKeys(timeKey, levelKey, messageKey string) Option {
	return func(f *Formatter) {
		f.timeKey = timeKey
		f.levelKey = levelKey
		f.messageKey = messageKey
	}
}

// New creates a Formatter writing to out.
func New(out io.Writer, opts ...Option) *Formatter {
	f := &Formatter{
		out:        out,
		format:     FormatText,
		timeKey:    "time",
		levelKey:   "level",
		messageKey: "msg",
	}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Write outputs a single log entry.
func (f *Formatter) Write(entry map[string]interface{}, raw string) error {
	if f.format == FormatJSON {
		_, err := fmt.Fprintln(f.out, raw)
		return err
	}
	return f.writeText(entry)
}

func (f *Formatter) writeText(entry map[string]interface{}) error {
	var parts []string

	if ts, ok := entry[f.timeKey]; ok {
		parts = append(parts, formatTime(ts))
	}
	if lvl, ok := entry[f.levelKey]; ok {
		parts = append(parts, fmt.Sprintf("[%s]", strings.ToUpper(fmt.Sprintf("%v", lvl))))
	}
	if msg, ok := entry[f.messageKey]; ok {
		parts = append(parts, fmt.Sprintf("%v", msg))
	}

	var extras []string
	for k, v := range entry {
		if k == f.timeKey || k == f.levelKey || k == f.messageKey {
			continue
		}
		extras = append(extras, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(extras)
	parts = append(parts, extras...)

	_, err := fmt.Fprintln(f.out, strings.Join(parts, " "))
	return err
}

func formatTime(v interface{}) string {
	switch t := v.(type) {
	case string:
		parsed, err := time.Parse(time.RFC3339, t)
		if err == nil {
			return parsed.Format("2006-01-02 15:04:05")
		}
		return t
	default:
		return fmt.Sprintf("%v", v)
	}
}
