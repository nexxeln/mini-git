package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nexxeln/mini-git/blob"
	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
)

func Add(fileNode []string) error {

	if len(fileNode) < 1 {
		return errors.New("Usage: mini-git add <file> [<file> ...]")
	}

	for _, node := range fileNode {
		path_info, err := os.Stat(node)
		if err != nil {
			return fmt.Errorf("Error adding file: %v", err)
		}

		if node == "." || path_info.IsDir() {
			if err := AddDirectory(node); err != nil {
				return fmt.Errorf("Error adding directory: %v", err)
			}
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("Error getting current working directory: %v", err)
			}

			if err := AddFile(cwd, node); err != nil {
				return fmt.Errorf("Error adding file: %v", err)
			}
		}
	}

	return nil
}

func AddDirectory(startPath string) error {
	return filepath.Walk(startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path: %v", err)
		}

		if info.IsDir() {
			return nil
		}

		if err := AddFile(startPath, path); err != nil {
			return fmt.Errorf("failed to add file: %v", err)
		}

		return nil
	})
}

func AddFile(startPath, filePath string) error {

	path_info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if filePath == "." || path_info.IsDir() {
		return AddDirectory(startPath)
	}

	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		fmt.Printf("not a mini-git repository (or any of the parent directories): %v \n Would you like to initialize ? [Y] Yes or any other key to exit: ", err)

		var response string
		_, err = fmt.Scanln(&response)
		if err != nil && err.Error() != "unexpected newline" {
			return fmt.Errorf("failed to read response: %v", err)
		}

		if (response == "Y" || response == "y") && response != "" {
			if err := Init(startPath); err != nil {
				return fmt.Errorf("failed to initialize repository: %v", err)
			}

			repoRoot, err = repository.FindRoot(startPath)
			if err != nil {
				return fmt.Errorf("failed to find repository root: %v", err)
			}

		} else {
			return errors.New("not a mini-git repository")
		}
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
