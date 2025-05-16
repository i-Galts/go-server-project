// Package ratelimiter provides a token bucket-based rate limiting implementation.
package ratelimiter

import (
	"sync"
	"time"

	"github.com/i-Galts/go-server-project/internal/app/logger"
	"github.com/i-Galts/go-server-project/internal/app/storage"
)

// manages rate limits for multiple clients (e.g., by IP address)
type RateLimiter struct {
	Buckets       map[string]*Bucket
	Mutex         sync.RWMutex
	Capacity      int
	RefillRate    int
	ClientStorage *storage.Storage
}

func NewLimiter(capacity, refillRate int) *RateLimiter {
	return &RateLimiter{
		Buckets:    make(map[string]*Bucket),
		Capacity:   capacity,
		RefillRate: refillRate,
	}
}

// checks whether a request from the specified client (identified by IP) is allowed
func (rl *RateLimiter) Permit(ip string) bool {
	return rl.getBucket(ip).Permit()
}

// retrieves or creates a token bucket for the given client IP
func (rl *RateLimiter) getBucket(ip string) *Bucket {
	rl.Mutex.Lock()
	defer rl.Mutex.Unlock()

	capacity := rl.Capacity
	refillRate := rl.RefillRate
	bucket, exists := rl.Buckets[ip]
	if !exists {
		if rl.ClientStorage != nil {
			conf, err := rl.ClientStorage.GetClientConfig(ip)
			if err != nil {
				logger.Log.Errorf("error getting client config by ip: ", err)
			} else {
				capacity = conf.Capacity
				refillRate = conf.RefillRate
			}
		}

		bucket = &Bucket{
			Capacity:     capacity,
			Tokens:       capacity,
			RefillRate:   refillRate,
			LastRefilled: time.Now(),
		}
		rl.Buckets[ip] = bucket
	}

	return bucket
}
