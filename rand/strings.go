package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes number of bytes
const RememberTokenBytes = 32

// Bytes returns a slice of bytes from an int param
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byte slice of size nBytes & then return a string that is the
// base64 encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function designed to generate remember tokens of a predetermined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
