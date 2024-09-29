package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

func SendTransaction(transaction *blockchain.Transaction, authorityNodeAddress string) error {
	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to serialize transaction: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/receiveTransaction", authorityNodeAddress), "application/json", bytes.NewBuffer(transactionBytes))
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %s", string(body))
	}

	fmt.Printf("Transaction sent successfully to %s\n", authorityNodeAddress)
	return nil
}

func ReceiveTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction blockchain.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "failed to decode transaction", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received transaction with hash %x\n", transaction.Hash)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction received successfully"))
}

func RequestBlocks(startBlockID uint64, authorityNodeAddress string) ([]*blockchain.Block, error) {
	url := fmt.Sprintf("http://%s/getBlocks?startBlockID=%d", authorityNodeAddress, startBlockID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request blocks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-OK response: %s", string(body))
	}

	var blocks []*blockchain.Block
	err = json.NewDecoder(resp.Body).Decode(&blocks)
	if err != nil {
		return nil, fmt.Errorf("failed to decode blocks: %v", err)
	}

	return blocks, nil
}

// SyncBlockchain synchronizes the local blockchain with the authority node's blockchain
func (n *Node) SyncBlockchain(authorityNodeAddress string) error {
	lastBlockID := uint64(0)
	if len(n.Blockchain.Blocks) > 0 {
		lastBlockID = n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1].ID
	}

	newBlocks, err := RequestBlocks(lastBlockID+1, authorityNodeAddress)
	if err != nil {
		return fmt.Errorf("failed to sync blockchain: %v", err)
	}

	// 3. Füge die neuen Blöcke zur lokalen Blockchain hinzu
	for _, block := range newBlocks {
		if err := n.Blockchain.AddBlock(block, n.TrustedPublicKey); err != nil {
			return fmt.Errorf("failed to add block: %v", err)
		}
		fmt.Printf("Added new block with ID %d and hash %x\n", block.ID, block.Hash)
	}

	return nil
}
