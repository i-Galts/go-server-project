# Go Load Balancer

A simple load balancer written in Go, designed to distribute incoming HTTP requests across multiple backend servers. The project includes a load balancer executable and a test client to simulate traffic.

## Features

- Configurable backend server pool
- Round-robin request distribution
- Basic health checking
- Rate limiting using Token Bucket per client (IP)
- SQLite-based client configuration

## Prerequisites

- Go 1.20 or higher
- GNU Make (for build automation)
- SQLite3 (CLI utility)
- Apache Bench (`ab`) for load testing (optional)

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:i-Galts/go-server-project.git
   cd go-server-project
   ```

2. Build the project:
   ```bash
   make
   ```

This will create two executables in the build/ directory:


lb - The load balancer
backend - A test backend
client - A test client


3. Configuration

The load balancer configuration is stored in build/lb_conf.json after building. It contains load balancer port, health check parameter, list of backend servers and default parameters for a rate limiter.

4. Usage

Running the Load Balancer and the Backends in different terminals:
   ```bash
   ./build/backend 700x
   ```

   ```
   ./build/lb
   ```
   
The load balancer will start on the port specified in the configuration (default: 8080).

5. Testing with the Client

First, start the load balancer as shown above.

In another terminal, run the test client:

   ```bash
   ./build/client
   ```

The client will send multiple requests to the load balancer, which will distribute them to the configured backends. If there is any alive backend server, its answer will be printed (the port of the backend).

6. Rate Limiting (Token Bucket)

Rate limiting is enabled per client IP. Each client has an independent "bucket" that defines:

- Capacity — how many requests can be made in a burst
- Refill rate — how quickly the bucket refills (tokens/second)

Database: clients.db
Each client must be pre-registered in the database with a custom rate limit (or fallback to global config).

Table Structure
The SQLite DB has a table client_limits:
```sql
    CREATE TABLE IF NOT EXISTS client_limits (
       client_id TEXT PRIMARY KEY,
       capacity INTEGER NOT NULL,
       refill_rate INTEGER NOT NULL
    );
```

Adding a New Client
To allow a new IP through the rate limiter:
```bash
   sqlite3 clients.db
```
and
```sql
   INSERT INTO client_limits (client_id, capacity, refill_rate)
   VALUES ('127.0.0.1', 10, 2);
```

7. Testing with Apache Bench

For more comprehensive load testing, you can use Apache Bench:

   ```bash
   ab -n 5000 -c 1000 http://localhost:8080/
   ```

This will send 5000 requests with 1000 concurrent connections to the load balancer.