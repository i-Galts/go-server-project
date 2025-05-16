package backend

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/i-Galts/go-server-project/internal/app/logger"
	"github.com/i-Galts/go-server-project/internal/app/server"
)

type IBackend interface {
	IsAlive() bool
	SetAlive(alive bool)
	Serve(w http.ResponseWriter, r *http.Request)
	GetURL() string
}

// represents a single backend server in the system
type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
	Mutex        sync.RWMutex
}

func (b *Backend) GetURL() string {
	return b.URL.String()
}

// updates the alive status of the backend in a thread-safe manner
func (b *Backend) SetAlive(alive bool) {
	b.Mutex.Lock()
	b.Alive = alive
	b.Mutex.Unlock()
}

// returns the current alive status of the backend in a thread-safe manner
func (b *Backend) IsAlive() (res bool) {
	b.Mutex.RLock()
	res = b.Alive
	b.Mutex.RUnlock()
	return
}

// serve forwards the incoming HTTP request to this backend using its reverse proxy
func (b *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	b.ReverseProxy.ServeHTTP(w, r)
}

// continuously checks the health of the backend at specified intervals
func MonitorBackend(backend *Backend, checkInterval time.Duration) {
	for range time.Tick(checkInterval) {
		res, err := http.Head(backend.URL.String())

		if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
			logger.Log.Warnf("%s is down", backend.URL)
			backend.SetAlive(false)
		} else {
			backend.SetAlive(true)
		}
	}
}

// initializes and starts monitoring for all backend servers defined in the config
func RunBackends(config *server.ServerConfig) []*Backend {
	var backends []*Backend

	checkInterval, err := time.ParseDuration(config.CheckInterval)
	if err != nil {
		fmt.Println("error parsing check interval duration:", err)
	}

	for _, u := range config.Backends {
		url, _ := url.Parse(u)

		b := &Backend{
			URL:          url,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		}

		backends = append(backends, b)

		go MonitorBackend(b, checkInterval)
	}

	return backends
}
