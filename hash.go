package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func createHash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	val := hex.EncodeToString(hash)
	hasher.Reset()
	return val
}
