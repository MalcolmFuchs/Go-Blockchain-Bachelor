package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

// Verschlüsselt die Patientendaten mit AES-GCM
func EncryptData(aesKey, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Verschlüssele den AES-Schlüssel mit dem Public Key des Arztes
func EncryptAESKeyWithDoctorKey(doctorPublicKey ed25519.PublicKey, aesKey []byte, patientPrivateKey ed25519.PrivateKey) ([]byte, error) {
	if len(doctorPublicKey) != ed25519.PublicKeySize || len(patientPrivateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid key size")
	}

	// Berechne den Shared Secret Key mit X25519
	var sharedSecret [32]byte
	_, err := curve25519.X25519(patientPrivateKey[:32], doctorPublicKey)
	if err != nil {
		return nil, err
	}

	// Verwende HKDF, um den shared secret in einen Verschlüsselungsschlüssel zu verwandeln
	hash := sha256.New
	hkdf := hkdf.New(hash, sharedSecret[:], nil, nil)

	// Erstelle einen AES-Schlüssel aus dem Shared Secret
	aesSharedKey := make([]byte, 32) // 256-bit AES-Schlüssel
	if _, err := io.ReadFull(hkdf, aesSharedKey); err != nil {
		return nil, err
	}

	// Verschlüssele den AES-Schlüssel mit AES-GCM
	encryptedAESKey, err := EncryptData(aesSharedKey, aesKey)
	if err != nil {
		return nil, err
	}

	return encryptedAESKey, nil
}

func DecodePublicKey(encodedKey string) (ed25519.PublicKey, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, errors.New("invalid base64 encoding")
	}
	if len(decodedKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}
	return ed25519.PublicKey(decodedKey), nil
}
