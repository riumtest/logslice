// Package ratelimit implements a token-bucket rate limiter for logslice.
//
// It is used to cap the number of log lines written to output per second,
// preventing runaway log sources from flooding the terminal.
//
// Usage:
//
//	limiter := ratelimit.New(100) // allow up to 100 lines/sec
//	for _, line := range lines {
//		if limiter.Allow() {
//			fmt.Println(line)
//		}
//	}
//
// A rate of 0 or less disables limiting entirely — every call to Allow
// returns true.
//
// The limiter is not safe for concurrent use without external locking.
package ratelimit
