package blockchain

type Blockchain struct {
	Blocks []Block
}

func (bc *Blockchain) AddBlock(newBlock Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) GetLastBlock() Block {
	if len(bc.Blocks) == 0 {
		return NewBlock([]Transaction{}, "")
	}
	return bc.Blocks[len(bc.Blocks)-1]
}

// Validierung der Blockchain (z.B. durch Prüfung der Hashes)
func (bc *Blockchain) Validate() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		prevBlock := bc.Blocks[i-1]
		currentBlock := bc.Blocks[i]
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
		if currentBlock.Hash != calculateHash(currentBlock) {
			return false
		}
	}
	return true
}
