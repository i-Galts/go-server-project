package backend

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/i-Galts/go-server-project/internal/app/server"
)

type IBackend interface {
	IsAlive() bool
	SetAlive(alive bool)
	Serve(w http.ResponseWriter, r *http.Request)
	GetURL() string
}

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
	Mutex        sync.RWMutex
}

func (b *Backend) GetURL() string {
	return b.URL.String()
}

func (b *Backend) SetAlive(alive bool) {
	b.Mutex.Lock()
	b.Alive = alive
	b.Mutex.Unlock()
}

func (b *Backend) IsAlive() (res bool) {
	b.Mutex.RLock()
	res = b.Alive
	b.Mutex.RUnlock()
	return
}

func (b *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	b.ReverseProxy.ServeHTTP(w, r)
}

func MonitorBackend(backend *Backend, checkInterval time.Duration) {
	for range time.Tick(checkInterval) {
		res, err := http.Head(backend.URL.String())

		if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
			fmt.Printf("%s is down\n", backend.URL)
			backend.SetAlive(false)
		} else {
			backend.SetAlive(true)
		}
	}
}

func RunBackends(config *server.ServerConfig) []*Backend {
	var backends []*Backend

	for _, u := range config.Backends {
		url, _ := url.Parse(u)

		b := &Backend{
			URL:          url,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		}

		backends = append(backends, b)
		checkInterval, err := time.ParseDuration(config.CheckInterval)
		if err != nil {
			fmt.Println("error parsing check interval duration:", err)
		}
		go MonitorBackend(b, checkInterval)
	}

	return backends
}
