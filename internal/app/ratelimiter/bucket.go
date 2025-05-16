// Package ratelimiter provides a token bucket-based rate limiting implementation.
// It allows controlling the rate at which actions (e.g., HTTP requests) are permitted,
// preventing abuse or overuse of system resources.
package ratelimiter

import (
	"math"
	"sync"
	"time"
)

// represents a token bucket used for rate limiting
type Bucket struct {
	Capacity     int
	Tokens       int
	RefillRate   int
	LastRefilled time.Time
	Mutex        sync.Mutex
}

// refills tokens based on elapsed time since last check
func (b *Bucket) Permit() bool {
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.LastRefilled).Seconds()
	newTokens := int(elapsed * float64(b.RefillRate))
	if newTokens > 0 {
		b.Tokens = int(math.Min(float64(b.Capacity), float64(b.Tokens+newTokens)))
		b.LastRefilled = now
	}

	if b.Tokens > 0 {
		b.Tokens--
		return true
	}

	return false
}
