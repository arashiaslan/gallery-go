package store

import (
	"context"
	"database/sql"
	"errors"
)

type Book struct {
	ID            string
	Title         string
	OwnerID       string
	CoverPhotoID  *string
	BGColor       string
	CreatedAt     string
	UpdatedAt     string
}

// CreateBook inserts a new book
func (d *DB) CreateBook(ctx context.Context, book *Book) error {
	query := `
		INSERT INTO books (id, title, owner_id, cover_photo_id, bg_color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`
	_, err := d.ExecContext(ctx, query, book.ID, book.Title, book.OwnerID, book.CoverPhotoID, book.BGColor)
	return err
}

// GetBookByID fetches a book by ID
func (d *DB) GetBookByID(ctx context.Context, id string) (*Book, error) {
	book := &Book{}
	query := `SELECT id, title, owner_id, cover_photo_id, bg_color, created_at, updated_at FROM books WHERE id = ?`
	err := d.QueryRowContext(ctx, query, id).Scan(
		&book.ID, &book.Title, &book.OwnerID, &book.CoverPhotoID, &book.BGColor, &book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return book, nil
}

// GetBooksByOwner returns all books for a user
func (d *DB) GetBooksByOwner(ctx context.Context, ownerID string) ([]*Book, error) {
	query := `SELECT id, title, owner_id, cover_photo_id, bg_color, created_at, updated_at FROM books WHERE owner_id = ? ORDER BY updated_at DESC`
	rows, err := d.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		book := &Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.OwnerID, &book.CoverPhotoID, &book.BGColor, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

// GetAllBooks returns every book (public guest view)
func (d *DB) GetAllBooks(ctx context.Context) ([]*Book, error) {
	query := `SELECT id, title, owner_id, cover_photo_id, bg_color, created_at, updated_at FROM books ORDER BY updated_at DESC`
	rows, err := d.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		book := &Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.OwnerID, &book.CoverPhotoID, &book.BGColor, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

// UpdateBook updates book metadata
func (d *DB) UpdateBook(ctx context.Context, book *Book) error {
	query := `
		UPDATE books SET title = ?, cover_photo_id = ?, bg_color = ?, updated_at = datetime('now')
		WHERE id = ?
	`
	result, err := d.ExecContext(ctx, query, book.Title, book.CoverPhotoID, book.BGColor, book.ID)
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

// DeleteBook removes a book and its photos atomically
func (d *DB) DeleteBook(ctx context.Context, id string) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete associated photos first
	if _, err := tx.ExecContext(ctx, "DELETE FROM photos WHERE book_id = ?", id); err != nil {
		return err
	}
	// Then delete the book
	if _, err := tx.ExecContext(ctx, "DELETE FROM books WHERE id = ?", id); err != nil {
		return err
	}
	return tx.Commit()
}
