package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

type API struct {
	AuthorityNode *AuthorityNode
}

func NewAPI(authorityNode *AuthorityNode) *API {
	return &API{AuthorityNode: authorityNode}
}

func (api *API) GetBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	blockchainData, err := api.AuthorityNode.Blockchain.GetBlockchainData()
	if err != nil {
		http.Error(w, "failed to get blockchain data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blockchainData)
}

func (api *API) AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction blockchain.Transaction

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "failed to decode transaction", http.StatusBadRequest)
		return
	}

	err = blockchain.SignTransaction(&transaction, api.AuthorityNode.Node.PrivateKey)
	if err != nil {
		http.Error(w, "failed to sign transaction", http.StatusInternalServerError)
		return
	}

	api.AuthorityNode.PendingTransactions = append(api.AuthorityNode.PendingTransactions, &transaction)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction added successfully"))
}

func (api *API) CreateBlockHandler(w http.ResponseWriter, r *http.Request) {
	block, err := api.AuthorityNode.CreateBlock()
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

func (api *API) SyncHandler(w http.ResponseWriter, r *http.Request) {
	var syncRequest SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&syncRequest); err != nil {
		http.Error(w, "failed to decode sync request", http.StatusBadRequest)
		return
	}

	var syncBlocks []*blockchain.Block
	syncStartIndex := -1

	for i, block := range api.AuthorityNode.Blockchain.Blocks {
		if fmt.Sprintf("%x", block.Hash) == syncRequest.LastBlockHash {
			syncStartIndex = i + 1
			break
		}
	}

	if syncStartIndex != -1 && syncStartIndex < len(api.AuthorityNode.Blockchain.Blocks) {
		syncBlocks = api.AuthorityNode.Blockchain.Blocks[syncStartIndex:]
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

func (api *API) SetupRoutes() {
	http.HandleFunc("/getBlockchain", api.GetBlockchainHandler)
	http.HandleFunc("/addTransaction", api.AddTransactionHandler)
	http.HandleFunc("/createBlock", api.CreateBlockHandler)
	http.HandleFunc("/sync", api.SyncHandler)
}
