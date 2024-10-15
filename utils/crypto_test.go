package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Create a temporary PEM file with a P-256 private key for testing.
func createTestPEMFile(t *testing.T) string {
	// Generate a new P-256 ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Convert the private key to ASN.1 DER encoded form
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal EC private key: %v", err)
	}

	// Create a PEM block with the private key
	pemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Write the PEM block to a temporary file
	tempFile, err := os.CreateTemp("", "test_key_*.pem")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tempFile.Close()

	if err := pem.Encode(tempFile, pemBlock); err != nil {
		t.Fatalf("Failed to encode PEM block: %v", err)
	}

	return tempFile.Name()
}

func TestLoadPrivateKey(t *testing.T) {
	// Create a temporary PEM file containing a P-256 private key
	pemFile := createTestPEMFile(t)
	defer os.Remove(pemFile) // Clean up the temporary file after the test

	// Load the private key from the PEM file
	privateKey, publicKey, err := LoadPrivateKey(pemFile)
	if err != nil {
		t.Fatalf("Failed to load private key: %v", err)
	}

	// Validate that the private key is not nil
	if privateKey == nil {
		t.Fatalf("Private key is nil")
	}

	// Validate that the public key is not nil
	if publicKey == nil {
		t.Fatalf("Public key is nil")
	}

	// Validate that the private key and public key are using the P-256 curve
	if privateKey.Curve != elliptic.P256() {
		t.Fatalf("Private key is not using P-256 curve")
	}
	if publicKey.Curve != elliptic.P256() {
		t.Fatalf("Public key is not using P-256 curve")
	}

	// Validate that the private and public keys are consistent
	if !privateKey.PublicKey.Equal(publicKey) {
		t.Fatalf("Private and public keys do not match")
	}

	t.Logf("Private and public keys loaded successfully")
}

// Test encryption and decryption
func TestEncryptDecryptData(t *testing.T) {
	// Generate key pairs for sender and recipient
	senderPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "Error generating sender's private key")

	recipientPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "Error generating recipient's private key")

	// Prepare a message to encrypt
	plaintext := []byte("This is a secret message")

	senderEcdhPrivKey, err := EcdsaPrivToEcdh(senderPrivKey)
	require.NoError(t, err, "Error during conversion from ecdsa to ecdh private key")

	senderEcdhPubKey, err := EcdsaPubToEcdh(&senderPrivKey.PublicKey)
	require.NoError(t, err, "Error during conversion from ecdsa to ecdh public key")

	recipientEcdhPrivKey, err := EcdsaPrivToEcdh(recipientPrivKey)
	require.NoError(t, err, "Error during conversion from ecdsa to ecdh private key")

	recipientEcdhPubKey, err := EcdsaPubToEcdh(&recipientPrivKey.PublicKey)
	require.NoError(t, err, "Error during conversion from ecdsa to ecdh public key")

	// Encrypt the message using the sender's private key and recipient's public key
	ciphertext, nonce, err := EncryptData(senderEcdhPrivKey, recipientEcdhPubKey, plaintext)
	require.NoError(t, err, "Error during encryption")

	// Decrypt the message using the recipient's private key and sender's public key
	decryptedMessage, err := DecryptData(recipientEcdhPrivKey, senderEcdhPubKey, ciphertext, nonce)
	require.NoError(t, err, "Error during decryption")

	// Ensure the decrypted message matches the original plaintext
	require.Equal(t, plaintext, decryptedMessage, "Decrypted message does not match the original")
}

// Test signing and verifying a transaction
func TestSignVerifyTransaction(t *testing.T) {
	// Generate a key pair for the sender
	senderPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "Error generating sender's private key")

	// Prepare a sample transaction (this could be any data, such as transaction data)
	transactionData := []byte("Transaction data to sign")

	// Sign the transaction
	rBytes, sBytes, err := SignTransaction(senderPrivKey, transactionData)
	require.NoError(t, err, "Error during transaction signing")

	// Verify the signature using the sender's public key
	isValid := VerifySignature(&senderPrivKey.PublicKey, transactionData, rBytes, sBytes)
	require.True(t, isValid, "Signature verification failed")
}
