package storage

import (
	"github.com/i-Galts/go-server-project/internal/app/model"
)

type UserRepo struct {
	storage *Storage
}

func (r *UserRepo) Create(u *model.User) (*model.User, error) {
	err := r.storage.database.QueryRow(
		"INSERT INTO users (email, enc_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncPassword,
	).Scan(&u.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.storage.database.QueryRow(
		"SELECT user_id, email, enc_password WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}
