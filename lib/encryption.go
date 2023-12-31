package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
)

// extendPassword will extend the password to the required length
func extendPassword(pass string) string {
	for len(pass) < 32 {
		for _, p := range pass {
			pass = pass + string(p)
			if len(pass) == 32 {
				break
			}
		}
	}
	return pass
}

// NewGCM will generate a GCM given a password
func NewGCM(password string) (cipher.AEAD, error) {
	block, err := aes.NewCipher([]byte(password))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// ReadCypherText will read and decrypt a cyphertext
func ReadCypherText(gcm cipher.AEAD, path string) ([]byte, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	return data[gcm.NonceSize():], data[:gcm.NonceSize()], nil
}
