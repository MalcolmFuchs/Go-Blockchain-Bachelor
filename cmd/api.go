package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

func (node *Node) GetBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	blockchainData, err := node.Blockchain.GetBlockchainData()
	if err != nil {
		http.Error(w, "failed to get blockchain data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blockchainData)
}

func (authorityNode *AuthorityNode) AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction blockchain.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "failed to decode transaction", http.StatusBadRequest)
		return
	}

	signature, err := blockchain.SignTransaction(&transaction, authorityNode.PrivateKey)
	if err != nil {
		http.Error(w, "failed to sign transaction", http.StatusInternalServerError)
		return
	}

	fmt.Println(signature)

	authorityNode.PendingTransactions = append(authorityNode.PendingTransactions, &transaction)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction added successfully"))
}

func (authorityNode *AuthorityNode) CreateBlockHandler(w http.ResponseWriter, r *http.Request) {
	block, err := authorityNode.CreateBlock()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create block: %v", err), http.StatusInternalServerError)
		return
	}

	blockData, err := json.Marshal(block)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to serialize block: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blockData)
}

func (authorityNode *AuthorityNode) SyncHandler(w http.ResponseWriter, r *http.Request) {
	var syncRequest SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&syncRequest); err != nil {
		http.Error(w, "failed to decode sync request", http.StatusBadRequest)
		return
	}

	var syncBlocks []*blockchain.Block
	syncStartIndex := -1

	for i, block := range authorityNode.Blockchain.Blocks {
		if fmt.Sprintf("%x", block.Hash) == syncRequest.LastBlockHash {
			syncStartIndex = i + 1
			break
		}
	}

	if syncStartIndex != -1 && syncStartIndex < len(authorityNode.Blockchain.Blocks) {
		syncBlocks = authorityNode.Blockchain.Blocks[syncStartIndex:]
	}

	syncResponse := SyncResponse{Blocks: syncBlocks}
	responseBody, err := json.Marshal(syncResponse)
	if err != nil {
		http.Error(w, "failed to serialize sync response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func (authorityNode *AuthorityNode) SetupAuthorityNodeRoutes() {
  authorityNode.SetupNodeRoutes()
	http.HandleFunc("/addTransaction", authorityNode.AddTransactionHandler)
	http.HandleFunc("/createBlock", authorityNode.CreateBlockHandler)
	http.HandleFunc("/sync", authorityNode.SyncHandler)
}

func (node *Node) SetupNodeRoutes() {
	http.HandleFunc("/getBlockchain", node.GetBlockchainHandler)
}
