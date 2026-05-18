// Package sampler provides three sampling strategies for reducing the volume
// of log records passed through the logslice pipeline.
//
// # Modes
//
// ModeRate keeps approximately 1-in-N records chosen at random, useful for
// reducing high-volume streams while preserving a representative sample.
//
// ModeHead keeps only the first N records seen, then discards the rest.
// This is useful when you want a quick snapshot of the beginning of a stream.
//
// ModeTail buffers the last N records and emits them all at Flush time.
// Records are not emitted from Feed; callers must call Flush after the source
// is exhausted to retrieve the tail window.
//
// # Usage
//
//	s := sampler.New(sampler.ModeHead, 100)
//	for _, rec := range records {
//		if out, ok := s.Feed(rec); ok {
//			// emit out
//		}
//	}
//	for _, out := range s.Flush() {
//		// emit tail records
//	}
package sampler
