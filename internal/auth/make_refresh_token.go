package auth 

import (
	"crypto/rand"
	"encoding/hex"
)

// generate a random 256-bit hex-encoded string
func MakeRefreshToken() (string, error) {
	// generate a random 256 bits of random data
	data := make([]byte, 32)
	rand.Read(data)		// fill values into 'data'
	

	// convert the byte slice into hex string
	hexString := hex.EncodeToString(data)
	return hexString, nil
}