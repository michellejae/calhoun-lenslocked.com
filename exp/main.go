package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"gitlab.com/michellejae/lenslocked.com/hash"
)

func main() {
	toHash := []byte("this is my string to hash")
	h := hmac.New(sha256.New, []byte("my-secrett"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(base64.URLEncoding.EncodeToString(b))

	hmac := hash.NewHmac("my-secrett")
	fmt.Println(hmac.Hash("this is my string to hash"))
}
