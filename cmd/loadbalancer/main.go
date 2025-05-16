package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/i-Galts/go-server-project/internal/app/backend"
	"github.com/i-Galts/go-server-project/internal/app/loadbalancer"
	"github.com/i-Galts/go-server-project/internal/app/logger"
	"github.com/i-Galts/go-server-project/internal/app/ratelimiter"
	"github.com/i-Galts/go-server-project/internal/app/server"
	"github.com/i-Galts/go-server-project/internal/app/storage"
)

var (
	configPath = flag.String("conf-path", "configs/lb_conf.json", "Path to server config file")
	dbName     = "clients.db"
)

func main() {
	flag.Parse()

	logger.Log.Infof("Starting load balancer...")

	file, err := os.Open(*configPath)
	if err != nil {
		logger.Log.Fatalf("failed to open server config file: %v", err)
		return
	}
	defer file.Close()

	conf, err := server.LoadConfig(*configPath)
	if err != nil {
		logger.Log.Fatalf("failed to load server config: %v", err)
	}

	if len(conf.Backends) == 0 {
		logger.Log.Fatalf("No backends provided in the configuration. Exiting.")
	}

	rlCapacity := conf.RateLimiterCap
	rlRefillRate := conf.RateLimiterRefillRate
	logger.Log.Infof("Initializing rate limiter with capacity=%d, refill_rate=%d",
		rlCapacity, rlRefillRate)

	// rate limiter
	db, err := storage.NewStorage(dbName)
	if err != nil {
		logger.Log.Fatalf("error opening clients database: %v", err)
	}

	rl := ratelimiter.NewLimiter(rlCapacity, rlRefillRate)
	rl.ClientStorage = db

	logger.Log.Infof("Running backends...")
	backends := backend.RunBackends(&conf)
	lb := loadbalancer.NewLoadBalancer(rl)

	for _, b := range backends {
		lb.Add(b)
	}

	// start load balancer
	logger.Log.Infof("Starting load balancer server on port %s", conf.Port)
	http.HandleFunc("/", lb.Serve)
	s := server.NewServer(&conf)
	err = s.Launch()
	if err != nil {
		logger.Log.Errorf("Error launching server: %v", err)
	}

	err = http.ListenAndServe(conf.Port, nil)
	if err != nil {
		logger.Log.Fatalf("Error running server: %v", err)
	}
}
