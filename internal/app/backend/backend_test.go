package backend_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/i-Galts/go-server-project/internal/app/backend"
	"github.com/stretchr/testify/assert"
)

func TestBackend_IsAlive_SetAlive(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080")
	backend := &backend.Backend{
		URL:          u,
		Alive:        false,
		ReverseProxy: nil,
	}

	assert.False(t, backend.IsAlive())

	backend.SetAlive(true)
	assert.True(t, backend.IsAlive())
}

func TestMonitorBackend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	url, _ := url.Parse(server.URL)
	b := &backend.Backend{
		URL:          url,
		Alive:        false,
		ReverseProxy: nil,
	}

	checkInterval := 100 * time.Millisecond
	go backend.MonitorBackend(b, checkInterval)

	time.Sleep(200 * time.Millisecond)

	assert.True(t, b.IsAlive())
}
