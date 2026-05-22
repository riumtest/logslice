// Package fieldstats provides a Collector that accumulates frequency
// statistics for the values of a named field across a stream of structured
// JSON log entries.
//
// Typical usage:
//
//	c := fieldstats.New("level")
//	for _, entry := range entries {
//		c.Add(entry)
//	}
//	for _, e := range c.Top(5) {
//		fmt.Printf("%s: %d\n", e.Value, e.Count)
//	}
//
// The Collector is not safe for concurrent use without external
// synchronisation.
package fieldstats
