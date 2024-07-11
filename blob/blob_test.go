package blob

import (
	"bytes"
	"testing"
)

func TestNewBlob(t *testing.T) {
	content := []byte("hello world!")
	blob, err := NewBlob(content)
	if err != nil {
		t.Fatalf("Failed to create new blob: %v", err)
	}

	if blob.Size != int64(len(content)) {
		t.Errorf("Incorrect blob size. Expected %d, got %d", len(content), blob.Size)
	}

	if !bytes.Equal(blob.Content, content) {
		t.Errorf("Blob content does not match input content")
	}

	expectedHash := "430ce34d020724ed75a196dfc2ad67c77772d169"
	if blob.Hash != expectedHash {
		t.Errorf("Incorrect hash. Expected %s, got %s", expectedHash, blob.Hash)
	}
}

func TestBlobSerialize(t *testing.T) {
	content := []byte("Hello, World!")
	blob, _ := NewBlob(content)

	serialized, err := blob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}

	expectedPrefix := []byte("blob 13\x00")
	if !bytes.HasPrefix(serialized, expectedPrefix) {
		t.Errorf("Serialized blob does not have correct prefix")
	}

	if !bytes.Equal(serialized[len(expectedPrefix):], content) {
		t.Errorf("Serialized blob content does not match original content")
	}
}

func TestBlobDeserialize(t *testing.T) {
	content := []byte("Hello, World!")
	originalBlob, _ := NewBlob(content)
	serialized, _ := originalBlob.Serialize()

	deserializedBlob, err := Deserialize(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize blob: %v", err)
	}

	if deserializedBlob.Size != originalBlob.Size {
		t.Errorf("Deserialized blob size does not match. Expected %d, got %d", originalBlob.Size, deserializedBlob.Size)
	}

	if !bytes.Equal(deserializedBlob.Content, originalBlob.Content) {
		t.Errorf("Deserialized blob content does not match original content")
	}

	if deserializedBlob.Hash != originalBlob.Hash {
		t.Errorf("Deserialized blob hash does not match. Expected %s, got %s", originalBlob.Hash, deserializedBlob.Hash)
	}
}
