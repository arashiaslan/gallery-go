package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt string
}

// CreateSession inserts a new session
func (d *DB) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, ?, datetime('now'))
	`
	_, err := d.ExecContext(ctx, query, session.ID, session.UserID, session.Token, session.ExpiresAt)
	return err
}

// GetSessionByToken fetches a session by token, ignoring expired sessions
func (d *DB) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	session := &Session{}
	query := `SELECT id, user_id, token, expires_at, created_at
	          FROM sessions WHERE token = ? AND expires_at > ?`
	err := d.QueryRowContext(ctx, query, token, time.Now()).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

// DeleteSession removes a session by token
func (d *DB) DeleteSession(ctx context.Context, token string) error {
	_, err := d.ExecContext(ctx, "DELETE FROM sessions WHERE token = ?", token)
	return err
}

// DeleteExpiredSessions removes all expired sessions
func (d *DB) DeleteExpiredSessions(ctx context.Context) error {
	_, err := d.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}
