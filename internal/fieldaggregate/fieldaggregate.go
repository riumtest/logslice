// Package fieldaggregate computes aggregate values (sum, min, max, avg)
// over a numeric field across all processed log entries.
package fieldaggregate

import "math"

// Op is the aggregation operation to apply.
type Op string

const (
	Sum Op = "sum"
	Min Op = "min"
	Max Op = "max"
	Avg Op = "avg"
)

// Aggregator accumulates values from a named field and exposes the result.
type Aggregator struct {
	field string
	op    Op
	sum   float64
	min   float64
	max   float64
	count int64
}

// Option configures an Aggregator.
type Option func(*Aggregator)

// New returns an Aggregator that reads field and applies op.
func New(field string, op Op, opts ...Option) *Aggregator {
	a := &Aggregator{
		field: field,
		op:    op,
		min:   math.MaxFloat64,
		max:   -math.MaxFloat64,
	}
	for _, o := range opts {
		o(a)
	}
	return a
}

// Add ingests a single log entry map.
func (a *Aggregator) Add(entry map[string]any) {
	v, ok := entry[a.field]
	if !ok {
		return
	}
	f, ok := toFloat(v)
	if !ok {
		return
	}
	a.sum += f
	a.count++
	if f < a.min {
		a.min = f
	}
	if f > a.max {
		a.max = f
	}
}

// Result returns the aggregate value. Returns 0 and false if no data was seen.
func (a *Aggregator) Result() (float64, bool) {
	if a.count == 0 {
		return 0, false
	}
	switch a.op {
	case Sum:
		return a.sum, true
	case Min:
		return a.min, true
	case Max:
		return a.max, true
	case Avg:
		return a.sum / float64(a.count), true
	default:
		return 0, false
	}
}

// Count returns the number of entries that contributed a value.
func (a *Aggregator) Count() int64 { return a.count }

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	}
	return 0, false
}
