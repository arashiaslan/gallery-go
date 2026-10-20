package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/squarebrain/gallery/internal/media"
	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
)

type BookHandler struct {
	db      *store.DB
	dataDir string
}

func NewBookHandler(db *store.DB, dataDir string) *BookHandler {
	return &BookHandler{db: db, dataDir: dataDir}
}

type CreateBookRequest struct {
	Title   string `json:"title"`
	BGColor string `json:"bg_color,omitempty"`
}

type BookResponse struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	OwnerID      string  `json:"owner_id"`
	CoverPhotoID *string `json:"cover_photo_id"`
	BGColor      string  `json:"bg_color"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

func bookToResponse(b *store.Book) BookResponse {
	return BookResponse{
		ID:           b.ID,
		Title:        b.Title,
		OwnerID:      b.OwnerID,
		CoverPhotoID: b.CoverPhotoID,
		BGColor:      b.BGColor,
		CreatedAt:    b.CreatedAt,
		UpdatedAt:    b.UpdatedAt,
	}
}

// CreateBook handles POST /api/books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)

	var req CreateBookRequest
	LimitBody(w, r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Title == "" {
		WriteError(w, http.StatusBadRequest, "Title is required")
		return
	}

	bgColor := req.BGColor
	if bgColor == "" {
		bgColor = "#8E8DB8"
	}

	book := &store.Book{
		ID:      utils.GenerateShortID(),
		Title:   req.Title,
		OwnerID: user.ID,
		BGColor: bgColor,
	}

	if err := h.db.CreateBook(r.Context(), book); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create book")
		return
	}

	WriteSuccess(w, http.StatusCreated, bookToResponse(book))
}

// ListBooks handles GET /api/books
func (h *BookHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)

	books, err := h.db.GetBooksByOwner(r.Context(), user.ID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch books")
		return
	}

	if books == nil {
		books = []*store.Book{}
	}

	responses := make([]BookResponse, len(books))
	for i, b := range books {
		responses[i] = bookToResponse(b)
	}

	WriteSuccess(w, http.StatusOK, responses)
}

// GetBook handles GET /api/books/:id
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	id := chi.URLParam(r, "id")

	book, err := h.db.GetBookByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch book")
		return
	}

	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}

	// Check authorization
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	WriteSuccess(w, http.StatusOK, bookToResponse(book))
}

type UpdateBookRequest struct {
	Title       string `json:"title,omitempty"`
	BGColor     string `json:"bg_color,omitempty"`
	CoverPhotoID *string `json:"cover_photo_id"`
}

// UpdateBook handles PATCH /api/books/:id
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	id := chi.URLParam(r, "id")

	book, err := h.db.GetBookByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch book")
		return
	}

	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}

	// Check authorization
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req UpdateBookRequest
	LimitBody(w, r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.BGColor != "" {
		book.BGColor = req.BGColor
	}
	if req.CoverPhotoID != nil {
		book.CoverPhotoID = req.CoverPhotoID
	}

	if err := h.db.UpdateBook(r.Context(), book); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update book")
		return
	}

	WriteSuccess(w, http.StatusOK, bookToResponse(book))
}

// DeleteBook handles DELETE /api/books/:id
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	id := chi.URLParam(r, "id")

	book, err := h.db.GetBookByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch book")
		return
	}

	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}

	// Check authorization
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Remove DB rows first (atomic), then the files on disk. If file removal
	// fails we log it: the DB is the source of truth and must stay consistent.
	if err := h.db.DeleteBook(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete book")
		return
	}

	if err := media.RemoveBookFiles(h.dataDir, id); err != nil {
		log.Printf("failed to remove files for book %s: %v", id, err)
	}

	WriteSuccess(w, http.StatusOK, map[string]string{"message": "Book deleted"})
}
