package blockchain

import (
	"time"

	patient "github.com/MalcolmFuchs/Go-Blockchain-Bachelor/build/patient"
)

type Block struct {
	Index        int
	Timestamp    time.Time
	Hash         string
	Transactions []Transaction
	PrevHash     string
}

type Blockchain struct {
	GenesisBlock   Block
	Chain          []Block
	AuthorityNodes []AuthorityNode
}

type Transaction struct {
	PatientID string
	Record    patient.PatientRecord
}

type AuthorityNode struct {
	Id         string
	PublicKey  string
	PrivateKey string
}
