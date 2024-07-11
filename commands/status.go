package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexxeln/mini-git/blob"
	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
)

func Status(startPath string) error {
	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	branch, err := getCurrentBranch(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	fmt.Printf("On branch %s\n", branch)

	latestCommitTree, err := getLatestCommitTree(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to get latest commit tree: %v", err)
	}

	staged, err := getStagedChanges(repoRoot, latestCommitTree)
	if err != nil {
		return fmt.Errorf("failed to get staged changes: %v", err)
	}

	unstaged, err := getUnstagedChanges(repoRoot, staged, latestCommitTree)
	if err != nil {
		return fmt.Errorf("failed to get unstaged changes: %v", err)
	}

	if len(staged) > 0 {
		fmt.Println("\nChanges to be committed:")
		for _, file := range staged {
			fmt.Printf("  new file: %s\n", file)
		}
	}

	if len(unstaged) > 0 {
		fmt.Println("\nChanges not staged for commit:")
		for _, file := range unstaged {
			fmt.Printf("  modified: %s\n", file)
		}
	}

	if len(staged) == 0 && len(unstaged) == 0 {
		if len(latestCommitTree) == 0 {
			fmt.Println("No commits yet")
		} else {
			fmt.Println("nothing to commit, working tree clean")
		}
	}

	return nil
}

func getCurrentBranch(repoRoot string) (string, error) {
	headPath := filepath.Join(repoRoot, ".mini-git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD file: %v", err)
	}

	ref := strings.TrimSpace(string(headContent))
	if strings.HasPrefix(ref, "ref: refs/heads/") {
		return strings.TrimPrefix(ref, "ref: refs/heads/"), nil
	}

	return "detached HEAD", nil
}

func getLatestCommitTree(repoRoot string) (map[string]string, error) {
	headPath := filepath.Join(repoRoot, ".mini-git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read HEAD file: %v", err)
	}

	ref := strings.TrimSpace(string(headContent))
	var commitHash string
	if strings.HasPrefix(ref, "ref: ") {
		refPath := filepath.Join(repoRoot, ".mini-git", strings.TrimPrefix(ref, "ref: "))
		commitHashBytes, err := os.ReadFile(refPath)
		if err != nil {
			if os.IsNotExist(err) {
				// If the ref file doesn't exist, it means there are no commits yet
				return make(map[string]string), nil
			}
			return nil, fmt.Errorf("failed to read ref file: %v", err)
		}
		commitHash = strings.TrimSpace(string(commitHashBytes))
	} else {
		commitHash = ref
	}

	// If commitHash is empty, it means there are no commits yet
	if commitHash == "" {
		return make(map[string]string), nil
	}

	commit, err := objects.RetrieveCommit(repoRoot, commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve commit: %v", err)
	}

	tree, err := objects.RetrieveTree(repoRoot, commit.TreeHash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tree: %v", err)
	}

	treeFiles := make(map[string]string)
	for _, entry := range tree.Entries {
		treeFiles[entry.Name] = entry.Hash
	}

	return treeFiles, nil
}

func getStagedChanges(repoPath string, latestCommitTree map[string]string) ([]string, error) {
	indexPath := filepath.Join(repoPath, ".mini-git", "index")
	indexFile, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to open index file: %v", err)
	}
	defer indexFile.Close()

	var staged []string
	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 2 {
			fileName := parts[1]
			fileHash := parts[0]
			if committedHash, exists := latestCommitTree[fileName]; !exists || committedHash != fileHash {
				staged = append(staged, fileName)
			}
		}
	}

	return staged, scanner.Err()
}

func getUnstagedChanges(repoRoot string, staged []string, committedFiles map[string]string) ([]string, error) {
	stagedMap := make(map[string]bool)
	for _, file := range staged {
		stagedMap[file] = true
	}

	var unstaged []string
	err := filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".mini-git" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}

		if !stagedMap[relPath] {
			if storedHash, exists := committedFiles[relPath]; exists {
				// File exists in the commit, check if it's modified
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				newBlob, err := blob.NewBlob(content)
				if err != nil {
					return err
				}
				if newBlob.Hash != storedHash {
					unstaged = append(unstaged, relPath)
				}
			} else {
				// New file
				unstaged = append(unstaged, relPath)
			}
		}

		return nil
	})

	return unstaged, err
}
