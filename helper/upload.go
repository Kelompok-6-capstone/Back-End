package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadFile handles the file upload process
func UploadFile(fileHeader *multipart.FileHeader, uploadDir string, filePrefix string, maxSize int64, allowedExtensions []string) (string, error) {
	// Validate file size
	if fileHeader.Size > maxSize {
		return "", fmt.Errorf("ukuran file maksimal %d MB", maxSize/(1024*1024))
	}

	// Validate file extension
	ext := filepath.Ext(fileHeader.Filename)
	allowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("file harus memiliki ekstensi %v", allowedExtensions)
	}

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("gagal membuat direktori upload: %v", err)
		}
	}

	// Generate unique file name
	filePath := fmt.Sprintf("%s/%s_%s", uploadDir, filePrefix, fileHeader.Filename)
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}
	defer dst.Close()

	// Copy file to the destination
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	// Return the file path
	return filePath, nil
}

// DeleteFile deletes a file from the given path
func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("gagal menghapus file: %v", err)
		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("file tidak ditemukan")
	} else {
		return fmt.Errorf("gagal memeriksa file: %v", err)
	}
	return nil
}
