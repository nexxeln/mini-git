package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexxeln/mini-git/commit"
	"github.com/nexxeln/mini-git/objects"
	"github.com/nexxeln/mini-git/repository"
	"github.com/nexxeln/mini-git/tree"
)

func Commit(startPath, message, author string) error {
	repoRoot, err := repository.FindRoot(startPath)
	if err != nil {
		return fmt.Errorf("not a mini-git repository (or any of the parent directories): %v", err)
	}

	indexPath := filepath.Join(repoRoot, ".mini-git", "index")
	indexFile, err := os.Open(indexPath)
	if err != nil {
		return fmt.Errorf("failed to open index file: %v", err)
	}
	defer indexFile.Close()

	rootTree := tree.NewTree()
	scanner := bufio.NewScanner(indexFile)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			return fmt.Errorf("invalid index entry: %s", scanner.Text())
		}
		hash, path := parts[0], parts[1]
		rootTree.AddEntry(path, hash, tree.EntryTypeBlob)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading index file: %v", err)
	}

	if err := objects.Store(repoRoot, rootTree); err != nil {
		return fmt.Errorf("failed to store root tree: %v", err)
	}

	headPath := filepath.Join(repoRoot, ".mini-git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD file: %v", err)
	}

	currentRef := strings.TrimSpace(string(headContent))
	var parentHash string
	if strings.HasPrefix(currentRef, "ref: ") {
		refPath := filepath.Join(repoRoot, ".mini-git", strings.TrimPrefix(currentRef, "ref: "))
		refContent, err := os.ReadFile(refPath)
		if err == nil {
			parentHash = strings.TrimSpace(string(refContent))
		}
	} else {
		parentHash = currentRef
	}

	newCommit := commit.NewCommit(rootTree.Hash(), parentHash, author, author, message)
	if err := objects.Store(repoRoot, newCommit); err != nil {
		return fmt.Errorf("failed to store commit: %v", err)
	}

	if strings.HasPrefix(currentRef, "ref: ") {
		refPath := filepath.Join(repoRoot, ".mini-git", strings.TrimPrefix(currentRef, "ref: "))
		if err := os.WriteFile(refPath, []byte(newCommit.Hash()), 0644); err != nil {
			return fmt.Errorf("failed to update ref file: %v", err)
		}
	} else {
		if err := os.WriteFile(headPath, []byte(newCommit.Hash()), 0644); err != nil {
			return fmt.Errorf("failed to update HEAD: %v", err)
		}
	}

	return nil
}
