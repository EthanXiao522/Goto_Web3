package repository

import (
	"database/sql"
	"fmt"

	"github.com/xyd/web3-learning-tracker/internal/model"
)

type UserRepo struct{ DB *sql.DB }

func (r *UserRepo) Create(user *model.User) (uint64, error) {
	result, err := r.DB.Exec(
		`INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`,
		user.Username, user.Email, user.PasswordHash,
	)
	if err != nil {
		return 0, fmt.Errorf("user create: %w", err)
	}
	id, _ := result.LastInsertId()
	return uint64(id), nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = ?`,
		email,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user find by email: %w", err)
	}
	return u, nil
}

func (r *UserRepo) FindByID(id uint64) (*model.User, error) {
	u := &model.User{}
	err := r.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user find by id: %w", err)
	}
	return u, nil
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	u := &model.User{}
	err := r.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = ?`,
		username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user find by username: %w", err)
	}
	return u, nil
}

func (r *UserRepo) Update(user *model.User) error {
	_, err := r.DB.Exec(
		`UPDATE users SET username = ?, email = ?, password_hash = ? WHERE id = ?`,
		user.Username, user.Email, user.PasswordHash, user.ID,
	)
	if err != nil {
		return fmt.Errorf("user update: %w", err)
	}
	return nil
}
