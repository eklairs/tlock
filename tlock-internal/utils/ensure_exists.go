package utils

import (
	"os"
	"path/filepath"
)

// Ensures if the directory for the file exists, and then creates the file
func EnsureExists(file string) (*os.File, error) {
	// Ensure that the parent of the file exists
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return nil, err
	}

	// Create file and return
	return os.Create(file)
}
