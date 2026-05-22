// Package fieldstats collects frequency statistics for values of a given
// JSON log field across a stream of log entries.
package fieldstats

import (
	"fmt"
	"sort"
)

// Entry holds a value and the number of times it was observed.
type Entry struct {
	Value string
	Count int
}

// Collector accumulates value counts for a single field.
type Collector struct {
	field  string
	counts map[string]int
	total  int
}

// New returns a Collector that tracks occurrences of values for field.
func New(field string) *Collector {
	return &Collector{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add records the value of the tracked field found in entry.
// If the field is absent the entry is silently skipped.
func (c *Collector) Add(entry map[string]any) {
	v, ok := entry[c.field]
	if !ok {
		return
	}
	c.counts[fmt.Sprintf("%v", v)]++
	c.total++
}

// Total returns the number of entries that contained the tracked field.
func (c *Collector) Total() int { return c.total }

// Top returns up to n entries sorted by descending count.
// If n <= 0 all entries are returned.
func (c *Collector) Top(n int) []Entry {
	entries := make([]Entry, 0, len(c.counts))
	for v, cnt := range c.counts {
		entries = append(entries, Entry{Value: v, Count: cnt})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	if n > 0 && n < len(entries) {
		return entries[:n]
	}
	return entries
}

// Percent returns the percentage of total entries that had the given value,
// rounded to two decimal places. Returns 0 if the value was never observed
// or if no entries have been collected.
func (c *Collector) Percent(value string) float64 {
	if c.total == 0 {
		return 0
	}
	return float64(c.counts[value]) / float64(c.total) * 100
}

// Reset clears all accumulated counts.
func (c *Collector) Reset() {
	c.counts = make(map[string]int)
	c.total = 0
}
