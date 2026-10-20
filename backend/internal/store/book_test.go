package store

import (
	"context"
	"testing"
)

func setupTestDB(t *testing.T) *DB {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestCreateAndGetBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a test user first
	user := &User{
		ID:           "user1",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         "member",
	}
	if err := db.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	book := &Book{
		ID:      "book1",
		Title:   "Test Album",
		OwnerID: "user1",
		BGColor: "#8E8DB8",
	}

	if err := db.CreateBook(context.Background(), book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	retrieved, err := db.GetBookByID(context.Background(), "book1")
	if err != nil {
		t.Fatalf("Failed to get book: %v", err)
	}

	if retrieved == nil {
		t.Error("Book not found")
	}

	if retrieved.Title != "Test Album" {
		t.Errorf("Expected title 'Test Album', got %s", retrieved.Title)
	}
}

func TestGetBooksByOwner(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:           "user1",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         "member",
	}
	if err := db.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	book1 := &Book{ID: "book1", Title: "Album 1", OwnerID: "user1", BGColor: "#8E8DB8"}
	book2 := &Book{ID: "book2", Title: "Album 2", OwnerID: "user1", BGColor: "#8E8DB8"}

	db.CreateBook(context.Background(), book1)
	db.CreateBook(context.Background(), book2)

	books, err := db.GetBooksByOwner(context.Background(), "user1")
	if err != nil {
		t.Fatalf("Failed to get books: %v", err)
	}

	if len(books) != 2 {
		t.Errorf("Expected 2 books, got %d", len(books))
	}
}

func TestDeleteBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &User{
		ID:           "user1",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         "member",
	}
	db.CreateUser(context.Background(), user)

	book := &Book{ID: "book1", Title: "Album 1", OwnerID: "user1", BGColor: "#8E8DB8"}
	db.CreateBook(context.Background(), book)

	if err := db.DeleteBook(context.Background(), "book1"); err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}

	retrieved, _ := db.GetBookByID(context.Background(), "book1")
	if retrieved != nil {
		t.Error("Book should be deleted")
	}
}
