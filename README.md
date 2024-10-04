
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

## API-Endpunkte

### 1. `/addTransaction`

- **Beschreibung**: Fügt eine neue Transaktion zur Blockchain hinzu.
- **Methode**: `POST`
- **URL**: `http://localhost:8080/addTransaction`
- **Header**: 
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
  "type": "Blood Test",
  "notes": "Routine blood test for cholesterol levels.",
  "results": "Cholesterol: 180 mg/dL",
  "doctor": "4f8b54c26d30dff4b932e769d1a9b0e7130c1d8c3c7209ecfa61ff01e8f4f07a",
  "patient": "3c5a2bcd98216d66dbac8b86e5f8d8a9c5405a98df27e9ed1e9fc74a91d8b5e4",
  "key": "d3c06ff6d01328abf346f35b4c6a25c45f92d06c3b4f48d40a7e4e06b1d29bb6"
  }
  ```
- **Beispiel mit `curl`**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{
  "type": "Blood Test",
  "notes": "Routine blood test for cholesterol levels.",
  "results": "Cholesterol: 180 mg/dL",
  "doctor": "4f8b54c26d30dff4b932e769d1a9b0e7130c1d8c3c7209ecfa61ff01e8f4f07a",
  "patient": "3c5a2bcd98216d66dbac8b86e5f8d8a9c5405a98df27e9ed1e9fc74a91d8b5e4",
  "key": "d3c06ff6d01328abf346f35b4c6a25c45f92d06c3b4f48d40a7e4e06b1d29bb6"
  }' http://localhost:8080/addTransaction
  ```

### 2. `/getBlockchain`

- **Beschreibung**: Gibt die gesamte Blockchain als JSON zurück.
- **Methode**: `GET`
- **URL**: `http://localhost:8080/getBlockchain`
- **Beispiel mit `curl`**:
  ```bash
  curl http://localhost:8080/getBlockchain
  ```

## Starten der Nodes und Testen der Endpunkte

1. **Authority Node starten**:
   ```bash
   ./Go-Blockchain-Bachelor node --port 8080
   ```

2. **Client Node starten und mit Authority Node verbinden**:
   ```bash
   ./Go-Blockchain-Bachelor node --authority localhost:8080 --port 8081
   ```

3. **Postman-Tests durchführen**:
   - Füge die Endpunkte `/addTransaction` und `/getBlockchain` als separate Requests in Postman hinzu.
   - Stelle sicher, dass die Header und Body-Inhalte den oben beschriebenen Beispielen entsprechen.

4. **Überprüfen der Ergebnisse**:
   - Verwende `curl` oder Postman, um die verschiedenen Endpunkte zu testen.
   - Überprüfe die Serverlogs und die Blockchain-Ausgaben, um sicherzustellen, dass die Transaktionen und Blöcke korrekt hinzugefügt und abgerufen werden.
