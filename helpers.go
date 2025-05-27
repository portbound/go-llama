package main

import (
	"os"
	"path/filepath"
)

func CheckDir(dir string) ([]os.DirEntry, error) {
	// If dir/ does not exist, create it
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)

	// Construct path to go-llama/<subdir>
	path := filepath.Join(exeDir, dir)

	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	}

	// Return any files
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return entries, err
}
