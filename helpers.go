package main

import "os"

func CheckDir(dir string) ([]os.DirEntry, error) {
	// If dir/ does not exist, create it
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// Return any files
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	return entries, err
}
