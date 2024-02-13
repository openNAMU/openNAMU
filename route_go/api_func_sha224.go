package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

	new_data := map[string]string{}
	new_data["data"] = hash_str

	json_data, _ := json.Marshal(new_data)
	fmt.Print(string(json_data))
}
