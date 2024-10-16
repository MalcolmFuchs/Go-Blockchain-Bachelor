# 🚀 **Go-Blockchain-Bachelor**

## 📋 **Übersicht**

Dieses Dokument beschreibt die Implementierung und Verwendung der `Go-Blockchain-Bachelor`-Anwendung. Es enthält Anweisungen zur **Installation**, zum **Starten der Nodes** und zur **Nutzung der verschiedenen Endpunkte**, die zur Interaktion mit der Blockchain verfügbar sind.

---

## 🛠️ **Installation**

1. **Voraussetzungen**:
   - ✅ Golang 1.18 oder höher
   - ✅ `curl` für API-Tests

2. **Klonen des Repositories**:
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

## Starten der Nodes und Testen der Endpunkte

1. **Authority Node starten**:
   ```bash
   ./Go-Blockchain-Bachelor node --port 8080
   ```

2. **Client Node starten und mit Authority Node verbinden**:
   ```bash
   ./Go-Blockchain-Bachelor node --authority localhost:8080 --port 8081
   ```

3. **Transaktion hinzufügen:**
   ```bash
   ./Go-Blockchain-Bachelor create --node_address localhost:8080 --type "medical" --notes "Routine Check-up" --results "All tests normal" --patient ./keys/patient_public_key.pem --key ./keys/doctor_private_key.pem

   ./Go-Blockchain-Bachelor create --node_address localhost:8080 --type "medical" --notes "Blood Test" --results "Cholesterol levels normal" --patient ./keys/patient_public_key.pem --key ./keys/doctor_private_key.pem

   ./Go-Blockchain-Bachelor create --node_address localhost:8080 --type "medical" --notes "X-Ray Examination" --results "No fractures detected" --patient ./keys/patient_public_key.pem --key ./keys/doctor_private_key.pem

   ./Go-Blockchain-Bachelor create --node_address localhost:8080 --type "prescription" --notes "Prescribed medication for hypertension" --results "Medication: Lisinopril 10mg daily" --patient ./keys/patient_public_key.pem --key ./keys/doctor_private_key.pem

   ./Go-Blockchain-Bachelor create --node_address localhost:8080 --type "referral" --notes "Referral to cardiologist" --results "Appointment scheduled for 2024-11-01" --patient ./keys/patient_public_key.pem --key ./keys/doctor_private_key.pem
    ```

4. **TransaktionsPool anzeigen:**
   ```bash
   curl "http://localhost:8080/getTransactionPool"
    ```

5. **Block erstellen:**
   ```bash
    curl "http://localhost:8080/createBlock"
    ```

4. **Blockchain anzeigen:**
   ```bash
   curl "http://localhost:8080/getBlockchain"
    ```

5. **Patienten Transaktionen anzeigen:**
   ```bash
      ./Go-Blockchain-Bachelor view --node_address localhost:8080 --key ./keys/patient_private_key.pem
    ```


