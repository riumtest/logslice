package fieldsample_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldsample"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

// alwaysZero makes the sampler always "keep" (rand returns 0).
func alwaysZero(_ int) int { return 0 }

// alwaysOne makes the sampler always "drop" (rand returns 1).
func alwaysOne(_ int) int { return 1 }

func TestIdentityNoRules(t *testing.T) {
	tr := fieldsample.New(nil)
	e := entry("msg", "hello")
	out := tr.Apply(e)
	if out == nil {
		t.Fatal("expected entry to pass through")
	}
	if out["msg"] != "hello" {
		t.Fatalf("unexpected value: %v", out["msg"])
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	tr := fieldsample.New(nil)
	e := entry("a", "1")
	out := tr.Apply(e)
	out["a"] = "mutated"
	if e["a"] != "1" {
		t.Fatal("original entry was mutated")
	}
}

func TestKeepAllWhenNIsOne(t *testing.T) {
	rules := []fieldsample.Rule{{Field: "level", Value: "error", N: 1}}
	tr := fieldsample.New(rules, fieldsample.WithRandFn(alwaysOne))
	e := entry("level", "error")
	if tr.Apply(e) == nil {
		t.Fatal("N=1 should keep every entry")
	}
}

func TestDropWhenRandNonZero(t *testing.T) {
	rules := []fieldsample.Rule{{Field: "level", Value: "debug", N: 10}}
	tr := fieldsample.New(rules, fieldsample.WithRandFn(alwaysOne))
	e := entry("level", "debug")
	if tr.Apply(e) != nil {
		t.Fatal("expected entry to be dropped")
	}
}

func TestKeepWhenRandZero(t *testing.T) {
	rules := []fieldsample.Rule{{Field: "level", Value: "debug", N: 10}}
	tr := fieldsample.New(rules, fieldsample.WithRandFn(alwaysZero))
	e := entry("level", "debug")
	if tr.Apply(e) == nil {
		t.Fatal("expected entry to be kept")
	}
}

func TestRuleSkippedWhenFieldMissing(t *testing.T) {
	rules := []fieldsample.Rule{{Field: "level", Value: "debug", N: 10}}
	tr := fieldsample.New(rules, fieldsample.WithRandFn(alwaysOne))
	e := entry("msg", "no level field here")
	// Rule does not apply, so entry should pass through.
	if tr.Apply(e) == nil {
		t.Fatal("expected entry without matching field to pass through")
	}
}

func TestRuleSkippedWhenValueMismatch(t *testing.T) {
	rules := []fieldsample.Rule{{Field: "level", Value: "debug", N: 10}}
	tr := fieldsample.New(rules, fieldsample.WithRandFn(alwaysOne))
	e := entry("level", "info")
	if tr.Apply(e) == nil {
		t.Fatal("expected non-matching value to pass through")
	}
}

func TestNormalisesNLessThanOne(t *testing.T) {
	rules := []fieldsample.Rule{{Field: "", N: 0}}
	tr := fieldsample.New(rules, fieldsample.WithRandFn(alwaysOne))
	e := entry("msg", "hi")
	// N normalised to 1 → keep all.
	if tr.Apply(e) == nil {
		t.Fatal("N<1 should be treated as N=1")
	}
}
