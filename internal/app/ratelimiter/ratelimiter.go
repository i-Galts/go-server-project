package ratelimiter

import (
	"fmt"
	"sync"
	"time"

	"github.com/i-Galts/go-server-project/internal/app/storage"
)

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

func (rl *RateLimiter) Permit(ip string) bool {
	return rl.getBucket(ip).Permit()
}

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
				fmt.Println("error getting client config by ip: ", err)
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
