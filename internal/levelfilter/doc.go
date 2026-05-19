// Package levelfilter implements severity-level based filtering for structured
// JSON log streams.
//
// It recognises the canonical level names used by popular Go logging libraries
// (zerolog, zap, logrus, slog):
//
//	trace < debug < info < warn < error < fatal
//
// Usage:
//
//	f := levelfilter.New(levelfilter.LevelWarn, "level")
//	if f.Match(entry) {
//	    // entry level is warn, error, or fatal
//	}
//
// Entries whose level field is absent, non-string, or unrecognised are always
// passed through so that unusual log formats are not silently dropped.
package levelfilter
