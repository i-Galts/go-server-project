package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/i-Galts/go-server-project/internal/app/backend"
	"github.com/i-Galts/go-server-project/internal/app/loadbalancer"
	"github.com/i-Galts/go-server-project/internal/app/ratelimiter"
	"github.com/i-Galts/go-server-project/internal/app/server"
	"github.com/i-Galts/go-server-project/internal/app/storage"
)

var configPath = flag.String("conf-path", "configs/lb_conf.json", "Path to server config file")

func main() {
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		fmt.Println("failed to open server config file:", err)
		log.Fatal(err)
		return
	}
	defer file.Close()

	conf, err := server.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("failed to load server config:", err)
	}

	if len(conf.Backends) == 0 {
		fmt.Println("no backends given in the configuration")
		os.Exit(0)
	}

	// rate limiter
	db, err := storage.NewStorage("clients.db")
	if err != nil {
		fmt.Println("error opening clients database:", err)
	}

	rl := ratelimiter.NewLimiter(
		conf.RateLimiterCap,
		conf.RateLimiterRefillRate,
	)
	rl.ClientStorage = db

	backends := backend.RunBackends(&conf)
	lb := loadbalancer.NewLoadBalancer(rl)
	for _, b := range backends {
		lb.Add(b)
	}

	http.HandleFunc("/", lb.Serve)

	s := server.NewServer(&conf)
	err = s.Launch()
	if err != nil {
		fmt.Println("error launching server:", err)
	}

	err = http.ListenAndServe(conf.Port, nil)
	if err != nil {
		fmt.Println("error running server", err)
	}
}
