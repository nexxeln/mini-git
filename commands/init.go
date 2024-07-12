package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func Init(path string) error {
	gitDir := filepath.Join(path, ".mini-git")

	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return fmt.Errorf("failed to create .mini-git directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(gitDir, "objects"), 0755); err != nil {
		return fmt.Errorf("failed to create objects directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(gitDir, "refs", "heads"), 0755); err != nil {
		return fmt.Errorf("failed to create refs directory: %v", err)
	}

	headPath := filepath.Join(gitDir, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	masterRef := filepath.Join(gitDir, "refs", "heads", "master")
	if err := os.WriteFile(masterRef, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create master branch: %v", err)
	}

	configPath := filepath.Join(gitDir, "config")
	configContent := `[core]
	repositoryformatversion = 0
	filemode = false
	bare = false
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	fmt.Println("Initialized empty Mini Git repository in", gitDir)
	return nil
}
