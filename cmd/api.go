package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

type API struct {
	Node *Node
}

func NewAPI(node *Node) *API {
	return &API{Node: node}
}

func (api *API) GetBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	blockchainData, err := api.Node.Blockchain.GetBlockchainData()
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

	err = blockchain.SignTransaction(&transaction, api.Node.PrivateKey)
	if err != nil {
		http.Error(w, "failed to sign transaction", http.StatusInternalServerError)
		return
	}

	api.Node.Blockchain.Blocks = append(api.Node.Blockchain.Blocks, &blockchain.Block{
		ID:           uint64(len(api.Node.Blockchain.Blocks) + 1),
		Transactions: []*blockchain.Transaction{&transaction},
		Timestamp:    transaction.Timestamp,
	})

	fmt.Println("Added transaction to blockchain with ID:", transaction.Hash)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction added successfully"))
}

func (api *API) CreateBlockHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Rufe die CreateBlock()-Funktion des Authority Nodes auf
	block, err := api.AuthorityNode.CreateBlock()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create block: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Sende den erstellten Block als JSON-Antwort zurück
	blockData, err := json.Marshal(block)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to serialize block: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blockData)
}

// SetupRoutes sets up the RESTful routes for the API
func (api *API) SetupRoutes() {
	http.HandleFunc("/getBlockchain", api.GetBlockchainHandler)
	http.HandleFunc("/addTransaction", api.AddTransactionHandler)
	http.HandleFunc("/createBlock", api.CreateBlockHandler) // Neuer Endpunkt zum Erstellen eines Blocks
}
