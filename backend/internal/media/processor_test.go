package media

import (
	"bytes"
	"image"
	"image/jpeg"
	"testing"
)

func createTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buf := &bytes.Buffer{}
	jpeg.Encode(buf, img, nil)
	return buf.Bytes()
}

func TestValidateImage(t *testing.T) {
	imageData := createTestImage(800, 600)

	format, width, height, err := ValidateImage(imageData)
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	if format != FormatJPEG {
		t.Errorf("Expected JPEG, got %s", format)
	}

	if width != 800 || height != 600 {
		t.Errorf("Expected 800x600, got %dx%d", width, height)
	}
}

func TestValidateImageInvalid(t *testing.T) {
	_, _, _, err := ValidateImage([]byte("not an image"))
	if err == nil {
		t.Error("Should fail on invalid image")
	}
}

func TestResizeImage(t *testing.T) {
	imageData := createTestImage(1600, 1200)
	img, _ := DecodeImage(imageData)

	resized := ResizeImage(img, 400, 400)
	bounds := resized.Bounds()

	// Should fit within 400x400
	if bounds.Dx() > 400 || bounds.Dy() > 400 {
		t.Errorf("Resized image too large: %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestCorrectOrientation(t *testing.T) {
	imageData := createTestImage(100, 200)
	img, _ := DecodeImage(imageData)

	// Rotate 90
	rotated := CorrectOrientation(img, 6)
	bounds := rotated.Bounds()

	// After 90 degree rotation, dimensions should swap
	if bounds.Dx() == 100 && bounds.Dy() == 200 {
		t.Error("Image should be rotated")
	}
}

func TestGenerateThumbnail(t *testing.T) {
	imageData := createTestImage(1600, 1200)
	img, _ := DecodeImage(imageData)

	thumbData, err := GenerateThumbnail(img)
	if err != nil {
		t.Fatalf("Failed to generate thumbnail: %v", err)
	}

	if len(thumbData) == 0 {
		t.Error("Thumbnail data is empty")
	}

	// Verify output is valid JPEG
	_, format, err := image.DecodeConfig(bytes.NewReader(thumbData))
	if err != nil {
		t.Errorf("Failed to decode thumbnail: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("Expected jpeg, got %s", format)
	}
}

func TestGenerateMedium(t *testing.T) {
	imageData := createTestImage(3200, 2400)
	img, _ := DecodeImage(imageData)

	mediumData, err := GenerateMedium(img)
	if err != nil {
		t.Fatalf("Failed to generate medium: %v", err)
	}

	if len(mediumData) == 0 {
		t.Error("Medium data is empty")
	}

	// Verify output is valid JPEG
	_, format, err := image.DecodeConfig(bytes.NewReader(mediumData))
	if err != nil {
		t.Errorf("Failed to decode medium: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("Expected jpeg, got %s", format)
	}
}
