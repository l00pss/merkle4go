package merkle4go

import (
	"bytes"
	"errors"
)

// MerkleTree represents a generic Merkle tree
type MerkleTree[T any] struct {
	root      *Node
	leaves    []*Node
	hasher    Hasher
	converter DataConverter[T]
}

// NewMerkleTree creates a new Merkle tree with the given data
func NewMerkleTree[T any](data []T, hasher Hasher, converter DataConverter[T]) (*MerkleTree[T], error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}

	tree := &MerkleTree[T]{
		hasher:    hasher,
		converter: converter,
	}

	leaves := make([]*Node, len(data))
	for i, item := range data {
		bytes := converter.ToBytes(item)
		hash := hasher.Hash(append([]byte{0x00}, bytes...))
		leaves[i] = NewLeafNode(hash)
	}
	tree.leaves = leaves

	tree.root = tree.buildTree(leaves)
	return tree, nil
}

func (mt *MerkleTree[T]) buildTree(nodes []*Node) *Node {
	if len(nodes) == 1 {
		return nodes[0]
	}

	nextLevel := make([]*Node, 0, (len(nodes)+1)/2)

	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		var right *Node

		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			right = nodes[i]
		}

		combined := append(left.Hash, right.Hash...)
		parentHash := mt.hasher.Hash(append([]byte{0x01}, combined...))
		parent := NewInternalNode(left, right, parentHash)
		nextLevel = append(nextLevel, parent)
	}

	return mt.buildTree(nextLevel)
}

// Root returns the root hash of the tree
func (mt *MerkleTree[T]) Root() []byte {
	if mt.root == nil {
		return nil
	}
	return mt.root.Hash
}

// GetProof generates a Merkle proof for the element at the given index
func (mt *MerkleTree[T]) GetProof(index int) (*Proof, error) {
	if index < 0 || index >= len(mt.leaves) {
		return nil, errors.New("index out of range")
	}

	proof := &Proof{
		Index:  index,
		Hashes: [][]byte{},
		Leaf:   mt.leaves[index].Hash,
	}

	if len(mt.leaves) == 1 {
		return proof, nil
	}

	nodes := mt.leaves
	currentIndex := index

	for len(nodes) > 1 {
		nextLevel := make([]*Node, 0, (len(nodes)+1)/2)

		for i := 0; i < len(nodes); i += 2 {
			left := nodes[i]
			var right *Node

			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				right = nodes[i]
			}

			if i == currentIndex || i+1 == currentIndex {
				if i == currentIndex {
					proof.Hashes = append(proof.Hashes, right.Hash)
				} else {
					proof.Hashes = append(proof.Hashes, left.Hash)
				}
			}

			combined := append(left.Hash, right.Hash...)
			parentHash := mt.hasher.Hash(append([]byte{0x01}, combined...))
			parent := NewInternalNode(left, right, parentHash)
			nextLevel = append(nextLevel, parent)
		}

		nodes = nextLevel
		currentIndex /= 2
	}

	return proof, nil
}

// Verify verifies that the given data produces the expected root hash
func (mt *MerkleTree[T]) Verify(data []T) bool {
	if len(data) != len(mt.leaves) {
		return false
	}

	for i, item := range data {
		bytesData := mt.converter.ToBytes(item)
		hash := mt.hasher.Hash(append([]byte{0x00}, bytesData...))
		if !bytes.Equal(hash, mt.leaves[i].Hash) {
			return false
		}
	}

	return true
}
