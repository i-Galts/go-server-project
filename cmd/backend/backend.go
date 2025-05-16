package main

import (
	"fmt"
	"net/http"
	"os"
)

// runs a backend server on a given port
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run backend.go :<port>")
		return
	}
	port := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend on port %s\n", port)
	})

	fmt.Printf("[backend]: starting backend on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
