// Package formatter renders structured log entries for human consumption
// or passes them through as raw JSON lines.
//
// # Usage
//
//	f := formatter.New(os.Stdout, formatter.WithFormat(formatter.FormatText))
//	f.Write(entry, rawLine)
//
// # Text format
//
// Text output is structured as:
//
//	<time> [<LEVEL>] <message> key=value key=value ...
//
// Fields that do not map to the configured time, level, or message keys are
// appended as key=value pairs in alphabetical order.
//
// # JSON format
//
// When FormatJSON is selected, the original raw JSON line is written verbatim.
// This is useful for piping output to other tools.
//
// # Custom keys
//
// Many log libraries use different field names (e.g. "timestamp" instead of
// "time"). Use WithKeys to configure the field names that correspond to the
// timestamp, severity level, and human-readable message.
package formatter
