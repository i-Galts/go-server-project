package storage_test

import (
	"testing"

	"github.com/i-Galts/go-server-project/internal/app/model"
	"github.com/i-Galts/go-server-project/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo_Create(t *testing.T) {
	s, truncate := storage.TestStorage(t, dbURL)
	defer truncate("users")

	u, err := s.CreateUserRepo().Create(&model.User{
		Email: "example@example.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepo_FindByEmail(t *testing.T) {
	s, truncate := storage.TestStorage(t, dbURL)
	defer truncate("users")

	email := "example@example.com"

	_, err := s.CreateUserRepo().FindByEmail(email)
	assert.Error(t, err)

	// s.CreateUserRepo().Create(&model.User{
	// 	Email: "example@example.com",
	// })
	// u, err = s.CreateUserRepo().FindByEmail(email)
	// assert.NoError(t, err)
	// assert.NotNil(t, u)
}
