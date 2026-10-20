package handlers

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/squarebrain/gallery/internal/auth"
	"github.com/squarebrain/gallery/internal/media"
	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
)

// createTestImage creates a valid JPEG image
func createTestImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	buf := &bytes.Buffer{}
	jpeg.Encode(buf, img, nil)
	return buf.Bytes()
}

func TestUploadPhotos(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create user and book
	hash, _ := auth.HashPassword("password123")
	user := &store.User{
		ID:           utils.GenerateShortID(),
		Email:        "test@example.com",
		PasswordHash: hash,
		Role:         "member",
	}
	db.CreateUser(context.Background(), user)

	book := &store.Book{
		ID:      utils.GenerateShortID(),
		Title:   "Test Album",
		OwnerID: user.ID,
	}
	db.CreateBook(context.Background(), book)

	// Create test image
	imageData := createTestImage()

	// Create multipart form
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, _ := writer.CreateFormFile("photos", "test.jpg")
	part.Write(imageData)
	writer.Close()

	// Create processor
	dataDir := t.TempDir()
	processor := media.NewProcessor(dataDir, 1, db)
	defer processor.Shutdown(context.Background())

	// Create handler
	photoHandler := NewPhotoHandler(db, dataDir, processor, 25)

	// Mock request with user context
	req := httptest.NewRequest("POST", "/api/books/"+book.ID+"/photos",
		bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(context.WithValue(req.Context(), "user", user))
	req = req.WithContext(setChiParam(req.Context(), "id", book.ID))

	w := httptest.NewRecorder()
	photoHandler.UploadPhotos(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestUploadPhotoUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create book owned by different user
	user1 := &store.User{
		ID:   utils.GenerateShortID(),
		Role: "member",
	}
	db.CreateUser(context.Background(), user1)

	user2 := &store.User{
		ID:   utils.GenerateShortID(),
		Role: "member",
	}
	db.CreateUser(context.Background(), user2)

	book := &store.Book{
		ID:      utils.GenerateShortID(),
		Title:   "Test",
		OwnerID: user1.ID,
	}
	db.CreateBook(context.Background(), book)

	dataDir := t.TempDir()
	processor := media.NewProcessor(dataDir, 1, db)
	defer processor.Shutdown(context.Background())

	photoHandler := NewPhotoHandler(db, dataDir, processor, 25)

	req := httptest.NewRequest("POST", "/api/books/"+book.ID+"/photos",
		bytes.NewReader([]byte("")))
	req = req.WithContext(context.WithValue(req.Context(), "user", user2))
	req = req.WithContext(setChiParam(req.Context(), "id", book.ID))

	w := httptest.NewRecorder()
	photoHandler.UploadPhotos(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", w.Code)
	}
}

func TestDeletePhoto(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := &store.User{
		ID:   utils.GenerateShortID(),
		Role: "member",
	}
	db.CreateUser(context.Background(), user)

	book := &store.Book{
		ID:      utils.GenerateShortID(),
		OwnerID: user.ID,
	}
	db.CreateBook(context.Background(), book)

	photo := &store.Photo{
		ID:     utils.GenerateShortID(),
		BookID: book.ID,
	}
	db.CreatePhoto(context.Background(), photo)

	dataDir := t.TempDir()
	processor := media.NewProcessor(dataDir, 1, db)
	defer processor.Shutdown(context.Background())

	photoHandler := NewPhotoHandler(db, dataDir, processor, 25)

	req := httptest.NewRequest("DELETE", "/api/photos/"+photo.ID, nil)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	// Set chi params
	req = req.WithContext(setChiParam(req.Context(), "id", photo.ID))

	w := httptest.NewRecorder()
	photoHandler.DeletePhoto(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Verify photo was deleted
	deleted, _ := db.GetPhotoByID(context.Background(), photo.ID)
	if deleted != nil {
		t.Error("Photo should be deleted")
	}
}

// Helper to set chi params in context
func setChiParam(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{key}, Values: []string{value}},
	})
}
