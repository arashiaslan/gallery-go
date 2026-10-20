package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string // "admin" or "member"
	CreatedAt    string
	UpdatedAt    string
}

// CreateUser inserts a new user
func (d *DB) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`
	_, err := d.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Role)
	return err
}

// GetUserByEmail fetches a user by email
func (d *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE email = ?`
	err := d.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID fetches a user by ID
func (d *DB) GetUserByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, password_hash, role, created_at, updated_at FROM users WHERE id = ?`
	err := d.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// CountUsers returns the number of users in the database
func (d *DB) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := d.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}
