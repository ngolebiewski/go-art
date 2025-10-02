package utils

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"log"

	// Register image formats without direct usage
	_ "image/gif"
	_ "image/png"

	"golang.org/x/image/draw"
)

// Define target sizes and limits
const (
	MaxThumbSize = 200        // 200x200 max size
	MaxImageSize = 400        // 400x400 max size, was 2000 pixels max on a side!
	MaxThumbBlob = 64 * 1024  // 64 KB
	MaxImageBlob = 200 * 1024 // 200 KB, Was 500KB
)

// ProcessImage generates a single JPEG BLOB for both thumbnail and full image.
func ProcessImage(file io.Reader) (thumbData []byte, imageData []byte, err error) {
	// 1. Decode the image (supports JPEG, PNG, GIF via standard library)
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Successfully decoded uploaded image (Format: %s)", format)

	// --- 2. Generate THUMBNAIL (Max 200x200, Max 64KB, JPEG) ---

	// Calculate target size (Casting int to uint for the function)
	w, h := getScaledDimensions(uint(img.Bounds().Dx()), uint(img.Bounds().Dy()), MaxThumbSize)

	// Create new destination image for the thumbnail
	thumb := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	// Scale the source image (img) to the destination image (thumb)
	draw.BiLinear.Scale(thumb, thumb.Bounds(), img, img.Bounds(), draw.Src, nil)

	// Encode and compress the thumbnail until it meets the 64KB limit
	thumbData, err = encodeAndCompressJPEG(thumb, MaxThumbBlob)
	if err != nil {
		log.Printf("Warning: Thumbnail compression failed to reach %dKB limit", MaxThumbBlob/1024)
		// Continue with best effort
	}

	// --- 3. Generate FULL-SIZE Image (Max 2000x2000, Max 500KB, JPEG) ---

	// Calculate target size (Casting int to uint for the function)
	w, h = getScaledDimensions(uint(img.Bounds().Dx()), uint(img.Bounds().Dy()), MaxImageSize)

	// Create new destination image for the full size
	fullImage := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	// Scale the source image (img) to the destination image (fullImage)
	draw.BiLinear.Scale(fullImage, fullImage.Bounds(), img, img.Bounds(), draw.Src, nil)

	// Encode and compress the full image until it meets the 500KB limit
	imageData, err = encodeAndCompressJPEG(fullImage, MaxImageBlob)
	if err != nil {
		log.Printf("Warning: Full image compression failed to reach %dKB limit", MaxImageBlob/1024)
	}

	return thumbData, imageData, nil
}

// getScaledDimensions calculates new dimensions to fit within maxDim while maintaining aspect ratio
func getScaledDimensions(width, height, maxDim uint) (uint, uint) {
	if width <= maxDim && height <= maxDim {
		return width, height
	}

	ratio := float64(width) / float64(height)
	if width > height {
		return maxDim, uint(float64(maxDim) / ratio)
	}
	return uint(float64(maxDim) * ratio), maxDim
}

// encodeAndCompressJPEG repeatedly encodes the image with decreasing quality until size limit is met
func encodeAndCompressJPEG(img image.Image, maxSize int) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Start at 90% quality and loop down
	for quality := 90; quality >= 40; quality -= 10 {
		buf.Reset()
		err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}

		if buf.Len() <= maxSize {
			return buf.Bytes(), nil
		}
	}

	// If it's still too big at 40% quality, return the best effort and an error
	return buf.Bytes(), errors.New("image size limit exceeded even after max compression")
}
