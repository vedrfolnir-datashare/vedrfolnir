package core

import (
	"Ved/lib/VedCrypto"
	"bytes"
	"crypto/rand"
	"fmt"
)

/* CreateKey
 * create a key for a keyword in a specific epoch/timestamp
 * input: keyword, the target keyword
 * 		  timestamp, the target epoch
 * output: key
 */

func CreateKey(keyword []byte, timestamp []byte) []byte {
	k_ff := "Sixteen byte key"
	k_f, err := VedCrypto.Encrypt(timestamp, []byte(k_ff))
	if err != nil {
		fmt.Println("fail to use timestamp to create key:", err)
		return nil
	}
	k, err := VedCrypto.Encrypt(keyword, []byte(k_f))
	if err != nil {
		fmt.Println("fail to use keyword to create key:", err)
		return nil
	}
	return k
}

/* CreateKey_All
 * create a group of keys for a keyword in all epoch/timestamp
 * input: keyword, the target keyword
 * output: a group of key
 */

func CreateKey_All(keyword []byte) [][]byte {
	k_ff := "Sixteen byte key"
	k_f, err := VedCrypto.Encrypt(keyword, []byte(k_ff))
	if err != nil {
		fmt.Println("fail to use timestamp to create key:", err)
		return nil
	}

	deep := 9
	var key [][]byte
	queue := [][]byte{[]byte(k_f)}
	for level := 0; level < deep; level++ {
		var nextLevel [][]byte
		for _, node := range queue {
			leftChild := VedCrypto.Hash(node)
			rightChild := VedCrypto.Hash(rightShift(node))
			nextLevel = append(nextLevel, leftChild)
			nextLevel = append(nextLevel, rightChild)
		}
		queue = nextLevel
		key = append(key, nextLevel...)
	}

	return key
}

func rightShift(input []byte) []byte {
	length := len(input)
	result := make([]byte, length)
	copy(result[1:], input[:length-1])
	return result
}

func CreateRandomNumber(n int) []byte {
	randomBytes := make([]byte, n) //128bit
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println("fail to create random number::", err)
		return nil
	}
	return randomBytes
}

func xor(a, b []byte) []byte {
	//print(len(a), " ", len(b), " ", a, " ", b)
	if len(a) != len(b) {
		panic("siles length do not match")
	}

	result := make([]byte, len(a))

	for i := range a {
		result[i] = a[i] ^ b[i]
	}

	return result
}

/* Encryption
 * encrypt a keyword to ciphertext
 * input: keyword, the target keyword
 * 		  key, owner's encryption key
 * output: the keyword's ciphertext
 */

func Encryption(keyword []byte, key []byte) []byte {

	// Wi = L || R
	// Rc -> S -> F
	// C = Wi xor (S||F)

	//print("encrypting...")
	H := VedCrypto.Hash(keyword) //256bit
	S := CreateRandomNumber(31)  //128bit
	F, err := VedCrypto.Encrypt(S, key)
	if err != nil {
		fmt.Println("fail to encrypt:", err)
		return nil
	}
	S = append(S, 0)
	Temp := append(S, F...)
	//print(len(Temp[:64]), " ", len(H), "\n")
	//print(Temp[:64], " ", H, "\n")
	//Ciphertext := xor(H, Temp[:64])
	//print(len(S), " ", len(F), " ", len(Temp), " ", len(H), "\n")
	//print(Temp, " ", H, "\n")
	Ciphertext := xor(H, Temp)
	return Ciphertext
}

/* Search
 * check whether the ciphertext belongs to the key
 * input: keyword, the target keyword
 * 		  key, owner's encryption key
 * 		  ciphertext, the target ciphertext
 * output: true/false, the ciphertext is belonged to the key?
 */

func Search(Ciphertext []byte, keyword []byte, key []byte) bool {
	//print("searching...")
	H := VedCrypto.Hash(keyword)
	Temp := xor(H, Ciphertext)
	S := Temp[:32]
	T := Temp[32:]
	//print(len(S), " ", len(T), "\n")
	//print(S, " ", T, "\n")
	Tnew, err := VedCrypto.Encrypt(S[:31], key)
	if err != nil {
		fmt.Println("fail to decrypt:", err)
		return false
	}
	//print(len(Tnew), " ", len(Tnew[:32]), "\n")
	//print(Tnew, " ", Tnew[:32], "\n")
	return bytes.Equal(Tnew[:32], T)
}
