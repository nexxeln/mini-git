package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nexxeln/mini-git/blob"
	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
)

func Add(startPath, filePath string) error {
	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	b, err := blob.NewBlob(content)
	if err != nil {
		return fmt.Errorf("failed to create blob: %v", err)
	}

	if err := objects.Store(repoRoot, b); err != nil {
		return fmt.Errorf("failed to store blob: %v", err)
	}

	indexPath := filepath.Join(repoRoot, ".mini-git", "index")
	f, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open index file: %v", err)
	}
	defer f.Close()

	relPath, err := filepath.Rel(repoRoot, absFilePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %v", err)
	}

	_, err = fmt.Fprintf(f, "%s %s\n", b.Hash, relPath)
	if err != nil {
		return fmt.Errorf("failed to write to index file: %v", err)
	}

	return nil
}
