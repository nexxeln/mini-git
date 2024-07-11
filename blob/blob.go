package blob

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
)

type Blob struct {
	Content []byte
	Size    int64
	Hash    string
}

func NewBlob(content []byte) (*Blob, error) {
	size := int64(len(content))
	hash := sha1.Sum(content)

	return &Blob{
		Content: content,
		Size:    size,
		Hash:    hex.EncodeToString(hash[:]),
	}, nil
}

func (b *Blob) Serialize() ([]byte, error) {
	header := fmt.Sprintf("blob %d\x00", b.Size)
	return append([]byte(header), b.Content...), nil
}

func Deserialize(data []byte) (*Blob, error) {
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid blob data: no null byte found")
	}

	header := string(data[:nullIndex])
	if !bytes.HasPrefix(data, []byte("blob ")) {
		return nil, fmt.Errorf("invalid blob data: incorrect header")
	}

	sizeStr := header[5:]
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid blob size: %v", err)
	}

	content := data[nullIndex+1:]
	if int64(len(content)) != size {
		return nil, fmt.Errorf("content size mismatch")
	}

	return NewBlob(content)
}
