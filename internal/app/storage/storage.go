package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	config   *StorageConfig
	database *sql.DB
	userRepo *UserRepo
}

func NewStorage(config *StorageConfig) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.config.URL)
	if err != nil {
		return nil
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	s.database = db

	return nil
}

func (s *Storage) Close() {
	s.database.Close()
}

func (s *Storage) CreateUserRepo() *UserRepo {
	if s.userRepo != nil {
		return nil
	}

	s.userRepo = &UserRepo{
		storage: s,
	}

	return s.userRepo
}
