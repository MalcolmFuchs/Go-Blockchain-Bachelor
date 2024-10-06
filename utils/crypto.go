package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

type EncryptedData struct {
	Ciphertext []byte `json:"ciphertext"`
	Nonce      []byte `json:"nonce"`
}

func EncryptData(plaintext []byte, key []byte) (EncryptedData, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptedData{}, fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptedData{}, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return EncryptedData{}, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	return EncryptedData{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	}, nil
}

func DecryptData(encrypted EncryptedData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	plaintext, err := aesGCM.Open(nil, encrypted.Nonce, encrypted.Ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plaintext, nil
}

func EncryptAESKeyWithPublicKey(aesKey []byte, patientPublicKey ed25519.PublicKey) ([]byte, error) {
	// Konvertiere Ed25519-Public-Key zu X25519-Public-Key
	x25519PubKey := ed25519PublicKeyToX25519(patientPublicKey)

	// Generiere einen ECDH-Privat-Key
	var ecdhPrivKey [32]byte
	if _, err := rand.Read(ecdhPrivKey[:]); err != nil {
		return nil, fmt.Errorf("failed to generate ECDH private key: %v", err)
	}

	// Berechne das gemeinsame Geheimnis
	sharedSecret, err := curve25519.X25519(ecdhPrivKey[:], x25519PubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute shared secret: %v", err)
	}

	// Hash des gemeinsamen Geheimnisses verwenden als AES-Schlüssel
	hashedSecret := sha256.Sum256(sharedSecret)

	// Verschlüssele den AES-Schlüssel mit dem gehashten gemeinsamen Geheimnis
	encryptedAESKey, err := EncryptData(aesKey, hashedSecret[:])
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt AES key: %v", err)
	}

	// In diesem Beispiel kombinieren wir den ECDH-Privat-Key und den verschlüsselten AES-Schlüssel
	// In der Praxis solltest du den ECDH-Privat-Key nicht übertragen
	return append(ecdhPrivKey[:], encryptedAESKey.Ciphertext...), nil
}

func ed25519PublicKeyToX25519(edPubKey ed25519.PublicKey) []byte {
	// Dies ist eine vereinfachte Konvertierung. In der Praxis solltest du einen korrekten Mechanismus verwenden.
	return edPubKey[:32]
}

func GenerateRandomHexKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keys: %v", err)
	}
	return publicKey, privateKey, nil
}

func GenerateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %v", err)
	}

	return randomBytes, nil
}

func GenerateECDSAKeys() (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key pair: %v", err)
	}
	return privateKey, nil
}

// Helper function to serialize a public key by concatenating X and Y coordinates
func serializePublicKey(publicKey *ecdsa.PublicKey) []byte {
	xBytes := publicKey.X.Bytes()
	yBytes := publicKey.Y.Bytes()
	return append(xBytes, yBytes...)
}

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256 Bits
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}
