// Package fieldmath provides a Transformer that applies arithmetic
// operations (add, sub, mul, div) between numeric fields in a log entry,
// writing the result to a destination field.
//
// Rules are evaluated in order. If either operand field is absent or
// non-numeric the rule is silently skipped. Division by zero is also
// skipped rather than producing NaN or an error in the output stream.
//
// Example usage:
//
//	tf := fieldmath.New([]fieldmath.Rule{
//		{Left: "bytes_sent", Right: "bytes_recv", Dest: "bytes_total", Op: fieldmath.OpAdd},
//		{Left: "bytes_sent", Right: "requests",   Dest: "avg_size",    Op: fieldmath.OpDiv},
//	})
//	out := tf.Apply(entry)
package fieldmath
