package handlers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/squarebrain/gallery/internal/media"
	"github.com/squarebrain/gallery/internal/store"
)

type MediaHandler struct {
	db      *store.DB
	dataDir string
}

func NewMediaHandler(db *store.DB, dataDir string) *MediaHandler {
	return &MediaHandler{db: db, dataDir: dataDir}
}

// ServeMedia handles GET /media/:book_id/:photo_id/:size
// Public read (required by the guest view); existence checks still apply.
func (h *MediaHandler) ServeMedia(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "book_id")
	photoID := chi.URLParam(r, "photo_id")
	size := chi.URLParam(r, "size")

	// Validate size
	if size != "thumb" && size != "medium" && size != "original" {
		WriteError(w, http.StatusBadRequest, "Invalid size")
		return
	}

	// Verify the book exists
	book, err := h.db.GetBookByID(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}

	// Check photo exists
	photo, err := h.db.GetPhotoByID(r.Context(), photoID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if photo == nil {
		WriteError(w, http.StatusNotFound, "Photo not found")
		return
	}

	// Build safe file path (prevent path traversal)
	filePath := media.GetPhotoPath(h.dataDir, bookID, photoID, size)
	if filePath == "" {
		WriteError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	// Check if file exists - if original not found, try with different extension
	if size == "original" {
		// Try common extensions
		baseDir := filepath.Join(h.dataDir, "uploads", bookID, photoID)
		matches, _ := filepath.Glob(filepath.Join(baseDir, "original.*"))
		if len(matches) > 0 {
			filePath = matches[0]
		}
	}

	// Serve file with caching headers
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		WriteError(w, http.StatusNotFound, "File not found")
		return
	}

	// Calculate ETag
	etag := generateETag(filePath, fileInfo)

	// Check If-None-Match
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Set caching headers
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year for immutable content

	// Determine MIME type
	mimeType := getMimeType(size)
	w.Header().Set("Content-Type", mimeType)

	// Serve file
	file, err := os.Open(filePath)
	if err != nil {
		WriteError(w, http.StatusNotFound, "File not found")
		return
	}
	defer file.Close()

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

// generateETag creates an ETag from file path and info
func generateETag(filePath string, fileInfo os.FileInfo) string {
	hash := md5.Sum([]byte(filePath + fileInfo.ModTime().String()))
	return fmt.Sprintf(`"%x"`, hash)
}

// getMimeType returns MIME type based on size
func getMimeType(size string) string {
	switch size {
	case "thumb", "medium":
		return "image/jpeg"
	case "original":
		return "image/jpeg" // Default, may vary
	default:
		return "application/octet-stream"
	}
}
