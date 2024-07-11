package tree

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type EntryType int

const (
	EntryTypeBlob EntryType = iota
	EntryTypeTree
)

type Entry struct {
	Name string
	Hash string
	Type EntryType
}

type Tree struct {
	Entries []Entry
}

func NewTree() *Tree {
	return &Tree{
		Entries: make([]Entry, 0),
	}
}

func (t *Tree) AddEntry(name, hash string, entryType EntryType) {
	t.Entries = append(t.Entries, Entry{
		Name: name,
		Hash: hash,
		Type: entryType,
	})
}

func (t *Tree) Serialize() ([]byte, error) {
	var buffer bytes.Buffer

	for _, entry := range t.Entries {
		entryString := fmt.Sprintf("%06o %s %s\t%s\n",
			func() int {
				if entry.Type == EntryTypeBlob {
					return 100644 // regular file
				}
				return 40000 // directory
			}(),
			func() string {
				if entry.Type == EntryTypeBlob {
					return "blob"
				}
				return "tree"
			}(),
			entry.Hash,
			entry.Name)
		buffer.WriteString(entryString)
	}

	content := buffer.Bytes()
	header := fmt.Sprintf("tree %d\x00", len(content))
	return append([]byte(header), content...), nil
}

func Deserialize(data []byte) (*Tree, error) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid tree data: no null byte found")
	}

	header := string(data[:nullIndex])
	if !strings.HasPrefix(header, "tree ") {
		return nil, fmt.Errorf("invalid tree data: incorrect header")
	}

	tree := NewTree()
	content := data[nullIndex+1:]
	lines := bytes.Split(content, []byte("\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := bytes.Split(line, []byte("\t"))
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tree entry: %s", string(line))
		}

		entryParts := bytes.Fields(parts[0])
		if len(entryParts) != 3 {
			return nil, fmt.Errorf("invalid tree entry parts: %s", string(parts[0]))
		}

		mode, err := strconv.ParseInt(string(entryParts[0]), 8, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid mode: %s", string(entryParts[0]))
		}

		entryType := EntryTypeBlob
		if mode == 40000 {
			entryType = EntryTypeTree
		}

		tree.AddEntry(string(parts[1]), string(entryParts[2]), entryType)
	}

	return tree, nil
}

func (t *Tree) Hash() string {
	serialized, _ := t.Serialize()
	hash := sha1.Sum(serialized)
	return hex.EncodeToString(hash[:])
}
