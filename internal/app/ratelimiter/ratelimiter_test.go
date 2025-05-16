package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/i-Galts/go-server-project/internal/app/ratelimiter"
	"github.com/i-Galts/go-server-project/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	getConfigFunc func(ip string) (storage.ClientConfig, error)
}

func (m *mockStorage) GetClientConfig(ip string) (storage.ClientConfig, error) {
	return m.getConfigFunc(ip)
}

func TestBucket_Permit_AllowsTokens(t *testing.T) {
	b := &ratelimiter.Bucket{
		Capacity:     3,
		Tokens:       3,
		RefillRate:   1,
		LastRefilled: time.Now(),
	}

	assert.True(t, b.Permit())
	assert.True(t, b.Permit())
	assert.True(t, b.Permit())
	assert.False(t, b.Permit())
}

func TestBucket_RefillsOverTime(t *testing.T) {
	now := time.Now()
	b := &ratelimiter.Bucket{
		Capacity:     3,
		Tokens:       0,
		RefillRate:   2,
		LastRefilled: now.Add(-2 * time.Second),
	}

	assert.True(t, b.Permit())
	assert.Equal(t, b.Tokens, 2)
}

func TestBucket_RefillDoesNotExceedCapacity(t *testing.T) {
	now := time.Now()
	b := &ratelimiter.Bucket{
		Capacity:     5,
		Tokens:       0,
		RefillRate:   2,
		LastRefilled: now.Add(-10 * time.Second),
	}

	assert.True(t, b.Permit())
	assert.Equal(t, b.Tokens, 4)
}

func TestRateLimiter_Permit_AllowsUpToLimit(t *testing.T) {
	rl := ratelimiter.NewLimiter(2, 1)

	ip := "192.168.1.1"

	assert.True(t, rl.Permit(ip))
	assert.True(t, rl.Permit(ip))
	assert.False(t, rl.Permit(ip))
}
