// Package output provides a Writer type for directing formatted log lines
// to configurable destinations (stdout, files, buffers, etc.).
//
// # Usage
//
// Create a Writer with options:
//
//	w := output.New(
//		output.WithDestination(os.Stdout),
//		output.WithColorize(true),
//	)
//
// Write formatted lines:
//
//	if err := w.Write(line); err != nil {
//		log.Fatal(err)
//	}
//
// # Colorization
//
// When WithColorize(true) is set, lines containing level keywords
// (error, warn, info) are wrapped with ANSI escape codes:
//   - error → red
//   - warn  → yellow
//   - info  → cyan
//
// Colorization is applied after formatting and before writing to
// the underlying destination.
package output
