// Package ratelimit provides per-job-type rate limiting using token
// buckets, so worker concurrenty and downstream requeust rate can be
// controlled independently
package ratelimit

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// Limiter manages one token-bucket rate limiter per job type. job types
// with no configured limit are allowed through unrestricted
type Limiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
}

// New creates and empty Limiter. use SetLimit to configure rates per job
// types before starting the worker pool
func New() *Limiter {
	return &Limiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// SetLimit configures a rate limit for jobType: ratePerSecond sustained
// rate, with burst allowing short spikes above that rate
func (l *Limiter) SetLimit(jobType string, ratePerSecond float64, burst int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limiters[jobType] = rate.NewLimiter(rate.Limit(ratePerSecond), burst)
}

// Wait blocks until a token is available for jobType, or ctx is
// cancled; job types with no configured limit return immediately
func (l *Limiter) Wait(ctx context.Context, jobType string) error {
	l.mu.RLock()
	lim, ok := l.limiters[jobType]
	l.mu.RUnlock()

	if !ok {
		return nil // no limit configured for this job type
	}

	return lim.Wait(ctx)
}
