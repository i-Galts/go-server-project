package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type ClientConfig struct {
	Capacity   int
	RefillRate int
}

type Storage struct {
	database *sql.DB
}

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

func (s *Storage) Close() {
	s.database.Close()
}
