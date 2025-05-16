package ratelimiter

import (
	"math"
	"sync"
	"time"
)

type Bucket struct {
	Capacity     int
	Tokens       int
	RefillRate   int
	LastRefilled time.Time
	Mutex        sync.Mutex
}

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
