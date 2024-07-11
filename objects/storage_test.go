package objects

import (
	"bytes"
	"os"
	"testing"

	"github.com/nexxeln/mini-git/blob"
	"github.com/nexxeln/mini-git/tree"
)

func TestStoreAndRetrieveBlob(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mini-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	content := []byte("Hello, World!")
	b, err := blob.NewBlob(content)
	if err != nil {
		t.Fatalf("Failed to create new blob: %v", err)
	}

	err = Store(tempDir, b)
	if err != nil {
		t.Fatalf("Failed to store blob: %v", err)
	}

	retrievedBlob, err := RetrieveBlob(tempDir, b.Hash)
	if err != nil {
		t.Fatalf("Failed to retrieve blob: %v", err)
	}

	if !bytes.Equal(retrievedBlob.Content, b.Content) {
		t.Errorf("Retrieved blob content does not match original")
	}

	if retrievedBlob.Size != b.Size {
		t.Errorf("Retrieved blob size does not match. Expected %d, got %d", b.Size, retrievedBlob.Size)
	}

	if retrievedBlob.Hash != b.Hash {
		t.Errorf("Retrieved blob hash does not match. Expected %s, got %s", b.Hash, retrievedBlob.Hash)
	}
}

func TestStoreAndRetrieveTree(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mini-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tr := tree.NewTree()
	b1, _ := blob.NewBlob([]byte("file1 content"))
	b2, _ := blob.NewBlob([]byte("file2 content"))
	tr.AddEntry("file1.txt", b1.Hash, tree.EntryTypeBlob)
	tr.AddEntry("file2.txt", b2.Hash, tree.EntryTypeBlob)

	err = Store(tempDir, tr)
	if err != nil {
		t.Fatalf("Failed to store tree: %v", err)
	}

	retrievedTree, err := RetrieveTree(tempDir, tr.Hash())
	if err != nil {
		t.Fatalf("Failed to retrieve tree: %v", err)
	}

	if len(retrievedTree.Entries) != len(tr.Entries) {
		t.Errorf("Retrieved tree has different number of entries. Expected %d, got %d", len(tr.Entries), len(retrievedTree.Entries))
	}

	for i, entry := range tr.Entries {
		retrievedEntry := retrievedTree.Entries[i]
		if entry.Name != retrievedEntry.Name || entry.Hash != retrievedEntry.Hash || entry.Type != retrievedEntry.Type {
			t.Errorf("Retrieved tree entry does not match original. Expected %+v, got %+v", entry, retrievedEntry)
		}
	}
}

func TestRetrieveNonExistentObjects(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mini-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	nonExistentHash := "0123456789abcdef0123456789abcdef01234567"

	_, err = RetrieveBlob(tempDir, nonExistentHash)
	if err == nil {
		t.Errorf("Expected an error when retrieving non-existent blob, but got nil")
	}

	_, err = RetrieveTree(tempDir, nonExistentHash)
	if err == nil {
		t.Errorf("Expected an error when retrieving non-existent tree, but got nil")
	}
}
