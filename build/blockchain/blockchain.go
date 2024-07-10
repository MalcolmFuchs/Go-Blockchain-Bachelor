package blockchain

import "time"

type Blockchain struct {
	genesisBlock   Block
	chain          []Block
	authorityNodes []AuthorityNode
}

type AuthorityNode struct {
	id         string
	publicKey  string
	privateKey string
}

func CreateBlockchain(authorityNodes []AuthorityNode) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		authorityNodes,
	}
}

func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
		prevBlock := b.chain[i]
		currBlock := b.chain[i+1]

		if currBlock.hash != currBlock.Hash() || currBlock.prevHash != prevBlock.hash {
			return false
		}
	}
	return true
}

func (b *Blockchain) addBlock(block Block) {
	b.chain = append(b.chain, block)
}
