package commit

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	TreeHash   string
	ParentHash string
	Author     string
	Committer  string
	AuthorDate time.Time
	CommitDate time.Time
	Message    string
}

func NewCommit(treeHash, parentHash, author, committer, message string) *Commit {
	now := time.Now()
	return &Commit{
		TreeHash:   treeHash,
		ParentHash: parentHash,
		Author:     author,
		Committer:  committer,
		AuthorDate: now,
		CommitDate: now,
		Message:    message,
	}
}

func (c *Commit) Serialize() ([]byte, error) {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("tree %s\n", c.TreeHash))
	if c.ParentHash != "" {
		buffer.WriteString(fmt.Sprintf("parent %s\n", c.ParentHash))
	}
	buffer.WriteString(fmt.Sprintf("author %s %d +0000\n", c.Author, c.AuthorDate.Unix()))
	buffer.WriteString(fmt.Sprintf("committer %s %d +0000\n", c.Committer, c.CommitDate.Unix()))
	buffer.WriteString("\n")
	buffer.WriteString(c.Message)

	content := buffer.Bytes()
	header := fmt.Sprintf("commit %d\x00", len(content))
	return append([]byte(header), content...), nil
}

func Deserialize(data []byte) (*Commit, error) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid commit data: no null byte found")
	}

	header := string(data[:nullIndex])
	if !strings.HasPrefix(header, "commit ") {
		return nil, fmt.Errorf("invalid commit data: incorrect header")
	}

	content := data[nullIndex+1:]
	commit := &Commit{}
	lines := bytes.Split(content, []byte("\n"))

	var message strings.Builder
	inMessage := false

	for _, line := range lines {
		if inMessage {
			message.Write(line)
			message.WriteString("\n")
			continue
		}

		if len(line) == 0 {
			inMessage = true
			continue
		}

		parts := bytes.SplitN(line, []byte(" "), 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid commit line: %s", string(line))
		}

		key := string(parts[0])
		value := string(parts[1])

		switch key {
		case "tree":
			commit.TreeHash = value
		case "parent":
			commit.ParentHash = value
		case "author":
			commit.Author, commit.AuthorDate = parseAuthorLine(value)
		case "committer":
			commit.Committer, commit.CommitDate = parseAuthorLine(value)
		default:
			return nil, fmt.Errorf("unknown commit field: %s", key)
		}
	}

	commit.Message = strings.TrimSpace(message.String())

	return commit, nil
}

func parseAuthorLine(line string) (string, time.Time) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return line, time.Time{}
	}

	timestamp, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
	if err != nil {
		return line, time.Time{}
	}

	return strings.Join(parts[:len(parts)-2], " "), time.Unix(timestamp, 0).UTC()
}

func (c *Commit) Hash() string {
	serialized, _ := c.Serialize()
	hash := sha1.Sum(serialized)
	return hex.EncodeToString(hash[:])
}
