package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
)

func (node *Node) GetBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	// Stelle sicher, dass die Blockchain vorhanden ist
	if node.Blockchain == nil {
		http.Error(w, "Blockchain not initialized", http.StatusInternalServerError)
		return
	}

	// Serialisiere die Blockchain
	blockchainData, err := json.MarshalIndent(node.Blockchain.Blocks, "", "  ")
	if err != nil {
		http.Error(w, "Failed to serialize blockchain", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blockchainData)
}

func (authorityNode *AuthorityNode) AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transactionData struct {
		Type    string `json:"type"`
		Notes   string `json:"notes"`
		Results string `json:"results"`
		Doctor  string `json:"doctor"`
		Patient string `json:"patient"`
		Key     string `json:"key"`
	}

	// Dekodiere die Transaktionsdaten aus der Anfrage
	err := json.NewDecoder(r.Body).Decode(&transactionData)
	if err != nil {
		http.Error(w, "failed to decode transaction data", http.StatusBadRequest)
		return
	}

	// Erstelle eine neue Transaktion aus den empfangenen Daten
	transaction, err := blockchain.NewTransaction(
		transactionData.Type,
		transactionData.Notes,
		transactionData.Results,
		transactionData.Doctor,
		transactionData.Patient,
		transactionData.Key,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Signiere die Transaktion mit dem Private Key des Authority Nodes
	signature, err := blockchain.SignTransaction(transaction, authorityNode.PrivateKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign transaction: %v", err), http.StatusInternalServerError)
		return
	}
	transaction.Signature = signature

	// Füge die Transaktion zum Transaktionspool hinzu
	// if err := authorityNode.TransactionPool.AddTransactionToPool(transaction, authorityNode.PublicKey); err != nil {
	// 	http.Error(w, fmt.Sprintf("failed to add transaction to pool: %v", err), http.StatusInternalServerError)
	// 	return
	// }

	if err := authorityNode.AddTransaction(transaction); err != nil {
		http.Error(w, fmt.Sprintf("failed to add transaction to pool: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction added to pool successfully"))
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
	if syncRequest.LastBlockHash == "" {
		// Client hat keine Blöcke, sende die gesamte Blockchain
		syncBlocks = authorityNode.Blockchain.Blocks
	} else {
		// Suche nach dem Block mit dem gegebenen Hash
		syncStartIndex := -1
		for i, block := range authorityNode.Blockchain.Blocks {
			if fmt.Sprintf("%x", block.Hash) == syncRequest.LastBlockHash {
				syncStartIndex = i + 1
				break
			}
		}
		if syncStartIndex == -1 {
			http.Error(w, "block not found", http.StatusNotFound)
			return
		}
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

func (a *AuthorityNode) GetPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	publicKeyHex := hex.EncodeToString(a.PublicKey)
	response := map[string]string{
		"publicKey": publicKeyHex,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (authorityNode *AuthorityNode) GetTransactionPoolHandler(w http.ResponseWriter, r *http.Request) {
	// Sperre den Zugriff auf den Pool
	authorityNode.mutex.Lock()
	defer authorityNode.mutex.Unlock()

	// Serialisiere den Transaktionspool
	poolData, err := json.Marshal(authorityNode.TransactionPool.Transactions)
	if err != nil {
		http.Error(w, "failed to serialize transaction pool", http.StatusInternalServerError)
		return
	}

	// Gib den Transaktionspool als JSON zurück
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(poolData)
}

func (a *AuthorityNode) SetupAuthorityNodeRoutes() {
	a.SetupNodeRoutes()
	http.HandleFunc("/addTransaction", a.AddTransactionHandler)
	http.HandleFunc("/createBlock", a.CreateBlockHandler)
	http.HandleFunc("/sync", a.SyncHandler)
	http.HandleFunc("/getTransactionPool", a.GetTransactionPoolHandler)
	http.HandleFunc("/getPublicKey", a.GetPublicKeyHandler)
}

func (node *Node) SetupNodeRoutes() {
	http.HandleFunc("/getBlockchain", node.GetBlockchainHandler)
}
