package storage_test

import (
	"os"
	"testing"
)

var (
	dbURL string
)

func TestMain(m *testing.M) {
	dbURL = os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=ilyagaltsov password=mysecretpassword dbname=project_test sslmode=disable"
	}

	os.Exit(m.Run())
}
