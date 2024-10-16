package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"

	"golang.org/x/crypto/hkdf"
)

type EncryptedData struct {
	Ciphertext []byte `json:"ciphertext"`
	Nonce      []byte `json:"nonce"`
}

func LoadPrivateKey(filename string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// Read the private key PEM file
	pemData, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Parse the ECDSA private key
	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse EC private key: %v", err)
	}

	// The public key is embedded in the private key
	pubKey := &privKey.PublicKey

	return privKey, pubKey, nil
}

// LoadPublicKey reads an ECDSA public key from a PEM file
func LoadPublicKey(filename string) (*ecdsa.PublicKey, error) {
	// Read the public key PEM file
	pemData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Assert the type of public key to be ECDSA
	pubKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return pubKey, nil
}

func EncryptData(senderPrivKey *ecdh.PrivateKey, recipientPubKey *ecdh.PublicKey, plaintext []byte) ([]byte, []byte, error) {
	// Perform ECDH key exchange to derive the shared secret
	sharedSecret, err := senderPrivKey.ECDH(recipientPubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("ECDH key exchange failed: %v", err)
	}

	// Derive symmetric key using HKDF
	salt := []byte("ECDH encryption")
	info := []byte("encryption key")
	hkdf := hkdf.New(sha256.New, sharedSecret, salt, info)
	symmetricKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, symmetricKey); err != nil {
		return nil, nil, err
	}

	// Encrypt data using AES-GCM
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func DecryptData(recipientPrivKey *ecdh.PrivateKey, senderPubKey *ecdh.PublicKey, ciphertext, nonce []byte) ([]byte, error) {
	// Perform ECDH key exchange to derive the shared secret
	sharedSecret, err := recipientPrivKey.ECDH(senderPubKey)
	if err != nil {
		return nil, fmt.Errorf("ECDH key exchange failed: %v", err)
	}

	// Derive symmetric key using HKDF
	salt := []byte("ECDH encryption")
	info := []byte("encryption key")
	hkdf := hkdf.New(sha256.New, sharedSecret, salt, info)
	symmetricKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, symmetricKey); err != nil {
		return nil, err
	}

	// Decrypt data using AES-GCM
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func SignTransaction(senderPrivKey *ecdsa.PrivateKey, transactionData []byte) ([]byte, []byte, error) {
	hash := sha256.Sum256(transactionData)
	r, s, err := ecdsa.Sign(rand.Reader, senderPrivKey, hash[:])
	if err != nil {
		return nil, nil, err
	}
	return r.Bytes(), s.Bytes(), nil
}

func VerifySignature(senderPubKey *ecdsa.PublicKey, transactionData, rBytes, sBytes []byte) bool {
	hash := sha256.Sum256(transactionData)
	var r, s big.Int
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)
	return ecdsa.Verify(senderPubKey, hash[:], &r, &s)
}

func EcdsaPrivToEcdh(ecdsaPrivKey *ecdsa.PrivateKey) (*ecdh.PrivateKey, error) {
	ecdhCurve := ecdh.P256()
	ecdhPrivKey, err := ecdhCurve.NewPrivateKey(ecdsaPrivKey.D.Bytes())
	if err != nil {
		fmt.Println("Error converting ECDSA private key to ECDH private key:", err)
		return nil, err
	}

	return ecdhPrivKey, nil
}

func EcdsaPubToEcdh(ecdsaPubKey *ecdsa.PublicKey) (*ecdh.PublicKey, error) {
	ecdhCurve := ecdh.P256()
	ecdhPubKey, err := ecdhCurve.NewPublicKey(SerializePublicKey(ecdsaPubKey))
	if err != nil {
		fmt.Println("Error converting ECDSA public key to ECDH public key:", err)
		return nil, err
	}

	return ecdhPubKey, nil
}

// Serialize the ECDSA public key in uncompressed form (X and Y coordinates concatenated)
func SerializePublicKey(pubKey *ecdsa.PublicKey) []byte {
	// Uncompressed public key format: 0x04 || X || Y
	return append([]byte{0x04}, append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)...)
}

// Deserialize a public key from uncompressed bytes and return an *ecdsa.PublicKey
func DeserializePublicKey(pubKeyBytes []byte) (*ecdsa.PublicKey, error) {
	// Step 1: Ensure the key is uncompressed and has a 0x04 prefix
	if len(pubKeyBytes) == 0 || pubKeyBytes[0] != 0x04 {
		return nil, fmt.Errorf("invalid uncompressed public key format")
	}

	// Step 2: Extract the X and Y coordinates from the byte slice
	curve := elliptic.P256()
	coordinateLength := (curve.Params().BitSize + 7) / 8 // Number of bytes per coordinate

	expectedLength := 1 + 2*coordinateLength // 1 byte for prefix, X and Y coordinates
	if len(pubKeyBytes) != expectedLength {
		return nil, fmt.Errorf("invalid public key length")
	}

	xBytes := pubKeyBytes[1 : 1+coordinateLength]
	yBytes := pubKeyBytes[1+coordinateLength:]

	// Step 3: Convert the coordinate bytes to big.Int
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	// Step 4: Create and return the ECDSA public key
	pubKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	return pubKey, nil
}
