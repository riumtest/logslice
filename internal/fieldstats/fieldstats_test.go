package fieldstats_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldstats"
)

func makeEntry(field, value string) map[string]any {
	return map[string]any{field: value}
}

func TestCollectorEmpty(t *testing.T) {
	c := fieldstats.New("level")
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
	if top := c.Top(0); len(top) != 0 {
		t.Fatalf("expected empty top, got %v", top)
	}
}

func TestCollectorCounts(t *testing.T) {
	c := fieldstats.New("level")
	for _, v := range []string{"info", "info", "warn", "error", "info"} {
		c.Add(makeEntry("level", v))
	}
	if got := c.Total(); got != 5 {
		t.Fatalf("expected total 5, got %d", got)
	}
	top := c.Top(0)
	if len(top) != 3 {
		t.Fatalf("expected 3 distinct values, got %d", len(top))
	}
	if top[0].Value != "info" || top[0].Count != 3 {
		t.Errorf("expected info=3, got %+v", top[0])
	}
}

func TestCollectorTopN(t *testing.T) {
	c := fieldstats.New("level")
	for _, v := range []string{"a", "b", "b", "c", "c", "c"} {
		c.Add(makeEntry("level", v))
	}
	top := c.Top(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Value != "c" {
		t.Errorf("expected c first, got %s", top[0].Value)
	}
}

func TestCollectorSkipsMissingField(t *testing.T) {
	c := fieldstats.New("level")
	c.Add(map[string]any{"msg": "hello"})
	if c.Total() != 0 {
		t.Fatal("expected total 0 when field absent")
	}
}

func TestCollectorReset(t *testing.T) {
	c := fieldstats.New("level")
	c.Add(makeEntry("level", "info"))
	c.Reset()
	if c.Total() != 0 {
		t.Fatal("expected total 0 after reset")
	}
	if len(c.Top(0)) != 0 {
		t.Fatal("expected empty top after reset")
	}
}

func TestCollectorNumericValue(t *testing.T) {
	c := fieldstats.New("code")
	c.Add(map[string]any{"code": 200})
	c.Add(map[string]any{"code": 200})
	c.Add(map[string]any{"code": 404})
	if c.Total() != 3 {
		t.Fatalf("expected 3, got %d", c.Total())
	}
	top := c.Top(1)
	if top[0].Value != "200" || top[0].Count != 2 {
		t.Errorf("unexpected top entry: %+v", top[0])
	}
}
