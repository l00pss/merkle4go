package merkle4go

// Node represents a node in the Merkle tree
type Node struct {
	Hash   []byte
	Left   *Node
	Right  *Node
	IsLeaf bool
}

// NewLeafNode creates a new leaf node
func NewLeafNode(hash []byte) *Node {
	return &Node{
		Hash:   hash,
		Left:   nil,
		Right:  nil,
		IsLeaf: true,
	}
}

// NewInternalNode creates a new internal node
func NewInternalNode(left, right *Node, hash []byte) *Node {
	return &Node{
		Hash:   hash,
		Left:   left,
		Right:  right,
		IsLeaf: false,
	}
}