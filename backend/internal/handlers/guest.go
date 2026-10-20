package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/squarebrain/gallery/internal/store"
)

// GuestHandler serves read-only endpoints that require no authentication.
// Anyone can browse albums and view photos; mutations stay behind auth.
type GuestHandler struct {
	db *store.DB
}

func NewGuestHandler(db *store.DB) *GuestHandler {
	return &GuestHandler{db: db}
}

// ListBooks handles GET /api/guest/books
func (h *GuestHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.db.GetAllBooks(r.Context())
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

// GetBook handles GET /api/guest/books/:id
func (h *GuestHandler) GetBook(w http.ResponseWriter, r *http.Request) {
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

	WriteSuccess(w, http.StatusOK, bookToResponse(book))
}

// GetBookPages handles GET /api/guest/books/:id/pages
func (h *GuestHandler) GetBookPages(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")

	book, err := h.db.GetBookByID(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}

	photos, err := h.db.GetPhotosByBook(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch photos")
		return
	}
	if photos == nil {
		photos = []*store.Photo{}
	}

	responses := make([]PhotoResponse, len(photos))
	for i, p := range photos {
		responses[i] = photoToResponse(p)
	}

	WriteSuccess(w, http.StatusOK, PagedPhotosResponse{
		Photos: responses,
		Count:  len(responses),
	})
}
