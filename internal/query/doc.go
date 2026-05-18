// Package query provides a lightweight DSL for filtering structured JSON log
// records. A filter expression takes the form:
//
//	<field> <op> <value>
//
// Supported operators:
//
//	=   equal (string or numeric)
//	!=  not equal
//	>   greater than (numeric)
//	>=  greater than or equal (numeric)
//	<   less than (numeric)
//	<=  less than or equal (numeric)
//	~   contains substring (string)
//
// Examples:
//
//	level = error
//	status >= 500
//	msg ~ timeout
//
// Fields are looked up in the top-level keys of the JSON object. Values are
// always treated as strings unless the operator implies numeric comparison, in
// which case both sides are parsed as float64.
package query
