package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// help us generate n random bytes or return error if was one
// uses crypto/rand to safe to use with remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// will generate a byte slice of size nBytes and then
// return a string that is the base64 URL encoded version
// of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		// can't return nil first cause string is not a pointer
		return "", err
	}
	// used to convert byte to string. cannot use native String method from strings package as not all byte's can be converted into strings with it
	// is not encryption method. just encoding / decoding
	return base64.URLEncoding.EncodeToString(b), nil
}

// helper function designed to generate remember tokens
// of a predetemined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
