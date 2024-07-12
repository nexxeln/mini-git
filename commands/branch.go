package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexxeln/mini-git/repository"
)

func Branch(startPath string, args []string) error {
	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	if len(args) == 0 {
		// List branches
		return listBranches(repoRoot)
	} else if len(args) == 1 {
		// Create new branch
		return createBranch(repoRoot, args[0])
	}

	return fmt.Errorf("invalid number of arguments for branch command")
}

func listBranches(repoRoot string) error {
	branchesDir := filepath.Join(repoRoot, ".mini-git", "refs", "heads")
	entries, err := os.ReadDir(branchesDir)
	if err != nil {
		return fmt.Errorf("failed to read branches directory: %v", err)
	}

	currentBranch, err := getCurrentBranch(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		branchName := entry.Name()
		if branchName == currentBranch {
			fmt.Printf("* %s\n", branchName)
		} else {
			fmt.Printf("  %s\n", branchName)
		}
	}

	return nil
}

func createBranch(repoRoot, branchName string) error {
	headCommitHash, err := getHEADCommitHash(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get HEAD commit hash: %v", err)
	}

	// If there are no commits yet, headCommitHash will be empty
	// In this case, we'll create an empty branch
	branchPath := filepath.Join(repoRoot, ".mini-git", "refs", "heads", branchName)
	if err := os.WriteFile(branchPath, []byte(headCommitHash), 0644); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	fmt.Printf("Created branch '%s'\n", branchName)
	return nil
}

func getHEADCommitHash(repoRoot string) (string, error) {
	headPath := filepath.Join(repoRoot, ".mini-git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD file: %v", err)
	}

	ref := strings.TrimSpace(string(headContent))
	if strings.HasPrefix(ref, "ref: ") {
		refPath := filepath.Join(repoRoot, ".mini-git", strings.TrimPrefix(ref, "ref: "))
		commitHash, err := os.ReadFile(refPath)
		if err != nil {
			if os.IsNotExist(err) {
				// If the ref file doesn't exist, it means there are no commits yet
				return "", nil
			}
			return "", fmt.Errorf("failed to read ref file: %v", err)
		}
		return strings.TrimSpace(string(commitHash)), nil
	}

	// If HEAD is not a ref, it's probably a commit hash (detached HEAD state)
	return ref, nil
}
