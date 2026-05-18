package query

import (
	"fmt"
	"strings"
)

// Query holds a list of filters that are ANDed together.
type Query struct {
	Filters []*Filter
}

// Match reports whether the entry satisfies all filters in the query.
func (q *Query) Match(entry map[string]any) bool {
	for _, f := range q.Filters {
		if !f.Match(entry) {
			return false
		}
	}
	return true
}

// Parse parses a query string into a Query.
// Syntax: field=value field!=value field>value field~substring ...
// Multiple expressions are separated by whitespace and ANDed.
func Parse(input string) (*Query, error) {
	if strings.TrimSpace(input) == "" {
		return &Query{}, nil
	}

	tokens := strings.Fields(input)
	filters := make([]*Filter, 0, len(tokens))

	for _, tok := range tokens {
		f, err := parseFilter(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid filter %q: %w", tok, err)
		}
		filters = append(filters, f)
	}

	return &Query{Filters: filters}, nil
}

var operators = []Op{OpGte, OpLte, OpNeq, OpGt, OpLt, OpEq, OpContains}

func parseFilter(token string) (*Filter, error) {
	for _, op := range operators {
		idx := strings.Index(token, string(op))
		if idx <= 0 {
			continue
		}
		field := token[:idx]
		value := token[idx+len(op):]
		if field == "" || value == "" {
			return nil, fmt.Errorf("field and value must not be empty")
		}
		return &Filter{Field: field, Op: op, Value: value}, nil
	}
	return nil, fmt.Errorf("no valid operator found")
}
