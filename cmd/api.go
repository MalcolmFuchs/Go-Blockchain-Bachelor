package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/blockchain"
	"github.com/MalcolmFuchs/Go-Blockchain-Bachelor/utils"
)

func (a *AuthorityNode) GetPatientTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patientID")
	if patientID == "" {
		http.Error(w, "patientID is required", http.StatusBadRequest)
		return
	}

	// Versuch, die Standard Base64-Dekodierung zu verwenden
	decodedPatientID, err := base64.StdEncoding.DecodeString(patientID)
	if err != nil {
		// Wenn die Standard-Dekodierung fehlschlägt, versuche die URL-sichere Dekodierung
		fmt.Println("Standard Base64-Dekodierung fehlgeschlagen:", err)
		decodedPatientID, err = base64.URLEncoding.DecodeString(patientID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode patientID: %v", err), http.StatusBadRequest)
			return
		}
	}

	patientHash := base64.StdEncoding.EncodeToString(decodedPatientID)
	patientData, exists := a.Patients[patientHash]
	if !exists {
		http.Error(w, "patient not found", http.StatusNotFound)
		return
	}

	var transactions []*blockchain.Transaction
	for _, tx := range patientData.Transactions {
		transactions = append(transactions, tx)
	}

	responseData, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to serialize transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

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
	var transaction blockchain.Transaction

	// Dekodiere die Transaktionsdaten aus der Anfrage
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "failed to decode transaction data", http.StatusBadRequest)
		return
	}

  if err != nil {
    http.Error(w, fmt.Sprintf("failed to create transaction: %v", err), http.StatusInternalServerError)
    return
  }

  // Füge die validierte und signierte Transaktion zum Pool hinzu
  if err := authorityNode.AddTransaction(&transaction); err != nil {
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
	publicKey := utils.SerializePublicKey(&a.PrivateKey.PublicKey)
	response := map[string]string{
		"publicKey": base64.StdEncoding.EncodeToString(publicKey),
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
	http.HandleFunc("/getPatientTransactions", a.GetPatientTransactionsHandler)
	http.HandleFunc("/getTransactionPool", a.GetTransactionPoolHandler)
	http.HandleFunc("/sync", a.SyncHandler)
	http.HandleFunc("/getPublicKey", a.GetPublicKeyHandler)
}

func (node *Node) SetupNodeRoutes() {
	http.HandleFunc("/getBlockchain", node.GetBlockchainHandler)
}
