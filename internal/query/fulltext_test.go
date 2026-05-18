package query

import (
	"testing"
)

func TestFullTextMatcherNoTerms(t *testing.T) {
	m := NewFullTextMatcher([]string{})
	record := map[string]interface{}{"msg": "hello world"}
	if !m.Match(record) {
		t.Error("expected match when no terms provided")
	}
}

func TestFullTextMatcherEmptyTermsFiltered(t *testing.T) {
	m := NewFullTextMatcher([]string{" ", "", "  "})
	if m.HasTerms() {
		t.Error("expected no terms after filtering whitespace")
	}
}

func TestFullTextMatcherSingleTerm(t *testing.T) {
	m := NewFullTextMatcher([]string{"error"})
	tests := []struct {
		record map[string]interface{}
		want   bool
	}{
		{map[string]interface{}{"msg": "an ERROR occurred"}, true},
		{map[string]interface{}{"msg": "all good"}, false},
		{map[string]interface{}{"level": "error", "msg": "boom"}, true},
	}
	for _, tt := range tests {
		got := m.Match(tt.record)
		if got != tt.want {
			t.Errorf("Match(%v) = %v, want %v", tt.record, got, tt.want)
		}
	}
}

func TestFullTextMatcherMultipleTerms(t *testing.T) {
	m := NewFullTextMatcher([]string{"database", "timeout"})
	tests := []struct {
		record map[string]interface{}
		want   bool
	}{
		{map[string]interface{}{"msg": "database timeout exceeded"}, true},
		{map[string]interface{}{"msg": "database connection ok"}, false},
		{map[string]interface{}{"msg": "timeout", "src": "database"}, true},
		{map[string]interface{}{"msg": "nothing relevant"}, false},
	}
	for _, tt := range tests {
		got := m.Match(tt.record)
		if got != tt.want {
			t.Errorf("Match(%v) = %v, want %v", tt.record, got, tt.want)
		}
	}
}

func TestFullTextMatcherNumericValue(t *testing.T) {
	m := NewFullTextMatcher([]string{"42"})
	record := map[string]interface{}{"code": float64(42), "msg": "ok"}
	if !m.Match(record) {
		t.Error("expected numeric value 42 to match term \"42\"")
	}
}
