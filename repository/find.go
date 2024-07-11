package repository

import (
	"os"
	"path/filepath"
)

func FindRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(current, ".mini-git")
		if _, err := os.Stat(gitDir); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", os.ErrNotExist
		}
		current = parent
	}
}
