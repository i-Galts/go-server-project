package loadbalancer

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/i-Galts/go-server-project/internal/app/backend"
	"github.com/i-Galts/go-server-project/internal/app/ratelimiter"
)

type ILoadBalancer interface {
	Add(backend.IBackend)
	Serve(w http.ResponseWriter, r *http.Request)
}

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

// use of atomics is faster than mutex
func (lb *LoadBalancer) GetNextIndex() uint64 {
	return atomic.AddUint64(&lb.Current, uint64(1)) % uint64(len(lb.Backends))
}

// round-robin fashion
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
		fmt.Printf("error: too many requests from client %s\n", clientIP)
		return
	}

	b := lb.GetNextBackend()
	if b == nil {
		fmt.Println("all backends are down")
		return
	}
	// w.Header().Add("X-Forwarded-Server", b.GetURL())
	b.Serve(w, r)
}
