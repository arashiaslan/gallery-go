package store

import (
	"context"
	"testing"
	"time"
)

func TestSessionExpiry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	user := &User{
		ID:           "user1",
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         "member",
	}
	if err := db.CreateUser(ctx, user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	valid := &Session{
		ID:        "s1",
		UserID:    "user1",
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := db.CreateSession(ctx, valid); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	expired := &Session{
		ID:        "s2",
		UserID:    "user1",
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	if err := db.CreateSession(ctx, expired); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	got, err := db.GetSessionByToken(ctx, "valid-token")
	if err != nil {
		t.Fatalf("GetSessionByToken(valid) error: %v", err)
	}
	if got == nil {
		t.Fatal("valid session should be found")
	}

	got, err = db.GetSessionByToken(ctx, "expired-token")
	if err != nil {
		t.Fatalf("GetSessionByToken(expired) error: %v", err)
	}
	if got != nil {
		t.Error("expired session must not be returned")
	}

	if err := db.DeleteExpiredSessions(ctx); err != nil {
		t.Fatalf("DeleteExpiredSessions error: %v", err)
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions").Scan(&count); err != nil {
		t.Fatalf("count error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 session after cleanup, got %d", count)
	}
}

func TestReorderRejectsForeignPhoto(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	user := &User{ID: "user1", Email: "a@example.com", PasswordHash: "hash", Role: "member"}
	if err := db.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	bookA := &Book{ID: "bookA", Title: "A", OwnerID: "user1", BGColor: "#fff"}
	bookB := &Book{ID: "bookB", Title: "B", OwnerID: "user1", BGColor: "#fff"}
	if err := db.CreateBook(ctx, bookA); err != nil {
		t.Fatalf("create bookA: %v", err)
	}
	if err := db.CreateBook(ctx, bookB); err != nil {
		t.Fatalf("create bookB: %v", err)
	}

	photoA := &Photo{ID: "photoA", BookID: "bookA", Filename: "a.jpg"}
	photoB := &Photo{ID: "photoB", BookID: "bookB", Filename: "b.jpg", Position: 3}
	if err := db.CreatePhoto(ctx, photoA); err != nil {
		t.Fatalf("create photoA: %v", err)
	}
	if err := db.CreatePhoto(ctx, photoB); err != nil {
		t.Fatalf("create photoB: %v", err)
	}

	// photoB does not belong to bookA -> must be rejected
	err := db.ReorderPhotos(ctx, "bookA", map[string]int{"photoB": 0})
	if err != ErrPhotoNotInBook {
		t.Errorf("expected ErrPhotoNotInBook, got %v", err)
	}

	// position of photoB must be unchanged
	got, err := db.GetPhotoByID(ctx, "photoB")
	if err != nil {
		t.Fatalf("get photoB: %v", err)
	}
	if got.Position != 3 {
		t.Errorf("photoB position should be unchanged, got %d", got.Position)
	}

	// own photo reorders fine
	if err := db.ReorderPhotos(ctx, "bookA", map[string]int{"photoA": 5}); err != nil {
		t.Errorf("reorder own photo: %v", err)
	}
}
