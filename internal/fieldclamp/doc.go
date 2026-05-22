// Package fieldclamp provides a transform that clamps numeric JSON log fields
// to configurable minimum and/or maximum bounds.
//
// Each Rule targets a single field by name and accepts optional Min and Max
// float64 pointers. When a bound is nil it is not enforced. Non-numeric field
// values and fields absent from the entry are left untouched.
//
// Example usage:
//
//	min := 0.0
//	max := 100.0
//	c := fieldclamp.New([]fieldclamp.Rule{
//		{Field: "score", Min: &min, Max: &max},
//	})
//	out := c.Apply(entry)
package fieldclamp
