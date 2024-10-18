package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func createHash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	val := hex.EncodeToString(hash)
	hasher.Reset()
	fmt.Println(val)
	return val
}
func verifHash(val string, enc string) {
	if createHash(val) == enc {
		fmt.Printf("son iguales")
	} else {
		fmt.Printf("no son iguales")
	}
}
func main() {
	enc := createHash("holamama")
	verifHash("holama", enc)
}
