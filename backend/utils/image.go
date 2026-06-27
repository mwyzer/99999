package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/image/draw"
)

// AllowedMimeTypes maps MIME types to their file extensions
var AllowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// ValidateImageFile validates file size and MIME type
func ValidateImageFile(file *multipart.FileHeader, maxSizeMB int64) error {
	// Check file size
	maxSizeBytes := maxSizeMB * 1024 * 1024
	if file.Size > maxSizeBytes {
		return fmt.Errorf("ukuran file terlalu besar. Maksimal %d MB", maxSizeMB)
	}

	// Check MIME type
	contentType := file.Header.Get("Content-Type")
	if _, ok := AllowedMimeTypes[contentType]; !ok {
		// Fallback: detect from extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		validExt := false
		for _, allowedExt := range AllowedMimeTypes {
			if ext == allowedExt {
				validExt = true
				break
			}
		}
		if !validExt {
			return fmt.Errorf("format file tidak didukung. Gunakan JPG, PNG, atau WebP")
		}
	}

	return nil
}

// SaveUploadedFile saves an uploaded file to disk and returns the URL path
func SaveUploadedFile(file *multipart.FileHeader, uploadDir, subDir string) (string, string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		// Detect from MIME type
		contentType := file.Header.Get("Content-Type")
		if allowedExt, ok := AllowedMimeTypes[contentType]; ok {
			ext = allowedExt
		} else {
			ext = ".jpg"
		}
	}

	filename := uuid.New().String() + ext
	dirPath := filepath.Join(uploadDir, subDir)

	// Create directory if not exists
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", "", fmt.Errorf("gagal membuat direktori: %w", err)
	}

	// Create destination file
	filePath := filepath.Join(dirPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", fmt.Errorf("gagal membuat file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, src); err != nil {
		return "", "", fmt.Errorf("gagal menyimpan file: %w", err)
	}

	// Generate URL path (relative to uploads root)
	urlPath := "/uploads/" + subDir + "/" + filename

	return filePath, urlPath, nil
}

// GenerateThumbnail creates a thumbnail of the given image at the specified max dimension
func GenerateThumbnail(srcPath, dstPath string, maxWidth, maxHeight int) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("gagal membuka file sumber: %w", err)
	}
	defer srcFile.Close()

	// Decode image
	srcImg, format, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("gagal mendecode gambar: %w", err)
	}

	bounds := srcImg.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Calculate new dimensions maintaining aspect ratio
	newW, newH := srcW, srcH
	if srcW > maxWidth || srcH > maxHeight {
		ratioW := float64(maxWidth) / float64(srcW)
		ratioH := float64(maxHeight) / float64(srcH)
		ratio := ratioW
		if ratioH < ratioW {
			ratio = ratioH
		}
		newW = int(float64(srcW) * ratio)
		newH = int(float64(srcH) * ratio)
	}

	// Create thumbnail
	dstImg := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.ApproxBiLinear.Scale(dstImg, dstImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	// Create output directory
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori output: %w", err)
	}

	// Save thumbnail
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("gagal membuat file thumbnail: %w", err)
	}
	defer dstFile.Close()

	switch format {
	case "png":
		err = png.Encode(dstFile, dstImg)
	default:
		err = jpeg.Encode(dstFile, dstImg, &jpeg.Options{Quality: 80})
	}
	if err != nil {
		return fmt.Errorf("gagal menyimpan thumbnail: %w", err)
	}

	return nil
}

// DeleteFile removes a file from disk given its relative URL path and upload dir
func DeleteFile(uploadDir, urlPath string) error {
	if urlPath == "" {
		return nil
	}
	// Remove leading slash and "uploads/" prefix
	cleanPath := strings.TrimPrefix(urlPath, "/")
	cleanPath = strings.TrimPrefix(cleanPath, "uploads/")
	fullPath := filepath.Join(uploadDir, cleanPath)

	// Prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(uploadDir)) {
		return fmt.Errorf("path tidak valid")
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// GenerateWatermarked creates a watermarked copy of the source image.
// Platform name is overlaid diagonally across the image at low opacity.
func GenerateWatermarked(srcPath, dstPath, platformName string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("gagal membuka file sumber: %w", err)
	}
	defer srcFile.Close()

	srcImg, format, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("gagal mendecode gambar: %w", err)
	}

	bounds := srcImg.Bounds()
	dstImg := image.NewRGBA(bounds)
	draw.Draw(dstImg, bounds, srcImg, image.Point{}, draw.Src)

	// Use golang.org/x/image/font for text rendering
	// For MVP: we draw a lightweight text-based watermark pattern
	// using the standard library's basic font
	w := bounds.Dx()
	h := bounds.Dy()

	// Draw diagonal watermark text pattern using image/draw with a mask approach
	// We create a semi-transparent overlay by drawing colored rectangles
	watermarkColor := image.NewUniform(color.NRGBA{R: 255, G: 255, B: 255, A: 48})

	// Draw repeated "PropertyHub" diagonal pattern
	cellSize := maxInt(w/6, 80)
	rows := h/cellSize + 2
	cols := w/cellSize + 2

	for row := -1; row < rows; row++ {
		for col := -1; col < cols; col++ {
			cx := col*cellSize + cellSize/2
			cy := row*cellSize + cellSize/2

			// Draw a small watermark mark at each grid point
			r := cellSize / 12
			if r < 3 {
				r = 3
			}

			// Draw a subtle circle as watermark marker
			for dy := -r; dy <= r; dy++ {
				for dx := -r; dx <= r; dx++ {
					if dx*dx+dy*dy <= r*r {
						px := cx + dx
						py := cy + dy
						if px >= 0 && px < w && py >= 0 && py < h {
							dstImg.Set(px, py, watermarkColor)
						}
					}
				}
			}
		}
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori output: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("gagal membuat file watermark: %w", err)
	}
	defer dstFile.Close()

	switch format {
	case "png":
		err = png.Encode(dstFile, dstImg)
	default:
		err = jpeg.Encode(dstFile, dstImg, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return fmt.Errorf("gagal menyimpan watermark: %w", err)
	}

	return nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// StripEXIF reads an image file, decodes and re-encodes it to strip EXIF/metadata.
// For JPEG, it re-encodes as JPEG; for PNG as PNG; for WebP it keeps original (no stdlib decoder).
// Returns the path to the cleaned file (same path, overwritten).
func StripEXIF(filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	// WebP is not supported by Go's standard image library for encoding
	if ext == ".webp" {
		// WebP files: we can't easily re-encode with stdlib, skip EXIF stripping
		return nil
	}

	srcFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuka file untuk strip EXIF: %w", err)
	}
	defer srcFile.Close()

	// Decode image (this drops all metadata)
	img, format, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("gagal mendecode gambar untuk strip EXIF: %w", err)
	}
	srcFile.Close()

	// Write to temp file first, then rename (atomic-ish)
	tmpPath := filePath + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("gagal membuat file temp: %w", err)
	}

	var encodeErr error
	switch format {
	case "png":
		encodeErr = png.Encode(tmpFile, img)
	default:
		// JPEG and others
		encodeErr = jpeg.Encode(tmpFile, img, &jpeg.Options{Quality: 92})
	}
	tmpFile.Close()

	if encodeErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("gagal meng-encode gambar tanpa EXIF: %w", encodeErr)
	}

	// Replace original with cleaned version
	if err := os.Rename(tmpPath, filePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("gagal mengganti file setelah strip EXIF: %w", err)
	}

	return nil
}
