// Package ratelimit provides a token-bucket rate limiter for log output,
// allowing users to cap the number of log lines emitted per second.
package ratelimit

import (
	"time"
)

// Limiter controls the rate of log line emission.
type Limiter struct {
	rate     int
	bucket   int
	cap      int
	lastTick time.Time
	now      func() time.Time
}

// Option configures a Limiter.
type Option func(*Limiter)

// WithNow overrides the clock used by the limiter (useful for testing).
func WithNow(fn func() time.Time) Option {
	return func(l *Limiter) {
		l.now = fn
	}
}

// New creates a Limiter that allows at most ratePerSec lines per second.
// A ratePerSec <= 0 disables rate limiting (Allow always returns true).
func New(ratePerSec int, opts ...Option) *Limiter {
	l := &Limiter{
		rate:   ratePerSec,
		cap:    ratePerSec,
		now:    time.Now,
	}
	for _, o := range opts {
		o(l)
	}
	l.bucket = l.cap
	l.lastTick = l.now()
	return l
}

// Allow reports whether a log line should be emitted.
// It refills the token bucket based on elapsed time and consumes one token.
func (l *Limiter) Allow() bool {
	if l.rate <= 0 {
		return true
	}
	now := l.now()
	elapsed := now.Sub(l.lastTick)
	l.lastTick = now

	// Refill tokens proportional to elapsed time.
	refill := int(elapsed.Seconds() * float64(l.rate))
	if refill > 0 {
		l.bucket += refill
		if l.bucket > l.cap {
			l.bucket = l.cap
		}
	}

	if l.bucket <= 0 {
		return false
	}
	l.bucket--
	return true
}
