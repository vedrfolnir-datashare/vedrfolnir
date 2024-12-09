package main

import (
	"Ved/core"
)

/*
func TestInitIndex() {

	var key0, key1 []byte
	key0, key1, _ = core.InitOwner()
	fmt.Printf(string(key0), string(key1))

	startTime := time.Now()
	core.InitIndex(key0)
	endTime := time.Now()
	executionTime := endTime.Sub(startTime)
	fmt.Printf("\n time：%s\n", executionTime)

}
*/
func TestUpdateIndex() {
	// key1 := []byte("ThisIs16ByteKey1")
	// hash := VedCrypto.Hash(key1)
	// hash2 := VedCrypto.Hash(hash)
	// hash3 := VedCrypto.Hash(hash2)
	// print(string(hash), "\n", string(hash2), "\n", string(hash3))

	//key0 := []byte("ThisIs16ByteKey0")
	//HSMT := core.InitIndex(key0)
	//print("current root: ", string(HSMT.Hash), "\n")

	key1 := []byte("ThisIs16ByteKey1")
	deep := 19
	newIndex := core.NewIndex(key1, deep)
	print(newIndex)
	// print(string(newIndex.Left.Hash))
	//deep=1, key = a90772b617cac27c9283990de61165bd2021eee6cb84c13ff438fc82e0d1a0b3
	//		  root (hashKey) = a25d25735556a2015a6e9c80b10e45999f2b71d131a39d48247b0978333008ee
	//deep=2, key = a25d25735556a2015a6e9c80b10e45999f2b71d131a39d48247b0978333008ee `25d25735556a2015a6e9c80b10e45999f2b71d131a39d48247b0978333008ee
	//			hashKey = a61b62264d261849adb0b55eb6efdac50540b1621169f044b9763c3b1e347765 480f2dfaf332de2175b7475651ab63402543d6b342220b882fd924c10ddd25b3
	//			root = 12157598575d967a5127f8f8056f55e6651696773e6f3863727f0ca9a5bdbac3

	/*
		keyword := []byte("lwjsz") //42cb8f26e90ca25c54da1489eed90e84fd254fee470564640fc506df057b8b75
		hash := VedCrypto.Hash(keyword)
		hashStr := string(hash)
		// byteStr := []byte(hashStr)
		// print("keyword hash: ", hash, " ", hashStr, " ", byteStr, " \n")

		finalChar := hashStr[len(hashStr)-1:]
		//print("the final character is: ", finalChar, "\n")
		ciphertext, _ := VedCrypto.Encrypt([]byte(finalChar), key0)
		// print("start generate proof...\n")
		inexistProof := VedCrypto.Hash(ciphertext, ciphertext)
		// print("start update index proof...\n")
		HSMT, _ = core.UpdateIndex(hash, HSMT, inexistProof, newIndex)
		print("root after 1st update: ", string(HSMT.Hash), "\n")
		HSMT, _ = core.UpdateIndex(hash, HSMT, inexistProof, newIndex)
		print("root after 2nd update: ", string(HSMT.Hash), "\n")
		proofs, leaf, _ := core.GetProof(hash, HSMT)
		res := core.VerifyProof(proofs, HSMT, leaf)
		print(res)

	*/
}
