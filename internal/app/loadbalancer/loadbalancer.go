package loadbalancer

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/i-Galts/go-server-project/internal/app/backend"
)

type ILoadBalancer interface {
	Add(backend.IBackend)
	Serve(w http.ResponseWriter, r *http.Request)
}

type LoadBalancer struct {
	Backends []backend.IBackend
	Current  uint64
	Mutex    sync.Mutex
}

func NewLoadBalancer(cur uint64) *LoadBalancer {
	return &LoadBalancer{Current: cur}
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
	b := lb.GetNextBackend()
	if b == nil {
		fmt.Println("all backends are down")
		return
	}
	w.Header().Add("X-Forwarded-Server", b.GetURL())
	b.Serve(w, r)
}
