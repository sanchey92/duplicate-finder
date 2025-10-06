package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
)

var supportedAlgorithms = map[string]bool{
	"md5":    true,
	"sha256": true,
}

type Hasher struct {
	algorithm string
}

func New(alg string) (*Hasher, error) {
	if !supportedAlgorithms[alg] {
		return nil, fmt.Errorf("unsupported hash algorithm: %s (supported md5 & sha256)", alg)
	}

	return &Hasher{algorithm: alg}, nil
}

func (h *Hasher) Calculate(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer file.Close()

	hasher := h.createHasher()
	if _, err = io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("cannot calculate hash for %s: %w", path, err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func (h *Hasher) GetAlgorithm() string {
	return h.algorithm
}

func (h *Hasher) IsSupported(alg string) bool {
	return supportedAlgorithms[alg]
}

func (h *Hasher) createHasher() hash.Hash {
	switch h.algorithm {
	case "md5":
		return md5.New()
	case "sha256":
		return sha256.New()
	default:
		panic(fmt.Errorf("usupported algorithm"))
	}
}
