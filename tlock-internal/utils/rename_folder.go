package utils

import (
	"os"
	"path"
)

func RenameFolder(folder string, name string) {
	// Get directory
	parent := path.Dir(folder)

	// Rename
	os.Rename(folder, path.Join(parent, name))
}
