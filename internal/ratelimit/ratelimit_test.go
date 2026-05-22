package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourusername/logslice/internal/ratelimit"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNoRateLimit(t *testing.T) {
	l := ratelimit.New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatal("expected Allow() == true when rate <= 0")
		}
	}
}

func TestRateLimitBucketExhausts(t *testing.T) {
	base := time.Now()
	l := ratelimit.New(3, ratelimit.WithNow(fixedClock(base)))

	if !l.Allow() { t.Fatal("token 1 should be allowed") }
	if !l.Allow() { t.Fatal("token 2 should be allowed") }
	if !l.Allow() { t.Fatal("token 3 should be allowed") }
	if l.Allow()  { t.Fatal("token 4 should be denied (bucket empty)") }
}

func TestRateLimitRefillsAfterTime(t *testing.T) {
	base := time.Now()
	current := base
	clock := func() time.Time { return current }

	l := ratelimit.New(2, ratelimit.WithNow(clock))

	// Exhaust the bucket.
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("bucket should be empty")
	}

	// Advance clock by 1 second — should refill 2 tokens.
	current = base.Add(time.Second)
	if !l.Allow() { t.Fatal("should allow after refill") }
	if !l.Allow() { t.Fatal("should allow second token after refill") }
	if l.Allow()  { t.Fatal("bucket should be empty again") }
}

func TestRateLimitCapNotExceeded(t *testing.T) {
	base := time.Now()
	current := base
	clock := func() time.Time { return current }

	l := ratelimit.New(2, ratelimit.WithNow(clock))

	// Exhaust bucket first.
	l.Allow()
	l.Allow()

	// Advance by 10 seconds — refill should be capped at 2.
	current = base.Add(10 * time.Second)
	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 2 {
		t.Fatalf("expected 2 allowed after cap refill, got %d", allowed)
	}
}
