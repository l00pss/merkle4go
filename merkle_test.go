package merkle4go

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestMerkleTreeCreation(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	if tree.Root() == nil {
		t.Fatal("Root hash is nil")
	}

	if len(tree.Root()) != 32 {
		t.Fatalf("Expected root hash length 32, got %d", len(tree.Root()))
	}
}

func TestMerkleTreeWithBytes(t *testing.T) {
	data := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("merkle"),
		[]byte("tree"),
	}
	hasher := SHA256Hasher{}
	converter := ByteConverter{}

	tree, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	if tree.Root() == nil {
		t.Fatal("Root hash is nil")
	}
}

func TestMerkleProof(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	// Test proof for each element
	for i := range data {
		proof, err := tree.GetProof(i)
		if err != nil {
			t.Fatalf("Failed to get proof for index %d: %v", i, err)
		}

		if proof.Index != i {
			t.Fatalf("Expected proof index %d, got %d", i, proof.Index)
		}

		// Verify proof
		if !proof.Verify(tree.Root(), hasher) {
			t.Fatalf("Proof verification failed for index %d", i)
		}
	}
}

func TestMerkleProofInvalidIndex(t *testing.T) {
	data := []string{"a", "b", "c"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	// Test invalid indices
	_, err = tree.GetProof(-1)
	if err == nil {
		t.Fatal("Expected error for negative index")
	}

	_, err = tree.GetProof(len(data))
	if err == nil {
		t.Fatal("Expected error for index out of range")
	}
}

func TestMerkleTreeVerify(t *testing.T) {
	originalData := []string{"a", "b", "c", "d"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree, err := NewMerkleTree(originalData, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree: %v", err)
	}

	// Verify with original data
	if !tree.Verify(originalData) {
		t.Fatal("Verification failed for original data")
	}

	// Verify with modified data
	modifiedData := []string{"a", "b", "c", "x"}
	if tree.Verify(modifiedData) {
		t.Fatal("Verification should fail for modified data")
	}

	// Verify with different length data
	differentLengthData := []string{"a", "b", "c"}
	if tree.Verify(differentLengthData) {
		t.Fatal("Verification should fail for different length data")
	}
}

func TestEmptyData(t *testing.T) {
	data := []string{}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	_, err := NewMerkleTree(data, hasher, converter)
	if err == nil {
		t.Fatal("Expected error for empty data")
	}
}

func TestSingleElement(t *testing.T) {
	data := []string{"single"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree with single element: %v", err)
	}

	if tree.Root() == nil {
		t.Fatal("Root hash is nil for single element tree")
	}

	proof, err := tree.GetProof(0)
	if err != nil {
		t.Fatalf("Failed to get proof for single element: %v", err)
	}

	if !proof.Verify(tree.Root(), hasher) {
		t.Fatal("Proof verification failed for single element")
	}
}

func TestOddNumberOfElements(t *testing.T) {
	data := []string{"a", "b", "c"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create Merkle tree with odd number of elements: %v", err)
	}

	// Test all proofs
	for i := range data {
		proof, err := tree.GetProof(i)
		if err != nil {
			t.Fatalf("Failed to get proof for index %d: %v", i, err)
		}

		if !proof.Verify(tree.Root(), hasher) {
			t.Fatalf("Proof verification failed for index %d", i)
		}
	}
}

func TestConsistentRootHash(t *testing.T) {
	data := []string{"test1", "test2", "test3", "test4"}
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	tree1, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create first Merkle tree: %v", err)
	}

	tree2, err := NewMerkleTree(data, hasher, converter)
	if err != nil {
		t.Fatalf("Failed to create second Merkle tree: %v", err)
	}

	if !bytes.Equal(tree1.Root(), tree2.Root()) {
		t.Fatal("Root hashes should be identical for same data")
	}
}

func ExampleMerkleTree() {
	// Create some sample data
	data := []string{"transaction1", "transaction2", "transaction3", "transaction4"}

	// Initialize hasher and converter
	hasher := SHA256Hasher{}
	converter := StringConverter{}

	// Create Merkle tree
	tree, _ := NewMerkleTree(data, hasher, converter)

	// Get root hash
	root := tree.Root()
	println("Root hash:", hex.EncodeToString(root))

	// Generate proof for first transaction
	proof, _ := tree.GetProof(0)

	// Verify proof
	isValid := proof.Verify(root, hasher)
	println("Proof valid:", isValid)
}
