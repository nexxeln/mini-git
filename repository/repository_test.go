package repository

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mini-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo, err := InitRepository(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %v", err)
	}

	miniGitDir := filepath.Join(tempDir, ".mini-git")
	if _, err := os.Stat(miniGitDir); os.IsNotExist(err) {
		t.Errorf(".mini-git directory was not created")
	}

	objectsDir := filepath.Join(miniGitDir, "objects")
	if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
		t.Errorf("objects directory was not created")
	}

	refsDir := filepath.Join(miniGitDir, "refs")
	if _, err := os.Stat(refsDir); os.IsNotExist(err) {
		t.Errorf("refs directory was not created")
	}

	headFile := filepath.Join(miniGitDir, "HEAD")
	if _, err := os.Stat(headFile); os.IsNotExist(err) {
		t.Errorf("HEAD file was not created")
	}

	if repo.WorkTree != tempDir {
		t.Errorf("WorkTree not set correctly. Expected %s, got %s", tempDir, repo.WorkTree)
	}
	if repo.GitDir != miniGitDir {
		t.Errorf("GitDir not set correctly. Expected %s, got %s", miniGitDir, repo.GitDir)
	}
}
