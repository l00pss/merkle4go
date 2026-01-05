package merkle4go

import (
	"crypto/sha256"
)

// Hasher defines the interface for hash functions
type Hasher interface {
	Hash(data []byte) []byte
}

// SHA256Hasher implements Hasher using SHA-256
type SHA256Hasher struct{}

// Hash computes SHA-256 hash of the given data
func (h SHA256Hasher) Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// DataConverter converts data of type T to bytes
type DataConverter[T any] interface {
	ToBytes(data T) []byte
}

// StringConverter converts string data to bytes
type StringConverter struct{}

func (c StringConverter) ToBytes(data string) []byte {
	return []byte(data)
}

// ByteConverter implements DataConverter for byte slices
type ByteConverter struct{}

func (c ByteConverter) ToBytes(data []byte) []byte {
	return data
}
