// Package dedupe implements consecutive-duplicate suppression for structured
// log streams.
//
// A [Deduper] computes a SHA-256 fingerprint of each log entry and compares it
// to the fingerprint of the immediately preceding entry. When the two match the
// entry is considered a duplicate and [Deduper.IsDuplicate] returns true.
//
// By default the fingerprint covers the entire JSON object. Use [WithFields] to
// restrict fingerprinting to a specific subset of fields — useful when you want
// to deduplicate on message text alone while ignoring timestamps or counters.
//
// Only *consecutive* duplicates are suppressed. If the same message appears
// again after a different message it will pass through normally.
//
// Example:
//
//	d := dedupe.New(dedupe.WithFields("msg", "level"))
//	for _, entry := range entries {
//		if d.IsDuplicate(entry) {
//			continue
//		}
//		process(entry)
//	}
package dedupe
