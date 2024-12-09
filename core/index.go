package core

import (
	"Ved/lib/VedCrypto"
	"bytes"
	"fmt"
)

/* InitIndex
 * init an index for an owner (without existency check)
 * input: key0 []byte("ThisIs16ByteKey0")
 * output: root node of a HSMT
 */

func InitIndex(key0 []byte) *MerkleNode {
	var ciphertextArray [][]byte
	var leaf []*MerkleNode
	var root *MerkleNode

	for i := 0; i <= 15; i++ {
		plaintext := []byte(fmt.Sprintf("%x", i))
		ciphertext, err := VedCrypto.Encrypt(plaintext, key0)
		if err != nil {
			fmt.Println("Error encrypting:", err)
			return nil
		}
		ciphertextArray = append(ciphertextArray, ciphertext)
	}

	for _, ciphertext := range ciphertextArray {
		hash := VedCrypto.Hash(ciphertext, ciphertext)
		node := NewMerkleNode(nil, nil, hash)
		leaf = append(leaf, node)
	}

	for i := 0; i < 64; i++ {
		//print("create ", i, " merkel tree\n")
		var middle []*MerkleNode
		root = BuildMerkleTree(leaf)

		for i := 0; i <= 15; i++ {
			//hash := VedCrypto.Hash(root.Hash, []byte(fmt.Sprintf("%x", i)))
			//hash := VedCrypto.Hash(root.Hash, root.Hash)
			//node := NewMerkleNode(root, root, nil, hash)
			middle = append(middle, root)
		}

		leaf = middle
	}

	return root
}

/* NewIndex
 * generate a new index for a keyword (without existence check)
 * input: key1 []byte("ThisIs16ByteKey1"), key for generating an index
 *	  deep 2, then the leaf node number is 2^2 = 4
 * output: root node of a BHT-CMT
 */

func NewIndex(key1 []byte, deep int) *MerkleNode {
	var key1hash = VedCrypto.Hash(key1)
	if deep == 1 {
		rootHash := VedCrypto.Hash(key1hash)
		return NewMerkleNode(nil, nil, rootHash)
	}

	var oldBHT [][]byte
	oldBHT = append(oldBHT, key1hash)

	for i := 1; i < deep; i++ {
		var newBHT [][]byte
		for _, parent := range oldBHT {
			left := VedCrypto.Hash(parent)
			right := make([]byte, len(left))
			copy(right, left)
			right[0] ^= 0x01
			//print("current: ", string(left), " ", string(right), "\n")
			newBHT = append(newBHT, left, right)
		}

		oldBHT = newBHT
	}

	var nodes []*MerkleNode
	for _, v := range oldBHT {
		hash := VedCrypto.Hash(v)
		node := NewMerkleNode(nil, nil, hash)
		nodes = append(nodes, node)
	}

	root := BuildMerkleTree(nodes)
	return root
}

/* UpdateIndex
 * add a new keyword index to a HSMT
 * input: proofs, for update route
 *		  index, candidate append index
 * output: parent, new root for HSMT
 */

func UpdateIndex(proofs []*Proof, index *MerkleNode) *MerkleNode {
	// print("index hash: ", string(index.Hash), "\n")
	parent := index

	// print("start update index...\n")
	// print("length of proofs: ", len(proofs), "\n")
	for i := len(proofs) - 1; i >= 0; i-- {
		//print(proofs[i].isLeft, " ")
		if proofs[i].isLeft {
			hash := VedCrypto.Hash(proofs[i].Node.Hash, parent.Hash)
			parent = NewMerkleNode(proofs[i].Node, parent, hash)
		} else {
			hash := VedCrypto.Hash(parent.Hash, proofs[i].Node.Hash)
			parent = NewMerkleNode(parent, proofs[i].Node, hash)
		}
	}
	return parent
}

/* CheckExistence
 * check a keyword's existence
 * input: keyword []byte("apple")
 *		  root, the HSMT where the keyword checks existence
 *		  nonexistenceProof, the proof of nonexistence.
 *			it should be Hash(AES(key0, char), AES(key0, char)), then return false
 * output: proofs, root to the target keyword's hash
 * 		   true/false: the keyword exist?
 */

func CheckExistence(keyword []byte, root *MerkleNode, nonexistenceProof []byte) ([]*Proof, bool, *MerkleNode) {
	// print("start get proof...\n")
	proofs, target := GetProof(keyword, root)
	if proofs == nil {
		fmt.Println("fail to search keyword")
		return nil, true, nil
	}

	// print("start verify proof...\n")
	// print("proof and index: ", string(nonexistenceProof), " ", string(target.Hash), "\n")
	if !bytes.Equal(nonexistenceProof, target.Hash) {
		return proofs, true, target
	}
	return proofs, false, target
}

/* GetProof
 * get a keyword's non/existence proof
 * input: keyword []byte("apple")
 *		  root, the HSMT where the keyword checks existence
 * output: proofs, route to the target keyword's hash
 * 		   pre, the leaf node that the keyword should lie in
 */

func GetProof(keyword []byte, root *MerkleNode) ([]*Proof, *MerkleNode) {
	var proofs []*Proof
	var pre = root
	keywordHash := VedCrypto.Hash(keyword)
	// print("current keyword hash: ", string(keywordHash), "\n")
	keywordHashStr := string(keywordHash)
	for _, t := range keywordHashStr {
		//print("current character: ", int(t), "\n")
		t16 := transferTo16(t)
		//print(i, " current character after transfer: ", t16, "\n")
		proof, cur := getProofForaChar(t16, pre)
		if proof == nil {
			fmt.Println("fail to search keyword")
			return nil, nil
		}
		pre = cur
		proofs = append(proofs, proof...)
	}
	return proofs, pre
}

func getProofForaChar(target int, node *MerkleNode) ([]*Proof, *MerkleNode) {
	low, high := 0, 15
	var proofs []*Proof

	for i := 0; i < 4; i++ {
		mid := (high + low) * 5
		if target*10 < mid {
			proof := NewProof(false, node.Right)
			proofs = append(proofs, proof)
			node = node.Left
			high = mid / 10
			//print("target ", target, " < middle ", mid/10, " find left child ", low, " ", high, " ", string(node.Hash), "\n")
		} else if target*10 > mid {
			proof := NewProof(true, node.Left)
			proofs = append(proofs, proof)
			node = node.Right
			low = mid/10 + 1
			//print("middle ", mid/10, " < target ", target, ", find right child ", low, " ", high, " ", string(node.Hash), "\n")
		} else {
			print("wrong calculate on char")
			return nil, nil
		}

	}

	return proofs, node
}

/* VerifyProof
 * Verify a keyword/key's non/existence proof
 * input: proofs, route to the target keyword's hash
 *		  root, the HSMT where the keyword checks existence
 *		  leaf, the target keyword/key
 * output: true/false, the leaf is correct?
 */

func VerifyProof(proofs []*Proof, root *MerkleNode, leaf *MerkleNode) bool {
	hash := leaf.Hash
	for i := len(proofs) - 1; i >= 0; i-- {
		if proofs[i].isLeft {
			hash = VedCrypto.Hash(proofs[i].Node.Hash, hash)
		} else {
			hash = VedCrypto.Hash(hash, proofs[i].Node.Hash)
		}
	}

	if bytes.Equal(hash, root.Hash) {
		return true
	}
	return false
}

//Helper Function

// MerkleNode
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

// NewMerkleNode
func NewMerkleNode(left, right *MerkleNode, hash []byte) *MerkleNode {
	return &MerkleNode{left, right, hash}
}

func BuildMerkleTree(nodes []*MerkleNode) *MerkleNode {
	deep := 0
	for len(nodes) > 1 {
		var levelNodes []*MerkleNode
		//print(len(nodes), "\n")
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				// print("current node hash: ", string(nodes[i].Hash), " ", string(nodes[i+1].Hash), "\n")
				hash := VedCrypto.Hash(nodes[i].Hash, nodes[i+1].Hash)
				// print("current node hash result: ", string(hash), "\n")
				node := NewMerkleNode(nodes[i], nodes[i+1], hash[:])
				levelNodes = append(levelNodes, node)
			} else {
				/*
					hash := VedCrypto.Hash(nodes[i].Hash, nodes[i].Hash)
					node := NewMerkleNode(nodes[i], nodes[i], nil, hash[:])
					nodes[i].Parent = node
					levelNodes = append(levelNodes, node)

				*/
				print("illegal node number")
				return nil
			}
		}

		nodes = levelNodes
		deep++
	}
	//print("depth ", deep, "\n")
	return nodes[0]
}

type Proof struct {
	isLeft bool
	Node   *MerkleNode
}

func NewProof(isLeft bool, node *MerkleNode) *Proof {
	return &Proof{isLeft, node}
}

func transferTo16(a int32) int {
	if a > 58 {
		return int(a) - 87
	}
	return int(a) - 48
}
