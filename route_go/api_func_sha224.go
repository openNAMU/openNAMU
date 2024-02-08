package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	call_arg := os.Args[1:]

	data := call_arg[0]

	hasher := sha256.New224()
	hasher.Write([]byte(data))
	hash_byte := hasher.Sum(nil)
	hash_str := hex.EncodeToString(hash_byte)

	fmt.Print(hash_str)
}
