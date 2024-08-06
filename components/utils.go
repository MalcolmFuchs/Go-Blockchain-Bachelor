package components

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

func (b *Block) calculateHash() string {
	timestampString := b.Timestamp.Format(time.RFC3339)
	record := strconv.Itoa(b.Index) + timestampString + b.PatientData.PersonalData.InsuranceNumber + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func (bc *Blockchain) calculatBcHash(block Block) string {
	record := fmt.Sprintf("%d%s%s", block.Index, block.Timestamp, block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

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

func VerifySignature(data, rText, sText string, publicKey ecdsa.PublicKey) bool {
	hash := sha256.Sum256([]byte(data))
	r := new(big.Int)
	s := new(big.Int)
	r.SetString(rText, 10)
	s.SetString(sText, 10)

	return ecdsa.Verify(&publicKey, hash[:], r, s)
}
