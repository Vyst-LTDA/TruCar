package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileStorageService interface {
	Save(file *multipart.FileHeader, subpath string) (string, error)
	Delete(filePath string) error
}

type localStorageService struct {
	basePath string
}

func NewLocalStorageService(basePath string) FileStorageService {
	return &localStorageService{basePath: basePath}
}

func (s *localStorageService) Save(file *multipart.FileHeader, subpath string) (string, error) {
	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create the full path
	dir := filepath.Join(s.basePath, subpath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	fullPath := filepath.Join(dir, uniqueFilename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return the public URL
	publicURL := filepath.Join("/static", subpath, uniqueFilename)
	return publicURL, nil
}

func (s *localStorageService) Delete(filePath string) error {
	// Construct the full path from the web path
	// This is a simplification and might need to be more robust
	fullPath := filepath.Join(s.basePath, filepath.Clean(filePath))

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File already deleted
	}

	return os.Remove(fullPath)
}
