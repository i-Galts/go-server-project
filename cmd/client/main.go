package main

import (
	"fmt"
	"io"
	"net/http"
)

const (
	url          = "http://localhost:8080"
	max_requests = 500
)

// a client sends many GET requests to a server
func main() {
	for i := 0; i < max_requests; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("[client]: error in request for %s: %s\n", url, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[client]: error reading response body: %s\n", err)
			return
		}

		if len(body) != 0 {
			fmt.Printf("[client]: the server answered:\n%s\n", string(body))
		}
	}
}
