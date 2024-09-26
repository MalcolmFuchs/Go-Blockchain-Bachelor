package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"

	blockchain "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/block"
)

// Generiere einen neuen AES-Schlüssel für die Verschlüsselung der Patientendaten
func GenerateAESKey() []byte {
	key := make([]byte, 32) // 256-Bit AES-Schlüssel
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	return key
}

// Signiere Transaktionen oder Blöcke mit ed25519
func SignTransaction(privateKey ed25519.PrivateKey, data []byte) []byte {
	return ed25519.Sign(privateKey, data)
}

func VerifySignature(publicKey ed25519.PublicKey, data, signature []byte) bool {
	return ed25519.Verify(publicKey, data, signature)
}

func GetPatientPrivateKey() ed25519.PrivateKey {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic("Error generating patient private key: " + err.Error())
	}
	return privateKey
}

func GetPatientPublicKey(privateKey ed25519.PrivateKey) ed25519.PublicKey {
	return privateKey.Public().(ed25519.PublicKey)
}

// Generiere und erhalte einen Arzt-PublicKey
func GetDoctorPublicKey() ed25519.PublicKey {
	// Generiere einen neuen ed25519 Schlüsselpaar
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Error generating doctor public key:", err)
		return nil
	}

	return publicKey
}

func SignBlock(block *blockchain.Block, authorityPrivateKey ed25519.PrivateKey) {
	block.AuthoritySign = ed25519.Sign(authorityPrivateKey, []byte(block.Hash))
}
