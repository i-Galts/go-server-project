package loadbalancer

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

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

// round-robin fashion
func (lb *LoadBalancer) GetNextBackend() backend.IBackend {
	lb.Mutex.Lock()
	defer lb.Mutex.Unlock()

	size := uint64(len(lb.Backends))

	for i := uint64(0); i < size; i++ {
		ind := lb.Current % size
		backend := lb.Backends[ind]
		lb.Current++

		if backend.IsAlive() {
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
