package objects

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nexxeln/mini-git/blob"
	"github.com/nexxeln/mini-git/commit"
	"github.com/nexxeln/mini-git/tree"
)

func Store(repoPath string, obj interface{}) error {
	var data []byte
	var hash string
	var err error

	switch o := obj.(type) {
	case *blob.Blob:
		data, err = o.Serialize()
		hash = o.Hash
	case *tree.Tree:
		data, err = o.Serialize()
		hash = o.Hash()
	case *commit.Commit:
		data, err = o.Serialize()
		hash = o.Hash()
	default:
		return fmt.Errorf("unsupported object type")
	}

	if err != nil {
		return fmt.Errorf("failed to serialize object: %v", err)
	}

	objectsDir := filepath.Join(repoPath, ".mini-git", "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create objects directory: %v", err)
	}

	objectPath := filepath.Join(objectsDir, hash[:2], hash[2:])
	objectDir := filepath.Dir(objectPath)
	if err := os.MkdirAll(objectDir, 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %v", err)
	}

	if err := os.WriteFile(objectPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write object to file: %v", err)
	}

	return nil
}

func RetrieveBlob(repoPath, hash string) (*blob.Blob, error) {
	data, err := retrieveObject(repoPath, hash)
	if err != nil {
		return nil, err
	}

	return blob.Deserialize(data)
}

func RetrieveTree(repoPath, hash string) (*tree.Tree, error) {
	data, err := retrieveObject(repoPath, hash)
	if err != nil {
		return nil, err
	}

	return tree.Deserialize(data)
}

func RetrieveCommit(repoPath, hash string) (*commit.Commit, error) {
	data, err := retrieveObject(repoPath, hash)
	if err != nil {
		return nil, err
	}

	return commit.Deserialize(data)
}

func retrieveObject(repoPath, hash string) ([]byte, error) {
	objectPath := filepath.Join(repoPath, ".mini-git", "objects", hash[:2], hash[2:])

	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read object file: %v", err)
	}

	return data, nil
}
