package sampler_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/sampler"
)

func record(msg string) map[string]any {
	return map[string]any{"message": msg}
}

func TestHeadSampler(t *testing.T) {
	s := sampler.New(sampler.ModeHead, 3)
	kept := 0
	for i := 0; i < 10; i++ {
		if r, ok := s.Feed(record("x")); ok {
			if r == nil {
				t.Fatal("expected non-nil record")
			}
			kept++
		}
	}
	if kept != 3 {
		t.Fatalf("expected 3 records, got %d", kept)
	}
	if flushed := s.Flush(); len(flushed) != 0 {
		t.Fatalf("head sampler should flush nothing, got %d", len(flushed))
	}
}

func TestTailSampler(t *testing.T) {
	s := sampler.New(sampler.ModeTail, 3)
	for i := 0; i < 10; i++ {
		if _, ok := s.Feed(record("x")); ok {
			t.Fatal("tail sampler should never emit immediately")
		}
	}
	flushed := s.Flush()
	if len(flushed) != 3 {
		t.Fatalf("expected 3 tail records, got %d", len(flushed))
	}
	if len(s.Flush()) != 0 {
		t.Fatal("second flush should be empty")
	}
}

func TestRateSampler_Deterministic(t *testing.T) {
	s := sampler.New(sampler.ModeRate, 2, sampler.WithSeed(1))
	kept := 0
	for i := 0; i < 100; i++ {
		if _, ok := s.Feed(record("x")); ok {
			kept++
		}
	}
	// With rate 1-in-2 and seed 1 over 100 records we expect roughly 50 kept.
	if kept < 30 || kept > 70 {
		t.Fatalf("rate sampler kept unexpected count %d/100", kept)
	}
}

func TestSamplerNLessThanOne(t *testing.T) {
	// n=0 should be clamped to 1, keeping every record in head mode.
	s := sampler.New(sampler.ModeHead, 0)
	if _, ok := s.Feed(record("x")); !ok {
		t.Fatal("expected record to be kept when n is clamped to 1")
	}
}
