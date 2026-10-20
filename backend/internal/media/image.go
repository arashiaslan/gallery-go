package media

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/disintegration/imaging"
)

// ImageFormat represents supported image formats
type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
	FormatWebP ImageFormat = "webp"
	FormatHEIC ImageFormat = "heic" // Not supported without CGO
)

// ValidateImage checks if file is valid image and returns format + dimensions
func ValidateImage(data []byte) (ImageFormat, int, int, error) {
	if len(data) == 0 {
		return "", 0, 0, errors.New("empty file")
	}

	// Sniff content type
	contentType := http.DetectContentType(data)

	// Try to decode image
	img, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, errors.New("invalid image file")
	}

	// Determine format based on content-type
	var imageFormat ImageFormat
	switch {
	case contains(contentType, "jpeg"):
		imageFormat = FormatJPEG
	case contains(contentType, "png"):
		imageFormat = FormatPNG
	case contains(contentType, "webp"):
		imageFormat = FormatJPEG // Convert WebP to JPEG for now
	case contains(contentType, "heic"):
		return "", 0, 0, errors.New("HEIC format requires CGO, not supported")
	default:
		// Fallback to format string from decode
		switch format {
		case "jpeg":
			imageFormat = FormatJPEG
		case "png":
			imageFormat = FormatPNG
		case "webp":
			imageFormat = FormatJPEG // Convert WebP to JPEG for now
		default:
			return "", 0, 0, errors.New("unsupported image format: " + format)
		}
	}

	return imageFormat, img.Width, img.Height, nil
}

// DecodeImage decodes image data
func DecodeImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// EncodePNG encodes image to PNG format
func EncodePNG(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeJPEG encodes image to JPEG format
func EncodeJPEG(img image.Image, quality int) ([]byte, error) {
	buf := &bytes.Buffer{}
	options := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(buf, img, options); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ResizeImage resizes image to fit within maxWidth x maxHeight, maintaining aspect ratio
func ResizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	return imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
}

// CorrectOrientation applies EXIF orientation to image
func CorrectOrientation(img image.Image, orientation int) image.Image {
	// EXIF orientation values
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Rotate270(imaging.FlipH(img))
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Rotate90(imaging.FlipH(img))
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

// GenerateThumbnail creates a thumbnail (~400px)
func GenerateThumbnail(img image.Image) ([]byte, error) {
	resized := ResizeImage(img, 400, 400)
	return EncodeJPEG(resized, 80)
}

// GenerateMedium creates a medium image (~1200px)
func GenerateMedium(img image.Image) ([]byte, error) {
	resized := ResizeImage(img, 1200, 1200)
	return EncodeJPEG(resized, 85)
}

// Contains checks if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return bytes.Contains(bytes.ToLower([]byte(s)), bytes.ToLower([]byte(substr)))
}

// GetFileExtension returns file extension for format
func GetFileExtension(format ImageFormat) string {
	switch format {
	case FormatJPEG:
		return ".jpg"
	case FormatPNG:
		return ".png"
	case FormatWebP:
		return ".webp"
	default:
		return ".webp"
	}
}

// GetMimeType returns MIME type for format
func GetMimeType(format ImageFormat) string {
	switch format {
	case FormatJPEG:
		return "image/jpeg"
	case FormatPNG:
		return "image/png"
	case FormatWebP:
		return "image/webp"
	default:
		return "image/webp"
	}
}
