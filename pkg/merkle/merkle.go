package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

// Prune finds the hash and removes all nodes uneccessary for proving it
// belongs to the merkle tree. Returns an error if the hash does not exist.
func (mt *Tree) Prune(hash string) error {
	n, err := mt.Root.bfs(hash)
	if err != nil {
		return err
	}

	prune(n)
	return nil
}

// Prune removes all nodes in a tree that aren't neccessary for proving n
// belongs to its merkle tree. The return value is the root of the tree.
func prune(n *Node) *Node {
	if n.parent == nil {
		return n
	}

	if n == n.parent.Left {
		n.parent.Right = nil
	} else {
		n.parent.Left = nil
	}

	return prune(n.parent)
}

// bfs breadth first search
func (n *Node) bfs(hash string) (*Node, error) {
	nodes := []*Node{n}
	for len(nodes) > 0 {
		node, nodes := nodes[0], nodes[1:]
		if node.Hash == hash {
			return node, nil
		}
		if node.Left != nil {
			nodes = append(nodes, node.Left)
		}
		if node.Right != nil {
			nodes = append(nodes, node.Right)
		}
	}
	return nil, errors.New("The hash %s does not exist")
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
