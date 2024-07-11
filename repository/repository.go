package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

type Repository struct {
	WorkTree string
	GitDir   string
}

func InitRepository(path string) (*Repository, error) {
	repo := &Repository{
		WorkTree: path,
		GitDir:   filepath.Join(path, ".mini-git"),
	}

	if err := os.MkdirAll(repo.GitDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .mini-git directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(repo.GitDir, "objects"), 0755); err != nil {
		return nil, fmt.Errorf("failed to create objects directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(repo.GitDir, "refs"), 0755); err != nil {
		return nil, fmt.Errorf("failed to create refs directory: %v", err)
	}

	headFile := filepath.Join(repo.GitDir, "HEAD")
	if err := os.WriteFile(headFile, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		return nil, fmt.Errorf("failed to create HEAD file: %v", err)
	}

	return repo, nil
}
