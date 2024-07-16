package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
)

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

func (hr HealthRecord) toJSON() string {
	jsonData, err := json.Marshal(hr)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
	}

	return string(jsonData)
}

func fromJSON(data string) HealthRecord {
	var hr HealthRecord
	err := json.Unmarshal([]byte(data), &hr)
	if err != nil {
		fmt.Println("Error unmarshalling from JSON:", err)
	}
	return hr
}
