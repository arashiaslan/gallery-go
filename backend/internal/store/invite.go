package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Invite struct {
	ID          string
	Email       string
	Token       string
	CreatedByID string
	UsedByID    *string
	ExpiresAt   time.Time
	UsedAt      *time.Time
	CreatedAt   string
}

// CreateInvite inserts a new invite
func (d *DB) CreateInvite(ctx context.Context, invite *Invite) error {
	query := `
		INSERT INTO invites (id, email, token, created_by_id, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'))
	`
	_, err := d.ExecContext(ctx, query, invite.ID, invite.Email, invite.Token, invite.CreatedByID, invite.ExpiresAt)
	return err
}

// GetInviteByToken fetches an invite by token
func (d *DB) GetInviteByToken(ctx context.Context, token string) (*Invite, error) {
	invite := &Invite{}
	query := `SELECT id, email, token, created_by_id, used_by_id, expires_at, used_at, created_at FROM invites WHERE token = ?`
	err := d.QueryRowContext(ctx, query, token).Scan(
		&invite.ID, &invite.Email, &invite.Token, &invite.CreatedByID, &invite.UsedByID,
		&invite.ExpiresAt, &invite.UsedAt, &invite.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return invite, nil
}

// MarkInviteAsUsed marks an invite as used
func (d *DB) MarkInviteAsUsed(ctx context.Context, token string, userID string) error {
	query := `UPDATE invites SET used_by_id = ?, used_at = datetime('now') WHERE token = ?`
	_, err := d.ExecContext(ctx, query, userID, token)
	return err
}
