// Package fieldcoalesce implements a log entry transformer that performs
// coalesce operations across multiple source fields.
//
// A coalesce rule specifies an ordered list of source fields and a destination
// field. The transformer evaluates each source in order and writes the first
// non-null, non-empty value to the destination field. If no source yields a
// usable value the destination field is left unset.
//
// Example usage:
//
//	rules := []fieldcoalesce.Rule{
//		{Sources: []string{"msg", "message", "text"}, Dest: "canonical_msg"},
//	}
//	tr, err := fieldcoalesce.New(rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := tr.Transform(entry)
package fieldcoalesce
