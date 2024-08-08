package components

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateBlockchain() *Blockchain {
	blockchain := &Blockchain{
		Blocks:        []Block{},
		Nodes:         []*AuthorityNode{},
		Patients:      make(map[string]PersonalData),
		LastNodeIndex: -1,
		TxChan:        make(chan struct{}),
	}
	blockchain.CreateGenesisBlock()

	return blockchain
}

func (bc *Blockchain) validateAndAddBlock(newBlock Block, node *AuthorityNode) {
	if len(bc.Blocks) > 0 {
		lastBlock := bc.Blocks[len(bc.Blocks)-1]
		if newBlock.PrevHash != lastBlock.Hash {
			fmt.Println("Invalid block: PrevHash does not match")
			return
		}
	}

	newBlock.Hash = newBlock.calculateHash()
	dataToSign := fmt.Sprintf("%d%s%s%s", newBlock.Index, newBlock.Timestamp, newBlock.PatientData.PersonalData.InsuranceNumber, newBlock.PrevHash)
	r, s := SignData(dataToSign, node.PrivateKey)
	newBlock.SignatureR = r
	newBlock.SignatureS = s

	if VerifySignature(dataToSign, newBlock.SignatureR, newBlock.SignatureS, node.PublicKey) {
		bc.addBlock(newBlock)
	} else {
		fmt.Println("Invalid signature. Block not added.")
	}
}

func (bc *Blockchain) GetNextNode() (*AuthorityNode, error) {
	if len(bc.Nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	bc.LastNodeIndex = (bc.LastNodeIndex + 1) % len(bc.Nodes)
	return bc.Nodes[bc.LastNodeIndex], nil
}
