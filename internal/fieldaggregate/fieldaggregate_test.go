package fieldaggregate_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldaggregate"
)

func entry(field string, value any) map[string]any {
	return map[string]any{field: value}
}

func TestSumAggregator(t *testing.T) {
	a := fieldaggregate.New("latency", fieldaggregate.Sum)
	a.Add(entry("latency", float64(10)))
	a.Add(entry("latency", float64(20)))
	a.Add(entry("latency", float64(30)))
	res, ok := a.Result()
	if !ok {
		t.Fatal("expected result")
	}
	if res != 60 {
		t.Fatalf("sum: want 60, got %v", res)
	}
}

func TestMinAggregator(t *testing.T) {
	a := fieldaggregate.New("latency", fieldaggregate.Min)
	a.Add(entry("latency", float64(5)))
	a.Add(entry("latency", float64(15)))
	a.Add(entry("latency", float64(3)))
	res, ok := a.Result()
	if !ok {
		t.Fatal("expected result")
	}
	if res != 3 {
		t.Fatalf("min: want 3, got %v", res)
	}
}

func TestMaxAggregator(t *testing.T) {
	a := fieldaggregate.New("latency", fieldaggregate.Max)
	a.Add(entry("latency", float64(5)))
	a.Add(entry("latency", float64(99)))
	a.Add(entry("latency", float64(42)))
	res, ok := a.Result()
	if !ok {
		t.Fatal("expected result")
	}
	if res != 99 {
		t.Fatalf("max: want 99, got %v", res)
	}
}

func TestAvgAggregator(t *testing.T) {
	a := fieldaggregate.New("latency", fieldaggregate.Avg)
	a.Add(entry("latency", float64(10)))
	a.Add(entry("latency", float64(20)))
	a.Add(entry("latency", float64(30)))
	res, ok := a.Result()
	if !ok {
		t.Fatal("expected result")
	}
	if res != 20 {
		t.Fatalf("avg: want 20, got %v", res)
	}
}

func TestAggregatorSkipsMissingField(t *testing.T) {
	a := fieldaggregate.New("latency", fieldaggregate.Sum)
	a.Add(map[string]any{"other": float64(99)})
	_, ok := a.Result()
	if ok {
		t.Fatal("expected no result when field is absent")
	}
	if a.Count() != 0 {
		t.Fatalf("count: want 0, got %d", a.Count())
	}
}

func TestAggregatorSkipsNonNumeric(t *testing.T) {
	a := fieldaggregate.New("latency", fieldaggregate.Sum)
	a.Add(entry("latency", "fast"))
	a.Add(entry("latency", float64(5)))
	res, ok := a.Result()
	if !ok {
		t.Fatal("expected result")
	}
	if res != 5 {
		t.Fatalf("sum: want 5, got %v", res)
	}
	if a.Count() != 1 {
		t.Fatalf("count: want 1, got %d", a.Count())
	}
}
