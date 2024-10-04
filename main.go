package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/cmd"
)

func main() {
	// Starte das Cobra-Root-Kommando

	dPub, _, _ := ed25519.GenerateKey(nil)

	hexDPub := hex.EncodeToString(dPub)

	fmt.Println(hexDPub)

	cmd.Execute()
}
