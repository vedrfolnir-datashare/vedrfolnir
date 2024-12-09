package VedCrypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(ms ...[]byte) []byte {
	h := sha256.New()
	for _, m := range ms {
		h.Write(m)
	}
	//return h.Sum(nil)
	hash := h.Sum(nil)
	hashString := hex.EncodeToString(hash[:])
	return []byte(hashString)
}
