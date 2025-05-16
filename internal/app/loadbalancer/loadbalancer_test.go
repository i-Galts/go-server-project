package loadbalancer_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/i-Galts/go-server-project/internal/app/backend"
	"github.com/i-Galts/go-server-project/internal/app/loadbalancer"
	"github.com/i-Galts/go-server-project/internal/app/ratelimiter"

	"github.com/stretchr/testify/assert"
)

type mockBackend struct {
	isAlive bool
	url     string
}

func (m *mockBackend) IsAlive() bool {
	return m.isAlive
}

func (m *mockBackend) SetAlive(alive bool) {
	m.isAlive = alive
}

func (m *mockBackend) Serve(w http.ResponseWriter, r *http.Request) {
	// left empty
}

func (m *mockBackend) GetURL() string {
	return m.url
}

func TestGetNextBackend_SkipsDownedBackends(t *testing.T) {
	lb := &loadbalancer.LoadBalancer{
		Backends: []backend.IBackend{
			&mockBackend{isAlive: false, url: "http://down1"},
			&mockBackend{isAlive: true, url: "http://up1"},
			&mockBackend{isAlive: false, url: "http://down2"},
			&mockBackend{isAlive: true, url: "http://up2"},
		},
	}

	next := lb.GetNextBackend()
	assert.Equal(t, "http://up1", next.GetURL())

	atomic.StoreUint64(&lb.Current, 1)
	next = lb.GetNextBackend()
	assert.Equal(t, "http://up2", next.GetURL())
}

func TestGetNextBackend_AllDown_ReturnsNil(t *testing.T) {
	lb := &loadbalancer.LoadBalancer{
		Backends: []backend.IBackend{
			&mockBackend{isAlive: false, url: "http://down1"},
			&mockBackend{isAlive: false, url: "http://down2"},
		},
	}

	next := lb.GetNextBackend()
	assert.Nil(t, next)
}

func TestServe_ForwardsToBackend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	be := &backend.Backend{
		URL:          u,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
		Mutex:        sync.RWMutex{},
	}

	lb := &loadbalancer.LoadBalancer{
		Backends:    []backend.IBackend{be},
		RateLimiter: ratelimiter.NewLimiter(10, 60),
	}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	lb.Serve(w, req)

	assert.Equal(t, "OK\n", w.Body.String())
}
