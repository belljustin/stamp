package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
)

// Tree is a binary Merkle tree. The leaves form a set of SHA256 hashes.
// Internal nodes are the hash of their children.
type Tree struct {
	Root *Node `json:"root"`
}

// Node is a node in a Merkle Tree
type Node struct {
	Hash   string `json:"hash"`
	Left   *Node  `json:"left"`
	Right  *Node  `json:"right"`
	parent *Node
}

// NewTree creates a MerkleTree with data as the leaf nodes
func NewTree(data []string) *Tree {
	// depth of tree without the leaves
	depth := int(math.Ceil(math.Log2(float64(len(data)))))

	/* The path to a leaf node in a binary tree can be represented by a
	bitstring. Going left is indicated by 0. Going right is indicated by 1.
	If we index the nodes at a depth from left to right, we can access a
	given node from the root using the binary representation of the index
	such that the number of digits is equal to the depth chosen.

			 		root
			0				1			// depth 1, 1 digits
		00		01		10		11		// depth 2, 2 digits

	So to get to index 2 at depth 2, [1, 0] = [right, left] */
	root := new(Node)
	for i := 0; i < len(data); i++ {
		bs := toBitstring(i, depth)

		node := root
		for _, b := range bs {
			if string(b) == "0" {
				node = node.addLeft()
			} else {
				node = node.addRight()
			}
		}
		node.Hash = data[i]
		node.parent.update()
	}
	return &Tree{root}
}

func (t *Tree) Validate() bool {
	return t.Root.validate()
}

func (n *Node) validate() bool {
	valid := true

	var data string
	if n.Left != nil {
		valid = valid && n.Left.validate()
		data = n.Left.Hash
	}

	if n.Right != nil {
		valid = valid && n.Right.validate()
		data = data + n.Right.Hash
	}

	if len(data) > 0 {
		hash := sha256.Sum256([]byte(data))
		valid = valid && (n.Hash == hex.EncodeToString(hash[:]))
	}

	return valid
}

func (n *Node) addLeft() *Node {
	if n.Left == nil {
		n.Left = new(Node)
		n.Left.parent = n
	}
	return n.Left
}

func (n *Node) addRight() *Node {
	if n.Right == nil {
		n.Right = new(Node)
		n.Right.parent = n
	}
	return n.Right
}

func (n *Node) update() {
	var data string
	if n.Left != nil {
		data = n.Left.Hash
	}

	if n.Right != nil {
		data = data + n.Right.Hash
	}

	if len(data) > 0 {
		hash := sha256.Sum256([]byte(data))
		n.Hash = hex.EncodeToString(hash[:])
	}

	if n.parent != nil {
		n.parent.update()
	}
}

func toBitstring(n int, l int) string {
	b := strconv.FormatInt(int64(n), 2)
	if len(b) < l {
		b = strings.Repeat("0", l-len(b)) + b
	}
	return b
}
