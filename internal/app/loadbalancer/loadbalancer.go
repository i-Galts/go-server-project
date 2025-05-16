package loadbalancer

import (
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/i-Galts/go-server-project/internal/app/backend"
	"github.com/i-Galts/go-server-project/internal/app/logger"
	"github.com/i-Galts/go-server-project/internal/app/ratelimiter"
)

type ILoadBalancer interface {
	Add(backend.IBackend)
	Serve(w http.ResponseWriter, r *http.Request)
}

// LoadBalancer is a round-robin load balancer with health checks and rate limiting
type LoadBalancer struct {
	Backends    []backend.IBackend
	Current     uint64
	Mutex       sync.Mutex
	RateLimiter *ratelimiter.RateLimiter
}

func NewLoadBalancer(rl *ratelimiter.RateLimiter) *LoadBalancer {
	return &LoadBalancer{
		RateLimiter: rl,
	}
}

func (lb *LoadBalancer) Add(backend backend.IBackend) {
	lb.Backends = append(lb.Backends, backend)
}

// calculates the next index for round-robin selection using atomic operations
// use of atomics is faster than mutex
func (lb *LoadBalancer) GetNextIndex() uint64 {
	return atomic.AddUint64(&lb.Current, uint64(1)) % uint64(len(lb.Backends))
}

// selects the next healthy backend using a round-robin strategy
func (lb *LoadBalancer) GetNextBackend() backend.IBackend {
	size := uint64(len(lb.Backends))
	next := lb.GetNextIndex()

	for i := next; i < size+next; i++ {
		ind := i % size
		backend := lb.Backends[ind]
		if backend.IsAlive() {
			if i != next {
				atomic.StoreUint64(&lb.Current, ind)
			}
			return backend
		}
	}

	return nil
}

// Serve forwards the incoming HTTP request to a healthy backend after:
// - Extracting client IP for rate limiting
// - Checking if the client has exceeded allowed request rate
// - Selecting the next available backend
func (lb *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	var extractIP = func(rr *http.Request) string {
		ip := rr.RemoteAddr
		if i := strings.LastIndex(ip, ":"); i != -1 {
			ip = ip[:i]
		}
		return ip
	}

	clientIP := extractIP(r)

	if !lb.RateLimiter.Permit(clientIP) {
		logger.Log.Errorf("error: too many requests from client %s\n", clientIP)
		return
	}

	b := lb.GetNextBackend()
	if b == nil {
		logger.Log.Errorf("all backends are down")
		return
	}
	// debug purposes
	// w.Header().Add("X-Forwarded-Server", b.GetURL())
	b.Serve(w, r)
}
