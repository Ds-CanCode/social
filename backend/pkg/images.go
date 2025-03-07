package funcs

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// SaveImage handles image file validation and saving
func SaveImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Max file size (20MB)
	const maxFileSize = 20 << 20 // 20MB

	// Check file size
	if handler.Size > maxFileSize {
		return "", ErrMaxSizeImage
	}

	// Read first 512 bytes to determine MIME type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}

	// Reset file pointer after reading
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Get the MIME type of the file
	mimeType := http.DetectContentType(buffer)

	// Allowed MIME types
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	// Check if MIME type is allowed
	if !allowedTypes[mimeType] {
		return "", ErrInvalidFile
	}

	// Generate a new name for the file
	newName, err := GenereteTocken() // Ensure this function exists
	if err != nil {
		return "", err
	}
	newName += getFileExtension(handler.Filename)

	// Ensure directory exists
	savePath := "images/"
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		return "", err
	}

	// Create destination file
	dstPath := savePath + newName
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return newName, nil
}

// getFileExtension extracts the file extension and ensures it is valid
func getFileExtension(filename string) string {
	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	if ext == "" || len(ext) > 5 { // Sanity check for extensions like `.jpeg`
		return ".jpg"
	}
	return ext
}
