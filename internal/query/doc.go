// Package query implements the logslice query DSL for filtering structured
// JSON log entries.
//
// # Syntax
//
// A query is a whitespace-separated list of filter expressions. All filters
// are ANDed together. Each filter has the form:
//
//	field<op>value
//
// Supported operators:
//
//	=    equal
//	!=   not equal
//	>    greater than (numeric)
//	<    less than (numeric)
//	>=   greater than or equal (numeric)
//	<=   less than or equal (numeric)
//	~    contains (substring match)
//
// # Examples
//
//	level=error
//	level=error status>=500
//	msg~timeout latency>200
//
package query
