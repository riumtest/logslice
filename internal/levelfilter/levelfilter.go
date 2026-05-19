// Package levelfilter provides log-level based filtering for structured JSON log entries.
package levelfilter

import "strings"

// Level represents a log severity level.
type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

// ParseLevel converts a string to a Level. Returns LevelUnknown if unrecognised.
func ParseLevel(s string) Level {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// Filter holds the minimum level threshold and the field name to inspect.
type Filter struct {
	min   Level
	field string
}

// New creates a Filter that passes entries whose level is >= min.
// field is the JSON key that holds the level string (e.g. "level").
func New(min Level, field string) *Filter {
	if field == "" {
		field = "level"
	}
	return &Filter{min: min, field: field}
}

// Match returns true when the entry's level is >= the configured minimum.
// Entries with an unrecognised or missing level value are always passed through.
func (f *Filter) Match(entry map[string]any) bool {
	raw, ok := entry[f.field]
	if !ok {
		return true
	}
	s, ok := raw.(string)
	if !ok {
		return true
	}
	lvl := ParseLevel(s)
	if lvl == LevelUnknown {
		return true
	}
	return lvl >= f.min
}
