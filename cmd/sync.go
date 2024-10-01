package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

type SyncRequest struct {
	LastBlockHash string `json:"lastBlockHash"`
}

type SyncResponse struct {
	Blocks []*blockchain.Block `json:"blocks"`
}

func (n *Node) SyncWithAuthorityNode(authorityNodeAddress string) error {

	var lastBlockHash string
	if len(n.Blockchain.Blocks) > 0 {
		lastBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]
		lastBlockHash = fmt.Sprintf("%x", lastBlock.Hash)
	} else {
		fmt.Println("Blockchain is empty")
		lastBlockHash = ""
	}

	syncRequest := SyncRequest{LastBlockHash: lastBlockHash}
	requestBody, err := json.Marshal(syncRequest)
	if err != nil {
		return fmt.Errorf("failed to serialize sync request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/sync", authorityNodeAddress), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to send sync request: %v", err)
	}
	defer resp.Body.Close()

	var syncResponse SyncResponse
	if err := json.NewDecoder(resp.Body).Decode(&syncResponse); err != nil {
		return fmt.Errorf("failed to decode sync response: %v", err)
	}

	n.Blockchain.Blocks = append(n.Blockchain.Blocks, syncResponse.Blocks...)

	for _, block := range syncResponse.Blocks {
		n.Blockchain.BlockMap[fmt.Sprintf("%x", block.Hash)] = block
	}

	return nil
}
