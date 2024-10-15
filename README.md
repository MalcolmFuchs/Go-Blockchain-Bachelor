# Go-Blockchain-Bachelor

## Übersicht

Dieses Dokument beschreibt die Implementierung und Verwendung der `Go-Blockchain-Bachelor`-Anwendung. Es enthält Anweisungen zur Installation, zum Starten der Nodes und zur Nutzung der verschiedenen HTTP-Endpunkte, die zur Interaktion mit der Blockchain verfügbar sind.

## Installation

1. **Voraussetzungen**
   - Golang 1.18 oder höher
   - Postman oder `curl` für API-Tests

2. **Klonen des Repositories**
   ```bash
   git clone https://github.com/MalcolmFuchs/Go-Blockchain-Bachelor.git
   cd Go-Blockchain-Bachelor
   ```

3. **Abhängigkeiten installieren**
   ```bash
   go mod tidy
   ```

4. **Build**
   ```bash
   go build -o Go-Blockchain-Bachelor
   ```

## Anwendung starten

### Authority Node starten

```bash
./Go-Blockchain-Bachelor node --port 8080
```

### Client Node starten und mit Authority Node verbinden

```bash
./Go-Blockchain-Bachelor node --authority localhost:8080 --port 8081
```

### Add Transaction
```bash
./Go-Blockchain-Bachelor create --node_address localhost:8080 --type "medical" --notes "Routine Check-up" --results "All tests normal" --patient ./keys/patient_public_key.pem --key ./keys/doctor_private_key.pem
```
### GetTransactionPool
```bash
curl "http://localhost:8080/getTransactionPool"
```

### CreateBlock
```bash
curl -X POST "http://localhost:8080/createBlock"
```

### GetBlockchain
```bash
curl "http://localhost:8080/getBlockchain"
```




