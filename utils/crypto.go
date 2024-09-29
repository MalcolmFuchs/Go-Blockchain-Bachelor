package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate keys: %v", err)
	}
	return publicKey, privateKey, nil
}

func EncryptData(plaintext, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GCM mode: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func DecryptData(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %v", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plaintext, nil
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

// EncryptAESKeyWithPublicKey encrypts an AES key using the recipient's public key via ECDH
func EncryptAESKeyWithPublicKey(aesKey []byte, publicKey *ecdsa.PublicKey) ([]byte, error) {
    // Create an ECDH curve object
    curve := ecdh.P256()

    // Serialize the ECDSA public key into a single byte slice (concatenating X and Y coordinates)
    recipientPubKeyBytes := serializePublicKey(publicKey)

    // Convert the serialized ECDSA public key to the ECDH public key format
    recipientPubKey, err := curve.NewPublicKey(recipientPubKeyBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to convert public key: %v", err)
    }

    // Generate an ephemeral private key for ECDH
    ephemeralPrivKey, err := curve.GenerateKey(rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("failed to generate ephemeral private key: %v", err)
    }

    // Derive a shared secret using ECDH with the recipient's public key and ephemeral private key
    sharedSecret, err := ephemeralPrivKey.ECDH(recipientPubKey)
    if err != nil {
        return nil, fmt.Errorf("failed to derive shared secret: %v", err)
    }

    // Hash the shared secret using SHA-256 to produce the AES key
    sharedSecretHash := sha256.Sum256(sharedSecret)

    // Use the shared secret hash as the AES key
    block, err := aes.NewCipher(sharedSecretHash[:])
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES-GCM: %v", err)
    }

    nonce := make([]byte, aesGCM.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %v", err)
    }

    // Encrypt the AES key
    encryptedAESKey := aesGCM.Seal(nil, nonce, aesKey, nil)

    // Include the ephemeral public key (in bytes) for the recipient to derive the shared secret later
    ephemeralPubKeyBytes := ephemeralPrivKey.PublicKey().Bytes()
    encryptedData := append(ephemeralPubKeyBytes, nonce...)
    encryptedData = append(encryptedData, encryptedAESKey...)

    return encryptedData, nil
}
