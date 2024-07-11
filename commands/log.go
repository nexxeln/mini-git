package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
)

func Log(startPath string) error {
	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	headPath := filepath.Join(repoRoot, ".mini-git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD file: %v", err)
	}

	currentRef := strings.TrimSpace(string(headContent))
	if strings.HasPrefix(currentRef, "ref: ") {
		refPath := filepath.Join(repoRoot, ".mini-git", strings.TrimPrefix(currentRef, "ref: "))
		refContent, err := os.ReadFile(refPath)
		if err != nil {
			return fmt.Errorf("failed to read ref file %s: %v", refPath, err)
		}
		currentRef = strings.TrimSpace(string(refContent))
	}

	for currentRef != "" {
		commit, err := objects.RetrieveCommit(repoRoot, currentRef)
		if err != nil {
			return fmt.Errorf("failed to retrieve commit %s: %v", currentRef, err)
		}

		fmt.Printf("commit %s\n", currentRef)
		fmt.Printf("Author: %s\n", commit.Author)
		fmt.Printf("Date: %s\n", commit.AuthorDate)
		fmt.Printf("\n    %s\n\n", commit.Message)

		currentRef = commit.ParentHash
	}

	return nil
}
