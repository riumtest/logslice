package fieldflip_test

import (
	"testing"

	"github.com/logslice/logslice/internal/fieldflip"
)

func entry(kv ...any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestIdentityNoRules(t *testing.T) {
	f := fieldflip.New(nil)
	in := entry("a", "hello", "b", "world")
	out := f.Apply(in)
	if out["a"] != "hello" || out["b"] != "world" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestFlipTwoFields(t *testing.T) {
	f := fieldflip.New([]fieldflip.Rule{{A: "x", B: "y"}})
	in := entry("x", "foo", "y", "bar")
	out := f.Apply(in)
	if out["x"] != "bar" || out["y"] != "foo" {
		t.Fatalf("expected swap, got x=%v y=%v", out["x"], out["y"])
	}
}

func TestFlipMissingFieldSkipped(t *testing.T) {
	f := fieldflip.New([]fieldflip.Rule{{A: "a", B: "missing"}})
	in := entry("a", "value")
	out := f.Apply(in)
	if out["a"] != "value" {
		t.Fatalf("field a should be unchanged, got %v", out["a"])
	}
}

func TestFlipMultipleRules(t *testing.T) {
	rules := []fieldflip.Rule{
		{A: "first", B: "second"},
		{A: "third", B: "fourth"},
	}
	f := fieldflip.New(rules)
	in := entry("first", 1, "second", 2, "third", 3, "fourth", 4)
	out := f.Apply(in)
	if out["first"] != 2 || out["second"] != 1 {
		t.Fatalf("first/second not swapped: %v", out)
	}
	if out["third"] != 4 || out["fourth"] != 3 {
		t.Fatalf("third/fourth not swapped: %v", out)
	}
}

func TestOriginalEntryNotMutated(t *testing.T) {
	f := fieldflip.New([]fieldflip.Rule{{A: "p", B: "q"}})
	in := entry("p", "alpha", "q", "beta")
	_ = f.Apply(in)
	if in["p"] != "alpha" || in["q"] != "beta" {
		t.Fatal("original entry was mutated")
	}
}

func TestFlipNumericValues(t *testing.T) {
	f := fieldflip.New([]fieldflip.Rule{{A: "latency", B: "timeout"}})
	in := entry("latency", 42.0, "timeout", 100.0)
	out := f.Apply(in)
	if out["latency"] != 100.0 || out["timeout"] != 42.0 {
		t.Fatalf("numeric swap failed: %v", out)
	}
}
