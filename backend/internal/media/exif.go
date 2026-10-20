package media

import (
	"bytes"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// ExtractEXIFData extracts EXIF metadata from image
type EXIFData struct {
	TakenAt     time.Time
	Orientation int
	Width       int
	Height      int
}

// ExtractEXIF extracts EXIF metadata from image bytes
func ExtractEXIF(data []byte) *EXIFData {
	exifData := &EXIFData{
		TakenAt:     time.Now(),
		Orientation: 1,
	}

	ex, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		// No EXIF data, return defaults
		return exifData
	}

	// Extract taken date
	if dt, err := ex.DateTime(); err == nil {
		exifData.TakenAt = dt
	}

	// Extract orientation
	if orientation, err := ex.Get(exif.Orientation); err == nil {
		if val, err := orientation.Int(0); err == nil {
			exifData.Orientation = int(val)
		}
	}

	// Extract width/height (ImageWidth/ImageLength tags)
	if w, err := ex.Get(exif.ImageWidth); err == nil {
		if val, err := w.Int(0); err == nil {
			exifData.Width = int(val)
		}
	}

	if h, err := ex.Get(exif.ImageLength); err == nil {
		if val, err := h.Int(0); err == nil {
			exifData.Height = int(val)
		}
	}

	return exifData
}

// Try to extract orientation safely
func getOrientationSafe(ex *exif.Exif) int {
	if ex == nil {
		return 1
	}

	tag, err := ex.Get(exif.Orientation)
	if err != nil {
		return 1
	}

	val, err := tag.Int(0)
	if err != nil {
		return 1
	}

	return int(val)
}
