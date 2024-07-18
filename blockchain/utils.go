package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func HashInsuranceNumber(insuranceNumber string) string {
	hasher := sha256.New()
	hasher.Write([]byte(insuranceNumber))
	hashed := hasher.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateKeyPair() (*ecdsa.PrivateKey, ecdsa.PublicKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
	}
	return privateKey, privateKey.PublicKey
}

func SignData(data string, privateKey *ecdsa.PrivateKey) (string, string) {
	hash := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		fmt.Println("Error signing data:", err)
		return "", ""
	}
	return r.String(), s.String()
}

func VeryfiySignature(data, rText, sText string, publicKey ecdsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(data))
	r := new(big.Int)
	s := new(big.Int)
	r.SetString(rText, 10)
	s.SetString(sText, 10)

	return ecdsa.Verify(&publicKey, hash[:], r, s)
}
