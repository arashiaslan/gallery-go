package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Photo struct {
	ID        string
	BookID    string
	Filename  string
	Width     int
	Height    int
	TakenAt   *time.Time
	Caption   *string
	Position  int
	CreatedAt string
}

// CreatePhoto inserts a new photo
func (d *DB) CreatePhoto(ctx context.Context, photo *Photo) error {
	query := `
		INSERT INTO photos (id, book_id, filename, width, height, taken_at, caption, position, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`
	_, err := d.ExecContext(ctx, query, photo.ID, photo.BookID, photo.Filename,
		photo.Width, photo.Height, photo.TakenAt, photo.Caption, photo.Position)
	return err
}

// GetPhotoByID fetches a photo by ID
func (d *DB) GetPhotoByID(ctx context.Context, id string) (*Photo, error) {
	photo := &Photo{}
	query := `SELECT id, book_id, filename, width, height, taken_at, caption, position, created_at
	          FROM photos WHERE id = ?`
	err := d.QueryRowContext(ctx, query, id).Scan(
		&photo.ID, &photo.BookID, &photo.Filename, &photo.Width, &photo.Height,
		&photo.TakenAt, &photo.Caption, &photo.Position, &photo.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return photo, nil
}

// GetPhotosByBook returns all photos for a book in order
func (d *DB) GetPhotosByBook(ctx context.Context, bookID string) ([]*Photo, error) {
	query := `SELECT id, book_id, filename, width, height, taken_at, caption, position, created_at
	          FROM photos WHERE book_id = ? ORDER BY position ASC`
	rows, err := d.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*Photo
	for rows.Next() {
		photo := &Photo{}
		if err := rows.Scan(&photo.ID, &photo.BookID, &photo.Filename, &photo.Width, &photo.Height,
			&photo.TakenAt, &photo.Caption, &photo.Position, &photo.CreatedAt); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, rows.Err()
}

// UpdatePhotoCaption updates just the caption
func (d *DB) UpdatePhotoCaption(ctx context.Context, id string, caption *string) error {
	query := `UPDATE photos SET caption = ? WHERE id = ?`
	result, err := d.ExecContext(ctx, query, caption, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UpdatePhotoTakenAt updates only the taken_at column (EXIF correction)
func (d *DB) UpdatePhotoTakenAt(ctx context.Context, id string, takenAt time.Time) error {
	query := `UPDATE photos SET taken_at = ? WHERE id = ?`
	result, err := d.ExecContext(ctx, query, takenAt, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeletePhoto removes a photo
func (d *DB) DeletePhoto(ctx context.Context, id string) error {
	_, err := d.ExecContext(ctx, "DELETE FROM photos WHERE id = ?", id)
	return err
}

// ReorderPhotos updates positions of multiple photos within a single book.
// It rejects any photo ID that does not belong to the given book (IDOR guard).
func (d *DB) ReorderPhotos(ctx context.Context, bookID string, updates map[string]int) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for photoID, position := range updates {
		query := `UPDATE photos SET position = ? WHERE id = ? AND book_id = ?`
		result, err := tx.ExecContext(ctx, query, position, photoID, bookID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return ErrPhotoNotInBook
		}
	}

	return tx.Commit()
}

// ErrPhotoNotInBook is returned when a reorder target does not belong to the book.
var ErrPhotoNotInBook = errors.New("photo does not belong to book")

// GetPhotoCount returns number of photos in a book
func (d *DB) GetPhotoCount(ctx context.Context, bookID string) (int, error) {
	var count int
	err := d.QueryRowContext(ctx, "SELECT COUNT(*) FROM photos WHERE book_id = ?", bookID).Scan(&count)
	return count, err
}
