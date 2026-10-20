package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/squarebrain/gallery/internal/media"
	"github.com/squarebrain/gallery/internal/store"
	"github.com/squarebrain/gallery/internal/utils"
)

type PhotoHandler struct {
	db          *store.DB
	dataDir     string
	processor   *media.Processor
	maxUploadMB int
}

func NewPhotoHandler(db *store.DB, dataDir string, processor *media.Processor, maxUploadMB int) *PhotoHandler {
	return &PhotoHandler{db: db, dataDir: dataDir, processor: processor, maxUploadMB: maxUploadMB}
}

// maxUploadFiles caps how many files a single request may carry.
const maxUploadFiles = 20

type PhotoResponse struct {
	ID       string     `json:"id"`
	BookID   string     `json:"book_id"`
	Filename string     `json:"filename"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	TakenAt  *time.Time `json:"taken_at"`
	Caption  *string    `json:"caption"`
	Position int        `json:"position"`
}

func photoToResponse(p *store.Photo) PhotoResponse {
	return PhotoResponse{
		ID:       p.ID,
		BookID:   p.BookID,
		Filename: p.Filename,
		Width:    p.Width,
		Height:   p.Height,
		TakenAt:  p.TakenAt,
		Caption:  p.Caption,
		Position: p.Position,
	}
}

// UploadPhotos handles POST /api/books/:id/photos
func (h *PhotoHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	bookID := chi.URLParam(r, "id")

	// Check authorization
	book, err := h.db.GetBookByID(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Hard cap on the total request body: per-file limit plus headroom,
	// times the maximum number of files per request.
	maxFileBytes := int64(h.maxUploadMB) << 20
	bodyLimit := (maxFileBytes + (1 << 20)) * maxUploadFiles
	r.Body = http.MaxBytesReader(w, r.Body, bodyLimit)

	// Parse multipart form (in-memory threshold = per-file limit)
	if err := r.ParseMultipartForm(maxFileBytes); err != nil {
		WriteError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	files := r.MultipartForm.File["photos"]
	if len(files) == 0 {
		WriteError(w, http.StatusBadRequest, "No files provided")
		return
	}
	if len(files) > maxUploadFiles {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("Too many files: max %d per request", maxUploadFiles))
		return
	}

	// Get current photo count for position
	photoCount, err := h.db.GetPhotoCount(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	var uploadedPhotos []PhotoResponse
	var uploadErrors []string

	for i, fileHeader := range files {
		func() {
			file, err := fileHeader.Open()
			if err != nil {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: open failed", i+1))
				return
			}
			defer file.Close()

			// Read file data with a per-file size cap
			fileData, err := io.ReadAll(io.LimitReader(file, maxFileBytes+1))
			if err != nil {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: read failed", i+1))
				return
			}
			if int64(len(fileData)) > maxFileBytes {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: exceeds %dMB limit", i+1, h.maxUploadMB))
				return
			}

			// Validate image
			_, width, height, err := media.ValidateImage(fileData)
			if err != nil {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: %v", i+1, err))
				return
			}

			// Create photo record
			photoID := utils.GenerateShortID()
			takenAt := time.Now()

			photo := &store.Photo{
				ID:       photoID,
				BookID:   bookID,
				Filename: fileHeader.Filename,
				Width:    width,
				Height:   height,
				TakenAt:  &takenAt,
				Position: photoCount + i,
			}

			if err := h.db.CreatePhoto(r.Context(), photo); err != nil {
				uploadErrors = append(uploadErrors, fmt.Sprintf("File %d: database error", i+1))
				return
			}

			// Queue for async processing. The request context is canceled as
			// soon as this handler returns, so the background update must use
			// a detached context.
			go func(pID string, data []byte) {
				bgCtx := context.WithoutCancel(r.Context())
				resultChan := h.processor.Enqueue(bookID, pID, data)
				result := <-resultChan

				// Update photo with extracted EXIF data if available
				if result.TakenAt.Unix() > 0 {
					if err := h.db.UpdatePhotoTakenAt(bgCtx, pID, result.TakenAt); err != nil {
						log.Printf("failed to update EXIF taken_at for photo %s: %v", pID, err)
					}
				}
			}(photoID, fileData)

			uploadedPhotos = append(uploadedPhotos, photoToResponse(photo))
		}()
	}

	if len(uploadedPhotos) == 0 {
		WriteError(w, http.StatusBadRequest, "No valid photos uploaded")
		return
	}

	// Response
	response := map[string]interface{}{
		"photos": uploadedPhotos,
	}
	if len(uploadErrors) > 0 {
		response["errors"] = uploadErrors
	}

	WriteSuccess(w, http.StatusCreated, response)
}

type UpdatePhotoRequest struct {
	Caption *string `json:"caption"`
}

// UpdatePhoto handles PATCH /api/photos/:id
func (h *PhotoHandler) UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	photoID := chi.URLParam(r, "id")

	photo, err := h.db.GetPhotoByID(r.Context(), photoID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if photo == nil {
		WriteError(w, http.StatusNotFound, "Photo not found")
		return
	}

	// Check authorization via book
	book, err := h.db.GetBookByID(r.Context(), photo.BookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req UpdatePhotoRequest
	LimitBody(w, r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Caption != nil {
		photo.Caption = req.Caption
	}

	if err := h.db.UpdatePhotoCaption(r.Context(), photoID, photo.Caption); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update photo")
		return
	}

	WriteSuccess(w, http.StatusOK, photoToResponse(photo))
}

// DeletePhoto handles DELETE /api/photos/:id
func (h *PhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	photoID := chi.URLParam(r, "id")

	photo, err := h.db.GetPhotoByID(r.Context(), photoID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if photo == nil {
		WriteError(w, http.StatusNotFound, "Photo not found")
		return
	}

	// Check authorization via book
	book, err := h.db.GetBookByID(r.Context(), photo.BookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Delete file
	if err := media.DeletePhotoFiles(h.dataDir, photo.BookID, photoID); err != nil {
		// Log but don't fail
		fmt.Printf("Failed to delete photo files: %v\n", err)
	}

	// Delete from database
	if err := h.db.DeletePhoto(r.Context(), photoID); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete photo")
		return
	}

	WriteSuccess(w, http.StatusOK, map[string]string{"message": "Photo deleted"})
}

type ReorderRequest struct {
	Order []string `json:"order"`
}

// ReorderPhotos handles PATCH /api/books/:id/reorder
func (h *PhotoHandler) ReorderPhotos(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	bookID := chi.URLParam(r, "id")

	// Check authorization
	book, err := h.db.GetBookByID(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req ReorderRequest
	LimitBody(w, r)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Build update map
	updates := make(map[string]int)
	for pos, photoID := range req.Order {
		updates[photoID] = pos
	}

	if err := h.db.ReorderPhotos(r.Context(), bookID, updates); err != nil {
		if errors.Is(err, store.ErrPhotoNotInBook) {
			WriteError(w, http.StatusBadRequest, "Invalid photo order")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to reorder photos")
		return
	}

	WriteSuccess(w, http.StatusOK, map[string]string{"message": "Photos reordered"})
}

type PagedPhotosResponse struct {
	Photos []PhotoResponse `json:"photos"`
	Count  int             `json:"count"`
}

// GetBookPages handles GET /api/books/:id/pages
func (h *PhotoHandler) GetBookPages(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*store.User)
	bookID := chi.URLParam(r, "id")

	// Check authorization
	book, err := h.db.GetBookByID(r.Context(), bookID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if book == nil {
		WriteError(w, http.StatusNotFound, "Book not found")
		return
	}
	if book.OwnerID != user.ID && user.Role != "admin" {
		WriteError(w, http.StatusForbidden, "Access denied")
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
