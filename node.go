package merkle4go

type Node struct {
	Hash   []byte
	Left   *Node
	Right  *Node
	IsLeaf bool
}

func NewLeafNode(hash []byte) *Node {
	return &Node{
		Hash:   hash,
		Left:   nil,
		Right:  nil,
		IsLeaf: true,
	}
}

func NewInternalNode(left, right *Node, hash []byte) *Node {
	return &Node{
		Hash:   hash,
		Left:   left,
		Right:  right,
		IsLeaf: false,
	}
}
