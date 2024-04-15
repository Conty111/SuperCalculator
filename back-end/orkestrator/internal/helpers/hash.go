package helpers

import (
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"hash"
	"strings"
)

type HashHelper struct {
	hasher hash.Hash
}

func NewHasher(hashAlg string) (*HashHelper, error) {
	var hasher hash.Hash
	switch strings.ToLower(hashAlg) {
	case "", "sha256":
		hasher = sha256.New()
	case "md5":
		hasher = md5.New() //nolint:gosec
	default:
		return nil, errors.New("Invalid hash alghorithm name")
	}
	return &HashHelper{
		hasher: hasher,
	}, nil
}

// HashString хэширует строку с и возвращает хэш в виде строки
func (h *HashHelper) HashString(input string) (string, error) {
	defer h.hasher.Reset()
	_, err := h.hasher.Write([]byte(input))
	if err != nil {
		return "", err
	}

	hashBytes := h.hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}
