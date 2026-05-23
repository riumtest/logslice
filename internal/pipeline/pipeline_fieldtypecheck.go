package pipeline

import (
	"github.com/user/logslice/internal/fieldtypecheck"
)

// TypeCheckRule is the pipeline-level representation of a field type rule.
type TypeCheckRule struct {
	Field    string
	Expected string
}

// buildFieldTypeChecker constructs a fieldtypecheck.Checker from pipeline
// configuration, or returns nil when no rules are defined.
func buildFieldTypeChecker(rules []TypeCheckRule, destField string, rejectMode bool) *fieldtypecheck.Checker {
	if len(rules) == 0 {
		return nil
	}
	var r []fieldtypecheck.Rule
	for _, tc := range rules {
		r = append(r, fieldtypecheck.Rule{Field: tc.Field, Expected: tc.Expected})
	}
	var opts []func(*fieldtypecheck.Checker)
	if destField != "" {
		opts = append(opts, fieldtypecheck.WithDestField(destField))
	}
	if rejectMode {
		opts = append(opts, fieldtypecheck.WithRejectMode())
	}
	return fieldtypecheck.New(r, opts...)
}
