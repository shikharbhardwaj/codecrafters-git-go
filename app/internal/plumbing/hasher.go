package plumbing

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
)

type Hash [20]byte

type Hasher struct {
	hash.Hash
}

func NewHasher(bytes []byte) Hasher {
	h := Hasher{sha1.New()}

	h.Write(bytes)

	return h
}

func (h Hasher) Sum() (hash Hash) {
	copy(hash[:], h.Hash.Sum(nil))
	return
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
