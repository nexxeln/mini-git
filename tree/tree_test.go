package tree

import (
	"bytes"
	"testing"

	"github.com/nexxeln/mini-git/blob"
)

func TestNewTree(t *testing.T) {
	tree := NewTree()

	if tree == nil {
		t.Fatal("NewTree returned nil")
	}

	if len(tree.Entries) != 0 {
		t.Errorf("New tree should have 0 entries, got %d", len(tree.Entries))
	}
}

func TestAddEntry(t *testing.T) {
	tree := NewTree()
	b, _ := blob.NewBlob([]byte("test content"))
	tree.AddEntry("test.txt", b.Hash, EntryTypeBlob)

	if len(tree.Entries) != 1 {
		t.Errorf("Tree should have 1 entry, got %d", len(tree.Entries))
	}

	entry := tree.Entries[0]

	if entry.Name != "test.txt" {
		t.Errorf("Entry name should be 'test.txt', got '%s'", entry.Name)
	}

	if entry.Hash != b.Hash {
		t.Errorf("Entry hash should be '%s', got '%s'", b.Hash, entry.Hash)
	}

	if entry.Type != EntryTypeBlob {
		t.Errorf("Entry type should be %d, got %d", EntryTypeBlob, entry.Type)
	}
}

func TestTreeSerialize(t *testing.T) {
	tree := NewTree()
	b1, _ := blob.NewBlob([]byte("file content"))
	b2, _ := blob.NewBlob([]byte("another file"))
	tree.AddEntry("file.txt", b1.Hash, EntryTypeBlob)
	tree.AddEntry("another.txt", b2.Hash, EntryTypeBlob)
	serialized, err := tree.Serialize()

	if err != nil {
		t.Fatalf("Failed to serialize tree: %v", err)
	}

	if !bytes.HasPrefix(serialized, []byte("tree ")) {
		t.Errorf("Serialized tree should start with 'tree ', got %s", string(serialized[:5]))
	}

	if !bytes.Contains(serialized, []byte("file.txt")) {
		t.Errorf("Serialized tree should contain 'file.txt'")
	}

	if !bytes.Contains(serialized, []byte("another.txt")) {
		t.Errorf("Serialized tree should contain 'another.txt'")
	}
}

func TestTreeDeserialize(t *testing.T) {
	originalTree := NewTree()
	b1, _ := blob.NewBlob([]byte("file content"))
	b2, _ := blob.NewBlob([]byte("another file"))
	originalTree.AddEntry("file.txt", b1.Hash, EntryTypeBlob)
	originalTree.AddEntry("another.txt", b2.Hash, EntryTypeBlob)

	serialized, _ := originalTree.Serialize()
	deserializedTree, err := Deserialize(serialized)

	if err != nil {
		t.Fatalf("Failed to deserialize tree: %v", err)
	}

	if len(deserializedTree.Entries) != len(originalTree.Entries) {
		t.Errorf("Deserialized tree should have %d entries, got %d", len(originalTree.Entries), len(deserializedTree.Entries))
	}

	for i, entry := range originalTree.Entries {
		desEntry := deserializedTree.Entries[i]

		if entry.Name != desEntry.Name {
			t.Errorf("Entry %d: name should be '%s', got '%s'", i, entry.Name, desEntry.Name)
		}

		if entry.Hash != desEntry.Hash {
			t.Errorf("Entry %d: hash should be '%s', got '%s'", i, entry.Hash, desEntry.Hash)
		}

		if entry.Type != desEntry.Type {
			t.Errorf("Entry %d: type should be %d, got %d", i, entry.Type, desEntry.Type)
		}
	}
}
