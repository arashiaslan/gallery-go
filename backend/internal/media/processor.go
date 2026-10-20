package media

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/squarebrain/gallery/internal/store"
)

// ProcessResult holds result of image processing
type ProcessResult struct {
	PhotoID     string
	Width       int
	Height      int
	TakenAt     time.Time
	OriginalErr error
	ThumbErr    error
	MediumErr   error
}

// Processor handles concurrent image processing
type Processor struct {
	dataDir    string
	maxWorkers int
	queue      chan *ProcessJob
	mu         sync.RWMutex
	closed     bool
	wg         sync.WaitGroup
	db         *store.DB
}

// ProcessJob represents a single image to process
type ProcessJob struct {
	BookID      string
	PhotoID     string
	OriginalData []byte
	ResultChan  chan *ProcessResult
}

// NewProcessor creates new image processor with worker pool
func NewProcessor(dataDir string, maxWorkers int, db *store.DB) *Processor {
	p := &Processor{
		dataDir:    dataDir,
		maxWorkers: maxWorkers,
		queue:      make(chan *ProcessJob, 100),
		db:         db,
	}

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	return p
}

// Enqueue adds job to processing queue. If the processor is shut down it
// returns a result carrying an error instead of blocking or panicking.
// The read lock guarantees the queue is never closed while a send is in
// flight, so no "send on closed channel" panic can happen.
func (p *Processor) Enqueue(bookID, photoID string, data []byte) <-chan *ProcessResult {
	resultChan := make(chan *ProcessResult, 1)

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		resultChan <- &ProcessResult{
			PhotoID:     photoID,
			OriginalErr: fmt.Errorf("processor is shut down"),
		}
		return resultChan
	}

	p.queue <- &ProcessJob{
		BookID:       bookID,
		PhotoID:      photoID,
		OriginalData: data,
		ResultChan:   resultChan,
	}
	return resultChan
}

// Shutdown gracefully stops the processor: it stops new jobs, closes the
// queue safely (no senders can be active) and waits for workers to drain.
func (p *Processor) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	if !p.closed {
		p.closed = true
		close(p.queue)
	}
	p.mu.Unlock()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// worker processes jobs until the queue is closed and drained
func (p *Processor) worker() {
	defer p.wg.Done()

	for job := range p.queue {
		job.ResultChan <- p.processImage(job)
	}
}

// processImage processes a single image
func (p *Processor) processImage(job *ProcessJob) *ProcessResult {
	result := &ProcessResult{
		PhotoID: job.PhotoID,
		TakenAt: time.Now(),
	}

	// Validate image
	format, width, height, err := ValidateImage(job.OriginalData)
	if err != nil {
		result.OriginalErr = fmt.Errorf("validation failed: %w", err)
		return result
	}

	result.Width = width
	result.Height = height

	// Extract EXIF data
	exifData := ExtractEXIF(job.OriginalData)
	if !exifData.TakenAt.IsZero() {
		result.TakenAt = exifData.TakenAt
	}

	// Decode image
	img, err := DecodeImage(job.OriginalData)
	if err != nil {
		result.OriginalErr = fmt.Errorf("decode failed: %w", err)
		return result
	}

	// Apply EXIF orientation correction
	if exifData.Orientation > 1 {
		img = CorrectOrientation(img, exifData.Orientation)
	}

	// Create upload directory structure
	photoDir := filepath.Join(p.dataDir, "uploads", job.BookID, job.PhotoID)
	if err := os.MkdirAll(photoDir, 0755); err != nil {
		result.OriginalErr = fmt.Errorf("mkdir failed: %w", err)
		return result
	}

	// Save original image
	originalPath := filepath.Join(photoDir, "original"+GetFileExtension(format))
	if err := os.WriteFile(originalPath, job.OriginalData, 0644); err != nil {
		result.OriginalErr = fmt.Errorf("save original failed: %w", err)
		return result
	}

	// Generate and save thumbnail
	thumbData, err := GenerateThumbnail(img)
	if err != nil {
		result.ThumbErr = fmt.Errorf("generate thumbnail failed: %w", err)
	} else {
		thumbPath := filepath.Join(photoDir, "thumb.jpg")
		if err := os.WriteFile(thumbPath, thumbData, 0644); err != nil {
			result.ThumbErr = fmt.Errorf("save thumbnail failed: %w", err)
		}
	}

	// Generate and save medium image
	mediumData, err := GenerateMedium(img)
	if err != nil {
		result.MediumErr = fmt.Errorf("generate medium failed: %w", err)
	} else {
		mediumPath := filepath.Join(photoDir, "medium.jpg")
		if err := os.WriteFile(mediumPath, mediumData, 0644); err != nil {
			result.MediumErr = fmt.Errorf("save medium failed: %w", err)
		}
	}

	return result
}

// GetPhotoPath returns path to photo file
func GetPhotoPath(dataDir, bookID, photoID, size string) string {
	// Validate size to prevent path traversal
	if size != "thumb" && size != "medium" && size != "original" {
		size = "medium"
	}

	// Validate IDs to prevent path traversal
	if !isValidID(bookID) || !isValidID(photoID) {
		return ""
	}

	var filename string
	switch size {
	case "thumb":
		filename = "thumb.jpg"
	case "medium":
		filename = "medium.jpg"
	case "original":
		filename = "original.jpg" // Default, actual extension may vary
	}

	return filepath.Join(dataDir, "uploads", bookID, photoID, filename)
}

// isValidID checks if ID contains only safe characters
func isValidID(id string) bool {
	if len(id) == 0 || len(id) > 100 {
		return false
	}
	for _, ch := range id {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return false
		}
	}
	return true
}

// PhotoFilePath builds path for specific photo file (use for deletion)
func PhotoFilePath(dataDir, bookID, photoID, size string) string {
	return filepath.Join(dataDir, "uploads", bookID, photoID, size)
}

// DeletePhotoFiles removes all photo files and directory
func DeletePhotoFiles(dataDir, bookID, photoID string) error {
	photoDir := filepath.Join(dataDir, "uploads", bookID, photoID)
	return os.RemoveAll(photoDir)
}

// RemoveBookFiles removes every file belonging to a book
func RemoveBookFiles(dataDir, bookID string) error {
	bookDir := filepath.Join(dataDir, "uploads", bookID)
	return os.RemoveAll(bookDir)
}
