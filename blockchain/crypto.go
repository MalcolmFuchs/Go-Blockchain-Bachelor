package blockchain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func Encrypt(data string, passphrase string) string {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return ""
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM:", err)
		return ""
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println("Error reading random data:", err)
		return ""
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return hex.EncodeToString(ciphertext)
}

func Decrypt(data string, passphrase string) string {
	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return ""
	}
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return ""
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM:", err)
		return ""
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return ""
	}

	return string(plaintext)
}
