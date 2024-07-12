package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
	"github.com/nexxeln/mini-git/tree"
)

func Checkout(startPath string, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mini-git checkout <branch-name>")
	}

	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	branchName := args[0]
	return checkoutBranch(repoRoot, branchName)
}

func checkoutBranch(repoRoot, branchName string) error {
	branchPath := filepath.Join(repoRoot, ".mini-git", "refs", "heads", branchName)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	commitHash, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("failed to read branch file: %v", err)
	}

	headPath := filepath.Join(repoRoot, ".mini-git", "HEAD")
	headContent := fmt.Sprintf("ref: refs/heads/%s", branchName)
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %v", err)
	}

	if len(strings.TrimSpace(string(commitHash))) == 0 {
		fmt.Printf("Switched to branch '%s'\n", branchName)
		return nil
	}

	commit, err := objects.RetrieveCommit(repoRoot, strings.TrimSpace(string(commitHash)))
	if err != nil {
		return fmt.Errorf("failed to retrieve commit: %v", err)
	}

	tree, err := objects.RetrieveTree(repoRoot, commit.TreeHash)
	if err != nil {
		return fmt.Errorf("failed to retrieve tree: %v", err)
	}

	if err := updateWorkingDirectory(repoRoot, tree); err != nil {
		return fmt.Errorf("failed to update working directory: %v", err)
	}

	fmt.Printf("Switched to branch '%s'\n", branchName)
	return nil
}

func updateWorkingDirectory(repoRoot string, tree *tree.Tree) error {
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".mini-git" {
			return filepath.SkipDir
		}
		if path != repoRoot {
			os.RemoveAll(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, entry := range tree.Entries {
		blob, err := objects.RetrieveBlob(repoRoot, entry.Hash)
		if err != nil {
			return fmt.Errorf("failed to retrieve blob: %v", err)
		}

		filePath := filepath.Join(repoRoot, entry.Name)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directories: %v", err)
		}

		if err := os.WriteFile(filePath, blob.Content, 0644); err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
	}

	return nil
}
