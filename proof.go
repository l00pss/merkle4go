package merkle4go

import "bytes"

// Proof represents a Merkle proof for a specific leaf
type Proof struct {
	Index  int
	Hashes [][]byte
	Leaf   []byte
}

// Verify verifies the Merkle proof against the given root hash
func (p *Proof) Verify(rootHash []byte, hasher Hasher) bool {
	if len(p.Hashes) == 0 && len(p.Leaf) > 0 {
		return bytes.Equal(p.Leaf, rootHash)
	}

	computedHash := p.Leaf
	index := p.Index

	for _, siblingHash := range p.Hashes {
		var combined []byte
		if index%2 == 0 {
			combined = append(computedHash, siblingHash...)
		} else {
			combined = append(siblingHash, computedHash...)
		}
		// Security: Prefix internal nodes with 0x01
		computedHash = hasher.Hash(append([]byte{0x01}, combined...))
		index /= 2
	}

	return bytes.Equal(computedHash, rootHash)
}
