package storage

import (
	"fmt"
	"strings"
	"testing"
)

func TestStorage(t *testing.T, url string) (*Storage, func(...string)) {
	t.Helper()

	conf := NewStorageConfig()
	conf.URL = url
	s := NewStorage(conf)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			_, err := s.database.Exec(
				fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}

		s.Close()
	}
}
