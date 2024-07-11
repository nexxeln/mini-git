package commit

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCommit(t *testing.T) {
	treeHash := "0123456789abcdef0123456789abcdef01234567"
	parentHash := "fedcba9876543210fedcba9876543210fedcba98"
	author := "John Doe <john@example.com>"
	committer := "Jane Doe <jane@example.com>"
	message := "Initial commit"

	commit := NewCommit(treeHash, parentHash, author, committer, message)

	if commit.TreeHash != treeHash {
		t.Errorf("Expected tree hash %s, got %s", treeHash, commit.TreeHash)
	}
	if commit.ParentHash != parentHash {
		t.Errorf("Expected parent hash %s, got %s", parentHash, commit.ParentHash)
	}
	if commit.Author != author {
		t.Errorf("Expected author %s, got %s", author, commit.Author)
	}
	if commit.Committer != committer {
		t.Errorf("Expected committer %s, got %s", committer, commit.Committer)
	}
	if commit.Message != message {
		t.Errorf("Expected message %s, got %s", message, commit.Message)
	}
}

func TestCommitSerialize(t *testing.T) {
	treeHash := "0123456789abcdef0123456789abcdef01234567"
	parentHash := "fedcba9876543210fedcba9876543210fedcba98"
	author := "John Doe <john@example.com>"
	committer := "Jane Doe <jane@example.com>"
	message := "Initial commit"
	timestamp := time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC)

	commit := NewCommit(treeHash, parentHash, author, committer, message)
	commit.AuthorDate = timestamp
	commit.CommitDate = timestamp

	serialized, err := commit.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize commit: %v", err)
	}

	expected := fmt.Sprintf(`commit 216%ctree 0123456789abcdef0123456789abcdef01234567
parent fedcba9876543210fedcba9876543210fedcba98
author John Doe <john@example.com> 1625097600 +0000
committer Jane Doe <jane@example.com> 1625097600 +0000

Initial commit`, 0)

	if string(serialized) != expected {
		t.Errorf("Serialized commit does not match expected content.\nExpected:\n%q\nGot:\n%q", expected, string(serialized))
	}
}

func TestCommitDeserialize(t *testing.T) {
	commitData := `commit 216` + "\x00" + `tree 0123456789abcdef0123456789abcdef01234567
parent fedcba9876543210fedcba9876543210fedcba98
author John Doe <john@example.com> 1625097600 +0000
committer Jane Doe <jane@example.com> 1625097600 +0000

Initial commit`

	commit, err := Deserialize([]byte(commitData))
	if err != nil {
		t.Fatalf("Failed to deserialize commit: %v", err)
	}

	if commit.TreeHash != "0123456789abcdef0123456789abcdef01234567" {
		t.Errorf("Incorrect tree hash")
	}
	if commit.ParentHash != "fedcba9876543210fedcba9876543210fedcba98" {
		t.Errorf("Incorrect parent hash")
	}
	if commit.Author != "John Doe <john@example.com>" {
		t.Errorf("Incorrect author")
	}
	if commit.Committer != "Jane Doe <jane@example.com>" {
		t.Errorf("Incorrect committer")
	}
	if commit.Message != "Initial commit" {
		t.Errorf("Incorrect message")
	}
	expectedTime := time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC)
	if !commit.AuthorDate.Equal(expectedTime) {
		t.Errorf("Incorrect author date")
	}
	if !commit.CommitDate.Equal(expectedTime) {
		t.Errorf("Incorrect commit date")
	}
}
