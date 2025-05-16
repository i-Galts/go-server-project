// Package storage provides persistent storage functionality for client-specific rate limiting configurations.
// It uses SQLite as the backend database and supports operations to get, set, and update client limits.
package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// represents the rate limit configuration for a specific client
type ClientConfig struct {
	Capacity   int
	RefillRate int
}

// manages interactions with the underlying SQLite database
type Storage struct {
	database *sql.DB
}

// initializes a new SQLite-based storage instance at the given file path
func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS client_limits (
			client_id TEXT PRIMARY KEY,
			capacity INTEGER NOT NULL,
			refill_rate INTEGER NOT NULL)
	`)
	if err != nil {
		return nil, err
	}
	return &Storage{database: db}, nil
}

// retrieves the rate limit configuration for a specific client identified by its IP address
func (s *Storage) GetClientConfig(ip string) (*ClientConfig, error) {
	var config ClientConfig
	err := s.database.QueryRow(`
		SELECT capacity, refill_rate FROM client_limits WHERE client_id = ?;`,
		ip,
	).Scan(&config.Capacity, &config.RefillRate)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &config, nil
}

// inserts or updates a client's rate limit configuration in the database
func (s *Storage) SetClientConfig(clientID string, capacity, refillRate int) error {
	_, err := s.database.Exec(`
		INSERT INTO client_limits (client_id, capacity, refill_rate)
		VALUES (?, ?, ?)
		ON CONFLICT(client_id) DO UPDATE SET
			capacity = excluded.capacity,
			refill_rate = excluded.refill_rate
	`, clientID, capacity, refillRate)
	return err
}

// gracefully closes the database connection
func (s *Storage) Close() {
	s.database.Close()
}
