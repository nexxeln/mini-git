package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
)

func Merge(startPath string, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mini-git merge <branch-name>")
	}

	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	branchToMerge := args[0]
	return mergeBranch(repoRoot, branchToMerge)
}

func mergeBranch(repoRoot, branchToMerge string) error {
	currentBranch, err := getCurrentBranch(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	currentCommitHash, err := getCommitHash(repoRoot, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get current commit hash: %v", err)
	}

	mergeCommitHash, err := getCommitHash(repoRoot, branchToMerge)
	if err != nil {
		return fmt.Errorf("failed to get merge commit hash: %v", err)
	}

	// Check if it's a fast-forward merge
	isAncestor, err := isAncestor(repoRoot, currentCommitHash, mergeCommitHash)
	if err != nil {
		return fmt.Errorf("failed to check ancestry: %v", err)
	}

	if isAncestor {
		return fastForwardMerge(repoRoot, currentBranch, branchToMerge, mergeCommitHash)
	}

	return fmt.Errorf("non-fast-forward merges are not yet implemented")
}

func fastForwardMerge(repoRoot, currentBranch, branchToMerge, mergeCommitHash string) error {
	branchPath := filepath.Join(repoRoot, ".mini-git", "refs", "heads", currentBranch)
	if err := os.WriteFile(branchPath, []byte(mergeCommitHash), 0644); err != nil {
		return fmt.Errorf("failed to update branch reference: %v", err)
	}

	if err := checkoutCommit(repoRoot, mergeCommitHash); err != nil {
		return fmt.Errorf("failed to update working directory: %v", err)
	}

	fmt.Printf("Fast-forward merge successful. %s merged into %s.\n", branchToMerge, currentBranch)
	return nil
}

func isAncestor(repoRoot, possibleAncestor, commit string) (bool, error) {
	for commit != "" {
		if commit == possibleAncestor {
			return true, nil
		}

		commitObj, err := objects.RetrieveCommit(repoRoot, commit)
		if err != nil {
			return false, fmt.Errorf("failed to retrieve commit: %v", err)
		}

		commit = commitObj.ParentHash
	}

	return false, nil
}

func getCommitHash(repoRoot, branchName string) (string, error) {
	branchPath := filepath.Join(repoRoot, ".mini-git", "refs", "heads", branchName)
	commitHash, err := os.ReadFile(branchPath)
	if err != nil {
		return "", fmt.Errorf("failed to read branch file: %v", err)
	}
	return strings.TrimSpace(string(commitHash)), nil
}

func checkoutCommit(repoRoot, commitHash string) error {
	commit, err := objects.RetrieveCommit(repoRoot, commitHash)
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

	return nil
}
