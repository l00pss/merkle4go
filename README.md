# Merkle4Go

A simple and flexible Merkle Tree library written in Go. It works with any data type.

## Installation

```bash
go get github.com/l00pss/merkle4go
```

## How to Use?

### Basic Example (With Strings)

```go
package main

import (
    "fmt"
    "github.com/l00pss/merkle4go"
)

func main() {
    // Our data
    data := []string{"data1", "data2", "data3", "data4"}
    
    // Create the tree
    hasher := merkle4go.SHA256Hasher{}
    converter := merkle4go.StringConverter{}
    
    tree, _ := merkle4go.NewMerkleTree(data, hasher, converter)
    
    // Get Root hash
    fmt.Printf("Root Hash: %x\n", tree.Root())
    
    // Generate and verify Proof
    proof, _ := tree.GetProof(0)
    valid := proof.Verify(tree.Root(), hasher)
    
    fmt.Println("Verification result:", valid)
}
```

### Working with Byte Data

If you have raw `[]byte` data, you can use `ByteConverter`:

```go
data := [][]byte{
    []byte("hello"),
    []byte("world"),
}

tree, _ := merkle4go.NewMerkleTree(data, merkle4go.SHA256Hasher{}, merkle4go.ByteConverter{})
```

### Working with Custom Data Types

For your own structs, you just need to implement the `ToBytes` method:

```go
type User struct {
    Name string
    Id   int
}

// Define Converter
type UserConverter struct{}

func (c UserConverter) ToBytes(u User) []byte {
    return []byte(fmt.Sprintf("%s:%d", u.Name, u.Id))
}

// Usage
users := []User{{Name: "Alice", Id: 1}, {Name: "Bob", Id: 2}}
tree, _ := merkle4go.NewMerkleTree(users, merkle4go.SHA256Hasher{}, UserConverter{})
```

## License

MIT
